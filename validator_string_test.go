package validator

import (
	"testing"

	"github.com/raoptimus/validator.go/rule"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestValidatorString_EmptyStringMin0_ReturnsNil(t *testing.T) {
	ctx := context.Background()
	dto := &testObject{Name: ""}
	rules := map[string][]RuleValidator{
		"Name": {
			rule.NewStringLength(0, 0),
		},
	}
	err := Validate(ctx, dto, rules, false)
	assert.Nil(t, err)
}

func TestValidatorString_EmptyStringMin1_ReturnsError(t *testing.T) {
	ctx := context.Background()
	dto := &testObject{Name: ""}
	rules := map[string][]RuleValidator{
		"Name": {
			rule.NewStringLength(1, 0),
		},
	}
	err := Validate(ctx, dto, rules, false)
	assert.NotNil(t, err)
}

func TestValidatorString_NotRequiredEmptyString_ReturnsNil(t *testing.T) {
	ctx := context.Background()
	dto := &testObject2{Name: nil}
	rules := map[string][]RuleValidator{
		"Name": {
			rule.NewStringLength(1, 0),
		},
	}
	err := Validate(ctx, dto, rules, false)
	assert.NoError(t, err)
}
