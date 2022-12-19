package rule

import (
	"strings"

	"github.com/raoptimus/validator.go/regexpc"
	"golang.org/x/net/idna"
)

var regxpDomain, _ = regexpc.Compile(`://([^/]+)`)

type Url struct {
	pattern      string
	validSchemes []string
	enableIDN    bool
	message      string
}

func NewUrl() Url {
	return Url{
		pattern:      `^{schemes}:\/\/(([a-zA-Z0-9][a-zA-Z0-9_-]*)(\.[a-zA-Z0-9][a-zA-Z0-9_-]*)+)(?::\d{1,5})?([?\/#].*$|$)`,
		validSchemes: []string{"http", "https"},
		enableIDN:    false,
		message:      "This value is not a valid URL.",
	}
}

func (u Url) WithPattern(pattern string) Url {
	u.pattern = pattern
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

	r, err := regexpc.Compile(u.getPattern())
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

func (u Url) getPattern() string {
	if !strings.Contains(u.pattern, "{schemes}") {
		return u.pattern
	}

	return strings.ReplaceAll(
		u.pattern,
		"{schemes}",
		"((?i)"+strings.Join(u.validSchemes, "|")+")",
	)
}
