/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// --- Callback When/SkipOnEmpty/SkipOnError ---

func TestCallback_When_ReturnsFalse_Skips(t *testing.T) {
	ctx := context.Background()
	r := NewCallback(func(_ context.Context, v int) error {
		return NewResult().WithError(NewValidationError("fail"))
	}).When(func(_ context.Context, _ any) bool { return false })
	err := ValidateValue(ctx, 42, r)
	assert.NoError(t, err)
}

func TestCallback_SkipOnEmpty_NilValue(t *testing.T) {
	ctx := context.Background()
	r := NewCallback(func(_ context.Context, v int) error {
		return NewResult().WithError(NewValidationError("fail"))
	}).SkipOnEmpty()
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestCallback_SkipOnError_Flag(t *testing.T) {
	r := NewCallback(func(_ context.Context, v int) error { return nil }).SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

// --- Compare When/SkipOnEmpty/SkipOnError/WithMessage ---

func TestCompare_WithMessage_CustomMessage(t *testing.T) {
	r := NewCompare(10, "", "==").WithMessage("custom")
	ctx := context.Background()
	err := r.ValidateValue(ctx, 5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "custom")
}

func TestCompare_When_ReturnsFalse_Skips(t *testing.T) {
	ctx := context.Background()
	r := NewCompare(10, "", "==").When(func(_ context.Context, _ any) bool { return false })
	err := ValidateValue(ctx, 5, r)
	assert.NoError(t, err)
}

func TestCompare_SkipOnEmpty_NilValue(t *testing.T) {
	ctx := context.Background()
	r := NewCompare(10, "", "==").SkipOnEmpty()
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestCompare_SkipOnError_Flag(t *testing.T) {
	r := NewCompare(10, "", "==").SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

func TestCompare_InvalidOperator(t *testing.T) {
	ctx := context.Background()
	err := NewCompare(10, "", "???").ValidateValue(ctx, 5)
	assert.ErrorIs(t, err, ErrUnknownOperator)
}

func TestCompare_AllOperators(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		op    string
		a, b  any
		valid bool
	}{
		{"==", 10, 10, true},
		{"==", 10, 5, false},
		{"!=", 10, 5, true},
		{"!=", 10, 10, false},
		{">", 10, 5, true},
		{">", 5, 10, false},
		{">=", 10, 10, true},
		{">=", 10, 5, true},
		{">=", 5, 10, false},
		{"<", 5, 10, true},
		{"<", 10, 5, false},
		{"<=", 10, 10, true},
		{"<=", 5, 10, true},
		{"<=", 10, 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.op, func(t *testing.T) {
			err := NewCompare(tt.b, "", tt.op).ValidateValue(ctx, tt.a)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestCompare_GT_TypeVariants(t *testing.T) {
	ctx := context.Background()
	// uint
	assert.NoError(t, NewCompare(uint(5), "", ">").ValidateValue(ctx, uint(10)))
	assert.Error(t, NewCompare(uint(10), "", ">").ValidateValue(ctx, uint(5)))
	// int64
	assert.NoError(t, NewCompare(int64(5), "", ">").ValidateValue(ctx, int64(10)))
	assert.Error(t, NewCompare(int64(10), "", ">").ValidateValue(ctx, int64(5)))
	// float64
	assert.NoError(t, NewCompare(float64(5), "", ">").ValidateValue(ctx, float64(10)))
	assert.Error(t, NewCompare(float64(10), "", ">").ValidateValue(ctx, float64(5)))
	// string
	assert.NoError(t, NewCompare("a", "", ">").ValidateValue(ctx, "b"))
	assert.Error(t, NewCompare("b", "", ">").ValidateValue(ctx, "a"))
}

func TestCompare_TargetAttribute_FromDataSet(t *testing.T) {
	type form struct {
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	ctx := context.Background()
	f := &form{Password: "secret", ConfirmPassword: "different"}
	err := Validate(ctx, f, RuleSet{
		"ConfirmPassword": {NewCompare(nil, "Password", "==")},
	})
	assert.Error(t, err)

	f2 := &form{Password: "secret", ConfirmPassword: "secret"}
	err = Validate(ctx, f2, RuleSet{
		"ConfirmPassword": {NewCompare(nil, "Password", "==")},
	})
	assert.NoError(t, err)
}

// --- Each When/SkipOnEmpty/SkipOnError/WithMessage ---

func TestEach_WithMessage_CustomMessage(t *testing.T) {
	r := NewEach(NewRequired()).WithMessage("custom")
	assert.NotNil(t, r)
}

func TestEach_WithIncorrectInputMessage_CustomMessage(t *testing.T) {
	r := NewEach(NewRequired()).WithIncorrectInputMessage("custom")
	assert.NotNil(t, r)
}

func TestEach_When_ReturnsFalse_Skips(t *testing.T) {
	ctx := context.Background()
	r := NewEach(NewRequired()).When(func(_ context.Context, _ any) bool { return false })
	err := ValidateValue(ctx, []string{""}, r)
	assert.NoError(t, err)
}

func TestEach_SkipOnEmpty_NilValue(t *testing.T) {
	ctx := context.Background()
	r := NewEach(NewRequired()).SkipOnEmpty()
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestEach_SkipOnError_Flag(t *testing.T) {
	r := NewEach(NewRequired()).SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

// --- InRange When/SkipOnEmpty/SkipOnError/WithMessage/Not ---

func TestInRange_WithMessage_CustomMessage(t *testing.T) {
	ctx := context.Background()
	r := NewInRange([]any{1, 2, 3}).WithMessage("custom")
	err := r.ValidateValue(ctx, 5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "custom")
}

func TestInRange_Not_ReverseLogic(t *testing.T) {
	ctx := context.Background()
	r := NewInRange([]any{1, 2, 3}).Not()
	assert.NoError(t, r.ValidateValue(ctx, 5))
	assert.Error(t, r.ValidateValue(ctx, 1))
}

func TestInRange_When_ReturnsFalse_Skips(t *testing.T) {
	ctx := context.Background()
	r := NewInRange([]any{1, 2, 3}).When(func(_ context.Context, _ any) bool { return false })
	err := ValidateValue(ctx, 5, r)
	assert.NoError(t, err)
}

func TestInRange_SkipOnEmpty_NilValue(t *testing.T) {
	ctx := context.Background()
	r := NewInRange([]any{1, 2, 3}).SkipOnEmpty()
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestInRange_SkipOnError_Flag(t *testing.T) {
	r := NewInRange([]any{1, 2, 3}).SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

// --- JSON When/SkipOnEmpty/SkipOnError ---

func TestJSON_When_ReturnsFalse_Skips(t *testing.T) {
	ctx := context.Background()
	r := NewJSON().When(func(_ context.Context, _ any) bool { return false })
	err := ValidateValue(ctx, "invalid", r)
	assert.NoError(t, err)
}

func TestJSON_SkipOnEmpty_NilValue(t *testing.T) {
	ctx := context.Background()
	r := NewJSON().SkipOnEmpty()
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestJSON_SkipOnError_Flag(t *testing.T) {
	r := NewJSON().SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

// --- Nested When/SkipOnEmpty/SkipOnError/WithMessage/notNormalizeRules ---

func TestNested_WithMessage_CustomMessage(t *testing.T) {
	r := NewNested(RuleSet{}).WithMessage("custom")
	assert.NotNil(t, r)
}

func TestNested_When_ReturnsFalse_Skips(t *testing.T) {
	ctx := context.Background()
	r := NewNested(RuleSet{"Name": {NewRequired()}}).When(func(_ context.Context, _ any) bool { return false })
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestNested_SkipOnEmpty_NilValue(t *testing.T) {
	ctx := context.Background()
	r := NewNested(RuleSet{"Name": {NewRequired()}}).SkipOnEmpty()
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestNested_SkipOnError_Flag(t *testing.T) {
	r := NewNested(RuleSet{}).SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

func TestNested_notNormalizeRules(t *testing.T) {
	r := NewNested(RuleSet{}).notNormalizeRules()
	assert.NotNil(t, r)
}

// --- Number When/SkipOnEmpty/SkipOnError ---

func TestNumber_When_ReturnsFalse_Skips(t *testing.T) {
	ctx := context.Background()
	r := NewNumber(1, 10).When(func(_ context.Context, _ any) bool { return false })
	err := ValidateValue(ctx, 0, r)
	assert.NoError(t, err)
}

func TestNumber_SkipOnEmpty_NilValue(t *testing.T) {
	ctx := context.Background()
	r := NewNumber(1, 10).SkipOnEmpty()
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestNumber_SkipOnError_Flag(t *testing.T) {
	r := NewNumber(1, 10).SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

// --- OR When/SkipOnEmpty/SkipOnError/WithMessage ---

func TestOR_WithMessage_CustomMessage(t *testing.T) {
	r := NewOR("original", NewRequired()).WithMessage("custom")
	assert.NotNil(t, r)
}

func TestOR_When_ReturnsFalse_Skips(t *testing.T) {
	ctx := context.Background()
	r := NewOR("msg", NewRequired()).When(func(_ context.Context, _ any) bool { return false })
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestOR_SkipOnEmpty_NilValue(t *testing.T) {
	ctx := context.Background()
	r := NewOR("msg", NewRequired()).SkipOnEmpty()
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestOR_SkipOnError_Flag(t *testing.T) {
	r := NewOR("msg", NewRequired()).SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

// --- Required When/SkipOnError/WithAllowZeroValue ---

func TestRequired_When_ReturnsFalse_Skips(t *testing.T) {
	ctx := context.Background()
	r := NewRequired().When(func(_ context.Context, _ any) bool { return false })
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestRequired_SkipOnError_Flag(t *testing.T) {
	r := NewRequired().SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

func TestRequired_WithAllowZeroValue_Deprecated(t *testing.T) {
	r := NewRequired().WithAllowZeroValue()
	assert.NotNil(t, r)
}

// --- SQL When/SkipOnEmpty/SkipOnError/WithMessage ---

func TestSQL_WithMessage_CustomMessage(t *testing.T) {
	r := NewSQL().WithMessage("custom")
	assert.NotNil(t, r)
}

func TestSQL_When_ReturnsFalse_Skips(t *testing.T) {
	ctx := context.Background()
	r := NewSQL().When(func(_ context.Context, _ any) bool { return false })
	err := ValidateValue(ctx, "DROP TABLE;", r)
	assert.NoError(t, err)
}

func TestSQL_SkipOnEmpty_NilValue(t *testing.T) {
	ctx := context.Background()
	r := NewSQL().SkipOnEmpty()
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestSQL_SkipOnError_Flag(t *testing.T) {
	r := NewSQL().SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

// --- StringLength WithMessage/WithTooShortMessage/WithTooLongMessage/SkipOnError ---

func TestStringLength_WithMessage_CustomMessage(t *testing.T) {
	r := NewStringLength(1, 10).WithMessage("custom")
	assert.NotNil(t, r)
}

func TestStringLength_WithTooShortMessage_CustomMessage(t *testing.T) {
	r := NewStringLength(5, 10).WithTooShortMessage("too short custom")
	ctx := context.Background()
	err := r.ValidateValue(ctx, "ab")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "too short custom")
}

func TestStringLength_WithTooLongMessage_CustomMessage(t *testing.T) {
	r := NewStringLength(1, 3).WithTooLongMessage("too long custom")
	ctx := context.Background()
	err := r.ValidateValue(ctx, "abcdef")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "too long custom")
}

func TestStringLength_SkipOnError_Flag(t *testing.T) {
	r := NewStringLength(1, 10).SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

// --- Time With*/When/SkipOnEmpty/SkipOnError + ValidateValue ---

func TestTime_WithMessage_CustomMessage(t *testing.T) {
	r := NewTime().WithMessage("custom")
	ctx := context.Background()
	err := r.ValidateValue(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "custom")
}

func TestTime_WithFormatMessage_CustomMessage(t *testing.T) {
	r := NewTime().WithFormatMessage("bad format")
	ctx := context.Background()
	err := r.ValidateValue(ctx, "not-a-date")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bad format")
}

func TestTime_WithFormat_CustomFormat(t *testing.T) {
	r := NewTime().WithFormat("2006-01-02")
	ctx := context.Background()
	assert.NoError(t, r.ValidateValue(ctx, "2024-01-15"))
	assert.Error(t, r.ValidateValue(ctx, "2024-01-15T10:00:00Z"))
}

func TestTime_WithMin_TooSmall(t *testing.T) {
	ctx := context.Background()
	r := NewTime().WithMin(func(_ context.Context) (time.Time, error) {
		return time.Parse(time.RFC3339, "2024-06-01T00:00:00Z")
	})
	err := r.ValidateValue(ctx, "2024-01-01T00:00:00Z")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no less than")
}

func TestTime_WithMax_TooBig(t *testing.T) {
	ctx := context.Background()
	r := NewTime().WithMax(func(_ context.Context) (time.Time, error) {
		return time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	})
	err := r.ValidateValue(ctx, "2024-06-01T00:00:00Z")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no greater than")
}

func TestTime_WithTooSmallMessage_CustomMessage(t *testing.T) {
	r := NewTime().WithTooSmallMessage("too small custom")
	assert.NotNil(t, r)
}

func TestTime_WithTooBigMessage_CustomMessage(t *testing.T) {
	r := NewTime().WithTooBigMessage("too big custom")
	assert.NotNil(t, r)
}

func TestTime_When_ReturnsFalse_Skips(t *testing.T) {
	ctx := context.Background()
	r := NewTime().When(func(_ context.Context, _ any) bool { return false })
	err := ValidateValue(ctx, "invalid", r)
	assert.NoError(t, err)
}

func TestTime_SkipOnEmpty_NilValue(t *testing.T) {
	ctx := context.Background()
	r := NewTime().SkipOnEmpty()
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestTime_SkipOnError_Flag(t *testing.T) {
	r := NewTime().SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

func TestTime_ValidateValue_NonStringNonTime(t *testing.T) {
	ctx := context.Background()
	err := NewTime().ValidateValue(ctx, 42)
	assert.Error(t, err)
}

// --- UniqueValues When/SkipOnEmpty/SkipOnError/WithMessage ---

func TestUniqueValues_WithMessage_CustomMessage(t *testing.T) {
	r := NewUniqueValues().WithMessage("custom")
	ctx := context.Background()
	err := r.ValidateValue(ctx, []int{1, 1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "custom")
}

func TestUniqueValues_When_ReturnsFalse_Skips(t *testing.T) {
	ctx := context.Background()
	r := NewUniqueValues().When(func(_ context.Context, _ any) bool { return false })
	err := ValidateValue(ctx, []int{1, 1}, r)
	assert.NoError(t, err)
}

func TestUniqueValues_SkipOnEmpty_NilValue(t *testing.T) {
	ctx := context.Background()
	r := NewUniqueValues().SkipOnEmpty()
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestUniqueValues_SkipOnError_Flag(t *testing.T) {
	r := NewUniqueValues().SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

// --- URL When/SkipOnEmpty/SkipOnError/WithPattern ---

func TestURL_When_ReturnsFalse_Skips(t *testing.T) {
	ctx := context.Background()
	r := NewURL().When(func(_ context.Context, _ any) bool { return false })
	err := ValidateValue(ctx, "invalid", r)
	assert.NoError(t, err)
}

func TestURL_SkipOnEmpty_NilValue(t *testing.T) {
	ctx := context.Background()
	r := NewURL().SkipOnEmpty()
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestURL_SkipOnError_Flag(t *testing.T) {
	r := NewURL().SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

func TestURL_WithPattern_CustomPattern(t *testing.T) {
	r := NewURL().WithPattern("^https://")
	assert.NotNil(t, r)
}

// --- DeepLinkURL When/SkipOnEmpty/SkipOnError/WithMessage/WithInvalidSchemes ---

func TestDeepLinkURL_WithMessage_CustomMessage(t *testing.T) {
	r := NewDeepLinkURL().WithMessage("custom")
	assert.NotNil(t, r)
}

func TestDeepLinkURL_WithInvalidSchemes_CustomSchemes(t *testing.T) {
	r := NewDeepLinkURL().WithInvalidSchemes([]string{"ftp"})
	assert.NotNil(t, r)
}

func TestDeepLinkURL_When_ReturnsFalse_Skips(t *testing.T) {
	ctx := context.Background()
	r := NewDeepLinkURL().When(func(_ context.Context, _ any) bool { return false })
	err := ValidateValue(ctx, "invalid", r)
	assert.NoError(t, err)
}

func TestDeepLinkURL_SkipOnEmpty_NilValue(t *testing.T) {
	ctx := context.Background()
	r := NewDeepLinkURL().SkipOnEmpty()
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestDeepLinkURL_SkipOnError_Flag(t *testing.T) {
	r := NewDeepLinkURL().SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

// --- UUID When/SkipOnEmpty/SkipOnError/WithMessage/WithInvalidVersionMessage ---

func TestUUID_WithMessage_CustomMessage(t *testing.T) {
	r := NewUUID().WithMessage("custom")
	ctx := context.Background()
	err := r.ValidateValue(ctx, "invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "custom")
}

func TestUUID_WithInvalidVersionMessage_CustomMessage(t *testing.T) {
	r := NewUUID().WithInvalidVersionMessage("bad version")
	assert.NotNil(t, r)
}

func TestUUID_When_ReturnsFalse_Skips(t *testing.T) {
	ctx := context.Background()
	r := NewUUID().When(func(_ context.Context, _ any) bool { return false })
	err := ValidateValue(ctx, "invalid", r)
	assert.NoError(t, err)
}

func TestUUID_SkipOnEmpty_NilValue(t *testing.T) {
	ctx := context.Background()
	r := NewUUID().SkipOnEmpty()
	err := ValidateValue(ctx, nil, r)
	assert.NoError(t, err)
}

func TestUUID_SkipOnError_Flag(t *testing.T) {
	r := NewUUID().SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

// --- MatchRegularExpression WithPattern/SkipOnError ---

func TestMatchRegularExpression_WithPattern_CustomPattern(t *testing.T) {
	r := NewMatchRegularExpression(`^\d+$`).WithPattern(`^[a-z]+$`)
	ctx := context.Background()
	assert.NoError(t, r.ValidateValue(ctx, "abc"))
	assert.Error(t, r.ValidateValue(ctx, "123"))
}

func TestMatchRegularExpression_SkipOnError_Flag(t *testing.T) {
	r := NewMatchRegularExpression(`^\d+$`).SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

// --- Result: AttributeErrorMessagesIndexedByPath, CommonErrorMessages ---

func TestResult_AttributeErrorMessagesIndexedByPath(t *testing.T) {
	r := NewResult().
		WithError(
			NewValidationError("err1").WithValuePath([]string{"field1", "sub"}),
			NewValidationError("err2").WithValuePath([]string{"field1"}),
			NewValidationError("err3").WithValuePath([]string{"field2"}),
		)
	m := r.AttributeErrorMessagesIndexedByPath("field1")
	assert.Len(t, m, 2)
	assert.Equal(t, []string{"err1"}, m["sub"])
	assert.Equal(t, []string{"err2"}, m[""])
}

func TestResult_CommonErrorMessages(t *testing.T) {
	r := NewResult().
		WithError(
			NewValidationError("common1"),
			NewValidationError("common2"),
			NewValidationError("field_err").WithValuePath([]string{"field"}),
		)
	msgs := r.CommonErrorMessages()
	assert.Equal(t, []string{"common1", "common2"}, msgs)
}

func TestResult_CommonErrorMessages_Empty(t *testing.T) {
	r := NewResult().
		WithError(NewValidationError("field_err").WithValuePath([]string{"field"}))
	msgs := r.CommonErrorMessages()
	assert.Empty(t, msgs)
}

// --- set* methods (private interface implementation) ---

func TestSetWhen_AllValidators(t *testing.T) {
	wf := WhenFunc(func(_ context.Context, _ any) bool { return true })

	cb := NewCallback(func(_ context.Context, v int) error { return nil })
	cb.setWhen(wf)
	assert.NotNil(t, cb.when())

	cmp := NewCompare(1, "", "==")
	cmp.setWhen(wf)
	assert.NotNil(t, cmp.when())

	each := NewEach(NewRequired())
	each.setWhen(wf)
	assert.NotNil(t, each.when())

	im := NewImageMeta()
	im.setWhen(wf)
	assert.NotNil(t, im.when())

	ir := NewInRange([]any{1})
	ir.setWhen(wf)
	assert.NotNil(t, ir.when())

	ip := NewIP()
	ip.setWhen(wf)
	assert.NotNil(t, ip.when())

	j := NewJSON()
	j.setWhen(wf)
	assert.NotNil(t, j.when())

	mac := NewMAC()
	mac.setWhen(wf)
	assert.NotNil(t, mac.when())

	mre := NewMatchRegularExpression(`\d+`)
	mre.setWhen(wf)
	assert.NotNil(t, mre.when())

	n := NewNested(RuleSet{})
	n.setWhen(wf)
	assert.NotNil(t, n.when())

	num := NewNumeric(0, 10)
	num.setWhen(wf)
	assert.NotNil(t, num.when())

	ogrn := NewOGRN()
	ogrn.setWhen(wf)
	assert.NotNil(t, ogrn.when())

	or := NewOR("msg", NewRequired())
	or.setWhen(wf)
	assert.NotNil(t, or.when())

	req := NewRequired()
	req.setWhen(wf)
	assert.NotNil(t, req.when())

	sql := NewSQL()
	sql.setWhen(wf)
	assert.NotNil(t, sql.when())

	tm := NewTime()
	tm.setWhen(wf)
	assert.NotNil(t, tm.when())

	uv := NewUniqueValues()
	uv.setWhen(wf)
	assert.NotNil(t, uv.when())

	url := NewURL()
	url.setWhen(wf)
	assert.NotNil(t, url.when())

	dl := NewDeepLinkURL()
	dl.setWhen(wf)
	assert.NotNil(t, dl.when())

	uuid := NewUUID()
	uuid.setWhen(wf)
	assert.NotNil(t, uuid.when())
}

func TestSetSkipOnEmpty_AllValidators(t *testing.T) {
	cb := NewCallback(func(_ context.Context, v int) error { return nil })
	cb.setSkipOnEmpty(true)
	assert.True(t, cb.skipOnEmpty())

	cmp := NewCompare(1, "", "==")
	cmp.setSkipOnEmpty(true)
	assert.True(t, cmp.skipOnEmpty())

	each := NewEach(NewRequired())
	each.setSkipOnEmpty(true)
	assert.True(t, each.skipOnEmpty())

	im := NewImageMeta()
	im.setSkipOnEmpty(true)
	assert.True(t, im.skipOnEmpty())

	ir := NewInRange([]any{1})
	ir.setSkipOnEmpty(true)
	assert.True(t, ir.skipOnEmpty())

	ip := NewIP()
	ip.setSkipOnEmpty(true)
	assert.True(t, ip.skipOnEmpty())

	j := NewJSON()
	j.setSkipOnEmpty(true)
	assert.True(t, j.skipOnEmpty())

	mac := NewMAC()
	mac.setSkipOnEmpty(true)
	assert.True(t, mac.skipOnEmpty())

	mre := NewMatchRegularExpression(`\d+`)
	mre.setSkipOnEmpty(true)
	assert.True(t, mre.skipOnEmpty())

	n := NewNested(RuleSet{})
	n.setSkipOnEmpty(true)
	assert.True(t, n.skipOnEmpty())

	num := NewNumeric(0, 10)
	num.setSkipOnEmpty(true)
	assert.True(t, num.skipOnEmpty())

	ogrn := NewOGRN()
	ogrn.setSkipOnEmpty(true)
	assert.True(t, ogrn.skipOnEmpty())

	or := NewOR("msg", NewRequired())
	or.setSkipOnEmpty(true)
	assert.True(t, or.skipOnEmpty())

	sql := NewSQL()
	sql.setSkipOnEmpty(true)
	assert.True(t, sql.skipOnEmpty())

	tm := NewTime()
	tm.setSkipOnEmpty(true)
	assert.True(t, tm.skipOnEmpty())

	uv := NewUniqueValues()
	uv.setSkipOnEmpty(true)
	assert.True(t, uv.skipOnEmpty())

	url := NewURL()
	url.setSkipOnEmpty(true)
	assert.True(t, url.skipOnEmpty())

	dl := NewDeepLinkURL()
	dl.setSkipOnEmpty(true)
	assert.True(t, dl.skipOnEmpty())

	uuid := NewUUID()
	uuid.setSkipOnEmpty(true)
	assert.True(t, uuid.skipOnEmpty())
}

func TestSetSkipOnError_AllValidators(t *testing.T) {
	cb := NewCallback(func(_ context.Context, v int) error { return nil })
	cb.setSkipOnError(true)
	assert.True(t, cb.shouldSkipOnError())

	cmp := NewCompare(1, "", "==")
	cmp.setSkipOnError(true)
	assert.True(t, cmp.shouldSkipOnError())

	each := NewEach(NewRequired())
	each.setSkipOnError(true)
	assert.True(t, each.shouldSkipOnError())

	im := NewImageMeta()
	im.setSkipOnError(true)
	assert.True(t, im.shouldSkipOnError())

	ir := NewInRange([]any{1})
	ir.setSkipOnError(true)
	assert.True(t, ir.shouldSkipOnError())

	ip := NewIP()
	ip.setSkipOnError(true)
	assert.True(t, ip.shouldSkipOnError())

	j := NewJSON()
	j.setSkipOnError(true)
	assert.True(t, j.shouldSkipOnError())

	mac := NewMAC()
	mac.setSkipOnError(true)
	assert.True(t, mac.shouldSkipOnError())

	mre := NewMatchRegularExpression(`\d+`)
	mre.setSkipOnError(true)
	assert.True(t, mre.shouldSkipOnError())

	n := NewNested(RuleSet{})
	n.setSkipOnError(true)
	assert.True(t, n.shouldSkipOnError())

	num := NewNumeric(0, 10)
	num.setSkipOnError(true)
	assert.True(t, num.shouldSkipOnError())

	ogrn := NewOGRN()
	ogrn.setSkipOnError(true)
	assert.True(t, ogrn.shouldSkipOnError())

	or := NewOR("msg", NewRequired())
	or.setSkipOnError(true)
	assert.True(t, or.shouldSkipOnError())

	sql := NewSQL()
	sql.setSkipOnError(true)
	assert.True(t, sql.shouldSkipOnError())

	tm := NewTime()
	tm.setSkipOnError(true)
	assert.True(t, tm.shouldSkipOnError())

	uv := NewUniqueValues()
	uv.setSkipOnError(true)
	assert.True(t, uv.shouldSkipOnError())

	url := NewURL()
	url.setSkipOnError(true)
	assert.True(t, url.shouldSkipOnError())

	dl := NewDeepLinkURL()
	dl.setSkipOnError(true)
	assert.True(t, dl.shouldSkipOnError())

	uuid := NewUUID()
	uuid.setSkipOnError(true)
	assert.True(t, uuid.shouldSkipOnError())
}
