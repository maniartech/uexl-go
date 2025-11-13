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

// Type returns the value type
func (v Value) Type() valueType {
	return v.Typ
}
