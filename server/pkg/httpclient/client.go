package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"quizizz.com/internal/logger"
)

// HeaderRequestID is the header name for request ID
const HeaderRequestID = "X-Request-ID"

// Client is a robust HTTP client with enhanced features
type Client struct {
	config       *Config
	httpClient   *http.Client
	baseURL      *url.URL
	breaker      *gobreaker.CircuitBreaker
	retryBackOff backoff.BackOff
	serviceName  string
	tracer       trace.Tracer
}

// Response wraps an HTTP response
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	RequestID  string
	Duration   time.Duration
}

// New creates a new HTTP client
func New(cfg *Config) (*Client, error) {
	baseURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	transport := createTransport(cfg)
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeouts.RequestTimeout,
	}

	var cbSettings gobreaker.Settings
	if cfg.CircuitBreaker.Enabled {
		cbSettings = gobreaker.Settings{
			Name:        cfg.CircuitBreaker.Name,
			MaxRequests: cfg.CircuitBreaker.MaxRequests,
			Interval:    cfg.CircuitBreaker.Interval,
			Timeout:     cfg.CircuitBreaker.Timeout,
			ReadyToTrip: cfg.CircuitBreaker.ReadyToTrip,
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				logger.Info("Circuit breaker state changed",
					zap.String("name", name),
					zap.String("from", from.String()),
					zap.String("to", to.String()),
				)
			},
		}
	}

	var retryBackOff backoff.BackOff
	if cfg.Retry.Enabled {
		exponentialBackOff := backoff.NewExponentialBackOff()
		exponentialBackOff.InitialInterval = cfg.Retry.InitialInterval
		exponentialBackOff.MaxInterval = cfg.Retry.MaxInterval
		exponentialBackOff.MaxElapsedTime = cfg.Retry.MaxElapsedTime
		exponentialBackOff.Multiplier = cfg.Retry.Multiplier
		retryBackOff = backoff.WithMaxRetries(exponentialBackOff, uint64(cfg.Retry.MaxRetries))
	} else {
		retryBackOff = &backoff.StopBackOff{}
	}

	// Get tracer
	tracer := otel.GetTracerProvider().Tracer(cfg.ServiceName)

	client := &Client{
		config:       cfg,
		httpClient:   httpClient,
		baseURL:      baseURL,
		breaker:      gobreaker.NewCircuitBreaker(cbSettings),
		retryBackOff: retryBackOff,
		serviceName:  cfg.ServiceName,
		tracer:       tracer,
	}

	return client, nil
}

// createTransport creates an HTTP transport with configured settings
func createTransport(cfg *Config) *http.Transport {
	dialer := &net.Dialer{
		Timeout:   cfg.Timeouts.DialTimeout,
		KeepAlive: cfg.Timeouts.DialKeepAlive,
	}

	transport := &http.Transport{
		Proxy:                 getProxyFunc(cfg),
		DialContext:           dialer.DialContext,
		MaxIdleConns:          cfg.Transport.MaxIdleConns,
		MaxIdleConnsPerHost:   cfg.Transport.MaxIdleConnsPerHost,
		MaxConnsPerHost:       cfg.Transport.MaxConnsPerHost,
		TLSHandshakeTimeout:   cfg.Timeouts.TLSHandshakeTimeout,
		ResponseHeaderTimeout: cfg.Timeouts.ResponseHeaderTimeout,
		ExpectContinueTimeout: cfg.Timeouts.ExpectContinueTimeout,
		IdleConnTimeout:       cfg.Timeouts.IdleConnTimeout,
		DisableCompression:    cfg.Transport.DisableCompression,
		DisableKeepAlives:     cfg.Transport.DisableKeepAlives,
		ForceAttemptHTTP2:     true,
	}

	return transport
}

// getProxyFunc returns a proxy function if a proxy URL is configured
func getProxyFunc(cfg *Config) func(*http.Request) (*url.URL, error) {
	if cfg.Transport.ProxyURL == "" {
		return http.ProxyFromEnvironment
	}

	proxyURL, err := url.Parse(cfg.Transport.ProxyURL)
	if err != nil {
		logger.Warn("Invalid proxy URL, using environment proxy", zap.Error(err))
		return http.ProxyFromEnvironment
	}

	return http.ProxyURL(proxyURL)
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, urlPath string, headers map[string]string) (*Response, error) {
	return c.Request(ctx, http.MethodGet, urlPath, nil, headers)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, urlPath string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Request(ctx, http.MethodPost, urlPath, body, headers)
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, urlPath string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Request(ctx, http.MethodPut, urlPath, body, headers)
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, urlPath string, headers map[string]string) (*Response, error) {
	return c.Request(ctx, http.MethodDelete, urlPath, nil, headers)
}

// Patch performs a PATCH request
func (c *Client) Patch(ctx context.Context, urlPath string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Request(ctx, http.MethodPatch, urlPath, body, headers)
}

// GetJSON performs a GET request and unmarshals the response into the given target
func (c *Client) GetJSON(ctx context.Context, urlPath string, headers map[string]string, target interface{}) error {
	resp, err := c.Get(ctx, urlPath, headers)
	if err != nil {
		return err
	}

	return json.Unmarshal(resp.Body, target)
}

// PostJSON performs a POST request and unmarshals the response into the given target
func (c *Client) PostJSON(ctx context.Context, urlPath string, body, target interface{}, headers map[string]string) error {
	resp, err := c.Post(ctx, urlPath, body, headers)
	if err != nil {
		return err
	}

	return json.Unmarshal(resp.Body, target)
}

// PutJSON performs a PUT request and unmarshals the response into the given target
func (c *Client) PutJSON(ctx context.Context, urlPath string, body, target interface{}, headers map[string]string) error {
	resp, err := c.Put(ctx, urlPath, body, headers)
	if err != nil {
		return err
	}

	return json.Unmarshal(resp.Body, target)
}

// Request performs an HTTP request with retries and circuit breaking
func (c *Client) Request(ctx context.Context, method, urlPath string, body interface{}, headers map[string]string) (*Response, error) {
	// Ensure we have headers map initialized
	if headers == nil {
		headers = make(map[string]string)
	}

	// Resolve the full URL
	fullURL := c.createURL(urlPath)
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Create a span for the outgoing request
	spanName := fmt.Sprintf("%s %s", method, urlPath)
	ctx, span := c.tracer.Start(
		ctx,
		spanName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			semconv.HTTPMethodKey.String(method),
			semconv.HTTPURLKey.String(fullURL),
			semconv.HTTPTargetKey.String(parsedURL.Path),
			semconv.HTTPSchemeKey.String(parsedURL.Scheme),
			semconv.NetPeerNameKey.String(parsedURL.Hostname()),
			semconv.NetPeerPortKey.String(parsedURL.Port()),
		),
	)
	defer span.End()

	requestFunc := func() (*Response, error) {
		return c.doRequest(ctx, method, urlPath, body, headers)
	}

	// Apply circuit breaker pattern
	if c.config.CircuitBreaker.Enabled {
		result, err := c.breaker.Execute(func() (interface{}, error) {
			return c.executeWithRetries(ctx, requestFunc)
		})

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}

		return result.(*Response), nil
	}

	// Just use retries without circuit breaker
	response, err := c.executeWithRetries(ctx, requestFunc)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return response, err
}

// generateID generates a unique ID for request tracking
func generateID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), randomString(8))
}

// randomString generates a random string of the specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// executeWithRetries performs a request with retries based on the configured backoff
func (c *Client) executeWithRetries(ctx context.Context, requestFunc func() (*Response, error)) (*Response, error) {
	var response *Response
	var err error
	var statusCode int

	// Track attempt count for logging
	attempt := 0

	if !c.config.Retry.Enabled {
		return requestFunc()
	}

	operation := func() error {
		attempt++
		span := trace.SpanFromContext(ctx)

		response, err = requestFunc()
		if err != nil {
			logger.WarnCtx(ctx, "Request retry due to error",
				zap.Error(err),
				zap.Int("attempt", attempt),
				zap.Int("maxAttempts", c.config.Retry.MaxRetries+1),
			)
			span.AddEvent("retry",
				trace.WithAttributes(
					attribute.Int("attempt", attempt),
					attribute.String("error", err.Error()),
				),
			)
			return err
		}

		statusCode = response.StatusCode
		if c.config.Retry.ShouldRetry != nil && c.config.Retry.ShouldRetry(nil, statusCode) {
			logger.WarnCtx(ctx, "Request retry due to status code",
				zap.Int("statusCode", statusCode),
				zap.Int("attempt", attempt),
				zap.Int("maxAttempts", c.config.Retry.MaxRetries+1),
			)
			span.AddEvent("retry",
				trace.WithAttributes(
					attribute.Int("attempt", attempt),
					attribute.Int("statusCode", statusCode),
				),
			)
			return fmt.Errorf("request failed with status code %d", statusCode)
		}

		return nil
	}

	retryErr := backoff.Retry(operation, c.retryBackOff)
	if retryErr != nil {
		logger.ErrorCtx(ctx, "Request failed after all retries",
			zap.Error(retryErr),
			zap.Int("maxAttempts", c.config.Retry.MaxRetries+1),
		)
		return response, retryErr
	}

	return response, nil
}

// doRequest performs a single HTTP request
func (c *Client) doRequest(ctx context.Context, method, urlPath string, body interface{}, headers map[string]string) (*Response, error) {
	// Create the URL
	fullURL := c.createURL(urlPath)

	// Create the request body if needed
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			logger.ErrorCtx(ctx, "Error marshaling request body", zap.Error(err))
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)

		// Add body details to the current span
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(semconv.HTTPRequestContentLengthKey.Int(len(jsonBody)))
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		logger.ErrorCtx(ctx, "Error creating request", zap.Error(err))
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add default headers
	for key, value := range c.config.DefaultHeaders {
		req.Header.Set(key, value)
	}

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set request ID if not present
	if req.Header.Get(HeaderRequestID) == "" {
		requestID := generateID()
		req.Header.Set(HeaderRequestID, requestID)
	}

	// Inject OpenTelemetry context into headers
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	requestID := req.Header.Get(HeaderRequestID)

	// Add request attributes to the current span
	span := trace.SpanFromContext(ctx)
	if requestID != "" {
		span.SetAttributes(attribute.String("request.id", requestID))
	}

	// Log the request
	logger.InfoCtx(ctx, "HTTP request",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
	)

	startTime := time.Now()

	// Perform the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.ErrorCtx(ctx, "Error performing request",
			zap.Error(err),
			zap.String("method", req.Method),
			zap.String("url", req.URL.String()),
		)
		return nil, fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorCtx(ctx, "Error reading response body", zap.Error(err))
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	duration := time.Since(startTime)

	// Update span with response information
	span.SetAttributes(
		semconv.HTTPStatusCodeKey.Int(resp.StatusCode),
		semconv.HTTPResponseContentLengthKey.Int(len(respBody)),
	)

	// Add duration as a metric
	span.SetAttributes(attribute.Int64("http.duration_ms", duration.Milliseconds()))

	// Set span status based on response code
	if resp.StatusCode >= 400 {
		span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", resp.StatusCode))
	} else {
		span.SetStatus(codes.Ok, "")
	}

	// Create the response
	response := &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       respBody,
		RequestID:  requestID,
		Duration:   duration,
	}

	// Log based on status code
	if resp.StatusCode >= 400 {
		logger.WarnCtx(ctx, "HTTP request failed",
			zap.Int("statusCode", resp.StatusCode),
			zap.Duration("duration", duration),
		)
	} else {
		logger.InfoCtx(ctx, "HTTP response",
			zap.Int("statusCode", resp.StatusCode),
			zap.Duration("duration", duration),
			zap.Int("bodySize", len(respBody)),
		)
	}

	return response, nil
}

// createURL creates a full URL from the base URL and path
func (c *Client) createURL(urlPath string) string {
	if urlPath == "" {
		return c.config.BaseURL
	}

	// Handle absolute URLs
	if strings.HasPrefix(urlPath, "http://") || strings.HasPrefix(urlPath, "https://") {
		return urlPath
	}

	// Create a copy of the base URL
	u := *c.baseURL

	// Join the base path and the requested path
	u.Path = path.Join(u.Path, urlPath)

	return u.String()
}
