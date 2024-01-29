package validator

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/raoptimus/validator.go/set"
)

const (
	separator    = "."
	eachShortcut = "*"
)

type Nested struct {
	normalizeRulesEnabled bool
	rules                 RuleSet
	message               string
}

func NewNested(rules RuleSet) Nested {
	return Nested{
		normalizeRulesEnabled: true,
		rules:                 rules,
		message:               "",
	}
}

func (n Nested) WithMessage(message string) Nested {
	n.message = message
	return n
}

func (n Nested) notNormalizeRules() Nested {
	n.normalizeRulesEnabled = false

	return n
}

// ValidateValue todo: make me
func (n Nested) ValidateValue(ctx context.Context, value any) error {
	if n.normalizeRulesEnabled {
		n.normalizeRulesEnabled = false // once
		if rules, err := n.normalizeRules(); err != nil {
			return err
		} else {
			n.rules = rules
		}
	}

	vt := reflect.TypeOf(value)
	if vt.Kind() == reflect.Pointer {
		vt = vt.Elem()
	}

	if len(n.rules) == 0 {
		if vt.Kind() != reflect.Struct {
			return fmt.Errorf("nested rule without rules could be used for structs only. %s given",
				vt.Kind().String(),
			)
		}

		var err error
		data, ok := value.(*set.DataSetStruct)
		if !ok {
			data, err = set.NewDataSetStruct(value)
			if err != nil {
				return err
			}
		}

		return Validate(ctx, data, n.rules)
	}

	if vt.Kind() != reflect.Struct {
		return NewResult().WithError(
			NewValidationError(fmt.Sprintf("value should be a struct. %T given.", value)).
				WithParams(map[string]any{
					"attribute": "", // todo: get attribute
					"value":     value,
				}),
		)
	}

	var err error
	data, ok := value.(*set.DataSetStruct)
	if !ok {
		data, err = set.NewDataSetStruct(value)
		if err != nil {
			return err
		}
	}

	compoundResult := NewResult()
	results := make([]Result, 0, len(n.rules))

	for fieldName, rules := range n.rules {
		// todo: parse valuePath

		validatedValue, err := data.FieldValue(fieldName)
		if err != nil { // todo: check after parsed
			return err
		}
		valuePath := data.FieldAliasName(fieldName)

		if err := ValidateValue(ctx, validatedValue, rules...); err != nil {
			var itemResult Result

			if errors.As(err, &itemResult) {
				result := NewResult()
				for _, itemError := range itemResult.Errors() {
					var errorValuePath []string
					if _, err := strconv.Atoi(valuePath); err != nil {
						errorValuePath = strings.Split(valuePath, separator)
					} else {
						errorValuePath = []string{valuePath}
					}
					if len(itemError.ValuePath) > 0 {
						errorValuePath = append(errorValuePath, itemError.ValuePath...)
					}
					itemError.ValuePath = errorValuePath
					result = result.WithError(itemError)
				}

				results = append(results, result)
				continue
			}

			return err
		}
	}

	for i := range results {
		compoundResult = compoundResult.WithError(results[i].Errors()...)
	}

	if !compoundResult.IsValid() {
		return compoundResult
	}

	return nil
}

func (n Nested) normalizeRules() (RuleSet, error) {
	nRules := n.rules

	for {
		rulesMap := make(map[string]RuleSet, len(nRules))
		needBreak := true

		for valuePath, rules := range nRules {
			if valuePath == eachShortcut {
				return nil, errors.New("bare shortcut is prohibited. Use 'Each' rule instead")
			}
			if valuePath == "" {
				continue
			}
			parts := strings.Split(valuePath, separator)
			if len(parts) == 1 {
				continue
			}

			needBreak = false

			lastValuePath := parts[len(parts)-1]
			remainingValuePath := strings.Join(parts, eachShortcut)
			remainingValuePath = strings.TrimRight(remainingValuePath, separator)
			if _, ok := rulesMap[remainingValuePath]; !ok {
				if _, ok := rulesMap[remainingValuePath]; ok {
					rulesMap[remainingValuePath][lastValuePath] = rules
				} else {
					rulesMap[remainingValuePath] = RuleSet{lastValuePath: rules}
				}

				delete(nRules, valuePath)
			}
		}

		for valuePath, nestedRules := range rulesMap {
			nRules[valuePath] = []Rule{
				NewEach(NewNested(nestedRules).notNormalizeRules()),
			}
		}

		if needBreak {
			break
		}
	}

	return nRules, nil
}
