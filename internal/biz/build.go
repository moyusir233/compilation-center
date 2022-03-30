package biz

import (
	"bytes"
	"gitee.com/moyusir/compilation-center/internal/biz/codegenerator"
	"gitee.com/moyusir/compilation-center/internal/biz/compiler"
	"gitee.com/moyusir/compilation-center/internal/conf"
	utilApi "gitee.com/moyusir/util/api/util/v1"
	"github.com/go-kratos/kratos/v2/errors"
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
	// IsValid 判断账号是否已经有效，是为了避免在账号已经注销的情况下使用了无效的编译缓存
	IsValid(username string) bool
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
	// 由于账号注册失败时保存的编译缓存是无效的，因此先判断账号是否有效，有效才允许使用缓存
	// （只有保存了客户端代码的账号，视作进行了稳定的编译，才能使用缓存）
	compiled, ok := u.dcCompiler.IsCompiled(username)
	if ok && u.repo.IsValid(username) {
		return compiled, nil
	}

	// 缓存无效或未查询到，则重新生成代码并编译
	dc, err := u.generator.GetDataCollectionServiceFiles(configs, states)
	if err != nil {
		return nil, errors.Newf(
			500, "Biz_Error", "生成数据收集服务代码时发生了错误:%v", err)
	}

	// 在后台保存客户端代码，以用户名为field在以client_code为键名中的hash中保存
	files := make(map[string]*bytes.Reader)
	for k, v := range dc {
		if strings.HasSuffix(k, ".proto") {
			files[k] = bytes.NewReader(v.Bytes())
		}
	}
	go func() {
		// 这里保存文件的错误不返回，避免影响服务容器的启动
		err = u.repo.SaveClientCode(username, files)
		if err != nil {
			u.logger.Error(err)
		}
	}()

	// 编译获得可执行程序
	result, err := u.dcCompiler.Compile(username, dc)
	if err != nil {
		return nil, errors.Newf(
			500, "Biz_Error", "编译数据收集服务代码时发生了错误:%v", err)
	}

	return result, nil
}

// BuildDPServiceExe 编译数据处理服务的二进制程序
func (u *BuildUsecase) BuildDPServiceExe(
	username string, states []*utilApi.DeviceStateRegisterInfo, configs []*utilApi.DeviceConfigRegisterInfo) (
	*bytes.Reader, error) {
	// 由于账号注册失败时保存的编译缓存是无效的，因此先判断账号是否有效，有效才允许使用缓存
	// （只有保存了客户端代码的账号，视作进行了稳定的编译，才能使用缓存）
	compiled, ok := u.dpCompiler.IsCompiled(username)
	if ok && u.repo.IsValid(username) {
		return compiled, nil
	}

	// 缓存无效或未查询到，则重新生成代码并编译
	dp, err := u.generator.GetDataProcessingServiceFiles(configs, states)
	if err != nil {
		return nil, errors.Newf(
			500, "Biz_Error", "生成数据处理服务代码时发生了错误:%v", err)
	}

	result, err := u.dpCompiler.Compile(username, dp)
	if err != nil {
		return nil, errors.Newf(
			500, "Biz_Error", "编译数据处理服务代码时发生了错误:%v", err)
	}

	return result, nil
}
