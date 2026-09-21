/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEach_ValidateValue_FirstValueIs1_NoError(t *testing.T) {
	ctx := context.Background()
	err := NewEach(NewNumber(1, 2)).ValidateValue(ctx, []int{1})
	assert.NoError(t, err)
}

func TestEach_ValidateValue_ConcurrentCalls_NoError(t *testing.T) {
	rule := NewEach(NewNumber(1, 2))

	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			assert.NoError(t, rule.ValidateValue(context.Background(), []int{1}))
		}()
	}
	wg.Wait()
}

func TestEach_ValidateValue_SharedNestedConcurrentCalls_NoError(t *testing.T) {
	t.Parallel()

	type item struct {
		Value string
	}

	sharedRule := NewNested(RuleSet{
		"Value": {NewRequired()},
	})

	const workers = 100
	ctx := t.Context()
	start := make(chan struct{})
	errs := make(chan error, workers)

	var wg sync.WaitGroup
	for range workers {
		rule := NewEach(sharedRule).SkipOnEmpty()
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			errs <- rule.ValidateValue(ctx, []item{{Value: "value"}})
		}()
	}

	close(start)
	wg.Wait()
	close(errs)

	for err := range errs {
		require.NoError(t, err)
	}
}

func TestEach_ValidateValue_OptionsApplyToElements_NoError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		ctx   context.Context
		rule  *Each
		value any
	}{
		{
			name:  "skip on empty",
			ctx:   t.Context(),
			rule:  NewEach(NewStringLength(1, 2)).SkipOnEmpty(),
			value: []string{""},
		},
		{
			name:  "skip on error",
			ctx:   withPreviousRulesErrored(t.Context()),
			rule:  NewEach(NewRequired()).SkipOnError(),
			value: []string{""},
		},
		{
			name: "when false",
			ctx:  t.Context(),
			rule: NewEach(NewRequired()).When(func(_ context.Context, _ any) bool {
				return false
			}),
			value: []string{""},
		},
		{
			name:  "nested each skip on empty",
			ctx:   t.Context(),
			rule:  NewEach(NewEach(NewStringLength(1, 2))).SkipOnEmpty(),
			value: [][]string{{""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, tt.rule.ValidateValue(tt.ctx, tt.value))
		})
	}
}

func TestEach_ValidateValue_FirstValueIs0_Error(t *testing.T) {
	ctx := context.Background()
	err := NewEach(NewNumber(1, 2)).ValidateValue(ctx, []int{0})
	assert.Error(t, err)
	expectedErr := Result{
		errors: []*ValidationError{
			{
				Message: "Value must be no less than 1.",
				Params: map[string]any{
					"max": int64(2),
					"min": int64(1),
				},
				ValuePath: []string{"0"},
			},
		},
	}

	assert.Equal(t, expectedErr, err)
}
