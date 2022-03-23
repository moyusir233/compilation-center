// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"gitee.com/moyusir/compilation-center/internal/biz"
	"gitee.com/moyusir/compilation-center/internal/conf"
	"gitee.com/moyusir/compilation-center/internal/data"
	"gitee.com/moyusir/compilation-center/internal/server"
	"gitee.com/moyusir/compilation-center/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(confServer *conf.Server, confService *conf.Service, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	clientCodeRepo := data.NewRedisRepo(dataData, logger)
	buildUsecase, cleanup2, err := biz.NewBuildUsecase(confService, clientCodeRepo, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	buildService := service.NewBuildService(buildUsecase)
	grpcServer := server.NewGRPCServer(confServer, buildService, logger)
	app := newApp(logger, grpcServer)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}
