package telemetry

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestShouldStoreAndRetrieveCorrelationIDFromContext(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedID := uuid.New()

	// Act
	ctxWithID := WithCorrelationID(ctx, expectedID)
	retrievedID, found := CorrelationIDFromContext(ctxWithID)

	// Assert
	assert.True(t, found)
	assert.Equal(t, expectedID, retrievedID)
}

func TestShouldReturnFalseWhenCorrelationIDNotInContext(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	retrievedID, found := CorrelationIDFromContext(ctx)

	// Assert
	assert.False(t, found)
	assert.Equal(t, uuid.UUID{}, retrievedID)
}

func TestShouldOverwriteExistingCorrelationID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	firstID := uuid.New()
	secondID := uuid.New()

	// Act
	ctxWithFirst := WithCorrelationID(ctx, firstID)
	ctxWithSecond := WithCorrelationID(ctxWithFirst, secondID)
	retrievedID, found := CorrelationIDFromContext(ctxWithSecond)

	// Assert
	assert.True(t, found)
	assert.Equal(t, secondID, retrievedID)
	assert.NotEqual(t, firstID, retrievedID)
}

func TestShouldStoreAndRetrieveCausationIDFromContext(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedID := uuid.New()

	// Act
	ctxWithID := WithCausationID(ctx, expectedID)
	retrievedID, found := CausationIDFromContext(ctxWithID)

	// Assert
	assert.True(t, found)
	assert.Equal(t, expectedID, retrievedID)
}

func TestShouldReturnFalseWhenCausationIDNotInContext(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	retrievedID, found := CausationIDFromContext(ctx)

	// Assert
	assert.False(t, found)
	assert.Equal(t, uuid.UUID{}, retrievedID)
}

func TestShouldOverwriteExistingCausationID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	firstID := uuid.New()
	secondID := uuid.New()

	// Act
	ctxWithFirst := WithCausationID(ctx, firstID)
	ctxWithSecond := WithCausationID(ctxWithFirst, secondID)
	retrievedID, found := CausationIDFromContext(ctxWithSecond)

	// Assert
	assert.True(t, found)
	assert.Equal(t, secondID, retrievedID)
	assert.NotEqual(t, firstID, retrievedID)
}

func TestShouldHandleBothCorrelationAndCausationIDsSimultaneously(t *testing.T) {
	// Arrange
	ctx := context.Background()
	correlationID := uuid.New()
	causationID := uuid.New()

	// Act
	ctxWithBoth := WithCorrelationID(WithCausationID(ctx, causationID), correlationID)
	retrievedCorrelationID, foundCorrelation := CorrelationIDFromContext(ctxWithBoth)
	retrievedCausationID, foundCausation := CausationIDFromContext(ctxWithBoth)

	// Assert
	assert.True(t, foundCorrelation)
	assert.True(t, foundCausation)
	assert.Equal(t, correlationID, retrievedCorrelationID)
	assert.Equal(t, causationID, retrievedCausationID)
}