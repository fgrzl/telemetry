package telemetry

import (
	"context"
	"encoding/binary"
	"log/slog"

	// OpenTelemetry

	"go.opentelemetry.io/otel/trace"
)

func CorrelateLogsAndTraces(ctx context.Context, log *slog.Logger) *slog.Logger {

	if log == nil {
		log = slog.Default()
	}

	// Get the trace/span IDs from the context:
	traceID := trace.SpanContextFromContext(ctx).TraceID()
	spanID := trace.SpanContextFromContext(ctx).SpanID()

	// DataDog expects 64-bit IDs in decimal, so you might do something like:
	ddTraceID := binary.BigEndian.Uint64(traceID[8:16]) // last 8 bytes
	ddSpanId := binary.BigEndian.Uint64(spanID[:])

	log.With(
		"trace_id", traceID.String(),
		"span_id", spanID.String(),
		"dd.trace_id", ddTraceID,
		"dd.span_id", ddSpanId,
	)

	return log
}
