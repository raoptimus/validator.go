package rule

import (
	"errors"
	"strings"
)

var ErrResultNotFound = errors.New("result not found")

type Result struct {
	errors []string
}

var emptyResult = Result{}

func NewResult() Result {
	return Result{
		errors: make([]string, 0),
	}
}

func (s Result) WithError(err string) Result {
	s.errors = append(s.errors, err)
	return s
}

func (s Result) Error() string {
	if len(s.errors) == 1 {
		return s.errors[0]
	}

	summary := strings.Builder{}
	for _, v := range s.errors {
		summary.WriteString(v)
		summary.WriteString(". ")
	}

	return strings.TrimRight(summary.String(), " ")
}

func (s Result) IsValid() bool {
	return len(s.errors) == 0
}

func (s Result) Errors() []string {
	r := s.errors
	return r
}

type ResultSet struct {
	results map[string]Result
}

func NewResultSet() ResultSet {
	return ResultSet{
		results: make(map[string]Result),
	}
}

func (s ResultSet) Error() string {
	var summary strings.Builder
	for k, v := range s.results {
		summary.WriteString(k)
		summary.WriteString(": ")
		summary.WriteString(v.Error())
		summary.WriteString("\n")
	}
	return strings.TrimRight(summary.String(), "\n")
}

func (s ResultSet) HasErrors() bool {
	return len(s.results) > 0
}

func (s ResultSet) Result(attribute string) (Result, error) {
	if r, ok := s.results[attribute]; !ok {
		return emptyResult, ErrResultNotFound
	} else {
		return r, nil
	}
}

func (s ResultSet) WithResult(attribute string, result Result) ResultSet {
	if result.IsValid() {
		return s
	}

	res, ok := s.results[attribute]

	if !ok {
		s.results[attribute] = result
		return s
	}

	if res.IsValid() {
		return s
	}

	for _, err := range result.Errors() {
		res = res.WithError(err)
	}

	s.results[attribute] = res
	return s
}

func (s ResultSet) Results() map[string]Result {
	ret := make(map[string]Result)
	for attr, res := range s.results {
		ret[attr] = res
	}
	return ret
}

func (s ResultSet) ResultErrors() map[string][]string {
	ret := make(map[string][]string)
	for attr, r := range s.results {
		ret[attr] = r.Errors()
	}

	return ret
}
