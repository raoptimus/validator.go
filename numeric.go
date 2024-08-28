package validator

import (
	"context"
)

type Numeric struct {
	min               float64
	max               float64
	notNumericMessage string
	tooBigMessage     string
	tooSmallMessage   string
	whenFunc          WhenFunc
	skipEmpty         bool
	skipError         bool
}

func NewNumeric(min, max float64) *Numeric {
	return &Numeric{
		min:               min,
		max:               max,
		notNumericMessage: "Value must be a numeric.",
		tooBigMessage:     "Value must be no greater than {max}.",
		tooSmallMessage:   "Value must be no less than {min}.",
	}
}

func (n *Numeric) WithTooBigMessage(message string) *Numeric {
	rc := *n
	rc.tooBigMessage = message

	return &rc
}

func (n *Numeric) WithTooSmallMessage(message string) *Numeric {
	rc := *n
	rc.tooSmallMessage = message

	return &rc
}

func (n *Numeric) WithNotNumericMessage(message string) *Numeric {
	rc := *n
	rc.notNumericMessage = message

	return &rc
}

func (n *Numeric) When(v WhenFunc) *Numeric {
	rc := *n
	rc.whenFunc = v

	return &rc
}

func (n *Numeric) when() WhenFunc {
	return n.whenFunc
}

func (n *Numeric) setWhen(v WhenFunc) {
	n.whenFunc = v
}

func (n *Numeric) SkipOnEmpty() *Numeric {
	rc := *n
	rc.skipEmpty = true

	return &rc
}

func (n *Numeric) skipOnEmpty() bool {
	return n.skipEmpty
}

func (n *Numeric) setSkipOnEmpty(v bool) {
	n.skipEmpty = v
}

func (n *Numeric) SkipOnError() *Numeric {
	rs := *n
	rs.skipError = true

	return &rs
}

func (n *Numeric) shouldSkipOnError() bool {
	return n.skipError
}
func (n *Numeric) setSkipOnError(v bool) {
	n.skipError = v
}

func (n *Numeric) ValidateValue(_ context.Context, value any) error {
	if value == nil {
		return NewResult().WithError(NewValidationError(n.notNumericMessage))
	}

	var i float64

	switch v := value.(type) {
	case *float32:
		i = float64(*v)
	case *float64:
		i = *v
	case float32:
		i = float64(v)
	case float64:
		i = v
	default:
		return NewResult().WithError(NewValidationError(n.notNumericMessage))
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
