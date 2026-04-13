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

func TestEmail_ValidateValue_ValidEmail_Successfully(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "simple email",
			value: "user@example.com",
		},
		{
			name:  "email with plus tag",
			value: "test.user+tag@domain.co",
		},
		{
			name:  "email with dots in local part",
			value: "first.last@example.org",
		},
		{
			name:  "email with numbers",
			value: "user123@test456.com",
		},
		{
			name:  "email with underscore",
			value: "user_name@example.com",
		},
		{
			name:  "email with percent",
			value: "user%name@example.com",
		},
		{
			name:  "email with hyphen in domain",
			value: "user@my-domain.com",
		},
		{
			name:  "email with subdomain",
			value: "user@sub.domain.com",
		},
		{
			name:  "single char local part",
			value: "a@example.com",
		},
		{
			name:  "single char tld",
			value: "user@example.a",
		},
		{
			name:  "hyphen in local part",
			value: "user-name@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewEmail().ValidateValue(ctx, tt.value)
			assert.NoError(t, err)
		})
	}
}

func TestEmail_ValidateValue_InvalidEmail_Failure(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "missing local part",
			value: "@domain.com",
		},
		{
			name:  "missing domain",
			value: "user@",
		},
		{
			name:  "plaintext no at sign",
			value: "plaintext",
		},
		{
			name:  "domain starts with dot",
			value: "user@.com",
		},
		{
			name:  "empty string",
			value: "",
		},
		{
			name:  "double at sign",
			value: "user@@example.com",
		},
		{
			name:  "space in address",
			value: "user @example.com",
		},
		{
			name:  "missing tld",
			value: "user@domain",
		},
		{
			name:  "uppercase letters",
			value: "USER@EXAMPLE.COM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewEmail().ValidateValue(ctx, tt.value)
			assert.Error(t, err)
		})
	}
}

func TestEmail_ValidateValue_NonStringValue_Failure(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{
			name:  "integer",
			value: 42,
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
			err := NewEmail().ValidateValue(ctx, tt.value)
			assert.Error(t, err)
		})
	}
}

func TestEmail_ValidateValue_ErrorMessage(t *testing.T) {
	ctx := context.Background()
	err := NewEmail().ValidateValue(ctx, "invalid")
	assert.Error(t, err)
	assert.Equal(t, "Email is not a valid email.", err.Error())
}

func TestEmail_ValidateValue_SkipOnEmpty_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, nil, NewEmail().SkipOnEmpty())
	assert.NoError(t, err)
}

func TestEmail_ValidateValue_WhenFalse_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, "invalid", NewEmail().When(func(_ context.Context, _ any) bool {
		return false
	}))
	assert.NoError(t, err)
}
