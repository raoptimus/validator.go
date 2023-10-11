package validator

import "golang.org/x/net/context"

type ContextKey uint8

const (
	ContextKeyFieldPrefix ContextKey = iota + 1
)

func FieldPrefixFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	prefix, ok := ctx.Value(ContextKeyFieldPrefix).(string)
	if !ok {
		return "", false
	}

	return prefix, true
}

func ContextWithFieldPrefix(ctx context.Context, prefix string) context.Context {
	if prevPrefix, ok := FieldPrefixFromContext(ctx); ok {
		prefix = prevPrefix + "." + prefix
	}

	return context.WithValue(ctx, ContextKeyFieldPrefix, prefix)
}
