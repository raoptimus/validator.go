/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidationError_CreatesWithMessage(t *testing.T) {
	err := NewValidationError("field is required")
	assert.NotNil(t, err)
	assert.Equal(t, "field is required", err.Message)
	assert.Nil(t, err.Params)
	assert.Nil(t, err.ValuePath)
}

func TestValidationError_Error_ReturnsMessage(t *testing.T) {
	err := NewValidationError("value is invalid")
	assert.Equal(t, "value is invalid", err.Error())
}

func TestValidationError_Error_EmptyMessage(t *testing.T) {
	err := NewValidationError("")
	assert.Equal(t, "", err.Error())
}

func TestValidationError_WithParams_SetsParams(t *testing.T) {
	params := map[string]any{"min": 1, "max": 10}
	err := NewValidationError("out of range").WithParams(params)
	assert.Equal(t, params, err.Params)
	assert.Equal(t, "out of range", err.Message)
}

func TestValidationError_WithParams_ReturnsSameInstance(t *testing.T) {
	err := NewValidationError("test")
	result := err.WithParams(map[string]any{"key": "val"})
	assert.Same(t, err, result)
}

func TestValidationError_WithValuePath_SetsPath(t *testing.T) {
	path := []string{"user", "address", "city"}
	err := NewValidationError("cannot be blank").WithValuePath(path)
	assert.Equal(t, path, err.ValuePath)
	assert.Equal(t, "cannot be blank", err.Message)
}

func TestValidationError_WithValuePath_ReturnsSameInstance(t *testing.T) {
	err := NewValidationError("test")
	result := err.WithValuePath([]string{"field"})
	assert.Same(t, err, result)
}

func TestValidationError_Chaining_ParamsAndValuePath(t *testing.T) {
	err := NewValidationError("invalid").
		WithParams(map[string]any{"limit": 5}).
		WithValuePath([]string{"items", "0", "name"})
	assert.Equal(t, "invalid", err.Message)
	assert.Equal(t, map[string]any{"limit": 5}, err.Params)
	assert.Equal(t, []string{"items", "0", "name"}, err.ValuePath)
}

func TestIsError_WithResultError_ReturnsIndexedMessages(t *testing.T) {
	result := NewResult().WithError(
		NewValidationError("cannot be blank").WithValuePath([]string{"name"}),
		NewValidationError("is too short").WithValuePath([]string{"name"}),
		NewValidationError("is invalid").WithValuePath([]string{"email"}),
	)
	msgs, ok := IsError(result)
	assert.True(t, ok)
	assert.Equal(t, map[string][]string{
		"name":  {"cannot be blank", "is too short"},
		"email": {"is invalid"},
	}, msgs)
}

func TestIsError_WithResultNoErrors_ReturnsEmptyMap(t *testing.T) {
	result := NewResult()
	msgs, ok := IsError(result)
	assert.True(t, ok)
	assert.Empty(t, msgs)
}

func TestIsError_WithNonResultError_ReturnsFalse(t *testing.T) {
	err := errors.New("some generic error")
	msgs, ok := IsError(err)
	assert.False(t, ok)
	assert.Nil(t, msgs)
}

func TestIsError_WithNilError_ReturnsFalse(t *testing.T) {
	msgs, ok := IsError(nil)
	assert.False(t, ok)
	assert.Nil(t, msgs)
}

func TestIsError_WithWrappedResultError_ReturnsTrue(t *testing.T) {
	result := NewResult().WithError(
		NewValidationError("required").WithValuePath([]string{"field"}),
	)
	wrapped := errors.Join(errors.New("wrapper"), result)
	msgs, ok := IsError(wrapped)
	assert.True(t, ok)
	assert.Equal(t, map[string][]string{
		"field": {"required"},
	}, msgs)
}
