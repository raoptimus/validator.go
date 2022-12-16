package rule

import (
	"fmt"
	"reflect"
	"strings"
)

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
	value = extractValue(value)

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

func extractValue(value reflect.Value) reflect.Value {
	if !value.IsValid() {
		return value
	}

	if value.Kind() == reflect.Pointer {
		return extractValue(reflect.Indirect(value))
	}

	return value
}
