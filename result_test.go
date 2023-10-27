package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult_Errors_Successfully(t *testing.T) {
	res := NewResult().WithError(NewValidationError("test err"))
	errs := res.Errors()
	assert.Equal(t, []*ValidationError{{Message: "test err"}}, errs)

	errs = append(errs, &ValidationError{Message: "invisible error"})
	assert.Equal(t, []*ValidationError{{Message: "test err"}}, res.Errors())

	res = res.WithError(NewValidationError("test2 err"))
	assert.Equal(t, []*ValidationError{{Message: "test err"}, {Message: "test2 err"}}, res.Errors())
}

func TestResult_IsValid_True(t *testing.T) {
	res := NewResult()
	assert.True(t, res.IsValid())
}

func TestResult_IsValid_False(t *testing.T) {
	res := NewResult().WithError(NewValidationError("test err"))
	assert.False(t, res.IsValid())
}
