package validator

import (
	"context"
	"fmt"
	"strings"
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
	for name, value := range params {
		attrPlaceholder := "{" + name + "}"
		if !strings.Contains(id, attrPlaceholder) {
			continue
		}

		id = strings.ReplaceAll(id, attrPlaceholder, fmt.Sprintf("%v", value))
	}

	return id
}

var DefaultTranslator Translator = &DummyTranslator{}

func SetTranslator(t Translator) {
	DefaultTranslator = t
}
