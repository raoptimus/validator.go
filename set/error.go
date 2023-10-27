package set

import (
	"errors"
	"fmt"
)

var BaseUndefinedFieldError = errors.New("undefined field")

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
	return BaseUndefinedFieldError.Error() + ": " + u.dataSetName + "." + u.attributeName
}

func (u *UndefinedFieldError) Unwrap() error {
	return BaseUndefinedFieldError
}
