package validator

import (
	"context"
	"testing"

	"github.com/raoptimus/validator.go/rule"
)

func BenchmarkValidatorString(b *testing.B) {
	ctx := context.Background()

	dto := &testObject{Name: ""}
	rules := map[string][]RuleValidator{
		"Name": {
			rule.NewStringLength(0, 0),
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
