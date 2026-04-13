/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator

const humanRegexp = "^[\\p{L}\\d !?-~-–—:;#()‘.,'\"«»„“’`´′″\\[\\]\\/]+$"

type HumanText struct {
	*MatchRegularExpression
}

func NewHumanText() *HumanText {
	return &HumanText{
		MatchRegularExpression: NewMatchRegularExpression(humanRegexp).
			WithMessage("This value must be a normal text."),
	}
}
