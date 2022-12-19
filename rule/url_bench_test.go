package rule

import (
	"testing"
)

func BenchmarkValidatorUrl(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := NewUrl().ValidateValue("https://example.com")
		if err != nil {
			b.Error(err)
		}
	}
}
