//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"cloudwegodemo/cmd/cloudwegodemo/internal/biz"
	"cloudwegodemo/cmd/cloudwegodemo/internal/config"
	"cloudwegodemo/cmd/cloudwegodemo/internal/repo"
	"cloudwegodemo/cmd/cloudwegodemo/internal/server"

	"cloudwegodemo/pkg/configor"
	"cloudwegodemo/pkg/database"

	"github.com/google/wire"
)

// wireApp init application.
func wireApp(*configor.Option) (*APP, func(), error) {
	panic(wire.Build(config.ConfigProvider, database.DatabaseProvider, repo.RepoProvider, biz.BizProvider, server.ServerProvider, newApp))
}
