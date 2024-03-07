package validator

import (
	"context"
	"errors"
	"reflect"

	"github.com/raoptimus/validator.go/v2/set"
)

func ValidateValue(ctx context.Context, value any, rules ...Rule) error {
	if len(rules) == 0 {
		return nil
	}

	if value == nil {
		requiredRule, ok := hasRequiredRule(rules)
		if !ok {
			return nil
		}

		return requiredRule.ValidateValue(ctx, value)
	}

	dataSet, err := normalizeDataSet(value)
	if err != nil {
		return err
	}

	if extDS, ok := ExtractDataSet[DataSet](ctx); !ok || value != extDS {
		ctx = withDataSet(ctx, dataSet)
	}

	rules = normalizeRules(rules)
	result := NewResult()

	for _, r := range rules {
		if err := r.ValidateValue(ctx, value); err != nil {
			var errRes Result
			if errors.As(err, &errRes) {
				result = result.WithError(errRes.Errors()...)

				continue
			}

			return err
		}
	}

	if result.IsValid() {
		return nil
	}

	for _, err := range result.Errors() {
		err.Message = DefaultTranslator.Translate(ctx, err.Message, err.Params)
	}

	return result
}

func Validate(ctx context.Context, dataSet any, rules RuleSet) error {
	normalizedDS, err := normalizeDataSet(dataSet) // 2 allocs
	if err != nil {
		return err
	}

	ctx = withDataSet(ctx, normalizedDS)
	results := make([]Result, 0, len(rules))

	for field, fieldRules := range rules {
		fieldValue, err := normalizedDS.FieldValue(field) // 2 allocs
		if err != nil {
			return err
		}
		aliasFieldName := normalizedDS.FieldAliasName(field)

		result := NewResult()
		fieldRules = normalizeRules(fieldRules)

		for _, validatorRule := range fieldRules {
			if isSkipValidate(ctx, fieldValue, validatorRule) {
				continue
			}

			if err := validatorRule.ValidateValue(ctx, fieldValue); err != nil {
				ctx = withPreviousRulesErrored(ctx)

				var errRes Result
				if errors.As(err, &errRes) {
					for _, rErr := range errRes.Errors() {
						if aliasFieldName != "" {
							valuePath := make([]string, 0, len(rErr.ValuePath)+1)
							valuePath = append(valuePath, aliasFieldName)
							valuePath = append(valuePath, rErr.ValuePath...)
							rErr.ValuePath = valuePath
						}
						result = result.WithError(rErr)
					}
				} else {
					return err
				}
			}
		}

		results = append(results, result)
	}

	summaryResult := NewResult()
	for i := range results {
		errs := (&results[i]).Errors()
		for _, err := range errs {
			err.Message = DefaultTranslator.Translate(ctx, err.Message, err.Params)
			summaryResult = summaryResult.WithError(err)
			// summaryResult = summaryResult.WithError(
			//	NewValidationError(DefaultTranslator.Translate(ctx, err.Message, err.Params)).
			//		WithParams(err.Params).
			//		WithValuePath(err.ValuePath),
			// )
		}
	}

	if !summaryResult.IsValid() {
		return summaryResult
	}

	return nil
}

func normalizeDataSet(ds any) (DataSet, error) {
	if ds == nil {
		return set.NewDataSetAny(ds), nil
	}

	rt := reflect.TypeOf(ds)
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}

	switch rt.Kind() {
	case reflect.Struct:
		if v, ok := ds.(DataSet); ok {
			return v, nil
		}

		return set.NewDataSetStruct(ds)
	case reflect.Map:
		if v, ok := ds.(map[string]any); ok {
			return set.NewDataSetMap(v), nil
		}
	}

	return set.NewDataSetAny(ds), nil
}

func normalizeRules(rules []Rule) []Rule {
	if len(rules) <= 1 {
		return rules
	}

	for i := range rules {
		if r, ok := rules[i].(*Required); ok {
			if i == 0 {
				break
			}
			ret := make([]Rule, 0, len(rules))
			ret = append(ret, r)
			ret = append(ret, rules[:i]...)
			ret = append(ret, rules[i:]...)

			return ret
		}
	}

	return rules
}

func isSkipValidate(ctx context.Context, value any, r Rule) bool {
	if rse, ok := r.(RuleSkipEmpty); ok {
		if rse.skipOnEmpty() && valueIsEmpty(reflect.ValueOf(value), false) {
			return true
		}
	}

	if rser, ok := r.(RuleSkipError); ok {
		if rser.shouldSkipOnError() && previousRulesErrored(ctx) {
			return true
		}
	}

	if rw, ok := r.(RuleWhen); ok {
		return rw.when() != nil && !rw.when()(ctx, value)
	}

	return false
}
