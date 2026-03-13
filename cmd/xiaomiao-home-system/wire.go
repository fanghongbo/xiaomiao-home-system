//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"
	"xiaomiao-home-system/internal/data"
	"xiaomiao-home-system/internal/server"
	"xiaomiao-home-system/internal/service"
	"xiaomiao-home-system/internal/task"

	register "github.com/go-kratos/kratos/contrib/registry/nacos/v2"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Config, *conf.Static, *register.Registry, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, task.ProviderSet, newApp))
}
