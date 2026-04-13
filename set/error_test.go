/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package set

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUndefinedFieldError_CreatesInstance(t *testing.T) {
	err := NewUndefinedFieldError("hello", "Field")

	assert.NotNil(t, err)
}

func TestUndefinedFieldError_Error_FormatsCorrectly(t *testing.T) {
	tests := []struct {
		name      string
		dataSet   any
		attribute string
		expected  string
	}{
		{
			name:      "string pointer type",
			dataSet:   "hello",
			attribute: "Name",
			expected:  "undefined field: string.Name",
		},
		{
			name:      "map type",
			dataSet:   map[string]any{"a": 1},
			attribute: "Key",
			expected:  "undefined field: map[string]interface {}.Key",
		},
		{
			name:      "struct type",
			dataSet:   struct{ X int }{X: 1},
			attribute: "Y",
			expected:  "undefined field: struct { X int }.Y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewUndefinedFieldError(tt.dataSet, tt.attribute)

			assert.Equal(t, tt.expected, err.Error())
		})
	}
}

func TestUndefinedFieldError_Unwrap_ReturnsBaseError(t *testing.T) {
	err := NewUndefinedFieldError("data", "Field")

	unwrapped := err.Unwrap()

	assert.Equal(t, ErrUndefinedField, unwrapped)
}

func TestUndefinedFieldError_ErrorsIs_MatchesBaseError(t *testing.T) {
	err := NewUndefinedFieldError("data", "Field")

	assert.True(t, errors.Is(err, ErrUndefinedField))
}

func TestUndefinedFieldError_ErrorsIs_DoesNotMatchOtherError(t *testing.T) {
	err := NewUndefinedFieldError("data", "Field")

	assert.False(t, errors.Is(err, errors.New("other error")))
}

func TestUndefinedFieldError_ErrorsAs_MatchesType(t *testing.T) {
	var target *UndefinedFieldError
	err := NewUndefinedFieldError("data", "Field")

	assert.True(t, errors.As(err, &target))
	assert.Equal(t, err, target)
}
