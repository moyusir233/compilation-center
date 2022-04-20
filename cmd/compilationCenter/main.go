package main

import (
	"flag"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"os"

	"gitee.com/moyusir/compilation-center/internal/conf"
	util "gitee.com/moyusir/util/logger"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "compilation-center"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server) *kratos.App {
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

func main() {
	flag.Parse()
	logger := util.NewJsonZapLoggerWarpper(Name)
	helper := log.NewHelper(logger)

	bc, err := conf.LoadConfig(flagconf, logger)
	if err != nil {
		helper.Fatalf("导入配置时发生了错误:%v", err)
	}

	app, cleanUp, err := initApp(bc.Server, bc.Service, bc.Data, logger)
	if err != nil {
		helper.Fatalf("应用初始化时发生了错误:%v", err)
	}
	defer cleanUp()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		helper.Fatalf("应用运行时发生了错误:%v", err)
	}
}
