/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package regexpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompile_ValidPattern_Successfully(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
	}{
		{
			name:    "simple literal",
			pattern: "abc",
		},
		{
			name:    "digit pattern",
			pattern: `^\d+$`,
		},
		{
			name:    "email-like pattern",
			pattern: `^[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z]{2,}$`,
		},
		{
			name:    "empty pattern",
			pattern: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Compile(tt.pattern)

			require.NoError(t, err)
			assert.NotNil(t, r)
		})
	}
}

func TestCompile_InvalidPattern_Failure(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
	}{
		{
			name:    "unclosed group",
			pattern: "abc(def",
		},
		{
			name:    "invalid repetition",
			pattern: "*abc",
		},
		{
			name:    "unclosed bracket",
			pattern: "[abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Compile(tt.pattern)

			require.Error(t, err)
			assert.Nil(t, r)
		})
	}
}

func TestCompile_SamePatternTwice_ReturnsCachedPointer(t *testing.T) {
	// Use a unique pattern to avoid interference from other tests
	pattern := `^cached_test_unique_\d+$`

	first, err := Compile(pattern)
	require.NoError(t, err)

	second, err := Compile(pattern)
	require.NoError(t, err)

	assert.Same(t, first, second, "second call should return the same pointer from cache")
}

func TestCompile_DifferentPatterns_ReturnsDifferentPointers(t *testing.T) {
	patternA := `^different_test_a_\d+$`
	patternB := `^different_test_b_\d+$`

	a, err := Compile(patternA)
	require.NoError(t, err)

	b, err := Compile(patternB)
	require.NoError(t, err)

	assert.NotSame(t, a, b, "different patterns should return different regexp pointers")
}
