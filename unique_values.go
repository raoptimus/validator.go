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
//   - Comparable types (primitives, fixed-size structs): O(n) hash map lookup.
//   - Non-comparable types (e.g. protobuf structs with slice fields):
//     O(n²) pairwise reflect.DeepEqual — avoids allocations from string
//     serialization and is faster for typical validation arrays (< 30 elements).
func (r *UniqueValues) ValidateValue(_ context.Context, value any) error {
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

	if elemType.Comparable() {
		return r.validateComparable(vs)
	}

	return r.validateDeepEqual(vs)
}

// validateComparable uses a hash map for O(n) duplicate detection.
// Pointer elements are dereferenced so that two *T with equal fields
// produce the same map key (the underlying value, not the address).
func (r *UniqueValues) validateComparable(vs reflect.Value) error {
	set := make(map[any]struct{}, vs.Len())

	for i := 0; i < vs.Len(); i++ {
		v := vs.Index(i)
		key := v.Interface()

		if v.Kind() == reflect.Ptr && !v.IsNil() {
			key = v.Elem().Interface()
		}

		if _, ok := set[key]; ok {
			return NewResult().WithError(NewValidationError(r.message))
		}

		set[key] = struct{}{}
	}

	return nil
}

// validateDeepEqual handles slices of non-comparable types (structs with
// slice/map fields, e.g. protobuf messages). Uses O(n²) pairwise comparison
// with reflect.DeepEqual instead of string serialization — this is cheaper
// for small slices because DeepEqual short-circuits on first field mismatch
// and allocates nothing per comparison.
func (r *UniqueValues) validateDeepEqual(vs reflect.Value) error {
	seen := make([]any, 0, vs.Len())

	for i := 0; i < vs.Len(); i++ {
		v := vs.Index(i)
		curr := v.Interface()

		if v.Kind() == reflect.Ptr && !v.IsNil() {
			curr = v.Elem().Interface()
		}

		for _, prev := range seen {
			if reflect.DeepEqual(prev, curr) {
				return NewResult().WithError(NewValidationError(r.message))
			}
		}

		seen = append(seen, curr)
	}

	return nil
}
