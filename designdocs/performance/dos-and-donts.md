# Performance Dos and Don'ts

## Quick Reference

A fast-lookup guide for making performance-conscious decisions during UExL development.

---

## General Performance

### ✅ DO

- **Profile before optimizing** - Always measure, never guess
- **Use benchmarks to validate** - Prove improvements with numbers
- **Optimize hot paths only** - Focus on code that runs frequently
- **Keep benchmarks realistic** - Mirror production usage patterns
- **Document trade-offs** - Explain why optimizations were made
- **Test correctness first** - Ensure optimizations don't break functionality
- **Use build tags for debug code** - Keep production builds fast
- **Benchmark comparisons statistically** - Use `benchstat` with 10+ runs

### ❌ DON'T

- **Optimize without profiling** - You'll often optimize the wrong thing
- **Micro-optimize cold paths** - Little impact, high complexity cost
- **Sacrifice correctness for speed** - Bugs cost more than milliseconds
- **Skip testing after optimization** - Regressions hide easily
- **Make code unreadable** - Maintainability matters
- **Assume allocations are free** - Heap allocations have real cost
- **Trust single benchmark runs** - Statistical noise can mislead
- **Optimize prematurely** - Get it working first, fast second

---

## Type System

### ✅ DO

```go
// ✅ Type-specific function signatures
func compareNumbers(op Opcode, left, right float64) bool {
    switch op {
    case OpEqual:
        return left == right
    // ...
    }
}

// ✅ Single type check, pass typed values
switch l := left.(type) {
case float64:
    if r, ok := right.(float64); ok {
        return compareNumbers(op, l, r)  // Typed
    }
}
```

```go
// ✅ Avoid repeated type assertions
func process(val any) {
    if f, ok := val.(float64); ok {  // Check once
        result := f * 2
        // Use f multiple times
        other := f + 10
    }
}
```

### ❌ DON'T

```go
// ❌ Accept any when you know the type
func compareNumbers(op Opcode, left, right any) bool {
    l := left.(float64)   // Redundant assertion
    r := right.(float64)
    // ...
}

// ❌ Repeated type assertions
func process(val any) {
    result := val.(float64) * 2      // First assertion
    other := val.(float64) + 10      // Second assertion (wasteful)
}
```

---

## Caching

### ✅ DO

```go
// ✅ Pre-compute frequently accessed data
func (vm *VM) setBaseInstructions(bc *ByteCode, ctx map[string]any) {
    // Build cache once
    vm.contextVarCache = make([]any, len(vm.contextVarNames))
    for i, name := range vm.contextVarNames {
        if val, ok := ctx[name]; ok {
            vm.contextVarCache[i] = val
        } else {
            vm.contextVarCache[i] = contextVarNotProvided  // Sentinel
        }
    }
}

// ✅ Use array access for O(1) lookup
func (vm *VM) getContextVar(index int) any {
    return vm.contextVarCache[index]  // Fast
}
```

```go
// ✅ Invalidate cache smartly (pointer comparison)
func (vm *VM) setBaseInstructions(bc *ByteCode, ctx map[string]any) {
    newPtr := reflect.ValueOf(ctx).Pointer()
    if newPtr == vm.lastContextPtr {
        return  // Cache still valid, skip rebuild
    }
    vm.lastContextPtr = newPtr
    // Rebuild cache...
}
```

### ❌ DON'T

```go
// ❌ Repeated map lookups
func (vm *VM) getContextVar(name string) any {
    return vm.context[name]  // Map lookup every time (slow)
}

// ❌ Always rebuild cache
func (vm *VM) setBaseInstructions(bc *ByteCode, ctx map[string]any) {
    // Rebuild cache even if ctx hasn't changed
    vm.contextVarCache = make([]any, len(vm.contextVarNames))
    // ...
}

// ❌ Deep equality checks for cache invalidation
func (vm *VM) needsCacheRebuild(ctx map[string]any) bool {
    return !reflect.DeepEqual(ctx, vm.lastContext)  // Expensive!
}
```

---

## Map Operations

### ✅ DO

```go
// ✅ Use built-in clear() for Go 1.21+
func (vm *VM) reset() {
    if len(vm.localVars) > 0 {
        clear(vm.localVars)  // Fast
    }
}

// ✅ Pre-allocate maps with known size
vm.contextVars = make(map[string]any, len(expectedVars))

// ✅ Reuse maps instead of recreating
var contextPool = sync.Pool{
    New: func() any {
        return make(map[string]any, 10)
    },
}

func getContext() map[string]any {
    ctx := contextPool.Get().(map[string]any)
    clear(ctx)  // Reuse
    return ctx
}
```

### ❌ DON'T

```go
// ❌ Iterate to clear (slow for large maps)
func (vm *VM) reset() {
    for k := range vm.localVars {
        delete(vm.localVars, k)  // Slow
    }
}

// ❌ Recreate map every time
func (vm *VM) reset() {
    vm.localVars = make(map[string]any)  // Allocation
}

// ❌ Use map when array would work
var contextVars map[int]any  // Map with int keys = waste

// Use slice instead
var contextVars []any
```

---

## Memory Management

### ✅ DO

```go
// ✅ Use sentinel values to avoid allocations
type contextVarMissing struct{}
var contextVarNotProvided = contextVarMissing{}  // Singleton

func (vm *VM) getVar(idx int) any {
    val := vm.cache[idx]
    if val == contextVarNotProvided {
        return nil
    }
    return val
}

// ✅ Reuse slices
func (vm *VM) reset() {
    vm.stack = vm.stack[:0]  // Reuse backing array
}

// ✅ Pre-allocate with capacity
vm.stack = make([]any, 0, 1024)  // Avoid growth
```

```go
// ✅ Use object pooling for frequent allocations
var framePool = sync.Pool{
    New: func() any {
        return &Frame{}
    },
}

func newFrame() *Frame {
    return framePool.Get().(*Frame)
}

func releaseFrame(f *Frame) {
    framePool.Put(f)
}
```

### ❌ DON'T

```go
// ❌ Allocate nil pointers to distinguish missing values
func (vm *VM) getVar(idx int) any {
    val := vm.cache[idx]
    if val == nil {
        return new(int)  // Allocation for missing value!
    }
    return val
}

// ❌ Create new slices when reusing works
func (vm *VM) reset() {
    vm.stack = make([]any, 0)  // Allocation
}

// ❌ Let slices grow incrementally
vm.stack = make([]any, 0)  // Will reallocate multiple times as it grows
```

---

## Stack Operations

### ✅ DO

```go
// ✅ Inline-friendly functions (for release builds)
// +build !debug

func (vm *VM) Push(val any) {
    vm.stack[vm.sp] = val
    vm.sp++
}

func (vm *VM) Pop() any {
    vm.sp--
    return vm.stack[vm.sp]
}

// ✅ Use build tags for safety checks
// +build debug

func (vm *VM) Push(val any) error {
    if vm.sp >= StackSize {
        return fmt.Errorf("stack overflow")
    }
    vm.stack[vm.sp] = val
    vm.sp++
    return nil
}
```

### ❌ DON'T

```go
// ❌ Error checks in hot path (production)
func (vm *VM) Push(val any) error {
    if vm.sp >= StackSize {  // Check every push
        return fmt.Errorf("stack overflow")
    }
    vm.stack[vm.sp] = val
    vm.sp++
    return nil
}

// ❌ Function calls that prevent inlining
func (vm *VM) Push(val any) {
    vm.validateStackSpace()  // Extra call
    vm.stack[vm.sp] = val
    vm.sp++
}
```

---

## Pointer Semantics

### ✅ DO

```go
// ✅ Use pointer comparison for identity checks
func (vm *VM) cacheValid(ctx map[string]any) bool {
    newPtr := reflect.ValueOf(ctx).Pointer()
    return newPtr == vm.lastPtr  // Fast pointer comparison
}

// ✅ Use pointers for large structs
func processFrame(f *Frame) {  // Pass by reference
    // Work with f
}

// ✅ Return pointers from constructors
func NewVM() *VM {
    return &VM{
        stack: make([]any, 0, 1024),
    }
}
```

### ❌ DON'T

```go
// ❌ Deep equality for identity
func (vm *VM) cacheValid(ctx map[string]any) bool {
    return reflect.DeepEqual(ctx, vm.lastCtx)  // Slow!
}

// ❌ Pass large structs by value
func processFrame(f Frame) {  // Copies entire struct!
    // Work with f
}

// ❌ Unnecessary pointer-to-pointer
func (vm *VM) getFrame() **Frame {  // Extra indirection
    return &vm.currentFrame
}
```

---

## String Operations

### ✅ DO

```go
// ✅ Use strings.Builder for concatenation
func buildExpr(parts []string) string {
    var b strings.Builder
    for _, p := range parts {
        b.WriteString(p)
    }
    return b.String()
}

// ✅ Direct string comparison (Go optimizes this)
if str1 == str2 {
    // Fast
}

// ✅ Reuse byte slices
var buf []byte
func appendData(data string) {
    buf = append(buf[:0], data...)  // Reuse backing array
}
```

### ❌ DON'T

```go
// ❌ String concatenation in loop
func buildExpr(parts []string) string {
    result := ""
    for _, p := range parts {
        result += p  // New allocation each iteration!
    }
    return result
}

// ❌ Unnecessary conversions
func process(data string) {
    bytes := []byte(data)     // Allocation
    str := string(bytes)      // Another allocation
    // ...
}
```

---

## Benchmarking

### ✅ DO

```go
// ✅ Setup outside benchmark loop
func BenchmarkVM_Run(b *testing.B) {
    // Parse & compile once
    node, _ := parser.ParseString("a && b")
    comp := compiler.New()
    comp.Compile(node)
    bytecode := comp.ByteCode()
    
    machine := vm.New(vm.LibContext{})
    ctx := map[string]any{"a": true, "b": false}
    
    // Reset timer before measurement
    b.ResetTimer()
    
    // Only measure execution
    for i := 0; i < b.N; i++ {
        machine.Run(bytecode, ctx)
    }
}

// ✅ Prevent compiler optimizations
var result any

func BenchmarkVM_Run(b *testing.B) {
    // ... setup ...
    var r any
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        r, _ = machine.Run(bytecode, ctx)
    }
    result = r  // Force computation
}
```

```go
// ✅ Run multiple iterations for statistics
// go test -bench=. -count=10 | benchstat
```

### ❌ DON'T

```go
// ❌ Setup inside benchmark loop
func BenchmarkVM_Run(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Parse & compile every iteration! (slow & unrealistic)
        node, _ := parser.ParseString("a && b")
        comp := compiler.New()
        comp.Compile(node)
        bytecode := comp.ByteCode()
        
        machine.Run(bytecode, ctx)
    }
}

// ❌ Unused results (may be optimized away)
func BenchmarkVM_Run(b *testing.B) {
    for i := 0; i < b.N; i++ {
        machine.Run(bytecode, ctx)  // Result unused - might be skipped!
    }
}

// ❌ Single benchmark run
// go test -bench=.  // Unreliable due to noise
```

---

## Profiling

### ✅ DO

```bash
# ✅ Profile with sufficient runtime
go test -bench=. -benchtime=10s -cpuprofile=cpu.prof

# ✅ Analyze with pprof
go tool pprof -http=:8080 cpu.prof

# ✅ Compare before/after
go tool pprof -base=before.prof after.prof

# ✅ Look at flame graphs for visual understanding
```

```go
// ✅ Add targeted benchmarks for profiling
func BenchmarkVM_ContextAccess(b *testing.B) {
    // Isolate specific operation for profiling
}
```

### ❌ DON'T

```bash
# ❌ Too short runtime (noisy results)
go test -bench=. -benchtime=100ms -cpuprofile=cpu.prof

# ❌ Profile without benchmark context
# Hard to correlate profile data to actual code paths

# ❌ Ignore flat vs cumulative time
# Misidentifies optimization targets

# ❌ Optimize based on gut feeling
# Profile first, then optimize
```

---

## Code Organization

### ✅ DO

```go
// ✅ Hot path in separate, focused function
func (vm *VM) executeComparison(op Opcode, left, right any) error {
    // Focused logic, easier to optimize
}

// ✅ Use //go:inline directive hints
//go:inline
func (vm *VM) peek() any {
    return vm.stack[vm.sp-1]
}

// ✅ Keep critical paths simple
func (vm *VM) run() error {
    for vm.ip < len(vm.instructions) {
        // Simple, inlineable operations
        op := vm.instructions[vm.ip]
        vm.executeOp(op)
    }
}
```

### ❌ DON'T

```go
// ❌ Complex hot path with branching
func (vm *VM) run() error {
    for vm.ip < len(vm.instructions) {
        op := vm.instructions[vm.ip]
        
        // Too much logic in hot path
        if vm.debugMode {
            vm.logOp(op)
        }
        if vm.shouldProfile {
            vm.profileStart()
        }
        
        vm.executeOp(op)
        
        if vm.shouldProfile {
            vm.profileEnd()
        }
        if vm.checkBreakpoints {
            vm.handleBreakpoint()
        }
    }
}

// ❌ Functions too large to inline
func (vm *VM) doEverything() {
    // 100+ lines of code
    // Compiler won't inline
}
```

---

## Error Handling

### ✅ DO

```go
// ✅ Return errors, don't panic
func (vm *VM) Run(bc *ByteCode, ctx map[string]any) (any, error) {
    if bc == nil {
        return nil, fmt.Errorf("bytecode is nil")
    }
    // ...
}

// ✅ Use build tags for expensive checks
// +build debug

func (vm *VM) validateState() error {
    if vm.sp < 0 || vm.sp >= StackSize {
        return fmt.Errorf("invalid sp: %d", vm.sp)
    }
    return nil
}
```

### ❌ DON'T

```go
// ❌ Panic in production code
func (vm *VM) Run(bc *ByteCode, ctx map[string]any) (any, error) {
    if bc == nil {
        panic("bytecode is nil")  // Crashes production!
    }
    // ...
}

// ❌ Expensive validation in hot path
func (vm *VM) Push(val any) error {
    // Validate entire VM state on every push
    if err := vm.validateState(); err != nil {  // Slow!
        return err
    }
    // ...
}
```

---

## Comparison Operations

### ✅ DO

```go
// ✅ Type switch once, dispatch to typed functions
func (vm *VM) compare(op Opcode, left, right any) (bool, error) {
    switch l := left.(type) {
    case float64:
        if r, ok := right.(float64); ok {
            return compareNumbers(op, l, r), nil
        }
    case string:
        if r, ok := right.(string); ok {
            return compareStrings(op, l, r), nil
        }
    }
    return false, fmt.Errorf("type mismatch")
}

func compareNumbers(op Opcode, left, right float64) bool {
    switch op {
    case OpEqual: return left == right
    case OpGreaterThan: return left > right
    // ...
    }
}
```

### ❌ DON'T

```go
// ❌ Type assertions inside comparison
func (vm *VM) compare(op Opcode, left, right any) (bool, error) {
    switch op {
    case OpEqual:
        // Type check for every operator
        if l, ok := left.(float64); ok {
            if r, ok := right.(float64); ok {
                return l == r, nil
            }
        }
    case OpGreaterThan:
        // Repeat type checks
        if l, ok := left.(float64); ok {
            if r, ok := right.(float64); ok {
                return l > r, nil
            }
        }
    }
}
```

---

## Testing Performance Changes

### ✅ DO

```bash
# ✅ Test correctness first
go test ./...

# ✅ Establish baseline
go test -bench=. -count=10 > baseline.txt

# Make changes...

# ✅ Compare statistically
go test -bench=. -count=10 > optimized.txt
benchstat baseline.txt optimized.txt

# ✅ Profile before and after
go test -bench=. -cpuprofile=before.prof
# ... make changes ...
go test -bench=. -cpuprofile=after.prof
go tool pprof -base=before.prof after.prof
```

### ❌ DON'T

```bash
# ❌ Skip correctness tests
# Optimize first, test later (recipe for bugs!)

# ❌ Single benchmark comparison
go test -bench=.  # Before
# ... changes ...
go test -bench=.  # After
# Compare single numbers (unreliable!)

# ❌ Trust gut feeling about improvement
# "Feels faster" != actually faster
```

---

## Summary: Top 10 Rules

1. **Profile First** - Measure before optimizing
2. **Test Always** - Correctness > Performance
3. **Focus Hot Paths** - Optimize code that runs often
4. **Cache Smart** - Pre-compute frequent lookups
5. **Avoid Allocations** - Reuse memory when possible
6. **Type Specifically** - Avoid unnecessary `any` conversions
7. **Benchmark Realistically** - Mirror production usage
8. **Use Statistics** - Multiple runs with `benchstat`
9. **Document Trade-offs** - Explain optimization decisions
10. **Keep it Simple** - Readable code is maintainable code

---

## Quick Decision Tree

```
Need to optimize?
│
├─ Is it a hot path? (>5% CPU in profile)
│  ├─ YES → Continue
│  └─ NO → Don't optimize (not worth complexity)
│
├─ Have you profiled?
│  ├─ YES → Continue
│  └─ NO → Profile first (don't guess!)
│
├─ Have you benchmarked baseline?
│  ├─ YES → Continue
│  └─ NO → Establish baseline (need comparison point)
│
├─ Is there a simple optimization?
│  ├─ YES → Apply & measure
│  └─ NO → Is complexity worth the gain?
│     ├─ YES → Document trade-offs & proceed
│     └─ NO → Keep it simple
│
└─ After optimization:
   ├─ Run tests (correctness)
   ├─ Run benchmarks (performance)
   ├─ Compare with benchstat (statistics)
   └─ Profile again (verify improvement)
```

---

## References

- [optimization-journey.md](optimization-journey.md) - Real-world examples
- [optimization-techniques.md](optimization-techniques.md) - Detailed techniques
- [best-practices.md](best-practices.md) - Philosophy and guidelines
- [profiling-guide.md](profiling-guide.md) - CPU profiling walkthrough
- [benchmarking-guide.md](benchmarking-guide.md) - Benchmark best practices

**Last Updated:** October 17, 2025
