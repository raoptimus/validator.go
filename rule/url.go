package rule

import (
	"net/url"
	"strings"

	"golang.org/x/net/idna"

	"github.com/raoptimus/validator.go/regexpc"
)

var regxpDomain, _ = regexpc.Compile(`://([^/]+)`)

const AllowAnyURLSchema = "*"

type Url struct {
	validSchemes []string
	enableIDN    bool
	message      string
}

func NewUrl() Url {
	return Url{
		validSchemes: []string{"http", "https"},
		enableIDN:    false,
		message:      "This value is not a valid URL.",
	}
}

func (u Url) WithValidScheme(scheme ...string) Url {
	u.validSchemes = scheme
	return u
}

func (u Url) WithMessage(message string) Url {
	u.message = message
	return u
}

func (u Url) WithEnableIDN() Url {
	u.enableIDN = true
	return u
}

func (u Url) ValidateValue(value any) error {
	v, ok := value.(string)
	// make sure the length is limited to avoid DOS attacks
	if !ok || len(v) >= 2000 {
		return NewResult().WithError(formatMessage(u.message))
	}

	if u.enableIDN {
		v = u.convertIDN(v)
	}

	uri, err := url.Parse(v)
	if err != nil {
		return NewResult().WithError(formatMessage(u.message))
	}

	if len(uri.Scheme) == 0 || (len(uri.Host) == 0 && len(uri.Opaque) == 0) {
		return NewResult().WithError(formatMessage(u.message))
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
			return NewResult().WithError(formatMessage(u.message))
		}
	}

	return nil
}

func (u Url) convertIDN(value string) string {
	if !strings.Contains(value, "://") {
		return u.idnToASCII(value)
	}

	return regxpDomain.ReplaceAllStringFunc(value, func(m string) string {
		p := regxpDomain.FindStringSubmatch(m)
		return "://" + u.idnToASCII(p[1])
	})
}

func (u Url) idnToASCII(idn string) string {
	if d, err := idna.ToASCII(idn); err == nil {
		return d
	} else {
		return idn
	}
}
