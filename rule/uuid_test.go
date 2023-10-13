package rule

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUUID_ValidateValue_Successfully(t *testing.T) {
	err := NewUUID().ValidateValue("00000000-0000-0000-0000-000000000001")
	assert.NoError(t, err)
}

func TestUUID_ValidateValue_InvalidValue_Failed(t *testing.T) {
	err := NewUUID().ValidateValue("12323435-343")
	assert.Error(t, err)
}

func TestUUID_ValidateValue_ZeroUUID_Failed(t *testing.T) {
	err := NewUUID().ValidateValue("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
}

func TestUUID_ValidateValue_Version(t *testing.T) {
	v4 := uuid.Must(uuid.NewV4())
	v7 := uuid.Must(uuid.NewV7())

	v := NewUUID().
		WithVersion(UUIDVersionV7)

	err := v.ValidateValue(v4.String())
	assert.Error(t, err)

	err = v.ValidateValue(v7.String())
	assert.NoError(t, err)
}
