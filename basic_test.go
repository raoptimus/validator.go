/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- indirectValue tests ---

func TestIndirectValue_Nil_ReturnsInvalid(t *testing.T) {
	val, ok := indirectValue(nil)
	assert.False(t, ok)
	assert.Nil(t, val)
}

func TestIndirectValue_ValidInt_ReturnsValue(t *testing.T) {
	val, ok := indirectValue(42)
	assert.True(t, ok)
	assert.Equal(t, 42, val)
}

func TestIndirectValue_ValidString_ReturnsValue(t *testing.T) {
	val, ok := indirectValue("hello")
	assert.True(t, ok)
	assert.Equal(t, "hello", val)
}

func TestIndirectValue_ValidPointer_ReturnsDereferencedValue(t *testing.T) {
	n := 99
	val, ok := indirectValue(&n)
	assert.True(t, ok)
	assert.Equal(t, 99, val)
}

func TestIndirectValue_NilPointer_ReturnsInvalid(t *testing.T) {
	var p *int
	val, ok := indirectValue(p)
	assert.False(t, ok)
	assert.Nil(t, val)
}

func TestIndirectValue_PointerToPointer_ReturnsDereferencedOnce(t *testing.T) {
	n := 7
	p := &n
	val, ok := indirectValue(&p)
	assert.True(t, ok)
	// reflect.Indirect dereferences one level, so we get *int
	assert.Equal(t, &n, val)
}

func TestIndirectValue_ZeroStruct_ReturnsValue(t *testing.T) {
	type dummy struct{ X int }
	val, ok := indirectValue(dummy{X: 0})
	assert.True(t, ok)
	assert.Equal(t, dummy{X: 0}, val)
}

// --- valueIsEmpty tests ---

func TestValueIsEmpty_InvalidValue_ReturnsTrue(t *testing.T) {
	var v reflect.Value
	assert.True(t, valueIsEmpty(v))
}

func TestValueIsEmpty_NilSlice_ReturnsTrue(t *testing.T) {
	var s []int
	assert.True(t, valueIsEmpty(reflect.ValueOf(s)))
}

func TestValueIsEmpty_EmptySlice_ReturnsTrue(t *testing.T) {
	s := make([]int, 0)
	assert.True(t, valueIsEmpty(reflect.ValueOf(s)))
}

func TestValueIsEmpty_NonEmptySlice_ReturnsFalse(t *testing.T) {
	s := []int{1, 2, 3}
	assert.False(t, valueIsEmpty(reflect.ValueOf(s)))
}

func TestValueIsEmpty_NilMap_ReturnsTrue(t *testing.T) {
	var m map[string]int
	assert.True(t, valueIsEmpty(reflect.ValueOf(m)))
}

func TestValueIsEmpty_EmptyMap_ReturnsTrue(t *testing.T) {
	m := make(map[string]int)
	assert.True(t, valueIsEmpty(reflect.ValueOf(m)))
}

func TestValueIsEmpty_NonEmptyMap_ReturnsFalse(t *testing.T) {
	m := map[string]int{"a": 1}
	assert.False(t, valueIsEmpty(reflect.ValueOf(m)))
}

func TestValueIsEmpty_EmptyString_ReturnsTrue(t *testing.T) {
	assert.True(t, valueIsEmpty(reflect.ValueOf("")))
}

func TestValueIsEmpty_WhitespaceString_ReturnsTrue(t *testing.T) {
	assert.True(t, valueIsEmpty(reflect.ValueOf("   \t\n")))
}

func TestValueIsEmpty_NonEmptyString_ReturnsFalse(t *testing.T) {
	assert.False(t, valueIsEmpty(reflect.ValueOf("hello")))
}

func TestValueIsEmpty_ZeroStruct_ReturnsTrue(t *testing.T) {
	type dummy struct {
		X int
		Y string
	}
	assert.True(t, valueIsEmpty(reflect.ValueOf(dummy{})))
}

func TestValueIsEmpty_NonZeroStruct_ReturnsFalse(t *testing.T) {
	type dummy struct {
		X int
		Y string
	}
	assert.False(t, valueIsEmpty(reflect.ValueOf(dummy{X: 1})))
}

func TestValueIsEmpty_ZeroArray_ReturnsTrue(t *testing.T) {
	var arr [3]int
	assert.True(t, valueIsEmpty(reflect.ValueOf(arr)))
}

func TestValueIsEmpty_NonZeroArray_ReturnsFalse(t *testing.T) {
	arr := [3]int{0, 1, 0}
	assert.False(t, valueIsEmpty(reflect.ValueOf(arr)))
}

func TestValueIsEmpty_NilPointer_ReturnsTrue(t *testing.T) {
	var p *int
	assert.True(t, valueIsEmpty(reflect.ValueOf(p)))
}

func TestValueIsEmpty_PointerToEmpty_ReturnsTrue(t *testing.T) {
	s := ""
	assert.True(t, valueIsEmpty(reflect.ValueOf(&s)))
}

func TestValueIsEmpty_PointerToNonEmpty_ReturnsFalse(t *testing.T) {
	s := "hello"
	assert.False(t, valueIsEmpty(reflect.ValueOf(&s)))
}

func TestValueIsEmpty_NonZeroInt_ReturnsFalse(t *testing.T) {
	assert.False(t, valueIsEmpty(reflect.ValueOf(42)))
}

func TestValueIsEmpty_ZeroInt_ReturnsFalse(t *testing.T) {
	// int falls through to default case, which returns false
	assert.False(t, valueIsEmpty(reflect.ValueOf(0)))
}

// stringerInt is an int-based type implementing fmt.Stringer.
// It hits the default branch in valueIsEmpty (not struct/slice/map/etc).
type stringerInt int

func (s stringerInt) String() string {
	if s == 0 {
		return ""
	}
	if s == -1 {
		return "   "
	}
	return "content"
}

// stringerNonEmpty is a struct-based Stringer (hits struct branch, not default).
type stringerNonEmpty struct{}

func (s stringerNonEmpty) String() string { return "content" }

func TestValueIsEmpty_StringerEmpty_ReturnsTrue(t *testing.T) {
	// stringerInt(0).String() returns "", hits default branch
	assert.True(t, valueIsEmpty(reflect.ValueOf(stringerInt(0))))
}

func TestValueIsEmpty_StringerWhitespace_ReturnsTrue(t *testing.T) {
	// stringerInt(-1).String() returns "   ", hits default branch
	assert.True(t, valueIsEmpty(reflect.ValueOf(stringerInt(-1))))
}

func TestValueIsEmpty_StringerNonEmpty_ReturnsFalse(t *testing.T) {
	// stringerInt(1).String() returns "content", hits default branch
	assert.False(t, valueIsEmpty(reflect.ValueOf(stringerInt(1))))
}

func TestValueIsEmpty_BoolFalse_ReturnsFalse(t *testing.T) {
	// bool does not implement Stringer, falls to default return false
	assert.False(t, valueIsEmpty(reflect.ValueOf(false)))
}

// --- toString tests ---

func TestToString_Nil_ReturnsFalse(t *testing.T) {
	s, ok := toString(nil)
	assert.False(t, ok)
	assert.Equal(t, "", s)
}

func TestToString_String_ReturnsValue(t *testing.T) {
	s, ok := toString("hello")
	assert.True(t, ok)
	assert.Equal(t, "hello", s)
}

func TestToString_EmptyString_ReturnsValue(t *testing.T) {
	s, ok := toString("")
	assert.True(t, ok)
	assert.Equal(t, "", s)
}

func TestToString_StringPointer_ReturnsValue(t *testing.T) {
	v := "world"
	s, ok := toString(&v)
	assert.True(t, ok)
	assert.Equal(t, "world", s)
}

func TestToString_Stringer_ReturnsStringValue(t *testing.T) {
	s, ok := toString(stringerNonEmpty{})
	assert.True(t, ok)
	assert.Equal(t, "content", s)
}

func TestToString_Int_ReturnsFalse(t *testing.T) {
	s, ok := toString(42)
	assert.False(t, ok)
	assert.Equal(t, "", s)
}

func TestToString_Float_ReturnsFalse(t *testing.T) {
	s, ok := toString(3.14)
	assert.False(t, ok)
	assert.Equal(t, "", s)
}

// Verify fmt.Stringer interface is used
type customStringer struct{ val string }

func (c customStringer) String() string { return c.val }

func TestToString_CustomStringer_ReturnsStringValue(t *testing.T) {
	s, ok := toString(customStringer{val: "custom"})
	assert.True(t, ok)
	assert.Equal(t, "custom", s)
}
