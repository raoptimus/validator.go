package validator

import (
	"context"

	"github.com/raoptimus/validator.go/v2/set"
)

type Key uint8

const (
	KeyDataSet Key = iota + 1
)

type DataSet interface {
	FieldValue(name string) (any, error)
	FieldAliasName(name string) string
	Name() set.Name
	Data() any
}

func withDataSet(ctx context.Context, ds DataSet) context.Context {
	return context.WithValue(ctx, KeyDataSet, ds)
}

func extractDataSet(ctx context.Context) (DataSet, bool) {
	if ctx == nil {
		return nil, false
	}
	if ds, ok := ctx.Value(KeyDataSet).(DataSet); ok {
		return ds, true
	}

	return nil, false
}
