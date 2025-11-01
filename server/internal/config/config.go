package config

import (
	"os"
	"strconv"
	"time"
)

// MongoDBConfig holds all MongoDB configuration
type MongoDBConfig struct {
	URI            string
	Database       string
	MaxPoolSize    uint64
	MinPoolSize    uint64
	ConnectTimeout time.Duration
	Timeout        time.Duration
}

// RedisConfig holds all Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Timeout  time.Duration
}

// OTELConfig holds configuration for OpenTelemetry
type OTELConfig struct {
	// Enabled determines if tracing is enabled
	Enabled bool

	// ServiceName is the name of the service
	ServiceName string

	// TracingExporterEndpoint is the endpoint for the tracing exporter (e.g. OTLP, Jaeger)
	TracingExporterEndpoint string

	// TracingExporterInsecure determines if the tracing exporter should use TLS
	TracingExporterInsecure bool

	// TracingSampleRatio is the ratio of traces to sample (0.0 - 1.0)
	TracingSampleRatio float64
}

// Config holds all configuration for the application
type Config struct {
	AppName  string
	Port     string
	LogLevel string
	Env      string

	// Resource configurations
	MongoDB MongoDBConfig
	Redis   RedisConfig
	OTEL    OTELConfig
}

// NewConfig creates a new Config
func NewConfig() *Config {
	return &Config{
		AppName:  getEnv("APP_NAME", "go-template-api"),
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Env:      getEnv("ENV", "development"),

		MongoDB: MongoDBConfig{
			URI:            getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database:       getEnv("MONGODB_DATABASE", "app"),
			MaxPoolSize:    uint64(getEnvAsInt("MONGODB_MAX_POOL_SIZE", 100)),
			MinPoolSize:    uint64(getEnvAsInt("MONGODB_MIN_POOL_SIZE", 10)),
			ConnectTimeout: getEnvAsDuration("MONGODB_CONNECT_TIMEOUT", 10*time.Second),
			Timeout:        getEnvAsDuration("MONGODB_TIMEOUT", 5*time.Second),
		},

		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			Timeout:  getEnvAsDuration("REDIS_TIMEOUT", 5*time.Second),
		},

		OTEL: OTELConfig{
			Enabled:                 getEnvAsBool("OTEL_ENABLED", true),
			ServiceName:             getEnv("OTEL_SERVICE_NAME", "go-template-api"),
			TracingExporterEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
			TracingExporterInsecure: getEnvAsBool("OTEL_EXPORTER_OTLP_INSECURE", true),
			TracingSampleRatio:      getEnvAsFloat("OTEL_TRACE_SAMPLER_ARG", 1.0),
		},
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// getEnvAsDuration retrieves an environment variable as a duration or returns a default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// getEnvAsBool retrieves an environment variable as a boolean or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// getEnvAsFloat retrieves an environment variable as a float or returns a default value
func getEnvAsFloat(key string, defaultValue float64) float64 {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return defaultValue
	}

	return value
}
