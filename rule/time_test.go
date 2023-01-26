package rule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/raoptimus/validator.go/ctype"
)

func TestTime_ValidateValue_ValidValueString_NoError(t *testing.T) {
	// utc
	err := NewTime().ValidateValue("2006-01-02T15:04:05Z")
	assert.NoError(t, err)

	// another timezone
	err = NewTime().ValidateValue("2006-01-02T15:04:05+07:00")
	assert.NoError(t, err)
}

func TestTime_ValidateValue_InvalidValueString_Error(t *testing.T) {
	err := NewTime().ValidateValue("2006-01-02 15:04:05")
	assert.Error(t, err)
}

func TestTime_ValidateValue_ValueIsNil_Error(t *testing.T) {
	err := NewTime().ValidateValue(nil)
	assert.Error(t, err)
}

func TestTime_ValidateValue_ValueIsObject_NoError(t *testing.T) {
	var tm ctype.Time
	err := tm.UnmarshalJSON([]byte("2006-01-02T15:04:05Z"))
	assert.NoError(t, err)

	err = NewTime().ValidateValue(tm)
	assert.NoError(t, err)

	parsed := tm.Time()

	assert.NotNil(t, parsed)
	assert.Equal(t, "2006-01-02T15:04:05Z", tm.String())
	assert.Equal(t, "2006-01-02T15:04:05Z", parsed.Format(time.RFC3339))
}
