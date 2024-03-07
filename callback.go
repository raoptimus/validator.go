package validator

import (
	"context"

	"github.com/pkg/errors"
)

type CallbackFunc[T any] func(ctx context.Context, value T) error
type Callback[T any] struct {
	f         CallbackFunc[T]
	whenFunc  WhenFunc
	skipEmpty bool
	skipError bool
}

func NewCallback[T any](f CallbackFunc[T]) *Callback[T] {
	return &Callback[T]{
		f: f,
	}
}

func (r *Callback[T]) When(v WhenFunc) *Callback[T] {
	rc := *r
	rc.whenFunc = v

	return &rc
}

func (r *Callback[T]) when() WhenFunc {
	return r.whenFunc
}

func (r *Callback[T]) setWhen(v WhenFunc) {
	r.whenFunc = v
}

func (r *Callback[T]) SkipOnEmpty() *Callback[T] {
	rc := *r
	rc.skipEmpty = true

	return &rc
}

func (r *Callback[T]) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *Callback[T]) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *Callback[T]) SkipOnError() *Callback[T] {
	rs := *r
	rs.skipError = true

	return &rs
}

func (r *Callback[T]) shouldSkipOnError() bool {
	return r.skipError
}
func (r *Callback[T]) setSkipOnError(v bool) {
	r.skipError = v
}

func (c *Callback[T]) ValidateValue(ctx context.Context, value any) error {
	v, ok := value.(T)
	if !ok {
		var v T
		return errors.Wrapf(CallbackUnexpectedValueTypeError, "got %T want %T", value, v)
	}

	if err := c.f(ctx, v); err != nil {
		var vErr *ValidationError
		if errors.As(err, &vErr) {
			return NewResult().WithError(vErr)
		}

		var result *Result
		if errors.As(err, &result) {
			return result
		}

		return err
	}

	return nil
}
