package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateValue_Int_Successfully(t *testing.T) {
	ctx := context.Background()
	rules := []Rule{
		NewRequired(),
		NewNumber(1, 3),
	}

	err := ValidateValue(ctx, 1, rules...)
	assert.NoError(t, err)
}

func TestValidateValue_Int_Failure(t *testing.T) {
	ctx := context.Background()
	rules := []Rule{
		NewRequired(),
		NewNumber(1, 3),
	}

	err := ValidateValue(ctx, 0, rules...)
	assert.Error(t, err)
	assert.Equal(t, "Value cannot be blank. Value must be no less than 1.", err.Error())

	assert.ErrorAs(t, err, &Result{})
	expectedResult := NewResult().
		WithError(
			NewValidationError("Value cannot be blank."),
			NewValidationError("Value must be no less than 1.").
				WithParams(map[string]any{"min": int64(1), "max": int64(3)}),
		)
	assert.Equal(t, expectedResult, err)
}

func TestValidateValue_IntNilPtrValue_Successfully(t *testing.T) {
	ctx := context.Background()
	rules := []Rule{
		NewNumber(1, 3),
	}

	err := ValidateValue(ctx, nil, rules...)
	assert.NoError(t, err)
}

func TestValidateValue_IntPtrValue_Successfully(t *testing.T) {
	ctx := context.Background()
	rules := []Rule{
		NewNumber(1, 3),
	}

	v := 2

	err := ValidateValue(ctx, &v, rules...)
	assert.NoError(t, err)
}

func TestValidate_Map_Successfully(t *testing.T) {
	ctx := context.Background()
	rules := RuleSet{
		"count": {
			NewRequired(),
			NewNumber(1, 3),
		},
	}
	data := map[string]any{
		"count": 1,
	}

	err := Validate(ctx, data, rules)
	assert.NoError(t, err)
}

func TestValidate_Nil_Failure(t *testing.T) {
	ctx := context.Background()
	rules := RuleSet{
		"count": {
			NewRequired(),
			NewNumber(1, 3),
		},
	}

	err := Validate(ctx, nil, rules)
	assert.Error(t, err)
}
