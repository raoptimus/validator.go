/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator

type (
	Rules   []Rule
	RuleSet map[string]Rules
)

func (rs Rules) SkipOnError() {
	for i, r := range rs {
		if rse, ok := r.(RuleSkipError); ok {
			rse.setSkipOnError(true)
		}

		rs[i] = r
	}
}
