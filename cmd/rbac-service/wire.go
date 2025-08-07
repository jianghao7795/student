//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"student/internal/conf"
	"student/internal/pkg/nacos"
	"student/internal/rbac-service/biz"
	"student/internal/rbac-service/data"
	"student/internal/rbac-service/server"
	"student/internal/rbac-service/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Bootstrap, log.Logger, *nacos.Discovery) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
