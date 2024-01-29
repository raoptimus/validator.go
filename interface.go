package validator

import "context"

type Rule interface {
	ValidateValue(ctx context.Context, value any) error
}
