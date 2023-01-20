package rule

import (
	"reflect"
)

type Validator interface {
	ValidateValue(value any) error
}

type Each struct {
	message   string
	validator Validator
}

func NewEach(r Validator) Each {
	return Each{
		message:   "Values is invalid",
		validator: r,
	}
}

func (e Each) WithMessage(message string) Each {
	e.message = message
	return e
}

func (e Each) ValidateValue(value any) error {
	if reflect.TypeOf(value).Kind() != reflect.Slice {
		return NewResult().WithError(formatMessage(e.message))
	}

	result := NewResult()

	vs := reflect.ValueOf(value)
	for i := 0; i < vs.Len(); i++ {
		v := vs.Index(i).Interface()
		if err := e.validator.ValidateValue(v); err != nil {
			result = result.WithError(err.Error())
		}
	}

	if result.IsValid() {
		return nil
	}

	return result
}
