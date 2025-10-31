// Package handlers provides HTTP request handlers for the API
package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"quizizz.com/internal/service"
)

// BaseHandler contains common dependencies and utilities for handlers
type BaseHandler struct {
	// Common services that most handlers might need
	AppService service.AppService

	// You can add more common dependencies here like:
	// - DB connections
	// - Cache clients
	// - Authentication services, etc.
}

// NewBaseHandler creates a new BaseHandler
func NewBaseHandler(appService service.AppService) *BaseHandler {
	return &BaseHandler{
		AppService: appService,
	}
}

// ShouldBindJSON wraps gin's binding with error handling
func (h *BaseHandler) ShouldBindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(err)
		return false
	}
	return true
}

// GetRequestLogger returns a logger with request context
func (h *BaseHandler) GetRequestLogger(c *gin.Context) *zap.Logger {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	if !exists {
		requestID = "unknown"
	}

	return zap.L().With(
		zap.String("requestID", requestID.(string)),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)
}
