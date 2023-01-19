package rule

import (
	"strings"

	"golang.org/x/net/idna"

	"github.com/raoptimus/validator.go/regexpc"
)

var regxpDomain, _ = regexpc.Compile(`://([^/]+)`)

type Url struct {
	urlPattern   string
	validSchemes []string
	enableIDN    bool
	message      string
}

func NewUrl() Url {
	return Url{
		urlPattern:   `^{schemes}:\/\/(([a-zA-Z0-9][a-zA-Z0-9_-]*)(\.[a-zA-Z0-9][a-zA-Z0-9_-]*)+)(?::\d{1,5})?([?\/#].*$|$)`,
		validSchemes: []string{"http", "https"},
		enableIDN:    false,
		message:      "This value is not a valid URL.",
	}
}

func (u Url) WithPattern(pattern string) Url {
	u.urlPattern = pattern
	return u
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

	r, err := regexpc.Compile(u.pattern())
	if err != nil {
		return NewResult().WithError(err.Error())
	}

	if !r.MatchString(v) {
		return NewResult().WithError(formatMessage(u.message))
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

func (u Url) pattern() string {
	if !strings.Contains(u.urlPattern, "{schemes}") {
		return u.urlPattern
	}

	return strings.ReplaceAll(
		u.urlPattern,
		"{schemes}",
		"((?i)"+strings.Join(u.validSchemes, "|")+")",
	)
}
