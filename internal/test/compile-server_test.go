package test

import (
	"bytes"
	"context"
	v1 "gitee.com/moyusir/compilation-center/api/compilationCenter/v1"
	utilApi "gitee.com/moyusir/util/api/util/v1"
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
			},
		},
	}

	// 调用编译的grpc函数
	replyStream, err := client.GetServiceProgram(context.Background(), &v1.BuildRequest{
		Username:                  "test",
		DeviceStateRegisterInfos:  stateInfo,
		DeviceConfigRegisterInfos: configInfo,
	})
	if err != nil {
		t.Fatal(err)
	}

	// 读取流中的二进制数据
	dcExe := bytes.NewBuffer(make([]byte, 0, 1024))
	dpExe := bytes.NewBuffer(make([]byte, 0, 1024))

	for {
		reply, err := replyStream.Recv()
		if err != nil {
			break
		}
		if len(reply.DcExe) > 0 {
			dcExe.Write(reply.DcExe)
		}
		if len(reply.DpExe) > 0 {
			dpExe.Write(reply.DpExe)
		}
	}

	// 将二进制数据保存
	dc, err := os.Create("/app/dc")
	if err != nil {
		t.Fatal(err)
	}
	defer dc.Close()
	_, err = dcExe.WriteTo(dc)
	if err != nil {
		t.Fatal(err)
	}

	dp, err := os.Create("/app/dp")
	if err != nil {
		t.Fatal(err)
	}
	defer dp.Close()
	_, err = dpExe.WriteTo(dp)
	if err != nil {
		t.Fatal(err)
	}

}
