package validator

import (
	"context"

	"github.com/raoptimus/validator.go/v2/regexpc"
)

type MatchRegularExpression struct {
	message string
	pattern string
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
