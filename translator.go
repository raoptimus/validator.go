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
	"fmt"
	"maps"
	"strings"
	"sync"
)

type Translator interface {
	Translate(
		ctx context.Context,
		id string,
		params map[string]any,
	) string
}

type DummyTranslator struct {
}

func (d *DummyTranslator) Translate(_ context.Context, id string, params map[string]any) string {
	return replacePlaceholders(id, params)
}

type CatalogTranslator struct {
	catalog *TranslationCatalog
}

func NewCatalogTranslator(catalog *TranslationCatalog) *CatalogTranslator {
	return &CatalogTranslator{catalog: catalog}
}

func (t *CatalogTranslator) Translate(ctx context.Context, id string, params map[string]any) string {
	lang := LanguageFromContext(ctx)

	if msg, ok := t.catalog.Get(lang, id); ok {
		return replacePlaceholders(msg, params)
	}

	return replacePlaceholders(id, params)
}

type TranslationCatalog struct {
	mu           sync.RWMutex
	translations map[Language]map[string]string
}

func NewTranslationCatalog() *TranslationCatalog {
	return &TranslationCatalog{
		translations: make(map[Language]map[string]string),
	}
}

func (c *TranslationCatalog) Register(lang Language, messages map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	existing, ok := c.translations[lang]
	if !ok {
		existing = make(map[string]string, len(messages))
		c.translations[lang] = existing
	}

	maps.Copy(existing, messages)
}

func (c *TranslationCatalog) Get(lang Language, id string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	msgs, ok := c.translations[lang]
	if !ok {
		return "", false
	}

	msg, ok := msgs[id]

	return msg, ok
}

func (c *TranslationCatalog) Missing(lang Language) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	enMsgs := c.translations[LanguageEN]
	langMsgs := c.translations[lang]

	var missing []string

	for id := range enMsgs {
		if _, ok := langMsgs[id]; !ok {
			missing = append(missing, id)
		}
	}

	return missing
}

var Translations = NewTranslationCatalog()

var DefaultTranslator Translator = NewCatalogTranslator(Translations)

func SetTranslator(t Translator) {
	DefaultTranslator = t
}

func replacePlaceholders(id string, params map[string]any) string {
	for name, value := range params {
		attrPlaceholder := "{" + name + "}"
		if !strings.Contains(id, attrPlaceholder) {
			continue
		}

		id = strings.ReplaceAll(id, attrPlaceholder, fmt.Sprintf("%v", value))
	}

	return id
}
