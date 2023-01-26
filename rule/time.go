package rule

import (
	"time"

	"github.com/raoptimus/validator.go/ctype"
)

type Time struct {
	message         string
	formatMessage   string
	tooBigMessage   string
	tooSmallMessage string
	format          string
	min             *time.Time
	max             *time.Time
}

func NewTime() Time {
	return Time{
		message:         "Value is invalid",
		formatMessage:   "Format of the time value must be equal {format}",
		tooBigMessage:   "Time must be no greater than {max}.",
		tooSmallMessage: "Time must be no less than {min}.",
		format:          time.RFC3339,
		min:             nil,
		max:             nil,
	}
}

func (t Time) WithMessage(message string) Time {
	t.message = message
	return t
}

func (t Time) WithFormatMessage(message string) Time {
	t.formatMessage = message
	return t
}

func (t Time) WithTooSmallMessage(message string) Time {
	t.tooSmallMessage = message
	return t
}

func (t Time) WithTooBigMessage(message string) Time {
	t.tooBigMessage = message
	return t
}

func (t Time) WithFormat(format string) Time {
	t.format = format
	return t
}

func (t Time) WithMin(min time.Time) Time {
	t.min = &min
	return t
}

func (t Time) WithMax(max time.Time) Time {
	t.max = &max
	return t
}

func (t Time) ValidateValue(value any) error {
	v, valid := indirectValue(value)
	if !valid {
		return NewResult().WithError(formatMessage(t.message))
	}

	vStr, okStr := v.(string)
	vObj, okObj := v.(ctype.Time)
	if !okStr && !okObj {
		return NewResult().WithError(formatMessage(t.message))
	}

	if okObj {
		vStr = vObj.String()
	}

	vt, err := time.Parse(t.format, vStr)
	if err != nil {
		return NewResult().WithError(
			formatMessageWithArgs(
				t.formatMessage,
				map[string]any{
					"format": t.format,
				},
			),
		)
	}

	result := NewResult()

	if t.min != nil && vt.Before(*t.min) {
		result = result.WithError(
			formatMessageWithArgs(
				t.tooSmallMessage,
				map[string]any{
					"min": t.min,
				},
			),
		)
	}

	if t.max != nil && vt.After(*t.max) {
		result = result.WithError(
			formatMessageWithArgs(
				t.tooBigMessage,
				map[string]any{
					"max": t.max,
				},
			),
		)
	}

	if result.IsValid() {
		if okObj {
			*vObj.Time() = vt
		}

		return nil
	}

	return result
}
