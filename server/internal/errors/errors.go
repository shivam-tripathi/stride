// Package errors provides application-specific error types and handling
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Standard errors that can be used directly
var (
	ErrNotFound           = errors.New("resource not found")
	ErrBadRequest         = errors.New("bad request")
	ErrInternal           = errors.New("internal server error")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrConflict           = errors.New("conflict")
	ErrServiceUnavailable = errors.New("service unavailable")
)

// AppError represents an application-specific error
type AppError struct {
	// Original is the original error (if any)
	Original error

	// StatusCode is the associated HTTP status code (if any)
	StatusCode int

	// Message is the user-facing error message
	Message string

	// Operational indicates whether the error is operational or programmer error
	Operational bool

	// Context contains additional metadata about the error
	Context map[string]interface{}
}

// Error makes AppError implement the error interface
func (e *AppError) Error() string {
	if e.Original != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Original)
	}
	return e.Message
}

// Unwrap enables errors.Is, errors.As functionality
func (e *AppError) Unwrap() error {
	return e.Original
}

// WithContext adds context to an error
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// New creates a new error with a message
func New(message string) error {
	return &AppError{
		Message: message,
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	return &AppError{
		Original: err,
		Message:  message,
	}
}

// Wrapf wraps an error with a formatted message
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &AppError{
		Original: err,
		Message:  fmt.Sprintf(format, args...),
	}
}

// HTTPError creates an error with HTTP status code
func HTTPError(statusCode int, message string) error {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// BadRequest creates a 400 error
func BadRequest(message string) error {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		Original:   ErrBadRequest,
	}
}

// NotFound creates a 404 error
func NotFound(message string) error {
	return &AppError{
		StatusCode: http.StatusNotFound,
		Message:    message,
		Original:   ErrNotFound,
	}
}

// Internal creates a 500 error
func Internal(message string) error {
	return &AppError{
		StatusCode:  http.StatusInternalServerError,
		Message:     message,
		Original:    ErrInternal,
		Operational: true,
	}
}

// GetStatusCode extracts the HTTP status code from an error
func GetStatusCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) && appErr.StatusCode != 0 {
		return appErr.StatusCode
	}

	// Map standard errors to appropriate status codes
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrConflict):
		return http.StatusConflict
	case errors.Is(err, ErrServiceUnavailable):
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// GetContextMap extracts the context map from an error
func GetContextMap(err error) map[string]interface{} {
	var appErr *AppError
	if errors.As(err, &appErr) && appErr.Context != nil {
		return appErr.Context
	}
	return nil
}

// GetUserMessage extracts a user-friendly message from an error
func GetUserMessage(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) && appErr.Message != "" {
		return appErr.Message
	}
	return err.Error()
}
