package validator

import (
	"context"
	"slices"
	"strings"
)

const (
	ImageMimeTypePNG = "image/png"
	ImageMimeTypeJPG = "image/jpeg"
	ImageMimeTypeGIF = "image/gif"
)

type ImageMetaData struct {
	Name     string
	MimeType string
	Size     uint64
	Width    int
	Height   int
}

type ImageMeta struct {
	message                string
	invalidMimeTypeMessage string
	tooLongSizeMessage     string
	mimeTypes              []string
	maxFileSizeBytes       uint64
	whenFunc               WhenFunc
	skipEmpty              bool
	skipError              bool
}

func NewImageMeta() *ImageMeta {
	return &ImageMeta{
		message:                "This value must be a ImageMetaData struct.",
		invalidMimeTypeMessage: "MimeType must be a equal one of [{mimeTypes}].",
		tooLongSizeMessage:     "File size must be greater than {maxFileSizeBytes} bytes.",
		mimeTypes: []string{
			ImageMimeTypeJPG,
			ImageMimeTypePNG,
		},
		maxFileSizeBytes: 0, // no limit
	}
}

func (r *ImageMeta) WithMessage(message string) *ImageMeta {
	rc := *r
	rc.message = message

	return &rc
}

func (r *ImageMeta) When(v WhenFunc) *ImageMeta {
	rc := *r
	rc.whenFunc = v

	return &rc
}

func (r *ImageMeta) when() WhenFunc {
	return r.whenFunc
}

func (r *ImageMeta) setWhen(v WhenFunc) {
	r.whenFunc = v
}

func (r *ImageMeta) SkipOnEmpty() *ImageMeta {
	rc := *r
	rc.skipEmpty = true

	return &rc
}

func (r *ImageMeta) skipOnEmpty() bool {
	return r.skipEmpty
}

func (r *ImageMeta) setSkipOnEmpty(v bool) {
	r.skipEmpty = v
}

func (r *ImageMeta) SkipOnError() *ImageMeta {
	rs := *r
	rs.skipError = true

	return &rs
}

func (r *ImageMeta) shouldSkipOnError() bool {
	return r.skipError
}
func (r *ImageMeta) setSkipOnError(v bool) {
	r.skipError = v
}

func (r *ImageMeta) ValidateValue(_ context.Context, value any) error {
	meta, ok := value.(*ImageMetaData)
	if !ok {
		return NewResult().WithError(NewValidationError(r.message))
	}

	if !slices.Contains(r.mimeTypes, meta.MimeType) {
		return NewResult().
			WithError(
				NewValidationError(r.invalidMimeTypeMessage).
					WithParams(map[string]any{
						"mimeTypes": strings.Join(r.mimeTypes, ", "),
					}),
			)
	}

	if meta.Size > r.maxFileSizeBytes {
		return NewResult().
			WithError(
				NewValidationError(r.tooLongSizeMessage).
					WithParams(map[string]any{
						"maxFileSizeBytes": r.maxFileSizeBytes,
					}),
			)
	}

	return nil
}
