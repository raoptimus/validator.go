package validator

import (
	"context"
	"encoding/json"
	"fmt"
)

type JSON struct {
	message string

	whenFunc  WhenFunc
	skipEmpty bool
	skipError bool
}

func NewJSON() *JSON {
	return &JSON{
		message: "Must be a valid JSON",
	}
}

func (j *JSON) ValidateValue(_ context.Context, value any) error {
	var bytes []byte

	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case *string:
		bytes = []byte(*v)
	case []byte:
		bytes = v
	case *[]byte:
		bytes = *v
	case json.RawMessage:
		bytes = v
	case *json.RawMessage:
		bytes = *v
	case fmt.Stringer:
		if i, ok := v.(fmt.Stringer); ok {
			bytes = []byte(i.String())
		}
	default:
		return NewResult().WithError(NewValidationError(j.message))
	}

	if isValid := json.Valid(bytes); !isValid {
		return NewResult().WithError(NewValidationError(j.message))
	}

	return nil
}

func (j *JSON) When(v WhenFunc) *JSON {
	rc := *j
	rc.whenFunc = v

	return &rc
}

func (j *JSON) when() WhenFunc {
	return j.whenFunc
}

func (j *JSON) setWhen(v WhenFunc) {
	j.whenFunc = v
}

func (j *JSON) SkipOnEmpty() *JSON {
	rc := *j
	rc.skipEmpty = true

	return &rc
}

func (j *JSON) skipOnEmpty() bool {
	return j.skipEmpty
}

func (j *JSON) setSkipOnEmpty(v bool) {
	j.skipEmpty = v
}

func (j *JSON) SkipOnError() *JSON {
	rs := *j
	rs.skipError = true

	return &rs
}

func (j *JSON) shouldSkipOnError() bool {
	return j.skipError
}

func (j *JSON) setSkipOnError(v bool) {
	j.skipError = v
}
