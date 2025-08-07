[![Dependabot Updates](https://github.com/fgrzl/telemetry/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/fgrzl/telemetry/actions/workflows/dependabot/dependabot-updates)
[![ci](https://github.com/fgrzl/telemetry/actions/workflows/ci.yml/badge.svg)](https://github.com/fgrzl/telemetry/actions/workflows/ci.yml)

# Telemetry

A Go library that provides OpenTelemetry utilities and extensions for tracing, logging correlation, and context management in distributed systems.

## Features

- **OpenTelemetry Configuration**: Easy setup of OTLP exporters with support for both gRPC and HTTP protocols
- **DataDog Integration**: Special handling for DataDog trace endpoints with proper authentication
- **Context Utilities**: Correlation and Causation ID management for distributed tracing
- **Log Correlation**: Automatic correlation of structured logs with trace and span IDs

## Installation

```bash
go get github.com/fgrzl/telemetry
```

## Quick Start

### Basic OpenTelemetry Setup

```go
package main

import (
    "context"
    "log"
    
    "github.com/fgrzl/telemetry"
)

func main() {
    ctx := context.Background()
    
    // Configure OpenTelemetry for your service
    shutdown := telemetry.ConfigureService(ctx, "my-service", "1.0.0")
    defer func() {
        if err := shutdown(ctx); err != nil {
            log.Printf("Error shutting down telemetry: %v", err)
        }
    }()
    
    // Your application code here...
}
```

### Using Correlation and Causation IDs

```go
package main

import (
    "context"
    "log/slog"
    
    "github.com/google/uuid"
    "github.com/fgrzl/telemetry"
)

func main() {
    ctx := context.Background()
    
    // Add correlation ID for distributed tracing
    correlationID := uuid.New()
    ctx = telemetry.WithCorrelationID(ctx, correlationID)
    
    // Add causation ID for immediate operation tracking
    causationID := uuid.New()
    ctx = telemetry.WithCausationID(ctx, causationID)
    
    // Retrieve IDs later
    if id, found := telemetry.CorrelationIDFromContext(ctx); found {
        slog.Info("Processing request", "correlation_id", id)
    }
}
```

### Log and Trace Correlation

```go
package main

import (
    "context"
    "log/slog"
    
    "github.com/fgrzl/telemetry"
    "go.opentelemetry.io/otel"
)

func processRequest(ctx context.Context) {
    // Start a span
    tracer := otel.GetTracerProvider().Tracer("my-service")
    ctx, span := tracer.Start(ctx, "process-request")
    defer span.End()
    
    // Create a logger correlated with the trace
    logger := telemetry.CorrelateLogsAndTraces(ctx, slog.Default())
    
    // Your logs will now include trace and span IDs
    logger.Info("Processing request started")
}
```

## Configuration

The library is configured through environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP endpoint URL | `localhost:4317` | No |
| `DD_API_KEY` | DataDog API key | - | Yes (if using DataDog) |
| `ENVIRONMENT` | Deployment environment | - | No |

### Example Configuration

```bash
# For local development with OTLP collector
export OTEL_EXPORTER_OTLP_ENDPOINT="localhost:4317"
export ENVIRONMENT="development"

# For DataDog integration
export OTEL_EXPORTER_OTLP_ENDPOINT="https://trace.agent.datadoghq.com"
export DD_API_KEY="your-datadog-api-key"
export ENVIRONMENT="production"
```

## API Reference

### ConfigureService

```go
func ConfigureService(ctx context.Context, serviceName, version string) func(context.Context) error
```

Initializes OpenTelemetry tracing and metrics for a service. Returns a shutdown function that should be called when the application terminates.

**Parameters:**
- `ctx`: Context for initialization
- `serviceName`: Name of your service (used in resource attributes)
- `version`: Version of your service (used in resource attributes)

**Returns:**
- Shutdown function that cleans up resources

### Context Utilities

```go
// Correlation ID management
func WithCorrelationID(ctx context.Context, correlationID uuid.UUID) context.Context
func CorrelationIDFromContext(ctx context.Context) (uuid.UUID, bool)

// Causation ID management  
func WithCausationID(ctx context.Context, causationID uuid.UUID) context.Context
func CausationIDFromContext(ctx context.Context) (uuid.UUID, bool)
```

### Log Correlation

```go
func CorrelateLogsAndTraces(ctx context.Context, log *slog.Logger) *slog.Logger
```

Returns a logger enriched with trace and span IDs from the context, including DataDog-compatible decimal representations.

## Development

### Prerequisites

- Go 1.23 or later
- Access to an OTLP collector or DataDog account for testing

### Running Tests

```bash
# Run all tests
go test -v

# Run tests with coverage
go test -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Building

```bash
go build ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for your changes
5. Ensure tests pass (`go test -v`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Code Standards

- Follow Go idioms and best practices
- Add tests for all new functionality
- Document exported functions with GoDoc comments
- Use structured logging with `slog`
- Handle errors appropriately

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
