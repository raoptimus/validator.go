package rule

import (
	"reflect"
)

type Required struct {
	message string
}

func NewRequired() *Required {
	return &Required{
		message: "Value cannot be blank.",
	}
}

func (s *Required) WithMessage(message string) *Required {
	s.message = message
	return s
}

func (s *Required) ValidateValue(value reflect.Value) *Result {
	if valueIsEmpty(value) {
		return NewResult().WithError(formatMessage(s.message))
	}

	return nil
}
