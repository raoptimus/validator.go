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

func TestMSISDN_ValidateValue_ValidMSISDN_Successfully(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "full phone number digits",
			value: "79001234567",
		},
		{
			name:  "short number",
			value: "1234",
		},
		{
			name:  "single digit",
			value: "0",
		},
		{
			name:  "long digit string",
			value: "123456789012345",
		},
		{
			name:  "all zeros",
			value: "0000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewMSISDN().ValidateValue(ctx, tt.value)
			assert.NoError(t, err)
		})
	}
}

func TestMSISDN_ValidateValue_InvalidMSISDN_Failure(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "alphabetic characters",
			value: "abc",
		},
		{
			name:  "plus prefix",
			value: "+79001234567",
		},
		{
			name:  "dashes in number",
			value: "123-456",
		},
		{
			name:  "spaces in number",
			value: "123 456",
		},
		{
			name:  "empty string",
			value: "",
		},
		{
			name:  "mixed digits and letters",
			value: "123abc",
		},
		{
			name:  "special characters",
			value: "12345!@#",
		},
		{
			name:  "parentheses format",
			value: "(123)4567890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewMSISDN().ValidateValue(ctx, tt.value)
			assert.Error(t, err)
		})
	}
}

func TestMSISDN_ValidateValue_NonStringValue_Failure(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{
			name:  "integer",
			value: 79001234567,
		},
		{
			name:  "nil",
			value: nil,
		},
		{
			name:  "float",
			value: 123.456,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewMSISDN().ValidateValue(ctx, tt.value)
			assert.Error(t, err)
		})
	}
}

func TestMSISDN_ValidateValue_ErrorMessage(t *testing.T) {
	ctx := context.Background()
	err := NewMSISDN().ValidateValue(ctx, "abc")
	assert.Error(t, err)
	assert.Equal(t, "MSISDN format is invalid.", err.Error())
}

func TestMSISDN_ValidateValue_SkipOnEmpty_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, nil, NewMSISDN().SkipOnEmpty())
	assert.NoError(t, err)
}

func TestMSISDN_ValidateValue_WhenFalse_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, "invalid", NewMSISDN().When(func(_ context.Context, _ any) bool {
		return false
	}))
	assert.NoError(t, err)
}
