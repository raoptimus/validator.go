package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEach_ValidateValue_HasExpectedValue_NoError(t *testing.T) {
	ctx := context.Background()
	err := NewInRange([]any{"one", "two"}).ValidateValue(ctx, "two")
	assert.NoError(t, err)
}

func TestEach_ValidateValue_NotHasValue_NoError(t *testing.T) {
	ctx := context.Background()
	err := NewInRange([]any{"two"}).ValidateValue(ctx, "one")
	assert.Error(t, err)
}

func TestEach_ValidateValue_ValueIsNil_ReturnsError(t *testing.T) {
	ctx := context.Background()
	err := NewInRange([]any{"two"}).ValidateValue(ctx, nil)
	assert.Error(t, err)
}
