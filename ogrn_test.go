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

func TestOGRN_ValidateValue_NonStringValue_Failure(t *testing.T) {
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
			err := NewOGRN().ValidateValue(ctx, tt.value)
			assert.Error(t, err)
		})
	}
}

func TestOGRN_ValidateValue_NonNumericString_Failure(t *testing.T) {
	ctx := context.Background()
	err := NewOGRN().ValidateValue(ctx, "abc1234567890")
	assert.Error(t, err)
	assert.Equal(t, "This value is not a valid OGRN.", err.Error())
}

func TestOGRN_ValidateValue_InvalidLength_Failure(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{name: "too short 12 digits", value: "102770013219"},
		{name: "14 digits", value: "10277001321950"},
		{name: "too long 16 digits", value: "1027700132195012"},
		{name: "empty string", value: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewOGRN().ValidateValue(ctx, tt.value)
			assert.Error(t, err)
		})
	}
}

func TestOGRN_ValidateValue_ValidOGRN_Successfully(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		// 1027700132195: first 12 digits = 102770013219, mod 11 = 5, control = 5
		{name: "valid OGRN starting with 1", value: "1027700132195"},
		// 5077746887312: first 12 = 507774688731, mod 11 = 2, control = 2
		{name: "valid OGRN starting with 5", value: "5077746887312"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewOGRN().ValidateValue(ctx, tt.value)
			assert.NoError(t, err)
		})
	}
}

func TestOGRN_ValidateValue_InvalidFirstDigitOGRN_Failure(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{name: "starts with 2", value: "2027700132195"},
		{name: "starts with 3", value: "3027700132195"},
		{name: "starts with 0", value: "0027700132195"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewOGRN().ValidateValue(ctx, tt.value)
			assert.Error(t, err)
			assert.Equal(t, "This value is not a valid OGRN.", err.Error())
		})
	}
}

func TestOGRN_ValidateValue_InvalidControlDigitOGRN_Failure(t *testing.T) {
	ctx := context.Background()
	// Valid is 1027700132195 (control=5), change last digit to 6
	err := NewOGRN().ValidateValue(ctx, "1027700132196")
	assert.Error(t, err)
	assert.Equal(t, "This value is not a valid OGRN.", err.Error())
}

func TestOGRN_ValidateValue_ValidOGRNIP_Successfully(t *testing.T) {
	// 304500116000157: first 14 = 30450011600015, mod 13 = 7, control = 7
	ctx := context.Background()
	err := NewOGRN().ValidateValue(ctx, "304500116000157")
	assert.NoError(t, err)
}

func TestOGRN_ValidateValue_InvalidFirstDigitOGRNIP_Failure(t *testing.T) {
	ctx := context.Background()
	// 15 digits but starts with 1 (not 3)
	err := NewOGRN().ValidateValue(ctx, "104500116000157")
	assert.Error(t, err)
	assert.Equal(t, "This value is not a valid OGRNIP.", err.Error())
}

func TestOGRN_ValidateValue_InvalidControlDigitOGRNIP_Failure(t *testing.T) {
	ctx := context.Background()
	// Valid is 304500116000157 (control=7), change last digit to 8
	err := NewOGRN().ValidateValue(ctx, "304500116000158")
	assert.Error(t, err)
	assert.Equal(t, "This value is not a valid OGRNIP.", err.Error())
}

func TestOGRN_WithOGRNMessage_CustomMessage(t *testing.T) {
	ctx := context.Background()
	err := NewOGRN().
		WithOGRNMessage("custom ogrn error").
		ValidateValue(ctx, "1027700132196") // invalid control digit
	assert.Error(t, err)
	assert.Equal(t, "custom ogrn error.", err.Error())
}

func TestOGRN_WithOGRNIPMessage_CustomMessage(t *testing.T) {
	ctx := context.Background()
	err := NewOGRN().
		WithOGRNIPMessage("custom ogrnip error").
		ValidateValue(ctx, "304500116000158") // invalid control digit
	assert.Error(t, err)
	assert.Equal(t, "custom ogrnip error.", err.Error())
}

func TestOGRN_ValidateValue_SkipOnEmpty_EmptyString_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, "", NewOGRN().SkipOnEmpty())
	assert.NoError(t, err)
}

func TestOGRN_ValidateValue_SkipOnEmpty_NilValue_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, nil, NewOGRN().SkipOnEmpty())
	assert.NoError(t, err)
}

func TestOGRN_ValidateValue_WhenReturnsFalse_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, "invalid",
		NewOGRN().When(func(_ context.Context, _ any) bool {
			return false
		}),
	)
	assert.NoError(t, err)
}

func TestOGRN_ValidateValue_WhenReturnsTrue_Failure(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, "invalid",
		NewOGRN().When(func(_ context.Context, _ any) bool {
			return true
		}),
	)
	assert.Error(t, err)
}

func TestOGRN_ValidateValue_SkipOnError_PreviousErrored_Successfully(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyPreviousRulesErrored, true)
	err := ValidateValue(ctx, "invalid", NewOGRN().SkipOnError())
	assert.NoError(t, err)
}

func TestOGRN_ValidateValue_SkipOnError_NoPreviousError_Failure(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, "invalid", NewOGRN().SkipOnError())
	assert.Error(t, err)
}

func TestOGRN_ValidateValue_ErrTypeOfResult(t *testing.T) {
	ctx := context.Background()
	err := NewOGRN().ValidateValue(ctx, "1027700132196")
	assert.Error(t, err)
	assert.IsType(t, Result{}, err)
}

func TestOGRN_ValidateValue_InvalidLengthMessage(t *testing.T) {
	ctx := context.Background()
	// 14 digits - not 13 or 15
	err := NewOGRN().ValidateValue(ctx, "10277001321950")
	assert.Error(t, err)
	assert.Equal(t, "This value should contain either 13 or 15 characters.", err.Error())
}
