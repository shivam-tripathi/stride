// Package logger provides a simplified interface to zap logging
package logger

import (
	"context"
	"os"
	"sync"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// global logger instance
	globalLogger *zap.Logger
	once         sync.Once
	// Default log level
	logLevel = zapcore.InfoLevel
)

// Init initializes the global logger
func Init(env string) {
	once.Do(func() {
		// Configure based on environment
		var config zap.Config
		if env == "development" {
			config = zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			logLevel = zapcore.DebugLevel
		} else {
			config = zap.NewProductionConfig()
			logLevel = zapcore.InfoLevel
		}

		// Common configuration
		config.EncoderConfig.TimeKey = "time"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		// Add this line to include function names in the logs
		config.EncoderConfig.FunctionKey = "function"

		// Set the level from our package variable
		config.Level = zap.NewAtomicLevelAt(logLevel)

		var err error
		// Add AddCallerSkip(1) to skip the logger wrapper and show the actual caller
		globalLogger, err = config.Build(zap.AddCallerSkip(1))
		if err != nil {
			// If we can't initialize the logger, use a simple stdout logger
			core := zapcore.NewCore(
				zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
				zapcore.AddSync(os.Stdout),
				zapcore.InfoLevel,
			)
			globalLogger = zap.New(core)
			globalLogger.Error("Failed to initialize the structured logger", zap.Error(err))
		}
	})
}

// SetLevel sets the logging level
func SetLevel(level zapcore.Level) {
	logLevel = level
	if globalLogger != nil {
		// Create a new atomicLevel and update the global logger
		atomicLevel := zap.NewAtomicLevelAt(level)

		// We need to recreate the logger with the new level
		config := zap.NewProductionConfig()
		config.Level = atomicLevel
		config.EncoderConfig.TimeKey = "time"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.FunctionKey = "function"

		newLogger, err := config.Build(zap.AddCallerSkip(1))
		if err == nil {
			// If successful, replace the global logger
			globalLogger = newLogger
		}
	}
}

// GetLevel returns the current logging level
func GetLevel() zapcore.Level {
	return logLevel
}

// Info logs an info level message with structured context
func Info(msg string, fields ...zap.Field) {
	ensureLogger()
	globalLogger.Info(msg, fields...)
}

// InfoCtx logs an info level message with structured context, including trace information
func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	ensureLogger()
	globalLogger.Info(msg, appendTraceFields(ctx, fields)...)
}

// Error logs an error level message with structured context
func Error(msg string, fields ...zap.Field) {
	ensureLogger()
	globalLogger.Error(msg, fields...)
}

// ErrorCtx logs an error level message with structured context, including trace information
func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	ensureLogger()
	globalLogger.Error(msg, appendTraceFields(ctx, fields)...)
}

// Debug logs a debug level message with structured context
func Debug(msg string, fields ...zap.Field) {
	ensureLogger()
	globalLogger.Debug(msg, fields...)
}

// DebugCtx logs a debug level message with structured context, including trace information
func DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	ensureLogger()
	globalLogger.Debug(msg, appendTraceFields(ctx, fields)...)
}

// Warn logs a warning level message with structured context
func Warn(msg string, fields ...zap.Field) {
	ensureLogger()
	globalLogger.Warn(msg, fields...)
}

// WarnCtx logs a warning level message with structured context, including trace information
func WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	ensureLogger()
	globalLogger.Warn(msg, appendTraceFields(ctx, fields)...)
}

// Fatal logs a fatal level message with structured context and exits
func Fatal(msg string, fields ...zap.Field) {
	ensureLogger()
	globalLogger.Fatal(msg, fields...)
}

// FatalCtx logs a fatal level message with structured context, including trace information, and exits
func FatalCtx(ctx context.Context, msg string, fields ...zap.Field) {
	ensureLogger()
	globalLogger.Fatal(msg, appendTraceFields(ctx, fields)...)
}

// With creates a child logger with additional context
func With(fields ...zap.Field) *zap.Logger {
	ensureLogger()
	return globalLogger.With(fields...)
}

// WithCtx creates a child logger with additional context, including trace information
func WithCtx(ctx context.Context, fields ...zap.Field) *zap.Logger {
	ensureLogger()
	return globalLogger.With(appendTraceFields(ctx, fields)...)
}

// ensureLogger initializes the logger if it hasn't been done yet
func ensureLogger() {
	if globalLogger == nil {
		Init("development") // Default to development if not explicitly initialized
	}
}

// Sync flushes any buffered log entries - useful for clean shutdown
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// appendTraceFields adds trace and span IDs from the context to the field list
func appendTraceFields(ctx context.Context, fields []zap.Field) []zap.Field {
	if ctx == nil {
		return fields
	}

	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return fields
	}

	// Create a new slice with a capacity for the original fields plus the trace fields
	newFields := make([]zap.Field, 0, len(fields)+2)
	newFields = append(newFields, fields...)

	if spanCtx.TraceID().IsValid() {
		newFields = append(newFields, zap.String("trace_id", spanCtx.TraceID().String()))
	}

	if spanCtx.SpanID().IsValid() {
		newFields = append(newFields, zap.String("span_id", spanCtx.SpanID().String()))
	}

	return newFields
}
