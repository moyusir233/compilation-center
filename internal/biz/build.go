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
	"time"
)

type BuildUsecase struct {
	generator  *codegenerator.CodeGenerator
	dcCompiler *compiler.Compiler
	dpCompiler *compiler.Compiler
	repo       ClientCodeRepo
	logger     *log.Helper
	// 缓存的过期时间
	expire time.Duration
}

// ClientCodeRepo 保存客户端代码使用的接口
type ClientCodeRepo interface {
	// SaveClientCode 保存客户端代码，其中map的键为文件名，值为文件对应的二进制数据
	SaveClientCode(key string, files map[string]*bytes.Reader) error
	// IsValid 判断账号是否已经有效，是为了避免在账号已经注销的情况下使用了无效的编译缓存
	IsValid(username string) bool
	// SaveExe 在redis中缓存编译得到的可执行文件
	SaveExe(key string, content []byte, expire time.Duration) error
	// GetExe 获得在redis缓存中保存的可执行文件
	GetExe(key string) ([]byte, error)
}

func NewBuildUsecase(service *conf.Service, repo ClientCodeRepo, logger log.Logger) (*BuildUsecase, error) {
	generator, err := codegenerator.NewCodeGenerator(
		service.CodeGenerator.DataProcessingTmplRoot, service.CodeGenerator.DataCollectionTmplRoot)
	if err != nil {
		return nil, err
	}

	timeout := service.Compiler.Timeout.AsDuration()
	dcCompiler := compiler.NewCompiler(service.Compiler.DataCollection)
	dpCompiler := compiler.NewCompiler(service.Compiler.DataProcessing)

	return &BuildUsecase{
		generator:  generator,
		dcCompiler: dcCompiler,
		dpCompiler: dpCompiler,
		repo:       repo,
		logger:     log.NewHelper(logger),
		expire:     timeout,
	}, nil
}

// BuildDCServiceExe 编译数据收集服务的二进制程序，并将生成的客户端proto文件以zip形式保存到数据库中
func (u *BuildUsecase) BuildDCServiceExe(
	username string, states []*utilApi.DeviceStateRegisterInfo, configs []*utilApi.DeviceConfigRegisterInfo) (
	*bytes.Reader, error) {
	// 由于账号注册失败时保存的编译缓存是无效的，因此先判断账号是否有效，有效才允许使用缓存
	// （只有保存了客户端代码的账号，视作进行了稳定的编译，才能使用缓存）
	key := username + "-dc"
	exe, err := u.repo.GetExe(key)
	if err == nil && u.repo.IsValid(username) {
		return bytes.NewReader(exe), nil
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
	result, err := u.dcCompiler.Compile(dc)
	if err != nil {
		return nil, errors.Newf(
			500, "Biz_Error", "编译数据收集服务代码时发生了错误:%v", err)
	}

	// 将可执行程序保存到redis缓存中
	go func() {
		err = u.repo.SaveExe(key, result, u.expire)
		if err != nil {
			u.logger.Error(err)
		}
	}()

	return bytes.NewReader(result), nil
}

// BuildDPServiceExe 编译数据处理服务的二进制程序
func (u *BuildUsecase) BuildDPServiceExe(
	username string, states []*utilApi.DeviceStateRegisterInfo, configs []*utilApi.DeviceConfigRegisterInfo) (
	*bytes.Reader, error) {
	// 由于账号注册失败时保存的编译缓存是无效的，因此先判断账号是否有效，有效才允许使用缓存
	// （只有保存了客户端代码的账号，视作进行了稳定的编译，才能使用缓存）
	key := username + "-dp"
	exe, err := u.repo.GetExe(key)
	if err == nil && u.repo.IsValid(username) {
		return bytes.NewReader(exe), nil
	}

	// 缓存无效或未查询到，则重新生成代码并编译
	dp, err := u.generator.GetDataProcessingServiceFiles(configs, states)
	if err != nil {
		return nil, errors.Newf(
			500, "Biz_Error", "生成数据处理服务代码时发生了错误:%v", err)
	}

	// 编译获得可执行程序
	result, err := u.dcCompiler.Compile(dp)
	if err != nil {
		return nil, errors.Newf(
			500, "Biz_Error", "编译数据收集服务代码时发生了错误:%v", err)
	}

	// 将可执行程序保存到redis缓存中
	go func() {
		err = u.repo.SaveExe(key, result, u.expire)
		if err != nil {
			u.logger.Error(err)
		}
	}()

	return bytes.NewReader(result), nil
}
