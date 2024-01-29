package validator

import "context"

const (
	msisdnRegexp = `^\d+$`
)

type MSISDN struct {
	basicRule MatchRegularExpression
}

func NewMSISDN() MSISDN {
	return MSISDN{
		basicRule: NewMatchRegularExpression(msisdnRegexp).
			WithMessage("MSISDN format is invalid."),
	}
}

func (s MSISDN) WithMessage(message string) MSISDN {
	s.basicRule = s.basicRule.WithMessage(message)
	return s
}

func (s MSISDN) ValidateValue(ctx context.Context, value any) error {
	return s.basicRule.ValidateValue(ctx, value)
}
