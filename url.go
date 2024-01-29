package validator

import (
	"context"
	"net/url"
	"strings"

	"golang.org/x/net/idna"

	"github.com/raoptimus/validator.go/regexpc"
)

var regxpDomain, _ = regexpc.Compile(`://([^/]+)`)

const AllowAnyURLSchema = "*"

type URL struct {
	validSchemes []string
	enableIDN    bool
	message      string
}

func NewURL() URL {
	return URL{
		validSchemes: []string{"http", "https"},
		enableIDN:    false,
		message:      "This value is not a valid URL.",
	}
}

func (u URL) WithValidScheme(scheme ...string) URL {
	u.validSchemes = scheme
	return u
}

func (u URL) WithMessage(message string) URL {
	u.message = message
	return u
}

func (u URL) WithEnableIDN() URL {
	u.enableIDN = true
	return u
}

func (u URL) ValidateValue(_ context.Context, value any) error {
	v, ok := value.(string)
	// make sure the length is limited to avoid DOS attacks
	if !ok || len(v) >= 2000 {
		return NewResult().WithError(NewValidationError(u.message))
	}

	if u.enableIDN {
		v = u.convertIDN(v)
	}

	uri, err := url.Parse(v)
	if err != nil {
		return NewResult().WithError(NewValidationError(u.message))
	}

	if len(uri.Scheme) == 0 || (len(uri.Host) == 0 && len(uri.Opaque) == 0) {
		return NewResult().WithError(NewValidationError(u.message))
	}

	if len(u.validSchemes) > 0 && u.validSchemes[0] != AllowAnyURLSchema {
		isValidScheme := false
		for _, s := range u.validSchemes {
			if s == uri.Scheme {
				isValidScheme = true
				break
			}
		}

		if !isValidScheme {
			return NewResult().WithError(NewValidationError(u.message))
		}
	}

	return nil
}

func (u URL) convertIDN(value string) string {
	if !strings.Contains(value, "://") {
		return u.idnToASCII(value)
	}

	return regxpDomain.ReplaceAllStringFunc(value, func(m string) string {
		p := regxpDomain.FindStringSubmatch(m)
		return "://" + u.idnToASCII(p[1])
	})
}

func (u URL) idnToASCII(idn string) string {
	if d, err := idna.ToASCII(idn); err == nil {
		return d
	} else {
		return idn
	}
}
