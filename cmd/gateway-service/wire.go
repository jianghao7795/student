//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"student/internal/conf"
	"student/internal/gateway-service/server"
	"student/internal/pkg/nacos"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Bootstrap, log.Logger, *nacos.Discovery) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, newApp))
}
