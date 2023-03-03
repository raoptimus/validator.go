package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raoptimus/validator.go/rule"
)

type testObject struct {
	Name string
}
type testObject2 struct {
	Name *string
}
type testObject3 struct {
	Name *string `json:"name"`
}
type testObject4 struct {
	Names []string
}

func TestValidatorRequired_EmptyString_ReturnsExpectedError(t *testing.T) {
	dto := &testObject{Name: ""}
	rules := map[string][]RuleValidator{
		"Name": {
			rule.NewRequired(),
		},
	}
	err := Validate(dto, rules, false)
	assert.NotNil(t, err)
}

func TestValidatorRequired_EmptyStringWithSpace_ReturnsExpectedError(t *testing.T) {
	dto := &testObject{Name: " "}
	rules := map[string][]RuleValidator{
		"Name": {
			rule.NewRequired(),
		},
	}
	err := Validate(dto, rules, false)
	assert.NotNil(t, err)
}

func TestValidatorRequired_NilPointerValue_ReturnsExpectedError(t *testing.T) {
	dto := &testObject2{Name: nil}
	rules := map[string][]RuleValidator{
		"Name": {
			rule.NewRequired(),
		},
	}
	err := Validate(dto, rules, false)
	assert.NotNil(t, err)
}

func TestValidatorRequired_EmptyPointerValue_ReturnsExpectedError(t *testing.T) {
	v := ""
	dto := &testObject2{Name: &v}
	rules := map[string][]RuleValidator{
		"Name": {
			rule.NewRequired().WithMessage("Required"),
		},
	}
	err := Validate(dto, rules, false)
	assert.NotNil(t, err)
	assert.Equal(t, "Required", err.Error())
	assert.Equal(t, map[string][]string{"Name": {"Required"}}, err.(rule.ResultSet).ResultErrors())
}

func TestValidatorRequired_EmptyPointerValueWithJsonTag_ReturnsExpectedError(t *testing.T) {
	v := ""
	dto := &testObject3{Name: &v}
	rules := map[string][]RuleValidator{
		"Name": {
			rule.NewRequired().WithMessage("Required"),
		},
	}
	err := Validate(dto, rules, false)
	assert.NotNil(t, err)
	assert.Equal(t, "Required", err.Error())
	assert.Equal(t, map[string][]string{"name": {"Required"}}, err.(rule.ResultSet).ResultErrors())
}

func TestValidatorRequired_NotEmptyString_ReturnsExpectedNil(t *testing.T) {
	dto := &testObject{Name: "test"}
	rules := map[string][]RuleValidator{
		"Name": {
			rule.NewRequired(),
		},
	}
	err := Validate(dto, rules, false)
	assert.Nil(t, err)
}

func TestValidatorRequired_NotEmptyPointerValue_ReturnsNil(t *testing.T) {
	v := "test"
	dto := &testObject2{Name: &v}
	rules := map[string][]RuleValidator{
		"Name": {
			rule.NewRequired(),
		},
	}
	err := Validate(dto, rules, false)
	assert.Nil(t, err)
}

func TestValidatorRequired_NotExistProperty_ReturnsExpectedError(t *testing.T) {
	dto := &testObject{Name: ""}
	rules := map[string][]RuleValidator{
		"NotExists": {
			rule.NewRequired().WithMessage("Required"),
		},
	}
	err := Validate(dto, rules, false)
	assert.NotNil(t, err)
	assert.Equal(t, "undefined property: validator.testObject.NotExists", err.Error())
	assert.Equal(t, ErrUndefinedField, errors.Unwrap(err))
}

func TestValidatorRequired_NotEmptySlice_ReturnsExpectedNil(t *testing.T) {
	dto := &testObject4{Names: []string{"123"}}
	rules := map[string][]RuleValidator{
		"Names": {
			rule.NewRequired(),
		},
	}
	err := Validate(dto, rules, false)
	assert.Nil(t, err)
}

func TestValidatorRequired_EmptySlice_ReturnsExpectedError(t *testing.T) {
	dto := &testObject4{Names: nil}
	rules := map[string][]RuleValidator{
		"Names": {
			rule.NewRequired(),
		},
	}
	err := Validate(dto, rules, false)
	assert.NotNil(t, err)
}
