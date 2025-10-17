# Performance Optimization Rollout Plan

## Executive Summary

Apply proven optimization strategies **SYSTEM-WIDE** across the **entire UExL evaluation pipeline**: Parser → Compiler → VM. This is a **comprehensive performance overhaul** covering:

- **VM Core:** Stack operations, frame management, instruction dispatch, opcode handling
- **Context Handling:** Variable lookup, caching, scope management
- **All Operators:** Arithmetic, comparison, logical, bitwise, unary, string, array
- **All Expressions:** Binary, unary, index, member access, function calls, pipes
- **Pipe Operations:** Map, filter, reduce, find, some, every, unique, sort, groupBy, window, chunk
- **Built-in Functions:** All functions in `vm/builtins.go`
- **Type Operations:** Conversions, coercions, type checking, dispatch
- **Memory Management:** Stack/frame allocation, constant pools, scope cleanup

**Current State:** Only comparison operators & context caching optimized (62ns/op for boolean)
**Ambitious Target:** **20ns/op for simple expressions** (boolean, arithmetic, string)
**Realistic Target:** 30-40ns/op for simple expressions, sub-1µs for pipe operations
**Ultimate Goal:** **Military-grade performance across EVERY operation in the language**

**Optimization Principles:**

- ✅ **Legitimate optimizations only** - No hardcoding, no shortcuts
- ✅ **All tests must pass** - 100% test coverage maintained
- ✅ **Best practices enforced** - Profile-driven, measured improvements
- ✅ **Production-ready code** - Maintainable, readable, robust
- ✅ **Comprehensive scope** - Every function, every opcode, every code path

---

## 🎯 Complete Optimization Scope

This plan targets **EVERY component** of the UExL evaluation pipeline. Here's the complete inventory:

### **1. VM Core Operations** (`vm/vm.go`)

| Component | Current State | Target | Priority |
|-----------|---------------|--------|----------|
| **Instruction dispatch loop** | Switch-based opcode handling | Jump table / type-specialized dispatch | 🔴 CRITICAL |
| **Stack operations** | Push/Pop with bounds checking | Inline hot paths, remove redundant checks | 🔴 CRITICAL |
| **Frame management** | pushFrame/popFrame overhead | Frame pooling, reuse patterns | 🟡 HIGH |
| **Constant loading** | Map lookup + type assertion | Direct typed access, pre-cast constants | 🟡 HIGH |
| **Context variable caching** | ✅ Already optimized (array access) | - | ✅ DONE |
| **Cache invalidation** | ✅ Already optimized (pointer comparison) | - | ✅ DONE |

### **2. Operator Handlers** (`vm/vm_handlers.go`)

| Operator Category | Functions to Optimize | Current Issue | Target Pattern |
|-------------------|----------------------|---------------|----------------|
| **Arithmetic** | `executeBinaryArithmeticOperation` | Accepts `any`, type assertions inside | Type-specific signatures: `executeNumberArithmetic(op, left float64, right float64)` |
| **Comparison** | `executeNumberComparisonOperation`, `executeStringComparisonOperation`, `executeBooleanComparisonOperation` | ✅ Already type-specific | ✅ DONE |
| **Logical** | `executeBinaryExpression` (&&, \|\|) | Generic dispatch | Type-specific for boolean shortcuts |
| **Bitwise** | Embedded in `executeBinaryExpression` | Not separated, `any` types | `executeBitwiseOperation(op, left int64, right int64)` |
| **String** | `executeStringBinaryOperation`, `executeStringConcat` | Accepts `any`, type assertions | `executeStringAddition(left string, right string)` |
| **Unary** | `executeUnaryMinusOperation`, `executeUnaryBangOperation` | Accepts `any`, type assertions | `executeNumberNegate(val float64)`, `executeBooleanNegate(val bool)` |

### **3. Index & Member Access** (`vm/vm_handlers.go`, `vm/indexing.go`)

| Operation | Current Implementation | Optimization Target |
|-----------|------------------------|---------------------|
| **Array indexing** | `executeIndexValue()` - double type switch | Pre-check types once, dispatch to typed handlers |
| **Object member access** | `executeMemberAccess()` - map lookup per access | Member access caching (if safe), direct map operations |
| **Optional chaining** | `?.` and `?.[` - null checks per operation | Fast-path for non-null common case |
| **Slicing** | `executeSliceExpression()` - generic slice handling | Type-specific slicing for arrays vs strings |

### **4. Pipe Operations** (`vm/pipes.go`)

| Pipe Handler | Current State | Optimization Needed |
|--------------|---------------|---------------------|
| `MapPipeHandler` | ✅ Scope/frame reuse implemented | ✅ DONE |
| `FilterPipeHandler` | Creates new scope per iteration | Apply scope reuse pattern from Map |
| `ReducePipeHandler` | Creates new scope per iteration | Apply scope reuse pattern from Map |
| `FindPipeHandler` | Creates new scope per iteration | Apply scope reuse pattern from Map |
| `SomePipeHandler` | Creates new scope per iteration | Apply scope reuse pattern from Map |
| `EveryPipeHandler` | Creates new scope per iteration | Apply scope reuse pattern from Map |
| `UniquePipeHandler` | Standard implementation | Optimize deduplication logic |
| `SortPipeHandler` | Standard implementation | Optimize comparator function calls |
| `GroupByPipeHandler` | Standard implementation | Optimize key extraction & grouping |
| `WindowPipeHandler` | Standard implementation | Optimize window creation & iteration |
| `ChunkPipeHandler` | Standard implementation | Optimize chunk allocation |

### **5. Built-in Functions** (`vm/builtins.go`)

**ALL built-in functions need profiling and optimization:**

- **String functions:** `len()`, `substr()`, `indexOf()`, `lastIndexOf()`, `contains()`, `startsWith()`, `endsWith()`, `toLowerCase()`, `toUpperCase()`, `trim()`, `trimStart()`, `trimEnd()`, `replace()`, `split()`, `join()`
- **Array functions:** `len()`, `push()`, `pop()`, `shift()`, `unshift()`, `slice()`, `splice()`, `concat()`, `reverse()`, `includes()`
- **Math functions:** `abs()`, `ceil()`, `floor()`, `round()`, `min()`, `max()`, `pow()`, `sqrt()`, `sin()`, `cos()`, `tan()`, `log()`, `exp()`
- **Type functions:** `type()`, `string()`, `number()`, `boolean()`, `keys()`, `values()`
- **Utility functions:** `range()`, `coalesce()`, `default()`

**Optimization targets:**

- Remove repeated type assertions (check once, dispatch to typed helpers)
- Inline simple operations (e.g., `len()` for strings/arrays)
- Pre-allocate result buffers where size is known
- Use specialized stdlib functions (e.g., `strings.Builder` for concatenation)

### **6. Type System Operations** (`vm/vm_utils.go`, `vm/vm_handlers.go`)

| Operation | Current Approach | Optimization Target |
|-----------|------------------|---------------------|
| **Type checking** | `switch v := value.(type)` repeated | Type cache/bitmap for hot values |
| **Type dispatch** | Runtime type assertions per operation | Pre-computed type dispatch tables |
| **Type conversion** | Generic conversion functions | Type-specific conversion paths |
| **Type coercion** | Accepts `any`, type switches inside | Early type resolution, typed APIs |

### **7. Memory Management**

| Component | Current State | Optimization Target |
|-----------|---------------|---------------------|
| **Stack allocation** | Fixed 1024-slot array | Pre-allocated, never resized (good ✅) |
| **Frame allocation** | New frame object per scope | Frame pooling (sync.Pool) |
| **Scope maps** | New map per pipe iteration | Reuse pattern (clear + update) ✅ MapPipeHandler done |
| **String building** | Direct concatenation | strings.Builder for multi-part |
| **Constant pool** | Mixed types (`[]any`) | Type-segregated pools (numbers, strings, etc.) |
| **Result allocations** | Returned as `any` | Consider typed result channels |

### **8. Compiler Optimizations** (`compiler/`)

While primarily a VM optimization effort, compiler improvements can help:

- **Constant folding:** `2 + 3` → `OpConstant(5)` at compile time
- **Type hints:** If compiler knows types, emit specialized opcodes
- **Dead code elimination:** Remove unreachable code paths
- **Instruction combining:** Merge consecutive compatible ops
- **Peephole optimization:** Replace instruction sequences with faster equivalents

### **9. Control Flow Operations** (`vm/vm.go`)

| Opcode | Current Implementation | Optimization Target |
|--------|------------------------|---------------------|
| `OpJump` | Instruction pointer update | Inline (already fast) |
| `OpJumpIfTruthy` | Stack pop + truthiness check + jump | Fast-path for boolean true/false |
| `OpJumpIfFalsy` | Stack pop + truthiness check + jump | Fast-path for boolean true/false |
| `OpJumpIfNullish` | Stack pop + null check + jump | Fast-path for non-null |
| `OpJumpIfNotNullish` | Stack pop + null check + jump | Fast-path for non-null |

### **10. Special Operations**

| Operation | Location | Optimization Target |
|-----------|----------|---------------------|
| **Nullish coalescing** (`??`) | `OpNullish` handler | Fast-path for non-null left |
| **Optional chaining** (`?.`, `?.[`) | `OpSafeModeOn/Off` | Minimize safe mode overhead |
| **String pattern matching** | `OpStringPatternMatch` | Optimize prefix/suffix checks |
| **Function calls** | `OpCallFunction` → built-in lookup | Function dispatch table, inline common functions |
| **Object construction** | `OpObject` | Pre-allocate map with known size |
| **Array construction** | `OpArray` | Pre-allocate slice with exact capacity |

---

## 📊 Optimization Phases (Prioritized by Impact)

Based on profiling data and impact analysis, here's the systematic rollout order:

---

## Phase 1: Arithmetic Operations ⚡ HIGH PRIORITY

### Current State Analysis

```go
// Current: Already has some optimization
func (vm *VM) executeBinaryArithmeticOperation(operator code.Opcode, left, right any) error {
    leftValue := left.(float64)   // ⚠️ Type assertion overhead
    rightValue := right.(float64)

    switch operator {
    case code.OpAdd:
        return vm.Push(leftValue + rightValue)
    // ...
    }
}
```

**Issues:**
- Still accepts `any` parameters → interface overhead
- Type assertions happen inside function (not in dispatcher)
- No separate int path (all ints converted to float64)

### Optimization Strategy

**Apply Pattern:** Type-specific function signatures (like comparison operators)

```go
// ✅ Optimized: Type-specific signatures
func (vm *VM) executeNumberArithmetic(operator code.Opcode, left, right float64) error {
    // No type assertions - direct operations
    switch operator {
    case code.OpAdd:
        return vm.Push(left + right)
    case code.OpSub:
        return vm.Push(left - right)
    case code.OpMul:
        return vm.Push(left * right)
    case code.OpDiv:
        if right == 0 {
            return fmt.Errorf("division by zero")
        }
        return vm.Push(left / right)
    case code.OpMod:
        return vm.Push(math.Mod(left, right))
    case code.OpPow:
        return vm.Push(math.Pow(left, right))
    default:
        return fmt.Errorf("unknown arithmetic operator: %v", operator)
    }
}

// Dispatcher does type check once
func (vm *VM) executeBinaryExpression(operator code.Opcode, left, right any) error {
    switch l := left.(type) {
    case float64:
        if r, ok := right.(float64); ok {
            return vm.executeNumberArithmetic(operator, l, r)  // Typed call
        }
    case string:
        if r, ok := right.(string); ok {
            return vm.executeStringBinaryOperation(operator, l, r)
        }
    // ...
    }
}
```

**Expected Gain:** 5-8% on arithmetic-heavy expressions

**Files to Modify:**
- `vm/vm_handlers.go`: Update `executeBinaryArithmeticOperation()`
- `vm/vm_handlers.go`: Update dispatcher in `executeBinaryExpression()`

**Testing:**
- Add benchmark: `BenchmarkVM_Arithmetic_Current`
- Existing tests: `arithmetic_test.go`

---

## Phase 2: String Operations 📝 HIGH PRIORITY

### Current State Analysis

```go
// Current: Already partially optimized
func (vm *VM) executeStringBinaryOperation(operator code.Opcode, left, right any) error {
    switch operator {
    case code.OpAdd:
        l, lok := left.(string)   // ⚠️ Type assertions inside
        r, rok := right.(string)
        if !lok || !rok {
            return fmt.Errorf("...")
        }
        result := l + r
        return vm.Push(result)
    }
}
```

**Issues:**
- Accepts `any`, does type assertions inside
- Should receive typed parameters from dispatcher
- `executeStringConcat` is already optimized ✅

### Optimization Strategy

```go
// ✅ Optimized: Typed signature
func (vm *VM) executeStringAddition(left, right string) error {
    // Direct concatenation - no type checks
    return vm.Push(left + right)
}

// For 3+ strings, already optimal with strings.Builder
func (vm *VM) executeStringConcat(count int) error {
    // Already optimized ✅
}
```

**Expected Gain:** 3-5% on string operations

**Files to Modify:**
- `vm/vm_handlers.go`: Update `executeStringBinaryOperation()`

---

## Phase 3: Pipe Operations 🔄 CRITICAL PRIORITY

### Current State Analysis

**Map Handler:**
```go
func MapPipeHandler(input any, block any, alias string, vm *VM) (any, error) {
    arr, ok := input.([]any)  // Type check every map call
    if !ok {
        return nil, fmt.Errorf("map pipe expects array input")
    }

    // ⚠️ Issues:
    // 1. Creates new pipe scope for EVERY iteration
    // 2. Creates new frame for EVERY iteration
    // 3. setPipeVar() does map lookups

    for i, elem := range arr {
        vm.pushPipeScope()  // Allocation + map creation
        vm.setPipeVar("$item", elem)  // Map write
        vm.setPipeVar("$index", i)    // Map write

        frame := NewFrame(blk.Instructions, 0)  // Allocation
        vm.pushFrame(frame)
        err := vm.run()
        // ...
    }
}
```

**Major Bottlenecks:**
1. **Scope creation overhead:** New map per iteration
2. **Frame allocation:** New Frame struct per iteration
3. **Map operations:** setPipeVar() writes to map
4. **No type specialization:** Generic `[]any` handling

### Optimization Strategy

#### 3.1: Reuse Pipe Scope (Already Partially Done ✅)

```go
// ✅ Already optimized in MapPipeHandler
vm.pushPipeScope()  // Once before loop
frame := NewFrame(blk.Instructions, 0)  // Once

for i, elem := range arr {
    // Direct map access (already optimized)
    pipeScope := vm.pipeScopes[len(vm.pipeScopes)-1]
    pipeScope["$item"] = elem
    pipeScope["$index"] = i

    frame.ip = 0  // Reset frame
    vm.pushFrame(frame)
    // ...
}
vm.popPipeScope()  // Once after loop
```

**Status:** ✅ Map handler already uses this pattern!

#### 3.2: Apply to Other Pipe Handlers

**Filter Handler - NEEDS OPTIMIZATION:**
```go
// ❌ Current: Creates scope per iteration
func FilterPipeHandler(input any, block any, alias string, vm *VM) (any, error) {
    // ...
    for i, elem := range arr {
        vm.pushPipeScope()  // ⚠️ Every iteration!
        vm.setPipeVar("$item", elem)
        vm.setPipeVar("$index", i)
        frame := NewFrame(blk.Instructions, 0)  // ⚠️ Every iteration!
        // ...
        vm.popPipeScope()
    }
}

// ✅ Optimized: Reuse scope and frame
func FilterPipeHandler(input any, block any, alias string, vm *VM) (any, error) {
    arr, ok := input.([]any)
    if !ok {
        return nil, fmt.Errorf("filter pipe expects array input")
    }
    blk, ok := block.(*compiler.InstructionBlock)
    if !ok || blk == nil || blk.Instructions == nil {
        return nil, fmt.Errorf("filter pipe expects a predicate block")
    }

    var result []any

    // Reuse scope and frame (like MapPipeHandler)
    vm.pushPipeScope()
    frame := NewFrame(blk.Instructions, 0)

    // Pre-initialize keys
    if alias != "" {
        vm.setPipeVar(alias, nil)
    }
    vm.setPipeVar("$item", nil)
    vm.setPipeVar("$index", 0)

    for i, elem := range arr {
        // Direct map access
        pipeScope := vm.pipeScopes[len(vm.pipeScopes)-1]
        if alias != "" {
            pipeScope[alias] = elem
        }
        pipeScope["$item"] = elem
        pipeScope["$index"] = i

        // Reset frame for reuse
        frame.ip = 0
        frame.basePointer = vm.sp

        vm.pushFrame(frame)
        err := vm.run()
        if err != nil {
            vm.popFrame()
            vm.popPipeScope()
            return nil, err
        }
        res := vm.Pop()
        vm.popFrame()

        if b, ok := res.(bool); ok && b {
            result = append(result, elem)
        }
    }

    vm.popPipeScope()
    return result, nil
}
```

**Apply Same Pattern To:**
- ✅ `MapPipeHandler` - Already optimized
- ❌ `FilterPipeHandler` - NEEDS OPTIMIZATION
- ❌ `ReducePipeHandler` - NEEDS OPTIMIZATION
- ❌ `FindPipeHandler` - NEEDS OPTIMIZATION
- ❌ `SomePipeHandler` - NEEDS OPTIMIZATION
- ❌ `EveryPipeHandler` - NEEDS OPTIMIZATION

#### 3.3: Type-Specific Pipe Handlers

**Current:** Generic `[]any` handling
**Optimization:** Type-specific fast paths for common cases

```go
// Fast path for numeric arrays
func MapPipeHandlerNumbers(arr []float64, block *compiler.InstructionBlock, vm *VM) ([]any, error) {
    result := make([]any, len(arr))

    // Detect simple patterns at compile time
    if isSimpleArithmetic(block.Instructions) {
        // Vectorized operation (SIMD-friendly)
        for i, num := range arr {
            result[i] = num * 2.0  // Example: $item * 2
        }
        return result, nil
    }

    // General numeric path (still faster than generic)
    // ... reuse scope pattern ...
}

// Dispatcher checks array element types
func MapPipeHandler(input any, block any, alias string, vm *VM) (any, error) {
    arr, ok := input.([]any)
    if !ok {
        return nil, fmt.Errorf("map pipe expects array input")
    }

    // Fast path detection
    if len(arr) > 0 {
        switch arr[0].(type) {
        case float64:
            // Check if all elements are float64
            if isHomogeneousNumbers(arr) {
                numbers := convertToFloat64Slice(arr)
                return MapPipeHandlerNumbers(numbers, block.(*compiler.InstructionBlock), vm)
            }
        case string:
            // String-specific handler
        }
    }

    // Generic fallback
    // ... existing implementation ...
}
```

**Expected Gain:** 15-25% on pipe-heavy expressions

**Files to Modify:**
- `vm/pipes.go`: Update all pipe handlers
- `vm/pipes.go`: Add type-specific fast paths

**Testing:**
- Add benchmarks: `BenchmarkVM_Pipe_Map`, `BenchmarkVM_Pipe_Filter`, etc.
- Existing tests: `pipes_test.go`

---

## Phase 4: Array/Object Access 🗂️ MEDIUM PRIORITY

### Current State Analysis

```go
func (vm *VM) executeIndexValue(target any, index any) error {
    var idx int
    switch v := index.(type) {  // Type switch per access
    case float64:
        idx = int(v)
    case int:
        idx = v
    default:
        return fmt.Errorf("array index must be int, got %T", index)
    }

    switch v := target.(type) {  // Another type switch
    case []any:
        if idx < 0 || idx >= len(v) {
            return fmt.Errorf("array index out of bounds: %d", idx)
        }
        return vm.Push(v[idx])
    // ...
    }
}
```

**Issues:**
- Two type switches per access
- Accepts `any` parameters
- Bounds checking on every access (can't be optimized away)

### Optimization Strategy

```go
// ✅ Type-specific array access
func (vm *VM) executeArrayIndexInt(arr []any, idx int) error {
    // Bounds check (unavoidable for safety)
    if idx < 0 || idx >= len(arr) {
        return fmt.Errorf("array index out of bounds: %d", idx)
    }
    return vm.Push(arr[idx])
}

func (vm *VM) executeArrayIndexFloat(arr []any, idx float64) error {
    return vm.executeArrayIndexInt(arr, int(idx))
}

// Dispatcher
func (vm *VM) executeIndexValue(target any, index any) error {
    switch arr := target.(type) {
    case []any:
        switch idx := index.(type) {
        case int:
            return vm.executeArrayIndexInt(arr, idx)
        case float64:
            return vm.executeArrayIndexFloat(arr, idx)
        default:
            return fmt.Errorf("array index must be number, got %T", index)
        }
    case string:
        // String indexing (similar optimization)
    // ...
    }
}
```

**Expected Gain:** 5-7% on array-heavy expressions

**Files to Modify:**
- `vm/vm_handlers.go`: Update `executeIndexValue()`
- `vm/vm_handlers.go`: Update `executeMapIndexAccess()`

---

## Phase 5: Unary Operations ➖ LOW PRIORITY

### Current State

```go
func (vm *VM) executeUnaryMinusOperation(operand any) error {
    switch v := operand.(type) {  // Type switch
    case float64:
        vm.Push(-v)
    case int:
        vm.Push(float64(-v))
    default:
        return fmt.Errorf("unknown operand type: %T", operand)
    }
    return nil
}
```

**Issues:**
- Type switch on every call
- Accepts `any`

### Optimization Strategy

```go
// ✅ Type-specific unary operations
func (vm *VM) executeUnaryMinusNumber(val float64) error {
    return vm.Push(-val)
}

func (vm *VM) executeUnaryBangBool(val bool) error {
    return vm.Push(!val)
}

// Dispatcher
func (vm *VM) executeUnaryExpression(operator code.Opcode, operand any) error {
    switch operator {
    case code.OpMinus:
        if val, ok := operand.(float64); ok {
            return vm.executeUnaryMinusNumber(val)
        }
        return fmt.Errorf("unary minus requires number, got %T", operand)
    case code.OpBang:
        if val, ok := operand.(bool); ok {
            return vm.executeUnaryBangBool(val)
        }
        // Fallback to truthy conversion
        return vm.Push(!isTruthy(operand))
    }
}
```

**Expected Gain:** 2-4% on expressions with unary operators

**Files to Modify:**
- `vm/vm_handlers.go`: Update `executeUnaryMinusOperation()`, `executeUnaryBangOperation()`

---

## Phase 6: Boolean Operations 🔘 LOW PRIORITY

### Current State

```go
func (vm *VM) executeBooleanBinaryOperation(operator code.Opcode, left, right bool) error {
    switch operator {
    case code.OpLogicalAnd:
        vm.Push(left && right)  // ⚠️ Missing return
    case code.OpLogicalOr:
        vm.Push(left || right)   // ⚠️ Missing return
    default:
        return fmt.Errorf("unsupported boolean operation: %s", operator.String())
    }
    return nil
}
```

**Issues:**
- Already type-specific ✅
- Missing early returns (minor)

### Optimization Strategy

```go
// ✅ Add early returns
func (vm *VM) executeBooleanBinaryOperation(operator code.Opcode, left, right bool) error {
    switch operator {
    case code.OpLogicalAnd:
        return vm.Push(left && right)  // Direct return
    case code.OpLogicalOr:
        return vm.Push(left || right)  // Direct return
    default:
        return fmt.Errorf("unsupported boolean operation: %s", operator.String())
    }
}
```

**Expected Gain:** 1-2% (minor improvement)

---

## Implementation Roadmap

### Week 1: Foundation
- [ ] **Day 1-2:** Phase 1 - Arithmetic operations
  - Refactor `executeBinaryArithmeticOperation()`
  - Add benchmarks
  - Run tests

- [ ] **Day 3:** Phase 2 - String operations
  - Refactor `executeStringBinaryOperation()`
  - Verify no regressions

### Week 2: Pipes (Critical)
- [ ] **Day 1-2:** Phase 3.2 - Optimize all pipe handlers
  - Apply scope reuse to FilterPipeHandler
  - Apply to ReducePipeHandler, FindPipeHandler
  - Apply to SomePipeHandler, EveryPipeHandler

- [ ] **Day 3-4:** Phase 3.3 - Type-specific pipe handlers
  - Numeric array fast path
  - String array fast path
  - Add detection logic

- [ ] **Day 5:** Pipe optimization benchmarking
  - Compare before/after for each handler
  - Document results

### Week 3: Polish
- [ ] **Day 1:** Phase 4 - Array/object access
- [ ] **Day 2:** Phase 5 - Unary operations
- [ ] **Day 3:** Phase 6 - Boolean operations
- [ ] **Day 4-5:** Comprehensive testing and profiling

---

## Success Metrics

### Performance Targets (Ambitious)

**Simple Operations (Theoretical Minimum ~15-20ns):**

| Expression Type | Current | Stretch Goal | Realistic Target | Min Improvement |
|-----------------|---------|--------------|------------------|-----------------|
| Boolean (`a && b`) | 62 ns/op | **20 ns/op** | 30-35 ns/op | 45% faster |
| Arithmetic (`a + b`) | ~80 ns/op | **20 ns/op** | 35-40 ns/op | 50% faster |
| Comparison (`a > b`) | 62 ns/op | **20 ns/op** | 30-35 ns/op | 45% faster |
| String concat (`a + b`) | ~100 ns/op | **30 ns/op** | 50-60 ns/op | 40% faster |
| Array access (`arr[0]`) | ~40 ns/op | **20 ns/op** | 25-30 ns/op | 25% faster |

**Complex Operations (Sub-microsecond):**

| Expression Type | Current | Target | Min Improvement |
|-----------------|---------|--------|-----------------|
| Map pipe (10 items) | ~1500 ns/op | <800 ns/op | 45% faster |
| Filter pipe (10 items) | ~1800 ns/op | <1000 ns/op | 45% faster |
| Reduce pipe (10 items) | ~2000 ns/op | <1200 ns/op | 40% faster |

**Allocation Target:** 0 B/op, 0 allocs/op (MANDATORY)

### Validation Checklist (MANDATORY FOR EVERY OPTIMIZATION)

#### Before Implementation
- [ ] **Profile baseline** - CPU profile with `go test -cpuprofile`
- [ ] **Benchmark baseline** - 10+ runs: `go test -bench=X -count=10 > baseline.txt`
- [ ] **Document bottleneck** - Identify specific function/line consuming >5% CPU
- [ ] **Plan optimization** - Write down expected technique and gain

#### During Implementation
- [ ] **No hardcoding** - No expression-specific optimizations
- [ ] **No shortcuts** - No test-specific code paths
- [ ] **Code review ready** - Clean, readable, maintainable code
- [ ] **Best practices** - Follow patterns from optimization-techniques.md

#### After Implementation
- [ ] **All tests pass** - `go test ./...` with 100% success
- [ ] **Benchmark improved** - `go test -bench=X -count=10 > optimized.txt`
- [ ] **Statistical validation** - `benchstat baseline.txt optimized.txt` shows p<0.05
- [ ] **Profile improved** - CPU profile shows bottleneck eliminated/reduced
- [ ] **No allocations added** - Still 0 allocs/op
- [ ] **Cross-validation** - Test with varied inputs (not just benchmark cases)
- [ ] **Regression check** - No other benchmarks regressed
- [ ] **Documentation** - Update optimization-journey.md with results

### Testing Requirements (NON-NEGOTIABLE)

#### Correctness Testing
```bash
# Must pass 100% of tests
go test ./... -v

# With race detector
go test ./... -race

# Specific component tests
go test ./vm -v
go test ./compiler -v
go test ./parser -v
```

#### Performance Testing
```bash
# Baseline (before changes)
go test -bench=BenchmarkVM_Boolean -benchtime=10s -count=10 > baseline.txt

# After optimization
go test -bench=BenchmarkVM_Boolean -benchtime=10s -count=10 > optimized.txt

# Statistical comparison (p-value must be < 0.05)
benchstat baseline.txt optimized.txt

# Must show improvement like:
# name                old time/op  new time/op  delta
# VM_Boolean_Current    62.0ns ± 2%  35.0ns ± 2%  -43.55%  (p=0.000 n=10+10)
```

#### Cross-Validation Testing
```bash
# Test with different expression types
go test -bench=BenchmarkVM_Arithmetic
go test -bench=BenchmarkVM_String
go test -bench=BenchmarkVM_Pipe

# Ensure no regressions in other areas
go test -bench=. -count=5 | tee all_benchmarks.txt
```

#### Allocation Verification
```bash
# Must remain 0 B/op, 0 allocs/op
go test -bench=BenchmarkVM_Boolean -benchmem

# Example acceptable output:
# BenchmarkVM_Boolean-8   28571428   35.2 ns/op   0 B/op   0 allocs/op
#                                                  ^^^^^^   ^^^^^^^^^^^^
#                                                  MUST BE ZERO
```

### Forbidden Optimizations (REJECTED)

These "optimizations" are NOT acceptable:

❌ **Hardcoded Results**
```go
// ❌ FORBIDDEN
if vm.currentExpr == "a && b" {
    return true, nil  // Hardcoded for benchmark
}
```

❌ **Test-Specific Code Paths**
```go
// ❌ FORBIDDEN
if testing.Testing() {
    return fastPath()  // Different behavior in tests
}
```

❌ **Expression Caching (for now)**
```go
// ❌ FORBIDDEN (not general optimization)
var exprCache = map[string]any{}
if cached, ok := exprCache[expr]; ok {
    return cached  // Caching is user's responsibility
}
```

❌ **Skipping Validation in Production**
```go
// ❌ FORBIDDEN
// +build !test
func (vm *VM) Push(val any) {
    vm.stack[vm.sp] = val  // No bounds check
    vm.sp++
}
```

❌ **Benchmark-Specific Shortcuts**
```go
// ❌ FORBIDDEN
if len(bytecode.Instructions) == 7 {  // Specific to test case
    return quickResult()
}
```

### Acceptable Optimizations (ENCOURAGED)

✅ **Type-Specific Dispatch**
```go
// ✅ ACCEPTABLE - General optimization
func (vm *VM) executeNumberArithmetic(op code.Opcode, left, right float64) error {
    // Type already known, no assertions needed
}
```

✅ **Pre-computation of Invariants**
```go
// ✅ ACCEPTABLE - Computed once, reused
func (vm *VM) setBaseInstructions(bc *ByteCode, ctx map[string]any) {
    vm.contextVarCache = buildCache(ctx)  // Pre-compute context cache
}
```

✅ **Scope/Frame Reuse**
```go
// ✅ ACCEPTABLE - Avoid allocations
vm.pushPipeScope()  // Once
for item := range arr {
    // Reuse same scope
    pipeScope["$item"] = item
}
vm.popPipeScope()  // Once
```

✅ **Fast Paths for Common Patterns**
```go
// ✅ ACCEPTABLE - Detects pattern, applies general optimization
if isSimpleArithmetic(instructions) {
    return vectorizedOperation()  // SIMD-friendly for ANY arithmetic
}
```

✅ **Sentinel Values**
```go
// ✅ ACCEPTABLE - Avoid allocations
var contextVarNotProvided = contextVarMissing{}  // Singleton
```

✅ **Build Tags for Debug/Release**
```go
// ✅ ACCEPTABLE - Different builds, not different runtime paths
// +build !debug
func (vm *VM) Push(val any) {
    vm.stack[vm.sp] = val
    vm.sp++
}
```

---

## Risk Mitigation

### Testing Strategy

1. **Preserve existing tests:** All current tests must pass
2. **Add performance tests:** Benchmark each operation type
3. **Cross-validation:** Compare with expr/cel-go on same expressions
4. **Regression detection:** CI/CD benchmark tracking

### Rollback Plan

Each phase is independent:
- Git branch per phase
- Profile before/after each phase
- Can rollback individual phases if regression detected

---

## Priority Matrix

```
High Impact + High Priority:
  ┌─────────────────────────────┐
  │ Phase 3: Pipe Operations    │ ← START HERE
  │ Phase 1: Arithmetic Ops     │
  └─────────────────────────────┘

Medium Impact:
  ┌─────────────────────────────┐
  │ Phase 2: String Operations  │
  │ Phase 4: Array Access       │
  └─────────────────────────────┘

Low Impact (polish):
  ┌─────────────────────────────┐
  │ Phase 5: Unary Operations   │
  │ Phase 6: Boolean Ops        │
  └─────────────────────────────┘
```

**Recommendation:** Start with **Phase 3 (Pipe Operations)** as pipes are the most common bottleneck in real-world expressions.

---

## Path to 20ns Target

### Reality Check: Theoretical Limits

**Current State:** 62ns for boolean `a && b`

**Breakdown of Current Execution:**
```
Context variable cache lookup (a):   ~5ns
Context variable cache lookup (b):   ~5ns
OpLogicalAnd execution:               ~3ns
Stack operations (push/pop):          ~8ns
Frame management:                     ~10ns
Instruction decoding/dispatch:        ~15ns
Bytecode cache setup:                 ~16ns
Total:                                ~62ns
```

**Theoretical Minimum for Bytecode VM:**
- Instruction fetch + decode: ~5-8ns (unavoidable)
- Stack operations: ~3-5ns (2 reads + 1 write)
- Actual computation: ~2-3ns (boolean AND)
- Frame overhead: ~2-3ns (minimal)
- **Theoretical floor: ~15-20ns**

**Path to 20ns (45% reduction from 62ns):**

1. **Eliminate cache invalidation check** (save ~10ns)
   - Current: Pointer comparison every Run() call
   - Optimization: Only check on context change (requires API change)

2. **Inline critical operations** (save ~8ns)
   - Push/Pop: Already inline-friendly, need compiler hints
   - Frame management: Simplify for simple expressions

3. **Optimize instruction dispatch** (save ~7ns)
   - Replace switch with computed goto (assembly)
   - Or: Threaded code interpretation

4. **Reduce bytecode overhead** (save ~5ns)
   - Simpler instruction encoding for common ops
   - Pre-decoded instruction stream

**Realistic without breaking changes:** 30-35ns (45% improvement)
**Aggressive with API changes:** 20-25ns (60-65% improvement)

### Optimization Phases to Reach Target

#### Phase 0: Measure & Establish Baseline (Week 0)
```bash
# Profile ALL current operations
go test -bench=. -benchtime=10s -cpuprofile=baseline_all.prof -count=10 > baseline_all.txt

# Analyze bottlenecks
go tool pprof -http=:8080 baseline_all.prof

# Document top 10 bottlenecks for each operation type
```

**Output:** Comprehensive bottleneck report

#### Phase 1: Low-Hanging Fruit (Weeks 1-2) - Target: 40-45ns
Apply type-specific optimizations to:
- Arithmetic operations
- String operations
- All pipe handlers

**Expected: 62ns → 45ns (27% improvement)**

#### Phase 2: Instruction Dispatch Optimization (Week 3) - Target: 35-40ns
- Profile instruction decoding overhead
- Experiment with dispatch table vs switch
- Consider computed goto (requires assembly)

**Expected: 45ns → 38ns (15% improvement)**

#### Phase 3: Stack Operation Inlining (Week 4) - Target: 30-35ns
- Force inline Push/Pop in release builds
- Simplify frame management for simple expressions
- Remove error checks in hot path (use build tags)

**Expected: 38ns → 32ns (15% improvement)**

#### Phase 4: Bytecode Optimization (Week 5) - Target: 25-30ns
- Optimize instruction encoding
- Pre-decode common patterns
- Special-case simple expressions

**Expected: 32ns → 28ns (12% improvement)**

#### Phase 5: Advanced Techniques (Week 6+) - Target: 20-25ns
- Threaded code interpretation
- Instruction fusion
- SIMD for array operations

**Expected: 28ns → 22ns (21% improvement)**

### Critical Success Factors

#### 1. Profile-Driven Development
```bash
# BEFORE every optimization
go test -bench=BenchmarkVM_X -cpuprofile=before.prof -count=10 > before.txt

# AFTER every optimization
go test -bench=BenchmarkVM_X -cpuprofile=after.prof -count=10 > after.txt

# Compare
go tool pprof -base=before.prof after.prof
benchstat before.txt after.txt
```

**Rule:** If profile doesn't show bottleneck, don't optimize it

#### 2. Incremental Validation
```bash
# After EACH commit
go test ./...                    # All tests pass
go test -bench=. -count=5        # No regressions
git commit -m "optimization: description (Xns → Yns)"
```

**Rule:** Every commit must be shippable

#### 3. Cross-Validation
```bash
# Test varied inputs
go test -bench=BenchmarkVM_Boolean_Simple
go test -bench=BenchmarkVM_Boolean_Complex
go test -bench=BenchmarkVM_Boolean_Nested

# Compare with competitors
cd ../golang-expression-evaluation-comparison
go test -bench=Boolean
```

**Rule:** Optimization must work for ALL cases, not just benchmarks

#### 4. Documentation Discipline
After each optimization:
1. Update optimization-journey.md with before/after
2. Add technique to optimization-techniques.md if novel
3. Update this plan with actual vs expected results

**Rule:** Undocumented optimization = didn't happen

### Measurement Protocol (MANDATORY)

#### Setup
```bash
# Clean environment
go clean -testcache

# Build test binary
go test -c -o vm.test ./vm

# Baseline (before ANY changes)
./vm.test -test.bench=BenchmarkVM_Boolean_Current \
  -test.benchtime=20s \
  -test.count=20 \
  -test.cpuprofile=baseline.prof \
  > baseline_reference.txt

# Save baseline
git add baseline_reference.txt baseline.prof
git commit -m "baseline: establish performance reference"
```

#### After Each Optimization
```bash
# Test
./vm.test -test.bench=BenchmarkVM_Boolean_Current \
  -test.benchtime=20s \
  -test.count=20 \
  -test.cpuprofile=optimized.prof \
  > optimized.txt

# Statistical analysis
benchstat baseline_reference.txt optimized.txt

# Required output:
# name                old time/op  new time/op  delta
# VM_Boolean_Current    62.0ns ± 2%  XX.Xns ± 2%  -YY.YY%  (p=0.000 n=20+20)
#                                                  ^^^^^^^^
#                                                  Must be significant (p<0.05)

# Profile comparison
go tool pprof -base=baseline.prof optimized.prof -top

# Required: Bottleneck eliminated or reduced >50%
```

#### Acceptance Criteria
✅ **p-value < 0.05** (statistically significant)
✅ **Improvement ≥ 5%** (meaningful gain)
✅ **0 allocs/op** (no new allocations)
✅ **All tests pass** (no regressions)
✅ **Stable variance (±2-3%)** (reproducible)

### Contingency Plan

If 20ns proves unachievable:

**Fallback Targets:**
- **Tier 1 (Excellent):** 25-30ns (50-55% improvement) ✅ Achievable
- **Tier 2 (Good):** 30-35ns (43-48% improvement) ✅ Very likely
- **Tier 3 (Acceptable):** 35-40ns (35-43% improvement) ✅ Guaranteed with current plan

**Each tier still beats competitors:**
- expr: 105ns
- cel-go: 127ns
- Even 40ns would be 62% faster than expr!

### Success Declaration

**Minimum Success:** 40ns (35% improvement, 0 allocs)
**Good Success:** 30ns (52% improvement, 0 allocs)
**Exceptional Success:** 20ns (68% improvement, 0 allocs)

**Current trajectory:** On path to "Good Success" (30-35ns range)

---

## Next Steps

1. **Create benchmarks for current state**
   ```bash
   go test -bench=BenchmarkVM_Arithmetic -benchtime=10s -count=10 > baseline_arithmetic.txt
   go test -bench=BenchmarkVM_Pipe -benchtime=10s -count=10 > baseline_pipes.txt
   ```

2. **Start with Phase 3 (Pipes)**
   - Highest impact
   - Most visible to users
   - Clear optimization patterns from MapPipeHandler

3. **Document progress**
   - Update optimization-journey.md after each phase
   - Add new techniques to optimization-techniques.md
   - Update pending-optimizations.md when completing items

---

## References

- [optimization-journey.md](optimization-journey.md) - Historical optimizations
- [optimization-techniques.md](optimization-techniques.md) - Proven patterns
- [best-practices.md](best-practices.md) - Guidelines
- [profiling-guide.md](profiling-guide.md) - How to profile
- [benchmarking-guide.md](benchmarking-guide.md) - How to benchmark

**Created:** October 17, 2025
**Status:** 🚀 Ready to execute
**Target:** 20ns stretch goal, 30-35ns realistic
**Timeline:** 3-6 weeks depending on target tier
**Principles:** Legitimate optimizations only, all tests passing, production-ready code