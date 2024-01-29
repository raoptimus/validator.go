package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchRegularExpression_ValidateValue_ValueNotString(t *testing.T) {
	ctx := context.Background()
	err := NewMatchRegularExpression(``).ValidateValue(ctx, 1)
	assert.Error(t, err)
	assert.Equal(t, "Value is invalid.", err.Error())
}
