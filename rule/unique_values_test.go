package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueValues_ValidateValue_HasExpectedValue_NoError(t *testing.T) {
	err := NewUniqueValues().ValidateValue([]string{"one", "two"})
	assert.NoError(t, err)

	err = NewUniqueValues().ValidateValue([]int{1, 2})
	assert.NoError(t, err)
}

func TestUniqueValues_ValidateValue_HasExpectedValue_ReturnsError(t *testing.T) {
	err := NewUniqueValues().ValidateValue([]string{"two", "two"})
	assert.Error(t, err)

	err = NewUniqueValues().ValidateValue([]int{1, 1})
	assert.Error(t, err)
}
