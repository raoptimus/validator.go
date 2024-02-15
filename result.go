package validator

import (
	"strings"
)

type Result struct {
	errors []*ValidationError
}

func NewResult() Result {
	return Result{
		errors: make([]*ValidationError, 0),
	}
}

func (s Result) WithError(errs ...*ValidationError) Result {
	s.errors = append(s.errors, errs...)
	return s
}

func (s Result) Error() string {
	if len(s.errors) == 1 {
		return s.errors[0].Message
	}

	var summary strings.Builder
	for _, v := range s.errors {
		if len(v.ValuePath) > 0 {
			summary.WriteString(strings.Join(v.ValuePath, "."))
			summary.WriteString(": ")
		}

		summary.WriteString(strings.TrimRight(v.Message, "."))
		summary.WriteString(". ")
	}

	return strings.TrimRight(summary.String(), " ")
}

func (s Result) IsValid() bool {
	return len(s.errors) == 0
}

func (s Result) Errors() []*ValidationError {
	r := s.errors
	return r
}

func (s Result) ErrorMessagesIndexedByPath() map[string][]string {
	errList := make(map[string][]string)
	for _, err := range s.errors {
		stringValuePath := strings.Join(err.ValuePath, separator)

		if _, ok := errList[stringValuePath]; !ok {
			errList[stringValuePath] = []string{err.Message}
		} else {
			errList[stringValuePath] = append(errList[stringValuePath], err.Message)
		}
	}

	return errList
}

func (s Result) AttributeErrorMessagesIndexedByPath(attribute string) map[string][]string {
	errList := make(map[string][]string)
	for _, err := range s.errors {
		var first string
		if len(err.ValuePath) > 0 {
			first = err.ValuePath[0]
		}
		if first != attribute {
			continue
		}
		stringValuePath := strings.Join(err.ValuePath[1:], separator)
		if _, ok := errList[stringValuePath]; !ok {
			errList[stringValuePath] = []string{err.Message}
		} else {
			errList[stringValuePath] = append(errList[stringValuePath], err.Message)
		}
	}

	return errList
}

func (s Result) CommonErrorMessages() []string {
	errList := make([]string, 0, len(s.errors))
	for _, err := range s.errors {
		var first string
		if len(err.ValuePath) > 0 {
			first = err.ValuePath[0]
		}
		if first != "" {
			continue
		}
		errList = append(errList, err.Message)
	}

	return errList
}
