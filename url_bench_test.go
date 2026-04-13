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

func BenchmarkValidatorUrl(b *testing.B) {
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := NewURL().ValidateValue(ctx, "https://example.com")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkValidatorDeepLinkURL(b *testing.B) {
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := NewDeepLinkURL().ValidateValue(ctx, "tg:resolve?domain={domain}")
		if err != nil {
			b.Error(err)
		}
	}
}
