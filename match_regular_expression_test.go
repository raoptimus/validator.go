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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchRegularExpression_ValidateValue_ValueNotString(t *testing.T) {
	ctx := context.Background()
	err := NewMatchRegularExpression(``).ValidateValue(ctx, 1)
	assert.Error(t, err)
	assert.Equal(t, "Value is invalid.", err.Error())
}
