# Development Guidelines for UExL

## Performance-Critical Rules

These rules MUST be followed to maintain UExL's zero-allocation performance.

### Rule 1: NEVER Use `Pop()` in Opcode Handlers

**❌ WRONG**:
```go
case code.OpMyNewOp:
    val := vm.Pop()  // ALLOCATES! Boxes Value → any
    result := process(val)
    vm.Push(result)
```

**✅ CORRECT**:
```go
case code.OpMyNewOp:
    val := vm.popValue()  // Zero-alloc, returns Value
    result := vm.processValue(val)  // Value-native processing
    vm.pushValue(result)
```

**Why**: `Pop()` calls `ToAny()` which allocates. Use `popValue()` instead.

**Exception**: Only use `Pop()` when:
- Calling external functions that require `any`
- Building arrays/objects (elements are `any`)
- At the final API boundary (return from `Run()`)

### Rule 2: NEVER Box Values in Hot Paths

**❌ WRONG**:
```go
func (vm *VM) executeMyOperation(val Value) error {
    anyVal := val.ToAny()  // ALLOCATES!
    return vm.doSomething(anyVal)
}
```

**✅ CORRECT**:
```go
func (vm *VM) executeMyOperation(val Value) error {
    switch val.Typ {
    case TypeFloat:
        return vm.doSomethingFloat(val.FloatVal)  // Direct access
    case TypeString:
        return vm.doSomethingString(val.StrVal)
    }
}
```

**How to identify hot paths**:
- Any opcode handler
- Stack operations (push/pop)
- Type checking/conversion
- Comparison operations

### Rule 3: NEVER Modify Value Struct Layout

**❌ WRONG**:
```go
type Value struct {
    Typ      valueType  // Small field first = bad!
    FloatVal float64    // 7 bytes padding added
    StrVal   string
    BoolVal  bool       // More padding
    AnyVal   any        // More padding
}
// Result: 56+ bytes with wasted padding
```

**✅ CORRECT**:
```go
type Value struct {
    AnyVal   any        // Largest fields first
    StrVal   string     // 16-byte types
    FloatVal float64    // 8-byte types
    Typ      valueType  // Small types last
    BoolVal  bool
}
// Result: 48 bytes, minimal padding
```

**Before modifying**:
1. Check current size: `unsafe.Sizeof(Value{})`
2. Make change
3. Check new size
4. If larger, reorder fields
5. Run benchmarks to verify performance

### Rule 4: ALWAYS Use Constructors

**❌ WRONG**:
```go
val := Value{
    Typ:      TypeFloat,
    FloatVal: 42.0,
    // Forgot other fields - bugs lurking!
}
```

**✅ CORRECT**:
```go
val := newFloatValue(42.0)  // Guaranteed correct initialization
```

**Available constructors**:
- `newFloatValue(float64)` - For numbers
- `newStringValue(string)` - For strings
- `newBoolValue(bool)` - For booleans
- `newNullValue()` - For null
- `newAnyValue(any)` - For complex types (arrays, maps, objects)

### Rule 5: ALWAYS Profile Before Optimizing

**Required process**:
```bash
# 1. Capture baseline
go test -bench=BenchmarkMyFeature -benchmem -count=10 > baseline.txt

# 2. Make changes

# 3. Capture new results
go test -bench=BenchmarkMyFeature -benchmem -count=10 > current.txt

# 4. Compare
benchstat baseline.txt current.txt

# 5. Check allocations specifically
go test -bench=BenchmarkMyFeature -benchmem | grep allocs
# MUST be 0 for primitive operations!
```

**Red flags**:
- Any increase in allocations
- Speed regression > 5%
- Increased memory usage without justification

### Rule 6: Test All Changes

**Minimum required tests**:
```bash
# 1. Unit tests
go test ./...

# 2. Specific package tests
go test ./vm -v

# 3. Benchmark comparison
go test -bench=. -benchmem

# 4. Integration test (comparison project)
cd ../golang-expression-evaluation-comparison
go test -bench=Benchmark_uexl -benchmem
```

**All must pass before merging**.

## Code Organization Rules

### Internal vs Public APIs

**Pattern**: Separate zero-alloc internal from boxed public APIs

```go
// Internal (zero-alloc) - lowercase, package-private
func (vm *VM) popValue() Value
func (vm *VM) pushValue(Value) error
func (vm *VM) executeComparisonOperationValues(code.Opcode, Value, Value) error

// Public (may box) - uppercase, exported
func (vm *VM) Pop() any {
    return vm.popValue().ToAny()  // Boxing only at API boundary
}
func (vm *VM) Push(any) error {
    return vm.pushValue(newAnyValue(val))
}
```

**Rule**: Opcode handlers ONLY use internal APIs

### File Organization

```
vm/
├── vm.go              # Main dispatch loop, opcode handlers
├── vm_utils.go        # Stack operations (push/pop/peek)
├── vm_handlers.go     # Operation implementations (comparison, binary, etc)
├── value.go           # Value type re-exports
├── builtins.go        # Built-in functions
└── pipes.go           # Pipe handlers
```

**Each file has specific responsibility - don't mix concerns.**

## Adding New Features

### Adding a New Opcode

**Checklist**:
1. [ ] Define opcode in `code/code.go`
2. [ ] Add compiler emission in `compiler/compiler.go`
3. [ ] Add VM handler in `vm/vm.go` using Value-native operations
4. [ ] Add helper function if needed in `vm/vm_handlers.go`
5. [ ] Write unit tests
6. [ ] Write benchmark
7. [ ] Verify 0 allocations for primitives

**Example**:
```go
// 1. code/code.go
const (
    // ...
    OpMyNewOp Opcode = 0x??
)

// 2. compiler/compiler.go
case *parser.MyNewNode:
    vm.emit(code.OpMyNewOp)

// 3. vm/vm.go
case code.OpMyNewOp:
    val := vm.popValue()  // Value-native!
    result := vm.executeMyNewOp(val)
    vm.pushValue(result)

// 4. vm/vm_handlers.go
func (vm *VM) executeMyNewOp(val Value) Value {
    switch val.Typ {
    case TypeFloat:
        // Handle float
    case TypeString:
        // Handle string
    }
}

// 5. Test
func TestMyNewOp(t *testing.T) { ... }

// 6. Benchmark
func BenchmarkMyNewOp(b *testing.B) {
    b.ReportAllocs()
    // ...
}
```

### Adding a New Built-in Function

**Checklist**:
1. [ ] Add function to `vm/builtins.go`
2. [ ] Handle Value types explicitly (don't rely on `any`)
3. [ ] Return appropriate type
4. [ ] Add tests
5. [ ] Document in `book/functions/`

**Example**:
```go
// vm/builtins.go
func builtinMyFunc(args ...any) (any, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("myFunc requires 1 argument")
    }

    // Type-specific handling
    switch val := args[0].(type) {
    case float64:
        return val * 2, nil
    case string:
        return strings.ToUpper(val), nil
    default:
        return nil, fmt.Errorf("myFunc requires number or string")
    }
}

// Register
var Builtins = VMFunctions{
    "myFunc": builtinMyFunc,
}
```

### Adding a New Pipe Handler

**Checklist**:
1. [ ] Add handler to `vm/pipes.go`
2. [ ] Use Value-native operations where possible
3. [ ] Handle pipe context variables correctly
4. [ ] Add tests
5. [ ] Document in `book/pipes/`

**Example**:
```go
// vm/pipes.go
func pipeMyHandler(input any, block any, alias string, vm *VM) (any, error) {
    // Validate input
    arr, ok := input.([]any)
    if !ok {
        return nil, fmt.Errorf("myPipe requires array")
    }

    // Process with Value context
    result := []any{}
    for i, item := range arr {
        vm.setPipeVar("$item", item)
        vm.setPipeVar("$index", float64(i))

        // Execute predicate block
        // ...
    }

    return result, nil
}

// Register
var DefaultPipeHandlers = PipeHandlers{
    "myPipe": pipeMyHandler,
}
```

## Common Mistakes

### Mistake 1: Using `any` When Value is Available

**❌ WRONG**:
```go
func compareValues(left, right any) bool {
    // Type assertions needed, slower
    l, _ := left.(float64)
    r, _ := right.(float64)
    return l == r
}
```

**✅ CORRECT**:
```go
func compareValuesValue(left, right Value) bool {
    if left.Typ != TypeFloat || right.Typ != TypeFloat {
        return false
    }
    return left.FloatVal == right.FloatVal  // Direct access
}
```

### Mistake 2: Forgetting to Check Allocations

**Always verify**:
```bash
go test -bench=BenchmarkMyFeature -benchmem
# Look for: X allocs/op
# For primitives, MUST be 0!
```

### Mistake 3: Not Using `pop2Values()` for Binary Ops

**❌ WRONG**:
```go
right := vm.popValue()
left := vm.popValue()
```

**✅ CORRECT**:
```go
right, left := vm.pop2Values()  // More efficient, clearer
```

### Mistake 4: Boxing in Comparison Operations

**❌ WRONG**:
```go
case code.OpEqual:
    right := vm.Pop()  // Boxes!
    left := vm.Pop()   // Boxes!
    vm.executeComparisonOperation(opcode, left, right)
```

**✅ CORRECT**:
```go
case code.OpEqual:
    right, left := vm.pop2Values()  // No boxing
    vm.executeComparisonOperationValues(opcode, left, right)
```

## Performance Budgets

### Must Maintain:
- **Main expression benchmark**: 180-230 ns/op
- **Primitive allocations**: 0 allocs/op (CRITICAL)
- **Boolean operations**: < 80 ns/op, 0 allocs
- **Arithmetic operations**: < 160 ns/op, 0 allocs
- **String operations**: < 300 ns/op
- **Map/filter operations**: < 5,000 ns/op

### Investigation Triggers:
- Any primitive operation allocates
- Main benchmark > 250 ns/op
- Any regression > 5% without clear benefit
- Struct size increases

## Code Review Checklist

Before submitting changes:

- [ ] All tests pass (`go test ./...`)
- [ ] Benchmarks run (`go test -bench=.`)
- [ ] No new allocations in primitive operations
- [ ] No regression > 5% in any benchmark
- [ ] Value struct size unchanged (or justified)
- [ ] Only use `Pop()` at API boundaries
- [ ] Use constructors for all Value creation
- [ ] Profile-guided optimization (not guesswork)
- [ ] Documentation updated if adding features
- [ ] Integration tests pass (comparison project)

## Getting Help

### When You See Allocations:

1. **Profile it**:
   ```bash
   go test -bench=BenchmarkProblem -memprofile=mem.prof
   go tool pprof -alloc_objects mem.prof
   top
   ```

2. **Find the allocation**:
   ```bash
   list <function_name>
   ```

3. **Common causes**:
   - Using `Pop()` instead of `popValue()`
   - Calling `ToAny()` unnecessarily
   - Boxing in comparisons
   - Building arrays/maps (sometimes unavoidable)

### When Performance Regresses:

1. **Compare carefully**:
   ```bash
   benchstat before.txt after.txt
   ```

2. **CPU profile**:
   ```bash
   go test -bench=BenchmarkProblem -cpuprofile=cpu.prof
   go tool pprof cpu.prof
   ```

3. **Look for**:
   - Increased function calls
   - More allocations
   - Larger struct copies
   - Missing inline optimizations

## Summary

### The Golden Rules:
1. ✅ Use `popValue()`, not `Pop()` in opcode handlers
2. ✅ Never box in hot paths
3. ✅ Always use constructors
4. ✅ Profile before optimizing
5. ✅ Test everything, especially allocations

### The Contract:
- Zero allocations for primitives (non-negotiable)
- Competitive performance (within 2× of fastest)
- Clean, maintainable code (no clever tricks)
- Comprehensive testing (prevent regressions)

### Remember:
> "Performance is a feature, but maintainability is a requirement."

Follow these guidelines and UExL will remain fast, efficient, and easy to work with.
