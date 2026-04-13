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

	"github.com/raoptimus/validator.go/v2/set"
	"github.com/stretchr/testify/assert"
)

func TestNewContext_ReturnsContextWithNilDataSet(t *testing.T) {
	ctx := NewContext(context.Background())

	assert.NotNil(t, ctx)
	assert.Nil(t, ctx.ds)
}

func TestContext_Value_KeyDataSet_ReturnsDataSet(t *testing.T) {
	ds := set.NewDataSetMap(map[string]any{"name": "alice"})
	ctx := NewContext(context.Background()).withDataSet(ds)

	got := ctx.Value(KeyDataSet)

	assert.Equal(t, ds, got)
}

func TestContext_Value_KeyDataSet_NilDataSet_ReturnsNil(t *testing.T) {
	ctx := NewContext(context.Background())

	got := ctx.Value(KeyDataSet)

	assert.Nil(t, got)
}

func TestContext_Value_OtherKey_DelegatesToParent(t *testing.T) {
	type customKey string
	parentCtx := context.WithValue(context.Background(), customKey("foo"), "bar")
	ctx := NewContext(parentCtx)

	got := ctx.Value(customKey("foo"))

	assert.Equal(t, "bar", got)
}

func TestContext_Value_OtherKey_NotFound_ReturnsNil(t *testing.T) {
	type customKey string
	ctx := NewContext(context.Background())

	got := ctx.Value(customKey("missing"))

	assert.Nil(t, got)
}

func TestContext_withDataSet_ReturnsNewContextWithDataSet(t *testing.T) {
	ds := set.NewDataSetMap(map[string]any{"key": "value"})
	original := NewContext(context.Background())

	result := original.withDataSet(ds)

	assert.Equal(t, ds, result.ds)
	assert.Nil(t, original.ds, "original context must not be modified")
	assert.NotSame(t, original, result)
}

func TestContext_dataSet_WithDataSet_ReturnsDataSetAndTrue(t *testing.T) {
	ds := set.NewDataSetMap(map[string]any{"a": 1})
	ctx := NewContext(context.Background()).withDataSet(ds)

	got, ok := ctx.dataSet()

	assert.True(t, ok)
	assert.Equal(t, ds, got)
}

func TestContext_dataSet_WithoutDataSet_ReturnsNilAndFalse(t *testing.T) {
	ctx := NewContext(context.Background())

	got, ok := ctx.dataSet()

	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestDataSetFromContext_MatchingType_Successfully(t *testing.T) {
	ds := set.NewDataSetMap(map[string]any{"x": "y"})
	ctx := NewContext(context.Background()).withDataSet(ds)

	got, ok := DataSetFromContext[*set.DataSetMap](ctx)

	assert.True(t, ok)
	assert.Equal(t, ds, got)
}

func TestDataSetFromContext_NonMatchingType_Failure(t *testing.T) {
	ds := set.NewDataSetMap(map[string]any{"x": "y"})
	ctx := NewContext(context.Background()).withDataSet(ds)

	got, ok := DataSetFromContext[*set.DataSetAny](ctx)

	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestDataSetFromContext_NilDataSet_Failure(t *testing.T) {
	ctx := NewContext(context.Background())

	got, ok := DataSetFromContext[*set.DataSetMap](ctx)

	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestWithDataSet_CreatesContextWithDataSet(t *testing.T) {
	ds := set.NewDataSetMap(map[string]any{"field": "val"})

	ctx := withDataSet(context.Background(), ds)

	got := ctx.Value(KeyDataSet)
	assert.Equal(t, ds, got)
}

func TestExtractDataSet_NilContext_Failure(t *testing.T) {
	got, ok := ExtractDataSet[*set.DataSetMap](nil)

	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestExtractDataSet_NoDataSetInContext_Failure(t *testing.T) {
	ctx := context.Background()

	got, ok := ExtractDataSet[*set.DataSetMap](ctx)

	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestExtractDataSet_MatchingDataSetType_Successfully(t *testing.T) {
	ds := set.NewDataSetMap(map[string]any{"k": "v"})
	ctx := withDataSet(context.Background(), ds)

	got, ok := ExtractDataSet[*set.DataSetMap](ctx)

	assert.True(t, ok)
	assert.Equal(t, ds, got)
}

func TestExtractDataSet_MatchingDataType_Successfully(t *testing.T) {
	data := map[string]any{"hello": "world"}
	ds := set.NewDataSetMap(data)
	ctx := withDataSet(context.Background(), ds)

	got, ok := ExtractDataSet[map[string]any](ctx)

	assert.True(t, ok)
	assert.Equal(t, data, got)
}

func TestExtractDataSet_NonMatchingType_Failure(t *testing.T) {
	ds := set.NewDataSetMap(map[string]any{"k": "v"})
	ctx := withDataSet(context.Background(), ds)

	got, ok := ExtractDataSet[string](ctx)

	assert.False(t, ok)
	assert.Equal(t, "", got)
}

func TestExtractDataSet_DataSetAny_MatchingDataType_Successfully(t *testing.T) {
	inner := "some-string-data"
	ds := set.NewDataSetAny(inner)
	ctx := withDataSet(context.Background(), ds)

	got, ok := ExtractDataSet[string](ctx)

	assert.True(t, ok)
	assert.Equal(t, "some-string-data", got)
}

func TestExtractDataSet_WrongValueType_InContext_Failure(t *testing.T) {
	// KeyDataSet holds a non-DataSet value
	ctx := context.WithValue(context.Background(), KeyDataSet, "not-a-dataset")

	got, ok := ExtractDataSet[*set.DataSetMap](ctx)

	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestWithPreviousRulesErrored_SetsFlag(t *testing.T) {
	ctx := withPreviousRulesErrored(context.Background())

	got, ok := ctx.Value(KeyPreviousRulesErrored).(bool)

	assert.True(t, ok)
	assert.True(t, got)
}

func TestPreviousRulesErrored_WhenSet_ReturnsTrue(t *testing.T) {
	ctx := withPreviousRulesErrored(context.Background())

	assert.True(t, previousRulesErrored(ctx))
}

func TestPreviousRulesErrored_WhenNotSet_ReturnsFalse(t *testing.T) {
	ctx := context.Background()

	assert.False(t, previousRulesErrored(ctx))
}

func TestPreviousRulesErrored_WrongType_ReturnsFalse(t *testing.T) {
	ctx := context.WithValue(context.Background(), KeyPreviousRulesErrored, "not-a-bool")

	assert.False(t, previousRulesErrored(ctx))
}

func TestContextWithRootDataSet_SetsRootDataSet(t *testing.T) {
	data := map[string]any{"root": true}

	ctx := contextWithRootDataSet(context.Background(), data)

	got, ok := RootDataSetFromContext[map[string]any](ctx)
	assert.True(t, ok)
	assert.Equal(t, data, got)
}

func TestContextWithRootDataSet_PreviousRootMovedToPrevNested(t *testing.T) {
	firstRoot := map[string]any{"first": true}
	secondRoot := map[string]any{"second": true}

	ctx := contextWithRootDataSet(context.Background(), firstRoot)
	ctx = contextWithRootDataSet(ctx, secondRoot)

	gotRoot, ok := RootDataSetFromContext[map[string]any](ctx)
	assert.True(t, ok)
	assert.Equal(t, secondRoot, gotRoot)

	gotPrev, ok := PrevNestedDataSetFromContext[map[string]any](ctx)
	assert.True(t, ok)
	assert.Equal(t, firstRoot, gotPrev)
}

func TestContextWithRootDataSet_NoPreviousRoot_NoPrevNested(t *testing.T) {
	data := map[string]any{"only": true}

	ctx := contextWithRootDataSet(context.Background(), data)

	_, ok := PrevNestedDataSetFromContext[map[string]any](ctx)
	assert.False(t, ok)
}

func TestRootDataSetFromContext_NotSet_Failure(t *testing.T) {
	ctx := context.Background()

	got, ok := RootDataSetFromContext[map[string]any](ctx)

	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestRootDataSetFromContext_WrongType_Failure(t *testing.T) {
	ctx := contextWithRootDataSet(context.Background(), "a-string")

	got, ok := RootDataSetFromContext[int](ctx)

	assert.False(t, ok)
	assert.Equal(t, 0, got)
}

func TestContextWithNestedDataSet_SetsNestedDataSet(t *testing.T) {
	data := map[string]any{"nested": true}

	ctx := contextWithNestedDataSet(context.Background(), data)

	got, ok := NestedDataSetFromContext[map[string]any](ctx)
	assert.True(t, ok)
	assert.Equal(t, data, got)
}

func TestNestedDataSetFromContext_NotSet_Failure(t *testing.T) {
	ctx := context.Background()

	got, ok := NestedDataSetFromContext[map[string]any](ctx)

	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestNestedDataSetFromContext_WrongType_Failure(t *testing.T) {
	ctx := contextWithNestedDataSet(context.Background(), 42)

	got, ok := NestedDataSetFromContext[string](ctx)

	assert.False(t, ok)
	assert.Equal(t, "", got)
}

func TestPrevNestedDataSetFromContext_NotSet_Failure(t *testing.T) {
	ctx := context.Background()

	got, ok := PrevNestedDataSetFromContext[map[string]any](ctx)

	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestPrevNestedDataSetFromContext_WrongType_Failure(t *testing.T) {
	ctx := context.WithValue(context.Background(), KeyPrevNestedDataSet, 99)

	got, ok := PrevNestedDataSetFromContext[string](ctx)

	assert.False(t, ok)
	assert.Equal(t, "", got)
}

func TestPrevNestedDataSetFromContext_MatchingType_Successfully(t *testing.T) {
	ctx := context.WithValue(context.Background(), KeyPrevNestedDataSet, "prev-data")

	got, ok := PrevNestedDataSetFromContext[string](ctx)

	assert.True(t, ok)
	assert.Equal(t, "prev-data", got)
}

func TestContextWithRootDataSet_ThreeNested_LastPrevIsSecond(t *testing.T) {
	first := map[string]any{"level": 1}
	second := map[string]any{"level": 2}
	third := map[string]any{"level": 3}

	ctx := contextWithRootDataSet(context.Background(), first)
	ctx = contextWithRootDataSet(ctx, second)
	ctx = contextWithRootDataSet(ctx, third)

	gotRoot, ok := RootDataSetFromContext[map[string]any](ctx)
	assert.True(t, ok)
	assert.Equal(t, third, gotRoot)

	gotPrev, ok := PrevNestedDataSetFromContext[map[string]any](ctx)
	assert.True(t, ok)
	assert.Equal(t, second, gotPrev)
}
