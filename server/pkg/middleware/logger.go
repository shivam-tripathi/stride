// Package middleware provides HTTP middleware functions
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"quizizz.com/internal/logger"
)

// requestLog contains the structured fields for request logging
type requestLog struct {
	ClientIP   string        `json:"clientIp"`
	Method     string        `json:"method"`
	Path       string        `json:"path"`
	Query      string        `json:"query,omitempty"`
	UserAgent  string        `json:"userAgent,omitempty"`
	StatusCode int           `json:"statusCode"`
	Latency    time.Duration `json:"latency"`
	Error      string        `json:"error,omitempty"`
	BodySize   int           `json:"bodySize"`
	RequestID  string        `json:"requestId,omitempty"`
}

// Logger returns a gin middleware for logging HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		requestID := c.GetHeader("X-Request-ID")

		// Add request ID to context for downstream handlers
		if requestID != "" {
			c.Set("requestID", requestID)
		}

		// Process request
		c.Next()

		// Collect log data
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Add HTTP headers for API responses
		c.Header("X-Response-Time", latency.String())

		// Build log structure
		logData := requestLog{
			ClientIP:   clientIP,
			Method:     c.Request.Method,
			Path:       path,
			Query:      query,
			UserAgent:  userAgent,
			StatusCode: statusCode,
			Latency:    latency,
			BodySize:   bodySize,
			RequestID:  requestID,
		}

		// Get error (if any)
		if len(c.Errors) > 0 {
			logData.Error = c.Errors.String()
		}

		// Determine log level based on status code
		var logFunc func(string, ...zap.Field)
		switch {
		case statusCode >= 500:
			logFunc = logger.Error
		case statusCode >= 400:
			logFunc = logger.Warn
		case statusCode >= 300:
			logFunc = logger.Info
		default:
			logFunc = logger.Info
		}

		// Create log fields
		logFields := []zap.Field{
			zap.String("clientIP", logData.ClientIP),
			zap.String("method", logData.Method),
			zap.String("path", logData.Path),
			zap.Int("statusCode", logData.StatusCode),
			zap.Duration("latency", logData.Latency),
			zap.Int("bodySize", logData.BodySize),
		}

		// Add optional fields
		if logData.Error != "" {
			logFields = append(logFields, zap.String("error", logData.Error))
		}
		if logData.Query != "" {
			logFields = append(logFields, zap.String("query", logData.Query))
		}
		if logData.UserAgent != "" {
			logFields = append(logFields, zap.String("userAgent", logData.UserAgent))
		}
		if logData.RequestID != "" {
			logFields = append(logFields, zap.String("requestID", logData.RequestID))
		}

		// Log the request
		logFunc("http-request", logFields...)
	}
}

// RequestID is a middleware that generates a unique ID for each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use X-Request-ID from the request if it exists
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// Generate a random request ID (in a real app, use a proper UUID generator)
			requestID = time.Now().Format("20060102150405.000000")
		}

		// Set the request ID in the context and response header
		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// Recovery returns a middleware that recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error with stack trace
				logger.Error("http-panic",
					zap.Any("error", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("clientIP", c.ClientIP()),
				)

				// Return a 500 error
				c.AbortWithStatusJSON(500, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "INTERNAL_ERROR",
						"message": "An unexpected error occurred",
					},
				})
			}
		}()

		// Process request
		c.Next()
	}
}
