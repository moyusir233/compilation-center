package test

import (
	"context"
	v1 "gitee.com/moyusir/compilation-center/api/compilationCenter/v1"
	utilApi "gitee.com/moyusir/util/api/util/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"os"
	"testing"
	"time"
)

func TestCompilationCenter(t *testing.T) {
	// 启动服务器，获得客户端
	client := StartCompilationCenterServer(t)

	// 定义注册信息
	configInfo := []*utilApi.DeviceConfigRegisterInfo{
		{
			Fields: []*utilApi.DeviceConfigRegisterInfo_Field{
				{
					Name: "id",
					Type: utilApi.Type_STRING,
				},
				{
					Name: "status",
					Type: utilApi.Type_BOOL,
				},
			},
		},
		{
			Fields: []*utilApi.DeviceConfigRegisterInfo_Field{
				{
					Name: "id",
					Type: utilApi.Type_STRING,
				},
				{
					Name: "status",
					Type: utilApi.Type_BOOL,
				},
			},
		},
	}
	stateInfo := []*utilApi.DeviceStateRegisterInfo{
		{
			Fields: []*utilApi.DeviceStateRegisterInfo_Field{
				{
					Name:        "id",
					Type:        utilApi.Type_STRING,
					WarningRule: nil,
				},
				{
					Name:        "time",
					Type:        utilApi.Type_TIMESTAMP,
					WarningRule: nil,
				},
				{
					Name: "current",
					Type: utilApi.Type_DOUBLE,
					WarningRule: &utilApi.DeviceStateRegisterInfo_WarningRule{
						CmpRule: &utilApi.DeviceStateRegisterInfo_CmpRule{
							Cmp: utilApi.DeviceStateRegisterInfo_GT,
							Arg: "1000",
						},
						AggregationOperation: utilApi.DeviceStateRegisterInfo_MIN,
						Duration:             durationpb.New(time.Minute),
					},
				},
				{
					Name: "voltage",
					Type: utilApi.Type_INT64,
					WarningRule: &utilApi.DeviceStateRegisterInfo_WarningRule{
						CmpRule: &utilApi.DeviceStateRegisterInfo_CmpRule{
							Cmp: utilApi.DeviceStateRegisterInfo_GT,
							Arg: "1000",
						},
						AggregationOperation: utilApi.DeviceStateRegisterInfo_MIN,
						Duration:             durationpb.New(time.Minute),
					},
				},
			},
		},
		{
			Fields: []*utilApi.DeviceStateRegisterInfo_Field{
				{
					Name:        "id",
					Type:        utilApi.Type_STRING,
					WarningRule: nil,
				},
				{
					Name:        "time",
					Type:        utilApi.Type_TIMESTAMP,
					WarningRule: nil,
				},
				{
					Name: "current",
					Type: utilApi.Type_DOUBLE,
					WarningRule: &utilApi.DeviceStateRegisterInfo_WarningRule{
						CmpRule: &utilApi.DeviceStateRegisterInfo_CmpRule{
							Cmp: utilApi.DeviceStateRegisterInfo_GT,
							Arg: "1000",
						},
						AggregationOperation: utilApi.DeviceStateRegisterInfo_MIN,
						Duration:             durationpb.New(time.Minute),
					},
				},
				{
					Name: "voltage",
					Type: utilApi.Type_INT64,
					WarningRule: &utilApi.DeviceStateRegisterInfo_WarningRule{
						CmpRule: &utilApi.DeviceStateRegisterInfo_CmpRule{
							Cmp: utilApi.DeviceStateRegisterInfo_GT,
							Arg: "1000",
						},
						AggregationOperation: utilApi.DeviceStateRegisterInfo_MIN,
						Duration:             durationpb.New(time.Minute),
					},
				},
			},
		},
	}

	// 读取流的辅助函数
	readExeStreamToFile := func(stream grpc.ClientStream, file *os.File) error {
		defer file.Close()
		// 读取流中的二进制数据
		reply := new(v1.BuildReply)
		for {
			err := stream.RecvMsg(reply)
			if err != nil {
				t.Log(err)
				break
			}
			if len(reply.Exe) > 0 {
				_, err := file.Write(reply.Exe)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	// 调用编译的grpc函数
	var streams []grpc.ClientStream
	var files []*os.File

	{
		dpStream, err := client.GetDataProcessingServiceProgram(context.Background(), &v1.BuildRequest{
			Username:                  "test",
			DeviceStateRegisterInfos:  stateInfo,
			DeviceConfigRegisterInfos: configInfo,
		})
		if err != nil {
			t.Fatal(err)
		}
		dpExe, err := os.Create("/app/dp")
		if err != nil {
			t.Fatal(err)
		}

		streams = append(streams, dpStream)
		files = append(files, dpExe)
	}
	{
		dcStream, err := client.GetDataCollectionServiceProgram(context.Background(), &v1.BuildRequest{
			Username:                  "test",
			DeviceStateRegisterInfos:  stateInfo,
			DeviceConfigRegisterInfos: configInfo,
		})
		if err != nil {
			t.Fatal(err)
		}

		dcExe, err := os.Create("/app/dc")
		if err != nil {
			t.Fatal(err)
		}

		streams = append(streams, dcStream)
		files = append(files, dcExe)
	}

	for i, s := range streams {
		err := readExeStreamToFile(s, files[i])
		if err != nil {
			t.Error(err)
		}
	}
}
