package validator

import (
	"fmt"
	"reflect"
	"strings"
)

func indirectValue(v any) (value any, isValid bool) {
	val := reflect.ValueOf(v)
	if !val.IsValid() || (val.Kind() == reflect.Ptr && val.IsNil()) {
		return nil, false
	}

	val = reflect.Indirect(val)
	if !val.IsValid() {
		return nil, false
	}

	return val.Interface(), true
}

func valueIsEmpty(value reflect.Value) bool {
	if !value.IsValid() || value.IsZero() {
		return true
	}

	kind := value.Kind()
	switch kind {
	case reflect.String:
		if len(strings.TrimSpace(value.String())) == 0 {
			return true
		}
	case reflect.Slice:
		if value.Len() == 0 {
			return true
		}
	case reflect.Ptr:
		if value.IsNil() {
			return true
		}

		return valueIsEmpty(value.Elem())
	}

	return false
}

func toString(v any) (string, bool) {
	if v == nil {
		return "", false
	}

	switch v := v.(type) {
	case string:
		return v, true
	case *string:
		return *v, true
	}

	if i, ok := v.(fmt.Stringer); ok {
		return i.String(), true
	}

	return "", false
}
