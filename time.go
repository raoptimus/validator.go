package validator

import (
	"context"
	"time"

	"github.com/raoptimus/validator.go/v2/vtype"
)

type TimeFunc func(ctx context.Context) (time.Time, error)

type Time struct {
	message         string
	formatMessage   string
	tooBigMessage   string
	tooSmallMessage string
	format          string
	min             TimeFunc
	max             TimeFunc
	whenFunc        WhenFunc
	skipEmpty       bool
	skipError       bool
}

func NewTime() *Time {
	return &Time{
		message:         "Value is invalid",
		formatMessage:   "Format of the time value must be equal {format}",
		tooBigMessage:   "Time must be no greater than {max}.",
		tooSmallMessage: "Time must be no less than {min}.",
		format:          time.RFC3339,
		min:             nil,
		max:             nil,
	}
}

func (r *Time) WithMessage(message string) *Time {
	rc := *r
	rc.message = message

	return &rc
}

func (r *Time) WithFormatMessage(message string) *Time {
	rc := *r
	rc.formatMessage = message

	return &rc
}

func (r *Time) WithTooSmallMessage(message string) *Time {
	rc := *r
	rc.tooSmallMessage = message

	return &rc
}

func (r *Time) WithTooBigMessage(message string) *Time {
	rc := *r
	rc.tooBigMessage = message

	return &rc
}

func (r *Time) WithFormat(format string) *Time {
	rc := *r
	rc.format = format

	return &rc
}

func (r *Time) WithMin(min TimeFunc) *Time {
	rc := *r
	rc.min = min

	return &rc
}

func (r *Time) WithMax(max TimeFunc) *Time {
	rc := *r
	rc.max = max

	return &rc
}

func (r *Time) When(v WhenFunc) *Time {
	rc := *r
	rc.whenFunc = v

	return &rc
}

func (r *Time) when() WhenFunc {
	return r.whenFunc
}

func (r *Time) setWhen(v WhenFunc) {
	r.whenFunc = v
}

func (r *Time) SkipOnEmpty() *Time {
	rc := *r
	rc.skipEmpty = true

	return &rc
}

func (r *Time) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *Time) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *Time) SkipOnError() *Time {
	rs := *r
	rs.skipError = true

	return &rs
}

func (r *Time) shouldSkipOnError() bool {
	return r.skipError
}
func (r *Time) setSkipOnError(v bool) {
	r.skipError = v
}

func (r *Time) ValidateValue(ctx context.Context, value any) error {
	v, valid := indirectValue(value)
	if !valid {
		return NewResult().WithError(NewValidationError(r.message))
	}

	vStr, okStr := toString(value)
	vObj, okObj := v.(vtype.Time)
	if !okStr && !okObj {
		return NewResult().WithError(NewValidationError(r.message))
	}

	if okObj {
		vStr = vObj.String()
	}

	vt, err := time.Parse(r.format, vStr)
	if err != nil {
		return NewResult().WithError(
			NewValidationError(r.formatMessage).
				WithParams(
					map[string]any{
						"format": r.format,
					},
				),
		)
	}

	result := NewResult()

	if r.min != nil {
		minTime, err := r.min(ctx)
		if err != nil {
			return err
		}
		if vt.Before(minTime) {
			result = result.WithError(
				NewValidationError(r.tooSmallMessage).
					WithParams(
						map[string]any{
							"min": minTime,
						},
					),
			)
		}
	}

	if r.max != nil {
		maxTime, err := r.max(ctx)
		if err != nil {
			return err
		}
		if vt.After(maxTime) {
			result = result.WithError(
				NewValidationError(r.tooBigMessage).
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
