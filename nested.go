/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/raoptimus/validator.go/v2/set"
)

const (
	separator      = "."
	NestedShortcut = "*"
)

type Nested struct {
	// normalizeOnce guards lazy rule normalization so that concurrent
	// ValidateValue calls on a shared *Nested are safe. Stored as a pointer
	// because the builder methods (WithMessage, SkipOnEmpty, etc.) copy the
	// struct — sync.Once must not be copied after first use.
	normalizeOnce *sync.Once
	normalizeErr  error
	rules         RuleSet
	message       string
	whenFunc      WhenFunc
	skipEmpty     bool
	skipError     bool
}

func NewNested(rules RuleSet) *Nested {
	return &Nested{
		normalizeOnce: &sync.Once{},
		rules:         rules,
	}
}

func (r *Nested) WithMessage(message string) *Nested {
	rc := *r
	rc.message = message

	return &rc
}

func (r *Nested) When(v WhenFunc) *Nested {
	rc := *r
	rc.whenFunc = v

	return &rc
}

func (r *Nested) when() WhenFunc {
	return r.whenFunc
}

func (r *Nested) setWhen(v WhenFunc) {
	r.whenFunc = v
}

func (r *Nested) SkipOnEmpty() *Nested {
	rc := *r
	rc.skipEmpty = true

	return &rc
}

func (r *Nested) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *Nested) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *Nested) notNormalizeRules() *Nested {
	rc := *r
	rc.normalizeOnce.Do(func() {}) // mark as already normalized

	return &rc
}

func (r *Nested) SkipOnError() *Nested {
	rs := *r
	rs.skipError = true

	return &rs
}

func (r *Nested) shouldSkipOnError() bool {
	return r.skipError
}
func (r *Nested) setSkipOnError(v bool) {
	r.skipError = v
}

func (r *Nested) ValidateValue(ctx context.Context, value any) error {
	r.normalizeOnce.Do(func() {
		rules, err := r.normalizeRules()
		if err != nil {
			r.normalizeErr = err

			return
		}

		r.rules = rules
	})

	if r.normalizeErr != nil {
		return r.normalizeErr
	}

	if value == nil {
		return NewResult().WithError(
			NewValidationError(fmt.Sprintf("value should be a struct. %T given.", value)),
		)
	}

	vt := reflect.TypeOf(value)
	if vt.Kind() == reflect.Pointer {
		vt = vt.Elem()
	}

	if len(r.rules) == 0 {
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

		return Validate(ctx, data, r.rules)
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

	ctx = contextWithNestedDataSet(ctx, value)

	var err error
	data, ok := value.(*set.DataSetStruct)
	if !ok {
		data, err = set.NewDataSetStruct(value)
		if err != nil {
			return err
		}
	}

	compoundResult := NewResult()
	results := make([]Result, 0, len(r.rules))

	fieldNames := make([]string, 0, len(r.rules))
	for fieldName := range r.rules {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)

	for _, fieldName := range fieldNames {
		rules := r.rules[fieldName]
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

func (r *Nested) normalizeRules() (RuleSet, error) {
	nRules := r.rules

	for {
		rulesMap := make(map[string]RuleSet, len(nRules))
		needBreak := true

		for valuePath, rules := range nRules {
			if valuePath == NestedShortcut {
				return nil, errors.New("bare shortcut is prohibited. Use 'Nested' rule instead")
			}
			if valuePath == "" {
				continue
			}
			idx := strings.LastIndex(valuePath, separator)
			if idx < 0 {
				continue
			}

			needBreak = false

			lastValuePath := valuePath[idx+1:]
			remainingValuePath := strings.ReplaceAll(valuePath[:idx], separator, NestedShortcut)
			if _, ok := rulesMap[remainingValuePath]; !ok {
				rulesMap[remainingValuePath] = RuleSet{lastValuePath: rules}
			} else {
				rulesMap[remainingValuePath][lastValuePath] = rules
			}

			delete(nRules, valuePath)
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
