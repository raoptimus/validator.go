/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package set

import (
	"errors"
	"fmt"
)

var ErrUndefinedField = errors.New("undefined field")

type UndefinedFieldError struct {
	dataSetName   string
	attributeName string
}

func NewUndefinedFieldError(dataSetPtr any, attributeName string) *UndefinedFieldError {
	return &UndefinedFieldError{
		dataSetName:   fmt.Sprintf("%T", dataSetPtr),
		attributeName: attributeName,
	}
}

func (u *UndefinedFieldError) Error() string {
	return ErrUndefinedField.Error() + ": " + u.dataSetName + "." + u.attributeName
}

func (u *UndefinedFieldError) Unwrap() error {
	return ErrUndefinedField
}
