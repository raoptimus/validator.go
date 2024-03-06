package validator

import (
	"context"

	"github.com/gofrs/uuid"
)

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
	whenFunc              WhenFunc
	skipEmpty             bool
}

func NewUUID() *UUID {
	return &UUID{
		message:               "Invalid UUID format.",
		invalidVersionMessage: "UUID version must be equal to {version}.",
	}
}

func (r *UUID) WithMessage(message string) *UUID {
	rc := *r
	rc.message = message

	return &rc
}

func (r *UUID) WithInvalidVersionMessage(message string) *UUID {
	rc := *r
	rc.invalidVersionMessage = message

	return &rc
}

func (r *UUID) WithVersion(version UUIDVersion) *UUID {
	rc := *r
	rc.version = version

	return &rc
}

func (r *UUID) When(v WhenFunc) *UUID {
	rc := *r
	rc.whenFunc = v

	return &rc
}

func (r *UUID) when() WhenFunc {
	return r.whenFunc
}

func (r *UUID) setWhen(v WhenFunc) {
	r.whenFunc = v
}

func (r *UUID) SkipOnEmpty(v bool) *UUID {
	rc := *r
	rc.skipEmpty = v

	return &rc
}

func (r *UUID) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *UUID) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *UUID) ValidateValue(_ context.Context, value any) error {
	v, ok := toString(value)
	if !ok {
		return NewResult().WithError(NewValidationError(r.message))
	}

	parsedUUID, err := uuid.FromString(v)
	if err != nil {
		return NewResult().WithError(NewValidationError(r.message))
	}

	if parsedUUID.IsNil() {
		return NewResult().WithError(NewValidationError(r.message))
	}

	if r.version > 0 && byte(r.version) != parsedUUID.Version() {
		return NewResult().
			WithError(
				NewValidationError(r.message).
					WithParams(map[string]any{
						"version": r.version,
					}),
			)
	}

	return nil
}
