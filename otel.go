package telemetry

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConfigureService(ctx context.Context, serviceName, version string) func(context.Context) error {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4317"
	}

	var (
		exp trace.SpanExporter
		err error
	)

	if strings.Contains(endpoint, "datadog") {
		apiKey := os.Getenv("DD_API_KEY")
		if apiKey == "" {
			slog.Error("DD_API_KEY not set for Datadog exporter")
			return func(context.Context) error { return nil }
		}
		exp, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithHeaders(map[string]string{
				"DD-API-KEY": apiKey,
			}),
			otlptracehttp.WithURLPath("/v1/traces"),
			otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
		)
	} else {
		exp, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		)
	}

	if err != nil {
		slog.Error("failed to create OTLP exporter", "error", err)
		return func(context.Context) error { return nil }
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("service.version", version),
			attribute.String("deployment.environment", os.Getenv("ENVIRONMENT")),
		),
	)
	if err != nil {
		slog.Error("failed to create resource", "error", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	mp := metric.NewMeterProvider()
	otel.SetMeterProvider(mp)

	return tp.Shutdown
}
