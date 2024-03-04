package rule

import "github.com/gofrs/uuid"

type UUIDVersion byte

const (
	UUIDVersionV1 UUIDVersion = 1
	UUIDVersionV3 UUIDVersion = 3
	UUIDVersionV4 UUIDVersion = 4
	UUIDVersionV5 UUIDVersion = 5
	UUIDVersionV6 UUIDVersion = 6
	UUIDVersionV7 UUIDVersion = 7
)

type UUID struct {
	message               string
	invalidVersionMessage string
	version               UUIDVersion
	skipOnEmpty           bool
}

func NewUUID() UUID {
	return UUID{
		message:               "Invalid UUID format.",
		invalidVersionMessage: "UUID version must be equal to {version}.",
	}
}

func (s UUID) SkipOnEmpty() UUID {
	s.skipOnEmpty = true
	return s
}

func (s UUID) WithMessage(message string) UUID {
	s.message = message
	return s
}

func (s UUID) WithInvalidVersionMessage(message string) UUID {
	s.invalidVersionMessage = message
	return s
}

func (s UUID) WithVersion(version UUIDVersion) UUID {
	s.version = version
	return s
}

func (s UUID) ValidateValue(value any) error {
	v, ok := toString(value)
	if !ok {
		return NewResult().WithError(formatMessage(s.message))
	}

	if v == "" && s.skipOnEmpty {
		return nil
	}

	parsedUUID, err := uuid.FromString(v)
	if err != nil {
		return NewResult().
			WithError(formatMessage(s.message))
	}

	if parsedUUID.IsNil() {
		return NewResult().
			WithError(formatMessage(s.message))
	}

	if s.version > 0 && byte(s.version) != parsedUUID.Version() {
		return NewResult().
			WithError(formatMessageWithArgs(s.invalidVersionMessage, map[string]any{"version": s.version}))
	}

	return nil
}
