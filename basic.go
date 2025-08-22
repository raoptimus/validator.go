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
	if !value.IsValid() {
		return true
	}

	switch value.Kind() {
	case reflect.Slice, reflect.Map:
		return value.IsNil() || value.Len() == 0
	case reflect.Array, reflect.Struct:
		return value.IsZero()
	case reflect.String:
		return len(strings.TrimSpace(value.String())) == 0
	case reflect.Ptr:
		return value.IsNil() || valueIsEmpty(value.Elem())
	default:
		if value.CanInterface() {
			if i, ok := value.Interface().(fmt.Stringer); ok {
				return len(strings.TrimSpace(i.String())) == 0
			}
		}

		return false
	}
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
