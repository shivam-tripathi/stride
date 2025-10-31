#!/bin/bash

echo "🚀 Setting up ca-service..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

echo "📦 Installing dependencies..."
go mod download
go mod tidy

echo "🔧 Installing wire tool..."
go install github.com/google/wire/cmd/wire@latest

echo "🔄 Generating wire dependencies..."
wire ./...

echo "📝 Creating .env file if it doesn't exist..."
if [ ! -f .env ]; then
    cat > .env << EOL
# Application Settings
APP_NAME=ca-service
PORT=8080
LOG_LEVEL=info
ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=app
DB_SSL_MODE=require
DB_MAX_OPEN=10
DB_MAX_IDLE=5
DB_TIMEOUT=5s

# Redis Configuration (disabled by default)
REDIS_ENABLED=false
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_TIMEOUT=5s

# OpenTelemetry Configuration
OTEL_ENABLED=true
OTEL_SERVICE_NAME=ca-service
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
OTEL_EXPORTER_OTLP_INSECURE=true
OTEL_TRACE_SAMPLER_ARG=1.0
EOL
    echo "✅ Created .env file with default configuration"
else
    echo "ℹ️ .env file already exists, skipping creation"
fi

echo "✨ Setup complete! You can now:"
echo "1. Update the .env file with your configuration"
echo "2. Run the application with: go run cmd/api/main.go"