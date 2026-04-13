/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumeric_ValidateValue_NilValue_Failure(t *testing.T) {
	ctx := context.Background()
	err := NewNumeric(0, 10).ValidateValue(ctx, nil)
	assert.Error(t, err)
	assert.Equal(t, "Value must be a numeric.", err.Error())
}

func TestNumeric_ValidateValue_NonFloatType_Failure(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{name: "string value", value: "abc"},
		{name: "int value", value: 42},
		{name: "bool value", value: true},
		{name: "slice value", value: []float64{1.0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewNumeric(0, 10).ValidateValue(ctx, tt.value)
			assert.Error(t, err)
			assert.Equal(t, "Value must be a numeric.", err.Error())
		})
	}
}

func TestNumeric_ValidateValue_Float32_Successfully(t *testing.T) {
	ctx := context.Background()
	err := NewNumeric(0, 10).ValidateValue(ctx, float32(5.0))
	assert.NoError(t, err)
}

func TestNumeric_ValidateValue_Float64_Successfully(t *testing.T) {
	ctx := context.Background()
	err := NewNumeric(0, 10).ValidateValue(ctx, float64(5.0))
	assert.NoError(t, err)
}

func TestNumeric_ValidateValue_PointerFloat32_Successfully(t *testing.T) {
	ctx := context.Background()
	v := float32(5.0)
	err := NewNumeric(0, 10).ValidateValue(ctx, &v)
	assert.NoError(t, err)
}

func TestNumeric_ValidateValue_PointerFloat64_Successfully(t *testing.T) {
	ctx := context.Background()
	v := float64(5.0)
	err := NewNumeric(0, 10).ValidateValue(ctx, &v)
	assert.NoError(t, err)
}

func TestNumeric_ValidateValue_ValueAtMinBoundary_Successfully(t *testing.T) {
	ctx := context.Background()
	err := NewNumeric(1.0, 10.0).ValidateValue(ctx, float64(1.0))
	assert.NoError(t, err)
}

func TestNumeric_ValidateValue_ValueAtMaxBoundary_Successfully(t *testing.T) {
	ctx := context.Background()
	err := NewNumeric(1.0, 10.0).ValidateValue(ctx, float64(10.0))
	assert.NoError(t, err)
}

func TestNumeric_ValidateValue_ValueBelowMin_Failure(t *testing.T) {
	ctx := context.Background()
	err := NewNumeric(1.0, 10.0).ValidateValue(ctx, float64(0.5))
	assert.Error(t, err)
	assert.Equal(t, "Value must be no less than {min}.", err.Error())
}

func TestNumeric_ValidateValue_ValueAboveMax_Failure(t *testing.T) {
	ctx := context.Background()
	err := NewNumeric(1.0, 10.0).ValidateValue(ctx, float64(10.5))
	assert.Error(t, err)
	assert.Equal(t, "Value must be no greater than {max}.", err.Error())
}

func TestNumeric_ValidateValue_ZeroMinMax_Successfully(t *testing.T) {
	ctx := context.Background()
	err := NewNumeric(0, 0).ValidateValue(ctx, float64(0))
	assert.NoError(t, err)
}

func TestNumeric_ValidateValue_NegativeRange_Successfully(t *testing.T) {
	ctx := context.Background()
	err := NewNumeric(-10.0, -1.0).ValidateValue(ctx, float64(-5.0))
	assert.NoError(t, err)
}

func TestNumeric_ValidateValue_NegativeRange_BelowMin_Failure(t *testing.T) {
	ctx := context.Background()
	err := NewNumeric(-10.0, -1.0).ValidateValue(ctx, float64(-11.0))
	assert.Error(t, err)
}

func TestNumeric_ValidateValue_ErrTypeOfResult(t *testing.T) {
	ctx := context.Background()
	err := NewNumeric(1, 10).ValidateValue(ctx, float64(11))
	assert.Error(t, err)
	assert.IsType(t, Result{}, err)
}

func TestNumeric_WithNotNumericMessage_CustomMessage(t *testing.T) {
	ctx := context.Background()
	err := NewNumeric(0, 10).
		WithNotNumericMessage("custom not numeric").
		ValidateValue(ctx, "abc")
	assert.Error(t, err)
	assert.Equal(t, "custom not numeric.", err.Error())
}

func TestNumeric_WithTooBigMessage_CustomMessage(t *testing.T) {
	ctx := context.Background()
	err := NewNumeric(0, 10).
		WithTooBigMessage("too big value").
		ValidateValue(ctx, float64(11))
	assert.Error(t, err)
	assert.Equal(t, "too big value.", err.Error())
}

func TestNumeric_WithTooSmallMessage_CustomMessage(t *testing.T) {
	ctx := context.Background()
	err := NewNumeric(5, 10).
		WithTooSmallMessage("too small value").
		ValidateValue(ctx, float64(1))
	assert.Error(t, err)
	assert.Equal(t, "too small value.", err.Error())
}

func TestNumeric_ValidateValue_SkipOnEmpty_NilValue_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, nil, NewNumeric(0, 10).SkipOnEmpty())
	assert.NoError(t, err)
}

func TestNumeric_ValidateValue_WhenReturnsFalse_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, float64(100),
		NewNumeric(0, 10).When(func(_ context.Context, _ any) bool {
			return false
		}),
	)
	assert.NoError(t, err)
}

func TestNumeric_ValidateValue_WhenReturnsTrue_Failure(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, float64(100),
		NewNumeric(0, 10).When(func(_ context.Context, _ any) bool {
			return true
		}),
	)
	assert.Error(t, err)
}

func TestNumeric_ValidateValue_SkipOnError_PreviousErrored_Successfully(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyPreviousRulesErrored, true)
	err := ValidateValue(ctx, float64(100), NewNumeric(0, 10).SkipOnError())
	assert.NoError(t, err)
}

func TestNumeric_ValidateValue_SkipOnError_NoPreviousError_Failure(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, float64(100), NewNumeric(0, 10).SkipOnError())
	assert.Error(t, err)
}

func TestNumeric_ValidateValue_PointerFloat32BelowMin_Failure(t *testing.T) {
	ctx := context.Background()
	v := float32(0.5)
	err := NewNumeric(1.0, 10.0).ValidateValue(ctx, &v)
	assert.Error(t, err)
}

func TestNumeric_ValidateValue_PointerFloat64AboveMax_Failure(t *testing.T) {
	ctx := context.Background()
	v := float64(20.0)
	err := NewNumeric(1.0, 10.0).ValidateValue(ctx, &v)
	assert.Error(t, err)
}
