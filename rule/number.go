package rule

import (
	"reflect"
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

func (n Number) ValidateValue(value any) error {
	var i int64

	switch reflect.TypeOf(value).Kind() {
	case reflect.Int8:
		i = int64(value.(int8))
	case reflect.Uint8:
		i = int64(value.(uint8))
	case reflect.Int:
		i = int64(value.(int))
	case reflect.Uint:
		i = int64(value.(uint))
	case reflect.Int16:
		i = int64(value.(int16))
	case reflect.Uint16:
		i = int64(value.(uint16))
	case reflect.Int32:
		i = int64(value.(int32))
	case reflect.Uint32:
		i = int64(value.(uint32))
	case reflect.Int64:
		i = value.(int64)
	case reflect.Uint64:
		i = int64(value.(uint64))
	default:
		return NewResult().WithError(n.notNumberMessage)
	}

	result := NewResult()

	if i < n.min {
		result = result.WithError(formatMessageWithArgs(n.tooSmallMessage,
			map[string]any{
				"min": n.min,
				"max": n.max,
			}),
		)
	}

	if i > n.max {
		result = result.WithError(formatMessageWithArgs(n.tooBigMessage,
			map[string]any{
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
