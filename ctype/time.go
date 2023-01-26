package ctype

import (
	"strings"
	"time"
)

func NewTime(unvalidatedTime string) Time {
	return Time{
		validatedTime:   &time.Time{},
		unvalidatedTime: unvalidatedTime,
	}
}

type Time struct {
	validatedTime   *time.Time
	unvalidatedTime string
}

func (t *Time) Time() (*time.Time, error) {
	if t.validatedTime == nil {
		return nil, ErrNotFilledTime
	}

	return t.validatedTime, nil
}

func (t *Time) String() string {
	return t.unvalidatedTime
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Time) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	t.unvalidatedTime = strings.Trim(string(data), "\"")
	t.validatedTime = &time.Time{}
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *Time) UnmarshalText(data []byte) error {
	t.unvalidatedTime = string(data)
	t.validatedTime = &time.Time{}
	return nil
}
