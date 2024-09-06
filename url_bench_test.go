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
