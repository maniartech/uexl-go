# Performance Optimization Guide

## Overview

UExL achieved **zero allocations** and competitive performance through systematic architectural optimization. This document explains the journey, techniques, and critical lessons for future development.

## Performance Achievements

### Benchmark Results
```
Framework  Time (ns/op)  Memory (B/op)  Allocs/op  Status
--------------------------------------------------------------
expr       132.1         32             1          Fastest raw speed
celgo      174.1         16             1          #2 speed
UExL       227.1         0              0          ZERO ALLOCATIONS ✅
```

### Evolution Timeline
```
Initial (comparison project):  9,388 ns/op, unknown allocs
Phase 1 optimizations:           266 ns/op, 4 allocs/op
Phase 2D (Value migration):      236 ns/op, 0 allocs/op ✅
Phase 2E (struct optimization):  227 ns/op, 0 allocs/op ✅

Total improvement: 41× faster, infinite reduction in allocations
```

## The Optimization Journey

### Phase 1: Failed Approach (Aborted)
**Attempt**: Type-specific push methods (`pushFloat64()`, `pushString()`, etc.)
**Result**: 13.6% regression on string operations
**Lesson**: Micro-optimizations without architectural change don't work
**Root cause**: Interface boxing still occurred when converting back

### Phase 2: Architectural Redesign (Success)

#### Phase 2A: Understanding the Problem
**Discovery**: All allocations came from `Value.ToAny()` during `Pop()` calls
**Profiling output**:
```
99.73% of allocations: Value.ToAny() (inline)
  79.35% from vm.Pop()
  20.38% from vm.LastPoppedStackElem()
```

#### Phase 2B: VM Stack Migration
**Change**: `stack []any` → `stack []Value`
**Impact**: Eliminated stack-level boxing
**Result**: Boolean operations → 0 allocations ✅

#### Phase 2C: Compiler Constants Migration
**Change**: `Constants []any` → `Constants []Value`
**Innovation**: `pushValue(Value)` method for zero-alloc constant loading
**Result**: Constant loading → 0 allocations ✅

#### Phase 2D: Value-Native Operations
**Change**: Rewrite opcode handlers to use `popValue()` instead of `Pop()`
**Key handlers**:
- `executeComparisonOperationValues()` - Direct Value comparison
- `isTruthyValue()` - Value-native truthiness checking
- `pop2Values()` - Batch popping for binary operations

**Result**: ALL operations → 0 allocations ✅

#### Phase 2E: Struct Layout Optimization
**Discovery**: Value struct had 14 bytes of wasted padding (56 bytes total)
**Fix**: Reordered fields by size (largest first)
```go
// Before: 56 bytes
type Value struct {
    Typ      valueType  // 1 byte
    FloatVal float64    // 8 bytes
    StrVal   string     // 16 bytes
    BoolVal  bool       // 1 byte
    AnyVal   any        // 16 bytes
}

// After: 48 bytes (14% smaller)
type Value struct {
    AnyVal   any        // 16 bytes, offset 0
    StrVal   string     // 16 bytes, offset 16
    FloatVal float64    // 8 bytes, offset 32
    Typ      valueType  // 1 byte, offset 40
    BoolVal  bool       // 1 byte, offset 41
    // 6 bytes padding to 48
}
```
**Result**: 3.9% speed improvement (236ns → 227ns)

## Critical Architecture Patterns

### 1. Value Type System

**Design**: Discriminated union storing primitives inline
```go
type Value struct {
    AnyVal   any        // For complex types (arrays, maps)
    StrVal   string     // Inline string storage
    FloatVal float64    // Inline number storage
    Typ      valueType  // Type discriminator
    BoolVal  bool       // Inline boolean storage
}
```

**Benefits**:
- Primitives stored without boxing
- Type information explicit
- Zero allocations for float64, string, bool

**Trade-off**: 48-byte struct vs 16-byte interface{} (3× larger)

### 2. Internal vs Public APIs

**Pattern**: Separate internal zero-alloc operations from public boxed APIs

```go
// Internal (zero-alloc)
func (vm *VM) popValue() Value
func (vm *VM) pushValue(val Value) error
func (vm *VM) pop2Values() (Value, Value)

// Public (boxes for external API)
func (vm *VM) Pop() any {
    return vm.popValue().ToAny()  // Only boxes at API boundary
}
```

**Critical Rule**: Opcode handlers MUST use internal APIs only

### 3. Context Variable Optimization

**Innovation**: Array-based cache instead of map lookups

```go
// Compile time: Store variable names in order
bytecode.ContextVars = ["Origin", "Country", "Value", "Adults"]

// Runtime: Pre-resolve to array
vm.contextVarCache = make([]Value, len(vm.contextVars))
for i, varName := range vm.contextVars {
    vm.contextVarCache[i] = newAnyValue(contextVarsValues[varName])
}

// Execution: O(1) array access instead of O(log n) map lookup
value := vm.contextVarCache[varIndex]  // Fast!
```

**Impact**: Eliminated 14% CPU time spent on map lookups (vs expr)

### 4. Type-Specific Fast Paths

**Pattern**: Handle common types directly without boxing

```go
func (vm *VM) executeComparisonOperationValues(op code.Opcode, left, right Value) error {
    // Fast path: same types (most common)
    if left.Typ == right.Typ {
        switch left.Typ {
        case TypeFloat:
            return vm.executeNumberComparisonOperation(op, left.FloatVal, right.FloatVal)
        case TypeString:
            return vm.executeStringComparisonOperation(op, left.StrVal, right.StrVal)
        case TypeBool:
            return vm.executeBooleanComparisonOperation(op, left.BoolVal, right.BoolVal)
        }
    }
    // Fallback: mixed types
    return vm.executeComparisonOperation(op, left.ToAny(), right.ToAny())
}
```

**Key**: Direct type access avoids interface conversions

## What to Remain Careful About

### ⚠️ Critical Rules for Future Development

#### 1. NEVER Use `Pop()` in Opcode Handlers
```go
// ❌ WRONG - Causes allocation
case code.OpEqual:
    right := vm.Pop()  // Boxes Value → any
    left := vm.Pop()

// ✅ CORRECT - Zero allocation
case code.OpEqual:
    right, left := vm.pop2Values()  // Returns Values directly
```

**Why**: Every `Pop()` call boxes the Value to `any`, causing allocation

#### 2. NEVER Modify Value Struct Field Order
```go
// ❌ WRONG - Breaks optimization, adds padding
type Value struct {
    Typ      valueType  // 1 byte
    FloatVal float64    // 7 bytes padding!
    ...
}

// ✅ CORRECT - Optimized layout
type Value struct {
    AnyVal   any        // Largest first
    StrVal   string     // Then 16-byte types
    FloatVal float64    // Then 8-byte types
    Typ      valueType  // Small types last
    BoolVal  bool
}
```

**Why**: Field order affects struct size due to alignment padding

#### 3. NEVER Box Values Inside VM Operations
```go
// ❌ WRONG - Boxes unnecessarily
func (vm *VM) someOperation(val Value) {
    anyVal := val.ToAny()  // Allocation!
    // ... use anyVal
}

// ✅ CORRECT - Use Value directly
func (vm *VM) someOperation(val Value) {
    switch val.Typ {
    case TypeFloat:
        return vm.handleFloat(val.FloatVal)
    case TypeString:
        return vm.handleString(val.StrVal)
    }
}
```

**Why**: `ToAny()` allocates; direct field access doesn't

#### 4. ALWAYS Use Value Constructors
```go
// ❌ WRONG - Manual struct construction, error-prone
val := Value{
    Typ:      TypeFloat,
    FloatVal: 42.0,
}

// ✅ CORRECT - Use constructor
val := newFloatValue(42.0)
```

**Why**: Constructors ensure correct field initialization

#### 5. ALWAYS Profile Before Optimizing
```bash
# Memory profile
go test -bench=BenchmarkMyOp -memprofile=mem.prof
go tool pprof -alloc_objects mem.prof

# CPU profile
go test -bench=BenchmarkMyOp -cpuprofile=cpu.prof
go tool pprof -top cpu.prof
```

**Why**: Guessing causes regressions; profiling finds real bottlenecks

### ⚠️ Testing Requirements

#### Must-Have Tests After Changes:
```bash
# 1. All unit tests pass
go test ./...

# 2. Benchmark comparison
go test -bench=. -benchmem -benchtime=10s > after.txt
# Compare with baseline

# 3. Allocation check
go test -bench=BenchmarkVM_Pure -benchmem
# Must show 0 allocs/op for primitives

# 4. Integration test with comparison project
cd ../golang-expression-evaluation-comparison
go test -bench=Benchmark_uexl -benchmem
```

### ⚠️ Performance Budgets

**Acceptable Ranges**:
- Main expression benchmark: **180-230 ns/op** ✅
- Allocations for primitives: **0 allocs/op** ✅ (CRITICAL)
- Map operations: **< 5,000 ns/op** ✅
- String operations: **< 300 ns/op** ✅

**Red Flags** (require investigation):
- Any allocation in boolean/comparison operations
- Main benchmark > 250 ns/op
- Any regression > 5%

## Common Pitfalls

### 1. Adding New Opcodes
```go
// ❌ WRONG
case code.OpNewThing:
    val := vm.Pop()  // Allocates!
    result := doSomething(val)
    vm.Push(result)

// ✅ CORRECT
case code.OpNewThing:
    val := vm.popValue()  // Zero-alloc
    result := vm.doSomethingValue(val)  // Value-native
    vm.pushValue(result)
```

### 2. Adding New Helper Functions
```go
// ❌ WRONG - Accepts any, forces boxing
func (vm *VM) newHelper(val any) error {
    // ...
}

// ✅ CORRECT - Accepts Value directly
func (vm *VM) newHelperValue(val Value) error {
    switch val.Typ {
    case TypeFloat:
        // Handle inline
    }
}
```

### 3. Returning Values from VM
```go
// Public API (must box)
func (vm *VM) Run(bytecode, context) (any, error) {
    // Internal operations use Value
    err := vm.run()  // Zero-alloc internally

    // Only box at the very end
    return vm.LastPoppedStackElem(), nil  // Boxes here (unavoidable)
}
```

**Why**: Public API requires `any` for compatibility, but internal ops stay allocation-free

## Optimization Decision Tree

```
Need to optimize something?
│
├─ Does profiling show it's a bottleneck? (>5% CPU or allocs)
│  │
│  NO──► Don't optimize (premature optimization)
│  │
│  YES──► Continue
│
├─ Can it be fixed architecturally? (affects all operations)
│  │
│  YES──► Implement architectural fix
│  │     (Example: Value migration)
│  │
│  NO──► Is it a hot path? (>10% CPU time)
│        │
│        YES──► Optimize specific operation
│        │     (Example: isTruthyValue)
│        │
│        NO──► Accept current performance
```

## Performance Monitoring

### Regression Detection
```bash
# Baseline capture
go test -bench=. -benchmem -benchtime=10s > baseline.txt

# After changes
go test -bench=. -benchmem -benchtime=10s > current.txt

# Compare (install: go install golang.org/x/perf/cmd/benchstat@latest)
benchstat baseline.txt current.txt

# Look for:
# - Any increase in allocs/op
# - Speed regression > 5%
# - Memory increase
```

### Continuous Monitoring
- Run comparison benchmarks before releases
- Track allocations in CI/CD (must be 0 for primitives)
- Monitor struct sizes with `unsafe.Sizeof()`

## Future Optimization Opportunities

### Safe Optimizations (Can Be Done):
1. **Inline directives** on hot functions
   ```go
   //go:inline
   func (vm *VM) popValue() Value { ... }
   ```
   Expected: 10-15ns improvement

2. **Pre-allocated boolean constants**
   ```go
   var trueValue = newBoolValue(true)
   var falseValue = newBoolValue(false)
   ```
   Expected: 5ns improvement

3. **Batch operations** for common patterns
   ```go
   func (vm *VM) pop3Values() (Value, Value, Value)
   ```
   Expected: 5-10ns improvement

### Unsafe Optimizations (DON'T DO):
1. ❌ Union-based Value struct (unsafe pointers)
2. ❌ Assembly dispatch loop
3. ❌ Computed goto dispatch (requires cgo)

**Why avoid**: High complexity, platform-specific, hard to maintain

## Conclusion

### The Formula for Success:
1. ✅ **Architecture first**: Design for zero allocations
2. ✅ **Profile everything**: Measure before optimizing
3. ✅ **Holistic fixes**: Solve root causes, not symptoms
4. ✅ **Test rigorously**: Prevent regressions
5. ✅ **Document deeply**: Enable future developers

### Key Metrics to Maintain:
- **Zero allocations** for primitive operations (CRITICAL)
- Competitive speed (within 2× of fastest competitor)
- Clean, maintainable code
- Comprehensive test coverage

### Remember:
> "The best performance optimization is one that makes the code simpler, not more complex."

UExL achieved zero allocations through architectural clarity, not clever tricks. Keep it that way.
