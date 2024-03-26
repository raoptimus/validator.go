package validator

import (
	"context"
	"reflect"
)

type Required struct {
	message   string
	whenFunc  WhenFunc
	skipError bool
}

func NewRequired() *Required {
	return &Required{
		message: "Value cannot be blank.",
	}
}

func (r *Required) When(v WhenFunc) *Required {
	rc := *r
	rc.whenFunc = v

	return &rc
}

func (r *Required) when() WhenFunc {
	return r.whenFunc
}

func (r *Required) setWhen(v WhenFunc) {
	r.whenFunc = v
}

func (r *Required) WithMessage(message string) *Required {
	rc := *r
	rc.message = message

	return &rc
}

// deprecated: should be removed
func (r *Required) WithAllowZeroValue() *Required {
	return r
}

func (r *Required) SkipOnError() *Required {
	rs := *r
	rs.skipError = true

	return &rs
}

func (r *Required) shouldSkipOnError() bool {
	return r.skipError
}
func (r *Required) setSkipOnError(v bool) {
	r.skipError = v
}

func (r *Required) ValidateValue(_ context.Context, value any) error {
	v := reflect.ValueOf(value)
	if valueIsEmpty(v) {
		return NewResult().WithError(NewValidationError(r.message))
	}

	return nil
}
