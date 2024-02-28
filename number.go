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
}

func NewNumber(min, max int64) Number {
	return Number{
		min:              min,
		max:              max,
		notNumberMessage: "Value must be a number.",
		tooBigMessage:    "Value must be no greater than {max}.",
		tooSmallMessage:  "Value must be no less than {min}.",
	}
}

func (n Number) WithTooBigMessage(message string) Number {
	n.tooBigMessage = message
	return n
}

func (n Number) WithTooSmallMessage(message string) Number {
	n.tooSmallMessage = message
	return n
}

func (n Number) WithNotNumberMessage(message string) Number {
	n.notNumberMessage = message
	return n
}

func (n Number) ValidateValue(_ context.Context, value any) error {
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
		return NewResult().WithError(NewValidationError(n.notNumberMessage))
	}

	result := NewResult()

	if i < n.min {
		result = result.WithError(
			NewValidationError(n.tooSmallMessage).
				WithParams(map[string]any{
					"min": n.min,
					"max": n.max,
				}),
		)
	}

	if i > n.max {
		result = result.WithError(
			NewValidationError(n.tooBigMessage).
				WithParams(map[string]any{
					"min": n.min,
					"max": n.max,
				}),
		)
	}

	if !result.IsValid() {
		return result
	}
	return nil
}
