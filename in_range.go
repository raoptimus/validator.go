package validator

import "context"

type InRange struct {
	message     string
	rangeValues []any
	not         bool
	whenFunc    WhenFunc
	skipEmpty   bool
}

func NewInRange(rangeValues []any) *InRange {
	return &InRange{
		message:     "This value is invalid",
		rangeValues: rangeValues,
		not:         false,
	}
}

func (r *InRange) WithMessage(message string) *InRange {
	rc := *r
	rc.message = message

	return &rc
}

func (r *InRange) Not() *InRange {
	rc := *r
	rc.not = true

	return &rc
}

func (r *InRange) When(v WhenFunc) *InRange {
	rc := *r
	rc.whenFunc = v

	return &rc
}

func (r *InRange) when() WhenFunc {
	return r.whenFunc
}

func (r *InRange) setWhen(v WhenFunc) {
	r.whenFunc = v
}

func (r *InRange) SkipOnEmpty(v bool) *InRange {
	rc := *r
	rc.skipEmpty = v

	return &rc
}

func (r *InRange) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *InRange) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *InRange) ValidateValue(_ context.Context, value any) error {
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
