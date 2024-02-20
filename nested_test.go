package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	TestInlineObject struct {
		Count int `json:"count"`
	}
	TestObject struct {
		Name   string           `json:"name"`
		Inline TestInlineObject `json:"inline"`
		Each   []string         `json:"each"`
	}
)

func TestNested_ValidateValue(t *testing.T) {
	rules := RuleSet{
		"Name": {
			NewRequired(),
		},
		"Inline": {
			NewNested(
				RuleSet{
					"Count": {
						NewRequired(),
						NewNumber(2, 3),
					},
				},
			),
		},
		"Each": {
			NewEach(NewStringLength(1, 255)),
		},
	}

	obj := TestObject{
		Name: "test",
		Inline: TestInlineObject{
			Count: 1,
		},
		Each: []string{"test", ""},
	}

	ctx := context.Background()
	err := Validate(ctx, &obj, rules)
	assert.Error(t, err)

	var result Result
	assert.ErrorAs(t, err, &result)

	expectedError := Result{errors: []*ValidationError{
		{
			Message:   "Value must be no less than 2.",
			Params:    map[string]any{"max": int64(3), "min": int64(2)},
			ValuePath: []string{"inline", "count"},
		},
		{
			Message:   "This value should contain at least 1.",
			Params:    map[string]any{"max": 255, "min": 1},
			ValuePath: []string{"each", "1"},
		},
	}}

	assert.Equal(t, expectedError, result)

	errorMessages := err.(Result).ErrorMessagesIndexedByPath()
	expectedMessages := map[string][]string{
		"inline.count": {"Value must be no less than 2."},
		"each.1":       {"This value should contain at least 1."},
	}
	assert.Equal(t, expectedMessages, errorMessages)
}
