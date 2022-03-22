package biz

import (
	"bytes"
	"gitee.com/moyusir/compilation-center/internal/biz/codegenerator"
	"gitee.com/moyusir/compilation-center/internal/biz/compiler"
	"gitee.com/moyusir/compilation-center/internal/conf"
	utilApi "gitee.com/moyusir/util/api/util/v1"
)

type BuildUsecase struct {
	generator  *codegenerator.CodeGenerator
	dcCompiler *compiler.Compiler
	dpCompiler *compiler.Compiler
}

func NewBuildUsecase(service *conf.Service) (*BuildUsecase, error) {
	generator, err := codegenerator.NewCodeGenerator(
		service.CodeGenerator.DataProcessingTmplRoot, service.CodeGenerator.DataCollectionTmplRoot)
	if err != nil {
		return nil, err
	}

	timeout := service.Compiler.Timeout.AsDuration()
	dcCompiler := compiler.NewCompiler(service.Compiler.DataCollection, timeout)
	dpCompiler := compiler.NewCompiler(service.Compiler.DataProcessing, timeout)

	return &BuildUsecase{
		generator:  generator,
		dcCompiler: dcCompiler,
		dpCompiler: dpCompiler,
	}, nil
}

// BuildDCServiceExe 编译数据收集服务的二进制程序
func (u *BuildUsecase) BuildDCServiceExe(
	username string, states []*utilApi.DeviceStateRegisterInfo, configs []*utilApi.DeviceConfigRegisterInfo) (
	*bytes.Buffer, error) {
	dc, err := u.generator.GetDataCollectionServiceFiles(configs, states)
	if err != nil {
		return nil, err
	}

	result, err := u.dcCompiler.Compile(username, dc)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// BuildDPServiceExe 编译数据处理服务的二进制程序
func (u *BuildUsecase) BuildDPServiceExe(
	username string, states []*utilApi.DeviceStateRegisterInfo, configs []*utilApi.DeviceConfigRegisterInfo) (
	*bytes.Buffer, error) {
	dp, err := u.generator.GetDataProcessingServiceFiles(configs, states)
	if err != nil {
		return nil, err
	}

	result, err := u.dpCompiler.Compile(username, dp)
	if err != nil {
		return nil, err
	}

	return result, nil
}
