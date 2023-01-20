package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTime_ValidateValue_ValidValue_NoError(t *testing.T) {
	// utc
	err := NewTime().ValidateValue("2006-01-02T15:04:05Z")
	assert.NoError(t, err)

	// another timezone
	err = NewTime().ValidateValue("2006-01-02T15:04:05+07:00")
	assert.NoError(t, err)
}

func TestTime_ValidateValue_InvalidValue_Error(t *testing.T) {
	err := NewTime().ValidateValue("2006-01-02 15:04:05")
	assert.Error(t, err)
}
