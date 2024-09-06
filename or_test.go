package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOR_ValidateValue_Successfully(t *testing.T) {
	ctx := context.Background()

	type args struct {
		rules []Rule
		value any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ip or mac rules for ip value",
			args: args{
				rules: []Rule{
					NewIP(),
					NewMAC(),
				},
				value: "127.0.0.1",
			},
		},
		{
			name: "ip or mac rules for mac value",
			args: args{
				rules: []Rule{
					NewIP(),
					NewMAC(),
				},
				value: "00:1b:63:84:45:e6",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewOR("Value is not in ip or mac format.", tt.args.rules...)
			err := o.ValidateValue(ctx, tt.args.value)
			require.NoError(t, err)
		})
	}
}

func TestOR_ValidateValue_Failure(t *testing.T) {
	ctx := context.Background()

	type args struct {
		rules []Rule
		value any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ip or mac rules for invalid value",
			args: args{
				rules: []Rule{
					NewIP(),
					NewMAC(),
				},
				value: "hello world",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewOR("Value is not in ip or mac format.", tt.args.rules...)
			err := o.ValidateValue(ctx, tt.args.value)
			require.Error(t, err)
		})
	}
}
