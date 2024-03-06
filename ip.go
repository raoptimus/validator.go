package validator

import (
	"context"
	"net"
)

type IP struct {
	message               string
	ipv4NotAllowedMessage string
	ipv6NotAllowedMessage string
	allowV4               bool
	allowV6               bool
	whenFunc              WhenFunc
	skipEmpty             bool
}

func NewIP() *IP {
	return &IP{
		message:               "Must be a valid IP address.",
		ipv4NotAllowedMessage: "Must not be an IPv4 address.",
		ipv6NotAllowedMessage: "Must not be an IPv6 address.",
		allowV4:               true,
		allowV6:               true,
	}
}

func (r *IP) WithMessage(v string) *IP {
	rc := *r
	rc.message = v

	return &rc
}

func (r *IP) When(v WhenFunc) *IP {
	rc := *r
	rc.whenFunc = v

	return &rc
}

func (r *IP) when() WhenFunc {
	return r.whenFunc
}

func (r *IP) setWhen(v WhenFunc) {
	r.whenFunc = v
}

func (r *IP) SkipOnEmpty(v bool) *IP {
	rc := *r
	rc.skipEmpty = v

	return &rc
}

func (r *IP) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *IP) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *IP) ValidateValue(_ context.Context, value any) error {
	v, ok := toString(value)
	if !ok {
		return NewResult().WithError(NewValidationError(r.message))
	}

	ip := net.ParseIP(v)
	if ip == nil {
		return NewResult().WithError(NewValidationError(r.message))
	}

	// TODO: implement ipv4 and ipv4 validations

	return nil
}
