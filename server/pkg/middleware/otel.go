package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"quizizz.com/internal/logger"
	"quizizz.com/pkg/httpclient"
)

// TracingContextKey is the key used to store the tracing context in the gin.Context
const TracingContextKey = "tracing-context"

// requestIDKey is the key used to store the request ID in the gin.Context
const requestIDKey = "requestID"

// OTEL returns a middleware that traces incoming requests
func OTEL(serviceName string) gin.HandlerFunc {
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(c *gin.Context) {
		// Extract tracing context from the incoming request headers
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// Start a new span for this request
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
		if c.FullPath() == "" {
			spanName = fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
		}

		// Define attributes for the span
		attrs := []attribute.KeyValue{
			semconv.HTTPMethodKey.String(c.Request.Method),
			semconv.HTTPRouteKey.String(c.FullPath()),
			semconv.HTTPURLKey.String(c.Request.URL.String()),
			semconv.HTTPRequestContentLengthKey.Int64(c.Request.ContentLength),
			semconv.HTTPSchemeKey.String(c.Request.URL.Scheme),
			attribute.String("http.host", c.Request.Host),
		}

		// Add user agent if available
		if ua := c.Request.UserAgent(); ua != "" {
			attrs = append(attrs, attribute.String("http.user_agent", ua))
		}

		// Add client IP information
		clientIP := c.ClientIP()
		if clientIP != "" {
			attrs = append(attrs, attribute.String("http.client_ip", clientIP))
			attrs = append(attrs, attribute.String("net.peer.ip", clientIP))
		}

		ctx, span := tracer.Start(
			ctx,
			spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(attrs...),
		)
		defer span.End()

		// Store the context in the request
		c.Request = c.Request.WithContext(ctx)

		// Also store it in gin.Context for easy retrieval
		c.Set(TracingContextKey, ctx)

		// Get request ID and add it to the span
		if requestID, exists := c.Get(requestIDKey); exists && requestID != nil {
			span.SetAttributes(attribute.String("request.id", requestID.(string)))
			// Also set the request ID in the response header for correlation
			c.Writer.Header().Set(httpclient.HeaderRequestID, requestID.(string))
		}

		// Record the start time
		startTime := time.Now()

		// Process the request
		c.Next()

		// Record the response status
		status := c.Writer.Status()
		respSize := c.Writer.Size()

		// Set response attributes
		span.SetAttributes(
			semconv.HTTPStatusCodeKey.Int(status),
			semconv.HTTPResponseContentLengthKey.Int(respSize),
		)

		// Record request duration
		duration := time.Since(startTime)
		span.SetAttributes(attribute.Int64("http.duration_ms", duration.Milliseconds()))

		// Set span status based on response code
		if status >= http.StatusInternalServerError {
			span.SetStatus(codes.Error, fmt.Sprintf("error status code: %d", status))
			if len(c.Errors) > 0 {
				span.RecordError(c.Errors[0])
			}
		} else {
			span.SetStatus(codes.Ok, "")
		}

		// Log the request with the trace context
		logger.InfoCtx(ctx, "HTTP request completed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.Int("size", respSize),
		)
	}
}

// ContextFromGin extracts the tracing context from a gin context
func ContextFromGin(c *gin.Context) context.Context {
	if tracingCtx, exists := c.Get(TracingContextKey); exists && tracingCtx != nil {
		return tracingCtx.(context.Context)
	}
	return c.Request.Context()
}

// StartSpanFromGin starts a new span from a gin context
func StartSpanFromGin(c *gin.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	ctx := ContextFromGin(c)
	tracer := otel.Tracer("")
	return tracer.Start(ctx, name, opts...)
}
