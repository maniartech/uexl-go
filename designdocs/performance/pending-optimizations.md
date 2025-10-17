# Pending Optimizations

## Overview

This document tracks optimization opportunities identified but not yet implemented. Each item includes rationale, estimated impact, complexity, and priority.

---

## High Priority (> 10% potential gain)

### P1: Stack Operation Inlining

**Current State:**
```go
func (vm *VM) Push(val any) error {
    if vm.sp >= StackSize {
        return fmt.Errorf("stack overflow")
    }
    vm.stack[vm.sp] = val
    vm.sp++
    return nil
}
```

**Profile Impact:** 7.94% CPU time in `Pop`, 6.21% in `Push`

**Proposed Optimization:**
1. Remove overflow check from release builds
2. Inline Push/Pop manually in hot paths
3. Use build tags for debug/release versions

```go
// +build !debug

func (vm *VM) Push(val any) {
    vm.stack[vm.sp] = val
    vm.sp++
}

// +build debug

func (vm *VM) Push(val any) error {
    if vm.sp >= StackSize {
        return fmt.Errorf("stack overflow at sp=%d", vm.sp)
    }
    vm.stack[vm.sp] = val
    vm.sp++
    return nil
}
```

**Expected Gain:** 5-10% (eliminate error checking overhead)

**Complexity:** Medium
- Build tag management
- Test coverage for both modes
- Error handling changes

**Trade-offs:**
- ❌ Lose runtime overflow detection in release
- ✅ Debug builds still have checks
- ✅ Significant performance gain

**Status:** Not started  
**Assigned:** -  
**Target:** v2.1

---

### P2: Comparison Operation Dispatch Table

**Current State:**
```go
func (vm *VM) executeComparisonOperation(op Opcode, left, right any) error {
    switch l := left.(type) {  // Type switch
    case float64:
        // ...
    case string:
        // ...
    case bool:
        // ...
    }
}
```

**Profile Impact:** 13.37% CPU time

**Proposed Optimization:**
Create dispatch table based on type IDs:

```go
type compareFunc func(Opcode, any, any) error

var compareDispatch [numTypes]compareFunc

func init() {
    compareDispatch[typeFloat64] = compareNumbers
    compareDispatch[typeString] = compareStrings
    compareDispatch[typeBool] = compareBools
}

func (vm *VM) executeComparisonOperation(op Opcode, left, right any) error {
    typeID := getTypeID(left)
    return compareDispatch[typeID](op, left, right)
}
```

**Expected Gain:** 3-5% (eliminate switch overhead)

**Complexity:** High
- Type ID management
- Dispatch table initialization
- Maintaining type safety

**Trade-offs:**
- ❌ More complex code
- ❌ Type ID overhead
- ✅ Faster dispatch

**Status:** Design phase  
**Assigned:** -  
**Target:** v2.2

---

### P3: Bytecode Instruction Batching

**Current State:**
Each opcode processed individually in `run()` loop

**Proposed Optimization:**
Detect patterns and execute batches:

```go
// Pattern: Push multiple constants
case OpConstant:
    // Look ahead for more OpConstant instructions
    count := countConsecutiveConstants(frame.ip)
    if count > 3 {
        vm.pushConstants(frame.ip, count)  // Batch operation
        frame.ip += count * 3
    } else {
        // Single push (existing code)
    }
```

**Expected Gain:** 10-15% for expressions with many constants

**Complexity:** High
- Pattern detection
- Batch operation implementations
- Correctness verification

**Status:** Research phase  
**Assigned:** -  
**Target:** v3.0

---

## Medium Priority (5-10% potential gain)

### P4: Instruction Decoding Optimization

**Current State:**
Sequential instruction reading with bounds checks

**Proposed Optimization:**
Optimize `ReadUint16` for hot paths:

```go
// Current: Safe but slower
func ReadUint16(instructions []byte, offset int) uint16 {
    return binary.BigEndian.Uint16(instructions[offset:offset+2])
}

// Proposed: Direct bit manipulation
//go:inline
func ReadUint16Fast(instructions []byte, offset int) uint16 {
    return uint16(instructions[offset])<<8 | uint16(instructions[offset+1])
}
```

**Expected Gain:** 2-3%

**Complexity:** Low
- Simple implementation
- Need thorough testing

**Status:** Ready to implement  
**Assigned:** -  
**Target:** v2.0

---

### P5: Context Variable Batch Access

**Current State:**
Variables accessed one at a time

**Proposed Optimization:**
Compiler detects patterns and emits batch load:

```go
// Pattern: a + b + c + d
// Current: OpContextVar, OpContextVar, OpAdd, OpContextVar, OpAdd, OpContextVar, OpAdd
// Optimized: OpContextVarBatch(4), OpAddBatch(4)
```

**Expected Gain:** 5-8% for expressions with many variables

**Complexity:** High (requires compiler changes)

**Status:** Design phase  
**Assigned:** -  
**Target:** v2.2

---

### P6: String Comparison Optimization

**Current State:**
String comparison uses Go's built-in comparison

**Profile Impact:** ~5% for string-heavy expressions

**Proposed Optimization:**
SIMD-based string comparison for common cases:

```go
// For equal-length strings under 32 bytes
func fastStringCompare(a, b string) bool {
    // Use SIMD instructions (via assembly or compiler intrinsics)
}
```

**Expected Gain:** 5-7% for string-heavy workloads

**Complexity:** Very High
- Assembly code
- Platform-specific
- Fallback required

**Status:** Not started  
**Assigned:** -  
**Target:** v3.0+

---

## Low Priority (< 5% potential gain)

### P7: Frame Management Optimization

**Current State:**
```go
func (vm *VM) pushFrame(frame *Frame) {
    vm.frames[vm.framesIdx] = frame
    vm.framesIdx++
}
```

**Proposed Optimization:**
Inline frame management in `run()`:

```go
// Eliminate function call overhead
vm.frames[vm.framesIdx] = &Frame{...}
vm.framesIdx++
```

**Expected Gain:** 1-2%

**Complexity:** Low

**Status:** Low priority (minimal gain)

---

### P8: Constant Pool Specialization

**Current State:**
All constants stored as `any`

**Proposed Optimization:**
Separate constant pools by type:

```go
type ByteCode struct {
    FloatConstants  []float64
    StringConstants []string
    BoolConstants   []bool
}
```

**Expected Gain:** 2-4% (eliminate interface boxing)

**Complexity:** High (requires compiler changes)

**Status:** Design phase  
**Target:** v2.3

---

## Research Opportunities

### R1: Register-Based VM

**Current:** Stack-based bytecode VM

**Proposal:** Convert to register-based VM

**Potential Gain:** 20-40%

**Complexity:** Very High (complete VM rewrite)

**References:**
- Lua VM (register-based)
- LuaJIT optimizations

**Status:** Long-term research  
**Target:** v4.0+

---

### R2: Just-In-Time Compilation

**Proposal:** JIT compile hot expressions to native code

**Potential Gain:** 50-100x for hot paths

**Complexity:** Extreme (JIT compiler)

**Considerations:**
- Code generation
- Platform support
- Security implications

**Status:** Research phase  
**Target:** v5.0+

---

### R3: Ahead-of-Time Compilation

**Proposal:** Compile expressions to Go functions

**Example:**
```go
// Expression: a + b * 2
// Compiled to:
func compiled(ctx map[string]any) float64 {
    a := ctx["a"].(float64)
    b := ctx["b"].(float64)
    return a + b * 2
}
```

**Potential Gain:** 80-100% (native code speed)

**Complexity:** High

**Status:** Proof of concept  
**Target:** v3.0

---

## Completed Optimizations

### ✅ C1: Type-Specific Comparison Functions

**Completed:** October 17, 2025  
**Actual Gain:** 3%  
**See:** [optimization-journey.md](optimization-journey.md#optimization-phase-1-type-system)

---

### ✅ C2: Context Variable Caching

**Completed:** October 17, 2025  
**Actual Gain:** 10%  
**See:** [optimization-journey.md](optimization-journey.md#optimization-phase-2-context-variable-caching)

---

### ✅ C3: Map Operation Optimization

**Completed:** October 17, 2025  
**Actual Gain:** 4%  
**See:** [optimization-journey.md](optimization-journey.md#optimization-phase-3-map-operations)

---

### ✅ C4: Smart Cache Invalidation

**Completed:** October 17, 2025  
**Actual Gain:** 31%  
**See:** [optimization-journey.md](optimization-journey.md#optimization-phase-4-smart-cache-invalidation)

---

## Optimization Selection Criteria

When choosing next optimization to implement:

**Impact Score = (Expected Gain × Priority Factor) / Complexity**

Where:
- Expected Gain: % performance improvement
- Priority Factor:
  - Hot path bottleneck: 3.0
  - Common use case: 2.0
  - Edge case: 0.5
- Complexity:
  - Low: 1.0
  - Medium: 2.0
  - High: 4.0
  - Very High: 8.0

**Example:**
```
P4: Instruction Decoding
  Expected Gain: 2-3% (average 2.5%)
  Priority: Hot path (3.0)
  Complexity: Low (1.0)
  
  Score = (2.5 × 3.0) / 1.0 = 7.5

P2: Comparison Dispatch
  Expected Gain: 3-5% (average 4%)
  Priority: Hot path (3.0)
  Complexity: High (4.0)
  
  Score = (4 × 3.0) / 4.0 = 3.0

Recommendation: P4 first (higher score)
```

---

## Other Areas to Explore

### Compiler Optimizations

1. **Constant Folding:** `2 + 3` → `5` at compile time
2. **Dead Code Elimination:** Remove unreachable branches
3. **Common Subexpression Elimination:** Reuse computed values
4. **Loop Unrolling:** For known-size pipe operations

### VM Optimizations

1. **Threaded Code:** Eliminate dispatch loop overhead (Go limitation)
2. **Superinstructions:** Combined opcodes for common patterns
3. **Inline Caching:** Cache type information for polymorphic sites
4. **Stack Allocation:** Replace stack slice with fixed array

### Memory Optimizations

1. **Object Pooling:** Reuse frame allocations
2. **String Interning:** Deduplicate string constants
3. **Compact Bytecode:** Reduce instruction size
4. **Memory Mapping:** Load bytecode directly without parsing

---

## Contributing

To propose a new optimization:

1. Profile and identify bottleneck (> 5% CPU time)
2. Design solution with trade-off analysis
3. Create proof of concept
4. Benchmark actual vs estimated gain
5. Add to this document with status "Proposed"
6. Discuss in team meeting
7. Get approval and assign

---

## References

- [optimization-journey.md](optimization-journey.md) - Historical optimizations
- [optimization-techniques.md](optimization-techniques.md) - Technique catalog
- [best-practices.md](best-practices.md) - Guidelines

**Last Updated:** October 17, 2025
