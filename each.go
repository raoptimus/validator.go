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
	"reflect"
	"strconv"
)

type Each struct {
	message               string
	incorrectInputMessage string
	rules                 Rules
	whenFunc              WhenFunc
	skipEmpty             bool
	skipError             bool
}

func NewEach(rules ...Rule) *Each {
	return &Each{
		message:               MessageInvalid,
		incorrectInputMessage: MessageEachIncorrectInput,
		rules:                 rules,
	}
}

func (r *Each) WithMessage(message string) *Each {
	rc := *r
	rc.message = message

	return &rc
}

func (r *Each) WithIncorrectInputMessage(incorrectInputMessage string) *Each {
	rc := *r
	rc.incorrectInputMessage = incorrectInputMessage

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

func (r *Each) SkipOnEmpty() *Each {
	rc := *r
	rc.skipEmpty = true

	return &rc
}

func (r *Each) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *Each) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *Each) SkipOnError() *Each {
	rs := *r
	rs.skipError = true

	return &rs
}

func (r *Each) shouldSkipOnError() bool {
	return r.skipError
}
func (r *Each) setSkipOnError(v bool) {
	r.skipError = v
}

func (r *Each) ValidateValue(ctx context.Context, value any) error {
	result := NewResult()
	if value == nil || reflect.TypeOf(value).Kind() != reflect.Slice {
		return result.WithError(
			NewValidationError(r.incorrectInputMessage).
				WithParams(map[string]any{
					// "attribute": "",//todo
					"value": value,
				}),
		)
	}

	overrides := &ruleOverrides{
		whenFunc:  r.whenFunc,
		skipEmpty: r.skipEmpty,
		skipError: r.skipError,
	}
	vs := reflect.ValueOf(value)
	for i := 0; i < vs.Len(); i++ {
		v := vs.Index(i).Interface()

		if err := validateValue(ctx, v, overrides, r.rules...); err != nil {
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
