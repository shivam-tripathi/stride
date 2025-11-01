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

// RepositorySet is a Wire provider set for repositories
var RepositorySet = wire.NewSet(
	provideUserRepository,
)

// ServiceSet is a Wire provider set for services
var ServiceSet = wire.NewSet(
	service.NewAppService,
	service.NewUserService,
)

// provideUserRepository provides a UserRepository
func provideUserRepository(db resources.DBResource) repository.UserRepository {
	return repository.NewUserRepository(db)
}

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

		// Repositories
		RepositorySet,

		// Services
		ServiceSet,

		// API Handlers
		api.NewHandler,

		// App
		app.NewApp,
	)
	return &app.App{}, nil
}

// InitializeAppWithResources wires up the dependencies with pre-initialized resources
// This is used when resources are initialized before Wire creates the app
func InitializeAppWithResources(cfg *config.Config, res *resources.Resources) (*app.App, error) {
	wire.Build(
		// Repositories - use the provided resources
		provideUserRepositoryFromResources,

		// Services
		ServiceSet,

		// API Handlers
		api.NewHandler,

		// App
		app.NewApp,
	)
	return &app.App{}, nil
}

// provideUserRepositoryFromResources creates a user repository from pre-initialized resources
func provideUserRepositoryFromResources(res *resources.Resources) repository.UserRepository {
	return repository.NewUserRepository(res.DB)
}
