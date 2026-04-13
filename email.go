/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
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
