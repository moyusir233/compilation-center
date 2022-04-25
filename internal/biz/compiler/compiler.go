package compiler

import (
	"bytes"
	"gitee.com/moyusir/compilation-center/internal/conf"
	"github.com/go-kratos/kratos/v2/errors"
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
	// 保证每次编译时项目文件夹只被一个协程占据的锁
	mutex *sync.Mutex
}

type tableNode struct {
	buffer *bytes.Buffer
	ticker *time.Ticker
}

func NewCompiler(dir *conf.Service_Compiler_CodeDir) *Compiler {
	return &Compiler{
		projectDir: dir.ProjectDir,
		apiDir:     dir.ApiDir,
		serviceDir: dir.ServiceDir,
		mutex:      &sync.Mutex{},
	}
}

// Compile 以指定的key强制执行编译，若key已存在，会覆盖保存的缓存
func (c *Compiler) Compile(code map[string]*bytes.Buffer) (
	exe []byte, err error) {
	// 执行编译，利用锁确保项目目录被单独的编译程序使用
	c.mutex.Lock()
	defer c.mutex.Unlock()

	buffer := bytes.NewBuffer(make([]byte, 0, 1024))

	// 执行编译
	err = c.compileTo(code, buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
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
