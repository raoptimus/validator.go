package ctype

import (
	"time"
)

type Time struct {
	validatedTime   *time.Time
	unvalidatedTime string
}

func (t *Time) Time() *time.Time {
	return t.validatedTime
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
	t.unvalidatedTime = string(data)
	tm := time.UnixMilli(0)
	t.validatedTime = &tm
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *Time) UnmarshalText(data []byte) error {
	t.unvalidatedTime = string(data)
	return nil
}
