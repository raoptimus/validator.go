package validator

import (
	"context"
	"time"
)

type Compare struct {
	targetValue     any
	targetAttribute string
	operator        string
	message         string
	operatorIsValid bool
	whenFunc        WhenFunc
	skipEmpty       bool
}

func NewCompare(targetValue any, targetAttribute, operator string) *Compare {
	c := &Compare{
		targetValue:     targetValue,
		targetAttribute: targetAttribute,
		operator:        operator,
		operatorIsValid: true,
	}

	switch operator {
	case "==":
		c.message = "Value must be equal to '{targetValueOrAttribute}'."
	case "!=":
		c.message = "Value must not be equal to '{targetValueOrAttribute}'."
	case ">":
		c.message = "Value must be greater than '{targetValueOrAttribute}'"
	case ">=":
		c.message = "Value must be greater than or equal to '{targetValueOrAttribute}'"
	case "<":
		c.message = "Value must be less than '{targetValueOrAttribute}'"
	case "<=":
		c.message = "Value must be less than or equal to '{targetValueOrAttribute}'"
	default:
		c.operatorIsValid = false
	}

	return c
}

func (r *Compare) WithMessage(v string) *Compare {
	rc := *r
	rc.message = v

	return &rc
}

func (r *Compare) When(v WhenFunc) *Compare {
	rc := *r
	rc.whenFunc = v

	return &rc
}

func (r *Compare) when() WhenFunc {
	return r.whenFunc
}

func (r *Compare) setWhen(v WhenFunc) {
	r.whenFunc = v
}

func (r *Compare) SkipOnEmpty(v bool) *Compare {
	rc := *r
	rc.skipEmpty = v

	return &rc
}

func (r *Compare) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *Compare) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *Compare) ValidateValue(ctx context.Context, value any) error {
	if !r.operatorIsValid {
		return UnknownOperatorError
	}

	var (
		targetValue       any
		targetValueOrAttr any
		err               error
	)
	targetValue = r.targetValue
	targetValueOrAttr = r.targetAttribute

	if r.targetValue == nil {
		dataSet, ok := ExtractDataSet[DataSet](ctx)
		if !ok {
			return NotExistsDataSetIntoContextError
		}
		targetValue, err = dataSet.FieldValue(r.targetAttribute)
		if err != nil {
			return err
		}
		targetValueOrAttr = targetValue
	}

	switch r.operator {
	case "==":
		if r.eq(value, targetValue) {
			return nil
		}
	case "!=":
		if !r.eq(value, targetValue) {
			return nil
		}
	case ">":
		if r.gt(value, targetValue) {
			return nil
		}
	case ">=":
		if r.eq(value, targetValue) || r.gt(value, targetValue) {
			return nil
		}
	case "<":
		if !r.eq(value, targetValue) && !r.gt(value, targetValue) {
			return nil
		}
	case "<=":
		if r.eq(value, targetValue) || !r.gt(value, targetValue) {
			return nil
		}
	}

	return NewResult().
		WithError(
			NewValidationError(r.message).
				WithParams(map[string]any{
					"targetValue":            r.targetValue,
					"targetAttribute":        r.targetAttribute,
					"targetValueOrAttribute": targetValueOrAttr,
				}),
		)
}

func (r *Compare) eq(a, b any) bool {
	if ia, ok := a.(int); ok {
		if ib, ok := b.(int); ok {
			return ia == ib
		}
	}

	if ia, ok := a.(uint); ok {
		if ib, ok := b.(uint); ok {
			return ia == ib
		}
	}

	if ia, ok := a.(int64); ok {
		if ib, ok := b.(int64); ok {
			return ia == ib
		}
	}

	if ia, ok := a.(float64); ok {
		if ib, ok := b.(float64); ok {
			return ia == ib
		}
	}

	if ia, ok := toString(a); ok {
		if ib, ok := toString(b); ok {
			return ia == ib
		}
	}

	if ia, ok := a.(time.Time); ok {
		if ib, ok := b.(time.Time); ok {
			return ia.Equal(ib)
		}
	}

	return a == b
}

func (r *Compare) gt(a, b any) bool {
	if ia, ok := a.(int); ok {
		if ib, ok := b.(int); ok {
			return ia > ib
		}
	}

	if ia, ok := a.(uint); ok {
		if ib, ok := b.(uint); ok {
			return ia > ib
		}
	}

	if ia, ok := a.(int64); ok {
		if ib, ok := b.(int64); ok {
			return ia > ib
		}
	}

	if ia, ok := a.(float64); ok {
		if ib, ok := b.(float64); ok {
			return ia > ib
		}
	}

	if ia, ok := toString(a); ok {
		if ib, ok := toString(b); ok {
			return ia > ib
		}
	}

	if ia, ok := a.(time.Time); ok {
		if ib, ok := b.(time.Time); ok {
			return ia.After(ib)
		}
	}

	return false
}
