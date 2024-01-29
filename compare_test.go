package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCompareFields struct {
	Field1 int `json:"field1"`
	Field2 int `json:"field2"`
}

func TestCompare_ValidateValue(t *testing.T) {
	obj := TestCompareFields{
		Field1: 1,
		Field2: 2,
	}

	rules := RuleSet{
		"Field1": {
			NewCompare(nil, "Field2", "=="),
		},
	}

	ctx := context.Background()
	err := Validate(ctx, &obj, rules)
	assert.Error(t, err)

	expectedError := Result{errors: []*ValidationError{
		{
			Message:   "Value must be equal to '2'.",
			Params:    map[string]any{"targetValue": nil, "targetAttribute": "Field2", "targetValueOrAttribute": 2},
			ValuePath: []string{"field1"},
		},
	}}
	assert.Equal(t, expectedError, err)

	errorMessages := err.(Result).ErrorMessagesIndexedByPath()
	expectedMessages := map[string][]string{
		"field1": {"Value must be equal to '2'."},
	}
	assert.Equal(t, expectedMessages, errorMessages)
}
