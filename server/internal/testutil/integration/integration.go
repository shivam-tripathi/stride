// Package integration provides utilities for integration testing
package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"quizizz.com/internal/api"
	"quizizz.com/internal/config"
	"quizizz.com/internal/logger"
	"quizizz.com/internal/repository"
	"quizizz.com/internal/resources"
	"quizizz.com/internal/service"
	"quizizz.com/pkg/middleware"
)

// TestEnv holds the test environment for integration tests
type TestEnv struct {
	Router      *gin.Engine
	Config      *config.Config
	Resources   *resources.Resources
	AppService  service.AppService
	UserService service.UserService
	UserRepo    repository.UserRepository
	Cleanup     func()
}

// Setup sets up the test environment for integration tests
func Setup(t *testing.T) *TestEnv {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize logger
	logger.Init("test")

	// Load configuration
	cfg := loadTestConfig(t)

	// Create test resources
	res := setupTestResources(t, cfg)

	// Create repositories
	// Initialize user repository
	// For integration tests, use MongoDB repository with test database
	// userRepo := repository.NewMongoUserRepository(resources.DB)
	// For now, use mock repository
	userRepo := repository.NewMockUserRepository()

	// Create services
	appService := service.NewAppService(cfg)
	userService := service.NewUserService(userRepo)

	apiHandler := api.NewHandler(appService, userService)

	// Create router
	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	// Register routes
	apiHandler.RegisterRoutes(router)

	// Return test environment
	return &TestEnv{
		Router:      router,
		Config:      cfg,
		Resources:   res,
		AppService:  appService,
		UserService: userService,
		UserRepo:    userRepo,
		Cleanup: func() {
			closeTestResources(t, res)
		},
	}
}

// loadTestConfig loads the test configuration
func loadTestConfig(t *testing.T) *config.Config {
	// Set test environment variables if needed
	os.Setenv("ENV", "test")
	os.Setenv("PORT", "8081")

	// Load configuration
	cfg := config.NewConfig()
	require.NotNil(t, cfg, "Failed to load test configuration")

	return cfg
}

// setupTestResources sets up test resources
func setupTestResources(t *testing.T, cfg *config.Config) *resources.Resources {
	// Create mock resources for integration tests to avoid external dependencies
	db := resources.NewMockDB(cfg)
	redis := resources.NewMockRedis(cfg)

	res := &resources.Resources{
		DB:    db,
		Redis: redis,
	}

	// Initialize resources
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := resources.InitResources(ctx, res)
	require.NoError(t, err, "Failed to initialize test resources")

	return res
}

// closeTestResources closes test resources
func closeTestResources(t *testing.T, res *resources.Resources) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resources.CloseResources(ctx, res)
}
