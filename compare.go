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
}

func NewCompare(targetValue any, targetAttribute, operator string) Compare {
	c := Compare{
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

func (c Compare) ValidateValue(ctx context.Context, value any) error {
	if !c.operatorIsValid {
		return UnknownOperatorError
	}

	var (
		targetValue       any
		targetValueOrAttr any
		err               error
	)
	targetValue = c.targetValue
	targetValueOrAttr = c.targetAttribute

	if c.targetValue == nil {
		dataSet, ok := extractDataSet(ctx)
		if !ok {
			return NotExistsDataSetIntoContextError
		}
		targetValue, err = dataSet.FieldValue(c.targetAttribute)
		if err != nil {
			return err
		}
		targetValueOrAttr = targetValue
	}

	switch c.operator {
	case "==":
		if c.eq(value, targetValue) {
			return nil
		}
	case "!=":
		if !c.eq(value, targetValue) {
			return nil
		}
	case ">":
		if c.gt(value, targetValue) {
			return nil
		}
	case ">=":
		if c.eq(value, targetValue) || c.gt(value, targetValue) {
			return nil
		}
	case "<":
		if !c.eq(value, targetValue) && !c.gt(value, targetValue) {
			return nil
		}
	case "<=":
		if c.eq(value, targetValue) || !c.gt(value, targetValue) {
			return nil
		}
	}

	return NewResult().
		WithError(
			NewValidationError(c.message).
				WithParams(map[string]any{
					"targetValue":            c.targetValue,
					"targetAttribute":        c.targetAttribute,
					"targetValueOrAttribute": targetValueOrAttr,
				}),
		)
}

func (c Compare) eq(a, b any) bool {
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

	if ia, ok := a.(string); ok {
		if ib, ok := b.(string); ok {
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

func (c Compare) gt(a, b any) bool {
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

	if ia, ok := a.(string); ok {
		if ib, ok := b.(string); ok {
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
