package validator

import (
	"math"
	"reflect"

	"google.golang.org/protobuf/proto"
)

// Inlined FNV-64a constants. Using the raw state instead of hash.Hash64
// avoids the interface allocation and per-field Reset cost of fnv.New64a.
const (
	fnvOffset64 = 14695981039346656037
	fnvPrime64  = 1099511628211
	lowByteMask = 0xff
)

// Type tags disambiguate values across kinds so that, e.g., the uint 65
// does not collide with the string "A". Written before each value.
const (
	tagNil byte = iota
	tagBool
	tagInt
	tagUint
	tagFloat
	tagComplex
	tagString
	tagStruct
	tagStructEnd
	tagSlice
	tagMap
	tagProto
)

// hasher is a zero-alloc FNV-64a streaming hasher. Its output is only
// used as a bucket key — equality is verified with reflect.DeepEqual on
// hash collisions in validateHashKey, so collisions are correctness-safe.
type hasher struct {
	state uint64
}

func newHasher() hasher {
	return hasher{state: fnvOffset64}
}

func (hw *hasher) reset() {
	hw.state = fnvOffset64
}

func (hw *hasher) writeByte(b byte) {
	hw.state = (hw.state ^ uint64(b)) * fnvPrime64
}

// writeUint64 absorbs the 8 little-endian bytes of v into the FNV-64a state.
// Manually unrolled: keeping state in a local register avoids reloading
// hw.state through the pointer on every byte, and elides the loop overhead
// (which is comparable to the per-byte work itself).
//
//nolint:mnd // bit-shift constants are inherent to a uint64 byte-by-byte unroll
func (hw *hasher) writeUint64(v uint64) {
	s := hw.state
	s = (s ^ (v & lowByteMask)) * fnvPrime64
	s = (s ^ ((v >> 8) & lowByteMask)) * fnvPrime64
	s = (s ^ ((v >> 16) & lowByteMask)) * fnvPrime64
	s = (s ^ ((v >> 24) & lowByteMask)) * fnvPrime64
	s = (s ^ ((v >> 32) & lowByteMask)) * fnvPrime64
	s = (s ^ ((v >> 40) & lowByteMask)) * fnvPrime64
	s = (s ^ ((v >> 48) & lowByteMask)) * fnvPrime64
	s = (s ^ ((v >> 56) & lowByteMask)) * fnvPrime64
	hw.state = s
}

// writeString hoists hw.state into a local so the compiler can keep it in
// a register across the loop. Through the pointer receiver it would be
// reloaded on every iteration.
// writeBytes absorbs b into the FNV-64a state. Same hoisting trick as
// writeString, but without forcing a []byte → string conversion on the
// caller (which would allocate).
func (hw *hasher) writeBytes(b []byte) {
	state := hw.state
	for i := 0; i < len(b); i++ {
		state = (state ^ uint64(b[i])) * fnvPrime64
	}
	hw.state = state
}

func (hw *hasher) writeString(s string) {
	state := hw.state
	for i := 0; i < len(s); i++ {
		state = (state ^ uint64(s[i])) * fnvPrime64
	}
	hw.state = state
}

// maxHashDepth caps both recursion and pointer/interface unwrapping so
// cyclic graphs (a map containing itself, a struct linked to itself via a
// pointer, an interface holding its own address) can neither stack-overflow
// nor spin forever. Any subtree past the cap collapses to tagNil — a
// bucket-key collision that reflect.DeepEqual still resolves correctly in
// validateHashKey, at the cost of degraded bucketing for pathologically
// deep inputs.
const maxHashDepth = 256

// hashValue streams v into hw by walking its structure via reflection.
// Each kind writes a type tag followed by its primitive bytes or its
// recursively-hashed children, so the final hw.state is an order-dependent
// FNV-64a digest of the whole value. Type tags prevent cross-kind
// collisions (e.g. the uint 65 vs the string "A").
func hashValue(hw *hasher, v reflect.Value) {
	hashValueAt(hw, v, 0)
}

func hashValueAt(hw *hasher, v reflect.Value, depth int) {
	if depth > maxHashDepth {
		hw.writeByte(tagNil)

		return
	}

	// Each unwrapping step counts towards depth just like a nesting level.
	// Without it a self-referential value (an interface holding a pointer to
	// itself) keeps this loop spinning forever: the kind stays Interface or
	// Pointer, nothing is ever nil, and maxHashDepth is never consulted.
	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		if v.IsNil() {
			hw.writeByte(tagNil)

			return
		}

		depth++
		if depth > maxHashDepth {
			hw.writeByte(tagNil)

			return
		}

		v = v.Elem()
	}

	if !v.IsValid() {
		hw.writeByte(tagNil)

		return
	}

	if hashProtoMessage(hw, v) {
		return
	}

	switch v.Kind() {
	case reflect.Bool:
		hw.writeByte(tagBool)
		if v.Bool() {
			hw.writeByte(1)
		} else {
			hw.writeByte(0)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		hw.writeByte(tagInt)
		//nolint:gosec // intentional int64 → uint64 bit reinterpretation
		hw.writeUint64(uint64(v.Int()))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		hw.writeByte(tagUint)
		hw.writeUint64(v.Uint())

	case reflect.Float32, reflect.Float64:
		hw.writeByte(tagFloat)
		hw.writeUint64(math.Float64bits(v.Float()))

	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		hw.writeByte(tagComplex)
		hw.writeUint64(math.Float64bits(real(c)))
		hw.writeUint64(math.Float64bits(imag(c)))

	case reflect.String:
		hw.writeByte(tagString)
		//nolint:gosec // Len is non-negative
		hw.writeUint64(uint64(v.Len()))
		hw.writeString(v.String())

	case reflect.Struct:
		// Walk every field, including unexported ones, to match
		// reflect.DeepEqual semantics. Safe because hashValue only uses
		// low-level accessors (.Bool, .Int, .Field, ...) which work on
		// read-only reflect.Values — it never calls .Interface(), which
		// is the only operation that panics on unexported fields.
		hw.writeByte(tagStruct)
		n := v.NumField()
		for i := 0; i < n; i++ {
			hashValueAt(hw, v.Field(i), depth+1)
		}
		hw.writeByte(tagStructEnd)

	case reflect.Slice, reflect.Array:
		hw.writeByte(tagSlice)
		n := v.Len()
		//nolint:gosec // Len is non-negative
		hw.writeUint64(uint64(n))
		for i := 0; i < n; i++ {
			hashValueAt(hw, v.Index(i), depth+1)
		}

	case reflect.Map:
		hw.writeByte(tagMap)
		//nolint:gosec // Len is non-negative
		hw.writeUint64(uint64(v.Len()))
		// XOR over per-entry sub-hashes for order independence.
		// The stack-allocated sub hasher keeps this alloc-free.
		var xored uint64
		iter := v.MapRange()
		for iter.Next() {
			sub := newHasher()
			hashValueAt(&sub, iter.Key(), depth+1)
			hashValueAt(&sub, iter.Value(), depth+1)
			xored ^= sub.state
		}
		hw.writeUint64(xored)

	case reflect.Invalid,
		reflect.Interface,
		reflect.Pointer,
		reflect.Chan,
		reflect.Func,
		reflect.UnsafePointer:
		// Invalid/Interface/Ptr are unreachable here (the top-of-function
		// unwrap loop and IsValid check already handled them). Chan, Func,
		// and UnsafePointer are opaque runtime handles that DeepEqual
		// compares by pointer identity, not structure — collapsing them
		// to tagNil forces DeepEqual to verify any bucket collision.
		hw.writeByte(tagNil)
	}
}

// assignProto stores v as a protobuf message in dst and reports whether it
// is one. Generated messages implement proto.Message on the pointer
// receiver, while the unwrap loop in hashValueAt hands us the dereferenced
// struct — hence the second attempt through the addressable form, which is
// what an element of a slice of messages ([]*T) looks like at this point.
//
// The message is handed back through a pointer argument rather than a
// return value so the helper does not hand out an interface.
func assignProto(dst *proto.Message, v reflect.Value) bool {
	if v.CanInterface() {
		if m, ok := v.Interface().(proto.Message); ok {
			*dst = m

			return true
		}
	}

	if v.CanAddr() && v.Addr().CanInterface() {
		if m, ok := v.Addr().Interface().(proto.Message); ok {
			*dst = m

			return true
		}
	}

	return false
}

// hashProtoMessage hashes a protobuf message by its deterministic wire
// encoding and reports whether it handled v.
//
// Walking a generated message field by field is not an option: its
// unexported state holds a live *MessageInfo once anything has touched the
// message through reflection (a logging interceptor, protojson, any
// ProtoReflect call). From there the walk reaches the message descriptor,
// the file descriptor, every message of that file and its imports — a
// densely connected graph whose path count explodes long before
// maxHashDepth is reached. The wire encoding, by contrast, covers exactly
// the data the message carries, and carries it in a stable field order.
//
// Equality still comes from protoEqual in the collision path: the encoding
// is a bucket key, and messages of different types can encode to the same
// bytes.
func hashProtoMessage(hw *hasher, v reflect.Value) bool {
	var m proto.Message
	if !assignProto(&m, v) {
		return false
	}

	hw.writeByte(tagProto)

	b, err := proto.MarshalOptions{Deterministic: true}.Marshal(m)
	if err != nil {
		// A message that cannot be encoded (a broken required field, a
		// cyclic Any) collapses to the bare tag: same bucket for every such
		// value, correctness restored by the equality check.
		return true
	}

	hw.writeUint64(uint64(len(b)))
	hw.writeBytes(b)

	return true
}

// protoEqual compares two values as protobuf messages, reporting whether
// both of them are messages at all. proto.Equal compares the data a
// message carries; reflect.DeepEqual would also compare the unexported
// state protobuf fills in lazily.
func protoEqual(a, b reflect.Value) (equal, ok bool) {
	var am, bm proto.Message
	if !assignProto(&am, a) || !assignProto(&bm, b) {
		return false, false
	}

	return proto.Equal(am, bm), true
}
