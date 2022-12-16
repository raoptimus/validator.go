package validator

import (
	"testing"

	"github.com/raoptimus/validator.go/rule"
)

func BenchmarkValidatorRequired(b *testing.B) {
	b.ReportAllocs()
	dto := &testObject{Name: "test"}
	rules := map[string][]RuleValidator{
		"Name": {
			rule.NewRequired().WithMessage("Required"),
		},
	}

	for i := 0; i < b.N; i++ {
		err := Validate(dto, rules, false)
		if err != nil {
			b.Error(err)
		}
	}
}
