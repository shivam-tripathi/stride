// Package httpclient provides a robust HTTP client with advanced features
// such as retries, circuit breaking, timeouts, and metrics.
package httpclient

import (
	"time"

	"github.com/sony/gobreaker"
)

// RetryConfig holds configuration for the retry mechanism
type RetryConfig struct {
	// Enabled determines if retries are enabled
	Enabled bool

	// MaxRetries is the maximum number of times to retry a request
	MaxRetries int

	// InitialInterval is the initial interval between retries
	InitialInterval time.Duration

	// MaxInterval is the maximum interval between retries
	MaxInterval time.Duration

	// MaxElapsedTime is the maximum elapsed time for retries
	MaxElapsedTime time.Duration

	// Multiplier is the factor by which the interval increases
	Multiplier float64

	// ShouldRetry is a function that determines if a request should be retried
	ShouldRetry func(err error, statusCode int) bool
}

// CircuitBreakerConfig holds configuration for the circuit breaker
type CircuitBreakerConfig struct {
	// Enabled determines if circuit breaking is enabled
	Enabled bool

	// Name is the name of the circuit breaker
	Name string

	// MaxRequests is the maximum number of requests allowed when the circuit is half-open
	MaxRequests uint32

	// Interval is the cyclic period of the closed state
	Interval time.Duration

	// Timeout is the period of the open state, after which the state becomes half-open
	Timeout time.Duration

	// ReadyToTrip is a function that determines if the circuit breaker should trip
	ReadyToTrip func(counts gobreaker.Counts) bool
}

// TimeoutConfig holds configuration for various timeouts
type TimeoutConfig struct {
	// RequestTimeout is the maximum time for the whole request
	RequestTimeout time.Duration

	// DialTimeout is the maximum time for connecting to the server
	DialTimeout time.Duration

	// DialKeepAlive is the interval between keep-alive probes
	DialKeepAlive time.Duration

	// TLSHandshakeTimeout is the maximum time for TLS handshake
	TLSHandshakeTimeout time.Duration

	// ResponseHeaderTimeout is the maximum time to wait for a server's response headers
	ResponseHeaderTimeout time.Duration

	// ExpectContinueTimeout is the maximum time to wait for a server's first
	// response headers after fully writing the request headers
	ExpectContinueTimeout time.Duration

	// IdleConnTimeout is the maximum time an idle connection will remain idle before closing
	IdleConnTimeout time.Duration
}

// TransportConfig holds configuration for the HTTP transport
type TransportConfig struct {
	// MaxIdleConns is the maximum number of idle connections
	MaxIdleConns int

	// MaxIdleConnsPerHost is the maximum number of idle connections per host
	MaxIdleConnsPerHost int

	// MaxConnsPerHost is the maximum number of connections per host
	MaxConnsPerHost int

	// DisableCompression disables compression
	DisableCompression bool

	// DisableKeepAlives disables keep-alives
	DisableKeepAlives bool

	// ProxyURL is the URL of the proxy to use
	ProxyURL string
}

// Config holds all configuration options for the HTTP client
type Config struct {
	// BaseURL is the base URL for all requests
	BaseURL string

	// ServiceName is the name of the service making the requests
	// This is used for logging and tracing
	ServiceName string

	// DefaultHeaders are headers that will be included in all requests
	DefaultHeaders map[string]string

	// Timeouts configuration
	Timeouts TimeoutConfig

	// Transport configuration
	Transport TransportConfig

	// Retry configuration
	Retry RetryConfig

	// CircuitBreaker configuration
	CircuitBreaker CircuitBreakerConfig

	// Tracing determines if tracing is enabled
	Tracing bool

	// Debug enables verbose logging
	Debug bool
}

// DefaultConfig returns a default configuration
func DefaultConfig(baseURL string) *Config {
	return &Config{
		BaseURL:     baseURL,
		ServiceName: "httpclient",
		DefaultHeaders: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			"User-Agent":   "go-httpclient/1.0",
		},
		Timeouts: TimeoutConfig{
			RequestTimeout:        30 * time.Second,
			DialTimeout:           5 * time.Second,
			DialKeepAlive:         15 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			IdleConnTimeout:       90 * time.Second,
		},
		Transport: TransportConfig{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     100,
			DisableCompression:  false,
			DisableKeepAlives:   false,
		},
		Retry: RetryConfig{
			Enabled:         true,
			MaxRetries:      3,
			InitialInterval: 100 * time.Millisecond,
			MaxInterval:     2 * time.Second,
			MaxElapsedTime:  10 * time.Second,
			Multiplier:      2.0,
			ShouldRetry: func(err error, statusCode int) bool {
				if err != nil {
					return true
				}
				return statusCode >= 500 || statusCode == 0 || statusCode == 429
			},
		},
		CircuitBreaker: CircuitBreakerConfig{
			Enabled:     true,
			Name:        "httpclient",
			MaxRequests: 100,
			Interval:    0, // disabled
			Timeout:     5 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.Requests >= 10 && failureRatio >= 0.6
			},
		},
		Tracing: true,
		Debug:   false,
	}
}

// WithBaseURL sets the base URL
func (c *Config) WithBaseURL(baseURL string) *Config {
	c.BaseURL = baseURL
	return c
}

// WithServiceName sets the service name
func (c *Config) WithServiceName(serviceName string) *Config {
	c.ServiceName = serviceName
	return c
}

// WithDefaultHeader adds a default header
func (c *Config) WithDefaultHeader(key, value string) *Config {
	if c.DefaultHeaders == nil {
		c.DefaultHeaders = make(map[string]string)
	}
	c.DefaultHeaders[key] = value
	return c
}

// WithRequestTimeout sets the request timeout
func (c *Config) WithRequestTimeout(timeout time.Duration) *Config {
	c.Timeouts.RequestTimeout = timeout
	return c
}

// WithRetryEnabled enables or disables retries
func (c *Config) WithRetryEnabled(enabled bool) *Config {
	c.Retry.Enabled = enabled
	return c
}

// WithMaxRetries sets the maximum number of retries
func (c *Config) WithMaxRetries(maxRetries int) *Config {
	c.Retry.MaxRetries = maxRetries
	return c
}

// WithCircuitBreakerEnabled enables or disables circuit breaking
func (c *Config) WithCircuitBreakerEnabled(enabled bool) *Config {
	c.CircuitBreaker.Enabled = enabled
	return c
}

// WithDebug enables or disables debug logging
func (c *Config) WithDebug(debug bool) *Config {
	c.Debug = debug
	return c
}

// WithDisableKeepAlives sets whether to disable keep-alives
func (c *Config) WithDisableKeepAlives(disable bool) *Config {
	c.Transport.DisableKeepAlives = disable
	return c
}

// WithMaxIdleConns sets the maximum number of idle connections
func (c *Config) WithMaxIdleConns(maxIdleConns int) *Config {
	c.Transport.MaxIdleConns = maxIdleConns
	return c
}
