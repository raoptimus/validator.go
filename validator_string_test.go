package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatorString_EmptyStringMin0_ReturnsNil(t *testing.T) {
	ctx := context.Background()
	dto := &testObject{Name: ""}
	rules := RuleSet{
		"Name": {
			NewStringLength(0, 0),
		},
	}
	err := Validate(ctx, dto, rules)
	assert.Nil(t, err)
}

func TestValidatorString_EmptyStringMin1_ReturnsError(t *testing.T) {
	ctx := context.Background()
	dto := &testObject{Name: ""}
	rules := RuleSet{
		"Name": {
			NewStringLength(1, 0),
		},
	}
	err := Validate(ctx, dto, rules)
	assert.NotNil(t, err)
}

func TestValidatorString_NotRequiredEmptyString_ReturnsNil(t *testing.T) {
	ctx := context.Background()
	dto := &testObject2{Name: nil}
	rules := RuleSet{
		"Name": {
			NewStringLength(1, 0),
		},
	}
	err := Validate(ctx, dto, rules)
	assert.NoError(t, err)
}
