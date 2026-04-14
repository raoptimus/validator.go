package validator

import (
	"math"
	"reflect"
)

// Inlined FNV-64a constants. Using the raw state instead of hash.Hash64
// avoids the interface allocation and per-field Reset cost of fnv.New64a.
const (
	fnvOffset64 = 14695981039346656037
	fnvPrime64  = 1099511628211
	bitsPerByte = 8
	uint64Bytes = 8
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
	hw.state ^= uint64(b)
	hw.state *= fnvPrime64
}

// writeUint64 absorbs the 8 little-endian bytes of v into the FNV-64a state.
func (hw *hasher) writeUint64(v uint64) {
	for i := 0; i < uint64Bytes; i++ {
		//nolint:gosec // intentional truncation to low 8 bits
		hw.writeByte(byte(v >> (i * bitsPerByte)))
	}
}

func (hw *hasher) writeString(s string) {
	for i := 0; i < len(s); i++ {
		hw.state ^= uint64(s[i])
		hw.state *= fnvPrime64
	}
}

// maxHashDepth caps recursion so cyclic graphs (a map containing itself,
// a struct linked to itself via a pointer) cannot stack-overflow. Any
// subtree past the cap collapses to tagNil — a bucket-key collision that
// reflect.DeepEqual still resolves correctly in validateHashKey, at the
// cost of degraded bucketing for pathologically deep inputs.
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

	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			hw.writeByte(tagNil)

			return
		}

		v = v.Elem()
	}

	if !v.IsValid() {
		hw.writeByte(tagNil)

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
		reflect.Ptr,
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
