package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEach_ValidateValue_FirstValueIs1_NoError(t *testing.T) {
	err := NewEach(NewNumber(1, 2)).ValidateValue([]int{1})
	assert.NoError(t, err)
}

func TestEach_ValidateValue_FirstValueIs0_Error(t *testing.T) {
	err := NewEach(NewNumber(1, 2)).ValidateValue([]int{0})
	assert.Error(t, err)
	assert.Equal(t, "Value must be no less than 1.", err.Error())
}
