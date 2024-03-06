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

func NewEach(rules ...Rule) *Each {
	return &Each{
		message:               "Value is invalid",
		incorrectInputMessage: "Value must be array",
		rules:                 rules,
		normalizeRulesEnabled: true,
	}
}

func (r *Each) WithMessage(message string) *Each {
	rc := *r
	rc.message = message

	return &rc
}

func (r *Each) When(v WhenFunc) *Each {
	rc := *r
	rc.whenFunc = v

	return &rc
}

func (r *Each) when() WhenFunc {
	return r.whenFunc
}

func (r *Each) setWhen(v WhenFunc) {
	r.whenFunc = v
}

func (r *Each) SkipOnEmpty(v bool) *Each {
	rc := *r
	rc.skipEmpty = v

	return &rc
}

func (r *Each) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *Each) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *Each) WithIncorrectInputMessage(incorrectInputMessage string) *Each {
	rc := *r
	rc.incorrectInputMessage = incorrectInputMessage

	return &rc
}

func (r *Each) ValidateValue(ctx context.Context, value any) error {
	r.normalizeRules()

	result := NewResult()
	if reflect.TypeOf(value).Kind() != reflect.Slice {
		return result.WithError(
			NewValidationError(r.incorrectInputMessage).
				WithParams(map[string]any{
					//"attribute": "",//todo
					"value": value,
				}),
		)
	}

	vs := reflect.ValueOf(value)
	for i := 0; i < vs.Len(); i++ {
		v := vs.Index(i).Interface()

		if err := ValidateValue(ctx, v, r.rules...); err != nil {
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

func (r *Each) normalizeRules() {
	if !r.normalizeRulesEnabled {
		return
	}
	r.normalizeRulesEnabled = false

	for i, rule := range r.rules {
		if rse, ok := rule.(RuleSkipEmpty); ok {
			rse.setSkipOnEmpty(r.skipEmpty)
		}

		if rw, ok := rule.(RuleWhen); ok {
			rw.setWhen(r.whenFunc)
		}

		r.rules[i] = rule
	}
}
