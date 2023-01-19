package validator

import (
	"errors"

	"github.com/raoptimus/validator.go/rule"
)

// IsError - проверяет на ошибку валидации и возвращает аттрибуты, где ключ равняется полю, а значения ошибкам валидации.
//
//	{
//		"client_id": [
//			"Value cannot be blank.",
//			"Value is invalid."
//		]
//	}
func IsError(err error) (map[string][]string, bool) {
	var attributes rule.ResultSet
	if errors.As(err, &attributes) {
		return attributes.ResultErrors(), true
	}

	return nil, false
}
