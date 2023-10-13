package validator

import (
	"context"
	"testing"

	"github.com/raoptimus/validator.go/rule"
)

func BenchmarkValidatorRequired(b *testing.B) {
	ctx := context.Background()
	dto := &testObject{Name: "test"}
	rules := map[string][]RuleValidator{
		"Name": {
			rule.NewRequired().WithMessage("Required"),
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := Validate(ctx, dto, rules, false)
		if err != nil {
			b.Error(err)
		}
	}
}
