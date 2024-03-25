package validator

import (
	"context"
	"errors"
	"testing"

	"github.com/raoptimus/validator.go/v2/set"
	"github.com/stretchr/testify/assert"
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
	ctx := context.Background()
	dto := &testObject{Name: ""}
	rules := RuleSet{
		"Name": {
			NewRequired(),
		},
	}
	err := Validate(ctx, dto, rules)
	assert.NotNil(t, err)
}

func TestValidatorRequired_EmptyStringWithSpace_ReturnsExpectedError(t *testing.T) {
	ctx := context.Background()
	dto := &testObject{Name: " "}
	rules := RuleSet{
		"Name": {
			NewRequired(),
		},
	}
	err := Validate(ctx, dto, rules)
	assert.NotNil(t, err)
}

func TestValidatorRequired_NilPointerValue_ReturnsExpectedError(t *testing.T) {
	ctx := context.Background()
	dto := &testObject2{Name: nil}
	rules := RuleSet{
		"Name": {
			NewRequired(),
		},
	}
	err := Validate(ctx, dto, rules)
	assert.NotNil(t, err)
}

func TestValidatorRequired_EmptyPointerValue_Successfully(t *testing.T) {
	ctx := context.Background()
	v := ""
	dto := &testObject2{Name: &v}
	rules := RuleSet{
		"Name": {
			NewRequired().WithMessage("Required"),
		},
	}
	err := Validate(ctx, dto, rules)
	assert.NoError(t, err)
}

func TestValidatorRequired_NotEmptyString_ReturnsExpectedNil(t *testing.T) {
	ctx := context.Background()
	dto := &testObject{Name: "test"}
	rules := RuleSet{
		"Name": {
			NewRequired(),
		},
	}
	err := Validate(ctx, dto, rules)
	assert.Nil(t, err)
}

func TestValidatorRequired_NotEmptyPointerValue_ReturnsNil(t *testing.T) {
	ctx := context.Background()
	v := "test"
	dto := &testObject2{Name: &v}
	rules := RuleSet{
		"Name": {
			NewRequired(),
		},
	}
	err := Validate(ctx, dto, rules)
	assert.Nil(t, err)
}

func TestValidatorRequired_NotExistField_ReturnsExpectedError(t *testing.T) {
	ctx := context.Background()
	dto := testObject{Name: ""}
	rules := RuleSet{
		"NotExists": {
			NewRequired().WithMessage("Required"),
		},
	}
	err := Validate(ctx, &dto, rules)
	assert.NotNil(t, err)
	assert.Equal(t, "undefined field: validator.testObject.NotExists", err.Error())
	assert.Equal(t, set.BaseUndefinedFieldError, errors.Unwrap(err))
}

func TestValidatorRequired_NotEmptySlice_ReturnsExpectedNil(t *testing.T) {
	ctx := context.Background()
	dto := &testObject4{Names: []string{"123"}}
	rules := RuleSet{
		"Names": {
			NewRequired(),
		},
	}
	err := Validate(ctx, dto, rules)
	assert.Nil(t, err)
}

func TestValidatorRequired_EmptySlice_ReturnsExpectedError(t *testing.T) {
	ctx := context.Background()
	dto := &testObject4{Names: nil}
	rules := RuleSet{
		"Names": {
			NewRequired(),
		},
	}
	err := Validate(ctx, dto, rules)
	assert.NotNil(t, err)
}
