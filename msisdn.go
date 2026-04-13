/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator

const msisdnRegexp = `^\d+$`

type MSISDN struct {
	*MatchRegularExpression
}

func NewMSISDN() MSISDN {
	return MSISDN{
		MatchRegularExpression: NewMatchRegularExpression(msisdnRegexp).
			WithMessage("MSISDN format is invalid."),
	}
}
