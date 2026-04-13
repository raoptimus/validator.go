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

func TestStringLength_ValidateValue(t *testing.T) {
	r := NewStringLength(2, 4)
	err := r.ValidateValue(context.Background(), "-")
	assert.Error(t, err)
}
