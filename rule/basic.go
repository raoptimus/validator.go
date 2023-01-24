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

func valueIsEmpty(value reflect.Value) bool {
	if !value.IsValid() || value.IsZero() {
		return true
	}

	if value.Kind() == reflect.String {
		if len(strings.Trim(value.String(), " ")) == 0 {
			return true
		}
	}

	return false
}
