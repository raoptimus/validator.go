package validator

import (
	"context"
	"reflect"
)

type Required struct {
	message        string
	allowZeroValue bool
}

func NewRequired() Required {
	return Required{
		message:        "Value cannot be blank.",
		allowZeroValue: false,
	}
}

func (s Required) WithMessage(message string) Required {
	s.message = message
	return s
}

func (s Required) WithAllowZeroValue() Required {
	s.allowZeroValue = true
	return s
}

func (s Required) ValidateValue(_ context.Context, value any) error {
	v := reflect.ValueOf(value)
	if valueIsEmpty(v, s.allowZeroValue) {
		return NewResult().WithError(NewValidationError(s.message))
	}

	return nil
}
