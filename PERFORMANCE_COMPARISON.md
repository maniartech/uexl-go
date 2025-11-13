# Performance Comparison: Before vs After Value Migration

## BEFORE (Phase 1 Baseline - with []any stack)
```
Operation       Time (ns/op)    Allocations (B/op)    Allocs/op
-----------------------------------------------------------------
Arithmetic         131.9              ?                   4
String              95.35             ?                   2
Boolean             75.69             ?                   0
```

## AFTER (Phase 2B - with []Value stack)
```
Operation       Time (ns/op)    Allocations (B/op)    Allocs/op
-----------------------------------------------------------------
Boolean             69.13              0                   0  ✅
Arithmetic         156.7              40                   5
Comparison         152.7              32                   4
String             285.7             104                   7
```

## Analysis

### ✅ WINS
1. **Boolean: 0 allocs achieved!** (75.69ns → 69.13ns, -8.7% faster, 0 allocs)
   - Booleans stored inline in Value struct
   - No interface boxing during operations
   - Slightly faster due to better cache locality

### ⚠️ PARTIAL SUCCESS
2. **Arithmetic: Still 5 allocs** (131.9ns → 156.7ns, +18.8% slower)
   - Stack operations are zero-alloc ✅
   - **BUT**: Constants loaded from `[]any` slice → boxes on access
   - **AND**: Final result boxes when returning from `Run()`
   - Time regression due to Value wrapper overhead

3. **String: 7 allocs** (95.35ns → 285.7ns, +200% slower!!!)
   - String concatenation creates new strings (expected)
   - Multiple constant loads from `[]any` array
   - Significant slowdown from wrapper overhead + boxing

### Root Cause of Remaining Allocations

**Constants Array Boxing:**
- `compiler.ByteCode.Constants` is `[]any`
- Each `OpConstant` instruction: `vm.Push(vm.constants[idx])`
- Accessing `vm.constants[idx]` **boxes the value** into interface
- Expression `(10.0 + 20.0) * 5.0` has 3 constants → 3 boxes

**Return Value Boxing:**
- `Run()` returns `any` interface
- Final `LastPoppedStackElem().ToAny()` boxes result
- This is unavoidable for public API

**Allocation Breakdown for `(10.0 + 20.0) * 5.0`:**
1. Load constant 10.0 from `[]any` → **1 alloc**
2. Load constant 20.0 from `[]any` → **1 alloc**
3. OpAdd (zero alloc) ✅
4. Load constant 5.0 from `[]any` → **1 alloc**
5. OpMul (zero alloc) ✅
6. Pop result → **1 alloc** (in ToAny)
7. Return from Run() → **1 alloc** (LastPoppedStackElem)
**Total: 5 allocs** ✅ matches observed

### Why Boolean Succeeded but Arithmetic Didn't

**Boolean expression: `true && false || true`**
- Uses `OpTrue`, `OpFalse` opcodes (not `OpConstant`)
- These directly call `vm.Push(true/false)` with literal values
- No constants array access → no boxing!
- Result: **0 allocations** ✅

**Arithmetic uses numeric literals:**
- Compiler stores numbers in Constants array
- VM loads via `OpConstant` from `[]any` slice
- Each load boxes the float64
- Result: **5 allocations** ❌

## Next Steps to Achieve 0 Allocs

### Phase 2C: Compiler Constants Migration
1. Change `compiler.ByteCode.Constants` from `[]any` to `[]Value`
2. Compiler stores Values directly when emitting constants
3. VM accesses constants without boxing
4. **Expected result: Arithmetic 0-1 allocs** (only final return boxes)

### Estimated Final Performance
```
Boolean:     0 allocs  ✅ (already achieved)
Arithmetic:  1 alloc   (only final return value)
String:      1 alloc   (only final return, concat still allocates strings)
Comparison:  1 alloc   (only final return)
```

The 1 remaining allocation is unavoidable unless we change `Run()` signature to return `Value` instead of `any`, which would break the public API.
