package rule

import (
	"testing"
)

func BenchmarkValidatorNumber(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := NewNumber(1, 3).ValidateValue(2)
		if err != nil {
			b.Error(err)
		}
	}
}
