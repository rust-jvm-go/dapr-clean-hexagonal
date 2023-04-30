//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"url-shortener/repository/mongodb"
	"url-shortener/xsetup"
)

func InitializeSetup(ctx context.Context) (xsetup.Environment, error) {
	wire.Build(xsetup.NewDaprClient, xsetup.NewEnvironment, xsetup.NewInitConfig, mongodb.NewMongoRepository)
	return xsetup.Environment{}, nil
}
