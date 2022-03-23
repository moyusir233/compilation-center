package biz

import (
	"bytes"
	"gitee.com/moyusir/compilation-center/internal/biz/codegenerator"
	"gitee.com/moyusir/compilation-center/internal/biz/compiler"
	"gitee.com/moyusir/compilation-center/internal/conf"
	utilApi "gitee.com/moyusir/util/api/util/v1"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type BuildUsecase struct {
	generator  *codegenerator.CodeGenerator
	dcCompiler *compiler.Compiler
	dpCompiler *compiler.Compiler
	repo       ClientCodeRepo
	logger     *log.Helper
}

// ClientCodeRepo 保存客户端代码使用的接口
type ClientCodeRepo interface {
	// SaveClientCode 保存客户端代码，其中map的键为文件名，值为文件对应的二进制数据
	SaveClientCode(key string, files map[string]*bytes.Reader) error
}

func NewBuildUsecase(service *conf.Service, repo ClientCodeRepo, logger log.Logger) (*BuildUsecase, func(), error) {
	generator, err := codegenerator.NewCodeGenerator(
		service.CodeGenerator.DataProcessingTmplRoot, service.CodeGenerator.DataCollectionTmplRoot)
	if err != nil {
		return nil, nil, err
	}

	timeout := service.Compiler.Timeout.AsDuration()
	dcCompiler := compiler.NewCompiler(service.Compiler.DataCollection, timeout)
	dpCompiler := compiler.NewCompiler(service.Compiler.DataProcessing, timeout)

	return &BuildUsecase{
			generator:  generator,
			dcCompiler: dcCompiler,
			dpCompiler: dpCompiler,
			repo:       repo,
			logger:     log.NewHelper(logger),
		}, func() {
			// 关闭编译器负责回收bytes.Buffer的协程
			dpCompiler.Close()
			dcCompiler.Close()
		}, nil
}

// BuildDCServiceExe 编译数据收集服务的二进制程序，并将生成的客户端proto文件以zip形式保存到数据库中
func (u *BuildUsecase) BuildDCServiceExe(
	username string, states []*utilApi.DeviceStateRegisterInfo, configs []*utilApi.DeviceConfigRegisterInfo) (
	*bytes.Reader, error) {
	// 先查询是否已经编译过相关文件
	compiled, ok := u.dcCompiler.IsCompiled(username)
	if ok {
		return compiled, nil
	}

	// 未查询到，则重新生成代码并编译
	dc, err := u.generator.GetDataCollectionServiceFiles(configs, states)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	// 编译获得可执行程序
	result, err := u.dcCompiler.Compile(username, dc)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	// 保存客户端代码，以用户名为field在以client_code为键名中的hash中保存
	files := make(map[string]*bytes.Reader)
	for k, v := range dc {
		if strings.HasSuffix(k, ".proto") {
			files[k] = bytes.NewReader(v.Bytes())
		}
	}

	// 这里保存文件的错误不返回，避免影响服务容器的启动
	err = u.repo.SaveClientCode(username, files)
	if err != nil {
		u.logger.Error(err)
	}

	return result, nil
}

// BuildDPServiceExe 编译数据处理服务的二进制程序
func (u *BuildUsecase) BuildDPServiceExe(
	username string, states []*utilApi.DeviceStateRegisterInfo, configs []*utilApi.DeviceConfigRegisterInfo) (
	*bytes.Reader, error) {
	// 先查询是否已经编译过相关文件
	compiled, ok := u.dpCompiler.IsCompiled(username)
	if ok {
		return compiled, nil
	}

	// 未查询到，则重新生成代码并编译
	dp, err := u.generator.GetDataProcessingServiceFiles(configs, states)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	result, err := u.dpCompiler.Compile(username, dp)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	return result, nil
}
