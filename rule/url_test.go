package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrl_ValidateValue_Successfully(t *testing.T) {
	err := NewUrl().ValidateValue("https://example.com")
	assert.NoError(t, err)
}

func TestUrl_ValidateValue_IDN_Successfully(t *testing.T) {
	err := NewUrl().WithEnableIDN().ValidateValue("https://президент.рф")
	assert.NoError(t, err)
}

func TestUrlValidateValue_EmptyValue_ReturnsError(t *testing.T) {
	err := NewUrl().ValidateValue("")
	assert.Error(t, err)
}

func TestUrlValidateValue_InvalidValue_ReturnsError(t *testing.T) {
	err := NewUrl().ValidateValue("http://")
	assert.Error(t, err)
}

func TestUrlValidateValue_InvalidValue_ReturnsExpectedErrorMessage(t *testing.T) {
	err := NewUrl().WithMessage("test error").ValidateValue("http://")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "test error")
}
