package validator

import (
	"context"
	"testing"
	"time"

	"github.com/raoptimus/validator.go/v2/vtype"
	"github.com/stretchr/testify/assert"
)

func TestTime_ValidateValue_ValidValueString_NoError(t *testing.T) {
	ctx := context.Background()

	// utc
	err := NewTime().ValidateValue(ctx, "2006-01-02T15:04:05Z")
	assert.NoError(t, err)

	// another timezone
	err = NewTime().ValidateValue(ctx, "2006-01-02T15:04:05+07:00")
	assert.NoError(t, err)
}

func TestTime_ValidateValue_InvalidValueString_Error(t *testing.T) {
	ctx := context.Background()
	err := NewTime().ValidateValue(ctx, "2006-01-02 15:04:05")
	assert.Error(t, err)
}

func TestTime_ValidateValue_ValueIsNil_Error(t *testing.T) {
	ctx := context.Background()
	err := NewTime().ValidateValue(ctx, nil)
	assert.Error(t, err)
}

func TestTime_ValidateValue_ValueIsObject_NoError(t *testing.T) {
	ctx := context.Background()
	var tm vtype.Time
	err := tm.UnmarshalJSON([]byte("2006-01-02T15:04:05Z"))
	assert.NoError(t, err)

	err = NewTime().ValidateValue(ctx, tm)
	assert.NoError(t, err)

	parsed := tm.Time()

	assert.NotNil(t, parsed)
	assert.Equal(t, "2006-01-02T15:04:05Z", tm.String())
	assert.Equal(t, "2006-01-02T15:04:05Z", parsed.Format(time.RFC3339))
}
