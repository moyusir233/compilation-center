//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"gitee.com/moyusir/compilation-center/internal/biz"
	"gitee.com/moyusir/compilation-center/internal/conf"
	"gitee.com/moyusir/compilation-center/internal/server"
	"gitee.com/moyusir/compilation-center/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Service, log.Logger) (*kratos.App, error) {
	panic(wire.Build(server.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
