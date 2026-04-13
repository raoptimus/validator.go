/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package vtype

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTime_CreatesTimeWithFields_Successfully(t *testing.T) {
	tests := []struct {
		name            string
		unvalidatedTime string
	}{
		{
			name:            "non-empty string",
			unvalidatedTime: "2024-01-15T10:30:00Z",
		},
		{
			name:            "empty string",
			unvalidatedTime: "",
		},
		{
			name:            "date only",
			unvalidatedTime: "2024-01-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewTime(tt.unvalidatedTime)

			assert.NotNil(t, result.validatedTime)
			assert.Equal(t, time.Time{}, *result.validatedTime)
			assert.Equal(t, tt.unvalidatedTime, result.unvalidatedTime)
		})
	}
}

func TestTime_Time_ReturnsValidatedTime_Successfully(t *testing.T) {
	tests := []struct {
		name string
		time Time
	}{
		{
			name: "validatedTime already set via NewTime",
			time: NewTime("2024-01-15T10:30:00Z"),
		},
		{
			name: "validatedTime is nil initializes to zero time",
			time: Time{validatedTime: nil, unvalidatedTime: "2024-01-15"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.time.Time()

			require.NotNil(t, result)
			assert.Equal(t, time.Time{}, *result)
		})
	}
}

func TestTime_Time_CalledTwiceReturnsSamePointer_Successfully(t *testing.T) {
	vt := NewTime("2024-01-15")

	first := vt.Time()
	second := vt.Time()

	assert.Same(t, first, second)
}

func TestTime_String_ReturnsUnvalidatedTime_Successfully(t *testing.T) {
	tests := []struct {
		name     string
		time     Time
		expected string
	}{
		{
			name:     "non-empty string",
			time:     NewTime("2024-01-15T10:30:00Z"),
			expected: "2024-01-15T10:30:00Z",
		},
		{
			name:     "empty string",
			time:     NewTime(""),
			expected: "",
		},
		{
			name:     "arbitrary text",
			time:     NewTime("not-a-date"),
			expected: "not-a-date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.time.String()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTime_UnmarshalJSON_NullData_Successfully(t *testing.T) {
	tests := []struct {
		name                    string
		initialUnvalidatedTime  string
		data                    []byte
		expectedUnvalidatedTime string
	}{
		{
			name:                    "null literal preserves existing state",
			initialUnvalidatedTime:  "2024-01-15",
			data:                    []byte("null"),
			expectedUnvalidatedTime: "2024-01-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vt := NewTime(tt.initialUnvalidatedTime)

			err := vt.UnmarshalJSON(tt.data)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedUnvalidatedTime, vt.unvalidatedTime)
		})
	}
}

func TestTime_UnmarshalJSON_QuotedString_Successfully(t *testing.T) {
	tests := []struct {
		name                    string
		data                    []byte
		expectedUnvalidatedTime string
	}{
		{
			name:                    "quoted datetime string",
			data:                    []byte(`"2024-01-15T10:30:00Z"`),
			expectedUnvalidatedTime: "2024-01-15T10:30:00Z",
		},
		{
			name:                    "quoted date only",
			data:                    []byte(`"2024-01-15"`),
			expectedUnvalidatedTime: "2024-01-15",
		},
		{
			name:                    "quoted empty string",
			data:                    []byte(`""`),
			expectedUnvalidatedTime: "",
		},
		{
			name:                    "unquoted string",
			data:                    []byte("2024-01-15"),
			expectedUnvalidatedTime: "2024-01-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vt := NewTime("")

			err := vt.UnmarshalJSON(tt.data)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedUnvalidatedTime, vt.unvalidatedTime)
			assert.NotNil(t, vt.validatedTime)
		})
	}
}

func TestTime_UnmarshalJSON_ValidatedTimeNil_InitializesIt_Successfully(t *testing.T) {
	vt := Time{validatedTime: nil, unvalidatedTime: ""}

	err := vt.UnmarshalJSON([]byte(`"2024-01-15"`))

	require.NoError(t, err)
	assert.Equal(t, "2024-01-15", vt.unvalidatedTime)
	require.NotNil(t, vt.validatedTime)
	assert.Equal(t, time.Time{}, *vt.validatedTime)
}

func TestTime_UnmarshalJSON_NullWithNilValidatedTime_DoesNotInitialize_Successfully(t *testing.T) {
	vt := Time{validatedTime: nil, unvalidatedTime: ""}

	err := vt.UnmarshalJSON([]byte("null"))

	require.NoError(t, err)
	assert.Nil(t, vt.validatedTime)
}

func TestTime_UnmarshalText_SetsUnvalidatedTime_Successfully(t *testing.T) {
	tests := []struct {
		name                    string
		data                    []byte
		expectedUnvalidatedTime string
	}{
		{
			name:                    "datetime string",
			data:                    []byte("2024-01-15T10:30:00Z"),
			expectedUnvalidatedTime: "2024-01-15T10:30:00Z",
		},
		{
			name:                    "empty bytes",
			data:                    []byte(""),
			expectedUnvalidatedTime: "",
		},
		{
			name:                    "arbitrary text",
			data:                    []byte("not-a-date"),
			expectedUnvalidatedTime: "not-a-date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vt := NewTime("")

			err := vt.UnmarshalText(tt.data)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedUnvalidatedTime, vt.unvalidatedTime)
			assert.NotNil(t, vt.validatedTime)
		})
	}
}

func TestTime_UnmarshalText_ValidatedTimeNil_InitializesIt_Successfully(t *testing.T) {
	vt := Time{validatedTime: nil, unvalidatedTime: ""}

	err := vt.UnmarshalText([]byte("2024-01-15"))

	require.NoError(t, err)
	assert.Equal(t, "2024-01-15", vt.unvalidatedTime)
	require.NotNil(t, vt.validatedTime)
	assert.Equal(t, time.Time{}, *vt.validatedTime)
}

func TestTime_MarshalJSON_ReturnsQuotedString_Successfully(t *testing.T) {
	tests := []struct {
		name     string
		time     Time
		expected []byte
	}{
		{
			name:     "datetime string",
			time:     NewTime("2024-01-15T10:30:00Z"),
			expected: []byte(`"2024-01-15T10:30:00Z"`),
		},
		{
			name:     "empty string produces empty quotes",
			time:     NewTime(""),
			expected: []byte(`""`),
		},
		{
			name:     "date only",
			time:     NewTime("2024-01-15"),
			expected: []byte(`"2024-01-15"`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.time.MarshalJSON()

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTime_MarshalText_ReturnsQuotedString_Successfully(t *testing.T) {
	tests := []struct {
		name     string
		time     Time
		expected []byte
	}{
		{
			name:     "datetime string",
			time:     NewTime("2024-01-15T10:30:00Z"),
			expected: []byte(`"2024-01-15T10:30:00Z"`),
		},
		{
			name:     "empty string produces empty quotes",
			time:     NewTime(""),
			expected: []byte(`""`),
		},
		{
			name:     "arbitrary text",
			time:     NewTime("hello-world"),
			expected: []byte(`"hello-world"`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.time.MarshalText()

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
