package types

// Value represents a stack value with type information to avoid interface boxing for primitives.
// Primitives (float64, string, bool) are stored directly without boxing.
// Complex types (arrays, maps, functions) still use any interface since they're already heap-allocated.
//
// Field order optimized to minimize padding (40 bytes vs 56 bytes unoptimized):
// - 16-byte aligned fields first (AnyVal, StrVal)
// - 8-byte aligned field next (FloatVal)
// - Small fields last (Typ, BoolVal)
type Value struct {
	// Complex types - boxed, but unavoidable for reference types like []any, map[string]any
	AnyVal any // 16 bytes, offset 0

	// String primitive - stored inline
	StrVal string // 16 bytes, offset 16

	// Numeric primitive - stored inline
	FloatVal float64 // 8 bytes, offset 32

	// Type discriminator and boolean primitive
	Typ     valueType // 1 byte, offset 40
	BoolVal bool      // 1 byte, offset 41
	// implicit 6 bytes padding to align to 48 bytes
}

type valueType uint8

const (
	TypeFloat valueType = iota
	TypeString
	TypeBool
	TypeAny // For arrays, maps, functions, and other complex types
	TypeNull

	// Temporal primitives. The canonical signed int64 millisecond count is stored inline in
	// FloatVal (never AnyVal): every value in the portable range (|ms| < ~2.5e14) is well under
	// 2^53, so the int64<->float64 round-trip is exact. This keeps datetime/duration zero-alloc.
	// See docs/specs/datetime-spec.md §2.
	TypeDateTime // an instant: ms since 1970-01-01T00:00:00.000Z (UTC)
	TypeDuration // an exact elapsed span in ms (may be negative)
)

// Constructors for primitive types - zero allocations

func NewFloatValue(f float64) Value {
	return Value{Typ: TypeFloat, FloatVal: f}
}

func NewStringValue(s string) Value {
	return Value{Typ: TypeString, StrVal: s}
}

func NewBoolValue(b bool) Value {
	return Value{Typ: TypeBool, BoolVal: b}
}

func NewNullValue() Value {
	return Value{Typ: TypeNull}
}

// NewDateTimeValue builds a datetime Value from a canonical millisecond instant.
// Stored inline in FloatVal (zero allocation). ms is expected within the portable range
// (datetime-spec §5.2) so the int64->float64 round-trip is exact.
func NewDateTimeValue(ms int64) Value {
	return Value{Typ: TypeDateTime, FloatVal: float64(ms)}
}

// NewDurationValue builds a duration Value from a canonical millisecond span (may be negative).
func NewDurationValue(ms int64) Value {
	return Value{Typ: TypeDuration, FloatVal: float64(ms)}
}

// Constructor for complex types - still boxes but only for non-primitives

func NewAnyValue(v any) Value {
	if v == nil {
		return NewNullValue()
	}

	// Fast path for primitives - avoid boxing
	switch val := v.(type) {
	case float64:
		return NewFloatValue(val)
	case string:
		return NewStringValue(val)
	case bool:
		return NewBoolValue(val)
	case int:
		return NewFloatValue(float64(val))
	case DateTime:
		return NewDateTimeValue(val.Millis)
	case Duration:
		return NewDurationValue(val.Millis)
	default:
		// For arrays, maps, functions, etc. - box them
		return Value{Typ: TypeAny, AnyVal: v}
	}
}

// Converters - extract values from Value

// ToAny converts value back to any interface (for compatibility)
func (v Value) ToAny() any {
	switch v.Typ {
	case TypeFloat:
		return v.FloatVal
	case TypeString:
		return v.StrVal
	case TypeBool:
		return v.BoolVal
	case TypeNull:
		return nil
	case TypeAny:
		return v.AnyVal
	case TypeDateTime:
		return DateTime{Millis: int64(v.FloatVal)}
	case TypeDuration:
		return Duration{Millis: int64(v.FloatVal)}
	default:
		return nil
	}
}

// Type-safe extractors

func (v Value) AsFloat() (float64, bool) {
	if v.Typ == TypeFloat {
		return v.FloatVal, true
	}
	return 0, false
}

func (v Value) AsString() (string, bool) {
	if v.Typ == TypeString {
		return v.StrVal, true
	}
	return "", false
}

func (v Value) AsBool() (bool, bool) {
	if v.Typ == TypeBool {
		return v.BoolVal, true
	}
	return false, false
}

func (v Value) AsAny() (any, bool) {
	if v.Typ == TypeAny {
		return v.AnyVal, true
	}
	return nil, false
}

// AsDateTimeMillis returns the canonical millisecond instant when v is a datetime.
func (v Value) AsDateTimeMillis() (int64, bool) {
	if v.Typ == TypeDateTime {
		return int64(v.FloatVal), true
	}
	return 0, false
}

// AsDurationMillis returns the canonical millisecond span when v is a duration.
func (v Value) AsDurationMillis() (int64, bool) {
	if v.Typ == TypeDuration {
		return int64(v.FloatVal), true
	}
	return 0, false
}

// Type checkers

func (v Value) IsFloat() bool {
	return v.Typ == TypeFloat
}

func (v Value) IsString() bool {
	return v.Typ == TypeString
}

func (v Value) IsBool() bool {
	return v.Typ == TypeBool
}

func (v Value) IsNull() bool {
	return v.Typ == TypeNull
}

func (v Value) IsAny() bool {
	return v.Typ == TypeAny
}

// IsDateTime reports whether v is a datetime instant.
func (v Value) IsDateTime() bool {
	return v.Typ == TypeDateTime
}

// IsDuration reports whether v is a duration span.
func (v Value) IsDuration() bool {
	return v.Typ == TypeDuration
}

// Type returns the value type
func (v Value) Type() valueType {
	return v.Typ
}
