package validator

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUUID_ValidateValue_Successfully(t *testing.T) {
	ctx := context.Background()
	err := NewUUID().ValidateValue(ctx, "00000000-0000-0000-0000-000000000001")
	assert.NoError(t, err)
}

func TestUUID_ValidateValue_InvalidValue_Failed(t *testing.T) {
	ctx := context.Background()
	err := NewUUID().ValidateValue(ctx, "12323435-343")
	assert.Error(t, err)
}

func TestUUID_ValidateValue_ZeroUUID_Failed(t *testing.T) {
	ctx := context.Background()
	err := NewUUID().ValidateValue(ctx, "00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
}

func TestUUID_ValidateValue_Version(t *testing.T) {
	ctx := context.Background()
	v4 := uuid.Must(uuid.NewV4())
	v7 := uuid.Must(uuid.NewV7())

	v := NewUUID().
		WithVersion(UUIDVersionV7)

	err := v.ValidateValue(ctx, v4.String())
	assert.Error(t, err)

	err = v.ValidateValue(ctx, v7.String())
	assert.NoError(t, err)
}
