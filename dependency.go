package validator

import (
	"reflect"
)

type RuleValidator interface {
	//Validate(value reflect.Value, previousRulesErrored bool)
	ValidateValue(value reflect.Value) error
}
