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
)

func BenchmarkValidatorMatchRegularExpression(b *testing.B) {
	ctx := context.Background()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := NewMatchRegularExpression(`[a-z]+`).ValidateValue(ctx, "hello")
		if err != nil {
			b.Error(err)
		}
	}
}
