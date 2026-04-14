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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithLanguage_And_LanguageFromContext(t *testing.T) {
	ctx := WithLanguage(context.Background(), LanguageRU)
	assert.Equal(t, LanguageRU, LanguageFromContext(ctx))
}

func TestLanguageFromContext_DefaultEN(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, LanguageEN, LanguageFromContext(ctx))
}

func TestWithLanguage_CustomLanguage(t *testing.T) {
	ctx := WithLanguage(context.Background(), Language("es"))
	assert.Equal(t, Language("es"), LanguageFromContext(ctx))
}
