package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrl_ValidateValue_Successfully(t *testing.T) {
	ctx := context.Background()
	r := NewURL()
	err := r.ValidateValue(ctx, "https://example.com")
	assert.NoError(t, err)
}

func TestUrl_ValidateValue_Schema(t *testing.T) {
	ctx := context.Background()
	r := NewURL().WithValidScheme("http", "myphotoapp")

	err := r.ValidateValue(ctx, "http://example.com")
	assert.NoError(t, err)

	err = r.ValidateValue(ctx, "myphotoapp:test?name=123")
	assert.NoError(t, err)

	err = r.ValidateValue(ctx, "https://example.com")
	assert.Error(t, err)
}

func TestUrl_ValidateValue_AnySchema(t *testing.T) {
	ctx := context.Background()
	r := NewURL().WithValidScheme(AllowAnyURLSchema)

	err := r.ValidateValue(ctx, "http://example.com")
	assert.NoError(t, err)

	err = r.ValidateValue(ctx, "myphotoapp:test?name=123")
	assert.NoError(t, err)

	err = r.ValidateValue(ctx, "https://example.com")
	assert.NoError(t, err)
}

func TestUrl_ValidateValue_IDN_Successfully(t *testing.T) {
	ctx := context.Background()
	err := NewURL().WithEnableIDN().ValidateValue(ctx, "https://президент.рф")
	assert.NoError(t, err)
}

func TestUrlValidateValue_EmptyValue_ReturnsError(t *testing.T) {
	ctx := context.Background()
	err := NewURL().ValidateValue(ctx, "")
	assert.Error(t, err)
}

func TestUrlValidateValue_InvalidValue_ReturnsError(t *testing.T) {
	ctx := context.Background()
	r := NewURL()

	err := r.ValidateValue(ctx, "http://")
	assert.Error(t, err)

	err = r.ValidateValue(ctx, "myphotoapp test?name=123")
	assert.Error(t, err)
}

func TestUrlValidateValue_InvalidValue_ReturnsExpectedErrorMessage(t *testing.T) {
	ctx := context.Background()
	err := NewURL().WithMessage("test error").ValidateValue(ctx, "http://")
	assert.Error(t, err)
	assert.Equal(t, "test error.", err.Error())
}
