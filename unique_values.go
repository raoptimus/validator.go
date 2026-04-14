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
)

type UniqueValues struct {
	message   string
	whenFunc  WhenFunc
	skipEmpty bool
	skipError bool
}

func NewUniqueValues() *UniqueValues {
	return &UniqueValues{
		message: "The list of values must be unique.",
	}
}

func (r *UniqueValues) WithMessage(message string) *UniqueValues {
	rc := *r
	rc.message = message

	return &rc
}

func (r *UniqueValues) When(v WhenFunc) *UniqueValues {
	rc := *r
	rc.whenFunc = v

	return &rc
}

func (r *UniqueValues) when() WhenFunc {
	return r.whenFunc
}

func (r *UniqueValues) setWhen(v WhenFunc) {
	r.whenFunc = v
}

func (r *UniqueValues) SkipOnEmpty() *UniqueValues {
	rc := *r
	rc.skipEmpty = true

	return &rc
}

func (r *UniqueValues) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *UniqueValues) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *UniqueValues) SkipOnError() *UniqueValues {
	rs := *r
	rs.skipError = true

	return &rs
}

func (r *UniqueValues) shouldSkipOnError() bool {
	return r.skipError
}
func (r *UniqueValues) setSkipOnError(v bool) {
	r.skipError = v
}

// ValidateValue checks that all elements in a slice are unique.
// For pointer slices ([]*T), elements are dereferenced before comparison
// so that two distinct pointers with equal values are detected as duplicates.
//
// The element type determines the comparison strategy:
//   - Statically comparable types (primitives, fixed-size structs, pointers):
//     O(n) hash map lookup via map[any]struct{}.
//   - Everything else (structs with slice/map fields, interface slices):
//     O(n) bucketed FNV-64a hash lookup with reflect.DeepEqual collision
//     fallback — see validateHashKey.
//
// Interface-typed slices ([]any, []io.Reader, ...) deliberately skip the
// comparable fast path: reflect.Type.Comparable reports true for interface
// kinds because they are syntactically comparable, but the runtime hash
// panics on non-comparable dynamic values. Routing them through the
// bucketed-hash path handles any dynamic type correctly.
func (r *UniqueValues) ValidateValue(_ context.Context, value any) error {
	// Untyped nil is rejected here (no type info means we cannot treat it
	// as a slice). A typed nil slice (var s []T = nil) has Kind == Slice
	// and reaches the validation loops below with Len == 0, which passes
	// as "trivially unique" — matches the empty-slice behavior.
	if value == nil {
		return NewResult().WithError(NewValidationError(r.message))
	}

	// Type-specialized fast paths skip reflection and the per-element
	// any-boxing that map[any]struct{} forces. Covers the most common
	// validator inputs; everything else falls through to reflect.
	switch s := value.(type) {
	case []string:
		return uniquePrimitive(s, r.message)
	case []int:
		return uniquePrimitive(s, r.message)
	case []int64:
		return uniquePrimitive(s, r.message)
	}

	if reflect.TypeOf(value).Kind() != reflect.Slice {
		return NewResult().WithError(NewValidationError(r.message))
	}

	vs := reflect.ValueOf(value)

	// Determine the actual element type, dereferencing one pointer level
	// to check comparability of the underlying struct, not the pointer.
	elemType := vs.Type().Elem()
	if elemType.Kind() == reflect.Pointer {
		elemType = elemType.Elem()
	}

	if elemType.Kind() != reflect.Interface && elemType.Comparable() {
		return r.validateComparable(vs)
	}

	return r.validateHashKey(vs)
}

// uniquePrimitive is a generic, allocation-minimal duplicate check for
// slices of comparable primitives. Avoids both reflect and the per-element
// any-boxing of validateComparable.
func uniquePrimitive[T comparable](s []T, message string) error {
	set := make(map[T]struct{}, len(s))
	for _, v := range s {
		if _, ok := set[v]; ok {
			return NewResult().WithError(NewValidationError(message))
		}

		set[v] = struct{}{}
	}

	return nil
}

// validateComparable uses a hash map for O(n) duplicate detection.
// Pointer elements are dereferenced so that two *T with equal fields
// produce the same map key (the underlying value, not the address).
func (r *UniqueValues) validateComparable(vs reflect.Value) error {
	n := vs.Len()
	isPtr := vs.Type().Elem().Kind() == reflect.Pointer
	set := make(map[any]struct{}, n)

	for i := 0; i < n; i++ {
		v := vs.Index(i)
		if isPtr && !v.IsNil() {
			v = v.Elem()
		}

		key := v.Interface()
		if _, ok := set[key]; ok {
			return NewResult().WithError(NewValidationError(r.message))
		}

		set[key] = struct{}{}
	}

	return nil
}

// validateHashKey handles slices of non-comparable or interface-typed
// elements using O(n) average-case bucketed hashing. Elements are grouped
// by an FNV-64a hash of their field structure (see hashvalue.go).
//
// The common case — every element produces a unique 64-bit hash — is
// handled by a single map[uint64]int (firstIdx). Only on a hash collision
// is the lazily-allocated overflow map populated, so n unique elements
// trigger n map writes and zero slice allocations.
//
// Because a 64-bit hash cannot prove equality, any collision is verified
// with reflect.DeepEqual: a hash collision between unequal values is not
// a false duplicate, it just degrades that bucket to O(k) probing.
func (r *UniqueValues) validateHashKey(vs reflect.Value) error {
	n := vs.Len()
	isPtr := vs.Type().Elem().Kind() == reflect.Pointer
	firstIdx := make(map[uint64]int, n)
	var overflow map[uint64][]int
	hw := newHasher()

	for i := 0; i < n; i++ {
		v := vs.Index(i)
		if isPtr && !v.IsNil() {
			v = v.Elem()
		}

		hw.reset()
		hashValue(&hw, v)
		key := hw.state

		j, seen := firstIdx[key]
		if !seen {
			firstIdx[key] = i

			continue
		}

		curr := v.Interface()
		if r.equalAt(vs, j, curr, isPtr) {
			return NewResult().WithError(NewValidationError(r.message))
		}

		for _, k := range overflow[key] {
			if r.equalAt(vs, k, curr, isPtr) {
				return NewResult().WithError(NewValidationError(r.message))
			}
		}

		if overflow == nil {
			overflow = make(map[uint64][]int)
		}
		overflow[key] = append(overflow[key], i)
	}

	return nil
}

// equalAt compares the element at index j against curr using DeepEqual,
// dereferencing j's pointer when the slice element type is a pointer.
func (r *UniqueValues) equalAt(vs reflect.Value, j int, curr any, isPtr bool) bool {
	prev := vs.Index(j)
	if isPtr && !prev.IsNil() {
		prev = prev.Elem()
	}

	return reflect.DeepEqual(prev.Interface(), curr)
}
