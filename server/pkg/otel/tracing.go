// Package otel provides OpenTelemetry instrumentation for the application
package otel

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"quizizz.com/internal/config"
	"quizizz.com/internal/logger"
)

var (
	// Global tracer provider
	tracerProvider *sdktrace.TracerProvider

	// Global tracer
	tracer trace.Tracer

	// To ensure we only initialize once
	once sync.Once
)

// Config holds configuration for OpenTelemetry
type Config struct {
	Enabled                 bool
	ServiceName             string
	TracingExporterEndpoint string
	TracingExporterInsecure bool
	TracingSampleRatio      float64
}

// InitTracer initializes the OpenTelemetry tracer
func InitTracer(ctx context.Context, cfg *config.Config) (*sdktrace.TracerProvider, error) {
	var err error

	once.Do(func() {
		logger.Info("Initializing OpenTelemetry tracer",
			zap.String("service", cfg.OTEL.ServiceName),
			zap.String("endpoint", cfg.OTEL.TracingExporterEndpoint),
		)

		// If OTEL is disabled, use a no-op tracer
		if !cfg.OTEL.Enabled {
			logger.Info("OpenTelemetry tracing is disabled")
			tracerProvider = sdktrace.NewTracerProvider()
			tracer = tracerProvider.Tracer(cfg.OTEL.ServiceName)
			otel.SetTracerProvider(tracerProvider)
			otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			))
			return
		}

		// Create a new OTLP exporter
		var traceExporter *otlptrace.Exporter
		var opts []otlptracegrpc.Option

		// Set endpoint
		opts = append(opts, otlptracegrpc.WithEndpoint(cfg.OTEL.TracingExporterEndpoint))

		// Use insecure credentials if configured, otherwise use the default secure credentials
		if cfg.OTEL.TracingExporterInsecure {
			logger.Info("Using insecure connection for OTLP exporter")
			opts = append(opts, otlptracegrpc.WithInsecure())
		}

		// Add gRPC dial option
		opts = append(opts, otlptracegrpc.WithDialOption(grpc.WithBlock()))

		traceExporter, err = otlptracegrpc.New(ctx, opts...)

		if err != nil {
			err = fmt.Errorf("failed to create trace exporter: %w", err)
			logger.Error("Failed to create trace exporter", zap.Error(err))
			return
		}

		// Create a resource that identifies your service
		res, resErr := resource.New(ctx,
			resource.WithAttributes(
				semconv.ServiceName(cfg.OTEL.ServiceName),
				attribute.String("environment", cfg.Env),
			),
		)

		if resErr != nil {
			err = fmt.Errorf("failed to create resource: %w", resErr)
			logger.Error("Failed to create resource", zap.Error(err))
			return
		}

		// Configure trace sampling
		samplingRatio := cfg.OTEL.TracingSampleRatio
		sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(samplingRatio))

		// Create a trace provider with the exporter
		tracerProvider = sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sampler),
			sdktrace.WithBatcher(traceExporter),
			sdktrace.WithResource(res),
		)

		// Set the global trace provider and propagator
		otel.SetTracerProvider(tracerProvider)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))

		// Create a tracer
		tracer = tracerProvider.Tracer(cfg.OTEL.ServiceName)

		logger.Info("OpenTelemetry tracer initialized successfully",
			zap.String("service", cfg.OTEL.ServiceName),
			zap.Float64("samplingRatio", samplingRatio),
		)
	})

	return tracerProvider, err
}

// Shutdown gracefully shuts down the tracer provider
func Shutdown(ctx context.Context) error {
	if tracerProvider == nil {
		return nil
	}

	// Allow some time for traces to be flushed
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger.Info("Shutting down OpenTelemetry tracer")
	err := tracerProvider.Shutdown(ctx)
	if err != nil {
		logger.Error("Error shutting down tracer provider", zap.Error(err))
	}

	return err
}

// Tracer returns the global tracer
func Tracer() trace.Tracer {
	return tracer
}

// SpanFromContext returns the current span from the context
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// ContextWithSpan adds a span to the context
func ContextWithSpan(ctx context.Context, span trace.Span) context.Context {
	return trace.ContextWithSpan(ctx, span)
}

// ExtractTraceInfo extracts trace and span IDs from context
func ExtractTraceInfo(ctx context.Context) (traceID, spanID string) {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return "", ""
	}

	return span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String()
}

// StartSpan starts a new span
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, name, opts...)
}
