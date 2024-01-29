package validator

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCallback struct {
	A int
	B int
}

func TestCallback_ValidateValue_NoError(t *testing.T) {
	rules := RuleSet{
		"A": {
			NewCallback(func(_ context.Context, value int) error {
				return nil
			}),
		},
	}

	ctx := context.Background()
	err := Validate(ctx, &TestCallback{A: 1, B: 2}, rules)
	assert.NoError(t, err)
}

func TestCallback_ValidateValue_Error(t *testing.T) {
	errAMustGreatB := errors.New("A must be great than B")

	rules := RuleSet{
		"A": {
			NewCallback(func(ctx context.Context, value int) error {
				if ds, ok := extractDataSet(ctx); ok {
					if obj, ok := ds.Data().(*TestCallback); ok {
						if obj.B > value {
							return errAMustGreatB
						}
					}
				}
				return nil
			}),
		},
	}

	ctx := context.Background()
	err := Validate(ctx, &TestCallback{A: 1, B: 2}, rules)
	assert.ErrorIs(t, err, errAMustGreatB)
}
