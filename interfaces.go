package validator

type RuleValidator interface {
	// ValidateValue - (value reflect.Value, previousRulesErrored bool)
	ValidateValue(value any) error
}
