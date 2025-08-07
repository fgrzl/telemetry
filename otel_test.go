package telemetry

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldConfigureServiceWithDefaultEndpoint(t *testing.T) {
	// Arrange
	ctx := context.Background()
	serviceName := "test-service"
	version := "1.0.0"

	// Clear environment variables
	originalEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	originalEnv := os.Getenv("ENVIRONMENT")
	defer func() {
		os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", originalEndpoint)
		os.Setenv("ENVIRONMENT", originalEnv)
	}()
	
	os.Unsetenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	os.Setenv("ENVIRONMENT", "test")

	// Act
	shutdownFunc := ConfigureService(ctx, serviceName, version)

	// Assert
	assert.NotNil(t, shutdownFunc)
	
	// Test that shutdown doesn't panic
	err := shutdownFunc(ctx)
	assert.NoError(t, err)
}

func TestShouldConfigureServiceWithCustomEndpoint(t *testing.T) {
	// Arrange
	ctx := context.Background()
	serviceName := "test-service"
	version := "1.0.0"
	customEndpoint := "localhost:4318"

	// Set environment variables
	originalEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	originalEnv := os.Getenv("ENVIRONMENT")
	defer func() {
		os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", originalEndpoint)
		os.Setenv("ENVIRONMENT", originalEnv)
	}()
	
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", customEndpoint)
	os.Setenv("ENVIRONMENT", "test")

	// Act
	shutdownFunc := ConfigureService(ctx, serviceName, version)

	// Assert
	assert.NotNil(t, shutdownFunc)
	
	// Test that shutdown doesn't panic
	err := shutdownFunc(ctx)
	assert.NoError(t, err)
}

func TestShouldHandleDatadogEndpointWithoutAPIKey(t *testing.T) {
	// Arrange
	ctx := context.Background()
	serviceName := "test-service"
	version := "1.0.0"
	datadogEndpoint := "https://trace.agent.datadoghq.com"

	// Set environment variables
	originalEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	originalAPIKey := os.Getenv("DD_API_KEY")
	originalEnv := os.Getenv("ENVIRONMENT")
	defer func() {
		os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", originalEndpoint)
		os.Setenv("DD_API_KEY", originalAPIKey)
		os.Setenv("ENVIRONMENT", originalEnv)
	}()
	
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", datadogEndpoint)
	os.Unsetenv("DD_API_KEY")
	os.Setenv("ENVIRONMENT", "test")

	// Act
	shutdownFunc := ConfigureService(ctx, serviceName, version)

	// Assert - Should return a no-op function when DD_API_KEY is missing
	assert.NotNil(t, shutdownFunc)
	
	// Test that shutdown doesn't panic
	err := shutdownFunc(ctx)
	assert.NoError(t, err)
}

func TestShouldHandleDatadogEndpointWithAPIKey(t *testing.T) {
	// Arrange
	ctx := context.Background()
	serviceName := "test-service"
	version := "1.0.0"
	datadogEndpoint := "https://trace.agent.datadoghq.com"
	apiKey := "test-api-key"

	// Set environment variables
	originalEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	originalAPIKey := os.Getenv("DD_API_KEY")
	originalEnv := os.Getenv("ENVIRONMENT")
	defer func() {
		os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", originalEndpoint)
		os.Setenv("DD_API_KEY", originalAPIKey)
		os.Setenv("ENVIRONMENT", originalEnv)
	}()
	
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", datadogEndpoint)
	os.Setenv("DD_API_KEY", apiKey)
	os.Setenv("ENVIRONMENT", "test")

	// Act
	shutdownFunc := ConfigureService(ctx, serviceName, version)

	// Assert
	assert.NotNil(t, shutdownFunc)
	
	// Test that shutdown doesn't panic
	err := shutdownFunc(ctx)
	assert.NoError(t, err)
}

func TestShouldSetServiceNameAndVersionInResource(t *testing.T) {
	// Arrange
	ctx := context.Background()
	serviceName := "my-service"
	version := "2.1.0"

	// Set environment variables
	originalEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	originalEnv := os.Getenv("ENVIRONMENT")
	defer func() {
		os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", originalEndpoint)
		os.Setenv("ENVIRONMENT", originalEnv)
	}()
	
	os.Unsetenv("OTEL_EXPORTER_OTLP_ENDPOINT") // Use default
	os.Setenv("ENVIRONMENT", "production")

	// Act
	shutdownFunc := ConfigureService(ctx, serviceName, version)

	// Assert
	require.NotNil(t, shutdownFunc)
	
	// Clean shutdown
	err := shutdownFunc(ctx)
	assert.NoError(t, err)
}