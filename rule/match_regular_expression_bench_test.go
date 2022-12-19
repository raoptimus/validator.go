package rule

import (
	"testing"
)

func BenchmarkValidatorMatchRegularExpression(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := NewMatchRegularExpression(`[a-z]+`).ValidateValue("hello")
		if err != nil {
			b.Error(err)
		}
	}
}
