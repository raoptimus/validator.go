package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringLength_ValidateValue(t *testing.T) {
	r := NewStringLength(2, 4)
	err := r.ValidateValue(context.Background(), "-")
	assert.Error(t, err)
}
