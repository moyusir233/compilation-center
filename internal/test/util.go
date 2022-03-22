package test

import (
	"context"
	v1 "gitee.com/moyusir/compilation-center/api/compilationCenter/v1"
	"gitee.com/moyusir/compilation-center/internal/conf"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"os"
	"testing"
	"time"
)

func newApp(logger log.Logger, gs *grpc.Server) *kratos.App {
	// go build -ldflags "-X main.Version=x.y.z"
	var (
		// Name is the name of the compiled software.
		Name string
		// Version is the version of the compiled software.
		Version string
		id, _   = os.Hostname()
	)

	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
		),
	)
}

func StartCompilationCenterServer(t *testing.T) v1.BuildClient {
	bootstrap, err := conf.LoadConfig("../../configs/config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	app, err := initApp(bootstrap.Server, bootstrap.Service, log.DefaultLogger)
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		err := app.Run()
		if err != nil {
			return
		}
	}()
	t.Cleanup(func() {
		app.Stop()
		<-done
	})

	for {
		select {
		case <-done:
			t.Fatal("failed to run server")
		default:
			conn, err := grpc.DialInsecure(context.Background(),
				grpc.WithEndpoint("localhost:9000"),
				grpc.WithTimeout(time.Hour),
			)
			if err != nil {
				continue
			}
			t.Cleanup(func() {
				conn.Close()
			})
			return v1.NewBuildClient(conn)
		}
	}
}
