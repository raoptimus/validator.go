/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator

import (
	"context"
	"net"
)

type MAC struct {
	message   string
	whenFunc  WhenFunc
	skipEmpty bool
	skipError bool
}

func NewMAC() *MAC {
	return &MAC{
		message: MessageInvalidMAC,
	}
}

func (m *MAC) WithMessage(v string) *MAC {
	rc := *m
	rc.message = v

	return &rc
}

func (m *MAC) When(v WhenFunc) *MAC {
	rc := *m
	rc.whenFunc = v

	return &rc
}

func (m *MAC) when() WhenFunc {
	return m.whenFunc
}

func (m *MAC) setWhen(v WhenFunc) {
	m.whenFunc = v
}

func (m *MAC) SkipOnEmpty() *MAC {
	rc := *m
	rc.skipEmpty = true

	return &rc
}

func (m *MAC) skipOnEmpty() bool {
	return m.skipEmpty
}

func (m *MAC) setSkipOnEmpty(v bool) {
	m.skipEmpty = v
}

func (m *MAC) SkipOnError() *MAC {
	rs := *m
	rs.skipError = true

	return &rs
}

func (m *MAC) shouldSkipOnError() bool {
	return m.skipError
}
func (m *MAC) setSkipOnError(v bool) {
	m.skipError = v
}

func (m *MAC) ValidateValue(_ context.Context, value any) error {
	v, ok := toString(value)
	if !ok {
		return NewResult().WithError(NewValidationError(m.message))
	}

	if _, err := net.ParseMAC(v); err != nil {
		return NewResult().WithError(NewValidationError(m.message))
	}

	return nil
}
