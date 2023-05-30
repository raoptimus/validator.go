package rule

const (
	emailRegexp = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)

type Email struct {
	basicRule MatchRegularExpression
}

func NewEmail() Email {
	return Email{
		basicRule: NewMatchRegularExpression(emailRegexp).
			WithMessage("Invalid Email format."),
	}
}

func (s Email) WithMessage(message string) Email {
	s.basicRule = s.basicRule.WithMessage(message)
	return s
}

func (s Email) ValidateValue(value any) error {
	return s.basicRule.ValidateValue(value)
}
