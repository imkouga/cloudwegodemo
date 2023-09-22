// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"cloudwegodemo/cmd/cloudwegodemo/internal/biz"
	"cloudwegodemo/cmd/cloudwegodemo/internal/config"
	"cloudwegodemo/cmd/cloudwegodemo/internal/repo"
	"cloudwegodemo/internal/server"
	"cloudwegodemo/pkg"
	"cloudwegodemo/pkg/bootstrap"
	"cloudwegodemo/pkg/configor"
	"cloudwegodemo/pkg/contrib/registry"
	"cloudwegodemo/pkg/database/mysql"
	"cloudwegodemo/pkg/database/redis"
	"github.com/google/wire"
)

// Injectors from wire.go:

// wireApp init application.
func wireApp(serviceName pkg.ServiceName, option *configor.Option) (*bootstrap.APP, func(), error) {
	configorConfigor, err := config.NewGlobalConfigor(option)
	if err != nil {
		return nil, nil, err
	}
	getHTTPServerOption := config.GetHTTPServerOptionFn()
	getRegistryOption := config.GetRegistryOptionFn()
	registryRegistry, err := registry.NewRegistry(serviceName, configorConfigor, getRegistryOption)
	if err != nil {
		return nil, nil, err
	}
	hertz, err := server.NewHTTPServer(getHTTPServerOption, registryRegistry)
	if err != nil {
		return nil, nil, err
	}
	getOptionFn := config.GetMySQLOptionFn()
	mySQL, err := mysql.NewMySQLPool(configorConfigor, getOptionFn)
	if err != nil {
		return nil, nil, err
	}
	redisGetOptionFn := config.GetRedisOptionFn()
	redisRedis, err := redis.NewRedisPool(configorConfigor, redisGetOptionFn)
	if err != nil {
		return nil, nil, err
	}
	baseRepo, err := repo.NewBaseRepo(mySQL, redisRedis)
	if err != nil {
		return nil, nil, err
	}
	baseBiz, err := biz.NewBaseBiz(baseRepo)
	if err != nil {
		return nil, nil, err
	}
	app, cleanup, err := newApp(configorConfigor, hertz, registryRegistry, baseBiz)
	if err != nil {
		return nil, nil, err
	}
	return app, func() {
		cleanup()
	}, nil
}

// wire.go:

var (
	commonProvder = wire.NewSet(registry.NewRegistry)
)
