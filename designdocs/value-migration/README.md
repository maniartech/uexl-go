# Value Migration: Zero-Allocation Architecture

## Overview

This directory documents UExL's **Value migration journey** - the architectural transformation that achieved **zero allocations** while maintaining competitive performance. This is a **completed optimization effort** (November 2025) that established the foundation for all future performance work.

## Quick Navigation

### 🚀 **Start Here**
- **[development-guidelines.md](development-guidelines.md)** - **CRITICAL RULES** for maintaining zero allocations
- **[performance-optimization.md](performance-optimization.md)** - Complete migration journey (266ns → 227ns)

### 🏗️ **Architecture Deep Dive**
- **[value-type-system.md](value-type-system.md)** - The Value struct and how it enables zero allocations

### 📚 **Related Documentation**
- **[../performance/](../performance/)** - Comprehensive future optimization plan (100+ targets, 20-35ns goal)

---

## Achievement Summary

### Final Result (November 2025)

```
Expression: (Origin == "MOW" || Country == "RU") && (Value >= 100 || Adults == 1)

UExL:   227.1 ns/op    0 B/op    0 allocs/op  ← ONLY ZERO-ALLOC FRAMEWORK ✅
expr:   132.1 ns/op   32 B/op    1 allocs/op
celgo:  174.1 ns/op   16 B/op    1 allocs/op
```

**Achievements:**
- ✅ **Zero allocations** - Only framework with 0 allocs/op
- ✅ **14.7% faster** - 266ns → 227ns (39ns improvement)
- ✅ **3× faster pipes** - 3,428ns vs 10,588ns (expr) on map operations
- ✅ **Type safety** - Maintained explicit semantics
- ✅ **Zero panics** - Robust error handling throughout

### Migration Phases

```
Phase         Target           Result          Achievement
─────────────────────────────────────────────────────────────
Initial       Basic VM         266 ns/4 allocs Starting point
Phase 2B      VM Stack         236 ns/0 allocs 11% faster, zero allocs! ✅
Phase 2C      Constants        236 ns/0 allocs Architecture cleanup
Phase 2D      Operations       236 ns/0 allocs Value-native ops
Phase 2E      Struct Layout    227 ns/0 allocs 3.9% faster ✅

Total: 39ns faster (14.7%), infinite reduction in allocations
```

---

## What Changed

### The Problem

```go
// Before: Everything used interface{}
stack := []any{...}
constants := []any{...}

// Pop() boxed Value → any (allocation!)
func (vm *VM) Pop() any {
    vm.sp--
    return vm.stack[vm.sp]  // Returns any
}

// Every operation allocated:
right := vm.Pop()  // Allocation 1
left := vm.Pop()   // Allocation 2
result := left.(float64) + right.(float64)
vm.Push(result)    // Allocation 3
```

**Result**: 4 allocations per operation

### The Solution

```go
// After: Value type system
type Value struct {
    AnyVal   any        // Complex types only
    StrVal   string     // Inline primitive
    FloatVal float64    // Inline primitive
    Typ      valueType  // Discriminator
    BoolVal  bool       // Inline primitive
}

stack := []Value{...}
constants := []Value{...}

// Internal API - zero-alloc
func (vm *VM) popValue() Value {
    vm.sp--
    return vm.stack[vm.sp]  // Returns Value, no boxing
}

// Operations use Value directly:
right := vm.popValue()  // No allocation
left := vm.popValue()   // No allocation
result := newFloatValue(left.FloatVal + right.FloatVal)
vm.pushValue(result)    // No allocation
```

**Result**: 0 allocations

---

## Key Architectural Decisions

### 1. Value Type System

**Trade-off**: 48-byte struct vs 16-byte interface (3× larger)

**Rationale**:
- Primitives stored inline (no boxing)
- Stack copies cost ~60ns but save ~40-80ns in allocations
- Zero GC pressure > raw speed
- Better long-term performance characteristics

### 2. Internal vs Public APIs

**Pattern**: Dual API layer

```go
// Internal (zero-alloc) - opcode handlers
func (vm *VM) popValue() Value
func (vm *VM) pushValue(Value) error

// Public (boxes) - external API
func (vm *VM) Pop() any
func (vm *VM) Push(any) error
```

**Rationale**:
- Keep internal operations allocation-free
- Box only at API boundaries (unavoidable)
- Maintain backward compatibility

### 3. Field Layout Optimization

**Original** (56 bytes):
```go
type Value struct {
    Typ      valueType  // 1 byte
    // 7 bytes padding
    FloatVal float64    // 8 bytes
    StrVal   string     // 16 bytes
    BoolVal  bool       // 1 byte
    // 7 bytes padding
    AnyVal   any        // 16 bytes
}
```

**Optimized** (48 bytes):
```go
type Value struct {
    AnyVal   any        // 16 bytes, largest first
    StrVal   string     // 16 bytes
    FloatVal float64    // 8 bytes
    Typ      valueType  // 1 byte
    BoolVal  bool       // 1 byte
    // 6 bytes padding (minimal)
}
```

**Rationale**: Minimize padding by ordering fields largest to smallest

---

## Critical Rules for Developers

### ⚠️ Never Break These

1. **NEVER use `Pop()` in opcode handlers** - Always use `popValue()`
2. **NEVER box in hot paths** - Direct field access only
3. **NEVER modify Value struct layout** - Breaks optimization
4. **ALWAYS profile before optimizing** - Measure, don't guess
5. **ALWAYS test allocations** - Must be 0 for primitives
6. **ALWAYS use constructors** - `newFloatValue()`, not manual struct

**See [development-guidelines.md](development-guidelines.md) for complete rules.**

---

## Document Structure

### For New Developers

**Day 1: Understand the architecture**
1. Read this README
2. Read [value-type-system.md](value-type-system.md)
3. Review [performance-optimization.md](performance-optimization.md)

**Day 2: Learn the rules**
1. Study [development-guidelines.md](development-guidelines.md)
2. Review existing opcode handlers (`vm/vm.go`)
3. Practice with test changes

### For Code Reviews

**Checklist**:
- [ ] No `Pop()` usage in opcode handlers
- [ ] No `ToAny()` calls in hot paths
- [ ] All allocations tested (must be 0)
- [ ] Value constructors used
- [ ] Tests pass, benchmarks stable

### For Optimizations

**Before optimizing**:
1. Profile to find bottlenecks
2. Verify >5% CPU time in target
3. Establish baseline benchmarks

**After optimizing**:
1. Run allocation tests
2. Compare benchmarks
3. Update documentation

---

## Testing Your Changes

### Required Tests

```bash
# 1. Correctness
go test ./...

# 2. Race detection
go test ./... -race

# 3. Allocations (CRITICAL)
go test -bench=BenchmarkYourFeature -benchmem
# MUST show: 0 B/op, 0 allocs/op

# 4. Performance
go test -bench=BenchmarkYourFeature -count=10 > current.txt
benchstat baseline.txt current.txt
```

**All must pass.**

---

## Performance Budgets

### Must Maintain:
- **Primitive allocations**: 0 allocs/op (CRITICAL)
- **Main benchmark**: 180-230 ns/op
- **Value struct size**: ≤ 48 bytes

### Investigation Triggers:
- Any allocation in primitive ops
- Main benchmark > 250 ns/op
- Regression > 5%
- Struct size increase

---

## Common Mistakes

### 1. Using Pop() Instead of popValue()

```go
// ❌ WRONG
case code.OpMyOp:
    val := vm.Pop()  // Allocates!

// ✅ CORRECT
case code.OpMyOp:
    val := vm.popValue()  // Zero-alloc
```

### 2. Boxing in Hot Paths

```go
// ❌ WRONG
func process(val Value) {
    anyVal := val.ToAny()  // Allocates!
    // ...
}

// ✅ CORRECT
func process(val Value) {
    switch val.Typ {
    case TypeFloat:
        // Use val.FloatVal directly
    }
}
```

### 3. Manual Struct Construction

```go
// ❌ WRONG
val := Value{Typ: TypeFloat, FloatVal: 42.0}

// ✅ CORRECT
val := newFloatValue(42.0)
```

---

## Future Work

### Completed in This Migration:
- ✅ Zero-allocation architecture
- ✅ Value type system
- ✅ VM stack migration
- ✅ Value-native operations
- ✅ Struct layout optimization

### Next Steps (See ../performance/):
- 🔴 VM core optimizations
- 🔴 Operator optimizations
- 🔴 Pipe optimizations
- 🔴 Built-in function optimizations
- 🔴 Compiler optimizations

**Target**: 20-35ns/op across ALL operations

---

## Success Metrics

**This migration achieved:**

✅ **Zero allocations** - Infinite improvement over initial 4 allocs/op
✅ **14.7% faster** - 266ns → 227ns
✅ **Industry-leading** - Only zero-alloc framework in benchmarks
✅ **Maintainable** - Clear patterns, documented trade-offs
✅ **Type-safe** - No unsafe operations
✅ **Zero panics** - Robust error handling

**All without breaking changes to public API.**

---

## Contributing

When modifying VM code:

1. Read [development-guidelines.md](development-guidelines.md)
2. Use Value-native operations
3. Profile if optimizing
4. Test allocations (must be 0)
5. Document trade-offs

---

## References

- **Value Type System**: [value-type-system.md](value-type-system.md)
- **Optimization Journey**: [performance-optimization.md](performance-optimization.md)
- **Development Rules**: [development-guidelines.md](development-guidelines.md)
- **Future Optimizations**: [../performance/README.md](../performance/README.md)

---

**Last Updated:** November 13, 2025
**Status:** ✅ Complete - Zero allocations achieved and documented
**Next Phase**: Comprehensive VM optimization (see ../performance/)
