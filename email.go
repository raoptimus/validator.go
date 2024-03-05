package validator

import "context"

const (
	emailRegexp = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)

type Email struct {
	basicRule MatchRegularExpression
	whenFunc  WhenFunc
	skipEmpty bool
}

func NewEmail() Email {
	return Email{
		basicRule: NewMatchRegularExpression(emailRegexp).
			WithMessage("Email is not a valid email."),
	}
}

func (e Email) When(v WhenFunc) Email {
	e.whenFunc = v

	return e
}

func (e Email) when() WhenFunc {
	return e.whenFunc
}

func (e Email) SkipOnEmpty(v bool) Email {
	e.skipEmpty = v

	return e
}

func (e Email) skipOnEmpty() bool {
	return e.skipEmpty
}

func (e Email) WithMessage(message string) Email {
	e.basicRule = e.basicRule.WithMessage(message)

	return e
}

func (e Email) ValidateValue(ctx context.Context, value any) error {
	return e.basicRule.ValidateValue(ctx, value)
}
