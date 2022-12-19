package validator

type RuleValidator interface {
	//Validate(value reflect.Value, previousRulesErrored bool)
	ValidateValue(value any) error
}
