/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageMeta_ValidateValue_ValidImageMeta_Successfully(t *testing.T) {
	tests := []struct {
		name  string
		value *ImageMetaData
	}{
		{
			name: "valid jpeg with zero size",
			value: &ImageMetaData{
				Name:     "photo.jpg",
				MimeType: ImageMimeTypeJPG,
				Size:     0,
				Width:    800,
				Height:   600,
			},
		},
		{
			name: "valid png with zero size",
			value: &ImageMetaData{
				Name:     "image.png",
				MimeType: ImageMimeTypePNG,
				Size:     0,
				Width:    1920,
				Height:   1080,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewImageMeta().ValidateValue(ctx, tt.value)
			assert.NoError(t, err)
		})
	}
}

func TestImageMeta_ValidateValue_NonImageMetaData_Failure(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{
			name:  "string value",
			value: "not an image",
		},
		{
			name:  "integer value",
			value: 42,
		},
		{
			name:  "nil value",
			value: nil,
		},
		{
			name:  "non-pointer ImageMetaData",
			value: ImageMetaData{Name: "test.jpg", MimeType: ImageMimeTypeJPG},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := NewImageMeta().ValidateValue(ctx, tt.value)
			assert.Error(t, err)
			assert.Equal(t, "This value must be a ImageMetaData struct.", err.Error())
		})
	}
}

func TestImageMeta_ValidateValue_InvalidMimeType_Failure(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
	}{
		{
			name:     "gif mime type",
			mimeType: ImageMimeTypeGIF,
		},
		{
			name:     "webp mime type",
			mimeType: "image/webp",
		},
		{
			name:     "empty mime type",
			mimeType: "",
		},
		{
			name:     "text mime type",
			mimeType: "text/plain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			meta := &ImageMetaData{
				Name:     "file.img",
				MimeType: tt.mimeType,
				Size:     0,
				Width:    100,
				Height:   100,
			}
			err := ValidateValue(ctx, meta, NewImageMeta())
			assert.Error(t, err)
			assert.Equal(t, "MimeType must be a equal one of [image/jpeg, image/png].", err.Error())
		})
	}
}

func TestImageMeta_ValidateValue_SizeExceedsDefault_Failure(t *testing.T) {
	tests := []struct {
		name string
		size uint64
	}{
		{
			name: "size is 1 byte exceeding default zero limit",
			size: 1,
		},
		{
			name: "large size exceeding default zero limit",
			size: 10485760,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			meta := &ImageMetaData{
				Name:     "photo.jpg",
				MimeType: ImageMimeTypeJPG,
				Size:     tt.size,
				Width:    800,
				Height:   600,
			}
			err := ValidateValue(ctx, meta, NewImageMeta())
			assert.Error(t, err)
			assert.Equal(t, "File size must be greater than 0 bytes.", err.Error())
		})
	}
}

func TestImageMeta_ValidateValue_SizeZero_Successfully(t *testing.T) {
	ctx := context.Background()
	meta := &ImageMetaData{
		Name:     "photo.jpg",
		MimeType: ImageMimeTypeJPG,
		Size:     0,
		Width:    800,
		Height:   600,
	}
	err := NewImageMeta().ValidateValue(ctx, meta)
	assert.NoError(t, err)
}

func TestImageMeta_WithMessage_CustomMessage_Failure(t *testing.T) {
	ctx := context.Background()
	err := NewImageMeta().WithMessage("custom error").ValidateValue(ctx, "not a meta")
	assert.Error(t, err)
	assert.Equal(t, "custom error.", err.Error())
}

func TestImageMeta_ValidateValue_SkipOnEmpty_NilValue_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, nil, NewImageMeta().SkipOnEmpty())
	assert.NoError(t, err)
}

func TestImageMeta_ValidateValue_WhenFuncReturnsFalse_Successfully(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, "invalid", NewImageMeta().When(func(_ context.Context, _ any) bool {
		return false
	}))
	assert.NoError(t, err)
}

func TestImageMeta_ValidateValue_WhenFuncReturnsTrue_Failure(t *testing.T) {
	ctx := context.Background()
	err := ValidateValue(ctx, "invalid", NewImageMeta().When(func(_ context.Context, _ any) bool {
		return true
	}))
	assert.Error(t, err)
}

func TestImageMeta_ValidateValue_SkipOnError_Successfully(t *testing.T) {
	r := NewImageMeta().SkipOnError()
	assert.True(t, r.shouldSkipOnError())
}

func TestImageMeta_ValidateValue_MimeTypeCheckedBeforeSize_Failure(t *testing.T) {
	ctx := context.Background()
	meta := &ImageMetaData{
		Name:     "file.gif",
		MimeType: ImageMimeTypeGIF,
		Size:     999999,
		Width:    100,
		Height:   100,
	}
	err := ValidateValue(ctx, meta, NewImageMeta())
	assert.Error(t, err)
	assert.Equal(t, "MimeType must be a equal one of [image/jpeg, image/png].", err.Error())
}
