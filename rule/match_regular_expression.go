package rule

import (
	"github.com/raoptimus/validator.go/regexpc"
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

func (s MatchRegularExpression) ValidateValue(value any) error {
	v, ok := value.(string)
	if !ok {
		return NewResult().WithError(formatMessage(s.message))
	}

	r, err := regexpc.Compile(s.pattern)
	if err != nil {
		return err
	}

	if !r.MatchString(v) {
		return NewResult().WithError(formatMessage(s.message))
	}
	return nil
}
