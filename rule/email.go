package rule

import (
	"strings"
)

type Email struct {
	message string
}

func NewEmail() Email {
	return Email{
		message: "Email is not a valid email",
	}
}

func (e *Email) ValidateValue(value any) error {
	v, ok := value.(string)
	if !ok {
		return NewResult().WithError(formatMessage(e.message))
	}

	if !strings.Contains(v, "@") {
		return NewResult().WithError(formatMessage(e.message))
	}

	result := NewResult()

	if !result.IsValid() {
		return result
	}
	return nil
}
