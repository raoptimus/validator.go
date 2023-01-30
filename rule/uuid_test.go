package rule

import (
	"testing"

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
