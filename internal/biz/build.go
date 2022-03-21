package biz

import (
	"bytes"
	"gitee.com/moyusir/compilation-center/internal/biz/codegenerator"
	"gitee.com/moyusir/compilation-center/internal/biz/compiler"
	"gitee.com/moyusir/compilation-center/internal/conf"
	utilApi "gitee.com/moyusir/util/api/util/v1"
	"golang.org/x/sync/errgroup"
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

// BuildServiceExe
func (u *BuildUsecase) BuildServiceExe(
	username string, states []*utilApi.DeviceStateRegisterInfo, configs []*utilApi.DeviceConfigRegisterInfo) (
	dcExe, dpExe *bytes.Buffer, err error) {
	dc, dp, err := u.generator.GetServiceFiles(configs, states)
	if err != nil {
		return nil, nil, err
	}

	eg := new(errgroup.Group)
	eg.Go(func() error {
		result, err := u.dcCompiler.Compile(username, dc)
		if err != nil {
			return err
		} else {
			dcExe = result
			return nil
		}
	})
	eg.Go(func() error {
		result, err := u.dpCompiler.Compile(username, dp)
		if err != nil {
			return err
		} else {
			dpExe = result
			return nil
		}
	})
	err = eg.Wait()

	if err != nil {
		return nil, nil, err
	}
	return
}
