package validator

import "context"

type InRange struct {
	message     string
	rangeValues []any
	not         bool
}

func NewInRange(rangeValues []any) InRange {
	return InRange{
		message:     "This value is invalid",
		rangeValues: rangeValues,
		not:         false,
	}
}

func (r InRange) WithMessage(message string) InRange {
	r.message = message
	return r
}

func (r InRange) Not() InRange {
	r.not = true
	return r
}

func (r InRange) ValidateValue(_ context.Context, value any) error {
	v, valid := indirectValue(value)
	if !valid {
		return NewResult().WithError(NewValidationError(r.message))
	}

	var in bool
	for _, rv := range r.rangeValues {
		if rv == v {
			in = true
			break
		}
	}

	if r.not == in {
		return NewResult().WithError(NewValidationError(r.message))
	}

	return nil
}
