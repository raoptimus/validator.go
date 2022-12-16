package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResultSet_GetResult_Successfully(t *testing.T) {
	rs := NewResultSet()
	rs = rs.WithResult("name", NewResult().WithError("test"))

	rsl := rs.GetResults()
	assert.Equal(t, map[string]Result{"name": {errors: []string{"test"}}}, rs.GetResults())

	rsl["invisible"] = NewResult().WithError("invisible error")
	assert.Equal(t, map[string]Result{"name": {errors: []string{"test"}}}, rs.GetResults())

	r, err := rs.GetResult("name")
	assert.NoError(t, err)

	rs = rs.WithResult("name", r.WithError("test2"))
	r = r.WithError("invisible error")
	assert.Equal(t, map[string]Result{"name": {errors: []string{"test", "test", "test2"}}}, rs.GetResults())
}

func TestResultSet_HasErrors_ReturnsTrue(t *testing.T) {
	rs := NewResultSet().WithResult("name", NewResult().WithError("test"))
	assert.True(t, rs.HasErrors())
}
func TestResultSet_HasErrors_ReturnsFalse(t *testing.T) {
	rs := NewResultSet().WithResult("name", NewResult())
	assert.False(t, rs.HasErrors())
}

func TestResult_GetErrors_Successfully(t *testing.T) {
	res := NewResult().WithError("test err")
	errs := res.GetErrors()
	assert.Equal(t, []string{"test err"}, errs)

	errs = append(errs, "invisible error")
	assert.Equal(t, []string{"test err"}, res.GetErrors())

	res = res.WithError("test2 err")
	assert.Equal(t, []string{"test err", "test2 err"}, res.GetErrors())
}

func TestResult_IsValid_True(t *testing.T) {
	res := NewResult()
	assert.True(t, res.IsValid())
}

func TestResult_IsValid_False(t *testing.T) {
	res := NewResult().WithError("test err")
	assert.False(t, res.IsValid())
}
