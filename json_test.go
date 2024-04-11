package validator

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type stringer string

func (s stringer) String() string {
	return string(s)
}

func TestJSON_ValidateValue_Successfully(t *testing.T) {
	// region data provider
	ctx := context.Background()

	str := `"hello world"`
	obj := `{"hello":"world"}`
	arr := `["hello","world"]`

	bytesWithString := []byte(str)
	bytesWithObject := []byte(obj)
	bytesWithArray := []byte(arr)

	jsonRawWithString := json.RawMessage(str)
	jsonRawWithObject := json.RawMessage(obj)
	jsonRawWithArray := json.RawMessage(arr)

	type args struct {
		value any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "stringer with string",
			args: args{value: stringer(str)},
		},
		{
			name: "string",
			args: args{value: str},
		},
		{
			name: "string object",
			args: args{value: obj},
		},
		{
			name: "string array",
			args: args{value: arr},
		},
		{
			name: "string pointer",
			args: args{value: &str},
		},
		{
			name: "string pointer object",
			args: args{value: &obj},
		},
		{
			name: "string pointer array",
			args: args{value: &arr},
		},
		{
			name: "bytes with string",
			args: args{value: bytesWithString},
		},
		{
			name: "bytes with object",
			args: args{value: bytesWithObject},
		},
		{
			name: "bytes with array",
			args: args{value: bytesWithArray},
		},
		{
			name: "bytes pointer with string",
			args: args{value: &bytesWithString},
		},
		{
			name: "bytes with object",
			args: args{value: &bytesWithObject},
		},
		{
			name: "bytes with array",
			args: args{value: &bytesWithArray},
		},
		{
			name: "json.RawMessage with string",
			args: args{value: jsonRawWithString},
		},
		{
			name: "json.RawMessage with object",
			args: args{value: jsonRawWithObject},
		},
		{
			name: "json.RawMessage with array",
			args: args{value: jsonRawWithArray},
		},
		{
			name: "json.RawMessage pointer with string",
			args: args{value: &jsonRawWithString},
		},
		{
			name: "json.RawMessage pointer with object",
			args: args{value: &jsonRawWithObject},
		},
		{
			name: "json.RawMessage pointer with array",
			args: args{value: &jsonRawWithArray},
		},
	}
	// endregion

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := NewJSON()
			err := j.ValidateValue(ctx, tt.args.value)
			require.NoError(t, err)
		})
	}
}

func TestJSON_ValidateValue_Failure(t *testing.T) {
	// region data provider
	ctx := context.Background()

	str := `hello world`
	obj := `{hello":"world"`
	arr := `["hello","world`

	bytesWithString := []byte(str)
	bytesWithObject := []byte(obj)
	bytesWithArray := []byte(arr)

	jsonRawWithString := json.RawMessage(str)
	jsonRawWithObject := json.RawMessage(obj)
	jsonRawWithArray := json.RawMessage(arr)

	type args struct {
		value any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "int",
			args: args{value: 6},
		},
		{
			name: "stringer with string",
			args: args{value: stringer(str)},
		},
		{
			name: "string",
			args: args{value: str},
		},
		{
			name: "string object",
			args: args{value: obj},
		},
		{
			name: "string array",
			args: args{value: arr},
		},
		{
			name: "string pointer",
			args: args{value: &str},
		},
		{
			name: "string pointer object",
			args: args{value: &obj},
		},
		{
			name: "string pointer array",
			args: args{value: &arr},
		},
		{
			name: "bytes with string",
			args: args{value: bytesWithString},
		},
		{
			name: "bytes with object",
			args: args{value: bytesWithObject},
		},
		{
			name: "bytes with array",
			args: args{value: bytesWithArray},
		},
		{
			name: "bytes pointer with string",
			args: args{value: &bytesWithString},
		},
		{
			name: "bytes with object",
			args: args{value: &bytesWithObject},
		},
		{
			name: "bytes with array",
			args: args{value: &bytesWithArray},
		},
		{
			name: "json.RawMessage with string",
			args: args{value: jsonRawWithString},
		},
		{
			name: "json.RawMessage with object",
			args: args{value: jsonRawWithObject},
		},
		{
			name: "json.RawMessage with array",
			args: args{value: jsonRawWithArray},
		},
		{
			name: "json.RawMessage pointer with string",
			args: args{value: &jsonRawWithString},
		},
		{
			name: "json.RawMessage pointer with object",
			args: args{value: &jsonRawWithObject},
		},
		{
			name: "json.RawMessage pointer with array",
			args: args{value: &jsonRawWithArray},
		},
	}
	// endregion

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := NewJSON()
			err := j.ValidateValue(ctx, tt.args.value)
			require.Error(t, err)
		})
	}
}
