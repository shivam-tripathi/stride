// Package health provides health check handlers
package health

import (
	"github.com/gin-gonic/gin"
	"quizizz.com/internal/api/handlers"
	"quizizz.com/internal/api/response"
)

// Handler handles health check requests
type Handler struct {
	*handlers.BaseHandler
	version string
}

// NewHandler creates a new health handler
func NewHandler(base *handlers.BaseHandler, version string) *Handler {
	return &Handler{
		BaseHandler: base,
		version:     version,
	}
}

// HealthCheck handles the health check endpoint
func (h *Handler) HealthCheck(c *gin.Context) {
	logger := h.GetRequestLogger(c)
	logger.Debug("Health check requested")

	response.Success(c, gin.H{
		"status":  "ok",
		"version": h.version,
	})
}

// LivenessCheck handles Kubernetes liveness probe
func (h *Handler) LivenessCheck(c *gin.Context) {
	response.Success(c, gin.H{
		"status": "alive",
	})
}

// ReadinessCheck handles Kubernetes readiness probe
func (h *Handler) ReadinessCheck(c *gin.Context) {
	// Here you might check database connections, cache availability, etc.
	// For simplicity, we're just returning success
	response.Success(c, gin.H{
		"status": "ready",
	})
}
