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

func TestMAC_ValidateValue_ValidMAC_Successfully(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{
			name:  "colon separated",
			value: "00:11:22:33:44:55",
		},
		{
			name:  "dash separated",
			value: "00-11-22-33-44-55",
		},
		{
			name:  "uppercase hex",
			value: "AA:BB:CC:DD:EE:FF",
		},
		{
			name:  "mixed case hex",
			value: "aA:bB:cC:dD:eE:fF",
		},
		{
			name:  "dot separated cisco format",
			value: "0011.2233.4455",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewMAC().ValidateValue(ctx, tt.value)
			assert.NoError(t, err)
		})
	}
}

func TestMAC_ValidateValue_InvalidMAC_Failure(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{
			name:  "random string",
			value: "not-a-mac",
		},
		{
			name:  "too short",
			value: "00:11:22",
		},
		{
			name:  "incomplete octets",
			value: "00:11:22:33:44",
		},
		{
			name:  "empty string",
			value: "",
		},
		{
			name:  "invalid hex characters",
			value: "GG:HH:II:JJ:KK:LL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewMAC().ValidateValue(ctx, tt.value)
			assert.Error(t, err)
		})
	}
}

func TestMAC_ValidateValue_NonStringValue_Failure(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{
			name:  "integer",
			value: 12345,
		},
		{
			name:  "nil",
			value: nil,
		},
		{
			name:  "boolean",
			value: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewMAC().ValidateValue(ctx, tt.value)
			assert.Error(t, err)
		})
	}
}

func TestMAC_ValidateValue_WithMessage_Failure(t *testing.T) {
	ctx := context.Background()
	err := NewMAC().WithMessage("custom mac error").ValidateValue(ctx, "invalid")
	assert.Error(t, err)
	assert.Equal(t, "custom mac error.", err.Error())
}

func TestMAC_ValidateValue_SkipOnEmpty_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, nil, NewMAC().SkipOnEmpty())
	assert.NoError(t, err)
}

func TestMAC_ValidateValue_WhenFalse_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, "invalid", NewMAC().When(func(_ context.Context, _ any) bool {
		return false
	}))
	assert.NoError(t, err)
}

func TestMAC_ValidateValue_WhenTrue_Failure(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, "invalid", NewMAC().When(func(_ context.Context, _ any) bool {
		return true
	}))
	assert.Error(t, err)
}

func TestMAC_SkipOnError_Successfully(t *testing.T) {
	r := NewMAC().SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

func TestMAC_ValidateValue_StringPointer_Successfully(t *testing.T) {
	ctx := context.Background()
	mac := "00:11:22:33:44:55"
	err := NewMAC().ValidateValue(ctx, &mac)
	assert.NoError(t, err)
}
