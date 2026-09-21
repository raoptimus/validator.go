/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator_test

import (
	"context"
	"errors"
	"testing"

	validator "github.com/raoptimus/validator.go/v2"
	"github.com/stretchr/testify/require"
)

var errExternalValidateValue = errors.New("external ValidateValue called")

type externalEach struct {
	*validator.Each
}

func (*externalEach) ValidateValue(context.Context, any) error {
	return errExternalValidateValue
}

type delegatingEach struct {
	*validator.Each
}

func (r *delegatingEach) ValidateValue(ctx context.Context, value any) error {
	return r.Each.ValidateValue(ctx, value)
}

func TestEach_ExternalDelegatingEachWrapperPropagatesOptions(t *testing.T) {
	rule := &delegatingEach{Each: validator.NewEach(validator.NewStringLength(1, 2))}

	err := validator.NewEach(rule).SkipOnEmpty().ValidateValue(t.Context(), [][]string{{""}})

	require.NoError(t, err)
}

func TestValidateValue_ExternalEachWrapperUsesPublicValidateValue(t *testing.T) {
	rule := &externalEach{Each: validator.NewEach()}

	err := validator.ValidateValue(t.Context(), "value", rule)

	require.ErrorIs(t, err, errExternalValidateValue)
}

func TestEach_ExternalEachWrapperUsesPublicValidateValueAsChild(t *testing.T) {
	rule := &externalEach{Each: validator.NewEach()}

	err := validator.NewEach(rule).ValidateValue(t.Context(), []string{"value"})

	require.ErrorIs(t, err, errExternalValidateValue)
}
