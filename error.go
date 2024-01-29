package validator

import (
	"errors"
)

var (
	NotExistsDataSetIntoContextError = errors.New("not exists data set into context")
	UnknownOperatorError             = errors.New("unknown operator")
	UnexpectedValueTypeError         = errors.New("unexpected value type")
)

type ValidationError struct {
	Message   string
	Params    map[string]any
	ValuePath []string
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		Message: message,
	}
}

func (v *ValidationError) Error() string {
	return v.Message
}

func (v *ValidationError) WithParams(params map[string]any) *ValidationError {
	v.Params = params

	return v
}

func (v *ValidationError) WithValuePath(valuePath []string) *ValidationError {
	v.ValuePath = valuePath

	return v
}

// ErrorMessagesIndexedByPath - проверяет на ошибку валидации и возвращает аттрибуты, где ключ равняется полю, а значения ошибкам валидации.
//
//	{
//		"client_id": [
//			"Value cannot be blank.",
//			"Value is invalid."
//		]
//	}
func ErrorMessagesIndexedByPath(err error) (map[string][]string, bool) {
	var result *Result
	if errors.As(err, &result) {
		return result.ErrorMessagesIndexedByPath(), true
	}

	return nil, false
}
