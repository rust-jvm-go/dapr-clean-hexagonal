// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"url-shortener/repository/mongodb"
	"url-shortener/xsetup"
)

// Injectors from wire.go:

func InitializeSetup(ctx context.Context) (*xsetup.Environment, error) {
	daprClient, err := xsetup.NewDaprClient()
	if err != nil {
		return nil, err
	}
	initConfig, err := xsetup.NewInitConfig(ctx, daprClient)
	if err != nil {
		return nil, err
	}
	iRedirectRepository, err := mongodb.NewMongoRepository(initConfig)
	if err != nil {
		return nil, err
	}
	environment := xsetup.NewEnvironment(daprClient, initConfig, iRedirectRepository)
	return environment, nil
}
