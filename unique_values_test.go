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

type testNonComparableItem struct {
	Name string
	Tags []string // slice field makes the struct non-comparable
}

func TestUniqueValues_ValidateValue_Successfully(t *testing.T) {
	ctx := context.Background()
	v := func(s string) *string { return &s }

	tests := []struct {
		name  string
		value any
	}{
		{
			name:  "unique strings",
			value: []string{"one", "two"},
		},
		{
			name:  "unique ints",
			value: []int{1, 2},
		},
		{
			name:  "empty slice",
			value: []string{},
		},
		{
			name:  "single element",
			value: []string{"one"},
		},
		{
			name:  "unique pointers to same value type",
			value: []*string{v("a"), v("b")},
		},
		{
			name:  "pointer with nil element",
			value: []*string{v("a"), nil},
		},
		{
			name:  "unique non-comparable structs",
			value: []testNonComparableItem{{Name: "a", Tags: []string{"x"}}, {Name: "b", Tags: []string{"x"}}},
		},
		{
			name:  "unique non-comparable struct pointers",
			value: []*testNonComparableItem{{Name: "a"}, {Name: "b"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewUniqueValues().ValidateValue(ctx, tt.value)
			assert.NoError(t, err)
		})
	}
}

func TestUniqueValues_ValidateValue_Failure(t *testing.T) {
	ctx := context.Background()
	v := func(s string) *string { return &s }

	tests := []struct {
		name  string
		value any
	}{
		{
			name:  "duplicate strings",
			value: []string{"two", "two"},
		},
		{
			name:  "duplicate ints",
			value: []int{1, 1},
		},
		{
			name:  "nil input",
			value: nil,
		},
		{
			name:  "non-slice string",
			value: "not a slice",
		},
		{
			name:  "non-slice int",
			value: 42,
		},
		{
			name:  "duplicate pointers by value",
			value: []*string{v("same"), v("same")},
		},
		{
			name:  "duplicate non-comparable structs",
			value: []testNonComparableItem{{Name: "a", Tags: []string{"x"}}, {Name: "a", Tags: []string{"x"}}},
		},
		{
			name:  "duplicate non-comparable struct pointers",
			value: []*testNonComparableItem{{Name: "a", Tags: []string{"x"}}, {Name: "a", Tags: []string{"x"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewUniqueValues().ValidateValue(ctx, tt.value)
			assert.Error(t, err)
		})
	}
}
