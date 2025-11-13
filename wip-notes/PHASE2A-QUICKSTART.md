# Phase 2 Quick Start Guide

**Start Date**: November 13, 2025
**Status**: Ready to begin Phase 2A

## Pre-Implementation Baseline (VERIFIED)

### ✅ All Tests Pass
```bash
$ go test ./... -race
ok      github.com/maniartech/uexl_go
ok      github.com/maniartech/uexl_go/compiler
ok      github.com/maniartech/uexl_go/parser
ok      github.com/maniartech/uexl_go/vm
```

### 📊 Performance Baseline
```
BenchmarkVM_Boolean_Current-16           81.69 ns/op      0 B/op    0 allocs/op
BenchmarkVM_Arithmetic_Current-16       125.7 ns/op     32 B/op    4 allocs/op
BenchmarkVM_String_Current-16           105.6 ns/op     32 B/op    2 allocs/op
BenchmarkVM_StringCompare_Current-16     60.51 ns/op     0 B/op    0 allocs/op
BenchmarkVM_Map_Current-16             2536 ns/op     2616 B/op  102 allocs/op
```

---

## 🔴 CRITICAL: Safety-First Implementation Protocol

### MANDATORY Pre-Work (DO NOT SKIP):

```bash
# 1. Capture baseline - MUST DO FIRST
cd /e/Projects/uexl/uexl-go
go test ./... -v > phase2a_baseline_tests.txt
go test ./... -race > phase2a_baseline_race.txt
go test -bench=BenchmarkVM -benchmem -benchtime=3s > phase2a_baseline_bench.txt
go test ./... -coverprofile=phase2a_coverage_before.out

# 2. Verify clean state
grep -c "PASS" phase2a_baseline_tests.txt  # Should show test count
grep "FAIL" phase2a_baseline_tests.txt     # Should be empty
grep "DATA RACE" phase2a_baseline_race.txt # Should be empty

# 3. Save current working state
git add -A
git commit -m "chore: checkpoint before Phase 2A"
git checkout -b phase2a-universal-push-optimization
```

### Testing After EACH Code Change:

**🚨 NEVER skip this - it catches bugs immediately**

```bash
# Quick test (30 seconds)
go test ./vm -v
go test ./vm -race

# If above passes, full test (2 minutes)
go test ./...
go test ./... -race

# If above passes, benchmark affected operation
go test -bench=BenchmarkVM_Arithmetic -benchmem

# If ANY test fails:
# 1. STOP immediately
# 2. DO NOT proceed to next change
# 3. Review the failure
# 4. Fix or revert
# 5. Re-run tests until clean
```

### Rollback Plan:

```bash
# If something goes wrong:
git status                    # See what changed
git diff vm/vm_handlers.go    # Review changes
git checkout vm/vm_handlers.go # Revert single file
git reset --hard HEAD          # Nuclear option - revert all
```

---

## Phase 2A: Complete Coverage Audit

### Missing Type-Specific Optimizations Found:

**vm/vm_handlers.go**:
- ❌ Line 226-228: `executeUnaryMinusOperation` - uses `Push()` instead of `pushFloat64()`
- ❌ Line 359: String indexing - allocates new string
- ❌ Line 449, 459: String concatenation - uses generic `Push()`
- ❌ Line 487: String builder result - uses generic `Push()`
- ❌ Line 154: Generic string binary operation - uses `Push()`

**vm/indexing.go**:
- ❌ Line 11: Push nil (optional indexing)
- ❌ Line 49: Array indexing result
- ❌ Line 70: Object key result
- ❌ Line 94: String indexing - allocates string

**vm/slicing.go**:
- ❌ Line 11, 60, 67, 74: Array slicing results
- ❌ Line 112, 119, 126: String slicing results

**vm/vm.go**:
- ❌ Line 79: Constants (mixed types - keep as-is)
- ❌ Line 215: OpNull (keep as-is)
- ❌ Line 223: OpArray (keep as-is)
- ❌ Line 234: OpObject (keep as-is)

**Total Impact**: ~20+ locations need optimization

## Phase 2A: Implementation Checklist

### Step 1: Optimize Unary Minus Operation ⭐ START HERE
**File**: `vm/vm_handlers.go` (lines 224-231)

**Current Code**:
```go
func (vm *VM) executeUnaryMinusOperation(operand any) error {
    switch v := operand.(type) {
    case float64:
        vm.Push(-v)  // ❌ Allocates
    case int:
        vm.Push(float64(-v))  // ❌ Allocates
    default:
        return fmt.Errorf("unknown operand type: %T", operand)
    }
    return nil
}
```

**Target Code**:
```go
func (vm *VM) executeUnaryMinusOperation(operand any) error {
    switch v := operand.(type) {
    case float64:
        return vm.pushFloat64(-v)  // ✅ Zero allocations
    case int:
        return vm.pushFloat64(float64(-v))  // ✅ Zero allocations
    default:
        return fmt.Errorf("unknown operand type: %T", operand)
    }
}
```

**Test Command**:
```bash
# BEFORE implementing:
go test ./vm -v -run TestUnary  # See current test coverage

# AFTER implementing:
go test ./vm -v -run TestUnary  # Must still pass
go test ./vm -race              # Must be race-free
go test -bench=BenchmarkVM_Arithmetic -benchmem -benchtime=3s

# Expected results:
# - All tests pass ✅
# - 4 allocs → 2 allocs (50% reduction immediately)
# - OR 4 allocs → 0 allocs (if string concat also fixed)

# Test edge cases manually:
# go test ./vm -v -run TestBinaryExpression
# go test ./vm -v -run TestArithmetic
```

**Validation Checklist**:
- [ ] Unit tests pass: `go test ./vm -v`
- [ ] Race detector clean: `go test ./vm -race`
- [ ] Full suite passes: `go test ./...`
- [ ] Benchmark shows improvement
- [ ] No new allocations introduced
- [ ] Edge cases tested (negative numbers, zero, overflow)

---

### Step 2: Optimize String Concatenation Operations
**Files**: `vm/vm_handlers.go` (lines 449, 459, 487)

**Current Code** (line 449, 459):
```go
return vm.Push(leftStr + rightStr)  // ❌ Interface boxing
```

**Target Code**:
```go
return vm.pushString(leftStr + rightStr)  // ✅ Zero allocation overhead
```

**Current Code** (line 487 - string builder):
```go
return vm.Push(builder.String())  // ❌ Interface boxing
```

**Target Code**:
```go
return vm.pushString(builder.String())  // ✅ Zero allocation overhead
```

**Impact**: Reduces allocations in string concatenation operations

---

### Step 3: Optimize String Indexing (Multiple Files)
**Files**: `vm/vm_handlers.go` (line 359), `vm/indexing.go` (line 94)

**Add Single-Char Cache** (in `vm/vm_utils.go`):
```go
// Pre-allocated single-character strings for ASCII
var singleCharCache [256]string

func init() {
    for i := 0; i < 256; i++ {
        singleCharCache[i] = string(byte(i))
    }
}
```

**Update String Indexing in vm/vm_handlers.go** (line 359):
```go
case string:
    if idx < 0 || idx >= len(v) {
        return fmt.Errorf("string index out of bounds: %d", idx)
    }
    ch := v[idx]
    if ch < 256 {
        return vm.pushString(singleCharCache[ch])  // ✅ Zero allocation for ASCII
    }
    return vm.pushString(string(ch))  // Rare case: non-ASCII
```

**Update String Indexing in vm/indexing.go** (line 94):
```go
func (vm *VM) executeStringIndex(str string, index any) error {
    // ... existing validation ...

    rune := runes[intIdx]
    if rune < 256 {
        return vm.pushString(singleCharCache[rune])  // ✅ Zero allocation for ASCII
    }
    return vm.pushString(string(rune))  // Non-ASCII
}
```

**Test Command**:
```bash
go test -bench=BenchmarkVM_String -benchmem -benchtime=3s
# Expected: 2 allocs → 0-1 allocs
```

---

### Step 4: Optimize Slicing Operations
**Files**: `vm/slicing.go` (multiple locations)

**Empty Array/String Results**:
```go
// Lines 60, 67 - Empty array results
return vm.Push([]any{})  // Keep as-is (empty slice is cheap)

// Lines 112, 119 - Empty string results
return vm.pushString("")  // ✅ Change to pushString
```

**Non-Empty Results**:
```go
// Line 74 - Array slice result
return vm.Push(result)  // Keep as-is ([]any already optimized)

// Line 126 - String slice result
return vm.pushString(string(result))  // ✅ Change to pushString
```

**Impact**: Minor - slicing is less frequent, but good for consistency

---

### Step 5: Optimize Generic String Binary Operations
**File**: `vm/vm_handlers.go` (line 154)

**Current Code**:
```go
func (vm *VM) executeStringBinaryOperation(operator code.Opcode, left, right any) error {
    switch operator {
    case code.OpAdd:
        l, lok := left.(string)
        r, rok := right.(string)
        if !lok || !rok {
            return fmt.Errorf("string addition requires string operands")
        }
        result := l + r
        return vm.Push(result)  // ❌ Interface boxing
    // ...
}
```

**Target Code**:
```go
        result := l + r
        return vm.pushString(result)  // ✅ Zero allocation overhead
```

**Note**: This function might be redundant with the optimized `executeStringAddition` - needs review

---

### Step 6: Review and Consolidate
**File**: `vm/vm_handlers.go`

**Potential Redundancy**:
- `executeStringAddition` (line 136-140) - optimized version
- `executeStringBinaryOperation` (line 143-157) - generic version with `OpAdd` case

**Action**: Verify if `executeStringBinaryOperation` is still called, or if all paths use `executeStringAddition`

---

### Step 7: Verify All Changes
```bash
# Run all tests with race detector
go test ./... -race

# Run benchmarks
go test -bench=BenchmarkVM -benchmem -benchtime=3s

# Compare results
# Expected outcomes:
# - Arithmetic: 125.7 ns → ~90-100 ns, 4 allocs → 0 allocs
# - String: 105.6 ns → ~80-90 ns, 2 allocs → 0 allocs
# - String operations across the board: reduced allocations
```

---

## Complete Optimization Summary

### Files to Modify:

1. **vm/vm_utils.go** - Add single-character string cache
2. **vm/vm_handlers.go** - 5 optimizations:
   - executeUnaryMinusOperation (line 226-228)
   - String concatenation (lines 449, 459)
   - String builder result (line 487)
   - String binary operation (line 154)
3. **vm/indexing.go** - 1 optimization:
   - executeStringIndex (line 94)
4. **vm/slicing.go** - 2 optimizations:
   - Empty string results (lines 112, 119)
   - String slice result (line 126)

### Optimization Categories:

| Category | Locations | Expected Impact |
|----------|-----------|-----------------|
| Unary operations | 1 | High (4→0 allocs) |
| String concat | 3 | Medium (reduce boxing) |
| String indexing | 2 | Medium (0-1 allocs) |
| String slicing | 3 | Low (consistency) |
| **TOTAL** | **9 changes** | **20-30% improvement** |

---

## Implementation Order (Revised Priority)

1. **Unary Minus** (vm_handlers.go:226-228) - Highest impact ⭐
2. **String Concatenation** (vm_handlers.go:449,459,487,154) - High impact
3. **Single-Char Cache** (vm_utils.go) - Add infrastructure
4. **String Indexing** (vm_handlers.go:359, indexing.go:94) - Medium impact
5. **String Slicing** (slicing.go:112,119,126) - Low impact, completeness
6. **Verify & Test** - Ensure no regressions

---

## Success Criteria (STRICT - ALL MUST PASS)

### 🔴 CRITICAL (Must Pass - No Exceptions):
- [ ] **Zero Test Failures**: `go test ./...` returns 100% pass
- [ ] **Zero Race Conditions**: `go test ./... -race` completely clean
- [ ] **Zero Regressions**: All existing behavior identical
- [ ] **Stress Test Passes**: `go test ./vm -count=100` (100 iterations)

### ⚠️ REQUIRED (Performance Targets):
- [ ] **Arithmetic**: 4 allocs → 0 allocs (100% reduction)
- [ ] **String ops**: 2 allocs → 0-1 allocs (50-100% reduction)
- [ ] **Performance**: 15-20% improvement in affected benchmarks
- [ ] **No Degradation**: No benchmark regresses >2%

### ✅ VALIDATION (Quality Gates):
- [ ] **Code Quality**: Maintains readability and simplicity
- [ ] **Coverage**: No decrease in test coverage %
- [ ] **Documentation**: All changes have inline comments
- [ ] **Edge Cases Tested**:
  - [ ] Negative numbers work correctly
  - [ ] Zero values handled
  - [ ] Nil/null handled
  - [ ] Empty strings work
  - [ ] Unicode characters work (for string ops)
  - [ ] Large numbers (no overflow)

### 🚨 ABORT CONDITIONS (Stop Immediately If):
- ❌ ANY test fails
- ❌ Race detector shows ANY race
- ❌ Benchmark regresses >5%
- ❌ Coverage decreases
- ❌ You don't understand a change you made

**Recovery**: `git reset --hard HEAD` (revert all changes)

## Commands Reference

### Test & Benchmark:
```bash
# Quick test
go test ./vm -v

# With race detector
go test ./vm -race

# Specific benchmark
go test -bench=BenchmarkVM_Arithmetic -benchmem -benchtime=3s

# All benchmarks
go test -bench=BenchmarkVM -benchmem -benchtime=3s

# Save baseline
go test -bench=BenchmarkVM -benchmem -benchtime=3s > phase2a_baseline.txt
```

### After Implementation:
```bash
# Compare before/after
go test -bench=BenchmarkVM -benchmem -benchtime=3s > phase2a_after.txt
benchstat phase2a_baseline.txt phase2a_after.txt
```

## Next Steps After Phase 2A

Once Phase 2A is complete:
1. Document results in `wip-notes/README.md`
2. Update baseline benchmarks
3. Proceed to Phase 2B (Pool-Based Resource Management)

---

**Ready to Start**: Yes ✅
**First Task**: Optimize `executeUnaryMinusOperation` in `vm/vm_handlers.go`
**Expected Time**: 30-60 minutes for Phase 2A
**Risk Level**: Low (simple, proven optimizations)
