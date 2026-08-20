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
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// hangGuard runs fn and fails the test if it does not return in time, so a
// regression in the depth accounting shows up as a failure instead of a
// suite that hangs until the global timeout.
func hangGuard(t *testing.T, fn func()) {
	t.Helper()

	done := make(chan struct{})
	go func() {
		defer close(done)
		fn()
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("hashing did not terminate: depth accounting no longer covers unwrapping")
	}
}

// selfReferentialAny returns an interface value holding a pointer to itself.
// Unwrapping it alternates Interface → Pointer forever: the kind never
// settles, nothing is ever nil, so only the depth cap can stop the walk.
func selfReferentialAny() any {
	var a any
	a = &a

	return a
}

type selfReferentialStruct struct {
	Name string
	Self *selfReferentialStruct
}

func TestHashValue_SelfReferentialInterface_Terminates(t *testing.T) {
	t.Parallel()

	a := selfReferentialAny()

	hangGuard(t, func() {
		hw := newHasher()
		hashValue(&hw, reflect.ValueOf(&a))
	})
}

func TestHashValue_SelfReferentialStruct_Terminates(t *testing.T) {
	t.Parallel()

	s := &selfReferentialStruct{Name: "loop"}
	s.Self = s

	hangGuard(t, func() {
		hw := newHasher()
		hashValue(&hw, reflect.ValueOf(s))
	})
}

func TestUniqueValues_ValidateValue_SelfReferentialElements_Terminates(t *testing.T) {
	t.Parallel()

	values := []any{selfReferentialAny(), selfReferentialAny()}
	r := NewUniqueValues()

	hangGuard(t, func() {
		_ = r.ValidateValue(context.Background(), values)
	})
}

// Values buried under more pointers than maxHashDepth collapse to the same
// bucket key. Correctness must still come from the DeepEqual fallback: the
// two chains hold different payloads, so the rule has to accept them as
// unique, and two chains holding equal payloads must still be rejected.
func TestUniqueValues_ValidateValue_BeyondDepthCap_FallsBackToDeepEqual(t *testing.T) {
	t.Parallel()

	chain := func(payload string) any {
		var v any = []string{payload} // non-comparable leaf keeps the hashed path
		for range maxHashDepth + 64 {
			next := v
			v = &next
		}

		return v
	}

	ctx := context.Background()
	r := NewUniqueValues()

	require.NoError(t, r.ValidateValue(ctx, []any{chain("one"), chain("two")}))
	require.Error(t, r.ValidateValue(ctx, []any{chain("same"), chain("same")}))
}

// hotProto returns a message whose unexported state has been populated the
// way any reflective pass over it does — a logging interceptor, protojson,
// or a plain ProtoReflect call. Before that the state holds a nil
// *MessageInfo and the structural walk stops there; afterwards it opens
// into the descriptor graph.
func hotProto(s string) *wrapperspb.StringValue {
	m := wrapperspb.String(s)
	_ = m.ProtoReflect().Descriptor().FullName()

	return m
}

func TestUniqueValues_ValidateValue_HotProtoMessages_Terminates(t *testing.T) {
	t.Parallel()

	values := make([]*wrapperspb.StringValue, 0, 100)
	for i := range 100 {
		values = append(values, hotProto(fmt.Sprintf("value-%d", i)))
	}

	r := NewUniqueValues()

	hangGuard(t, func() {
		require.NoError(t, r.ValidateValue(context.Background(), values))
	})
}

func TestUniqueValues_ValidateValue_HotProtoMessages_DetectsDuplicates(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	r := NewUniqueValues()

	hangGuard(t, func() {
		require.Error(t, r.ValidateValue(ctx, []*wrapperspb.StringValue{
			hotProto("one"), hotProto("two"), hotProto("one"),
		}))
	})
}

// Equal payloads must collide regardless of whether a message has been
// through reflection: the state differs, the data does not.
func TestUniqueValues_ValidateValue_MixedHotAndColdProto_DetectsDuplicates(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	r := NewUniqueValues()

	hangGuard(t, func() {
		require.Error(t, r.ValidateValue(ctx, []*wrapperspb.StringValue{
			wrapperspb.String("same"), hotProto("same"),
		}))
		require.NoError(t, r.ValidateValue(ctx, []*wrapperspb.StringValue{
			wrapperspb.String("cold"), hotProto("hot"),
		}))
	})
}

// A message nested inside another one must not reopen the structural walk.
func TestUniqueValues_ValidateValue_NestedProtoMessages_Terminates(t *testing.T) {
	t.Parallel()

	nested := func(key, value string) *structpb.Struct {
		s, err := structpb.NewStruct(map[string]any{key: value})
		require.NoError(t, err)
		_ = s.ProtoReflect().Descriptor().FullName()

		return s
	}

	ctx := context.Background()
	r := NewUniqueValues()

	hangGuard(t, func() {
		require.NoError(t, r.ValidateValue(ctx, []*structpb.Struct{
			nested("a", "1"), nested("b", "2"),
		}))
		require.Error(t, r.ValidateValue(ctx, []*structpb.Struct{
			nested("a", "1"), nested("a", "1"),
		}))
	})
}

// Different message types can encode to identical bytes (StringValue and
// BytesValue both put their payload in field 1, wire type 2). The bucket
// key then collides, and only the equality check keeps them apart.
func TestUniqueValues_ValidateValue_ProtoWireCollision_NotADuplicate(t *testing.T) {
	t.Parallel()

	values := []proto.Message{
		wrapperspb.String("x"),
		wrapperspb.Bytes([]byte("x")),
	}

	require.NoError(t, NewUniqueValues().ValidateValue(context.Background(), values))
}
