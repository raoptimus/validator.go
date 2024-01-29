package validator

import (
	"context"
	"testing"
)

// BenchmarkValidatorRequired-10    	 6303788	       175.8 ns/op	      80 B/op	       4 allocs/op
// BenchmarkValidatorRequired-10         17474499          338.6 ns/op       184 B/op          7 allocs/op
func BenchmarkValidatorRequired(b *testing.B) {
	ctx := context.Background()
	dto := &testObject{Name: "test"}

	rules := RuleSet{
		"Name": {
			NewRequired().WithMessage("Required"),
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := Validate(ctx, dto, rules)
		if err != nil {
			b.Error(err)
		}
	}
}
