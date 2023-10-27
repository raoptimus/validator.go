package validator

import (
	"context"
	"testing"
)

func BenchmarkValidatorNumber(b *testing.B) {
	ctx := context.Background()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := NewNumber(1, 3).ValidateValue(ctx, 2)
		if err != nil {
			b.Error(err)
		}
	}
}
