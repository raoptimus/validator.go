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
