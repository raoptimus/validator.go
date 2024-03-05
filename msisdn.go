package validator

const (
	msisdnRegexp = `^\d+$`
)

type MSISDN struct {
	MatchRegularExpression
}

func NewMSISDN() MSISDN {
	return MSISDN{
		MatchRegularExpression: NewMatchRegularExpression(msisdnRegexp).
			WithMessage("MSISDN format is invalid."),
	}
}
