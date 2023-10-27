package validator

import (
	"context"
	"testing"
)

func BenchmarkValidatorString(b *testing.B) {
	ctx := context.Background()

	dto := &testObject{Name: ""}
	rules := RuleSet{
		"Name": {
			NewStringLength(0, 0),
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
