package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEach_ValidateValue_FirstValueIs1_NoError(t *testing.T) {
	ctx := context.Background()
	err := NewEach(NewNumber(1, 2)).ValidateValue(ctx, []int{1})
	assert.NoError(t, err)
}

func TestEach_ValidateValue_FirstValueIs0_Error(t *testing.T) {
	ctx := context.Background()
	err := NewEach(NewNumber(1, 2)).ValidateValue(ctx, []int{0})
	assert.Error(t, err)
	expectedErr := Result{
		errors: []*ValidationError{
			{
				Message: "Value must be no less than 1.",
				Params: map[string]any{
					"max": int64(2),
					"min": int64(1),
				},
				ValuePath: []string{"0"},
			},
		},
	}

	assert.Equal(t, expectedErr, err)
}
