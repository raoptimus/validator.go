/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDataSetAny_CreatesInstance(t *testing.T) {
	ds := NewDataSetAny("hello")

	assert.NotNil(t, ds)
}

func TestDataSetAny_FieldValue_AlwaysReturnsError(t *testing.T) {
	tests := []struct {
		name  string
		field string
	}{
		{
			name:  "any field name",
			field: "Name",
		},
		{
			name:  "empty field name",
			field: "",
		},
	}

	ds := NewDataSetAny(42)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := ds.FieldValue(tt.field)

			require.Error(t, err)
			assert.Nil(t, val)
			assert.Equal(t, "not supported", err.Error())
		})
	}
}

func TestDataSetAny_FieldAliasName_ReturnsInputName(t *testing.T) {
	ds := NewDataSetAny("data")

	assert.Equal(t, "myField", ds.FieldAliasName("myField"))
	assert.Equal(t, "", ds.FieldAliasName(""))
}

func TestDataSetAny_Name_ReturnsNameAny(t *testing.T) {
	ds := NewDataSetAny("data")

	assert.Equal(t, NameAny, ds.Name())
}

func TestDataSetAny_Data_ReturnsOriginalData(t *testing.T) {
	tests := []struct {
		name string
		data any
	}{
		{
			name: "string data",
			data: "hello",
		},
		{
			name: "int data",
			data: 42,
		},
		{
			name: "nil data",
			data: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewDataSetAny(tt.data)

			assert.Equal(t, tt.data, ds.Data())
		})
	}
}
