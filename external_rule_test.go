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
