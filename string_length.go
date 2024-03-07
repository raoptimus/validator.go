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
	skipError      bool
}

func NewStringLength(min, max int) *StringLength {
	return &StringLength{
		message:         "This value must be a string.",
		tooShortMessage: "This value should contain at least {min}.",
		tooLongMessage:  "This value should contain at most {max}.",
		min:             min,
		max:             max,
	}
}

func (r *StringLength) WithMessage(message string) *StringLength {
	rc := *r
	rc.message = message

	return &rc
}

func (r *StringLength) WithTooShortMessage(message string) *StringLength {
	rc := *r
	rc.tooShortMessage = message

	return &rc
}

func (r *StringLength) WithTooLongMessage(message string) *StringLength {
	rc := *r
	rc.tooLongMessage = message

	return &rc
}

func (r *StringLength) When(v WhenFunc) *StringLength {
	rc := *r
	rc.whenFunc = v

	return &rc
}

func (r *StringLength) when() WhenFunc {
	return r.whenFunc
}

func (r *StringLength) setWhen(v WhenFunc) {
	r.whenFunc = v
}

func (r *StringLength) SkipOnEmpty() *StringLength {
	rc := *r
	rc.skipEmpty = true

	return &rc
}

func (r *StringLength) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *StringLength) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *StringLength) SkipOnError() *StringLength {
	rs := *r
	rs.skipError = true

	return &rs
}

func (r *StringLength) shouldSkipOnError() bool {
	return r.skipError
}
func (r *StringLength) setSkipOnError(v bool) {
	r.skipError = v
}

func (r *StringLength) ValidateValue(_ context.Context, value any) error {
	v, ok := toString(value)
	if !ok {
		return NewResult().WithError(NewValidationError(r.message))
	}

	result := NewResult()
	v = strings.Trim(v, " ")
	l := utf8.RuneCountInString(v)

	if l < r.min {
		result = NewResult().
			WithError(
				NewValidationError(r.tooShortMessage).
					WithParams(map[string]any{
						"min": r.min,
						"max": r.max,
					}),
			)
	}

	if l > r.max {
		result = NewResult().
			WithError(
				NewValidationError(r.tooLongMessage).
					WithParams(map[string]any{
						"min": r.min,
						"max": r.max,
					}),
			)
	}

	if !result.IsValid() {
		return result
	}
	return nil
}
