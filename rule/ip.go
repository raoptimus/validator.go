package rule

import (
	"net"
)

type IP struct {
	message               string
	ipv4NotAllowedMessage string
	ipv6NotAllowedMessage string
	allowV4               bool
	allowV6               bool
}

func NewIP() IP {
	return IP{
		message:               "Must be a valid IP address.",
		ipv4NotAllowedMessage: "Must not be an IPv4 address.",
		ipv6NotAllowedMessage: "Must not be an IPv6 address.",
		allowV4:               true,
		allowV6:               true,
	}
}

func (s IP) ValidateValue(value any) error {
	v, ok := value.(string)
	if !ok {
		return NewResult().WithError(formatMessage(s.message))
	}

	ip := net.ParseIP(v)
	if ip == nil {
		return NewResult().WithError(formatMessage(s.message))
	}

	// TODO: implement ipv4 and ipv4 validations

	return nil
}
