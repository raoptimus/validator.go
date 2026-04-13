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
	"github.com/stretchr/testify/require"
)

func TestNewDataSetMap_CreatesInstance(t *testing.T) {
	data := map[string]any{"key": "value"}

	ds := NewDataSetMap(data)

	assert.NotNil(t, ds)
}

func TestDataSetMap_FieldValue_ExistingKey_Successfully(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected any
	}{
		{
			name:     "string value",
			key:      "name",
			expected: "Alice",
		},
		{
			name:     "int value",
			key:      "age",
			expected: 30,
		},
		{
			name:     "nil value",
			key:      "nothing",
			expected: nil,
		},
	}

	data := map[string]any{
		"name":    "Alice",
		"age":     30,
		"nothing": nil,
	}
	ds := NewDataSetMap(data)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := ds.FieldValue(tt.key)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, val)
		})
	}
}

func TestDataSetMap_FieldValue_MissingKey_Failure(t *testing.T) {
	ds := NewDataSetMap(map[string]any{"a": 1})

	val, err := ds.FieldValue("missing")

	require.Error(t, err)
	assert.Nil(t, val)

	var undefinedErr *UndefinedFieldError
	assert.True(t, errors.As(err, &undefinedErr))
}

func TestDataSetMap_FieldValue_EmptyMap_Failure(t *testing.T) {
	ds := NewDataSetMap(map[string]any{})

	val, err := ds.FieldValue("any")

	require.Error(t, err)
	assert.Nil(t, val)
}

func TestDataSetMap_FieldAliasName_ReturnsInputName(t *testing.T) {
	ds := NewDataSetMap(map[string]any{})

	assert.Equal(t, "myField", ds.FieldAliasName("myField"))
	assert.Equal(t, "", ds.FieldAliasName(""))
}

func TestDataSetMap_Name_ReturnsNameMap(t *testing.T) {
	ds := NewDataSetMap(map[string]any{})

	assert.Equal(t, NameMap, ds.Name())
}

func TestDataSetMap_Data_ReturnsOriginalMap(t *testing.T) {
	data := map[string]any{"key": "value"}
	ds := NewDataSetMap(data)

	result := ds.Data()

	assert.Equal(t, data, result)
}
