package validator

import (
	"reflect"

	"github.com/raoptimus/validator.go/rule"
)

type RuleValidator interface {
	//Validate(value reflect.Value, previousRulesErrored bool)
	ValidateValue(value reflect.Value) *rule.Result
}
