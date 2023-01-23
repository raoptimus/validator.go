package rule

import (
	"reflect"
)

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

func (r InRange) ValidateValue(value any) error {
	v := reflect.ValueOf(value)
	if !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) {
		return NewResult().WithError(formatMessage(r.message))
	}

	v = reflect.Indirect(v)
	if !v.IsValid() {
		return NewResult().WithError(formatMessage(r.message))
	}

	var in bool
	for _, rv := range r.rangeValues {
		if rv == v.Interface() {
			in = true
			break
		}
	}

	if r.not == in {
		return NewResult().WithError(formatMessage(r.message))
	}

	return nil
}
