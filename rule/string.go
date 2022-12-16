package rule

import (
	"reflect"
	"strings"
)

type StringLength struct {
	// string user-defined error message used when the value is not a string.
	message string
	// string user-defined error message used when the length of the value is smaller than {see min}.
	tooShortMessage string
	// string user-defined error message used when the length of the value is greater than {see max}.
	tooLongMessage string
	min, max       int
}

func NewStringLength(min, max int) StringLength {
	return StringLength{
		message:         "This value must be a string.",
		tooShortMessage: "This value should contain at least {min}.",
		tooLongMessage:  "This value should contain at most {max}.",
		min:             min,
		max:             max,
	}
}

func (s StringLength) WithMessage(message string) StringLength {
	s.message = message
	return s
}

func (s StringLength) WithTooShortMessage(message string) StringLength {
	s.tooShortMessage = message
	return s
}

func (s StringLength) WithTooLongMessage(message string) StringLength {
	s.tooLongMessage = message
	return s
}

func (s StringLength) ValidateValue(value reflect.Value) error {
	if !value.IsValid() || value.Kind() != reflect.String {
		return NewResult().WithError(formatMessage(s.message))
	}

	result := NewResult()
	v := strings.Trim(value.String(), " ")
	l := len(v)

	if l < s.min {
		result = result.WithError(formatMessageWithArgs(s.tooShortMessage,
			map[string]any{
				"min": s.min,
				"max": s.max,
			}),
		)
	}

	if l > s.max {
		result = result.WithError(formatMessageWithArgs(s.tooLongMessage,
			map[string]any{
				"min": s.min,
				"max": s.max,
			}),
		)
	}

	if !result.IsValid() {
		return result
	}

	return nil
}
