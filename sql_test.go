package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQL_ValidateValue_Successfully(t *testing.T) {
	// region data provider
	ctx := context.Background()

	type args struct {
		value any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "simple query",
			args: args{
				value: "select * from news",
			},
		},
		{
			name: "with where",
			args: args{
				value: "select * from news where (news.news_view_title = 'Новость 🧚‍♀️')",
			},
		},
	}
	// endregion

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSQL()
			require.NoError(t, r.ValidateValue(ctx, tt.args.value))
		})
	}
}

func TestSQL_ValidateValue_Failure(t *testing.T) {
	// region data provider
	ctx := context.Background()

	type args struct {
		value any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "simple query without columns",
			args: args{
				value: "select from news",
			},
		},
		{
			name: "simple query without table declaration",
			args: args{
				value: "select * from",
			},
		},
		{
			name: "simple query without select",
			args: args{
				value: "* from",
			},
		},
	}
	// endregion

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSQL()
			require.Error(t, r.ValidateValue(ctx, tt.args.value))
		})
	}
}

func TestSQLAsWhereClause_ValidateValue_Successfully(t *testing.T) {
	// region data provider
	ctx := context.Background()

	type args struct {
		value any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "single condition",
			args: args{
				value: "(news.news_view_title = 'Новость 🧚‍♀️')",
			},
		},
		{
			name: "two conditions",
			args: args{
				value: "(news.news_view_title = 'Новость 🧚‍♀️') OR (news.news_view_title = 'Новость')",
			},
		},
	}
	// endregion

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSQL().AsWhereClause()
			require.NoError(t, r.ValidateValue(ctx, tt.args.value))
		})
	}
}

func TestSQLAsWhereClause_ValidateValue_Failure(t *testing.T) {
	// region data provider
	ctx := context.Background()

	type args struct {
		value any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "no conditions",
			args: args{
				value: "",
			},
		},
		{
			name: "does not close bracket",
			args: args{
				value: "(news.news_view_title = 'Новость 🧚‍♀️'",
			},
		},
		{
			name: "does not close bracket in second condition",
			args: args{
				value: "(news.news_view_title = 'Новость 🧚‍♀️') OR (news.news_view_title = 'Новость'",
			},
		},
		{
			name: "does not have operand",
			args: args{
				value: "(news.news_view_title = 'Новость 🧚‍♀️') (news.news_view_title = 'Новость')",
			},
		},
	}
	// endregion

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSQL().AsWhereClause()
			require.Error(t, r.ValidateValue(ctx, tt.args.value))
		})
	}
}
