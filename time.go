package validator

import (
	"context"
	"time"

	"github.com/raoptimus/validator.go/v2/vtype"
)

type TimeFunc func() time.Time

type Time struct {
	message         string
	formatMessage   string
	tooBigMessage   string
	tooSmallMessage string
	format          string
	min             TimeFunc
	max             TimeFunc
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

func (t Time) WithMin(min TimeFunc) Time {
	t.min = min
	return t
}

func (t Time) WithMax(max TimeFunc) Time {
	t.max = max
	return t
}

func (t Time) ValidateValue(_ context.Context, value any) error {
	v, valid := indirectValue(value)
	if !valid {
		return NewResult().WithError(NewValidationError(t.message))
	}

	vStr, okStr := toString(value)
	vObj, okObj := v.(vtype.Time)
	if !okStr && !okObj {
		return NewResult().WithError(NewValidationError(t.message))
	}

	if okObj {
		vStr = vObj.String()
	}

	vt, err := time.Parse(t.format, vStr)
	if err != nil {
		return NewResult().WithError(
			NewValidationError(t.formatMessage).
				WithParams(
					map[string]any{
						"format": t.format,
					},
				),
		)
	}

	result := NewResult()

	if t.min != nil {
		minTime := t.min()
		if vt.Before(minTime) {
			result = result.WithError(
				NewValidationError(t.tooSmallMessage).
					WithParams(
						map[string]any{
							"min": minTime,
						},
					),
			)
		}
	}

	if t.max != nil {
		maxTime := t.max()
		if vt.After(maxTime) {
			result = result.WithError(
				NewValidationError(t.tooBigMessage).
					WithParams(
						map[string]any{
							"max": maxTime,
						},
					),
			)
		}
	}

	if result.IsValid() {
		if okObj {
			*vObj.Time() = vt
		}

		return nil
	}

	return result
}
