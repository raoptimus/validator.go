package rule

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

func formatMessageWithArgs(message string, arguments map[string]any) string {
	for attr, valAny := range arguments {
		attrPlaceholder := "{" + attr + "}"
		if !strings.Contains(message, attrPlaceholder) {
			continue
		}

		message = strings.ReplaceAll(message, attrPlaceholder, fmt.Sprintf("%v", valAny))
	}

	return formatMessage(message)

}

func formatMessage(message string) string {
	return message
}

func valueIsEmpty(value reflect.Value, allowZeroValue bool) bool {
	if !value.IsValid() || (!allowZeroValue && value.IsZero()) {
		return true
	}

	kind := value.Kind()
	switch kind {
	case reflect.String:
		if len(strings.Trim(value.String(), " ")) == 0 {
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
	}

	return false
}

func toString(v any) (string, bool) {
	str, ok := v.(string)
	if !ok {
		i, ok := v.(fmt.Stringer)
		if !ok {
			return "", false
		}

		str = i.String()
	}

	return str, true
}
