//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"quizizz.com/internal/api"
	"quizizz.com/internal/app"
	"quizizz.com/internal/config"
	"quizizz.com/internal/repository"
	"quizizz.com/internal/resources"
	"quizizz.com/internal/service"
)

// ResourcesSet is a Wire provider set for resources
var ResourcesSet = wire.NewSet(
	resources.NewDB,
	resources.NewRedis,
	provideResources,
)

// ServiceSet is a Wire provider set for services
var ServiceSet = wire.NewSet(
	service.NewAppService,
	service.NewUserService,
	repository.NewMockUserRepository,
)

// provideResources provides a resources.Resources struct with all resources
func provideResources(db resources.DBResource, redis resources.RedisResource) *resources.Resources {
	return &resources.Resources{
		DB:    db,
		Redis: redis,
	}
}

// InitializeApp wires up the dependencies and returns an App
func InitializeApp() (*app.App, error) {
	wire.Build(
		// Configuration
		config.NewConfig,

		// Resources
		ResourcesSet,

		// Services
		ServiceSet,

		// API Handlers
		api.NewHandler,

		// App
		app.NewApp,
	)
	return &app.App{}, nil
}
