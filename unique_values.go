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
	if value == nil || reflect.TypeOf(value).Kind() != reflect.Slice {
		return NewResult().WithError(NewValidationError(r.message))
	}

	vs := reflect.ValueOf(value)

	// Determine the actual element type, dereferencing one pointer level
	// to check comparability of the underlying struct, not the pointer.
	elemType := vs.Type().Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	if elemType.Kind() != reflect.Interface && elemType.Comparable() {
		return r.validateComparable(vs)
	}

	return r.validateHashKey(vs)
}

// validateComparable uses a hash map for O(n) duplicate detection.
// Pointer elements are dereferenced so that two *T with equal fields
// produce the same map key (the underlying value, not the address).
func (r *UniqueValues) validateComparable(vs reflect.Value) error {
	n := vs.Len()
	isPtr := vs.Type().Elem().Kind() == reflect.Ptr
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
// into buckets by an FNV-64a hash of their field structure (see
// hashvalue.go). Because a 64-bit hash cannot prove equality, any bucket
// with more than one member is verified with reflect.DeepEqual — a hash
// collision between unequal values is not a false duplicate, it just
// degrades that bucket to O(k) DeepEqual probing per insert.
func (r *UniqueValues) validateHashKey(vs reflect.Value) error {
	n := vs.Len()
	isPtr := vs.Type().Elem().Kind() == reflect.Ptr
	buckets := make(map[uint64][]int, n)
	hw := newHasher()

	for i := 0; i < n; i++ {
		v := vs.Index(i)
		if isPtr && !v.IsNil() {
			v = v.Elem()
		}

		hw.reset()
		hashValue(&hw, v)
		key := hw.state

		if indices, ok := buckets[key]; ok {
			curr := v.Interface()
			for _, j := range indices {
				prev := vs.Index(j)
				if isPtr && !prev.IsNil() {
					prev = prev.Elem()
				}

				if reflect.DeepEqual(prev.Interface(), curr) {
					return NewResult().WithError(NewValidationError(r.message))
				}
			}
		}

		buckets[key] = append(buckets[key], i)
	}

	return nil
}
