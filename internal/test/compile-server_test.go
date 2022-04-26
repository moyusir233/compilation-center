package test

import (
	"bytes"
	"context"
	v1 "gitee.com/moyusir/compilation-center/api/compilationCenter/v1"
	utilApi "gitee.com/moyusir/util/api/util/v1"
	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestCompilationCenter(t *testing.T) {
	// 启动服务器，获得客户端
	client := StartCompilationCenterServer(t)

	//// 连接远程的编译中心进行测试
	//conn, err := g.DialInsecure(context.Background(),
	//	g.WithEndpoint("compilation-center.test.svc.cluster.local:9000"))
	//if err != nil {
	//	t.Fatal(err)
	//}
	//client := v1.NewBuildClient(conn)

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
					Name:        "status",
					Type:        utilApi.Type_BOOL,
					WarningRule: nil,
				},
				{
					Name:        "device_type",
					Type:        utilApi.Type_STRING,
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
					Name:        "status",
					Type:        utilApi.Type_BOOL,
					WarningRule: nil,
				},
				{
					Name:        "device_type",
					Type:        utilApi.Type_STRING,
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
	readExeStreamToWriter := func(stream grpc.ClientStream, writer io.Writer) error {
		// 读取流中的二进制数据
		reply := new(v1.BuildReply)
		for {
			err := stream.RecvMsg(reply)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					t.Error(err)
				}
				break
			}
			if len(reply.Exe) > 0 {
				_, err := writer.Write(reply.Exe)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	// 获得可执行文件
	getExes := func() ([]*bytes.Buffer, error) {
		// 调用编译的grpc函数
		var streams []grpc.ClientStream
		var buffers []*bytes.Buffer

		{
			dcStream, err := client.GetDataCollectionServiceProgram(context.Background(), &v1.BuildRequest{
				Username:                  "test",
				DeviceStateRegisterInfos:  stateInfo,
				DeviceConfigRegisterInfos: configInfo,
			})
			if err != nil {
				return nil, err
			}

			dcExe := bytes.NewBuffer(make([]byte, 0, 1024))

			streams = append(streams, dcStream)
			buffers = append(buffers, dcExe)
		}
		{
			dpStream, err := client.GetDataProcessingServiceProgram(context.Background(), &v1.BuildRequest{
				Username:                  "test",
				DeviceStateRegisterInfos:  stateInfo,
				DeviceConfigRegisterInfos: configInfo,
			})
			if err != nil {
				return nil, err
			}

			dpExe := bytes.NewBuffer(make([]byte, 0, 1024))

			streams = append(streams, dpStream)
			buffers = append(buffers, dpExe)
		}

		for i, s := range streams {
			err := readExeStreamToWriter(s, buffers[i])
			if err != nil {
				return nil, err
			}
		}

		return buffers, nil
	}

	exes, err := getExes()
	if err != nil {
		t.Fatal(err)
	}

	// 进行第二次，以确定通过缓存得到的可执行文件和原来的一致
	caches, err := getExes()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(exes); i++ {
		if !bytes.Equal(exes[i].Bytes(), caches[i].Bytes()) {
			t.Fatal("第一次编译获得的可执行文件和通过缓存获得的可执行文件不一致")
		}
	}
}

func TestTmp(t *testing.T) {
	file, err := os.Open("util.go")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	all, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(all))

	all, err = ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(all))
}
