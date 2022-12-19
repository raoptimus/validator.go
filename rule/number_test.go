package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumber_ValidateValue_ValueNotNumber(t *testing.T) {
	err := NewNumber(0, 0).ValidateValue("abc")
	assert.Error(t, err)
	assert.Equal(t, "Value must be a number.", err.Error())
}

func TestNumber_ValidateValue_ValueNotNumberCustomError(t *testing.T) {
	err := NewNumber(0, 0).
		WithNotNumberMessage("test").
		ValidateValue("abc")
	assert.Error(t, err)
	assert.Equal(t, "test", err.Error())
}

func TestNumber_ValidateValue_ValueLessThanMin(t *testing.T) {
	err := NewNumber(1, 10).ValidateValue(0)
	assert.Error(t, err)
	assert.Equal(t, "Value must be no less than 1.", err.Error())
}

func TestNumber_ValidateValue_ValueLessThanMinCustomError(t *testing.T) {
	err := NewNumber(1, 10).WithTooSmallMessage("test {min}").ValidateValue(0)
	assert.Error(t, err)
	assert.Equal(t, "test 1", err.Error())
}

func TestNumber_ValidateValue_ValueGreatThanMax(t *testing.T) {
	err := NewNumber(1, 10).ValidateValue(11)
	assert.Error(t, err)
	assert.Equal(t, "Value must be no greater than 10.", err.Error())
}

func TestNumber_ValidateValue_ValueGreatThanMaxCustomError(t *testing.T) {
	err := NewNumber(1, 10).
		WithTooBigMessage("test {max}").
		ValidateValue(11)
	assert.Error(t, err)
	assert.Equal(t, "test 10", err.Error())
}

func TestNumber_ValidateValue_ErrTypeOfResultSet(t *testing.T) {
	err := NewNumber(1, 10).ValidateValue(11)
	assert.Error(t, err)
	assert.IsType(t, Result{}, err)
}

func TestNumber_ValidateValue_Successfully(t *testing.T) {
	err := NewNumber(1, 10).ValidateValue(2)
	assert.NoError(t, err)
}

func TestNumber_ValidateValue_SuccessfullyZero(t *testing.T) {
	err := NewNumber(0, 0).ValidateValue(0)
	assert.NoError(t, err)
}
