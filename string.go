package validator

import (
	"context"
	"strings"
	"unicode/utf8"
)

type StringLength struct {
	// string user-defined error message used when the value is not a string.
	message string
	// string user-defined error message used when the length of the value is smaller than {see min}.
	tooShortMessage string
	// string user-defined error message used when the length of the value is greater than {see max}.
	tooLongMessage string
	min, max       int
	whenFunc       WhenFunc
	skipEmpty      bool
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

func (s StringLength) When(v WhenFunc) StringLength {
	s.whenFunc = v

	return s
}

func (s StringLength) when() WhenFunc {
	return s.whenFunc
}

func (s StringLength) SkipOnEmpty(v bool) StringLength {
	s.skipEmpty = v

	return s
}

func (s StringLength) skipOnEmpty() bool {
	return s.skipEmpty
}

func (s StringLength) ValidateValue(_ context.Context, value any) error {
	v, ok := toString(value)
	if !ok {
		return NewResult().WithError(NewValidationError(s.message))
	}

	result := NewResult()
	v = strings.Trim(v, " ")
	l := utf8.RuneCountInString(v)

	if l < s.min {
		result = NewResult().
			WithError(
				NewValidationError(s.tooShortMessage).
					WithParams(map[string]any{
						"min": s.min,
						"max": s.max,
					}),
			)
	}

	if l > s.max {
		result = NewResult().
			WithError(
				NewValidationError(s.tooLongMessage).
					WithParams(map[string]any{
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
