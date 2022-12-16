package rule

import (
	"errors"
	"strings"
	"sync"
)

var ErrResultNotFound = errors.New("result not found")

type Result struct {
	mu sync.Locker

	errors []string
}

func NewResult() *Result {
	return &Result{
		mu:     &sync.Mutex{},
		errors: make([]string, 0),
	}
}

func (s *Result) WithError(err string) *Result {
	s.AddError(err)
	return s
}

func (s *Result) Error() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	summary := strings.Builder{}
	for _, v := range s.errors {
		summary.WriteString(v)
		summary.WriteString(". ")
	}
	return strings.TrimRight(summary.String(), " ")
}

func (s *Result) IsValid() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.errors) == 0
}

func (s *Result) AddError(err string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.errors = append(s.errors, err)
}

func (s *Result) GetErrors() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := s.errors
	return r
}

type ResultSet struct {
	mu sync.Locker

	results map[string]*Result
}

func NewResultSet() *ResultSet {
	return &ResultSet{
		mu:      &sync.Mutex{},
		results: make(map[string]*Result),
	}
}

func (s *ResultSet) Error() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	summary := strings.Builder{}
	for _, v := range s.results {
		summary.WriteString(v.Error())
		summary.WriteString("\n")
	}
	return strings.TrimRight(summary.String(), "\n")
}

func (s *ResultSet) HasErrors() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.results) > 0
}

func (s *ResultSet) GetResult(attribute string) (*Result, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if r, ok := s.results[attribute]; !ok {
		return nil, ErrResultNotFound
	} else {
		return r, nil
	}
}

func (s *ResultSet) AddResult(attribute string, result *Result) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.results[attribute]

	if !ok {
		s.results[attribute] = result
		return
	}

	if r.IsValid() {
		return
	}

	for _, err := range result.GetErrors() {
		r.AddError(err)
	}
}

func (s *ResultSet) GetResults() map[string]*Result {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := s.results
	return r
}

func (s *ResultSet) GetResultsErrors() map[string][]string {
	s.mu.Lock()
	defer s.mu.Unlock()

	ret := make(map[string][]string)
	for attr, r := range s.results {
		ret[attr] = r.GetErrors()
	}

	return ret
}
