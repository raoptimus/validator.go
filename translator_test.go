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

func TestDummyTranslator_Translate_NoParams(t *testing.T) {
	ctx := context.Background()
	tr := &DummyTranslator{}
	result := tr.Translate(ctx, "Value cannot be blank.", nil)
	assert.Equal(t, "Value cannot be blank.", result)
}

func TestDummyTranslator_Translate_WithSingleParam(t *testing.T) {
	ctx := context.Background()
	tr := &DummyTranslator{}
	params := map[string]any{"min": 3}
	result := tr.Translate(ctx, "Value must be at least {min} characters.", params)
	assert.Equal(t, "Value must be at least 3 characters.", result)
}

func TestDummyTranslator_Translate_WithMultipleParams(t *testing.T) {
	ctx := context.Background()
	tr := &DummyTranslator{}
	params := map[string]any{"min": 1, "max": 100}
	result := tr.Translate(ctx, "Value must be between {min} and {max}.", params)
	assert.Equal(t, "Value must be between 1 and 100.", result)
}

func TestDummyTranslator_Translate_ParamNotInTemplate(t *testing.T) {
	ctx := context.Background()
	tr := &DummyTranslator{}
	params := map[string]any{"unused": "data"}
	result := tr.Translate(ctx, "Value is invalid.", params)
	assert.Equal(t, "Value is invalid.", result)
}

func TestDummyTranslator_Translate_EmptyString(t *testing.T) {
	ctx := context.Background()
	tr := &DummyTranslator{}
	result := tr.Translate(ctx, "", nil)
	assert.Equal(t, "", result)
}

func TestDummyTranslator_Translate_EmptyParams(t *testing.T) {
	ctx := context.Background()
	tr := &DummyTranslator{}
	result := tr.Translate(ctx, "Hello {name}.", map[string]any{})
	assert.Equal(t, "Hello {name}.", result)
}

func TestDummyTranslator_Translate_RepeatedPlaceholder(t *testing.T) {
	ctx := context.Background()
	tr := &DummyTranslator{}
	params := map[string]any{"val": "X"}
	result := tr.Translate(ctx, "{val} and {val} again.", params)
	assert.Equal(t, "X and X again.", result)
}

func TestDummyTranslator_Translate_StringParam(t *testing.T) {
	ctx := context.Background()
	tr := &DummyTranslator{}
	params := map[string]any{"field": "email"}
	result := tr.Translate(ctx, "The {field} is required.", params)
	assert.Equal(t, "The email is required.", result)
}

func TestSetTranslator_ChangesDefaultTranslator(t *testing.T) {
	original := DefaultTranslator
	defer SetTranslator(original)

	custom := &DummyTranslator{}
	SetTranslator(custom)
	assert.Same(t, custom, DefaultTranslator)
}

func TestSetTranslator_DefaultTranslatorIsCatalog(t *testing.T) {
	assert.IsType(t, &CatalogTranslator{}, DefaultTranslator)
}

func TestCatalogTranslator_Translate_EN(t *testing.T) {
	ctx := WithLanguage(context.Background(), LanguageEN)
	result := DefaultTranslator.Translate(ctx, MessageRequired, nil)
	assert.Equal(t, "Value cannot be blank.", result)
}

func TestCatalogTranslator_Translate_RU(t *testing.T) {
	ctx := WithLanguage(context.Background(), LanguageRU)
	result := DefaultTranslator.Translate(ctx, MessageRequired, nil)
	assert.Equal(t, "Значение не должно быть пустым.", result)
}

func TestCatalogTranslator_Translate_RU_WithParams(t *testing.T) {
	ctx := WithLanguage(context.Background(), LanguageRU)
	params := map[string]any{"min": 3}
	result := DefaultTranslator.Translate(ctx, MessageTooShort, params)
	assert.Equal(t, "Значение должно содержать минимум 3 символов.", result)
}

func TestCatalogTranslator_Translate_FallbackToMessageID(t *testing.T) {
	ctx := WithLanguage(context.Background(), Language("fr"))
	result := DefaultTranslator.Translate(ctx, MessageRequired, nil)
	assert.Equal(t, "Value cannot be blank.", result)
}

func TestCatalogTranslator_Translate_NoLanguageInContext(t *testing.T) {
	ctx := context.Background()
	result := DefaultTranslator.Translate(ctx, MessageRequired, nil)
	assert.Equal(t, "Value cannot be blank.", result)
}

func TestTranslationCatalog_Register_CustomLanguage(t *testing.T) {
	catalog := NewTranslationCatalog()
	catalog.Register(Language("es"), map[string]string{
		MessageRequired: "El valor no puede estar vacío.",
	})

	tr := NewCatalogTranslator(catalog)
	ctx := WithLanguage(context.Background(), Language("es"))
	result := tr.Translate(ctx, MessageRequired, nil)
	assert.Equal(t, "El valor no puede estar vacío.", result)
}

func TestTranslationCatalog_Register_OverrideExisting(t *testing.T) {
	catalog := NewTranslationCatalog()
	catalog.Register(LanguageEN, map[string]string{
		MessageRequired: "Field is required.",
	})

	tr := NewCatalogTranslator(catalog)
	ctx := WithLanguage(context.Background(), LanguageEN)
	result := tr.Translate(ctx, MessageRequired, nil)
	assert.Equal(t, "Field is required.", result)
}

func TestTranslationCatalog_Missing(t *testing.T) {
	catalog := NewTranslationCatalog()
	catalog.Register(LanguageEN, map[string]string{
		"msg1": "Message 1",
		"msg2": "Message 2",
		"msg3": "Message 3",
	})
	catalog.Register(Language("es"), map[string]string{
		"msg1": "Mensaje 1",
	})

	missing := catalog.Missing(Language("es"))
	assert.Len(t, missing, 2)
	assert.Contains(t, missing, "msg2")
	assert.Contains(t, missing, "msg3")
}

func TestTranslationCatalog_Missing_FullCoverage(t *testing.T) {
	catalog := NewTranslationCatalog()
	catalog.Register(LanguageEN, map[string]string{
		"msg1": "Message 1",
	})
	catalog.Register(Language("de"), map[string]string{
		"msg1": "Nachricht 1",
	})

	missing := catalog.Missing(Language("de"))
	assert.Empty(t, missing)
}

func TestAllMessagesTranslated_RU(t *testing.T) {
	missing := Translations.Missing(LanguageRU)
	assert.Empty(t, missing, "missing RU translations: %v", missing)
}
