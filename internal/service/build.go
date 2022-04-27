package service

import (
	v1 "gitee.com/moyusir/compilation-center/api/compilationCenter/v1"
	"gitee.com/moyusir/compilation-center/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
)

type BuildService struct {
	v1.UnimplementedBuildServer
	uc *biz.BuildUsecase
}

func NewBuildService(uc *biz.BuildUsecase) *BuildService {
	return &BuildService{uc: uc}
}

func (b *BuildService) GetDataCollectionServiceProgram(request *v1.BuildRequest, stream v1.Build_GetDataCollectionServiceProgramServer) error {
	exe, err := b.uc.BuildDCServiceExe(
		request.Username, request.DeviceStateRegisterInfos, request.DeviceConfigRegisterInfos)
	if err != nil {
		return err
	}
	defer exe.Close()

	reply := &v1.BuildReply{
		Exe: make([]byte, 1024),
	}
	for {
		n, err := exe.Read(reply.Exe)
		if err != nil {
			break
		}
		reply.Exe = reply.Exe[:n]
		err = stream.Send(reply)
		if err != nil {
			return errors.Newf(
				500, "Service_Error", "发送数据收集服务可执行程序时发生了错误:%v", err)
		}
	}
	return nil
}

func (b *BuildService) GetDataProcessingServiceProgram(request *v1.BuildRequest, stream v1.Build_GetDataProcessingServiceProgramServer) error {
	exe, err := b.uc.BuildDPServiceExe(
		request.Username, request.DeviceStateRegisterInfos, request.DeviceConfigRegisterInfos)
	if err != nil {
		return err
	}
	defer exe.Close()

	reply := &v1.BuildReply{
		Exe: make([]byte, 1024),
	}
	for {
		n, err := exe.Read(reply.Exe)
		if err != nil {
			break
		}
		reply.Exe = reply.Exe[:n]
		err = stream.Send(reply)
		if err != nil {
			return errors.Newf(
				500, "Service_Error", "发送数据处理服务可执行程序时发生了错误:%v", err)
		}
	}
	return nil
}
