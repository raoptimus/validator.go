package validator

import (
	"context"
	"reflect"
)

type UniqueValues struct {
	message   string
	whenFunc  WhenFunc
	skipEmpty bool
}

func NewUniqueValues() UniqueValues {
	return UniqueValues{
		message: "The list of values must be unique.",
	}
}

func (r UniqueValues) WithMessage(message string) UniqueValues {
	r.message = message

	return r
}

func (r UniqueValues) When(v WhenFunc) UniqueValues {
	r.whenFunc = v

	return r
}

func (r UniqueValues) when() WhenFunc {
	return r.whenFunc
}

func (r UniqueValues) SkipOnEmpty(v bool) UniqueValues {
	r.skipEmpty = v

	return r
}

func (r UniqueValues) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r UniqueValues) ValidateValue(_ context.Context, value any) error {
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
