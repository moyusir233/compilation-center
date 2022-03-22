package service

import (
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

func (s *BuildService) GetServiceProgram(req *pb.BuildRequest, stream pb.Build_GetServiceProgramServer) error {
	dc, dp, err := s.uc.BuildServiceExe(req.Username, req.DeviceStateRegisterInfos, req.DeviceConfigRegisterInfos)
	if err != nil {
		return err
	}

	dcEnd, dpEnd := false, false
	reply := &pb.BuildReply{
		DcExe: make([]byte, 1024),
		DpExe: make([]byte, 1024),
	}

	for !dcEnd || !dpEnd {
		if !dcEnd {
			length, err := dc.Read(reply.DcExe)
			if err != nil {
				dcEnd = true
			} else {
				reply.DcExe = reply.DcExe[:length]
			}
		}
		if !dpEnd {
			length, err := dp.Read(reply.DpExe)
			if err != nil {
				dpEnd = true
			} else {
				reply.DpExe = reply.DpExe[:length]
			}
		}
		err := stream.Send(reply)
		if err != nil {
			return err
		}
	}
	return nil
}
