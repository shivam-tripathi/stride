// Package routes defines the application routes
package routes

import (
	"github.com/gin-gonic/gin"
	"quizizz.com/internal/api/handlers"
	"quizizz.com/internal/api/handlers/health"
	"quizizz.com/internal/api/handlers/ping"
	"quizizz.com/internal/api/handlers/user"
)

// API defines the API routes
type API struct {
	BaseHandler   *handlers.BaseHandler
	HealthHandler *health.Handler
	PingHandler   *ping.Handler
	UserHandler   *user.Handler
}

// NewAPI creates a new API routes instance
func NewAPI(
	baseHandler *handlers.BaseHandler,
	healthHandler *health.Handler,
	pingHandler *ping.Handler,
	userHandler *user.Handler,
) *API {
	return &API{
		BaseHandler:   baseHandler,
		HealthHandler: healthHandler,
		PingHandler:   pingHandler,
		UserHandler:   userHandler,
	}
}

// RegisterRoutes registers all the API routes
func (a *API) RegisterRoutes(router *gin.Engine) {
	// Health check routes
	router.GET("/_meta/health", a.HealthHandler.HealthCheck)
	router.GET("/livez", a.HealthHandler.LivenessCheck)
	router.GET("/readyz", a.HealthHandler.ReadinessCheck)

	// API group with versioning
	apiGroup := router.Group("/api")
	{
		v1 := apiGroup.Group("/v1")
		{
			// Ping endpoint
			v1.GET("/ping", a.PingHandler.Ping)

			// User routes
			users := v1.Group("/users")
			{
				users.GET("", a.UserHandler.ListUsers)
				users.POST("", a.UserHandler.CreateUser)
				users.GET("/:id", a.UserHandler.GetUser)
				users.PUT("/:id", a.UserHandler.UpdateUser)
				users.DELETE("/:id", a.UserHandler.DeleteUser)
			}
		}
	}
}
