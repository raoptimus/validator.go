package validator

import (
	"context"
)

type WhenFunc func(ctx context.Context, value any) bool

type Rule interface {
	ValidateValue(ctx context.Context, value any) error
}

type RuleWhen interface {
	when() WhenFunc
	setWhen(v WhenFunc)
}

type RuleSkipEmpty interface {
	skipOnEmpty() bool
	setSkipOnEmpty(v bool)
}

type RuleSkipError interface {
	shouldSkipOnError(ctx context.Context) bool
	setSkipOnError(v bool)
}
