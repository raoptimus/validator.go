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

func TestIP_ValidateValue_NonStringValue_Failure(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{name: "int value", value: 123},
		{name: "nil value", value: nil},
		{name: "bool value", value: true},
		{name: "float value", value: 1.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewIP().ValidateValue(ctx, tt.value)
			assert.Error(t, err)
			assert.Equal(t, "Must be a valid IP address.", err.Error())
		})
	}
}

func TestIP_ValidateValue_InvalidIPString_Failure(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{name: "random string", value: "not-an-ip"},
		{name: "incomplete IP", value: "192.168.1"},
		{name: "out of range octet", value: "256.1.1.1"},
		{name: "with port", value: "192.168.1.1:8080"},
		{name: "extra dots", value: "192.168.1.1.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewIP().ValidateValue(ctx, tt.value)
			assert.Error(t, err)
			assert.Equal(t, "Must be a valid IP address.", err.Error())
		})
	}
}

func TestIP_ValidateValue_EmptyString_Failure(t *testing.T) {
	ctx := context.Background()
	err := NewIP().ValidateValue(ctx, "")
	assert.Error(t, err)
	assert.Equal(t, "Must be a valid IP address.", err.Error())
}

func TestIP_ValidateValue_ValidIPv4_Successfully(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{name: "localhost", value: "127.0.0.1"},
		{name: "private network", value: "192.168.1.1"},
		{name: "all zeros", value: "0.0.0.0"},
		{name: "broadcast", value: "255.255.255.255"},
		{name: "class A", value: "10.0.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewIP().ValidateValue(ctx, tt.value)
			assert.NoError(t, err)
		})
	}
}

func TestIP_ValidateValue_ValidIPv6_Successfully(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{name: "loopback", value: "::1"},
		{name: "full address", value: "2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
		{name: "abbreviated", value: "2001:db8::1"},
		{name: "all zeros", value: "::"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewIP().ValidateValue(ctx, tt.value)
			assert.NoError(t, err)
		})
	}
}

func TestIP_WithMessage_CustomMessage(t *testing.T) {
	ctx := context.Background()
	err := NewIP().
		WithMessage("custom ip error").
		ValidateValue(ctx, "invalid")
	assert.Error(t, err)
	assert.Equal(t, "custom ip error.", err.Error())
}

func TestIP_ValidateValue_ErrTypeOfResult(t *testing.T) {
	ctx := context.Background()
	err := NewIP().ValidateValue(ctx, "invalid")
	assert.Error(t, err)
	assert.IsType(t, Result{}, err)
}

func TestIP_ValidateValue_SkipOnEmpty_EmptyString_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, "", NewIP().SkipOnEmpty())
	assert.NoError(t, err)
}

func TestIP_ValidateValue_SkipOnEmpty_NilValue_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, nil, NewIP().SkipOnEmpty())
	assert.NoError(t, err)
}

func TestIP_ValidateValue_WhenReturnsFalse_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, "invalid",
		NewIP().When(func(_ context.Context, _ any) bool {
			return false
		}),
	)
	assert.NoError(t, err)
}

func TestIP_ValidateValue_WhenReturnsTrue_Failure(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, "invalid",
		NewIP().When(func(_ context.Context, _ any) bool {
			return true
		}),
	)
	assert.Error(t, err)
}

func TestIP_ValidateValue_SkipOnError_PreviousErrored_Successfully(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyPreviousRulesErrored, true)
	err := ValidateValue(ctx, "invalid", NewIP().SkipOnError())
	assert.NoError(t, err)
}

func TestIP_ValidateValue_SkipOnError_NoPreviousError_Failure(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, "invalid", NewIP().SkipOnError())
	assert.Error(t, err)
}

func TestIP_ValidateValue_StringPointer_Successfully(t *testing.T) {
	ctx := context.Background()
	v := "192.168.1.1"
	err := NewIP().ValidateValue(ctx, &v)
	assert.NoError(t, err)
}
