package validator

const humanRegexp = `^[А-Яа-яЁёa-zA-Z0-9 ,.-]+$`

type HumanText struct {
	*MatchRegularExpression
}

func NewHumanText() *HumanText {
	return &HumanText{
		MatchRegularExpression: NewMatchRegularExpression(humanRegexp).
			WithMessage("This value must be a normal text."),
	}
}
