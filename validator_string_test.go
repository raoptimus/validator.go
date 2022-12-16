package validator

import (
	"testing"

	"github.com/raoptimus/validator.go/rule"
	"github.com/stretchr/testify/assert"
)

func TestValidatorString_EmptyStringMin0_ReturnsNil(t *testing.T) {
	dto := &testObject{Name: ""}
	rules := map[string][]RuleValidator{
		"Name": {
			rule.NewStringLength(0, 0),
		},
	}
	err := Validate(dto, rules, false)
	assert.Nil(t, err)
}

func TestValidatorString_EmptyStringMin1_ReturnsError(t *testing.T) {
	dto := &testObject{Name: ""}
	rules := map[string][]RuleValidator{
		"Name": {
			rule.NewStringLength(1, 0),
		},
	}
	err := Validate(dto, rules, false)
	assert.NotNil(t, err)
}
