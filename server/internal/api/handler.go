// Package api provides the API layer for the application
package api

import (
	"github.com/gin-gonic/gin"
	"quizizz.com/internal/api/handlers"
	"quizizz.com/internal/api/handlers/health"
	"quizizz.com/internal/api/handlers/ping"
	"quizizz.com/internal/api/handlers/user"
	"quizizz.com/internal/api/routes"
	"quizizz.com/internal/service"
)

// Version represents the API version
const Version = "1.0.0"

// Handler is the main API handler
type Handler struct {
	api *routes.API
}

func (h *Handler) API() *routes.API {
	return h.api
}

// NewHandler creates a new Handler
func NewHandler(appService service.AppService, userService service.UserService) *Handler {
	// Create base handler with common dependencies
	baseHandler := handlers.NewBaseHandler(appService)

	// Create specific handlers
	healthHandler := health.NewHandler(baseHandler, Version)
	pingHandler := ping.NewHandler(baseHandler)
	userHandler := user.NewHandler(baseHandler, userService)

	// Create API routes
	api := routes.NewAPI(
		baseHandler,
		healthHandler,
		pingHandler,
		userHandler,
	)

	return &Handler{
		api: api,
	}
}

// RegisterRoutes registers all API routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Register all routes from the API
	h.api.RegisterRoutes(router)
}
