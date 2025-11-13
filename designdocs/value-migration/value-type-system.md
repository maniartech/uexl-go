# Value Type System Architecture

## Overview

The Value type system is the foundation of UExL's zero-allocation performance. This document explains its design, implementation, and critical usage patterns.

## Core Design

### The Value Struct

```go
// Location: types/value.go
type Value struct {
    AnyVal   any        // 16 bytes, offset 0  - Complex types
    StrVal   string     // 16 bytes, offset 16 - String primitive
    FloatVal float64    // 8 bytes, offset 32  - Number primitive
    Typ      valueType  // 1 byte, offset 40   - Type discriminator
    BoolVal  bool       // 1 byte, offset 41   - Boolean primitive
    // implicit 6 bytes padding to 48 bytes total
}
```

**Size**: 48 bytes (optimized from original 56 bytes)

### Type Discriminator

```go
type valueType uint8

const (
    TypeFloat  valueType = iota  // Numeric values
    TypeString                   // String values
    TypeBool                     // Boolean values
    TypeAny                      // Arrays, maps, functions, objects
    TypeNull                     // null literal
)
```

## Design Philosophy

### Why Not Just Use `interface{}`?

**Problem with `interface{}`**:
```go
// Every operation boxes/unboxes
stack := make([]interface{}, 1024)
stack[sp] = 42.0           // Boxes float64 → interface{} (allocation!)
val := stack[sp].(float64) // Unboxes interface{} → float64 (allocation!)
```

**Solution with Value**:
```go
// Primitives stored inline, no boxing
stack := make([]Value, 1024)
stack[sp] = Value{Typ: TypeFloat, FloatVal: 42.0}  // No allocation!
val := stack[sp].FloatVal                          // Direct access!
```

### The Trade-off

| Aspect | `interface{}` | `Value` struct |
|--------|---------------|----------------|
| Size | 16 bytes | 48 bytes |
| Primitive storage | Boxed (allocates) | Inline (zero-alloc) |
| Type info | Implicit | Explicit |
| Copy cost | Low | Higher |
| Allocation cost | High | Zero |

**Result**: Accept 3× larger stack values to achieve zero allocations

## Field Layout Optimization

### Why Field Order Matters

Go aligns struct fields to their natural boundaries:
- 1-byte types (bool, uint8): any alignment
- 8-byte types (float64, int64): 8-byte boundary
- 16-byte types (string, any): 8-byte boundary (contains pointers)

**Bad layout** (original):
```go
type Value struct {
    Typ      valueType  // offset 0, size 1
    // 7 bytes padding for FloatVal alignment!
    FloatVal float64    // offset 8, size 8
    StrVal   string     // offset 16, size 16
    BoolVal  bool       // offset 32, size 1
    // 7 bytes padding for AnyVal alignment!
    AnyVal   any        // offset 40, size 16
}
// Total: 56 bytes (14 bytes wasted!)
```

**Good layout** (current):
```go
type Value struct {
    AnyVal   any        // offset 0, size 16 (largest first)
    StrVal   string     // offset 16, size 16
    FloatVal float64    // offset 32, size 8
    Typ      valueType  // offset 40, size 1
    BoolVal  bool       // offset 41, size 1
    // 6 bytes padding (minimal)
}
// Total: 48 bytes (14% smaller!)
```

**Rule**: Order fields from largest to smallest to minimize padding

## Constructors

### Why Use Constructors?

**❌ Manual construction** (DON'T DO):
```go
val := Value{
    Typ:      TypeFloat,
    FloatVal: 42.0,
    StrVal:   "",      // Forgot to initialize! Bug waiting to happen
    BoolVal:  false,
    AnyVal:   nil,
}
```

**✅ Constructor** (ALWAYS USE):
```go
val := newFloatValue(42.0)
// Guarantees correct initialization
```

### Available Constructors

```go
// Primitives (zero-alloc)
func NewFloatValue(f float64) Value
func NewStringValue(s string) Value
func NewBoolValue(b bool) Value
func NewNullValue() Value

// Complex types (boxes, but unavoidable)
func NewAnyValue(v any) Value
```

**Critical**: VM package re-exports these as lowercase for internal use:
```go
// vm/value.go
var (
    newFloatValue  = types.NewFloatValue
    newStringValue = types.NewStringValue
    newBoolValue   = types.NewBoolValue
    newNullValue   = types.NewNullValue
    newAnyValue    = types.NewAnyValue
)
```

## Type Checking

### Fast Type Checks

```go
// Check type
if val.Typ == TypeFloat {
    result := val.FloatVal * 2  // Direct access
}

// Or use helper methods
if val.IsFloat() {
    result := val.AsFloat() * 2  // Safe accessor
}
```

### Type Check Methods

```go
func (v Value) IsFloat() bool   { return v.Typ == TypeFloat }
func (v Value) IsString() bool  { return v.Typ == TypeString }
func (v Value) IsBool() bool    { return v.Typ == TypeBool }
func (v Value) IsNull() bool    { return v.Typ == TypeNull }
func (v Value) IsAny() bool     { return v.Typ == TypeAny }
```

### Safe Accessor Methods

```go
func (v Value) AsFloat() float64 {
    if v.Typ == TypeFloat {
        return v.FloatVal
    }
    return 0  // Default for wrong type
}

// Similar for AsString(), AsBool(), AsAny()
```

## Conversion to `any`

### When to Convert

**✅ ONLY at these boundaries**:
1. Final return from `VM.Run()`
2. Public API methods (`Pop()`, `Top()`)
3. Building arrays/objects (contain `any` elements)
4. Calling external functions

**❌ NEVER inside opcode handlers**

### The `ToAny()` Method

```go
func (v Value) ToAny() any {
    switch v.Typ {
    case TypeFloat:
        return v.FloatVal    // Boxes float64
    case TypeString:
        return v.StrVal      // Boxes string
    case TypeBool:
        return v.BoolVal     // Boxes bool
    case TypeNull:
        return nil           // No boxing
    case TypeAny:
        return v.AnyVal      // Already boxed
    }
    return nil
}
```

**Critical**: This allocates! Use sparingly.

## Usage Patterns

### Pattern 1: Loading Constants

```go
// Compiler stores as Value
bytecode.Constants = []Value{
    newFloatValue(42.0),
    newStringValue("hello"),
}

// VM loads without boxing
case code.OpConstant:
    constIndex := code.ReadUint16(...)
    vm.pushValue(vm.constants[constIndex])  // Zero-alloc!
```

### Pattern 2: Context Variables

```go
// VM converts context at startup
for i, varName := range vm.contextVars {
    vm.contextVarCache[i] = newAnyValue(contextVarsValues[varName])
}

// VM loads without boxing
case code.OpContextVar:
    varIndex := code.ReadUint16(...)
    vm.pushValue(vm.contextVarCache[varIndex])  // Zero-alloc!
```

### Pattern 3: Binary Operations

```go
// Opcode handler uses Value-native operations
case code.OpEqual:
    right, left := vm.pop2Values()  // Returns Values directly
    vm.executeComparisonOperationValues(opcode, left, right)
```

### Pattern 4: Type-Specific Operations

```go
func (vm *VM) executeComparisonOperationValues(op code.Opcode, left, right Value) error {
    if left.Typ == right.Typ {
        switch left.Typ {
        case TypeFloat:
            // Direct field access, no boxing
            return vm.executeNumberComparisonOperation(op, left.FloatVal, right.FloatVal)
        case TypeString:
            return vm.executeStringComparisonOperation(op, left.StrVal, right.StrVal)
        // ...
        }
    }
    // Fallback for mixed types (boxes here)
    return vm.executeComparisonOperation(op, left.ToAny(), right.ToAny())
}
```

## Critical Don'ts

### ❌ DON'T: Reorder Fields

```go
// This breaks optimization!
type Value struct {
    Typ      valueType  // Moving small field first
    FloatVal float64    // Adds 7 bytes padding!
    // ...
}
```

**Why**: Increases struct size, slower copies, breaks cache efficiency

### ❌ DON'T: Add New Fields Without Analysis

```go
// Don't do this without checking struct size!
type Value struct {
    AnyVal   any
    StrVal   string
    FloatVal float64
    IntVal   int64      // NEW - might add padding!
    Typ      valueType
    BoolVal  bool
}
```

**Required**: Run `unsafe.Sizeof(Value{})` and check padding

### ❌ DON'T: Use Zero Value Directly

```go
// Ambiguous - what type is this?
var val Value  // All zeros - could be TypeFloat with 0.0, or TypeNull?
```

**Always use constructors** to make type explicit

### ❌ DON'T: Box in Hot Paths

```go
// Inside opcode handler (hot path)
val := vm.popValue()
anyVal := val.ToAny()  // ❌ ALLOCATES!
```

## Testing Value Types

### Unit Test Pattern

```go
func TestValueConstruction(t *testing.T) {
    // Test primitive types
    floatVal := newFloatValue(42.0)
    assert.Equal(t, TypeFloat, floatVal.Typ)
    assert.Equal(t, 42.0, floatVal.FloatVal)

    // Test conversion
    anyVal := floatVal.ToAny()
    assert.Equal(t, 42.0, anyVal.(float64))
}
```

### Allocation Test Pattern

```go
func TestValueNoAlloc(t *testing.T) {
    // Zero-alloc operations shouldn't allocate
    allocs := testing.AllocsPerRun(1000, func() {
        val := newFloatValue(42.0)
        _ = val.FloatVal  // Direct access
    })
    assert.Equal(t, 0.0, allocs, "Should not allocate")
}
```

### Benchmark Pattern

```go
func BenchmarkValueDirect(b *testing.B) {
    val := newFloatValue(42.0)
    b.ReportAllocs()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        _ = val.FloatVal  // Should be 0 allocs
    }
}
```

## Size Monitoring

### Check Struct Size

```go
package main

import (
    "fmt"
    "unsafe"
    "github.com/maniartech/uexl/types"
)

func main() {
    var v types.Value
    fmt.Printf("Value size: %d bytes\n", unsafe.Sizeof(v))

    // Should be 48 bytes
    if unsafe.Sizeof(v) > 48 {
        panic("Value struct too large!")
    }
}
```

**Add this to CI/CD** to prevent regressions

## Future Considerations

### If We Need More Types

**Option 1: Add to existing struct** (simple)
```go
type Value struct {
    AnyVal   any
    StrVal   string
    FloatVal float64
    IntVal   int64     // NEW - adds 8 bytes (total: 56 bytes)
    Typ      valueType
    BoolVal  bool
}
```

**Option 2: Union with unsafe** (complex)
```go
type Value struct {
    typ  valueType
    data [16]byte  // Union of all types
}
```

**Recommendation**: Stick with Option 1 unless struct exceeds 64 bytes

### Performance Monitoring

Track these metrics:
- Value struct size (must be ≤ 48 bytes)
- Allocation count (must be 0 for primitives)
- Copy overhead (benchmark stack operations)

## Summary

### Key Principles:
1. ✅ Primitives inline, complex types boxed
2. ✅ Largest fields first (minimize padding)
3. ✅ Always use constructors
4. ✅ Never box in hot paths
5. ✅ Monitor struct size continuously

### The Contract:
- Value struct enables zero allocations
- In exchange: larger stack values (48 vs 16 bytes)
- Trade-off accepted: architectural advantage > micro-optimization

### Remember:
> "The Value type is the foundation. Protect it, optimize it, but don't break it."
