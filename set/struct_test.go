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

type testStruct struct {
	Name    string  `json:"name"`
	Age     int     `json:"age,omitempty"`
	Address string
	Ptr     *string `json:"ptr"`
}

func TestNewDataSetStruct_WithStructPointer_Successfully(t *testing.T) {
	data := &testStruct{Name: "Alice", Age: 30}

	ds, err := NewDataSetStruct(data)

	require.NoError(t, err)
	assert.NotNil(t, ds)
	assert.Equal(t, NameStruct, ds.Name())
	assert.Equal(t, data, ds.Data())
}

func TestNewDataSetStruct_WithStructValue_Successfully(t *testing.T) {
	data := testStruct{Name: "Bob", Age: 25}

	ds, err := NewDataSetStruct(data)

	require.NoError(t, err)
	assert.NotNil(t, ds)
	assert.Equal(t, NameStruct, ds.Name())
}

func TestNewDataSetStruct_WithNonStruct_Failure(t *testing.T) {
	tests := []struct {
		name string
		data any
	}{
		{
			name: "string value",
			data: "hello",
		},
		{
			name: "int value",
			data: 42,
		},
		{
			name: "string pointer",
			data: func() *string { s := "hello"; return &s }(),
		},
		{
			name: "slice",
			data: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds, err := NewDataSetStruct(tt.data)

			require.Error(t, err)
			assert.Nil(t, ds)
			assert.Contains(t, err.Error(), ErrDataMustBeStructPointer.Error())
		})
	}
}

func TestDataSetStruct_FieldValue_ExistingField_Successfully(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected any
	}{
		{
			name:     "string field",
			field:    "Name",
			expected: "Alice",
		},
		{
			name:     "int field",
			field:    "Age",
			expected: 30,
		},
		{
			name:     "field without json tag",
			field:    "Address",
			expected: "123 Main St",
		},
	}

	data := &testStruct{Name: "Alice", Age: 30, Address: "123 Main St"}
	ds, err := NewDataSetStruct(data)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := ds.FieldValue(tt.field)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, val)
		})
	}
}

func TestDataSetStruct_FieldValue_NonExistingField_Failure(t *testing.T) {
	data := &testStruct{Name: "Alice"}
	ds, err := NewDataSetStruct(data)
	require.NoError(t, err)

	val, err := ds.FieldValue("NonExistent")

	require.Error(t, err)
	assert.Nil(t, val)

	var undefinedErr *UndefinedFieldError
	assert.True(t, errors.As(err, &undefinedErr))
}

func TestDataSetStruct_FieldValue_NilPointerField_ReturnsNil(t *testing.T) {
	data := &testStruct{Name: "Alice", Ptr: nil}
	ds, err := NewDataSetStruct(data)
	require.NoError(t, err)

	val, err := ds.FieldValue("Ptr")

	require.NoError(t, err)
	assert.Nil(t, val)
}

func TestDataSetStruct_FieldValue_NonNilPointerField_Successfully(t *testing.T) {
	s := "hello"
	data := &testStruct{Ptr: &s}
	ds, err := NewDataSetStruct(data)
	require.NoError(t, err)

	val, err := ds.FieldValue("Ptr")

	require.NoError(t, err)
	assert.Equal(t, &s, val)
}

func TestDataSetStruct_FieldAliasName_WithJsonTag_Successfully(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{
			name:     "simple json tag",
			field:    "Name",
			expected: "name",
		},
		{
			name:     "json tag with omitempty",
			field:    "Age",
			expected: "age",
		},
		{
			name:     "pointer field with json tag",
			field:    "Ptr",
			expected: "ptr",
		},
	}

	ds, err := NewDataSetStruct(&testStruct{})
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alias := ds.FieldAliasName(tt.field)

			assert.Equal(t, tt.expected, alias)
		})
	}
}

func TestDataSetStruct_FieldAliasName_WithoutJsonTag_ReturnsFieldName(t *testing.T) {
	ds, err := NewDataSetStruct(&testStruct{})
	require.NoError(t, err)

	alias := ds.FieldAliasName("Address")

	assert.Equal(t, "Address", alias)
}

func TestDataSetStruct_FieldAliasName_NonExistingField_ReturnsInputName(t *testing.T) {
	ds, err := NewDataSetStruct(&testStruct{})
	require.NoError(t, err)

	alias := ds.FieldAliasName("NonExistent")

	assert.Equal(t, "NonExistent", alias)
}

func TestDataSetStruct_Name_ReturnsNameStruct(t *testing.T) {
	ds, err := NewDataSetStruct(&testStruct{})
	require.NoError(t, err)

	assert.Equal(t, NameStruct, ds.Name())
}

func TestDataSetStruct_Data_ReturnsPointer(t *testing.T) {
	data := &testStruct{Name: "test"}
	ds, err := NewDataSetStruct(data)
	require.NoError(t, err)

	assert.Equal(t, data, ds.Data())
}
