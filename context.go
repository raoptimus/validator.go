package validator

import (
	"context"

	"github.com/raoptimus/validator.go/v2/set"
)

type Key uint8

const (
	KeyDataSet Key = iota + 1
	PreviousRulesErrored
)

type DataSet interface {
	FieldValue(name string) (any, error)
	FieldAliasName(name string) string
	Name() set.Name
	Data() any
}

type Context struct {
	context.Context
	ds DataSet
}

func NewContext(ctx context.Context) *Context {
	return &Context{Context: ctx}
}

func (c *Context) Value(key any) any {
	if key == KeyDataSet {
		return c.ds
	}

	return c.Context.Value(key)
}

func (c *Context) withDataSet(ds DataSet) *Context {
	cc := *c
	cc.ds = ds

	return &cc
}

func (c *Context) dataSet() (DataSet, bool) {
	return c.ds, c.ds != nil
}

func DataSetFromContext[T DataSet](ctx *Context) (T, bool) {
	if ds, ok := ctx.dataSet(); ok {
		if dsT, ok2 := ds.(T); ok2 {
			return dsT, true
		}
	}
	var v T

	return v, false
}

// todo: write funcs if context.Context interface

func withDataSet(ctx context.Context, ds DataSet) context.Context {
	return NewContext(ctx).withDataSet(ds)
	//return context.WithValue(ctx, KeyDataSet, ds)
}

func ExtractDataSet[T DataSet](ctx context.Context) (T, bool) {
	var v T
	if ctx == nil {
		return v, false
	}

	ds, ok := ctx.Value(KeyDataSet).(T)
	if !ok {
		return v, false
	}

	return ds, true
}

//func withPreviousRulesErrored(ctx context.Context) context.Context {
//	return context.WithValue(ctx, PreviousRulesErrored, true)
//}

func previousRulesErrored(ctx context.Context) bool {
	if y, ok := ctx.Value(PreviousRulesErrored).(bool); ok {
		return y
	}
	return false
}
