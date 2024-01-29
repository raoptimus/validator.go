package vtype

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

func (t *Time) Time() *time.Time {
	if t.validatedTime == nil {
		t.validatedTime = &time.Time{}
	}

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

	t.unvalidatedTime = strings.Trim(string(data), "\"")
	if t.validatedTime == nil {
		t.validatedTime = &time.Time{}
	}

	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *Time) UnmarshalText(data []byte) error {
	t.unvalidatedTime = string(data)
	if t.validatedTime == nil {
		t.validatedTime = &time.Time{}
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (t *Time) MarshalJSON() ([]byte, error) {
	data := []byte("\"" + t.unvalidatedTime + "\"")
	return data, nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (t *Time) MarshalText() ([]byte, error) {
	data := []byte("\"" + t.unvalidatedTime + "\"")
	return data, nil
}
