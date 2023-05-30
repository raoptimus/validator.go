package rule

const (
	msisdnRegexp = `^\d+$`
)

type MSISDN struct {
	basicRule MatchRegularExpression
}

func NewMSISDN() MSISDN {
	return MSISDN{
		basicRule: NewMatchRegularExpression(msisdnRegexp).
			WithMessage("Invalid MSISDN format."),
	}
}

func (s MSISDN) WithMessage(message string) MSISDN {
	s.basicRule = s.basicRule.WithMessage(message)
	return s
}

func (s MSISDN) ValidateValue(value any) error {
	return s.basicRule.ValidateValue(value)
}
