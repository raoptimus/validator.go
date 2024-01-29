package validator

import (
	"context"
	"testing"
)

func BenchmarkValidatorUrl(b *testing.B) {
	ctx := context.Background()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := NewURL().ValidateValue(ctx, "https://example.com")
		if err != nil {
			b.Error(err)
		}
	}
}
