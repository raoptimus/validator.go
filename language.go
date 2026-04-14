/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator

import "context"

type Language string

const (
	LanguageEN Language = "en"
	LanguageRU Language = "ru"
)

type ctxKeyLanguage struct{}

func WithLanguage(ctx context.Context, lang Language) context.Context {
	return context.WithValue(ctx, ctxKeyLanguage{}, lang)
}

func LanguageFromContext(ctx context.Context) Language {
	if lang, ok := ctx.Value(ctxKeyLanguage{}).(Language); ok {
		return lang
	}

	return LanguageEN
}
