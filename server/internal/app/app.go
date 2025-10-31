package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"quizizz.com/internal/api"
	"quizizz.com/internal/config"
	"quizizz.com/internal/logger"
	"quizizz.com/internal/resources"
	"quizizz.com/pkg/middleware"
	"quizizz.com/pkg/otel"
)

// App represents the application
type App struct {
	router         *gin.Engine
	config         *config.Config
	server         *http.Server
	resources      *resources.Resources
	tracerProvider *sdktrace.TracerProvider
}

// NewApp creates a new App
func NewApp(config *config.Config, handler *api.Handler, resources *resources.Resources) *App {
	// Initialize logger
	logger.Init(config.Env)

	// Set Gin mode based on environment
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create a new Gin engine without default middleware
	router := gin.New()

	// Add custom middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	// Add OpenTelemetry middleware if enabled
	if config.OTEL.Enabled {
		router.Use(middleware.OTEL(config.OTEL.ServiceName))
	}

	// Register routes
	handler.RegisterRoutes(router)

	// Configure HTTP server
	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		router:    router,
		config:    config,
		server:    server,
		resources: resources,
	}
}

// Run starts the application
func (a *App) Run() error {
	ctx := context.Background()

	// Initialize OpenTelemetry
	if a.config.OTEL.Enabled {
		logger.Info("Initializing OpenTelemetry")
		tracerProvider, err := otel.InitTracer(ctx, a.config)
		if err != nil {
			return fmt.Errorf("failed to initialize OpenTelemetry: %w", err)
		}
		a.tracerProvider = tracerProvider
	}

	// Initialize resources before starting the server
	if err := resources.InitResources(ctx, a.resources); err != nil {
		return fmt.Errorf("failed to initialize resources: %w", err)
	}

	// Log startup
	logger.Info("Starting server",
		zap.String("port", a.config.Port),
		zap.String("env", a.config.Env),
	)

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		logger.Info("Server is listening", zap.String("port", a.config.Port))
		serverErrors <- a.server.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown or server errors.
	select {
	case err := <-serverErrors:
		logger.Error("Server error", zap.Error(err))
		return fmt.Errorf("error: starting server: %w", err)

	case sig := <-shutdown:
		logger.Info("Server is shutting down", zap.String("signal", sig.String()))

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Close all resources
		resources.CloseResources(ctx, a.resources)

		// Shutdown tracing
		if a.tracerProvider != nil {
			if err := otel.Shutdown(ctx); err != nil {
				logger.Error("Error shutting down tracer provider", zap.Error(err))
			}
		}

		// Asking listener to shut down and shed load.
		if err := a.server.Shutdown(ctx); err != nil {
			logger.Error("Could not stop server gracefully", zap.Error(err))
			a.server.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	// Flush any buffered log entries before exit
	logger.Sync()

	return nil
}
