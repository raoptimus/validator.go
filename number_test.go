package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumber_ValidateValue_ValueNotNumber(t *testing.T) {
	ctx := context.Background()
	err := NewNumber(0, 0).ValidateValue(ctx, "abc")
	assert.Error(t, err)
	assert.Equal(t, "Value must be a number.", err.Error())
}

func TestNumber_ValidateValue_ValueNotNumberCustomError(t *testing.T) {
	ctx := context.Background()
	err := NewNumber(0, 0).
		WithNotNumberMessage("test").
		ValidateValue(ctx, "abc")
	assert.Error(t, err)
	assert.Equal(t, "test.", err.Error())
}

func TestNumber_ValidateValue_ValueLessThanMin(t *testing.T) {
	ctx := context.Background()
	err := NewNumber(1, 10).ValidateValue(ctx, 0)
	assert.Error(t, err)
	assert.Equal(t, "Value must be no less than {min}.", err.Error())
}

func TestNumber_ValidateValue_ValueLessThanMinCustomError(t *testing.T) {
	ctx := context.Background()
	err := NewNumber(1, 10).WithTooSmallMessage("test {min}").
		ValidateValue(ctx, 0)
	assert.Error(t, err)
	assert.Equal(t, "test {min}.", err.Error())
}

func TestNumber_ValidateValue_ValueGreatThanMax(t *testing.T) {
	ctx := context.Background()
	err := NewNumber(1, 10).ValidateValue(ctx, 11)
	assert.Error(t, err)
	assert.Equal(t, "Value must be no greater than {max}.", err.Error())
}

func TestNumber_ValidateValue_ValueGreatThanMaxCustomError(t *testing.T) {
	ctx := context.Background()
	err := NewNumber(1, 10).
		WithTooBigMessage("test {max}").
		ValidateValue(ctx, 11)
	assert.Error(t, err)
	assert.Equal(t, "test {max}.", err.Error())
}

func TestNumber_ValidateValue_ErrTypeOfResultSet(t *testing.T) {
	ctx := context.Background()
	err := NewNumber(1, 10).ValidateValue(ctx, 11)
	assert.Error(t, err)
	assert.IsType(t, Result{}, err)
}

func TestNumber_ValidateValue_Successfully(t *testing.T) {
	ctx := context.Background()
	err := NewNumber(1, 10).ValidateValue(ctx, 2)
	assert.NoError(t, err)
}

func TestNumber_ValidateValue_SuccessfullyZero(t *testing.T) {
	ctx := context.Background()
	err := NewNumber(0, 0).ValidateValue(ctx, 0)
	assert.NoError(t, err)
}
