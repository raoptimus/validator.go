package validator

import (
	"context"
	"testing"
)

// BenchmarkValidatorRequired-10    	 6303788	       175.8 ns/op	      80 B/op	       4 allocs/op
// BenchmarkValidatorRequired-10    	 1566454	       757.4 ns/op	    1200 B/op	      18 allocs/op
// BenchmarkValidatorRequired-10         12395048          276.4 ns/op       184 B/op          7 allocs/op
// BenchmarkValidatorRequired-10         3717618           965.3 ns/op      1176 B/op         18 allocs/op
// BenchmarkValidatorRequired-10         3617426           1010 ns/op       1224 B/op         19 allocs/op
// BenchmarkValidatorRequired-10    	 2317490	       517.0 ns/op	     280 B/op	      10 allocs/op
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
