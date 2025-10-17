# Performance Best Practices

## Philosophy

UExL's performance philosophy rests on three foundational principles:

### 1. Zero-Cost Abstractions
> "What you don't use, you don't pay for. What you do use, you couldn't hand code any better."

- Abstractions compile away (type system)
- Runtime costs only for actual work
- No hidden allocations
- No surprise overhead

### 2. Data-Oriented Design
> "Memory access patterns matter more than algorithm complexity in modern CPUs."

- Optimize for cache locality
- Sequential > Random access
- Arrays > Maps for hot paths
- Minimize pointer chasing

### 3. Profile-Driven Optimization
> "In God we trust, all others bring data."

- Never optimize without profiling
- Measure before and after
- Real workloads > Microbenchmarks
- Profile guides priorities

---

## Core Guidelines

### G1: Profile Before Optimizing

**Always profile to find bottlenecks:**

```bash
# Generate CPU profile
go test -bench=BenchmarkName -cpuprofile=cpu.prof -benchtime=10s

# Analyze top functions
go tool pprof -top -cum cpu.prof
```

**Focus on hot paths (> 5% CPU time):**
- ✅ Optimize functions consuming significant time
- ❌ Ignore one-time setup or error paths
- ❌ Don't optimize based on assumptions

**Document profile findings:**
```
Before optimization:
  functionX: 45% CPU time
  functionY: 20% CPU time
  
Target: functionX (high impact)
```

---

### G2: Maintain Zero Allocations

**UExL's zero-allocation guarantee:**

Every benchmark should show:
```
BenchmarkOperation-16    N ns/op    0 B/op    0 allocs/op
                                     ^^        ^^^^^^^^^^^^
                                     Must be zero
```

**Common allocation sources:**
- Interface boxing (avoid `any` in hot paths)
- Map operations (preallocate)
- Slice growth (preallocate capacity)
- Error creation (reuse error values)
- String concatenation (use strings.Builder)

**Check allocations:**
```bash
go test -bench=. -benchmem
```

---

### G3: Preserve Type Safety

**Never sacrifice safety for speed:**

❌ **Wrong:**
```go
// Unsafe pointer arithmetic for "speed"
func unsafePeek(stack []any, sp int) any {
    ptr := unsafe.Pointer(&stack[0])
    return *(*any)(unsafe.Pointer(uintptr(ptr) + uintptr(sp)))
}
```

✅ **Right:**
```go
// Simple, safe, and fast enough
func peek(stack []any, sp int) any {
    if sp > 0 {
        return stack[sp-1]
    }
    return nil
}
```

**Type safety benefits:**
- Compile-time error detection
- Clear contracts
- Maintainable code
- No undefined behavior

---

### G4: Optimize Hot Paths Only

**Hot path definition:**
- Executed many times per operation
- Consumes > 5% CPU time
- In critical execution loop

**Priority order:**
1. **Critical** (> 20% CPU): Must optimize
2. **Important** (5-20% CPU): Should optimize
3. **Minor** (1-5% CPU): Consider if easy win
4. **Negligible** (< 1% CPU): Don't optimize

**Example from UExL:**
```
run() loop:              77% CPU  → Critical
setBaseInstructions():   18% CPU  → Important  
error formatting:        0.1% CPU → Negligible (skip)
```

---

### G5: Amortize Expensive Operations

**Setup costs are acceptable if amortized:**

```go
// Expensive setup (once)
func setBaseInstructions(...) {
    // O(N) map lookups to build cache
    for i, varName := range vm.contextVars {
        vm.contextVarCache[i] = contextValues[varName]
    }
}

// Cheap access (many times)
case OpContextVar:
    value := vm.contextVarCache[varIndex]  // O(1) array access
```

**Amortization calculation:**
```
Setup cost: N × 10 ns  (e.g., 4 variables = 40 ns)
Access cost saved: 8 ns per access (10 ns map - 2 ns array)
Break-even: 40 ns / 8 ns = 5 accesses

Typical expression: 4 variables × 2 accesses each = 8 accesses
Net savings: 8 × 8 ns - 40 ns = 24 ns
```

---

### G6: Understand Trade-offs

Every optimization has trade-offs. Document them:

**Example: Smart Cache Invalidation**

**Gains:**
- 30% performance improvement
- ~40 ns saved per cache hit
- 95%+ hit rate in benchmarks

**Costs:**
- 2-3 ns pointer comparison overhead
- Additional struct field (8 bytes)
- User must not mutate maps between calls

**Decision:**
- ✅ Accept: Gains far outweigh costs
- ✅ Document: User responsibility clear
- ✅ Verify: Test coverage for edge cases

---

### G7: Write Inline-Friendly Code

**Keep hot path functions small and simple:**

```go
// Good: Simple, will inline
func (vm *VM) Push(val any) error {
    vm.stack[vm.sp] = val
    vm.sp++
    return nil
}

// Bad: Too complex, won't inline
func (vm *VM) PushWithLogging(val any) error {
    if vm.debug {  // Branch
        log.Printf("Pushing %v", val)  // Call
    }
    if vm.sp >= StackSize {  // Branch
        return fmt.Errorf("overflow at sp=%d", vm.sp)  // Allocation
    }
    vm.stack[vm.sp] = val
    vm.sp++
    return nil
}
```

**Inlining budget:** ~80 cost units
- Assignment: ~1 unit
- Arithmetic: ~1 unit
- Branch: ~5 units
- Function call: ~20+ units
- Allocation: ~50+ units

**Check inlining:**
```bash
go build -gcflags='-m -m' 2>&1 | grep 'function too complex'
```

---

### G8: Prefer Stack to Heap

**Stack allocation is faster:**
- No GC pressure
- Better cache locality
- Automatic cleanup

**Force stack allocation:**
```go
// Heap allocation (escape analysis)
func makeSlice() []int {
    s := make([]int, 100)
    return s  // Escapes to heap
}

// Stack allocation
func makeSlice(buf []int) []int {
    if cap(buf) >= 100 {
        return buf[:100]  // Reuse stack buffer
    }
    return make([]int, 100)  // Fallback
}
```

**Check escape analysis:**
```bash
go build -gcflags='-m' 2>&1 | grep escape
```

---

### G9: Batch Operations

**Reduce overhead by batching:**

❌ **Inefficient:**
```go
for i := 0; i < 10; i++ {
    vm.Push(values[i])  // Function call overhead × 10
}
```

✅ **Better:**
```go
// Batch push operation
func (vm *VM) PushN(values []any) error {
    if vm.sp + len(values) > StackSize {
        return fmt.Errorf("stack overflow")
    }
    copy(vm.stack[vm.sp:], values)  // Single operation
    vm.sp += len(values)
    return nil
}
```

---

### G10: Document Performance Decisions

**Every non-obvious performance decision needs comments:**

```go
// Pre-resolve context variables into a lookup slice for O(1) access.
// This eliminates expensive map lookups in the hot path.
// Trade-off: O(N) setup cost amortized over multiple OpContextVar accesses.
// Benchmark: 4 variables × 2 accesses = 8 accesses > 5 break-even point.
if len(vm.contextVars) > 0 {
    vm.contextVarCache = make([]any, len(vm.contextVars))
    for i, varName := range vm.contextVars {
        vm.contextVarCache[i] = contextValues[varName]
    }
}
```

**Good comments explain:**
- Why (rationale)
- Trade-offs (costs vs benefits)
- Measurements (benchmark data)
- Assumptions (when safe to use)

---

## Design Patterns

### Pattern 1: Cache with Invalidation

**When to use:**
- Expensive computation/lookup
- Results reused multiple times
- Clear invalidation point

**Implementation:**
```go
type Cache struct {
    data []any
    lastKey uintptr  // For pointer-based invalidation
}

func (c *Cache) Get(key map[string]any) []any {
    newKey := reflect.ValueOf(key).Pointer()
    if newKey != c.lastKey {
        c.rebuild(key)
        c.lastKey = newKey
    }
    return c.data
}
```

---

### Pattern 2: Type-Specific Dispatch

**When to use:**
- Operations on different types
- Type known at runtime
- Hot path

**Implementation:**
```go
func dispatch(op Opcode, left, right any) error {
    switch l := left.(type) {
    case float64:
        r := right.(float64)  // Assert once
        return numberOp(op, l, r)  // Typed function
    case string:
        r := right.(string)
        return stringOp(op, l, r)
    }
}
```

---

### Pattern 3: Sentinel Values

**When to use:**
- Need to distinguish nil from missing
- Three-state logic required

**Implementation:**
```go
type sentinel struct{}
var notProvided = sentinel{}

cache[i] = value              // Actual value (can be nil)
cache[j] = notProvided        // Not set

if _, missing := cache[i].(sentinel); missing {
    // Handle not-set case
}
```

---

### Pattern 4: Pre-allocation with Reuse

**When to use:**
- Known maximum size
- Frequent allocations
- GC pressure

**Implementation:**
```go
type VM struct {
    stack []any  // Pre-allocated to max size
    sp    int    // Current position
}

func New() *VM {
    return &VM{
        stack: make([]any, MaxStackSize),  // Allocate once
        sp:    0,
    }
}

func (vm *VM) Reset() {
    vm.sp = 0  // Reuse allocation
    for i := 0; i < vm.sp; i++ {
        vm.stack[i] = nil  // Clear references for GC
    }
}
```

---

## Anti-Patterns

### Anti-Pattern 1: Premature Optimization

❌ **Wrong:**
```go
// Optimizing before profiling
func (vm *VM) complexOptimizedPush(val any) error {
    // 50 lines of micro-optimizations
    // Saves 0.5 ns but unreadable
}
```

✅ **Right:**
```go
// Simple, measure, then optimize if needed
func (vm *VM) Push(val any) error {
    vm.stack[vm.sp] = val
    vm.sp++
    return nil
}
// Profile shows: 0.1% CPU time → Don't optimize
```

---

### Anti-Pattern 2: Micro-Optimizing Cold Paths

❌ **Wrong:**
```go
// Optimizing error paths (called rarely)
func (vm *VM) formatError(err error) string {
    // Complex string pooling
    // Manual memory management
    // Saves allocations in error path
}
```

✅ **Right:**
```go
// Simple error formatting is fine
func (vm *VM) formatError(err error) string {
    return fmt.Sprintf("vm error: %v", err)
}
// Profile shows: 0.01% CPU time → Simple is fine
```

---

### Anti-Pattern 3: Over-Engineering

❌ **Wrong:**
```go
// Custom memory allocator for everything
type CustomAllocator struct {
    pools [100]*MemoryPool
    sizes []int
    // ... 500 lines of complexity
}
```

✅ **Right:**
```go
// Let Go's allocator do its job
// Optimize only proven bottlenecks
stack := make([]any, StackSize)  // Simple, works well
```

---

### Anti-Pattern 4: Sacrificing Clarity

❌ **Wrong:**
```go
// Unreadable bit hacks
func hash(s string) uint32 {
    h := uint32(0x1505)
    for _, c := range s {
        h = ((h<<5)+h)^uint32(c)
    }
    return h&0x7fffffff>>3^h<<9
}
```

✅ **Right:**
```go
// Clear, with comment if complex
func hash(s string) uint32 {
    // FNV-1a hash algorithm
    h := uint32(0x811c9dc5)
    for _, c := range s {
        h ^= uint32(c)
        h *= 0x01000193
    }
    return h
}
```

---

## Testing and Validation

### Performance Tests

**Every optimization needs benchmarks:**

```go
func BenchmarkBefore(b *testing.B) {
    // Old implementation
    for i := 0; i < b.N; i++ {
        oldOperation()
    }
}

func BenchmarkAfter(b *testing.B) {
    // New implementation
    for i := 0; i < b.N; i++ {
        newOperation()
    }
}
```

**Run comparisons:**
```bash
go test -bench=. -benchmem -count=10 > new.txt
git checkout main
go test -bench=. -benchmem -count=10 > old.txt
benchstat old.txt new.txt
```

---

### Correctness Tests

**Performance must not break correctness:**

```go
func TestOptimizedVsOriginal(t *testing.T) {
    testCases := []struct{
        input    string
        expected any
    }{
        {"1 + 1", 2.0},
        // ... comprehensive cases
    }
    
    for _, tc := range testCases {
        resultOpt := optimizedPath(tc.input)
        resultOrig := originalPath(tc.input)
        
        if !reflect.DeepEqual(resultOpt, resultOrig) {
            t.Errorf("Results differ: %v vs %v", resultOpt, resultOrig)
        }
    }
}
```

---

### Regression Testing

**Track performance over time:**

```go
// performance_test.go
func TestPerformanceRegression(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping performance test")
    }
    
    result := testing.Benchmark(BenchmarkCriticalPath)
    nsPerOp := result.NsPerOp()
    
    const threshold = 70  // ns/op
    if nsPerOp > threshold {
        t.Errorf("Performance regression: %v ns/op > %v threshold", 
                 nsPerOp, threshold)
    }
}
```

---

## Code Review Checklist

When reviewing performance changes:

- [ ] Profiling data included (before/after)
- [ ] Benchmark results show improvement
- [ ] No allocations added (`-benchmem`)
- [ ] All tests pass
- [ ] Code clarity maintained
- [ ] Trade-offs documented
- [ ] Hot path focused (> 5% CPU)
- [ ] No premature optimization
- [ ] Edge cases tested
- [ ] Error handling preserved

---

## Common Mistakes

### 1. Optimizing Without Measuring

**Mistake:**
```go
// Assumption: map is slow
// Reality: map is fine for 10 items
```

**Solution:** Profile first, optimize second

---

### 2. Ignoring Big-O

**Mistake:**
```go
// Optimizing constant factors in O(n²) algorithm
// When O(n log n) algorithm exists
```

**Solution:** Algorithm > micro-optimization

---

### 3. Breaking Encapsulation

**Mistake:**
```go
// Exposing internal fields for "performance"
type VM struct {
    Stack []any  // Exported!
    SP    int    // Exported!
}
```

**Solution:** Keep abstractions, optimize inside

---

### 4. Assuming GC is Free

**Mistake:**
```go
// Creating garbage in hot path
for i := 0; i < n; i++ {
    temp := make([]int, 100)  // Allocation
    // ... use temp
}
```

**Solution:** Reuse allocations, minimize GC pressure

---

## When to Stop Optimizing

**Diminishing returns:** Stop when:

1. **Performance target met:** Hit SLA/requirement
2. **Cost > Benefit:** 1% gain for 100% complexity
3. **Hot paths exhausted:** Only cold paths remain
4. **External bottlenecks:** Network/disk dominant
5. **Maintainability risk:** Code becoming unmaintainable

**Example:**
```
Target: < 100 ns/op
Current: 62 ns/op
Status: ✅ Target exceeded by 38%

Action: Stop optimizing, focus elsewhere
```

---

## Future-Proofing

### Write Benchmark-Friendly Code

```go
// Easy to benchmark
func Process(input Data) Result {
    return compute(input)
}

// Hard to benchmark
func ProcessWithSideEffects() {
    global := readGlobal()
    result := compute(global)
    writeGlobal(result)
}
```

---

### Leave Optimization Hooks

```go
type VM struct {
    // Current: map lookup
    contextVarsValues map[string]any
    
    // Future: could switch to slice
    // contextVarsArray []any
    
    // Interface allows swapping
    contextProvider ContextProvider
}
```

---

### Document Assumptions

```go
// ASSUMPTION: Context maps are not mutated between Run() calls.
// If this assumption breaks, remove pointer-based cache invalidation.
// Fallback: Always rebuild cache (40 ns cost).
if newPtr != vm.lastContextPtr {
    vm.rebuildCache()
}
```

---

## Summary

**Golden Rules:**

1. **Profile before optimizing**
2. **Maintain zero allocations**
3. **Preserve type safety**
4. **Optimize hot paths only**
5. **Amortize expensive operations**
6. **Understand trade-offs**
7. **Write inline-friendly code**
8. **Prefer stack to heap**
9. **Batch operations**
10. **Document decisions**

**Remember:**
> "Premature optimization is the root of all evil, but profiler-guided optimization is the path to enlightenment."
> — Adapted from Donald Knuth

---

**Last Updated:** October 17, 2025
