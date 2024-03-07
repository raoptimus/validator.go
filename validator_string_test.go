package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatorString_ValidRulesForPtr_Successfully(t *testing.T) {
	var empty string
	one := "a"
	two := "ab"

	tests := []struct {
		name       string
		rule       *StringLength
		testObject *testObject2
	}{
		{
			name:       "empty string",
			rule:       NewStringLength(0, 0),
			testObject: &testObject2{Name: &empty},
		},
		{
			name:       "1 len string",
			rule:       NewStringLength(1, 2),
			testObject: &testObject2{Name: &one},
		},
		{
			name:       "2 len string",
			rule:       NewStringLength(1, 2),
			testObject: &testObject2{Name: &two},
		},
		{
			name:       "empty string with skip on empty, min 1",
			rule:       NewStringLength(1, 2).SkipOnEmpty(),
			testObject: &testObject2{Name: &empty},
		},
		{
			name: "not valid string, but not when",
			rule: NewStringLength(1, 1).
				When(func(ctx context.Context, value any) bool {
					return false
				}),
			testObject: &testObject2{Name: &two},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			rules := RuleSet{"Name": {tt.rule}}
			err := Validate(ctx, tt.testObject, rules)
			assert.NoError(t, err)
		})
	}
}

func TestValidatorString_ValidRules_Successfully(t *testing.T) {
	tests := []struct {
		name       string
		rule       *StringLength
		testObject *testObject
	}{
		{
			name:       "empty string",
			rule:       NewStringLength(0, 0),
			testObject: &testObject{Name: ""},
		},
		{
			name:       "1 len string",
			rule:       NewStringLength(1, 2),
			testObject: &testObject{Name: "a"},
		},
		{
			name:       "2 len string",
			rule:       NewStringLength(1, 2),
			testObject: &testObject{Name: "ab"},
		},
		{
			name:       "empty string with skip on empty, min 1",
			rule:       NewStringLength(1, 2).SkipOnEmpty(),
			testObject: &testObject{Name: ""},
		},
		{
			name: "not valid string, but not when",
			rule: NewStringLength(1, 1).
				When(func(ctx context.Context, value any) bool {
					return false
				}),
			testObject: &testObject{Name: "ab"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			rules := RuleSet{"Name": {tt.rule}}
			err := Validate(ctx, tt.testObject, rules)
			assert.NoError(t, err)
		})
	}
}

func TestValidatorString_NotValidRulesForPtr_Failure(t *testing.T) {
	var empty string
	two := "ab"

	tests := []struct {
		name           string
		rule           *StringLength
		testObject     *testObject2
		wantErrMessage string
	}{
		{
			name:           "nil string, min 0",
			rule:           NewStringLength(0, 1),
			testObject:     &testObject2{Name: nil},
			wantErrMessage: "Name: This value must be a string.",
		},
		{
			name:           "empty string, min 1",
			rule:           NewStringLength(1, 2),
			testObject:     &testObject2{Name: &empty},
			wantErrMessage: "Name: This value should contain at least 1.",
		},
		{
			name:           "empty string, min 1",
			rule:           NewStringLength(1, 1),
			testObject:     &testObject2{Name: &two},
			wantErrMessage: "Name: This value should contain at most 1.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			rules := RuleSet{"Name": {tt.rule}}
			err := Validate(ctx, tt.testObject, rules)
			assert.Error(t, err)
			assert.Equal(t, tt.wantErrMessage, err.Error())
		})
	}
}

func TestValidatorString_NotValidRules_Failure(t *testing.T) {
	tests := []struct {
		name           string
		rule           *StringLength
		testObject     *testObject
		wantErrMessage string
	}{
		{
			name:           "empty string, min 1",
			rule:           NewStringLength(1, 2),
			testObject:     &testObject{Name: ""},
			wantErrMessage: "Name: This value should contain at least 1.",
		},
		{
			name:           "empty string, min 1",
			rule:           NewStringLength(1, 1),
			testObject:     &testObject{Name: "ab"},
			wantErrMessage: "Name: This value should contain at most 1.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			rules := RuleSet{"Name": {tt.rule}}
			err := Validate(ctx, tt.testObject, rules)
			assert.Error(t, err)
			assert.Equal(t, tt.wantErrMessage, err.Error())
		})
	}
}
