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

func (r *UniqueValues) ValidateValue(_ context.Context, value any) error {
	if reflect.TypeOf(value).Kind() != reflect.Slice {
		return NewResult().WithError(NewValidationError(r.message))
	}

	vs := reflect.ValueOf(value)
	set := make(map[any]struct{}, vs.Len())

	for i := 0; i < vs.Len(); i++ {
		v := vs.Index(i).Interface()
		if _, ok := set[v]; ok {
			return NewResult().WithError(NewValidationError(r.message))
		}

		set[v] = struct{}{}
	}

	return nil
}
