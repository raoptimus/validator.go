package rule

import (
	"reflect"
)

type Validator interface {
	ValidateValue(value any) error
}

type Each struct {
	message string
	rules   []Validator
}

func NewEach(rules ...Validator) Each {
	return Each{
		message: "Value is invalid",
		rules:   rules,
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

		for _, r := range e.rules {
			// todo: use like validator.Validate
			if err := r.ValidateValue(v); err != nil {
				result = result.WithError(err.Error())
			}
		}
	}

	if result.IsValid() {
		return nil
	}

	return result
}
