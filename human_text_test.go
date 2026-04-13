/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHumanText_ValidateValue_ValidText_Successfully(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "english text",
			value: "Hello World",
		},
		{
			name:  "cyrillic text",
			value: "Привет мир",
		},
		{
			name:  "text with digits",
			value: "Test 123",
		},
		{
			name:  "text with comma and exclamation",
			value: "Hello, World!",
		},
		{
			name:  "text with question mark",
			value: "How are you?",
		},
		{
			name:  "text with colon and semicolon",
			value: "Note: items; details",
		},
		{
			name:  "text with parentheses",
			value: "Hello (World)",
		},
		{
			name:  "text with square brackets",
			value: "Hello [World]",
		},
		{
			name:  "text with single quotes",
			value: "it's a test",
		},
		{
			name:  "text with double quotes",
			value: `say "hello"`,
		},
		{
			name:  "text with guillemets",
			value: "«Привет»",
		},
		{
			name:  "text with dash",
			value: "well-known",
		},
		{
			name:  "text with em dash",
			value: "hello — world",
		},
		{
			name:  "text with en dash",
			value: "hello – world",
		},
		{
			name:  "text with hash",
			value: "#hashtag",
		},
		{
			name:  "text with forward slash",
			value: "yes/no",
		},
		{
			name:  "text with period",
			value: "Hello.",
		},
		{
			name:  "single letter",
			value: "A",
		},
		{
			name:  "single digit",
			value: "7",
		},
		{
			name:  "text with tilde",
			value: "approx~value",
		},
		{
			name:  "text with backtick",
			value: "code" + "`" + "block",
		},
		{
			name:  "text with german low double quote",
			value: "\u201eHallo\u201c",
		},
		{
			name:  "text with curly single quotes",
			value: "\u2018hello\u2019",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewHumanText().ValidateValue(ctx, tt.value)
			assert.NoError(t, err)
		})
	}
}

func TestHumanText_ValidateValue_InvalidText_Failure(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "empty string",
			value: "",
		},
		{
			name:  "tab character",
			value: "Hello\tWorld",
		},
		{
			name:  "newline character",
			value: "Hello\nWorld",
		},
		{
			name:  "carriage return",
			value: "Hello\rWorld",
		},
		{
			name:  "null byte",
			value: "Hello\x00World",
		},
		{
			name:  "dollar sign",
			value: "price $100",
		},
		{
			name:  "percent sign",
			value: "100%",
		},
		{
			name:  "ampersand",
			value: "A & B",
		},
		{
			name:  "asterisk",
			value: "hello*world",
		},
		{
			name:  "plus sign",
			value: "a+b",
		},
		{
			name:  "less than sign",
			value: "a<b",
		},
		{
			name:  "equals sign",
			value: "a=b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewHumanText().ValidateValue(ctx, tt.value)
			assert.Error(t, err)
			assert.Equal(t, "This value must be a normal text.", err.Error())
		})
	}
}

func TestHumanText_ValidateValue_NonStringValue_Failure(t *testing.T) {
	ctx := context.Background()
	err := NewHumanText().ValidateValue(ctx, 123)
	assert.Error(t, err)
	assert.Equal(t, "This value must be a normal text.", err.Error())
}

func TestHumanText_ValidateValue_NilValue_Failure(t *testing.T) {
	ctx := context.Background()
	err := NewHumanText().ValidateValue(ctx, nil)
	assert.Error(t, err)
}

func TestHumanText_ValidateValue_SkipOnEmpty_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, nil, NewHumanText().SkipOnEmpty())
	assert.NoError(t, err)
}

func TestHumanText_ValidateValue_WhenFalse_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, "Hello\tWorld", NewHumanText().When(func(_ context.Context, _ any) bool {
		return false
	}))
	assert.NoError(t, err)
}

func TestHumanText_ValidateValue_WhenTrue_Failure(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, "\t", NewHumanText().When(func(_ context.Context, _ any) bool {
		return true
	}))
	assert.Error(t, err)
}
