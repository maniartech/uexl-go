# Phase 2D: Zero-Allocation VM Operations

## Problem Analysis

**Current State**: 4 allocations per expression evaluation
- expr: 1 alloc (final return only)
- celgo: 1 alloc
- UExL: 4 allocs

**Root Cause**: `Pop()` boxes `Value→any` on every call

**Bytecode Analysis** for `(Origin == "MOW" || Country == "RU") && (Value >= 100 || Adults == 1)`:
```
OpContextVar [0]      (Origin) - pushValue ✅ zero-alloc
OpConstant [0]        ("MOW")  - pushValue ✅ zero-alloc
OpEqual               - Pop()×2 ❌ 2 allocs
OpJumpIfTruthy        - Pop()×1 ❌ 1 alloc
OpJumpIfFalsy         - Pop()×1 ❌ 1 alloc
Total: 4 allocations
```

## Holistic Solution

### Phase 1: Internal Value Operations
Add zero-alloc internal stack operations:
- `popValue() Value`
- `pop2Values() (Value, Value)`
- `peekValue() Value`
- `pushBoolValue(bool) error`

### Phase 2: Rewrite Opcode Handlers
Convert to Value-native operations:
1. **Comparison ops** (OpEqual, OpNotEqual, etc)
   - Use `pop2Values()` instead of `Pop()`
   - Add `executeComparisonOperationValues(op, left Value, right Value)`
   - Type-switch on `Value.Typ` for zero-alloc comparison

2. **Jump ops** (OpJumpIfTruthy, OpJumpIfFalsy)
   - Use `popValue()` instead of `Pop()`
   - Check truthiness on Value directly
   - Use `pushValue()` to keep value on stack

3. **Binary ops** (OpAdd, OpSub, etc)
   - Use `pop2Values()`
   - Add Value-native arithmetic handlers

4. **Unary ops** (OpMinus, OpBang, OpBitwiseNot)
   - Use `popValue()`
   - Add Value-native unary handlers

### Phase 3: Final Return
Keep `LastPoppedStackElem()` as-is - this 1 allocation is unavoidable

## Expected Results
- **Target**: 1 alloc/op (match expr/celgo)
- **Reduction**: 75% fewer allocations (4→1)
- **Speed**: Likely 10-20% faster (fewer interface conversions)

## Implementation Order
1. ✅ Context variable cache (already done - loads are zero-alloc)
2. Add internal Value stack methods
3. Rewrite comparison operations
4. Rewrite jump operations
5. Benchmark and validate
6. Rewrite binary/unary operations (if needed)

## Risk Assessment
- **Low risk**: Changes are internal to VM
- **High confidence**: Pattern already proven with constants/context vars
- **Validation**: Existing test suite (248 tests)
