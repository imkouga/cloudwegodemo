//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"cloudwegodemo/cmd/cloudwegodemo/internal/biz"
	"cloudwegodemo/cmd/cloudwegodemo/internal/config"
	"cloudwegodemo/cmd/cloudwegodemo/internal/repo"
	"cloudwegodemo/pkg"

	"cloudwegodemo/internal/server"

	"cloudwegodemo/pkg/bootstrap"
	"cloudwegodemo/pkg/configor"
	"cloudwegodemo/pkg/contrib/registry"
	"cloudwegodemo/pkg/database"

	"github.com/google/wire"
)

var (
	commonProvder = wire.NewSet(registry.NewRegistry)
)

// wireApp init application.
func wireApp(pkg.ServiceName, *configor.Option) (*bootstrap.APP, func(), error) {
	panic(wire.Build(commonProvder, config.ConfigProvider, database.DatabaseProvider, repo.RepoProvider, biz.BizProvider, server.ServerProvider, newApp))
}
