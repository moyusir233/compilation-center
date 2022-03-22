package service

import (
	v1 "gitee.com/moyusir/compilation-center/api/compilationCenter/v1"
	"gitee.com/moyusir/compilation-center/internal/biz"
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
			return err
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
			return err
		}
	}
	return nil
}
