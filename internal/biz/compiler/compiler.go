package compiler

import (
	"bytes"
	"context"
	"gitee.com/moyusir/compilation-center/internal/conf"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Compiler 负责编译服务代码，具有缓存功能
type Compiler struct {
	// 存放数据收集项目与数据处理项目的目录的绝对路径
	dcProjectDir, dpProjectDir string
	// 相对于数据收集项目的根目录，存放其api定义和服务实现的目录的相对路径
	dcApiDir, dcServiceDir string
	// 相对于数据处理项目的根目录，存放其api定义和服务实现的目录的相对路径
	dpApiDir, dpServiceDir string
	// 缓存的过期时间
	timeout time.Duration
	// 用于控制并发的读写锁
	rwMutex *sync.RWMutex
	// 存放编译结果的表，优先存放dc的编译结果，然后存放dp的编译结果
	table map[string][]*bytes.Buffer
	// 管理协程的协程组
	eg     *errgroup.Group
	ctx    context.Context
	cancel func()
	// 用于复用bytes.Buffer的池
	pool *sync.Pool
}

func NewCompiler(compiler *conf.Service_Compiler) *Compiler {
	ctx, cancel := context.WithCancel(context.Background())

	return &Compiler{
		dcProjectDir: compiler.DataCollection.ProjectDir,
		dpProjectDir: compiler.DataProcessing.ProjectDir,
		dcApiDir:     compiler.DataCollection.ApiDir,
		dcServiceDir: compiler.DataCollection.ServiceDir,
		dpApiDir:     compiler.DataProcessing.ApiDir,
		dpServiceDir: compiler.DataProcessing.ServiceDir,
		timeout:      compiler.Timeout.AsDuration(),
		rwMutex:      new(sync.RWMutex),
		table:        make(map[string][]*bytes.Buffer),
		eg:           new(errgroup.Group),
		ctx:          ctx,
		cancel:       cancel,
		pool: &sync.Pool{New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 1024))
		}},
	}
}
func (c *Compiler) autoClearCache(key string) {
	select {
	case <-time.After(c.timeout):
		// 超时，将使用的buffer放回池中，并删除掉key
		c.rwMutex.Lock()
		for _, buffer := range c.table[key] {
			if buffer != nil {
				buffer.Reset()
				c.pool.Put(buffer)
			}
		}
		delete(c.table, key)
		c.rwMutex.Unlock()
	case <-c.ctx.Done():
		return
	}
}

// Close 关闭还处于运行的所有计时器协程
func (c *Compiler) Close() error {
	c.cancel()
	return c.eg.Wait()
}

// Compile 以指定的key执行编译
func (c *Compiler) Compile(key string, dcCode, dpCode map[string]*bytes.Buffer) (
	dc *bytes.Buffer, dp *bytes.Buffer, err error) {
	c.rwMutex.RLock()
	// 先查询缓存
	if buffers, ok := c.table[key]; ok {
		c.rwMutex.RUnlock()
		return buffers[0], buffers[1], nil
	}

	// 由于缓存中未查询到，需要执行编译
	// 占据写锁，确保始终只有一次编译过程
	c.rwMutex.RUnlock()
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	// 获得存放编译结果的buffer
	dc = c.pool.Get().(*bytes.Buffer)
	dp = c.pool.Get().(*bytes.Buffer)
	defer func() {
		if err != nil {
			dc.Reset()
			dp.Reset()
			c.pool.Put(dc)
			c.pool.Put(dp)
		}
	}()

	// 并行执行编译
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return compileTo(c.dcProjectDir, c.dcApiDir, c.dcServiceDir, dcCode, dc)
	})
	eg.Go(func() error {
		return compileTo(c.dpProjectDir, c.dpApiDir, c.dpServiceDir, dpCode, dp)
	})
	err = eg.Wait()
	if err != nil {
		return nil, nil, err
	}

	// 将结果放入缓存表中
	c.table[key] = []*bytes.Buffer{dc, dp}

	// 注册自动清理的协程
	c.eg.Go(func() error {
		c.autoClearCache(key)
		return nil
	})
	return
}

// 通过执行shell完成编译
func compileTo(projectDir, apiDir, serviceDir string, code map[string]*bytes.Buffer, result *bytes.Buffer) error {
	// 为了进行编译，首先需要将生成的代码先写入到文件中
	apiDirPath := filepath.Join(projectDir, apiDir)
	serviceDirPath := filepath.Join(projectDir, serviceDir)

	for k, v := range code {
		var path string
		if strings.HasSuffix(k, ".go") {
			path = filepath.Join(serviceDirPath, k)
		} else if strings.HasSuffix(k, ".proto") {
			path = filepath.Join(apiDirPath, k)
		}
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		_, err = v.WriteTo(file)
		if err != nil {
			return err
		}
		err = file.Close()
		if err != nil {
			return err
		}
	}

	// 执行shell
	targetPath := filepath.Join(os.TempDir(), "server")
	cmd := exec.Command("/shell/build.sh", projectDir, targetPath)
	err := cmd.Run()
	if err != nil {
		return err
	}

	// 执行完毕后将编译结果写入buffer
	file, err := os.Open(targetPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(result, file)
	if err != nil {
		return err
	}

	return nil
}
