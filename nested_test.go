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
	"reflect"
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
	TestNestedItem struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	TestNestedSliceObject struct {
		Items []TestNestedItem `json:"items"`
	}
)

func TestNested_ValidateValue(t *testing.T) {
	t.Parallel()

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
			Message:   "This value should contain at least 1.",
			Params:    map[string]any{"max": 255, "min": 1},
			ValuePath: []string{"each", "1"},
		},
		{
			Message:   "Value must be no less than 2.",
			Params:    map[string]any{"max": int64(3), "min": int64(2)},
			ValuePath: []string{"inline", "count"},
		},
	}}

	if !reflect.DeepEqual(expectedError, result) {
		assert.Equal(t, expectedError, result)
	}

	errorMessages := err.(Result).ErrorMessagesIndexedByPath()
	expectedMessages := map[string][]string{
		"inline.count": {"Value must be no less than 2."},
		"each.1":       {"This value should contain at least 1."},
	}

	if !reflect.DeepEqual(expectedMessages, errorMessages) {
		assert.Equal(t, expectedMessages, errorMessages)
	}
}

func TestNested_ValidateValue_InvalidInput_Failure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
	}{
		{
			name:  "nil",
			value: nil,
		},
		{
			name:  "string",
			value: "not a struct",
		},
		{
			name:  "int",
			value: 42,
		},
	}

	nested := NewNested(RuleSet{
		"Name": {NewRequired()},
	})

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := nested.ValidateValue(context.Background(), tt.value)
			assert.Error(t, err)

			var result Result
			assert.ErrorAs(t, err, &result)
		})
	}
}

func TestNested_ValidateValue_BareShortcut_Failure(t *testing.T) {
	t.Parallel()

	nested := NewNested(RuleSet{
		"*": {NewRequired()},
	})

	obj := TestObject{Name: "test"}
	err := nested.ValidateValue(context.Background(), &obj)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "bare shortcut is prohibited")
}

func TestNested_ValidateValue_AllFieldsValid_Successfully(t *testing.T) {
	t.Parallel()

	nested := NewNested(RuleSet{
		"Name": {NewRequired()},
	})

	obj := TestObject{Name: "valid"}
	err := nested.ValidateValue(context.Background(), &obj)
	assert.NoError(t, err)
}

func TestNested_ValidateValue_DotNotationSingleField_Failure(t *testing.T) {
	t.Parallel()

	nested := NewNested(RuleSet{
		"Items.Name": {NewRequired()},
	})

	obj := TestNestedSliceObject{
		Items: []TestNestedItem{
			{Name: "", Value: "ok"},
		},
	}

	err := nested.ValidateValue(context.Background(), &obj)
	assert.Error(t, err)

	var result Result
	assert.ErrorAs(t, err, &result)

	msgs := result.ErrorMessagesIndexedByPath()
	assert.Equal(t, []string{"Value cannot be blank."}, msgs["items.0.name"])
}

func TestNested_ValidateValue_DotNotationSiblings_Successfully(t *testing.T) {
	t.Parallel()

	// Two dot-separated paths sharing the same parent ("Items") must be
	// merged into a single NewEach(NewNested(...)) during normalizeRules.
	// Before the fix, strings.Join joined ALL parts (producing unique keys
	// like "Items*Name"), so sibling rules were never grouped together.
	nested := NewNested(RuleSet{
		"Items.Name":  {NewRequired()},
		"Items.Value": {NewRequired()},
	})

	obj := TestNestedSliceObject{
		Items: []TestNestedItem{
			{Name: "ok", Value: ""},
			{Name: "", Value: "ok"},
		},
	}

	ctx := context.Background()
	err := nested.ValidateValue(ctx, &obj)
	assert.Error(t, err)

	var result Result
	assert.ErrorAs(t, err, &result)

	msgs := result.ErrorMessagesIndexedByPath()
	// Both fields must produce errors — if only one appears, merging failed.
	assert.Contains(t, msgs, "items.0.value")
	assert.Contains(t, msgs, "items.1.name")
	assert.Equal(t, []string{"Value cannot be blank."}, msgs["items.0.value"])
	assert.Equal(t, []string{"Value cannot be blank."}, msgs["items.1.name"])
}
