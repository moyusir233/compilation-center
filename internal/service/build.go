package service

import (
	"context"
	"gitee.com/moyusir/compilation-center/internal/biz"

	pb "gitee.com/moyusir/compilation-center/api/compilationCenter/v1"
)

type BuildService struct {
	pb.UnimplementedBuildServer
	uc *biz.BuildUsecase
}

func NewBuildService(uc *biz.BuildUsecase) *BuildService {
	return &BuildService{uc: uc}
}

func (s *BuildService) GetServiceProgram(ctx context.Context, req *pb.BuildRequest) (*pb.BuildReply, error) {
	dcExe, dpExe, err := s.uc.BuildServiceExe(req.Username, req.DeviceStateRegisterInfos, req.DeviceConfigRegisterInfos)
	if err != nil {
		return nil, err
	}

	return &pb.BuildReply{
		DcExe: dcExe.Bytes(),
		DpExe: dpExe.Bytes(),
	}, nil
}
