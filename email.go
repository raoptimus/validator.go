package validator

const emailRegexp = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{1,61}$`

type Email struct {
	*MatchRegularExpression
}

func NewEmail() *Email {
	return &Email{
		MatchRegularExpression: NewMatchRegularExpression(emailRegexp).
			WithMessage("Email is not a valid email."),
	}
}
