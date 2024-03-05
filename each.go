package validator

import (
	"context"
	"errors"
	"reflect"
	"strconv"
)

type Each struct {
	message               string
	incorrectInputMessage string
	rules                 Rules
	normalizeRulesEnabled bool
	whenFunc              WhenFunc
	skipEmpty             bool
}

func NewEach(rules ...Rule) Each {
	return Each{
		message:               "Value is invalid",
		incorrectInputMessage: "Value must be array",
		rules:                 rules,
	}
}

func (e Each) WithMessage(message string) Each {
	e.message = message

	return e
}

func (e Each) When(v WhenFunc) Each {
	e.whenFunc = v

	return e
}

func (e Each) when() WhenFunc {
	return e.whenFunc
}

func (e Each) setWhen(v WhenFunc) {
	e.whenFunc = v
}

func (e Each) SkipOnEmpty(v bool) Each {
	e.skipEmpty = v
	return e
}

func (e Each) skipOnEmpty() bool {
	return e.skipEmpty
}

func (e Each) setSkipOnEmpty(v bool) {
	e.skipEmpty = v // because copy
}

func (e Each) WithIncorrectInputMessage(incorrectInputMessage string) Each {
	e.incorrectInputMessage = incorrectInputMessage

	return e
}

func (e Each) ValidateValue(ctx context.Context, value any) error {
	e.normalizeRules()

	result := NewResult()
	if reflect.TypeOf(value).Kind() != reflect.Slice {
		return result.WithError(
			NewValidationError(e.incorrectInputMessage).
				WithParams(map[string]any{
					//"attribute": "",//todo
					"value": value,
				}),
		)
	}

	vs := reflect.ValueOf(value)
	for i := 0; i < vs.Len(); i++ {
		v := vs.Index(i).Interface()

		if err := ValidateValue(ctx, v, e.rules...); err != nil {
			var r Result
			if errors.As(err, &r) {
				for _, err := range r.Errors() {
					valuePath := []string{strconv.Itoa(i)}
					if len(err.ValuePath) > 0 {
						valuePath = append(valuePath, err.ValuePath...)
					}
					err.ValuePath = valuePath
					result = result.WithError(err)
				}

				continue
			}

			return err
		}
	}

	if result.IsValid() {
		return nil
	}

	return result
}

func (e Each) normalizeRules() {
	if !e.normalizeRulesEnabled {
		return
	}
	e.normalizeRulesEnabled = false // once

	for i, r := range e.rules {
		if rse, ok := r.(RuleSkipEmpty); ok {
			rse.setSkipOnEmpty(e.skipEmpty)
		}

		if rw, ok := r.(RuleWhen); ok {
			rw.setWhen(e.whenFunc)
		}

		e.rules[i] = r
	}
}
