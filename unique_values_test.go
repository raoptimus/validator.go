package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueValues_ValidateValue_HasExpectedValue_NoError(t *testing.T) {
	ctx := context.Background()
	err := NewUniqueValues().ValidateValue(ctx, []string{"one", "two"})
	assert.NoError(t, err)

	err = NewUniqueValues().ValidateValue(ctx, []int{1, 2})
	assert.NoError(t, err)
}

func TestUniqueValues_ValidateValue_HasExpectedValue_ReturnsError(t *testing.T) {
	ctx := context.Background()
	err := NewUniqueValues().ValidateValue(ctx, []string{"two", "two"})
	assert.Error(t, err)

	err = NewUniqueValues().ValidateValue(ctx, []int{1, 1})
	assert.Error(t, err)
}
