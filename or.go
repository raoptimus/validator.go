package validator

import (
	"context"
)

type OR struct {
	rules     []Rule
	message   string
	whenFunc  WhenFunc
	skipEmpty bool
	skipError bool
}

func NewOR(message string, rules ...Rule) *OR {
	return &OR{
		message: message,
		rules:   rules,
	}
}

func (o *OR) WithMessage(message string) *OR {
	rc := *o
	rc.message = message

	return &rc
}

func (o *OR) When(v WhenFunc) *OR {
	rc := *o
	rc.whenFunc = v

	return &rc
}

func (o *OR) when() WhenFunc {
	return o.whenFunc
}

func (o *OR) setWhen(v WhenFunc) {
	o.whenFunc = v
}

func (o *OR) SkipOnEmpty() *OR {
	rc := *o
	rc.skipEmpty = true

	return &rc
}

func (o *OR) skipOnEmpty() bool {
	return o.skipEmpty
}

func (o *OR) setSkipOnEmpty(v bool) {
	o.skipEmpty = v
}

func (o *OR) SkipOnError() *OR {
	rs := *o
	rs.skipError = true

	return &rs
}

func (o *OR) shouldSkipOnError() bool {
	return o.skipError
}
func (o *OR) setSkipOnError(v bool) {
	o.skipError = v
}

func (o *OR) ValidateValue(ctx context.Context, value any) error {
	for _, r := range o.rules {
		if err := r.ValidateValue(ctx, value); err == nil {
			return nil
		}
	}

	return NewResult().WithError(NewValidationError(o.message))
}
