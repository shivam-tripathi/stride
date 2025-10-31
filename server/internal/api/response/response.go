// Package response provides standardized API response structures
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"quizizz.com/internal/errors"
)

// Response is the standard API response envelope
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error represents the error details in a response
type Error struct {
	Code    string                 `json:"code,omitempty"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Success sends a successful response with data
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// Created sends a 201 created response with data
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// NoContent sends a 204 no content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Fail sends an error response
func Fail(c *gin.Context, err error) {
	// Get status code from the error
	statusCode := errors.GetStatusCode(err)

	// Get context from the error
	contextMap := errors.GetContextMap(err)

	// Get user-friendly message
	message := errors.GetUserMessage(err)

	// Create error response
	errorResponse := Error{
		Message: message,
		Details: contextMap,
	}

	// Create a code based on the error if possible
	if statusCode == http.StatusBadRequest {
		errorResponse.Code = "BAD_REQUEST"
	} else if statusCode == http.StatusNotFound {
		errorResponse.Code = "NOT_FOUND"
	} else if statusCode == http.StatusInternalServerError {
		errorResponse.Code = "INTERNAL_ERROR"
	}

	c.JSON(statusCode, Response{
		Success: false,
		Error:   &errorResponse,
	})
}

// BadRequest sends a 400 bad request response
func BadRequest(c *gin.Context, message string) {
	Fail(c, errors.BadRequest(message))
}

// NotFound sends a 404 not found response
func NotFound(c *gin.Context, message string) {
	Fail(c, errors.NotFound(message))
}

// InternalError sends a 500 internal server error response
func InternalError(c *gin.Context, message string) {
	Fail(c, errors.Internal(message))
}

// InternalServerError sends a 500 internal server error response
// This is an alias for InternalError for better API consistency
func InternalServerError(c *gin.Context, message string) {
	InternalError(c, message)
}
