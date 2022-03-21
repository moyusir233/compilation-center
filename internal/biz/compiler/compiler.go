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
	// 存放服务项目的目录的绝对路径
	projectDir string
	// 相对于数据收集项目的根目录，存放其api定义的目录的相对路径
	apiDir string
	// 相对于数据处理项目的根目录，存放其服务实现的目录的相对路径
	serviceDir string
	// 缓存的过期时间
	timeout time.Duration
	// 用于控制并发的读写锁
	rwMutex *sync.RWMutex
	// 存放编译结果的表
	table map[string]*bytes.Buffer
	// 管理协程的协程组
	eg     *errgroup.Group
	ctx    context.Context
	cancel func()
	// 用于复用bytes.Buffer的池
	pool *sync.Pool
}

func NewCompiler(dir *conf.Service_Compiler_CodeDir, timeout time.Duration) *Compiler {
	ctx, cancel := context.WithCancel(context.Background())

	return &Compiler{
		projectDir: dir.ProjectDir,
		apiDir:     dir.ApiDir,
		serviceDir: dir.ServiceDir,
		timeout:    timeout,
		rwMutex:    new(sync.RWMutex),
		table:      make(map[string]*bytes.Buffer),
		eg:         new(errgroup.Group),
		ctx:        ctx,
		cancel:     cancel,
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
		c.table[key].Reset()
		c.pool.Put(c.table[key])
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
func (c *Compiler) Compile(key string, code map[string]*bytes.Buffer) (
	exe *bytes.Buffer, err error) {
	c.rwMutex.RLock()
	// 先查询缓存
	if buffer, ok := c.table[key]; ok {
		c.rwMutex.RUnlock()
		return buffer, nil
	}

	// 由于缓存中未查询到，需要执行编译
	// 占据写锁，确保始终只有一次编译过程
	c.rwMutex.RUnlock()
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	// 获得存放编译结果的buffer
	exe = c.pool.Get().(*bytes.Buffer)
	defer func() {
		if err != nil {
			exe.Reset()
			c.pool.Put(exe)
		}
	}()

	// 执行编译
	err = c.compileTo(code, exe)
	if err != nil {
		return nil, err
	}

	// 将结果放入缓存表中
	c.table[key] = exe

	// 注册自动清理的协程
	c.eg.Go(func() error {
		c.autoClearCache(key)
		return nil
	})
	return
}

// 通过执行shell完成编译
func (c *Compiler) compileTo(code map[string]*bytes.Buffer, result *bytes.Buffer) error {
	// 为了进行编译，首先需要将生成的代码先写入到文件中
	apiDirPath := filepath.Join(c.projectDir, c.apiDir)
	serviceDirPath := filepath.Join(c.projectDir, c.serviceDir)

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

	// 创建存放编译结果的临时目录
	tempDir, err := os.MkdirTemp("", "app")
	if err != nil {
		tempDir = os.TempDir()
	}
	defer func() {
		if tempDir != os.TempDir() {
			os.RemoveAll(tempDir)
		}
	}()
	targetPath := filepath.Join(tempDir, "server")

	// 执行shell
	cmd := exec.Command("/shell/build.sh", c.projectDir, targetPath)
	err = cmd.Run()
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
