package validator

import (
	"context"
	"net/url"
	"strings"

	"golang.org/x/net/idna"

	"github.com/raoptimus/validator.go/v2/regexpc"
)

var regexpDomain, _ = regexpc.Compile(`://([^/]+)`)

const AllowAnyURLSchema = "*"

type URL struct {
	validSchemes []string
	enableIDN    bool
	message      string
	whenFunc     WhenFunc
	skipEmpty    bool
}

func NewURL() *URL {
	return &URL{
		validSchemes: []string{"http", "https"},
		enableIDN:    false,
		message:      "This value is not a valid URL.",
	}
}

func (r *URL) WithValidScheme(scheme ...string) *URL {
	rc := *r
	rc.validSchemes = scheme

	return &rc
}

func (r *URL) WithMessage(message string) *URL {
	rc := *r
	rc.message = message

	return &rc
}

func (r *URL) WithEnableIDN() *URL {
	rc := *r
	rc.enableIDN = true

	return &rc
}

func (r *URL) When(v WhenFunc) *URL {
	rc := *r
	rc.whenFunc = v

	return &rc
}

func (r *URL) when() WhenFunc {
	return r.whenFunc
}

func (r *URL) setWhen(v WhenFunc) {
	r.whenFunc = v
}

func (r *URL) SkipOnEmpty(v bool) *URL {
	rc := *r
	rc.skipEmpty = v

	return &rc
}

func (r *URL) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *URL) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *URL) ValidateValue(_ context.Context, value any) error {
	v, ok := toString(value)
	// make sure the length is limited to avoid DOS attacks
	if !ok || len(v) >= 2000 {
		return NewResult().WithError(NewValidationError(r.message))
	}

	if r.enableIDN {
		v = r.convertIDN(v)
	}

	uri, err := url.Parse(v)
	if err != nil {
		return NewResult().WithError(NewValidationError(r.message))
	}

	if len(uri.Scheme) == 0 || (len(uri.Host) == 0 && len(uri.Opaque) == 0) {
		return NewResult().WithError(NewValidationError(r.message))
	}

	if len(r.validSchemes) > 0 && r.validSchemes[0] != AllowAnyURLSchema {
		isValidScheme := false
		for _, s := range r.validSchemes {
			if s == uri.Scheme {
				isValidScheme = true
				break
			}
		}

		if !isValidScheme {
			return NewResult().WithError(NewValidationError(r.message))
		}
	}

	return nil
}

func (r *URL) convertIDN(value string) string {
	if !strings.Contains(value, "://") {
		return r.idnToASCII(value)
	}

	return regexpDomain.ReplaceAllStringFunc(value, func(m string) string {
		p := regexpDomain.FindStringSubmatch(m)
		return "://" + r.idnToASCII(p[1])
	})
}

func (r *URL) idnToASCII(idn string) string {
	if d, err := idna.ToASCII(idn); err == nil {
		return d
	} else {
		return idn
	}
}
