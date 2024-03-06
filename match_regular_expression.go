package validator

import (
	"context"

	"github.com/raoptimus/validator.go/v2/regexpc"
)

type MatchRegularExpression struct {
	message   string
	pattern   string
	whenFunc  WhenFunc
	skipEmpty bool
}

func NewMatchRegularExpression(pattern string) *MatchRegularExpression {
	return &MatchRegularExpression{
		message: "Value is invalid.",
		pattern: pattern,
	}
}

func (r *MatchRegularExpression) WithMessage(message string) *MatchRegularExpression {
	rc := *r
	rc.message = message

	return &rc
}

func (r *MatchRegularExpression) When(v WhenFunc) *MatchRegularExpression {
	rc := *r
	rc.whenFunc = v

	return &rc
}

func (r *MatchRegularExpression) when() WhenFunc {
	return r.whenFunc
}

func (r *MatchRegularExpression) setWhen(v WhenFunc) {
	r.whenFunc = v
}

func (r *MatchRegularExpression) SkipOnEmpty(v bool) *MatchRegularExpression {
	rc := *r
	rc.skipEmpty = v

	return &rc
}

func (r *MatchRegularExpression) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *MatchRegularExpression) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *MatchRegularExpression) ValidateValue(_ context.Context, value any) error {
	v, ok := toString(value)
	if !ok {
		return NewResult().WithError(NewValidationError(r.message))
	}

	rg, err := regexpc.Compile(r.pattern)
	if err != nil {
		return err
	}

	if !rg.MatchString(v) {
		return NewResult().WithError(NewValidationError(r.message))
	}

	return nil
}
