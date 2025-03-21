package validator

import (
	"context"

	"github.com/xwb1989/sqlparser"
)

type SQL struct {
	message       string
	whenFunc      WhenFunc
	asWhereClause bool
	skipEmpty     bool
	skipError     bool
}

func NewSQL() *SQL {
	return &SQL{
		message: "Value is invalid sql.",
	}
}

func (r *SQL) AsWhereClause() *SQL {
	rc := *r
	rc.asWhereClause = true

	return &rc
}

func (r *SQL) WithMessage(message string) *SQL {
	rc := *r
	rc.message = message

	return &rc
}

func (r *SQL) When(v WhenFunc) *SQL {
	rc := *r
	rc.whenFunc = v

	return &rc
}

func (r *SQL) when() WhenFunc {
	return r.whenFunc
}

func (r *SQL) setWhen(v WhenFunc) {
	r.whenFunc = v
}

func (r *SQL) SkipOnEmpty() *SQL {
	rc := *r
	rc.skipEmpty = true

	return &rc
}

func (r *SQL) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *SQL) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *SQL) SkipOnError() *SQL {
	rs := *r
	rs.skipError = true

	return &rs
}

func (r *SQL) shouldSkipOnError() bool {
	return r.skipError
}

func (r *SQL) setSkipOnError(v bool) {
	r.skipError = v
}

func (r *SQL) ValidateValue(_ context.Context, value any) error {
	v, ok := toString(value)
	if !ok {
		return NewResult().WithError(NewValidationError(r.message))
	}

	if r.asWhereClause {
		v = "select * from x where " + v
	}

	if _, err := sqlparser.Parse(v); err != nil {
		return NewResult().WithError(NewValidationError(r.message))
	}

	return nil
}
