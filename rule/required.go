package rule

import (
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

func (s Required) ValidateValue(value any) error {
	v := reflect.ValueOf(value)
	if valueIsEmpty(v, s.allowZeroValue) {
		return NewResult().WithError(formatMessage(s.message))
	}

	return nil
}
