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
