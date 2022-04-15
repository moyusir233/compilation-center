package compiler

import (
	"bytes"
	"context"
	"gitee.com/moyusir/compilation-center/internal/conf"
	"github.com/go-kratos/kratos/v2/errors"
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
	table map[string]*tableNode
	// 管理协程的协程组
	eg     *errgroup.Group
	ctx    context.Context
	cancel func()
	// 用于复用bytes.Buffer的池
	pool *sync.Pool
}

type tableNode struct {
	buffer *bytes.Buffer
	ticker *time.Ticker
}

func NewCompiler(dir *conf.Service_Compiler_CodeDir, timeout time.Duration) *Compiler {
	ctx, cancel := context.WithCancel(context.Background())

	return &Compiler{
		projectDir: dir.ProjectDir,
		apiDir:     dir.ApiDir,
		serviceDir: dir.ServiceDir,
		timeout:    timeout,
		rwMutex:    new(sync.RWMutex),
		table:      make(map[string]*tableNode),
		eg:         new(errgroup.Group),
		ctx:        ctx,
		cancel:     cancel,
		pool: &sync.Pool{New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 1024))
		}},
	}
}
func (c *Compiler) autoClearCache(key string) {
	if node, ok := c.table[key]; ok {
		defer node.ticker.Stop()
		select {
		case <-node.ticker.C:
			// 超时，将使用的buffer放回池中，并删除掉key
			c.rwMutex.Lock()
			if node.buffer != nil {
				node.buffer.Reset()
				c.pool.Put(node.buffer)
			}
			delete(c.table, key)
			c.rwMutex.Unlock()
		case <-c.ctx.Done():
			return
		}
	}
}

// Close 关闭还处于运行的所有计时器协程
func (c *Compiler) Close() error {
	c.cancel()
	return c.eg.Wait()
}

// IsCompiled 查询指定的key是否已经编译过,若是则返回key对应的文件和true，否则返回nil和false
func (c *Compiler) IsCompiled(key string) (*bytes.Reader, bool) {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	// 先查询缓存
	if node, ok := c.table[key]; ok {
		node.ticker.Reset(c.timeout)
		return bytes.NewReader(node.buffer.Bytes()), true
	}
	return nil, false
}

// Compile 以指定的key强制执行编译，若key已存在，会覆盖保存的缓存
func (c *Compiler) Compile(key string, code map[string]*bytes.Buffer) (
	exe *bytes.Reader, err error) {
	// 执行编译，利用锁确保项目目录被单独的编译程序使用
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	// 先查询缓存表，看是否需要再从池中取得buffer
	var buffer *bytes.Buffer
	if node, ok := c.table[key]; ok {
		node.ticker.Reset(c.timeout)
		buffer = node.buffer
		buffer.Reset()
	} else {
		buffer = c.pool.Get().(*bytes.Buffer)
		c.table[key] = &tableNode{
			buffer: buffer,
			ticker: time.NewTicker(c.timeout),
		}
		// 注册自动清理的协程
		c.eg.Go(func() error {
			c.autoClearCache(key)
			return nil
		})
	}

	// 执行编译
	err = c.compileTo(code, buffer)
	if err != nil {
		return nil, err
	}

	defer func() {
		// 当程序发生错误时，直接通过调用计时器清理掉节点和buffer
		if err != nil {
			c.table[key].ticker.Reset(time.Nanosecond)
		}
	}()

	return bytes.NewReader(buffer.Bytes()), nil
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

	targetPath := filepath.Join(c.projectDir, "bin", "server")

	// 执行shell
	cmd := exec.Command("/shell/build.sh", c.projectDir, targetPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Newf(500,
			"failed to exec the shell of build",
			"%s\n%s",
			err, string(output),
		)
	}

	// 执行完毕后将编译结果写入buffer
	file, err := os.Open(targetPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(result, file)
	if err != nil {
		return err
	}

	return nil
}
