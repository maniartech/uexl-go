# Phase 2: Universal VM Performance Optimization Plan

**Status**: Planning Phase
**Goal**: Apply type-specific optimizations universally across ALL VM operations
**Target**: Eliminate ALL unnecessary allocations and interface boxing overhead
**Quality Standards**: KISS, SRP, DRY, Go best practices, maintainability, thread-safety

---

## Current Performance Baseline (Nov 13, 2025)

### Benchmark Results:
```
BenchmarkVM_Boolean_Current-16           81.69 ns/op      0 B/op    0 allocs/op  ✅
BenchmarkVM_Arithmetic_Current-16       125.7 ns/op     32 B/op    4 allocs/op  ⚠️
BenchmarkVM_String_Current-16           105.6 ns/op     32 B/op    2 allocs/op  ⚠️
BenchmarkVM_StringCompare_Current-16     60.51 ns/op     0 B/op    0 allocs/op  ✅
BenchmarkVM_Map_Current-16             2536 ns/op     2616 B/op  102 allocs/op  ⚠️
```

### Status Analysis:
- ✅ **Boolean operations**: OPTIMIZED (0 allocs)
- ✅ **String comparisons**: OPTIMIZED (0 allocs)
- ⚠️ **Arithmetic**: 4 allocations (need optimization)
- ⚠️ **String concat**: 2 allocations (need optimization)
- ⚠️ **Pipe operations**: 102 allocations (need optimization)

---

## Phase 1 Achievements (Completed)

### What Was Done:
1. ✅ VM Pool implementation
2. ✅ VM Reset optimization
3. ✅ Context variable caching (O(1) array access)
4. ✅ Type-specific push methods: `pushFloat64()`, `pushString()`, `pushBool()`
5. ✅ Optimized arithmetic operations (fast-path for common ops)
6. ✅ Optimized comparison operations (type-specific dispatch)
7. ✅ Excel compatibility features

### What Was NOT Done:
- ❌ Universal application of optimizations (only partial coverage)
- ❌ Unary operations still use generic `Push()`
- ❌ String indexing still allocates
- ❌ Array indexing returns interface{}
- ❌ String concatenation has unnecessary allocations

---

## Phase 2 Optimization Targets

### Target 1: Zero-Allocation Arithmetic (125.7 ns → < 100 ns, 4 allocs → 0 allocs)

**Current Issues**:
```go
// vm_handlers.go:226-228 - Uses generic Push()
func (vm *VM) executeUnaryMinusOperation(operand any) error {
    switch v := operand.(type) {
    case float64:
        vm.Push(-v)  // ❌ Interface boxing
    case int:
        vm.Push(float64(-v))  // ❌ Interface boxing
```

**Solution**:
```go
func (vm *VM) executeUnaryMinusOperation(operand any) error {
    switch v := operand.(type) {
    case float64:
        return vm.pushFloat64(-v)  // ✅ Zero allocations
    case int:
        return vm.pushFloat64(float64(-v))  // ✅ Zero allocations
```

**Expected Impact**: 4 allocs → 0 allocs, ~20-30 ns/op improvement

---

### Target 2: Optimized String Operations (105.6 ns → < 80 ns, 2 allocs → 0-1 allocs)

**Current Issues**:
```go
// vm_handlers.go:359 - String indexing allocates
return vm.Push(string(v[idx]))  // ❌ Creates new string allocation
```

**Solution**: Use string builder pool or byte slice reuse
```go
// For single-character strings, use a pre-allocated cache
var singleCharCache [256]string  // Cache for ASCII characters
func init() {
    for i := 0; i < 256; i++ {
        singleCharCache[i] = string(byte(i))
    }
}

func (vm *VM) executeIndexValue(target any, index any) error {
    // ... existing code ...
    case string:
        if idx < 0 || idx >= len(v) {
            return fmt.Errorf("string index out of bounds: %d", idx)
        }
        ch := v[idx]
        if ch < 256 {
            return vm.pushString(singleCharCache[ch])  // ✅ Zero allocation
        }
        return vm.pushString(string(ch))  // Rare case: non-ASCII
```

**Expected Impact**: 2 allocs → 0-1 allocs (0 for ASCII, 1 for non-ASCII)

---

### Target 3: Pipe Operations Optimization (2536 ns → < 2000 ns, 102 allocs → < 80 allocs)

**Current Issues**:
- Map allocations in pipe scope management
- Frame allocations per pipe operation
- Result array allocations

**Solution**: Implement object pooling pattern
```go
type PipeContext struct {
    scopes     []map[string]any
    resultPool []any  // Reuse result arrays
}

var pipeContextPool = sync.Pool{
    New: func() interface{} {
        return &PipeContext{
            scopes: make([]map[string]any, 0, 8),
            resultPool: make([]any, 0, 100),
        }
    },
}

func (vm *VM) executePipe(...) error {
    ctx := pipeContextPool.Get().(*PipeContext)
    defer func() {
        // Clear and return to pool
        ctx.scopes = ctx.scopes[:0]
        ctx.resultPool = ctx.resultPool[:0]
        pipeContextPool.Put(ctx)
    }()
    // Use ctx instead of allocating new maps
}
```

**Expected Impact**: 102 allocs → 60-80 allocs, ~20% speed improvement

---

## Quality Criteria & Constraints

### MUST HAVE:
1. **Zero Breaking Changes**: All existing tests must pass
2. **Thread Safety**: Use sync.Pool correctly, no data races
3. **Maintainability**: Code must remain readable and well-documented
4. **KISS Principle**: Simple solutions over complex optimizations
5. **SRP**: Each function does one thing well
6. **DRY**: No code duplication

### MUST AVOID:
1. **Premature Optimization**: Only optimize hot paths (proven by benchmarks)
2. **Over-Engineering**: No complex caching schemes without measurement
3. **Go Anti-Patterns**: No naked returns, no init() abuse, no goroutine leaks
4. **Unsafe Code**: No use of `unsafe` package
5. **Global Mutable State**: Prefer sync.Pool over global caches

### TESTING REQUIREMENTS:
1. **Benchmark Coverage**: Every optimization must have before/after benchmarks
2. **Realistic Data**: Test with real-world expression patterns
3. **Edge Cases**: Test boundary conditions (empty strings, zero, nil, etc.)
4. **Concurrency**: Run tests with `-race` flag
5. **Regression**: Verify existing tests still pass

---

## TESTING PROTOCOL (MANDATORY - NON-NEGOTIABLE)

### 🔴 CRITICAL: Zero Tolerance for Breaking Changes

**Philosophy**: "Make it work, make it right, make it fast" - IN THAT ORDER

### Pre-Implementation Requirements:

1. **Baseline Capture** (MUST DO FIRST):
   ```bash
   # Save current test results
   go test ./... -v > baseline_tests.txt
   go test ./... -race > baseline_race.txt
   go test -bench=BenchmarkVM -benchmem -benchtime=3s > baseline_bench.txt
   go test ./... -coverprofile=coverage_before.out
   ```

2. **Test Inventory**:
   - Count total tests: `go test ./... -v | grep -c "PASS"`
   - Review test coverage: `go tool cover -html=coverage_before.out`
   - Identify weak spots in coverage

3. **Document Expected Behavior**:
   - List all affected operations
   - Document edge cases for each
   - Review existing bug reports/issues

### Per-Change Testing (INCREMENTAL):

**After EVERY single code change**:

```bash
# Step 1: Unit tests for changed file
go test ./vm -v -run TestNamePattern

# Step 2: Race detector on changed package
go test ./vm -race

# Step 3: Full test suite
go test ./...

# Step 4: Race detector full suite (if vm tests pass)
go test ./... -race

# Step 5: Benchmarks for affected operations
go test -bench=BenchmarkVM_Arithmetic -benchmem -benchtime=3s

# Step 6: If ANY failure - STOP, REVERT, ANALYZE
```

### Comprehensive Testing Matrix:

| Test Level | Command | Frequency | Failure Action |
|------------|---------|-----------|----------------|
| Unit Tests | `go test ./vm -v` | After each change | **STOP & REVERT** |
| Integration | `go test ./...` | After each change | **STOP & REVERT** |
| Race (VM) | `go test ./vm -race` | After each change | **STOP & REVERT** |
| Race (All) | `go test ./... -race` | Before commit | **STOP & FIX** |
| Benchmarks | `go test -bench=BenchmarkVM` | After each change | **ANALYZE** |
| Coverage | `go test -coverprofile=...` | Before commit | **REVIEW** |
| Memory | `go test -memprofile=...` | Before commit | **PROFILE** |
| Stress | `go test ./... -count=100` | Before PR | **CRITICAL** |

### Edge Case Testing (MUST VERIFY):

**For each optimization, test**:
- ✅ Nil/null values
- ✅ Zero values (0, 0.0, "")
- ✅ Negative numbers
- ✅ Very large numbers (overflow)
- ✅ Very small numbers (underflow)
- ✅ Empty strings, arrays, objects
- ✅ Single-element collections
- ✅ Unicode/multi-byte characters
- ✅ Boundary conditions (min/max indices)
- ✅ Invalid input (type mismatches)
- ✅ Concurrent access (if applicable)

### Regression Prevention:

```bash
# Before committing ANY change:

# 1. Diff test results
diff baseline_tests.txt current_tests.txt

# 2. Compare benchmarks (must improve or explain)
benchstat baseline_bench.txt current_bench.txt

# 3. Race detector stress test (run 10 times)
for i in {1..10}; do go test ./... -race || exit 1; done

# 4. Memory leak check
go test -bench=BenchmarkVM_Arithmetic -memprofile=mem.prof
go tool pprof -alloc_space mem.prof
# Verify no unexpected allocations

# 5. Coverage comparison
go test ./... -coverprofile=coverage_after.out
go tool cover -func=coverage_after.out > coverage_after.txt
diff coverage_before.txt coverage_after.txt
# Coverage must NOT decrease
```

### Test-Driven Development (TDD) Approach:

1. **Write test FIRST** for the optimization
2. **Verify test FAILS** with current code
3. **Implement optimization**
4. **Verify test PASSES**
5. **Run full suite** (must pass)
6. **Commit** with message: "test: add test for X optimization"
7. **Implement optimization** (may already be done in step 3)
8. **Commit** with message: "perf: optimize X operation"

### Continuous Integration Checks:

- [ ] All tests pass on Go 1.21+
- [ ] Race detector clean
- [ ] No new linter warnings
- [ ] Coverage >= baseline
- [ ] Benchmarks improved or justified
- [ ] Documentation updated

---

## Implementation Phases

### Phase 2A: Universal Type-Specific Push (Week 1)
**Scope**: Replace ALL remaining `vm.Push()` calls with type-specific variants

**Tasks**:
1. ✅ Audit all `vm.Push()` calls in vm/ directory (DONE - 20 locations found)
2. 🔄 Replace unary operations with `pushFloat64()`
3. 🔄 Replace string indexing with optimized approach
4. 🔄 Replace array access where applicable
5. 🔄 Add single-character string cache

**Success Metrics**:
- Arithmetic: 4 allocs → 0 allocs
- String ops: 2 allocs → 0-1 allocs
- All existing tests pass with `-race`

**Estimated Impact**: 15-25% performance improvement on arithmetic/string ops

---

### Phase 2B: Pool-Based Resource Management (Week 2)
**Scope**: Implement sync.Pool for frequently allocated objects

**Tasks**:
1. 🔄 Implement result array pool for pipes
2. 🔄 Implement map pool for pipe scopes
3. 🔄 Implement string builder pool for concatenation
4. 🔄 Add pool metrics for monitoring

**Success Metrics**:
- Pipe ops: 102 allocs → < 80 allocs
- No memory leaks (verify with pprof)
- Thread-safe (verify with `-race`)

**Estimated Impact**: 20-30% improvement on pipe operations

---

### Phase 2C: Hot Path Micro-Optimizations (Week 3)
**Scope**: Fine-tune critical execution paths

**Tasks**:
1. 🔄 Inline small frequently-called functions
2. 🔄 Optimize switch statement ordering (most common cases first)
3. 🔄 Reduce bounds checks where safe
4. 🔄 Optimize stack operations (reduce pointer arithmetic)

**Success Metrics**:
- Overall 5-10% improvement across all benchmarks
- Code complexity remains low (cyclomatic complexity < 10)

**Estimated Impact**: 5-10% improvement across all operations

---

### Phase 2D: Validation & Documentation (Week 4)
**Scope**: Ensure quality and prepare for production

**Tasks**:
1. 🔄 Run full benchmark suite (compare with Phase 1)
2. 🔄 Profile memory allocations (pprof)
3. 🔄 Run race detector on all tests
4. 🔄 Update documentation
5. 🔄 Create performance comparison report

**Success Metrics**:
- All tests pass (including `-race`)
- Performance improved across the board
- Documentation up-to-date
- No memory leaks

---

## Expected Final Results

### Performance Targets:
```
Operation              Current    Target     Improvement
------------------------------------------------------
Boolean               81.69 ns    < 70 ns      15%
Arithmetic           125.7 ns    < 100 ns     20%
String Concat        105.6 ns    < 80 ns      25%
String Compare        60.51 ns   < 55 ns      10%
Map (pipe)           2536 ns    < 2000 ns     20%
```

### Allocation Targets:
```
Operation              Current    Target     Improvement
------------------------------------------------------
Boolean                0 allocs   0 allocs     -
Arithmetic             4 allocs   0 allocs   100%
String Concat          2 allocs   0-1 allocs  50-100%
String Compare         0 allocs   0 allocs     -
Map (pipe)           102 allocs  60-80 allocs 20-40%
```

---

## Risk Assessment

### Low Risk (Safe to Implement):
- ✅ Type-specific push methods (already proven)
- ✅ Unary operation optimization (simple change)
- ✅ Single-character string cache (static, no concurrency issues)

### Medium Risk (Needs Testing):
- ⚠️ sync.Pool for pipe contexts (needs race testing)
- ⚠️ String builder pool (needs proper reset)
- ⚠️ Result array reuse (needs careful bounds checking)

### High Risk (Avoid for Now):
- ❌ Unsafe pointer operations (violates quality criteria)
- ❌ JIT compilation (too complex, future phase)
- ❌ Assembly optimizations (not maintainable)

---

## Go Best Practices Applied

### 1. Effective Use of sync.Pool:
```go
// ✅ CORRECT: Stateless pool usage
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func useBuffer() string {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()  // Clean before returning
        bufferPool.Put(buf)
    }()
    // Use buf...
}
```

### 2. Avoiding Premature Optimization:
```go
// ✅ CORRECT: Optimize proven hot paths only
// Run benchmark first, optimize if needed

// ❌ WRONG: Optimizing without measurement
// Don't add complexity without proof
```

### 3. Clear Error Handling:
```go
// ✅ CORRECT: Clear error paths
if err := vm.pushFloat64(result); err != nil {
    return fmt.Errorf("push float64: %w", err)
}

// ❌ WRONG: Swallowing errors
vm.pushFloat64(result)  // Ignores error
```

### 4. Readable Code Over Clever Code:
```go
// ✅ CORRECT: Clear intent
if ch < 256 {
    return vm.pushString(singleCharCache[ch])
}
return vm.pushString(string(ch))

// ❌ WRONG: Clever but obscure
return vm.pushString([]string{singleCharCache[ch], string(ch)}[btoi(ch >= 256)])
```

---

## Implementation Checklist

### Before Starting Any Phase:
- [ ] Create feature branch: `phase2/<phase-name>`
- [ ] Run full test suite to establish baseline (`go test ./... -v > baseline.txt`)
- [ ] Run race detector to establish clean state (`go test ./... -race > race_baseline.txt`)
- [ ] Run benchmarks to establish baseline (`go test -bench=. -benchmem > bench_baseline.txt`)
- [ ] Generate coverage report (`go test ./... -coverprofile=coverage_before.out`)
- [ ] Document current behavior and edge cases
- [ ] Review existing tests for affected functionality
- [ ] Identify potential breakage points

### During Implementation:
- [ ] Write tests FIRST (TDD approach)
- [ ] Implement changes ONE AT A TIME (single responsibility)
- [ ] After EACH change:
  - [ ] Run unit tests: `go test ./vm -v`
  - [ ] Run race detector: `go test ./vm -race`
  - [ ] Run affected benchmarks
  - [ ] Git commit with clear message
- [ ] Profile to verify improvement (not degradation)
- [ ] Document code changes inline
- [ ] If ANY test fails: STOP, revert, analyze
- [ ] Never commit broken code

### Before Merging (STRICT CHECKLIST):
- [ ] **ALL tests pass**: `go test ./... -v` (no skips, no failures)
- [ ] **Race detector clean**: `go test ./... -race -count=10` (run 10 times)
- [ ] **Benchmarks verified**: Performance improved OR justified if same
  - [ ] Compare: `benchstat baseline_bench.txt current_bench.txt`
  - [ ] Document any regressions with explanation
- [ ] **Coverage maintained or improved**: Compare coverage reports
- [ ] **Memory profiling**: No new leaks detected
  - [ ] `go test -bench=. -memprofile=mem.prof`
  - [ ] `go tool pprof -alloc_space mem.prof`
- [ ] **Code quality checks**:
  - [ ] No `TODO` or `FIXME` comments without tickets
  - [ ] No commented-out code
  - [ ] All functions documented
  - [ ] No panics in production code paths
- [ ] **Integration testing**:
  - [ ] Test with realistic expressions (see test suite)
  - [ ] Test edge cases (nil, empty, zero, negative)
  - [ ] Test error paths (invalid input, out of bounds)
- [ ] **Code reviewed**: At least one other developer
- [ ] **Documentation updated**: README, inline comments, CHANGELOG
- [ ] **Git hygiene**: Clean commit history, meaningful messages

---

## Monitoring & Validation

### Benchmark Command:
```bash
# Run all benchmarks with memory profiling
go test -bench=. -benchmem -benchtime=3s -memprofile=mem.prof -cpuprofile=cpu.prof

# Compare with baseline
benchstat phase1_baseline.txt phase2_current.txt
```

### Race Detection:
```bash
# Run with race detector
go test ./... -race -count=10
```

### Memory Profiling:
```bash
# Generate memory profile
go test -bench=BenchmarkVM_Arithmetic -memprofile=mem.prof
go tool pprof -alloc_space mem.prof
```

---

## Success Definition

Phase 2 is considered **SUCCESSFUL** when:

1. ✅ **Performance**: 15-25% improvement across core operations
2. ✅ **Allocations**: 50%+ reduction in allocation counts
3. ✅ **Quality**: All tests pass including race detector
4. ✅ **Maintainability**: Code complexity remains low
5. ✅ **Documentation**: All changes documented
6. ✅ **Stability**: No regressions in existing functionality

Phase 2 is considered **COMPLETE** when:
- All phases (2A-2D) finished
- Performance targets met
- Quality criteria satisfied
- Production-ready (no known issues)

---

## Notes

- This plan follows the principle of **incremental improvement**
- Each phase is **independently valuable** (can ship after any phase)
- Focus on **measured improvements** (no speculation)
- Maintain **code quality** throughout (no technical debt)
- **Thread safety** is non-negotiable (UExL must be concurrent-safe)

---

**Last Updated**: November 13, 2025
**Status**: Ready for Phase 2A implementation
**Next Action**: Begin Phase 2A - Universal Type-Specific Push optimization
