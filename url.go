package validator

import (
	"context"
	"strings"

	"golang.org/x/net/idna"

	"github.com/raoptimus/validator.go/v2/regexpc"
)

var regexpDomain, _ = regexpc.Compile(`://([^/]+)`)

const AllowAnyURLSchema = "*"
const defaultURLRegexpPattern = `^{schemes}:\/\/(([a-zA-Z0-9][a-zA-Z0-9_-]*)(\.[a-zA-Z0-9][a-zA-Z0-9_-]*)+)(?::\d{1,5})?([?\/#].*$|$)`

type URL struct {
	pattern      string
	validSchemes []string
	enableIDN    bool
	message      string
	whenFunc     WhenFunc
	skipEmpty    bool
	skipError    bool
}

func NewURL() *URL {
	return &URL{
		pattern:      defaultURLRegexpPattern,
		validSchemes: []string{"http", "https"},
		enableIDN:    false,
		message:      "This value is not a valid URL.",
	}
}

func (r *URL) WithPattern(pattern string) *URL {
	rc := *r
	rc.pattern = pattern

	return &rc
}

func (r *URL) WithValidScheme(scheme ...string) *URL {
	rc := *r
	for i, sh := range scheme {
		if sh == AllowAnyURLSchema {
			scheme[i] = ".*?"
		}
	}
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

func (r *URL) SkipOnEmpty() *URL {
	rc := *r
	rc.skipEmpty = true

	return &rc
}

func (r *URL) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *URL) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *URL) SkipOnError() *URL {
	rs := *r
	rs.skipError = true

	return &rs
}

func (r *URL) shouldSkipOnError() bool {
	return r.skipError
}
func (r *URL) setSkipOnError(v bool) {
	r.skipError = v
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

	pattern := r.getPattern()
	rgxp, err := regexpc.Compile(pattern)
	if err != nil {
		return NewResult().WithError(NewValidationError(r.message))
	}

	if !rgxp.MatchString(v) {
		return NewResult().WithError(NewValidationError(r.message))
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

func (r *URL) getPattern() string {
	if !strings.Contains(r.pattern, "{schemes}") {
		return r.pattern
	}

	return strings.ReplaceAll(
		r.pattern,
		"{schemes}",
		"((?i)"+strings.Join(r.validSchemes, "|")+")",
	)
}
