package validator

import (
	"errors"
	"reflect"

	"github.com/raoptimus/validator.go/rule"
)

var ErrUndefinedField = errors.New("undefined property")

type UndefinedFieldErr struct {
	DataSetName   string
	AttributeName string
}

func (u *UndefinedFieldErr) Error() string {
	return ErrUndefinedField.Error() + ": " + u.DataSetName + "." + u.AttributeName
}

func (u *UndefinedFieldErr) Unwrap() error {
	return ErrUndefinedField
}

func Validate(dataSet any, rules map[string][]RuleValidator, skipOnError bool) error {
	resultSet := rule.NewResultSet()

	pm := reflect.ValueOf(dataSet)
	vm := reflect.Indirect(pm)

	t := reflect.TypeOf(dataSet)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	var requiredIndex int

	for attr, validatorRules := range rules {
		value := vm.FieldByName(attr)
		if !value.IsValid() {
			return &UndefinedFieldErr{vm.Type().String(), attr}
		}

		// find required validator
		requiredIndex = -1
		for i, validator := range validatorRules {
			if _, ok := validator.(rule.Required); ok {
				requiredIndex = i
				break
			}
		}

		if value.Kind() == reflect.Pointer {
			if value.IsNil() {
				// if value is not required and is nil
				if requiredIndex == -1 {
					continue
				}
			} else {
				value = reflect.Indirect(value)
			}
		}

		fieldName := attr
		if field, ok := t.FieldByName(attr); ok {
			if v, ok := field.Tag.Lookup("json"); ok {
				fieldName = v
			}
		}

		if requiredIndex != -1 {
			required := validatorRules[requiredIndex]
			if _, ok := required.(rule.Required); ok {
				if err := required.ValidateValue(value.Interface()); err != nil {
					var errRes rule.Result
					if errors.As(err, &errRes) {
						resultSet = resultSet.WithResult(fieldName, errRes)
					}

					continue
				}
			}
		}

		for i, validator := range validatorRules {
			if requiredIndex == i {
				continue
			}

			if err := validator.ValidateValue(value.Interface()); err != nil {
				var errRes rule.Result
				if errors.As(err, &errRes) {
					resultSet = resultSet.WithResult(fieldName, errRes)
				}

				if skipOnError {
					break
				}
			}
		}
	}

	if resultSet.HasErrors() {
		return resultSet
	}
	return nil
}
