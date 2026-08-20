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
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
