package validator

import (
	"context"
	"reflect"
)

type UniqueValues struct {
	message string
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
