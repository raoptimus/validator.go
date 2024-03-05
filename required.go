package validator

import (
	"context"
	"reflect"
)

type Required struct {
	message        string
	allowZeroValue bool
	whenFunc       WhenFunc
}

func NewRequired() Required {
	return Required{
		message:        "Value cannot be blank.",
		allowZeroValue: false,
	}
}

func (r Required) When(v WhenFunc) Required {
	r.whenFunc = v

	return r
}

func (r Required) when() WhenFunc {
	return r.whenFunc
}

func (s Required) WithMessage(message string) Required {
	s.message = message
	return s
}

func (r Required) WithAllowZeroValue() Required {
	r.allowZeroValue = true

	return r
}

func (r Required) ValidateValue(_ context.Context, value any) error {
	v := reflect.ValueOf(value)
	if valueIsEmpty(v, r.allowZeroValue) {
		return NewResult().WithError(NewValidationError(r.message))
	}

	return nil
}
