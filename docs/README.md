# UExL Documentation

## Overview

This directory contains user-facing documentation and design documentation.

## Documentation Structure

### Design Documentation
- **[../designdocs/value-migration/](../designdocs/value-migration/)** - Zero-allocation architecture (completed)
- **[../designdocs/performance/](../designdocs/performance/)** - Comprehensive optimization plan (100+ targets, future work)

### Topic-Specific Documentation
- **[handling-nullish-col.md](handling-nullish-col.md)** - Nullish coalescing semantics

---

## Performance Achievement Summary

### Current State (November 2025)

```
Expression: (Origin == "MOW" || Country == "RU") && (Value >= 100 || Adults == 1)

UExL:   227.1 ns/op    0 B/op    0 allocs/op  ← ONLY ZERO-ALLOC FRAMEWORK ✅
expr:   132.1 ns/op   32 B/op    1 allocs/op
celgo:  174.1 ns/op   16 B/op    1 allocs/op
```

**Key Metrics:**
- ✅ **Zero allocations** - Only framework with 0 allocs/op
- ✅ **Competitive speed** - Within 72% of fastest (acceptable trade-off)
- ✅ **3× faster pipes** - 3,428ns vs 10,588ns (expr) on map operations
- ✅ **Type safety** - Explicit nullish/boolish semantics
- ✅ **Zero panics** - Robust error handling throughout

### The Journey

```
Phase         Target           Result          Achievement
─────────────────────────────────────────────────────────────
Initial       Basic VM         266 ns/4 allocs Starting point
Phase 2B      VM Stack         236 ns/0 allocs 11% faster, zero allocs! ✅
Phase 2C      Constants        236 ns/0 allocs Architecture cleanup
Phase 2D      Operations       236 ns/0 allocs Value-native ops
Phase 2E      Struct Layout    227 ns/0 allocs 3.9% faster ✅

Total Improvement: 39ns faster (14.7%), infinite reduction in allocations
```

---

## Document Structure

### For Development Work

**Before coding:**
1. Read **[development-guidelines.md](development-guidelines.md)** - Learn the 6 golden rules
2. Review **[value-type-system.md](value-type-system.md)** - Understand the Value architecture

**During coding:**
- Keep **[development-guidelines.md](development-guidelines.md)** open for quick reference
- Check **[performance-optimization.md](performance-optimization.md)** for patterns

**After coding:**
- Run all validation checks from **[development-guidelines.md](development-guidelines.md)**
- Verify zero allocations maintained

### For Understanding Architecture

1. **[value-type-system.md](value-type-system.md)** - Why Value struct exists, how it works
2. **[performance-optimization.md](performance-optimization.md)** - How we achieved zero allocations
3. **[development-guidelines.md](development-guidelines.md)** - How to maintain it

### For Planning Future Work

- **[../designdocs/performance/](../designdocs/performance/)** - Comprehensive optimization roadmap
  - 100+ optimization targets across 10 categories
  - VM core, operators, pipes, built-ins, etc.
  - Targets: 20-35ns/op across ALL operations

---

## Critical Rules (Never Forget)

### ⚠️ Rule 1: NEVER Use `Pop()` in Opcode Handlers

```go
// ❌ WRONG - Allocates!
case code.OpMyOp:
    val := vm.Pop()  // Boxes Value → any

// ✅ CORRECT - Zero-alloc
case code.OpMyOp:
    val := vm.popValue()  // Returns Value directly
```

### ⚠️ Rule 2: NEVER Box in Hot Paths

```go
// ❌ WRONG - Allocates!
anyVal := val.ToAny()  // Inside opcode handler

// ✅ CORRECT - Direct access
switch val.Typ {
case TypeFloat:
    result := val.FloatVal * 2  // No boxing
}
```

### ⚠️ Rule 3: NEVER Modify Value Struct Layout

```go
// Current optimized layout (48 bytes):
type Value struct {
    AnyVal   any        // 16 bytes, largest first
    StrVal   string     // 16 bytes
    FloatVal float64    // 8 bytes
    Typ      valueType  // 1 byte
    BoolVal  bool       // 1 byte
}

// ❌ DON'T reorder fields - breaks optimization!
```

### ⚠️ Rule 4: ALWAYS Profile Before Optimizing

```bash
# Required before any optimization:
go test -bench=BenchmarkX -cpuprofile=before.prof
go tool pprof -top before.prof
# Only optimize if function shows >5% CPU time
```

### ⚠️ Rule 5: ALWAYS Test Allocations

```bash
# Every change must verify:
go test -bench=BenchmarkX -benchmem
# Must show: 0 B/op, 0 allocs/op for primitives
```

### ⚠️ Rule 6: ALWAYS Use Constructors

```go
// ❌ WRONG - Error-prone
val := Value{Typ: TypeFloat, FloatVal: 42.0}

// ✅ CORRECT - Guaranteed correct
val := newFloatValue(42.0)
```

**See [development-guidelines.md](development-guidelines.md) for complete rules.**

---

## Performance Budgets

### Must Maintain:
- **Primitive allocations**: 0 allocs/op (CRITICAL - non-negotiable)
- **Main benchmark**: 180-230 ns/op (currently 227ns ✅)
- **Struct size**: ≤ 48 bytes (Value struct)

### Investigation Triggers:
- Any allocation in primitive operations
- Main benchmark > 250 ns/op
- Any regression > 5% without clear benefit
- Value struct size increase

---

## Common Pitfalls

### 1. Adding New Opcodes

```go
// ❌ WRONG
case code.OpNewThing:
    val := vm.Pop()  // Allocates!

// ✅ CORRECT
case code.OpNewThing:
    val := vm.popValue()  // Zero-alloc
    vm.processValue(val)
```

### 2. Adding New Helpers

```go
// ❌ WRONG - Accepts any
func (vm *VM) helper(val any) error

// ✅ CORRECT - Accepts Value
func (vm *VM) helperValue(val Value) error
```

### 3. Returning Values

```go
// Public API must box (unavoidable):
func (vm *VM) Run(...) (any, error) {
    // Internal uses Value (zero-alloc)
    err := vm.run()

    // Only box at API boundary
    return vm.LastPoppedStackElem(), nil
}
```

---

## How We Got Here

### The Problem (Initial State)

```
Benchmark comparison showed:
UExL:   266 ns/op   ? B/op    4 allocs/op  ← Allocations!
expr:   129 ns/op  32 B/op    1 allocs/op
celgo:  174 ns/op  16 B/op    1 allocs/op
```

**Root cause**: `Pop()` was boxing `Value → any` four times per operation

### The Solution (Value Migration)

**Phase 2B - VM Stack Migration:**
- Changed `stack []any` → `stack []Value`
- Eliminated stack-level boxing
- Result: 0 allocations for boolean ops ✅

**Phase 2C - Compiler Constants:**
- Changed `Constants []any` → `Constants []Value`
- Zero-alloc constant loading
- Result: Architecture cleanup ✅

**Phase 2D - Value-Native Operations:**
- Opcode handlers use `popValue()` instead of `Pop()`
- Direct Value comparisons
- Result: ALL operations zero-alloc ✅

**Phase 2E - Struct Optimization:**
- Reordered Value fields (largest first)
- Eliminated 14 bytes padding (56→48 bytes)
- Result: 3.9% speed improvement ✅

### The Result

```
UExL:   227 ns/op    0 B/op    0 allocs/op  ← ZERO ALLOCATIONS! 🏆
```

**Only framework in comparison benchmarks with zero allocations.**

---

## What Makes UExL Different

### 1. Value Type System

Instead of `interface{}` (boxes primitives), UExL uses a discriminated union:

```go
type Value struct {
    AnyVal   any        // Complex types (arrays, maps)
    StrVal   string     // String primitive (inline)
    FloatVal float64    // Number primitive (inline)
    Typ      valueType  // Type discriminator
    BoolVal  bool       // Boolean primitive (inline)
}
```

**Trade-off**: 48-byte struct vs 16-byte interface (3× larger)
**Benefit**: Zero allocations for primitives (infinite improvement)

### 2. Internal vs Public APIs

```go
// Internal (zero-alloc) - opcode handlers use this
func (vm *VM) popValue() Value

// Public (boxes) - external API uses this
func (vm *VM) Pop() any {
    return vm.popValue().ToAny()  // Boxing only at boundary
}
```

**Pattern**: Keep internal operations allocation-free, box only at API boundaries

### 3. Context Variable Caching

```go
// Compiler stores variable names in order:
bytecode.ContextVars = ["Origin", "Country", "Value", "Adults"]

// Runtime pre-resolves to array:
vm.contextVarCache[0] = contextValues["Origin"]  // Once
vm.contextVarCache[1] = contextValues["Country"]
// ...

// Execution uses O(1) array access:
value := vm.contextVarCache[varIndex]  // Fast!
```

**Benefit**: Eliminated 14% map lookup overhead vs competitors

---

## Future Work

### Completed (This Documentation Covers):
- ✅ Zero-allocation architecture
- ✅ Value type system
- ✅ VM stack migration
- ✅ Value-native operations
- ✅ Struct layout optimization

### Planned (See designdocs/performance/):
- 🔴 VM core optimizations (instruction dispatch, frame pooling)
- 🔴 Operator optimizations (arithmetic, string, bitwise)
- 🔴 Pipe optimizations (filter, reduce, etc.)
- 🔴 Built-in function optimizations (50+ functions)
- 🔴 Compiler optimizations (constant folding, type hints)

**Target**: 20-35ns/op across ALL operations (from current 227ns)

See **[../designdocs/performance/README.md](../designdocs/performance/README.md)** for complete roadmap.

---

## Quick Reference

**Need to know what NOT to do?**
→ [development-guidelines.md](development-guidelines.md)

**Want to understand the architecture?**
→ [value-type-system.md](value-type-system.md)

**Looking for optimization history?**
→ [performance-optimization.md](performance-optimization.md)

**Planning future optimizations?**
→ [../designdocs/performance/](../designdocs/performance/)

---

## Testing Your Changes

### Minimum Required Tests:

```bash
# 1. Correctness
go test ./...

# 2. Race detection
go test ./... -race

# 3. Allocation check (CRITICAL)
go test -bench=BenchmarkYourFeature -benchmem
# MUST show: 0 B/op, 0 allocs/op for primitives

# 4. Performance comparison
go test -bench=BenchmarkYourFeature -count=10 > current.txt
benchstat baseline.txt current.txt
# MUST show: no regression > 5%
```

**All must pass before committing.**

---

## Contributing

When adding features or fixes:

1. **Read development-guidelines.md first**
2. **Use Value-native operations** (not `any`)
3. **Profile if optimizing** (never guess)
4. **Test allocations** (must be 0 for primitives)
5. **Document trade-offs** (explain why)

---

## Summary

### The Formula for Zero Allocations:
1. ✅ **Value struct** - Inline primitives, box only complex types
2. ✅ **Internal APIs** - Work with Value directly
3. ✅ **Public APIs** - Box only at boundaries
4. ✅ **Type-specific paths** - Direct field access, no assertions
5. ✅ **Profile-driven** - Measure, optimize, verify

### Key Metrics to Protect:
- **Zero allocations** for primitives (CRITICAL)
- Competitive speed (within 2× of fastest)
- Clean, maintainable code
- Comprehensive test coverage

### Remember:
> "The best optimization is architecture that makes allocation unnecessary."

UExL achieved zero allocations through careful design, not clever tricks.

**Keep it that way.**

---

**Last Updated:** November 13, 2025
**Status:** Zero allocations achieved ✅, comprehensive documentation complete ✅
