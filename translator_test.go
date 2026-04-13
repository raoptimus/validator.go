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

func TestSetTranslator_DefaultTranslatorIsDummy(t *testing.T) {
	assert.IsType(t, &DummyTranslator{}, DefaultTranslator)
}
