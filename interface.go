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
	shouldSkipOnError() bool
	setSkipOnError(v bool)
}
