package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchRegularExpression_ValidateValue_ValueNotString(t *testing.T) {
	err := NewMatchRegularExpression(``).ValidateValue(1)
	assert.Error(t, err)
	assert.Equal(t, "Value is invalid.", err.Error())
}
