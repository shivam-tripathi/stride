// Package ping provides a simple ping handler
package ping

import (
	"github.com/gin-gonic/gin"
	"quizizz.com/internal/api/handlers"
	"quizizz.com/internal/api/response"
)

// Handler handles ping requests
type Handler struct {
	*handlers.BaseHandler
}

// NewHandler creates a new ping handler
func NewHandler(base *handlers.BaseHandler) *Handler {
	return &Handler{
		BaseHandler: base,
	}
}

// Ping handles the ping endpoint
func (h *Handler) Ping(c *gin.Context) {
	message := h.AppService.GetPingMessage()
	response.Success(c, gin.H{
		"message": message,
	})
}
