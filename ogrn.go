package validator

import (
	"context"
	"strconv"

	"github.com/raoptimus/validator.go/v2/regexpc"
)

const (
	ogrnNumberRegExp = "^[0-9]+$"
	ogrnNumberLen    = 13
	ogrnipNumberLen  = 15
)

type OGRN struct {
	ogrnMessage              string
	ogrnipMessage            string
	invalidOGRNLengthMessage string
	whenFunc                 WhenFunc
	skipEmpty                bool
	skipError                bool
}

func NewOGRN() *OGRN {
	return &OGRN{
		ogrnMessage:              "This value is not a valid OGRN.",
		ogrnipMessage:            "This value is not a valid OGRNIP.",
		invalidOGRNLengthMessage: "This value should contain either 13 or 15 characters.",
	}
}

func (o *OGRN) WithOGRNMessage(v string) *OGRN {
	rc := *o
	rc.ogrnMessage = v

	return &rc
}

func (o *OGRN) WithOGRNIPMessage(v string) *OGRN {
	rc := *o
	rc.ogrnipMessage = v

	return &rc
}

func (o *OGRN) When(v WhenFunc) *OGRN {
	rc := *o
	rc.whenFunc = v

	return &rc
}

func (o *OGRN) when() WhenFunc {
	return o.whenFunc
}

func (o *OGRN) setWhen(v WhenFunc) {
	o.whenFunc = v
}

func (o *OGRN) SkipOnEmpty() *OGRN {
	rc := *o
	rc.skipEmpty = true

	return &rc
}

func (o *OGRN) skipOnEmpty() bool {
	return o.skipEmpty
}

func (o *OGRN) setSkipOnEmpty(v bool) {
	o.skipEmpty = v
}

func (o *OGRN) SkipOnError() *OGRN {
	rs := *o
	rs.skipError = true

	return &rs
}

func (o *OGRN) shouldSkipOnError() bool {
	return o.skipError
}
func (o *OGRN) setSkipOnError(v bool) {
	o.skipError = v
}

func (o *OGRN) ValidateValue(_ context.Context, value any) error {
	v, ok := toString(value)
	if !ok {
		return NewResult().WithError(NewValidationError(o.ogrnMessage))
	}

	rg, err := regexpc.Compile(ogrnNumberRegExp)
	if err != nil {
		return err
	}

	if !rg.MatchString(v) {
		return NewResult().WithError(NewValidationError(o.ogrnMessage))
	}

	switch len(v) {
	case ogrnNumberLen:
		return o.validateOGRN(v)
	case ogrnipNumberLen:
		return o.validateOGRNIP(v)
	default:
		return NewResult().WithError(NewValidationError(o.invalidOGRNLengthMessage))
	}
}

func (o *OGRN) validateOGRN(ogrn string) error {
	firstDigit := ogrn[0]
	if firstDigit != '1' && firstDigit != '5' {
		return NewResult().WithError(NewValidationError(o.ogrnMessage))
	}

	num, err := strconv.ParseInt(ogrn[:12], 10, 64)
	if err != nil {
		return NewResult().WithError(NewValidationError(o.ogrnMessage))
	}

	remainder := num % 11
	controlDigit := int(remainder % 10)

	// Если остаток равен 10, то контрольное число должно быть 0
	if remainder == 10 {
		controlDigit = 0
	}

	lastDigit, err := strconv.Atoi(ogrn[12:])
	if err != nil {
		return NewResult().WithError(NewValidationError(o.ogrnMessage))
	}

	if controlDigit != lastDigit {
		return NewResult().WithError(NewValidationError(o.ogrnMessage))
	}

	return nil
}

func (o *OGRN) validateOGRNIP(ogrnip string) error {
	if ogrnip[0] != '3' {
		return NewResult().WithError(NewValidationError(o.ogrnipMessage))
	}

	num, err := strconv.ParseInt(ogrnip[:14], 10, 64)
	if err != nil {
		return NewResult().WithError(NewValidationError(o.ogrnipMessage))
	}

	remainder := num % 13
	controlDigit := int(remainder % 10)

	lastDigit, err := strconv.Atoi(ogrnip[14:])
	if err != nil {
		return NewResult().WithError(NewValidationError(o.ogrnipMessage))
	}

	if controlDigit != lastDigit {
		return NewResult().WithError(NewValidationError(o.ogrnipMessage))
	}

	return nil
}
