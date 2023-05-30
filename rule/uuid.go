package rule

const (
	uuidRegexp = `^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}$`
	zeroUUID   = "00000000-0000-0000-0000-000000000000"
)

type UUID struct {
	basicRule MatchRegularExpression
}

func NewUUID() UUID {
	return UUID{
		basicRule: NewMatchRegularExpression(uuidRegexp).
			WithMessage("Invalid UUID format."),
	}
}

func (s UUID) WithMessage(message string) UUID {
	s.basicRule = s.basicRule.WithMessage(message)
	return s
}

func (s UUID) ValidateValue(value any) error {
	v, ok := toString(value)
	if !ok {
		return NewResult().WithError(formatMessage(s.basicRule.message))
	}

	if err := s.basicRule.ValidateValue(v); err != nil {
		return err
	}

	if v == zeroUUID {
		return NewResult().WithError(formatMessage(s.basicRule.message))
	}

	return nil
}
