# Phase 2A: Testing & Validation Checklist

**Purpose**: Ensure ZERO breaking changes during optimization
**Philosophy**: "Make it work, make it right, make it fast" - IN THAT ORDER

---

## 🔴 MANDATORY Pre-Work (DO THIS FIRST)

### Step 1: Capture Baseline

```bash
cd /e/Projects/uexl/uexl-go

# Test baseline
go test ./... -v | tee phase2a_baseline_tests.txt
echo "Total tests:" $(grep -c "PASS:" phase2a_baseline_tests.txt)

# Race baseline
go test ./... -race 2>&1 | tee phase2a_baseline_race.txt

# Benchmark baseline
go test -bench=BenchmarkVM -benchmem -benchtime=3s | tee phase2a_baseline_bench.txt

# Coverage baseline
go test ./... -coverprofile=phase2a_coverage_before.out
go tool cover -func=phase2a_coverage_before.out > phase2a_coverage_before.txt
```

### Step 2: Verify Clean State

```bash
# Must show zero failures
grep "FAIL" phase2a_baseline_tests.txt
# Output should be empty

# Must show zero races
grep "DATA RACE" phase2a_baseline_race.txt
# Output should be empty

# Record test count
BASELINE_TEST_COUNT=$(grep -c "PASS:" phase2a_baseline_tests.txt)
echo "Baseline: $BASELINE_TEST_COUNT tests passing"
```

### Step 3: Create Feature Branch

```bash
git status  # Must be clean
git add -A
git commit -m "chore: checkpoint before Phase 2A optimization"
git checkout -b phase2a-universal-push-optimization
```

---

## 📋 Per-Change Testing Protocol

### After EVERY Code Change (No Exceptions):

```bash
# 1. Quick Smoke Test (30 seconds)
go test ./vm -v
# MUST pass - if not, STOP and debug

# 2. Race Detector (VM only, 1 minute)
go test ./vm -race
# MUST pass - if not, STOP and debug

# 3. Full Test Suite (2 minutes)
go test ./...
# MUST pass - if not, STOP and debug

# 4. Full Race Detector (3 minutes)
go test ./... -race
# MUST pass - if not, STOP and debug

# 5. Affected Benchmark (30 seconds)
go test -bench=BenchmarkVM_Arithmetic -benchmem
# Should show improvement or stay same

# 6. Git Commit (if all pass)
git add vm/vm_handlers.go  # Or whichever file changed
git commit -m "perf(vm): optimize unary minus operation"
```

**If ANY test fails**:
1. DO NOT proceed
2. DO NOT commit
3. Analyze the failure
4. Fix or revert
5. Re-run all tests
6. Only proceed when 100% clean

---

## 🧪 Comprehensive Test Matrix

### Level 1: Unit Tests (VM Package)

```bash
# All VM tests
go test ./vm -v

# Specific test patterns
go test ./vm -v -run TestBinaryExpression
go test ./vm -v -run TestUnary
go test ./vm -v -run TestArithmetic
go test ./vm -v -run TestString
go test ./vm -v -run TestIndex
go test ./vm -v -run TestSlice

# White-box tests
go test ./vm -v -run Test.*_wb
```

**Expected**: 100% pass, zero failures

### Level 2: Integration Tests (All Packages)

```bash
# Full suite
go test ./...

# Verbose mode
go test ./... -v

# With timeout
go test ./... -timeout=5m
```

**Expected**: All packages pass

### Level 3: Race Detection

```bash
# VM package
go test ./vm -race

# Full suite
go test ./... -race

# Specific operations
go test ./vm -race -run TestConcurrent
```

**Expected**: Zero data races detected

### Level 4: Stress Testing

```bash
# Run tests 100 times
go test ./vm -count=100

# Run race detector 10 times
for i in {1..10}; do
  echo "Iteration $i"
  go test ./vm -race || break
done

# Parallel execution
go test ./... -parallel=8
```

**Expected**: All iterations pass

### Level 5: Benchmarks

```bash
# Core benchmarks
go test -bench=BenchmarkVM -benchmem -benchtime=3s

# Specific operation
go test -bench=BenchmarkVM_Arithmetic -benchmem -benchtime=5s

# With CPU profile
go test -bench=BenchmarkVM_Arithmetic -cpuprofile=cpu.prof

# With memory profile
go test -bench=BenchmarkVM_Arithmetic -memprofile=mem.prof
```

**Expected**: Performance improvement or same

### Level 6: Coverage Analysis

```bash
# Generate coverage
go test ./... -coverprofile=coverage.out

# View as HTML
go tool cover -html=coverage.out -o coverage.html

# View as text
go tool cover -func=coverage.out

# Coverage by package
go test ./vm -coverprofile=vm_coverage.out
go tool cover -func=vm_coverage.out
```

**Expected**: Coverage >= baseline

### Level 7: Memory Profiling

```bash
# Generate memory profile
go test -bench=BenchmarkVM_Arithmetic -memprofile=mem.prof -benchtime=10s

# Analyze allocations
go tool pprof -alloc_space mem.prof

# Interactive analysis
go tool pprof mem.prof
# Commands: top, list, web
```

**Expected**: Allocations reduced, no leaks

---

## 🎯 Edge Case Testing

### Test Each Optimization With:

```bash
# Create test file: vm/phase2a_edge_cases_test.go
```

```go
package vm

import "testing"

func TestUnaryMinus_EdgeCases(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  float64
    }{
        {"zero", "-0", 0},
        {"positive", "-5", -5},
        {"negative", "-(-5)", 5},
        {"large", "-999999999", -999999999},
        {"small", "-0.0001", -0.0001},
        {"scientific", "-1e10", -1e10},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Run expression and verify
            result, err := Eval(tt.input, nil)
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if result.(float64) != tt.want {
                t.Errorf("got %v, want %v", result, tt.want)
            }
        })
    }
}

func TestStringIndexing_EdgeCases(t *testing.T) {
    tests := []struct {
        name  string
        expr  string
        ctx   map[string]any
        want  string
    }{
        {"ascii", `"hello"[0]`, nil, "h"},
        {"unicode", `"😀🎉"[0]`, nil, "😀"},
        {"empty", `""[0]`, nil, ""}, // Should error
        {"negative", `"hello"[-1]`, nil, "o"},
    }
    // ... similar pattern
}

func TestStringConcat_EdgeCases(t *testing.T) {
    tests := []struct {
        name  string
        expr  string
        want  string
    }{
        {"empty", `"" + ""`, ""},
        {"ascii", `"hello" + " world"`, "hello world"},
        {"unicode", `"😀" + "🎉"`, "😀🎉"},
        {"numbers", `"test" + "123"`, "test123"},
    }
    // ... similar pattern
}
```

**Run edge cases**:
```bash
go test ./vm -v -run EdgeCases
```

---

## ✅ Final Validation Checklist

### Before Creating PR:

- [ ] **All Tests Pass**
  ```bash
  go test ./... -v | tee final_tests.txt
  diff phase2a_baseline_tests.txt final_tests.txt
  ```

- [ ] **Race Detector Clean (10x)**
  ```bash
  for i in {1..10}; do
    echo "Race test iteration $i"
    go test ./... -race || exit 1
  done
  ```

- [ ] **Stress Test (100x)**
  ```bash
  go test ./vm -count=100
  ```

- [ ] **Benchmarks Improved**
  ```bash
  go test -bench=BenchmarkVM -benchmem -benchtime=3s > final_bench.txt
  benchstat phase2a_baseline_bench.txt final_bench.txt
  ```

- [ ] **Coverage Maintained**
  ```bash
  go test ./... -coverprofile=final_coverage.out
  go tool cover -func=final_coverage.out > final_coverage.txt
  diff phase2a_coverage_before.txt final_coverage.txt
  ```

- [ ] **No Memory Leaks**
  ```bash
  go test -bench=BenchmarkVM_Arithmetic -memprofile=final_mem.prof
  go tool pprof -alloc_space final_mem.prof
  # Verify allocation count matches expectation
  ```

- [ ] **Code Review**
  - [ ] All changes have inline comments
  - [ ] No debug code left (fmt.Println, etc.)
  - [ ] No TODO comments without tickets
  - [ ] Consistent with existing style
  - [ ] No panics in production code

- [ ] **Documentation**
  - [ ] Inline comments for complex logic
  - [ ] Commit messages follow convention
  - [ ] Update CHANGELOG if needed

---

## 🚨 Emergency Rollback Procedures

### If Tests Fail After Changes:

```bash
# Option 1: Revert specific file
git status
git diff vm/vm_handlers.go  # Review changes
git checkout -- vm/vm_handlers.go  # Revert

# Option 2: Revert to last commit
git reset --hard HEAD

# Option 3: Revert to baseline
git reset --hard <commit-hash-before-phase2a>

# Option 4: Nuclear - abandon branch
git checkout master
git branch -D phase2a-universal-push-optimization
```

### If Race Detected:

```bash
# Get detailed race report
go test ./vm -race -v > race_report.txt 2>&1

# Analyze the race
cat race_report.txt | grep "WARNING: DATA RACE"

# Common causes in our changes:
# 1. Shared string cache without sync
# 2. Concurrent map access in pipe scopes
# 3. Stack manipulation race

# Fix approach:
# 1. Identify the shared resource
# 2. Add proper synchronization (mutex, channels, atomic)
# 3. OR redesign to avoid sharing
```

### If Performance Regresses:

```bash
# Profile to find regression
go test -bench=BenchmarkVM_Arithmetic -cpuprofile=cpu.prof
go tool pprof cpu.prof
# Commands: top, list <function>, web

# Common issues:
# 1. Added unnecessary type assertions
# 2. Introduced new allocations
# 3. Added locking overhead
```

---

## 📊 Success Metrics

### Required Outcomes:

| Metric | Baseline | Target | Actual | Pass? |
|--------|----------|--------|--------|-------|
| Arithmetic allocs | 4 | 0 | ___ | ☐ |
| String allocs | 2 | 0-1 | ___ | ☐ |
| All tests pass | 100% | 100% | ___ | ☐ |
| Race detector | 0 races | 0 races | ___ | ☐ |
| Coverage % | __% | ≥__% | ___ | ☐ |
| Benchmark time | 125.7 ns | <110 ns | ___ | ☐ |

### Sign-Off:

- [ ] All metrics meet or exceed targets
- [ ] No known issues or bugs
- [ ] Ready for code review
- [ ] Ready for merge to master

---

**Last Updated**: November 13, 2025
**Status**: Ready for rigorous implementation
**Risk Level**: LOW (with this testing protocol)
