# Phase 2 Performance Optimization Plan

## Current Performance vs Targets

| Benchmark | Current | Target | Gap | Priority |
|-----------|---------|--------|-----|----------|
| `Benchmark_uexl` | 67.63 ns/op | < 10 ns | 87% reduction needed | 🟡 Medium |
| `Benchmark_uexl_startswith` | 260.7 ns/op | < 20 ns | 92% reduction needed | 🔴 High |
| `Benchmark_uexl_func` | 114.4 ns/op | < 20 ns | 83% reduction needed | 🔴 High |
| `Benchmark_uexl_map` | 30,026 ns/op | < 100 ns | 99.7% reduction needed | 🔴 Critical |

## Root Cause Analysis

### 1. String Operations Performance Issues (260.7 ns vs 20 ns target)

**Problem**: `name == "/groups/" + group + "/bar"` requires:
- Variable lookup: `name` (context access)
- String concatenation: `"/groups/" + group + "/bar"` (3 operations)
- String comparison: `==`

**Bottlenecks Identified**:
1. **Context Variable Access**: `getContextValue()` involves map lookup + interface{} boxing
2. **String Concatenation**: Multiple `executeStringBinaryOperation()` calls with stack operations
3. **Intermediate Stack Operations**: Push/Pop for each intermediate result

**Current Execution Path**:
```
OpContextVar (name) -> Push -> OpConstant ("/groups/") -> Push ->
OpContextVar (group) -> Push -> OpAdd -> OpConstant ("/bar") -> Push ->
OpAdd -> OpContextVar (name) -> Push -> OpEqual
```

### 2. Function/String Concatenation Issues (114.4 ns vs 20 ns target)

**Problem**: `"hello" + ", world"` should be trivial but takes 114 ns

**Bottlenecks Identified**:
1. **Bytecode Overhead**: Simple operations going through full VM execution
2. **Stack Operations**: Unnecessary Push/Pop for constants
3. **Type Checking**: Runtime type assertions for string operations

### 3. Map Operations Critical Issue (30,026 ns vs 100 ns target)

**Problem**: `array |map: $item * 2` on 100-element array is extremely slow

**Bottlenecks Identified**:
1. **Pipe Scope Creation**: `pushPipeScope()` + `popPipeScope()` per iteration (100x)
2. **Frame Management**: `pushFrame()` + `popFrame()` per iteration (100x)
3. **VM Context Switching**: Full `vm.run()` call per iteration (100x)
4. **Variable Setting**: `setPipeVar()` calls per iteration (100x)
5. **Memory Allocations**: `make([]any, len(arr))` + repeated scope maps

**Per-iteration Overhead**:
- Scope creation: ~50-100 ns
- Frame management: ~50-100 ns
- VM execution: ~100-200 ns
- **Total per item**: ~200-400 ns × 100 items = 20,000-40,000 ns

## Optimization Strategy

### Phase 2A: Fast Path String Operations (Target: 90% reduction)

**1. String Concatenation Fast Path**
```go
// Detect pattern: constant + variable + constant
// Compile to specialized OpStringTemplate instruction
case code.OpStringTemplate:
    template := vm.constants[constantIndex].(StringTemplate)
    result := template.Execute(vm.contextValues)
    vm.Push(result)
```

**2. Context Variable Caching**
```go
type FastContext struct {
    vars   []any              // Indexed access instead of map
    lookup map[string]int     // Variable name -> index mapping
}
```

**3. String Operation Specialization**
```go
case code.OpStringEquals:
    // Specialized string comparison without generic binary operation
    left := vm.stack[vm.sp-2].(string)
    right := vm.stack[vm.sp-1].(string)
    vm.sp -= 2
    vm.stack[vm.sp] = left == right
    vm.sp++
```

### Phase 2B: Map Operation Optimization (Target: 99% reduction)

**1. Bulk Operation Mode**
```go
type BulkMapOperation struct {
    inputArray    []any
    expression    *compiler.ByteCode
    resultBuffer  []any  // Pre-allocated
}

func (vm *VM) executeBulkMap(op *BulkMapOperation) error {
    // Single scope setup
    vm.pushPipeScope()
    defer vm.popPipeScope()

    // Reuse same execution context
    for i, item := range op.inputArray {
        vm.setPipeVar("$item", item)    // Fast variable update
        vm.setPipeVar("$index", i)

        // Direct expression evaluation without frame overhead
        result, err := vm.evaluateExpression(op.expression)
        if err != nil {
            return err
        }
        op.resultBuffer[i] = result
    }
    return nil
}
```

**2. Expression Pre-compilation**
```go
// Compile map expressions to specialized bytecode
case parser.PipeExpression:
    if isMapPipe(node) {
        return compiler.compileBulkMapOperation(node)
    }
```

**3. Memory Pool for Arrays**
```go
var arrayPool = sync.Pool{
    New: func() interface{} {
        return make([]any, 0, 256)  // Pre-sized buffer
    },
}
```

### Phase 2C: General VM Optimizations

**1. Instruction Fusion**
```go
// Combine common patterns into single instructions
case code.OpContextVarAdd:  // context_var + constant -> single op
case code.OpContextVarEquals: // context_var == constant -> single op
case code.OpStringTemplate:   // template string with variables
```

**2. Stack Operation Optimization**
```go
// Eliminate unnecessary Push/Pop cycles
func (vm *VM) pushDirect(value any) {
    vm.stack[vm.sp] = value
    vm.sp++
}
```

## Implementation Priority

### Priority 1: Map Operations (Biggest Impact)
- **Target**: 30,026 ns → 100 ns (99.7% reduction)
- **Approach**: Bulk operation mode + memory pooling
- **Expected ROI**: Highest - this is the biggest performance gap

### Priority 2: String Operations (High Impact)
- **Target**: 260.7 ns → 20 ns (92% reduction)
- **Approach**: Fast path string operations + context caching
- **Expected ROI**: High - significant improvement potential

### Priority 3: Simple Operations (Moderate Impact)
- **Target**: 114.4 ns → 20 ns (83% reduction)
- **Approach**: Instruction fusion + stack optimization
- **Expected ROI**: Moderate - but easier to implement

### Priority 4: Boolean Operations (Polish)
- **Target**: 67.63 ns → 10 ns (85% reduction)
- **Approach**: Advanced optimizations (JIT, template specialization)
- **Expected ROI**: Low - already very competitive

## Success Metrics

### Phase 2A Goals (Week 1):
- ✅ String operations: < 50 ns/op (80% reduction)
- ✅ Function operations: < 30 ns/op (75% reduction)

### Phase 2B Goals (Week 2):
- ✅ Map operations: < 1,000 ns/op (97% reduction)
- ✅ Bulk operation infrastructure working

### Phase 2C Goals (Week 3):
- ✅ All targets achieved:
  - Map: < 100 ns/op
  - String: < 20 ns/op
  - Function: < 20 ns/op
  - Boolean: < 10 ns/op (stretch goal)

## Risk Assessment

- **Low Risk**: String operation optimization (well-understood bottlenecks)
- **Medium Risk**: Map bulk operations (requires significant VM changes)
- **High Risk**: Sub-10 ns boolean operations (may require JIT)

## Next Steps

1. **Implement Priority 1**: Map operation bulk mode
2. **Create benchmarks**: Track progress for each optimization
3. **Validate approach**: Ensure correctness while improving performance
4. **Iterate**: Measure and refine based on results