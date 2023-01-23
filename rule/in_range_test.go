package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEach_ValidateValue_HasExpectedValue_NoError(t *testing.T) {
	err := NewInRange([]any{"one", "two"}).ValidateValue("two")
	assert.NoError(t, err)
}

func TestEach_ValidateValue_NotHasValue_NoError(t *testing.T) {
	err := NewInRange([]any{"two"}).ValidateValue("one")
	assert.Error(t, err)
}
