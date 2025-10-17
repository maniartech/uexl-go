# Optimization Techniques Reference

## Overview

This document provides detailed technical explanations of optimization patterns used in UExL, serving as a reference for applying similar techniques in future work.

---

## Table of Contents

1. [Type System Optimizations](#type-system-optimizations)
2. [Caching Strategies](#caching-strategies)
3. [Map Operation Optimizations](#map-operation-optimizations)
4. [Memory Access Patterns](#memory-access-patterns)
5. [Pointer Semantics](#pointer-semantics)
6. [Inlining and Compiler Hints](#inlining-and-compiler-hints)
7. [Sentinel Values](#sentinel-values)

---

## Type System Optimizations

### Eliminate Redundant Type Assertions

**Pattern:** Type assertion at call site + type assertion in callee = waste

**Before:**
```go
func executeComparisonOperation(op code.Opcode, left, right any) error {
    switch left.(type) {
    case float64:
        return executeNumberComparison(op, left, right)  // Will assert again
    }
}

func executeNumberComparison(op code.Opcode, left, right any) error {
    l := left.(float64)   // Redundant assertion
    r := right.(float64)  // Redundant assertion
    // ... comparison logic
}
```

**After:**
```go
func executeComparisonOperation(op code.Opcode, left, right any) error {
    switch l := left.(type) {  // Type switch + capture
    case float64:
        r, ok := right.(float64)
        if !ok {
            return fmt.Errorf("type mismatch")
        }
        return executeNumberComparison(op, l, r)  // Pass typed values
    }
}

func executeNumberComparison(op code.Opcode, left, right float64) error {
    // No assertions needed - parameters are already typed
    return vm.Push(left == right)
}
```

**Benefits:**
- Eliminates duplicate type assertions (~1-2 ns each)
- Enables better compiler inlining
- Clearer function contracts (types in signature)
- Safer (compile-time type checking)

**When to Apply:**
- Functions called from type-switched dispatchers
- Hot path functions with interface parameters
- Operations on homogeneous types

**Caution:**
- Don't apply to functions that truly need `any` for heterogeneity
- Balance type safety vs API flexibility

---

## Caching Strategies

### Pre-Resolved Value Caching

**Pattern:** Convert expensive lookups into pre-computed arrays

**Problem:**
```go
// Expensive: O(1) but high constant factor
case OpContextVar:
    varName := vm.contextVars[varIndex]
    value := vm.contextVarsValues[varName]  // Map lookup: hash + search
    vm.Push(value)
```

**Solution:**
```go
// Setup phase (once per Run):
func setBaseInstructions(bytecode, contextValues) {
    vm.contextVarCache = make([]any, len(vm.contextVars))
    for i, varName := range vm.contextVars {
        vm.contextVarCache[i] = contextValues[varName]  // Pre-resolve
    }
}

// Execution phase (many times):
case OpContextVar:
    value := vm.contextVarCache[varIndex]  // Array access: direct offset
    vm.Push(value)
}
```

**Cost Analysis:**

| Operation | Setup Cost | Access Cost | Break-Even |
|-----------|------------|-------------|------------|
| Map Lookup | None | ~10-20 ns | 1 access |
| Array Cache | N × 10 ns | ~1-2 ns | N/8 accesses |

**Formula:** Cache is profitable when:
```
AccessCount > SetupCost / (MapCost - ArrayCost)
AccessCount > N / 8  (typical values)
```

**Benefits:**
- Setup: O(N) map lookups once
- Access: O(1) array index, ~10x faster
- Cache-friendly (sequential memory)

**When to Apply:**
- Value set known at start
- Multiple accesses per value
- Values don't change during execution

**Implementation Checklist:**
1. Add cache array to struct
2. Build cache in setup phase
3. Use cache in hot path
4. Handle cache invalidation
5. Consider nil/missing value semantics

---

### Smart Cache Invalidation

**Pattern:** Detect when cache is still valid to skip rebuild

**Naive Approach:**
```go
func Run(bytecode, contextValues) {
    rebuildCache(contextValues)  // Always rebuild
    return execute()
}
```

**Smart Approach:**
```go
type VM struct {
    contextVarCache []any
    lastContextPtr  uintptr  // Track last context map pointer
}

func Run(bytecode, contextValues) {
    newPtr := reflect.ValueOf(contextValues).Pointer()
    if newPtr != vm.lastContextPtr || cacheSizeMismatch() {
        rebuildCache(contextValues)  // Rebuild only when needed
        vm.lastContextPtr = newPtr
    }
    return execute()
}
```

**Trade-off Analysis:**

| Scenario | Rebuild Cost | Check Cost | Net Savings |
|----------|--------------|------------|-------------|
| Same map (hit) | Saved: 40-60 ns | Paid: 2-3 ns | ~38-57 ns |
| Different map (miss) | Paid: 40-60 ns | Paid: 2-3 ns | -2-3 ns (acceptable) |

**Hit Rate Calculation:**
```
Benchmark loops: 100% hit rate (same map)
Production typical: 80-95% hit rate (maps reused in loops)
Break-even: 5% hit rate
```

**Benefits:**
- 95%+ cache hits: Massive savings
- Pointer comparison: ~2-3 ns
- Works for reference types (maps, slices)

**Limitations:**
- Doesn't detect map mutations (same pointer, different values)
- User responsibility: Don't mutate between calls
- Only works for reference types

**When to Apply:**
- Values passed by reference (maps, slices, pointers)
- High cache hit rate expected
- Expensive cache rebuild cost

**Anti-patterns:**
- Value types (structs by value) - pointer always changes
- Rarely reused caches - overhead not worth it
- Mutable caches where pointer check insufficient

---

## Map Operation Optimizations

### Conditional Map Clearing

**Problem:**
```go
// Clears map even when empty
for k := range myMap {
    delete(myMap, k)
}
```

**Cost:**
- Empty map: Still iterates (checks bucket structure)
- Iteration overhead: ~10-20 ns for empty map

**Solution:**
```go
if len(myMap) > 0 {
    clear(myMap)  // Go 1.21+ built-in
}
```

**Benefits:**
- Skip clearing when map empty (common case)
- `clear()` optimized by runtime
- Zero cost for empty maps

**When to Apply:**
- Maps frequently empty
- Maps cleared in hot paths
- Setup/teardown code

---

### Map vs Slice Trade-offs

**Use Map When:**
- Sparse key space
- String keys
- Unknown key set at compile time
- Infrequent access (< 10 per operation)

**Use Slice When:**
- Dense integer keys (0, 1, 2, ...)
- Known key range at compile time
- Frequent access (> 10 per operation)
- Key is index-like

**Conversion Example:**
```go
// Before: Map with integer keys
contextVars := map[int]any{}
value := contextVars[42]  // Hash + lookup

// After: Slice with index access
contextVars := make([]any, maxIndex+1)
value := contextVars[42]  // Direct offset
```

**Performance:**
```
Map access:   ~10-20 ns
Slice access: ~1-2 ns
Speedup:      5-10x
```

---

## Memory Access Patterns

### Sequential vs Random Access

**Cache-Friendly (Sequential):**
```go
// Array iteration - CPU prefetcher loves this
for i := 0; i < len(array); i++ {
    process(array[i])
}
```

**Cache-Hostile (Random):**
```go
// Map iteration - scattered memory
for k, v := range someMap {
    process(v)
}
```

**L1 Cache Performance:**
- Hit: ~1-2 cycles
- Miss: ~10-20 cycles (L2)
- Miss: ~50-100 cycles (L3)
- Miss: ~200+ cycles (RAM)

**Optimization Strategy:**
1. Group related data in structs
2. Use arrays/slices for iteration
3. Access memory sequentially when possible
4. Avoid pointer chasing

---

## Pointer Semantics

### Reference Type Comparison

**Go Reference Types:**
- Maps
- Slices (caveat: header vs backing array)
- Channels
- Pointers

**Pointer Comparison Pattern:**
```go
func sameMap(a, b map[K]V) bool {
    if a == nil && b == nil {
        return true
    }
    if a == nil || b == nil {
        return false
    }
    // Compare pointers via reflect
    return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
}
```

**Use Cases:**
- Cache invalidation
- Change detection
- Avoiding deep comparisons

**Gotchas:**
- Same pointer ≠ same contents (if mutated)
- Slice pointer comparison: Compares slice header, not backing array
- Cost: ~2-3 ns for reflect.ValueOf().Pointer()

---

## Inlining and Compiler Hints

### Inline-Friendly Functions

**Inline Budget:**
Go compiler inlines functions under ~80 "cost units"

**Inline-Friendly:**
```go
func (vm *VM) Push(val any) error {
    vm.stack[vm.sp] = val  // Simple, ~5 cost units
    vm.sp++
    return nil
}
```

**Inline-Hostile:**
```go
func (vm *VM) Push(val any) error {
    if vm.sp >= StackSize {  // Branch: +5 units
        return fmt.Errorf("overflow")  // Call: +20 units
    }
    vm.stack[vm.sp] = val
    vm.sp++
    return nil
}
// Total: ~30+ units - may not inline
```

**Force Inlining:**
```go
//go:inline
func (vm *VM) fastPush(val any) {
    vm.stack[vm.sp] = val
    vm.sp++
}
```

**Check Inlining:**
```bash
go build -gcflags='-m' 2>&1 | grep inline
```

---

## Sentinel Values

### Distinguishing Nil from Missing

**Problem:**
```go
cache := make([]any, n)
cache[0] = nil  // Valid nil value

// Later: Is this "not set" or "set to nil"?
if cache[0] == nil {
    // Ambiguous!
}
```

**Solution:**
```go
type sentinel struct{}
var notSet = sentinel{}

cache := make([]any, n)
cache[0] = nil  // Valid nil
cache[1] = notSet  // Not provided

// Later: Clear distinction
if _, missing := cache[1].(sentinel); missing {
    // Definitely not set
}
if cache[0] == nil {
    // Definitely nil value
}
```

**Benefits:**
- Type-safe sentinel
- Zero memory overhead (empty struct)
- Compile-time checking

**When to Apply:**
- Caches where nil is valid
- Optional values
- Three-state logic (set/unset/nil)

---

## Performance Measurement

### Benchmark Template

```go
func BenchmarkOperation(b *testing.B) {
    // Setup (excluded from timing)
    data := setupTestData()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Operation under test
        result := operation(data)
        
        // Prevent compiler optimization
        if result == nil {
            b.Fatal("unexpected nil")
        }
    }
    b.StopTimer()
    
    // Validation (excluded from timing)
    validateResult(result)
}
```

### CPU Profiling

```bash
# Generate profile
go test -bench=. -cpuprofile=cpu.prof -benchtime=10s

# Analyze top functions
go tool pprof -top cpu.prof

# Analyze cumulative time
go tool pprof -top -cum cpu.prof

# Interactive analysis
go tool pprof cpu.prof
```

### Optimization Workflow

1. **Profile** current code
2. **Identify** top 3 bottlenecks
3. **Hypothesize** optimization
4. **Implement** change
5. **Benchmark** before and after
6. **Profile** again
7. **Validate** correctness
8. **Document** findings

---

## Common Pitfalls

### 1. Premature Optimization

❌ **Don't:**
```go
// Optimizing before profiling
func process() {
    // Complex manual memory pool
    // Custom allocator
    // Assembly optimizations
}
```

✅ **Do:**
```go
// Profile first, find real bottleneck
// $ go test -cpuprofile=cpu.prof
// $ go tool pprof cpu.prof
// (pprof) top
// 70% in unexpectedFunction()

// Optimize what actually matters
```

### 2. Micro-Optimizations Breaking Clarity

❌ **Don't:**
```go
// Saves 2 ns, loses readability
func obscure(x int) int {
    return x&1<<3 | x&^1>>2
}
```

✅ **Do:**
```go
// Clear with comment explaining perf consideration
func transform(x int) int {
    // Bit operations for ~2ns improvement in hot path
    // Equivalent to: if x%2 == 0 { return x*8 } else { return x/4 }
    return x&1<<3 | x&^1>>2
}
```

### 3. Ignoring Allocation Costs

❌ **Don't:**
```go
// Creates new slice every call
func process(items []Item) []Item {
    result := make([]Item, 0)  // Allocation
    for _, item := range items {
        if item.Valid() {
            result = append(result, item)
        }
    }
    return result
}
```

✅ **Do:**
```go
// Reuse slice capacity
func process(items []Item, result []Item) []Item {
    result = result[:0]  // Reset length, keep capacity
    for _, item := range items {
        if item.Valid() {
            result = append(result, item)
        }
    }
    return result
}
```

### 4. Over-Optimizing Cold Paths

Focus on hot paths (> 5% CPU time):
- ✅ Optimize VM execution loop
- ✅ Optimize comparison operations
- ❌ Don't optimize error formatting
- ❌ Don't optimize one-time setup

---

## Optimization Checklist

Before implementing optimization:

- [ ] Profiled and identified bottleneck
- [ ] Bottleneck is in hot path (> 5% time)
- [ ] Estimated performance gain > 3%
- [ ] Considered code complexity increase
- [ ] Benchmarked current performance
- [ ] Designed optimization approach
- [ ] Considered edge cases
- [ ] Planned validation strategy

After implementing optimization:

- [ ] Benchmarked new performance
- [ ] Measured actual gain vs estimate
- [ ] Ran all tests (no regressions)
- [ ] Profiled to verify bottleneck reduced
- [ ] Documented code changes
- [ ] Updated this document if novel technique
- [ ] Committed with before/after metrics

---

## References

### Internal Documents
- [optimization-journey.md](optimization-journey.md) - Historical context
- [profiling-guide.md](profiling-guide.md) - Profiling details
- [best-practices.md](best-practices.md) - Philosophy and guidelines

### External Resources
- [Go Performance Tips](https://github.com/golang/go/wiki/Performance)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Compiler Optimization Decisions](https://github.com/golang/go/blob/master/src/cmd/compile/internal/inline/inl.go)

---

**Last Updated:** October 17, 2025
