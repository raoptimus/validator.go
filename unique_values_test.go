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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testNonComparableItem struct {
	Name string
	Tags []string // slice field makes the struct non-comparable
}

func makeUniqueNonComparableSlice(n int) []testNonComparableItem {
	items := make([]testNonComparableItem, 0, n)
	for i := 0; i < n; i++ {
		items = append(items, testNonComparableItem{
			Name: fmt.Sprintf("item-%d", i),
			Tags: []string{"tag"},
		})
	}

	return items
}

func makeDuplicateNonComparableSlice(n int) []testNonComparableItem {
	items := makeUniqueNonComparableSlice(n)
	items[n-1] = items[0]

	return items
}

func TestUniqueValues_ValidateValue_Successfully(t *testing.T) {
	t.Parallel()

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
			// Typed nil slice: Kind == Slice, Len == 0, treated as empty.
			// Diverges from untyped nil (which errors) — see ValidateValue.
			name:  "typed nil slice",
			value: []string(nil),
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
		{
			// Larger-n case to exercise bucketed-hash collision handling
			// at scale; unique inputs should never produce false positives.
			name:  "unique non-comparable structs at scale",
			value: makeUniqueNonComparableSlice(50),
		},
		{
			name: "unique interface slice with non-comparable dynamic types",
			value: []any{
				testNonComparableItem{Name: "a", Tags: []string{"x"}},
				testNonComparableItem{Name: "b", Tags: []string{"y"}},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := NewUniqueValues().ValidateValue(ctx, tt.value)
			assert.NoError(t, err)
		})
	}
}

func TestUniqueValues_ValidateValue_Failure(t *testing.T) {
	t.Parallel()

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
			// Two nil pointers compare equal after the nil-guarded deref:
			// locks in "nil equals nil" semantics in the comparable path
			// (elem is *string → deref to string → routed to validateComparable).
			name:  "duplicate nil pointers (comparable path)",
			value: []*string{nil, nil},
		},
		{
			// Same semantics but routed to validateHashKey because the
			// pointee is a non-comparable struct. Exercises the hash path's
			// nil-guard and DeepEqual fallback on typed-nil pointers.
			name:  "duplicate nil pointers (hash path)",
			value: []*testNonComparableItem{nil, nil},
		},
		{
			name:  "duplicate non-comparable structs",
			value: []testNonComparableItem{{Name: "a", Tags: []string{"x"}}, {Name: "a", Tags: []string{"x"}}},
		},
		{
			name:  "duplicate non-comparable struct pointers",
			value: []*testNonComparableItem{{Name: "a", Tags: []string{"x"}}, {Name: "a", Tags: []string{"x"}}},
		},
		{
			// Larger-n case to confirm duplicates are still detected
			// correctly after the bucketed hash path.
			name:  "duplicate non-comparable structs at scale",
			value: makeDuplicateNonComparableSlice(50),
		},
		{
			name: "duplicate interface slice with non-comparable dynamic types",
			value: []any{
				testNonComparableItem{Name: "a", Tags: []string{"x"}},
				testNonComparableItem{Name: "a", Tags: []string{"x"}},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := NewUniqueValues().ValidateValue(ctx, tt.value)
			assert.Error(t, err)
		})
	}
}

// TestUniqueValues_ValidateValue_CyclicStruct_Successfully proves the
// maxHashDepth cap in hashvalue.go terminates the structural walk on
// self-referential inputs without stack overflow.
func TestUniqueValues_ValidateValue_CyclicStruct_Successfully(t *testing.T) {
	t.Parallel()

	type cyclicNode struct {
		Tags []string // non-comparable → routed to validateHashKey
		Self *cyclicNode
	}

	n := &cyclicNode{Tags: []string{"x"}}
	n.Self = n

	err := NewUniqueValues().ValidateValue(context.Background(), []*cyclicNode{n})
	assert.NoError(t, err)
}
