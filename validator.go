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

	for attr, r := range rules {
		value := reflect.Indirect(vm.FieldByName(attr))
		if !value.IsValid() {
			return &UndefinedFieldErr{vm.Type().String(), attr}
		}

		for _, validator := range r {
			if _, ok := validator.(rule.Required); ok {
				if err := validator.ValidateValue(value.Interface()); err != nil {
					var errRes rule.Result
					if errors.As(err, &errRes) {
						resultSet = resultSet.WithResult(attr, errRes)
					}

					if skipOnError {
						goto next
					}
					break
				}
			}
		}

		for _, validator := range r {
			if _, ok := validator.(rule.Required); ok {
				continue
			}
			if err := validator.ValidateValue(value.Interface()); err != nil {
				var errRes rule.Result
				if errors.As(err, &errRes) {
					resultSet = resultSet.WithResult(attr, errRes)
				}

				if skipOnError {
					break
				}
			}
		}
	next:
	}

	if resultSet.HasErrors() {
		return resultSet
	}
	return nil
}
