package telemetry

import (
	"context"

	"github.com/google/uuid"
)

// correlationKeyType is the context key for correlation ID.
type correlationKeyType struct{}

// causationKeyType is the context key for causation ID.
type causationKeyType struct{}

// WithCorrelationID returns a new context with the given correlation ID.
func WithCorrelationID(ctx context.Context, correlationID uuid.UUID) context.Context {
	return context.WithValue(ctx, correlationKeyType{}, correlationID)
}

// CorrelationIDFromContext extracts the correlation ID from the context.
func CorrelationIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(correlationKeyType{}).(uuid.UUID)
	return id, ok
}

// WithCausationID returns a new context with the given causation ID.
func WithCausationID(ctx context.Context, causationID uuid.UUID) context.Context {
	return context.WithValue(ctx, causationKeyType{}, causationID)
}

// CausationIDFromContext extracts the causation ID from the context.
func CausationIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(causationKeyType{}).(uuid.UUID)
	return id, ok
}
