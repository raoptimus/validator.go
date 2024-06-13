package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrl_ValidateValue_ValidURL_Successfully(t *testing.T) {
	tests := []struct {
		Name         string
		ValidSchemes []string
		URL          string
	}{
		{
			Name: "valid example url with scheme https",
			URL:  "https://example.com",
		},
		{
			Name: "valid example url with scheme http",
			URL:  "http://example.com",
		},
		{
			Name: "valid example url",
			URL:  "http://example.com/hkjhiojhio/knhoijhio?a=1&s=2",
		},
		{
			Name:         "valid example url with custom scheme test.app",
			URL:          "test://example.com",
			ValidSchemes: []string{"test"},
		},
		{
			Name:         "any protocol",
			URL:          "test://example.com",
			ValidSchemes: []string{AllowAnyURLSchema},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctx := context.Background()
			r := NewURL()
			if len(tt.ValidSchemes) > 0 {
				r = r.WithValidScheme(tt.ValidSchemes...)
			}
			err := r.ValidateValue(ctx, tt.URL)
			assert.NoError(t, err)
		})
	}
}

func TestUrl_ValidateValue_InvalidURL_Failure(t *testing.T) {
	tests := []struct {
		Name         string
		ValidSchemes []string
		URL          string
	}{
		{
			Name: "invalid protocol",
			URL:  "httpz://example.com",
		},
		{
			Name: "invalid domain",
			URL:  "http://examplecom",
		},
		{
			Name: "domain is empty",
			URL:  "http://",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctx := context.Background()
			r := NewURL()
			if len(tt.ValidSchemes) > 0 {
				r.WithValidScheme(tt.ValidSchemes...)
			}
			err := r.ValidateValue(ctx, tt.URL)
			assert.Error(t, err)
		})
	}
}

func TestUrl_ValidateValue_IDN_Successfully(t *testing.T) {
	ctx := context.Background()
	err := NewURL().WithEnableIDN().ValidateValue(ctx, "https://президент.рф")
	assert.NoError(t, err)
}

func TestUrlValidateValue_EmptyValue_ReturnsError(t *testing.T) {
	ctx := context.Background()
	err := NewURL().ValidateValue(ctx, "")
	assert.Error(t, err)
}

func TestUrlValidateValue_InvalidValue_ReturnsExpectedErrorMessage(t *testing.T) {
	ctx := context.Background()
	err := NewURL().WithMessage("test error").ValidateValue(ctx, "http://")
	assert.Error(t, err)
	assert.Equal(t, "test error.", err.Error())
}
