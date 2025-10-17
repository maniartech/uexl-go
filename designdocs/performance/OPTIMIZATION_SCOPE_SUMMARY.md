# UExL System-Wide Optimization Scope

> **Complete Inventory of 100+ Optimization Targets**

**Last Updated:** October 17, 2025
**Status:** Planning Phase Complete - Ready for Systematic Implementation

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
| **VM Core** | 6 components | 2 ✅ | 4 🔴 | 33% |
| **Operators** | 6 categories | 1 ✅ | 5 🔴 | 17% |
| **Index/Access** | 4 operations | 0 | 4 🔴 | 0% |
| **Pipes** | 11 handlers | 1 ✅ | 10 🔴 | 9% |
| **Built-ins** | 50+ functions | 0 | 50+ 🔴 | 0% |
| **Type System** | 4 operations | 0 | 4 🔴 | 0% |
| **Memory Mgmt** | 6 components | 1 ✅ | 5 🔴 | 17% |
| **Compiler** | 5 optimizations | 0 | 5 🟡 | 0% |
| **Control Flow** | 5 opcodes | 0 | 5 🔴 | 0% |
| **Special Ops** | 6 operations | 0 | 6 🔴 | 0% |
| **TOTAL** | **100+** | **5** | **95+** | **~5%** |

### **Current Performance Baseline**

```
Boolean expressions:     62 ns/op   0 allocs   ✅ OPTIMIZED
Arithmetic operations:   ~80 ns/op  0 allocs   🔴 NOT OPTIMIZED
String operations:       ~100 ns/op 0 allocs   🔴 NOT OPTIMIZED
Pipe operations (map):   ~1000 ns/op 0 allocs  ✅ OPTIMIZED
Pipe operations (other): ~1500 ns/op 0 allocs  🔴 NOT OPTIMIZED
Array indexing:          ~50 ns/op  0 allocs   🔴 NOT OPTIMIZED
Function calls:          varies     varies     🔴 NOT OPTIMIZED
```

**Competitive Position:**
- ✅ Boolean: **41% faster** than expr (105ns), **51% faster** than cel-go (127ns)
- 🔴 Other operations: **NOT YET BENCHMARKED** against competitors

---

## 🗂️ Complete Optimization Inventory

### **1. VM Core Operations** (`vm/vm.go`)

**Priority:** 🔴 CRITICAL - Affects ALL operations

| # | Component | Current State | Optimization Target | Impact | Status |
|---|-----------|---------------|---------------------|--------|--------|
| 1.1 | **Instruction dispatch loop** | Switch-based opcode handling | Jump table or type-specialized dispatch | HIGH | 🔴 TODO |
| 1.2 | **Stack operations** | Push/Pop with bounds checking | Inline hot paths, eliminate redundant checks | HIGH | 🔴 TODO |
| 1.3 | **Frame management** | pushFrame/popFrame overhead | Frame pooling with sync.Pool | MEDIUM | 🔴 TODO |
| 1.4 | **Constant loading** | Map lookup + type assertion | Direct typed access, pre-cast constants | MEDIUM | 🔴 TODO |
| 1.5 | **Context variable caching** | Array-based cache (optimized) | - | - | ✅ DONE |
| 1.6 | **Cache invalidation** | Pointer comparison (optimized) | - | - | ✅ DONE |

**Expected Gain:** 10-20% improvement across ALL operations

---

### **2. Operator Handlers** (`vm/vm_handlers.go`)

**Priority:** 🟡 HIGH - Direct user impact

| # | Operator Category | Functions | Current Issue | Target Pattern | Impact | Status |
|---|-------------------|-----------|---------------|----------------|--------|--------|
| 2.1 | **Arithmetic** | `executeBinaryArithmeticOperation` | Accepts `any`, type assertions inside | Type-specific: `executeNumberArithmetic(op, l float64, r float64)` | HIGH | 🔴 TODO |
| 2.2 | **Comparison** | `executeNumberComparisonOperation`, `executeStringComparisonOperation`, `executeBooleanComparisonOperation` | ✅ Already type-specific | - | - | ✅ DONE |
| 2.3 | **Logical** | `executeBinaryExpression` (&&, \|\|) | Generic dispatch | Boolean-specific shortcuts | MEDIUM | 🔴 TODO |
| 2.4 | **Bitwise** | Embedded in `executeBinaryExpression` | Not separated, `any` types | `executeBitwiseOperation(op, l int64, r int64)` | LOW | 🔴 TODO |
| 2.5 | **String** | `executeStringBinaryOperation`, `executeStringConcat` | Accepts `any`, type assertions | `executeStringAddition(l string, r string)` | MEDIUM | 🔴 TODO |
| 2.6 | **Unary** | `executeUnaryMinusOperation`, `executeUnaryBangOperation` | Accepts `any`, type assertions | `executeNumberNegate(v float64)`, `executeBooleanNegate(v bool)` | LOW | 🔴 TODO |

**Expected Gain:** 5-15% per operator category

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

**Priority:** 🟡 HIGH - User-visible feature

| # | Pipe Handler | Current State | Optimization Needed | Impact | Status |
|---|--------------|---------------|---------------------|--------|--------|
| 4.1 | `MapPipeHandler` | ✅ Scope/frame reuse implemented | - | - | ✅ DONE |
| 4.2 | `FilterPipeHandler` | Creates new scope per iteration | Apply scope reuse pattern from Map | HIGH | 🔴 TODO |
| 4.3 | `ReducePipeHandler` | Creates new scope per iteration | Apply scope reuse pattern from Map | HIGH | 🔴 TODO |
| 4.4 | `FindPipeHandler` | Creates new scope per iteration | Apply scope reuse pattern from Map | MEDIUM | 🔴 TODO |
| 4.5 | `SomePipeHandler` | Creates new scope per iteration | Apply scope reuse pattern from Map | MEDIUM | 🔴 TODO |
| 4.6 | `EveryPipeHandler` | Creates new scope per iteration | Apply scope reuse pattern from Map | MEDIUM | 🔴 TODO |
| 4.7 | `UniquePipeHandler` | Standard implementation | Optimize deduplication logic | LOW | 🔴 TODO |
| 4.8 | `SortPipeHandler` | Standard implementation | Optimize comparator function calls | LOW | 🔴 TODO |
| 4.9 | `GroupByPipeHandler` | Standard implementation | Optimize key extraction & grouping | MEDIUM | 🔴 TODO |
| 4.10 | `WindowPipeHandler` | Standard implementation | Optimize window creation & iteration | LOW | 🔴 TODO |
| 4.11 | `ChunkPipeHandler` | Standard implementation | Optimize chunk allocation | LOW | 🔴 TODO |

**Expected Gain:** 15-30% for pipe operations (1500ns → 1000ns target)

---

### **5. Built-in Functions** (`vm/builtins.go`)

**Priority:** 🟠 MEDIUM - 50+ functions to optimize

#### **5.1 String Functions** (16 functions)

| Function | Current Implementation | Optimization Target | Status |
|----------|------------------------|---------------------|--------|
| `len()` | Type assertion per call | Inline or type-specific | 🔴 TODO |
| `substr()` | Type assertions + bounds checks | Pre-check types, optimize bounds | 🔴 TODO |
| `indexOf()` | Generic string search | Use stdlib optimized `strings.Index` | 🔴 TODO |
| `lastIndexOf()` | Generic string search | Use stdlib optimized `strings.LastIndex` | 🔴 TODO |
| `contains()` | Type assertions | Inline with `strings.Contains` | 🔴 TODO |
| `startsWith()` | Type assertions | Inline with `strings.HasPrefix` | 🔴 TODO |
| `endsWith()` | Type assertions | Inline with `strings.HasSuffix` | 🔴 TODO |
| `toLowerCase()` | Type assertions | Inline with `strings.ToLower` | 🔴 TODO |
| `toUpperCase()` | Type assertions | Inline with `strings.ToUpper` | 🔴 TODO |
| `trim()` | Type assertions | Inline with `strings.TrimSpace` | 🔴 TODO |
| `trimStart()` | Type assertions | Inline with `strings.TrimLeft` | 🔴 TODO |
| `trimEnd()` | Type assertions | Inline with `strings.TrimRight` | 🔴 TODO |
| `replace()` | Type assertions | Use `strings.ReplaceAll` | 🔴 TODO |
| `split()` | Type assertions + allocation | Pre-allocate result slice if possible | 🔴 TODO |
| `join()` | Type assertions + concatenation | Use `strings.Builder` | 🔴 TODO |
| `repeat()` | Type assertions | Use `strings.Repeat` | 🔴 TODO |

#### **5.2 Array Functions** (10 functions)

| Function | Current Implementation | Optimization Target | Status |
|----------|------------------------|---------------------|--------|
| `len()` | Type assertion per call | Inline or type-specific | 🔴 TODO |
| `push()` | Type assertions + append | Pre-allocate capacity if known | 🔴 TODO |
| `pop()` | Type assertions + slice | Optimize slice operations | 🔴 TODO |
| `shift()` | Type assertions + slice | Optimize slice operations | 🔴 TODO |
| `unshift()` | Type assertions + prepend | Optimize prepend pattern | 🔴 TODO |
| `slice()` | Type assertions + bounds | Pre-check bounds, type-specific | 🔴 TODO |
| `splice()` | Type assertions + complex ops | Optimize splice logic | 🔴 TODO |
| `concat()` | Type assertions + append | Pre-allocate result capacity | 🔴 TODO |
| `reverse()` | Type assertions + in-place | Optimize reversal algorithm | 🔴 TODO |
| `includes()` | Type assertions + linear search | Consider fast-path for primitives | 🔴 TODO |

#### **5.3 Math Functions** (15+ functions)

| Function | Current Implementation | Optimization Target | Status |
|----------|------------------------|---------------------|--------|
| `abs()` | Type assertion + math.Abs | Inline for simple cases | 🔴 TODO |
| `ceil()` | Type assertion + math.Ceil | Inline or type-specific | 🔴 TODO |
| `floor()` | Type assertion + math.Floor | Inline or type-specific | 🔴 TODO |
| `round()` | Type assertion + math.Round | Inline or type-specific | 🔴 TODO |
| `min()` | Type assertions + comparison | Type-specific comparison | 🔴 TODO |
| `max()` | Type assertions + comparison | Type-specific comparison | 🔴 TODO |
| `pow()` | Type assertion + math.Pow | Already inlined in VM? Check | 🔴 TODO |
| `sqrt()` | Type assertion + math.Sqrt | Inline or type-specific | 🔴 TODO |
| `sin()`, `cos()`, `tan()` | Type assertions + math funcs | Inline or type-specific | 🔴 TODO |
| `log()`, `exp()` | Type assertions + math funcs | Inline or type-specific | 🔴 TODO |

#### **5.4 Type & Utility Functions** (10+ functions)

| Function | Current Implementation | Optimization Target | Status |
|----------|------------------------|---------------------|--------|
| `type()` | Type switch | Optimize type detection | 🔴 TODO |
| `string()` | Type assertions + conversion | Type-specific conversion paths | 🔴 TODO |
| `number()` | Type assertions + conversion | Type-specific conversion paths | 🔴 TODO |
| `boolean()` | Type assertions + conversion | Type-specific conversion paths | 🔴 TODO |
| `keys()` | Map iteration + allocation | Pre-allocate result slice | 🔴 TODO |
| `values()` | Map iteration + allocation | Pre-allocate result slice | 🔴 TODO |
| `range()` | Loop + allocation | Pre-allocate exact capacity | 🔴 TODO |
| `coalesce()` | Multiple type checks | Optimize nullish checking | 🔴 TODO |
| `default()` | Type checks + fallback | Optimize fallback logic | 🔴 TODO |

**Expected Gain:** 2-10% per function category

---

### **6. Type System Operations** (`vm/vm_utils.go`, `vm/vm_handlers.go`)

**Priority:** 🟠 MEDIUM - Affects all operations

| # | Operation | Current Approach | Optimization Target | Impact | Status |
|---|-----------|------------------|---------------------|--------|--------|
| 6.1 | **Type checking** | `switch v := value.(type)` repeated | Type cache/bitmap for hot values | MEDIUM | 🔴 TODO |
| 6.2 | **Type dispatch** | Runtime type assertions per operation | Pre-computed type dispatch tables | MEDIUM | 🔴 TODO |
| 6.3 | **Type conversion** | Generic conversion functions | Type-specific conversion paths | LOW | 🔴 TODO |
| 6.4 | **Type coercion** | Accepts `any`, type switches inside | Early type resolution, typed APIs | LOW | 🔴 TODO |

**Expected Gain:** 3-8% across type-heavy operations

---

### **7. Memory Management**

**Priority:** 🟠 MEDIUM - Allocation reduction

| # | Component | Current State | Optimization Target | Impact | Status |
|---|-----------|---------------|---------------------|--------|--------|
| 7.1 | **Stack allocation** | Fixed 1024-slot array | Pre-allocated, never resized (good ✅) | - | ✅ DONE |
| 7.2 | **Frame allocation** | New frame object per scope | Frame pooling with sync.Pool | MEDIUM | 🔴 TODO |
| 7.3 | **Scope maps** | New map per pipe iteration | Reuse pattern (clear + update) | HIGH | 🔴 TODO |
| 7.4 | **String building** | Direct concatenation | strings.Builder for multi-part | LOW | 🔴 TODO |
| 7.5 | **Constant pool** | Mixed types (`[]any`) | Type-segregated pools (numbers, strings) | LOW | 🔴 TODO |
| 7.6 | **Result allocations** | Returned as `any` | Consider typed result channels | LOW | 🔴 TODO |

**Expected Gain:** 0 allocs/op maintained (already achieved), but reduce GC pressure

---

### **8. Compiler Optimizations** (`compiler/`)

**Priority:** 🟢 LOW - Future improvements

| # | Optimization | Current State | Target | Impact | Status |
|---|--------------|---------------|--------|--------|--------|
| 8.1 | **Constant folding** | No compile-time evaluation | `2 + 3` → `OpConstant(5)` | MEDIUM | 🟡 FUTURE |
| 8.2 | **Type hints** | No type information | If compiler knows types, emit specialized opcodes | HIGH | 🟡 FUTURE |
| 8.3 | **Dead code elimination** | All code compiled | Remove unreachable code paths | LOW | 🟡 FUTURE |
| 8.4 | **Instruction combining** | Each operation separate | Merge consecutive compatible ops | MEDIUM | 🟡 FUTURE |
| 8.5 | **Peephole optimization** | No pattern replacement | Replace instruction sequences with faster equivalents | LOW | 🟡 FUTURE |

**Expected Gain:** 5-15% potential (future work)

---

### **9. Control Flow Operations** (`vm/vm.go`)

**Priority:** 🟠 MEDIUM - Common in complex expressions

| # | Opcode | Current Implementation | Optimization Target | Impact | Status |
|---|--------|------------------------|---------------------|--------|--------|
| 9.1 | `OpJump` | Instruction pointer update | Already fast (inline) | - | 🟢 OK |
| 9.2 | `OpJumpIfTruthy` | Stack pop + truthiness check + jump | Fast-path for boolean true/false | MEDIUM | 🔴 TODO |
| 9.3 | `OpJumpIfFalsy` | Stack pop + truthiness check + jump | Fast-path for boolean true/false | MEDIUM | 🔴 TODO |
| 9.4 | `OpJumpIfNullish` | Stack pop + null check + jump | Fast-path for non-null | LOW | 🔴 TODO |
| 9.5 | `OpJumpIfNotNullish` | Stack pop + null check + jump | Fast-path for non-null | LOW | 🔴 TODO |

**Expected Gain:** 2-5% for expressions with short-circuit evaluation

---

### **10. Special Operations**

**Priority:** 🟢 LOW - Less frequently used

| # | Operation | Location | Optimization Target | Impact | Status |
|---|-----------|----------|---------------------|--------|--------|
| 10.1 | **Nullish coalescing** (`??`) | `OpNullish` handler | Fast-path for non-null left | LOW | 🔴 TODO |
| 10.2 | **Optional chaining** (`?.`, `?.[`) | `OpSafeModeOn/Off` | Minimize safe mode overhead | LOW | 🔴 TODO |
| 10.3 | **String pattern matching** | `OpStringPatternMatch` | Optimize prefix/suffix checks | LOW | 🔴 TODO |
| 10.4 | **Function calls** | `OpCallFunction` → built-in lookup | Function dispatch table, inline common functions | MEDIUM | 🔴 TODO |
| 10.5 | **Object construction** | `OpObject` | Pre-allocate map with known size | LOW | 🔴 TODO |
| 10.6 | **Array construction** | `OpArray` | Pre-allocate slice with exact capacity | LOW | 🔴 TODO |

**Expected Gain:** 1-5% per operation

---

## 🚀 Implementation Strategy

### **Phase Priority Order**

Based on **impact analysis** and **code dependencies**:

1. **Phase 1: Arithmetic Operations** (5-8% gain) - High user visibility
2. **Phase 2: String Operations** (3-5% gain) - Common operations
3. **Phase 3: Pipe Operations** (15-25% gain) - **HIGHEST IMPACT** ✅ Start here
4. **Phase 4: Array/Index Access** (5-7% gain) - Common operations
5. **Phase 5: Unary Operations** (2-4% gain) - Quick wins
6. **Phase 6: Boolean/Logical** (1-2% gain) - Already partially optimized
7. **Phase 7: VM Core** (10-20% gain) - **RISKY** - requires extensive testing
8. **Phase 8: Built-in Functions** (varies) - Case-by-case optimization
9. **Phase 9: Memory Management** (GC reduction) - Long-term optimization
10. **Phase 10: Compiler** (future) - Requires language feature stabilization

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

### **Tier 1: Stretch Goals (20-25ns/op)**

Requires perfect execution across all phases + compiler optimizations.

- Boolean/comparison: **20ns** (from 62ns) - 68% improvement
- Arithmetic: **22ns** (from ~80ns) - 72% improvement
- String ops: **25ns** (from ~100ns) - 75% improvement

### **Tier 2: Realistic Goals (30-35ns/op)**

Achievable with systematic VM optimization.

- Boolean/comparison: **30ns** (from 62ns) - 52% improvement
- Arithmetic: **32ns** (from ~80ns) - 60% improvement
- String ops: **35ns** (from ~100ns) - 65% improvement

### **Tier 3: Minimum Goals (35-40ns/op)**

Guaranteed with current optimization plan.

- Boolean/comparison: **35ns** (from 62ns) - 44% improvement
- Arithmetic: **38ns** (from ~80ns) - 52% improvement
- String ops: **40ns** (from ~100ns) - 60% improvement

**ALL tiers beat competitors (expr: 105ns, cel-go: 127ns)** ✅

---

## 📝 Next Steps

1. ✅ **Fix failing tests** - DONE (bitwise edge case test corrected)
2. 🔴 **Choose starting phase** - Recommend Phase 3 (Pipe Operations) or Phase 1 (Arithmetic)
3. 🔴 **Profile baseline** - Establish before-optimization metrics
4. 🔴 **Implement optimization** - Follow dos-and-donts.md patterns
5. 🔴 **Validate thoroughly** - All tests pass, benchstat confirms improvement
6. 🔴 **Document results** - Update optimization-journey.md
7. 🔴 **Repeat** - Move to next optimization target

**Follow:** [0-optimization-guidelines.md](0-optimization-guidelines.md) for daily workflow.

---

## 🎯 Success Criteria

**Project complete when:**

- ✅ ALL 100+ optimization targets addressed
- ✅ ALL tests passing (100% pass rate maintained throughout)
- ✅ Performance targets achieved (at least Tier 3, aim for Tier 2)
- ✅ Zero allocations maintained (0 allocs/op for all operations)
- ✅ Documentation complete (optimization-journey.md updated for each phase)
- ✅ Competitive benchmarks show UExL faster than expr & cel-go across all operations

**Timeline:** Estimated 2-4 weeks of focused optimization work

---

**Ready to start?** → Open [0-optimization-guidelines.md](0-optimization-guidelines.md) and begin! 🚀
