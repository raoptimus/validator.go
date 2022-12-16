package validator

import (
	"testing"

	"github.com/raoptimus/validator.go/rule"
)

func BenchmarkValidatorString(b *testing.B) {
	b.ReportAllocs()

	dto := &testObject{Name: ""}
	rules := map[string][]RuleValidator{
		"Name": {
			rule.NewStringLength(0, 0),
		},
	}

	for i := 0; i < b.N; i++ {
		err := Validate(dto, rules, false)
		if err != nil {
			b.Error(err)
		}
	}
}
