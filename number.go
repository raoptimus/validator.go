package validator

import (
	"context"
)

type Number struct {
	min              int64
	max              int64
	notNumberMessage string
	tooBigMessage    string
	tooSmallMessage  string
	whenFunc         WhenFunc
	skipEmpty        bool
	skipError        bool
}

func NewNumber(min, max int64) *Number {
	return &Number{
		min:              min,
		max:              max,
		notNumberMessage: "Value must be a number.",
		tooBigMessage:    "Value must be no greater than {max}.",
		tooSmallMessage:  "Value must be no less than {min}.",
	}
}

func (r *Number) WithTooBigMessage(message string) *Number {
	rc := *r
	rc.tooBigMessage = message

	return &rc
}

func (r *Number) WithTooSmallMessage(message string) *Number {
	rc := *r
	rc.tooSmallMessage = message

	return &rc
}

func (r *Number) WithNotNumberMessage(message string) *Number {
	rc := *r
	rc.notNumberMessage = message

	return &rc
}

func (r *Number) When(v WhenFunc) *Number {
	rc := *r
	rc.whenFunc = v

	return &rc
}

func (r *Number) when() WhenFunc {
	return r.whenFunc
}

func (r *Number) setWhen(v WhenFunc) {
	r.whenFunc = v
}

func (r *Number) SkipOnEmpty() *Number {
	rc := *r
	rc.skipEmpty = true

	return &rc
}

func (r *Number) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *Number) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *Number) SkipOnError() *Number {
	rs := *r
	rs.skipError = true

	return &rs
}

func (r *Number) shouldSkipOnError() bool {
	return r.skipError
}
func (r *Number) setSkipOnError(v bool) {
	r.skipError = v
}

func (r *Number) ValidateValue(_ context.Context, value any) error {
	var i int64

	switch v := value.(type) {
	case *int8:
		i = int64(*v)
	case *uint8:
		i = int64(*v)
	case *int:
		i = int64(*v)
	case *uint:
		i = int64(*v)
	case *int16:
		i = int64(*v)
	case *uint16:
		i = int64(*v)
	case *int32:
		i = int64(*v)
	case *uint32:
		i = int64(*v)
	case *int64:
		i = *v
	case *uint64:
		i = int64(*v)
	case int8:
		i = int64(v)
	case uint8:
		i = int64(v)
	case int:
		i = int64(v)
	case uint:
		i = int64(v)
	case int16:
		i = int64(v)
	case uint16:
		i = int64(v)
	case int32:
		i = int64(v)
	case uint32:
		i = int64(v)
	case int64:
		i = v
	case uint64:
		i = int64(v)
	default:
		return NewResult().WithError(NewValidationError(r.notNumberMessage))
	}

	result := NewResult()

	if i < r.min {
		result = result.WithError(
			NewValidationError(r.tooSmallMessage).
				WithParams(map[string]any{
					"min": r.min,
					"max": r.max,
				}),
		)
	}

	if i > r.max {
		result = result.WithError(
			NewValidationError(r.tooBigMessage).
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
