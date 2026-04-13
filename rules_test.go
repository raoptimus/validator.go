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

// mockRule is a simple Rule that does NOT implement RuleSkipError
type mockRule struct{}

func (m mockRule) ValidateValue(_ context.Context, _ any) error {
	return nil
}

func TestRules_SkipOnError_SetsSkipOnRulesImplementingInterface(t *testing.T) {
	num := NewNumber(0, 10)
	assert.False(t, num.shouldSkipOnError())

	rules := Rules{num}
	rules.SkipOnError()

	assert.True(t, num.shouldSkipOnError())
}

func TestRules_SkipOnError_MultipleRules(t *testing.T) {
	num1 := NewNumber(1, 100)
	num2 := NewNumber(0, 50)
	assert.False(t, num1.shouldSkipOnError())
	assert.False(t, num2.shouldSkipOnError())

	rules := Rules{num1, num2}
	rules.SkipOnError()

	assert.True(t, num1.shouldSkipOnError())
	assert.True(t, num2.shouldSkipOnError())
}

func TestRules_SkipOnError_SkipsRulesNotImplementingInterface(t *testing.T) {
	mock := mockRule{}
	num := NewNumber(0, 10)

	rules := Rules{mock, num}
	rules.SkipOnError()

	// mockRule does not implement RuleSkipError, so it should be unaffected
	// Number should still have skipOnError set
	assert.True(t, num.shouldSkipOnError())
}

func TestRules_SkipOnError_EmptyRules(t *testing.T) {
	rules := Rules{}
	// Should not panic
	rules.SkipOnError()
	assert.Empty(t, rules)
}

func TestRules_SkipOnError_NilRules(t *testing.T) {
	var rules Rules
	// Should not panic on nil slice
	rules.SkipOnError()
	assert.Nil(t, rules)
}

func TestRules_SkipOnError_MixedRuleTypes(t *testing.T) {
	num := NewNumber(0, 10)
	mock := mockRule{}
	req := NewRequired()

	rules := Rules{num, mock, req}
	rules.SkipOnError()

	assert.True(t, num.shouldSkipOnError())
	assert.True(t, req.shouldSkipOnError())
}
