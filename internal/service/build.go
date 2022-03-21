package service

import (
	"context"

	pb "gitee.com/moyusir/compilation-center/api/compilationCenter/v1"
)

type BuildService struct {
	pb.UnimplementedBuildServer
}

func NewBuildService() *BuildService {
	return &BuildService{}
}

func (s *BuildService) GetDataCollectionServiceProgram(ctx context.Context, req *pb.BuildRequest) (*pb.BuildReply, error) {
	return &pb.BuildReply{}, nil
}
func (s *BuildService) GetDataProcessingServiceProgram(ctx context.Context, req *pb.BuildRequest) (*pb.BuildReply, error) {
	return &pb.BuildReply{}, nil
}
