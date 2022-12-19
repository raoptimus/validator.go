package rule

import (
	"regexp"
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

	r, err := regexp.Compile(s.pattern)
	if err != nil {
		return err
	}

	if !r.MatchString(v) {
		return NewResult().WithError(formatMessage(s.message))
	}
	return nil
}
