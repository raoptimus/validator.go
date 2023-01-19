package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResultSet_Result_Successfully(t *testing.T) {
	rs := NewResultSet()
	rs = rs.WithResult("name", NewResult().WithError("test"))

	rsl := rs.Results()
	assert.Equal(t, map[string]Result{"name": {errors: []string{"test"}}}, rs.Results())

	rsl["invisible"] = NewResult().WithError("invisible error")
	assert.Equal(t, map[string]Result{"name": {errors: []string{"test"}}}, rs.Results())

	r, err := rs.Result("name")
	assert.NoError(t, err)

	rs = rs.WithResult("name", r.WithError("test2"))
	r = r.WithError("invisible error")
	assert.Equal(t, map[string]Result{"name": {errors: []string{"test", "test", "test2"}}}, rs.Results())
}

func TestResultSet_HasErrors_ReturnsTrue(t *testing.T) {
	rs := NewResultSet().WithResult("name", NewResult().WithError("test"))
	assert.True(t, rs.HasErrors())
}
func TestResultSet_HasErrors_ReturnsFalse(t *testing.T) {
	rs := NewResultSet().WithResult("name", NewResult())
	assert.False(t, rs.HasErrors())
}

func TestResult_Errors_Successfully(t *testing.T) {
	res := NewResult().WithError("test err")
	errs := res.Errors()
	assert.Equal(t, []string{"test err"}, errs)

	errs = append(errs, "invisible error")
	assert.Equal(t, []string{"test err"}, res.Errors())

	res = res.WithError("test2 err")
	assert.Equal(t, []string{"test err", "test2 err"}, res.Errors())
}

func TestResult_IsValid_True(t *testing.T) {
	res := NewResult()
	assert.True(t, res.IsValid())
}

func TestResult_IsValid_False(t *testing.T) {
	res := NewResult().WithError("test err")
	assert.False(t, res.IsValid())
}
