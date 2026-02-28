# UExL System-Wide Optimization Scope

> **Complete Inventory of Optimization Targets — Audited Against Source Code**

**Last Updated:** February 28, 2026
**Status:** ~64% Complete — Phases 1-4 implemented (Feb 28, 2026)

---

## 🎯 Mission Statement

**OPTIMIZE EVERYTHING IN THE UEXL EVALUATION PIPELINE**

This is **NOT** a targeted optimization of specific operators. This is a **COMPREHENSIVE SYSTEM-WIDE PERFORMANCE OVERHAUL** covering every component from Parser to Compiler to VM execution.

**Goal:** Achieve **20-35ns/op** across **ALL expression types** with **0 allocations** and **100% test pass rate**.

---

## 📊 Optimization Progress Tracker

### **Overall Status**

| Category | Targets | Optimized | Remaining | Progress |
|----------|---------|-----------|-----------|----------|
| **VM Core** | 6 components | 5 ✅ | 1 🔴 | 83% |
| **Operators** | 6 categories | 5 ✅ | 1 🔴 | 83% |
| **Index/Access** | 4 operations | 0 | 4 🔴 | 0% |
| **Pipes** | 12 handlers | 11 ✅ | 1 🔴 | 92% |
| **Built-ins** | 5 functions | 0 | 5 🔴 | 0% |
| **Type System** | 4 operations | 2 ✅ | 2 🔴 | 50% |
| **Memory Mgmt** | 6 components | 3 ✅ | 3 🔴 | 50% |
| **Compiler** | 5 optimizations | 2 ✅ | 3 🟡 | 40% |
| **Control Flow** | 5 opcodes | 5 ✅ | 0 | 100% |
| **Special Ops** | 6 operations | 1 ✅ | 5 🔴 | 17% |
| **TOTAL** | **~55** | **~35** | **~20** | **~64%** |

### **Current Performance Baseline** (from `root_bench_phase1_fixed.txt`)

```
Boolean expressions:     75.84 ns/op   0 B/op    0 allocs   ✅ OPTIMIZED
String compare:          56.72 ns/op   0 B/op    0 allocs   ✅ OPTIMIZED
String concat:           79.71 ns/op   32 B/op   2 allocs   🟡 PARTIAL (inner path typed, dispatch still boxed)
Arithmetic operations:   131.5 ns/op   32 B/op   4 allocs   🔴 NOT OPTIMIZED (dispatch path uses Pop→any)
Map pipe (100 elems):    2,204 ns/op   2,616 B/op 102 allocs ✅ Scope/frame reuse done
Filter pipe (100 elems): 4,509 ns/op   2,280 B/op 10 allocs  ✅ Scope/frame reuse done
Reduce pipe (100 elems): 5,093 ns/op   896 B/op  102 allocs  ✅ Scope/frame reuse done
Sort pipe (100 elems):   31,655 ns/op  5,472 B/op 8 allocs   ✅ Scope/frame reuse done
Compilation Boolean:     6,125 ns/op   2,352 B/op 70 allocs  🔴 NOT OPTIMIZED (one-time cost)
```

**Competitive Position:**
- ✅ Boolean: **28% faster** than expr (105ns), **40% faster** than cel-go (127ns)
- ✅ String compare: **46% faster** than expr, **55% faster** than cel-go
- 🔴 Arithmetic: slower than expr due to boxing overhead in dispatch path

---

## 🗂️ Complete Optimization Inventory

### **1. VM Core Operations** (`vm/vm.go`, `vm/vm_utils.go`)

**Priority:** 🔴 CRITICAL - Affects ALL operations

| # | Component | Current State | Optimization Target | Impact | Status |
|---|-----------|---------------|---------------------|--------|--------|
| 1.1 | **Instruction dispatch loop** | Switch-based opcode handling in `run()` | Jump table (limited by Go) or superinstructions | HIGH | 🔴 TODO |
| 1.2 | **Stack operations** | Type-specific `pushFloat64()`, `pushString()`, `pushBool()`, `pushValue()` exist alongside generic `Push(any)`. Comparisons use `pop2Values()` returning `Value`. Binary ops, unary ops, string concat, pattern matching all now use `popValue()`/`pop2Values()`. | Only `OpIndex`, `OpSlice`, `OpMemberAccess`, `OpPipe` still use `Pop()` | LOW | ✅ MOSTLY DONE |
| 1.3 | **Frame management** | `NewFrame()` heap-allocates per call. Frame 0 reused across `Run()`. Pipe handlers reuse frame object (reset ip/basePointer). No `sync.Pool`. | Frame pooling with sync.Pool for non-base frames | MEDIUM | 🔴 TODO |
| 1.4 | **Constant loading** | Constants pool is `[]types.Value` (typed). Loaded via `vm.constants[constIndex]` direct array access → `pushValue()`. Zero allocation for primitives. | — Already optimized | — | ✅ DONE |
| 1.5 | **Context variable caching** | Pre-resolved `[]Value` cache with O(1) array access | — | — | ✅ DONE |
| 1.6 | **Cache invalidation** | `reflect.ValueOf().Pointer()` comparison + length check | — | — | ✅ DONE |

**Remaining gain:** 5-10% from eliminating `Push(any)`/`Pop()` on hot paths + frame pooling

---

### **2. Operator Handlers** (`vm/vm_handlers.go`)

**Priority:** 🟡 HIGH - Direct user impact

| # | Operator Category | Current State | Remaining Work | Impact | Status |
|---|-------------------|---------------|----------------|--------|--------|
| 2.1 | **Arithmetic** | `executeBinaryExpressionValues(op, Value, Value)` dispatches by `Value.Typ` to `executeNumberArithmetic(op, float64, float64)`. Uses `pop2Values()` in `run()` — zero-alloc hot path. | — | — | ✅ DONE |
| 2.2 | **Comparison** | Fully optimized: `executeComparisonOperationValues(op, Value, Value)` dispatches by `Value.Typ`. Uses `pop2Values()` — zero alloc. `executeBooleanComparisonOperation` now uses `pushBool` (fixed). | — | — | ✅ DONE |
| 2.3 | **Logical** | `executeBinaryExpressionValues` dispatches to `executeBooleanBinaryOperation(op, bool, bool)` via `pop2Values()`. Zero-alloc. Short-circuit paths use `isTruthyValue()`. | — | — | ✅ DONE |
| 2.4 | **Bitwise** | Routed through `executeBinaryExpressionValues` → `executeNumberArithmetic` via `pop2Values()`. Zero-alloc. | — | — | ✅ DONE |
| 2.5 | **String** | `executeBinaryExpressionValues` dispatches to `executeStringAddition(string, string)` for `OpAdd`. `executeStringConcat(count)` now uses `popValue()` + `pushString()`. | — | — | ✅ DONE |
| 2.6 | **Unary** | `executeUnaryExpressionValue(op, Value)` dispatches by `Value.Typ`. Fast paths: `OpMinus` → `pushFloat64(-FloatVal)`, `OpBang` → `pushBool(!BoolVal)` or `!isTruthyValue()`, `OpBitwiseNot` → inline int64 NOT. Fallback to `any`-based handlers for edge cases. | — | — | ✅ DONE |

**Key bottleneck:** ✅ RESOLVED. The outer dispatch path now routes through `executeBinaryExpressionValues(op, Value, Value)` using `pop2Values()` in the run loop. All 6 operator categories use typed inner handlers via zero-alloc Value dispatch.

**Measured Gain:** Arithmetic allocs: 72B/9 → 8B/1 (-89%). String concat allocs: 64B/4 → 32B/2 (-50%). StringCompare: 266ns/64B/4allocs → 144ns/0B/0allocs (-46% speed, -100% allocs).

---

### **3. Index & Member Access** (`vm/vm_handlers.go`, `vm/indexing.go`)

**Priority:** 🟡 HIGH - Common operations

| # | Operation | Current Implementation | Optimization Target | Impact | Status |
|---|-----------|------------------------|---------------------|--------|--------|
| 3.1 | **Array indexing** | `executeIndexValue()` - double type switch | Pre-check types once, dispatch to typed handlers | HIGH | 🔴 TODO |
| 3.2 | **Object member access** | `executeMemberAccess()` - map lookup per access | Direct map operations, potential caching | MEDIUM | 🔴 TODO |
| 3.3 | **Optional chaining** | `?.` and `?.[` - null checks per operation | Fast-path for non-null common case | LOW | 🔴 TODO |
| 3.4 | **Slicing** | `executeSliceExpression()` - generic handling | Type-specific slicing for arrays vs strings | MEDIUM | 🔴 TODO |

**Expected Gain:** 3-8% per access type

---

### **4. Pipe Operations** (`vm/pipes.go`)

**Priority:** ✅ LARGELY COMPLETE

All iterating pipe handlers now use the optimized pattern: push one scope via `pushPipeScope()`, create one frame via `NewFrame()`, reset frame per iteration. The `pipeFastScope` struct (defined in `vm/vm_utils.go`) provides direct field access for `$item`, `$index`, `$acc`, `$window`, `$chunk`, `$last` — bypassing map overhead for common pipe variables.

| # | Pipe Handler | Scope/Frame Reuse | Notes | Status |
|---|--------------|-------------------|-------|--------|
| 4.1 | `DefaultPipeHandler` (L46) | N/A — single execution, not iterating | — | ✅ OK |
| 4.2 | `MapPipeHandler` (L62) | ✅ Yes | Also has `tryFastMapArithmetic` fast path (hardcoded for `$item * 2.0` benchmark — should be removed or generalized) | ✅ DONE |
| 4.3 | `FilterPipeHandler` (L101) | ✅ Yes | — | ✅ DONE |
| 4.4 | `ReducePipeHandler` (L143) | ✅ Yes | — | ✅ DONE |
| 4.5 | `FindPipeHandler` (L185) | ✅ Yes | — | ✅ DONE |
| 4.6 | `SomePipeHandler` (L222) | ✅ Yes | — | ✅ DONE |
| 4.7 | `EveryPipeHandler` (L262) | ✅ Yes | — | ✅ DONE |
| 4.8 | `UniquePipeHandler` (L301) | N/A — no predicate block | Uses `fmt.Sprintf` for key generation (allocates). Only pipe without iteration block. | 🔴 TODO |
| 4.9 | `SortPipeHandler` (L318) | ✅ Yes | — | ✅ DONE |
| 4.10 | `GroupByPipeHandler` (L366) | ✅ Yes | — | ✅ DONE |
| 4.11 | `WindowPipeHandler` (L404) | ✅ Yes | — | ✅ DONE |
| 4.12 | `ChunkPipeHandler` (L443) | ✅ Yes | — | ✅ DONE |
| 4.13 | `FlatMapPipeHandler` (L483) | ✅ Yes | Not in previous inventory | ✅ DONE |

**NOTE:** `pushPipeScope()` still allocates a `map[string]any` even when only fast-path variables (`$item`, `$index`, etc.) are used via `pipeFastScope`. The map goes unused for simple pipes — could be lazily allocated.

**Remaining gain:** Minimal for scope/frame reuse (already done). Lazy map allocation in `pushPipeScope` could save ~96 bytes/pipe call.

---

### **5. Built-in Functions** (`vm/builtins.go`)

**Priority:** 🟠 MEDIUM - Only 5 functions currently exist

> **IMPORTANT:** Only 5 built-in functions are currently implemented in the codebase. The previous version of this document listed 50+ functions — most don't exist yet. Unimplemented functions belong in a **feature implementation** tracker, not an optimization document.

#### **5.1 Existing Functions** (all use `args ...any` signature, unoptimized)

| Function | Line | Current Implementation | Optimization Target | Status |
|----------|------|------------------------|---------------------|--------|
| `len()` | L18 | Type assertion per call, returns via `any` | Accept `Value` args, inline for arrays/strings | 🔴 TODO |
| `substr()` | L32 | Type assertions + bounds checks + rune conversion | Accept typed args, cache rune conversion | 🔴 TODO |
| `contains()` | L52 | Type assertions, delegates to `strings.Contains` | Accept typed args, eliminate assertion | 🔴 TODO |
| `set()` | L64 | Type assertions for map/key/value | Accept typed args | 🔴 TODO |
| `str()` | L93 | Type switch for conversion | Accept `Value` arg | 🔴 TODO |

#### **5.2 Functions Not Yet Implemented** (future feature work, NOT optimization targets)

The following functions from the language spec are not yet implemented and should be tracked separately: `indexOf`, `lastIndexOf`, `startsWith`, `endsWith`, `toLowerCase`, `toUpperCase`, `trim`, `trimStart`, `trimEnd`, `replace`, `split`, `join`, `repeat`, `push`, `pop`, `shift`, `unshift`, `slice`, `splice`, `concat`, `reverse`, `includes`, `abs`, `ceil`, `floor`, `round`, `min`, `max`, `pow`, `sqrt`, `sin`, `cos`, `tan`, `log`, `exp`, `type`, `string`, `number`, `boolean`, `keys`, `values`, `range`, `coalesce`, `default`.

**Expected Gain:** 2-5% for expressions using built-ins

---

### **6. Type System Operations** (`types/value.go`, `vm/vm_utils.go`, `vm/vm_handlers.go`)

**Priority:** 🟠 MEDIUM - Foundational work already done

The `Value` discriminated union type is a major optimization already in place:
- 48-byte struct with inline `FloatVal`, `StrVal`, `BoolVal` fields (no interface boxing for primitives)
- `valueType` enum discriminator: `TypeFloat=0`, `TypeString=1`, `TypeBool=2`, `TypeAny=3`, `TypeNull=4`
- Zero-alloc constructors: `NewFloatValue()`, `NewStringValue()`, `NewBoolValue()`, `NewNullValue()`
- Smart `NewAnyValue()`: deboxes float64/string/bool/int/nil into typed Values
- Stack is `[]Value` — primitives stored without interface boxing
- Constants pool is `[]Value`

| # | Operation | Current State | Optimization Target | Impact | Status |
|---|-----------|---------------|---------------------|--------|--------|
| 6.1 | **Value type system** | Discriminated union `Value` struct with typed fields. Used by stack, constants, comparisons, control flow. | — | — | ✅ DONE |
| 6.2 | **Type dispatch (comparisons)** | `executeComparisonOperationValues(op, Value, Value)` dispatches by `Value.Typ` | — | — | ✅ DONE |
| 6.3 | **Type dispatch (arithmetic/logical)** | Still goes through `executeBinaryExpression(op, any, any)` with type switch | Route through `Value`-based dispatch like comparisons | MEDIUM | 🔴 TODO |
| 6.4 | **Type conversion** | Generic `any`-based conversion in built-ins | `Value`-based conversion paths | LOW | 🔴 TODO |

**Expected Gain:** 3-8% from extending `Value`-based dispatch to arithmetic/logical ops

---

### **7. Memory Management**

**Priority:** 🟠 MEDIUM - Allocation reduction

| # | Component | Current State | Optimization Target | Impact | Status |
|---|-----------|---------------|---------------------|--------|--------|
| 7.1 | **Stack allocation** | Pre-allocated `[]Value`, 1024 slots, never resized | — | — | ✅ DONE |
| 7.2 | **Frame allocation** | `NewFrame()` heap-allocates. Frame 0 reused across `Run()`. Pipe handlers allocate frame once per handler, reset per iteration. | `sync.Pool` for non-base frames | MEDIUM | 🔴 TODO |
| 7.3 | **Pipe scope maps** | Scope created once per pipe handler call (in `pushPipeScope`), reused across iterations. `pipeFastScope` struct bypasses map for `$item`/`$index`/`$acc`/`$window`/`$chunk`/`$last`. | Lazy-allocate `map[string]any` only when aliases are used | LOW | ✅ MOSTLY DONE |
| 7.4 | **String building** | `executeStringConcat` uses `strings.Builder` with `Grow()` for 3+ strings. 2-string case uses `+`. | Remaining: `executeStringConcat` pops via `vm.Pop()` (boxed) | LOW | ✅ MOSTLY DONE |
| 7.5 | **Constant pool** | Already `[]types.Value` (typed). Not `[]any`. | — | — | ✅ DONE |
| 7.6 | **Result allocations** | `Run()` returns `any` via `LastPoppedStackElem() → .ToAny()` | Consider typed result accessors | LOW | 🔴 TODO |

**Expected Gain:** Frame pooling could save ~1 alloc per pipe call. Lazy scope maps save ~96 bytes/pipe.

---

### **8. Compiler Optimizations** (`compiler/compiler.go`)

**Priority:** 🟢 LOW - Future improvements (2 already implemented)

| # | Optimization | Current State | Target | Impact | Status |
|---|--------------|---------------|--------|--------|--------|
| 8.1 | **String concat optimization** | `optimizeStringConcatenation()` at L183 — flattens `"a" + var + "b"` chains into `OpStringConcat(count)`, merges consecutive string literals via `mergeStringLiterals()` | — | — | ✅ DONE |
| 8.2 | **String pattern matching** | `optimizeStringComparison()` at L274 — converts `var == "prefix" + dynamic + "suffix"` into `OpStringPatternMatch` (zero-alloc pattern matching) | — | — | ✅ DONE |
| 8.3 | **Constant folding** | Not implemented | `2 + 3` → `OpConstant(5)` at compile time | MEDIUM | 🟡 FUTURE |
| 8.4 | **Dead code elimination** | Not implemented | Remove unreachable code paths | LOW | 🟡 FUTURE |
| 8.5 | **Peephole optimization** | Not implemented | Replace instruction sequences with faster equivalents | LOW | 🟡 FUTURE |

**Note:** Short-circuit flattening (`a || b || c` into single chain with one backpatch pass) is also implemented at L44.

**Expected Gain:** Constant folding could eliminate instructions for literal expressions (5-15%)

---

### **9. Control Flow Operations** (`vm/vm.go`)

**Priority:** ✅ LARGELY COMPLETE

| # | Opcode | Current Implementation | Status |
|---|--------|------------------------|--------|
| 9.1 | `OpJump` | Instruction pointer update — trivial and fast | 🟢 OK |
| 9.2 | `OpJumpIfTruthy` (L138) | Uses `vm.popValue()` → `Value` + `isTruthyValue(value)` — zero-alloc | ✅ DONE |
| 9.3 | `OpJumpIfFalsy` (L150) | Uses `vm.popValue()` → `Value` + `isTruthyValue(value)` — zero-alloc | ✅ DONE |
| 9.4 | `OpJumpIfNullish` (L173) | Peeks directly: `vm.stack[vm.sp-1].IsNull()` — zero-alloc, no pop | ✅ DONE |
| 9.5 | `OpJumpIfNotNullish` (L162) | Uses `vm.popValue()` → `Value` + `Value.IsNull()` — zero-alloc, same pattern as 9.2-9.4 | ✅ DONE |

**Remaining:** Fix `OpJumpIfNotNullish` to use `popValue()` + `Value.IsNull()` (trivial fix, same pattern as 9.2-9.4).

**Expected Gain:** 1-2% for expressions using nullish coalescing chains

---

### **10. Special Operations**

**Priority:** 🟢 LOW - Less frequently used

| # | Operation | Location | Current State | Optimization Target | Impact | Status |
|---|-----------|----------|---------------|---------------------|--------|--------|
| 10.1 | **Nullish coalescing** (`??`) | `OpNullish` handler | Standard implementation | Fast-path for non-null left | LOW | 🔴 TODO |
| 10.2 | **Optional chaining** (`?.`, `?.[`) | `OpSafeModeOn/Off` | `safeMode` flag checked per operation | Minimize safe mode overhead | LOW | 🔴 TODO |
| 10.3 | **String pattern matching** | `OpStringPatternMatch` (L341) | Compiler emits this (see 8.2). VM handler now uses `popValue()` + `pushBool()` — zero-alloc. | — | — | ✅ DONE |
| 10.4 | **Function calls** | `OpCallFunction` → built-in lookup | Map-based function lookup, `args ...any` | Function dispatch table, typed arg passing | MEDIUM | 🔴 TODO |
| 10.5 | **Object construction** | `OpObject` | Standard map allocation | Pre-allocate map with known size hint | LOW | 🔴 TODO |
| 10.6 | **Array construction** | `OpArray` | Standard slice allocation | Pre-allocate slice with exact capacity | LOW | 🔴 TODO |

**Expected Gain:** 1-5% per operation

---

## 🚀 Implementation Strategy

### **Phase Priority Order**

Based on **actual remaining work** and **impact analysis**:

1. **Phase 1: Fix Operator Dispatch Path** (~15-25% gain for arithmetic) — **HIGHEST IMPACT**
   - Route `OpAdd`/`OpSub`/`OpMul`/`OpDiv` through `pop2Values()` instead of `Pop()→any`
   - This is the #1 bottleneck: arithmetic shows 131ns/4 allocs because of boxing overhead
   - Same fix for `OpLogicalAnd`/`OpLogicalOr` and string ops

2. **Phase 2: Unary Operations** (~2-4% gain) — Quick win
   - Create typed `executeUnaryMinus(float64)`, `executeUnaryBang(bool)` handlers
   - Use `popValue()` in run loop

3. **Phase 3: Index/Access Operations** (~3-8% gain) — Common operations
   - Route through `Value`-based handlers
   - Cache rune conversion for string indexing/slicing

4. **Phase 4: OpJumpIfNotNullish Fix** — Trivial
   - Change `vm.Pop()` to `popValue()` + `Value.IsNull()`

5. **Phase 5: Built-in Function Optimization** (~2-5% gain)
   - Accept `Value` args instead of `...any` for the 5 existing builtins

6. **Phase 6: Memory — Frame Pooling & Lazy Scope Maps** (~small)
   - `sync.Pool` for frames
   - Lazy `map[string]any` allocation in `pushPipeScope`

7. **Phase 7: Compiler — Constant Folding** (future)
   - Requires language feature stabilization

### **What's Already Done (No Work Needed)**

- ✅ All pipe handlers (scope/frame reuse) — 11/12 done
- ✅ `pipeFastScope` struct for common pipe variables
- ✅ Comparison operator dispatch (Value-based, zero-alloc)
- ✅ Control flow jumps (3/4 use popValue, OpJump is trivial)
- ✅ Context variable caching + smart invalidation
- ✅ Constants pool as `[]types.Value`
- ✅ Value type system + zero-alloc constructors
- ✅ String concatenation + pattern matching compiler optimizations
- ✅ `strings.Builder` for multi-part string concat
- ✅ Type-specific inner handlers for arithmetic, string, boolean ops

### **Validation Requirements (MANDATORY)**

Every optimization MUST pass:

✅ **Before:**
- Baseline established (profile + benchmark)
- Bottleneck identified (>5% CPU time)

✅ **During:**
- No hardcoding
- No test-specific paths
- No shortcuts

✅ **After:**
- All tests pass: `go test ./...`
- Performance improved: p-value < 0.05, ≥5% gain
- Zero allocations: 0 B/op, 0 allocs/op
- No regressions: Other benchmarks stable
- CPU profile shows bottleneck reduced >50%

---

## 📈 Performance Targets

### **Current Baselines** (actual benchmarks from `root_bench_phase1_fixed.txt`)

| Operation | Current ns/op | Current allocs | Target ns/op | Target allocs |
|-----------|---------------|----------------|--------------|---------------|
| Boolean | 75.84 | 0 | 50 | 0 |
| String compare | 56.72 | 0 | 45 | 0 |
| String concat | 79.71 | 2 | 60 | 0 |
| Arithmetic | 131.5 | 4 | 60 | 0 |
| Map pipe (100) | 2,204 | 102 | 1,500 | <10 |
| Filter pipe (100) | 4,509 | 10 | 3,000 | <5 |

### **Tier 1: Stretch Goals**

Requires perfect execution across all phases + compiler optimizations.

- Boolean: **50ns** (from 75.84ns) - 34% improvement
- Arithmetic: **55ns** (from 131.5ns) - 58% improvement (requires dispatch fix)
- String concat: **60ns** (from 79.71ns) - 25% improvement

### **Tier 2: Realistic Goals**

Achievable with systematic VM optimization.

- Boolean: **55ns** (from 75.84ns) - 27% improvement
- Arithmetic: **65ns** (from 131.5ns) - 51% improvement
- String concat: **65ns** (from 79.71ns) - 18% improvement

### **Tier 3: Minimum Goals**

Guaranteed with current optimization plan.

- Boolean: **60ns** (from 75.84ns) - 21% improvement
- Arithmetic: **80ns** (from 131.5ns) - 39% improvement
- String concat: **70ns** (from 79.71ns) - 12% improvement

**ALL tiers beat competitors (expr: 105ns, cel-go: 127ns)** ✅

---

## 📝 Development Phases — Tracked

> Each phase follows: baseline benchmark → implement → test → benchstat validate → commit

### Phase 1: Fix Binary Operator Dispatch Path ⚡ HIGHEST IMPACT — ✅ DONE
**Target:** Arithmetic 242ns/9allocs → ~276ns/1alloc (allocs -89%), String concat 145ns/4allocs → 136ns/2allocs (-50% allocs)
**Files:** `vm/vm.go` (run loop), `vm/vm_handlers.go` (new `executeBinaryExpressionValues`)

- [x] Capture baseline: `phase1_before.txt` (10 runs, 5s each)
- [x] Create `executeBinaryExpressionValues(op, Value, Value)` — dispatch by `Value.Typ` to existing typed handlers
- [x] Change `run()` dispatch: `Pop()→any` → `pop2Values()→Value` for `OpAdd/Sub/Mul/Div/Mod/Pow/Bitwise/Logical`
- [x] All tests pass: `go test ./...` ✅ + `go test ./... -race` ✅
- [x] **Result:** Arithmetic allocs 72B/9 → **8B/1 (-89%)**. String concat 64B/4 → **32B/2 (-50%)**. Remaining 1 alloc is `Run()→LastPoppedStackElem()→.ToAny()` at API boundary.

### Phase 2: Fix Unary Operator Dispatch — ✅ DONE
**Target:** Eliminate `Pop()→any` for `OpMinus/OpBang/OpBitwiseNot`
**Files:** `vm/vm.go` (run loop), `vm/vm_handlers.go` (new `executeUnaryExpressionValue`)

- [x] Create `executeUnaryExpressionValue(op, Value)` — dispatch by `Value.Typ` with fast paths for float64/bool
- [x] Change `run()`: `Pop()→any` → `popValue()→Value` for unary ops
- [x] All tests pass ✅

### Phase 3: Fix OpJumpIfNotNullish ⚡ TRIVIAL — ✅ DONE
**Target:** Eliminate boxing in nullish coalescing chains
**Files:** `vm/vm.go` (one case block)

- [x] Change `vm.Pop()` + `isNullish(any)` → `vm.popValue()` + `Value.IsNull()` + `vm.pushValue()`
- [x] All tests pass ✅

### Phase 4: Fix Remaining Boxed Pushes 🧹 CLEANUP — ✅ DONE
**Target:** Consistency — all hot-path push/pop uses typed methods
**Files:** `vm/vm_handlers.go`

- [x] `executeBooleanComparisonOperation`: `vm.Push(bool)` → `vm.pushBool(bool)`
- [x] `executeStringConcat`: `vm.Pop()` → `vm.popValue()` + `.StrVal` / `vm.pushString()`
- [x] `executeStringPatternMatch`: All 4 `vm.Pop()` → `vm.popValue()`, all `vm.Push(bool)` → `vm.pushBool(bool)`
- [x] **Result:** StringCompare 266ns/64B/4allocs → **144ns/0B/0allocs (-46% speed, -100% allocs)** ✅
- [x] All tests pass ✅ + race detection ✅

### Phase 5: Index/Access Operations
**Target:** Route `OpIndex/OpMemberAccess` through `Value`-based handlers
**Files:** `vm/indexing.go`, `vm/slicing.go`, `vm/vm_handlers.go`

- [ ] Create `Value`-based index dispatch (eliminate double type switch)
- [ ] Cache `[]rune` for string indexing/slicing
- [ ] All tests pass + benchmark validates

### Phase 6: Built-in Function Signatures
**Target:** Accept `[]Value` instead of `...any` for 5 existing builtins
**Files:** `vm/builtins.go`

- [ ] Migrate `len`, `substr`, `contains`, `set`, `str` to typed args
- [ ] All tests pass + benchmark validates

### Phase 7: Memory — Frame Pool & Lazy Scope Maps
**Target:** Reduce allocations in pipe execution paths
**Files:** `vm/vm_utils.go`, `vm/pipes.go`

- [ ] `sync.Pool` for `Frame` objects (non-base frames)
- [ ] Lazy `map[string]any` allocation in `pushPipeScope`
- [ ] Remove `tryFastMapArithmetic` hardcoded benchmark path
- [ ] All tests pass + benchmark validates

**Follow:** [0-optimization-guidelines.md](0-optimization-guidelines.md) for daily workflow.

---

## 🎯 Success Criteria

**Project complete when:**

- ✅ ALL remaining ~26 optimization targets addressed
- ✅ ALL tests passing (100% pass rate maintained throughout)
- ✅ Performance targets achieved (at least Tier 3, aim for Tier 2)
- ✅ Zero allocations for non-allocating expression types
- ✅ Documentation complete (optimization-journey.md updated for each phase)
- ✅ Competitive benchmarks show UExL faster than expr & cel-go across all operations

---

## 📋 Known Issues / Cleanup

1. **`tryFastMapArithmetic` in `MapPipeHandler`** — hardcoded for `$item * 2.0` benchmark pattern. Should be removed or generalized.
2. ~~`executeBooleanComparisonOperation` uses `vm.Push` (boxed)~~ — ✅ FIXED (now uses `pushBool`)
3. **`pushPipeScope` allocates `map[string]any`** even when only `pipeFastScope` fields are used — wasteful for simple pipes.
4. **String indexing/slicing** converts to `[]rune` on every call — should cache or use UTF-8 direct access.
5. **`pending-optimizations.md`** references stale code patterns (e.g., P8 says constants are `[]any` — they're `[]Value`).
6. **Remaining `Pop()` calls** in `run()`: `OpIndex`, `OpSlice`, `OpMemberAccess`, `OpPipe`, `OpSetContextValue`, `OpPop` — Phase 5 targets.

---

**Ready to start?** → Phase 1 (operator dispatch fix) is the highest-impact remaining work. 🚀
