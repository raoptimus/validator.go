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

func NewMatchRegularExpression(pattern string) MatchRegularExpression {
	return MatchRegularExpression{
		message: "Value is invalid.",
		pattern: pattern,
	}
}

func (s MatchRegularExpression) WithMessage(message string) MatchRegularExpression {
	s.message = message

	return s
}

func (s MatchRegularExpression) When(v WhenFunc) MatchRegularExpression {
	s.whenFunc = v

	return s
}

func (s MatchRegularExpression) when() WhenFunc {
	return s.whenFunc
}

func (s MatchRegularExpression) SkipOnEmpty(v bool) MatchRegularExpression {
	s.skipEmpty = v

	return s
}

func (s MatchRegularExpression) skipOnEmpty() bool {
	return s.skipEmpty
}

func (s MatchRegularExpression) ValidateValue(_ context.Context, value any) error {
	v, ok := toString(value)
	if !ok {
		return NewResult().WithError(NewValidationError(s.message))
	}

	r, err := regexpc.Compile(s.pattern)
	if err != nil {
		return err
	}

	if !r.MatchString(v) {
		return NewResult().WithError(NewValidationError(s.message))
	}

	return nil
}
