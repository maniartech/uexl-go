# 🎯 UExL Performance Optimization Guidelines

> **START HERE** - Your primary guide for implementing **SYSTEM-WIDE** performance optimizations

**Status:** 🚀 Active Development - **COMPREHENSIVE OPTIMIZATION EFFORT**
**Current Performance:** 62ns/op (boolean expressions)
**Target:** 20ns stretch goal, 30-35ns realistic
**Scope:** **ENTIRE UExL PIPELINE** - Parser, Compiler, VM, All Operators, All Functions, All Expressions
**Last Updated:** October 17, 2025

---

## 🎯 Mission: Military-Grade Performance Everywhere

This is **NOT** just about optimizing boolean expressions or strings. This is a **COMPLETE SYSTEM-WIDE PERFORMANCE OVERHAUL** covering:

### **Every Component in the Pipeline:**

- ✅ **VM Core:** Stack operations, frame management, instruction dispatch, opcode handling
- ✅ **Context Handling:** Variable lookup, caching, scope management (partially done)
- ✅ **All Operators:** Arithmetic, comparison (done), logical, bitwise, unary, string, array
- ✅ **All Expressions:** Binary, unary, index, member access, function calls, pipes
- ✅ **Pipe Operations:** Map (done), filter, reduce, find, some, every, unique, sort, groupBy, window, chunk
- ✅ **Built-in Functions:** ALL functions in `vm/builtins.go` (50+ functions)
- ✅ **Type Operations:** Conversions, coercions, type checking, dispatch
- ✅ **Memory Management:** Stack/frame allocation, constant pools, scope cleanup
- ✅ **Control Flow:** All jump operations, short-circuit evaluation
- ✅ **Special Operations:** Nullish coalescing, optional chaining, pattern matching
- ✅ **Compiler:** Constant folding, type hints, dead code elimination (future)

**This means optimizing LITERALLY EVERYTHING that gets evaluated in UExL expressions.**

See **[optimization-rollout-plan.md](optimization-rollout-plan.md#-complete-optimization-scope)** for the complete 10-category inventory of 100+ optimization targets.

---

## 📚 Documentation Structure

This directory contains a comprehensive performance optimization suite. Here's how to use it:

### **Essential Documents (Use Daily)**

| Document | Purpose | When to Use |
|----------|---------|-------------|
| **[0-optimization-guidelines.md](0-optimization-guidelines.md)** | **YOU ARE HERE** - Start here, daily workflow | Always open |
| **[optimization-rollout-plan.md](optimization-rollout-plan.md)** | **PRIMARY ROADMAP** - Phase-by-phase implementation | Main development guide |
| **[dos-and-donts.md](dos-and-donts.md)** | **QUICK REFERENCE** - Code patterns & decisions | Keep in side panel |

### **Measurement & Analysis**

| Document | Purpose | When to Use |
|----------|---------|-------------|
| **[profiling-guide.md](profiling-guide.md)** | CPU profiling walkthrough | Before/after each optimization |
| **[benchmarking-guide.md](benchmarking-guide.md)** | Benchmark best practices | Writing/running benchmarks |

### **Reference Material**

| Document | Purpose | When to Use |
|----------|---------|-------------|
| **[optimization-techniques.md](optimization-techniques.md)** | Pattern library & code examples | Implementing specific optimizations |
| **[best-practices.md](best-practices.md)** | Philosophy & guidelines | Architectural decisions, code review |
| **[optimization-journey.md](optimization-journey.md)** | Historical record | After completing optimizations |
| **[pending-optimizations.md](pending-optimizations.md)** | Future work & research | Planning next optimizations |
| **[README.md](README.md)** | Overview & navigation | Quick stats, links |

---

## 🚀 Quick Start: Your First Day

### 1. **Understand the Context** (15 minutes)

Read these sections (in order):
1. [README.md](README.md) - Performance achievements so far
2. [optimization-journey.md](optimization-journey.md) - What's been optimized already
3. **This document** - How to proceed

### 2. **Set Up Your Workspace** (5 minutes)

```bash
# Navigate to project
cd e:\Projects\uexl\uexl-go

# Verify tests pass
go test ./...

# Establish baseline
go test -bench=BenchmarkVM_Boolean_Current -benchtime=10s -count=10 > baseline.txt
```

**Open in VS Code:**
- **LEFT PANEL:** [optimization-rollout-plan.md](optimization-rollout-plan.md) (main roadmap)
- **RIGHT PANEL:** [dos-and-donts.md](dos-and-donts.md) (quick reference)
- **BOTTOM:** Terminal for testing/profiling

### 3. **Choose Your First Optimization** (10 minutes)

**Recommended Starting Point:** **Check the rollout plan priority order**

The **[optimization-rollout-plan.md](optimization-rollout-plan.md)** lists **10 optimization categories** with **100+ specific targets**:

1. **VM Core Operations** (CRITICAL) - Instruction dispatch, stack ops, frame management
2. **Operator Handlers** (HIGH) - Arithmetic, bitwise, string, unary (comparison done ✅)
3. **Index & Member Access** (HIGH) - Array indexing, object access, optional chaining
4. **Pipe Operations** (HIGH) - Filter, reduce, find, some, every, etc. (Map done ✅)
5. **Built-in Functions** (MEDIUM) - 50+ functions in `vm/builtins.go`
6. **Type System Operations** (MEDIUM) - Type checking, dispatch, conversion
7. **Memory Management** (MEDIUM) - Frame pooling, scope reuse, string building
8. **Compiler Optimizations** (LOW) - Constant folding, type hints, dead code
9. **Control Flow Operations** (MEDIUM) - Fast-paths for jump operations
10. **Special Operations** (LOW) - Nullish coalescing, optional chaining, pattern matching

**Start with highest-impact areas first.** Current recommendation: **Phase 3 (Pipe Operations)** or **Phase 1 (Arithmetic)**.

Open: [optimization-rollout-plan.md](optimization-rollout-plan.md#-complete-optimization-scope)

---

## 📋 Daily Development Workflow

### **Morning: Planning Phase** (15 min)

```bash
□ Open optimization-rollout-plan.md
  → Review current phase objectives
  → Check validation checklist
  → Note files to modify

□ Profile current state (before changes)
  → go test -bench=BenchmarkVM_X -cpuprofile=before.prof -count=10 > before.txt
  → go tool pprof -http=:8080 before.prof
  → Identify bottleneck (must be >5% CPU time)

□ Plan optimization
  → Read relevant section in optimization-techniques.md
  → Check dos-and-donts.md for patterns
  → Write down expected improvement
```

### **Development: Coding Phase** (2-4 hours)

```bash
□ Implement optimization
  → Reference dos-and-donts.md for code patterns
  → Follow "Acceptable Optimizations" guidelines
  → Avoid "Forbidden Optimizations"

□ During coding, ask yourself:
  ✓ Is this a general optimization? (not hardcoded)
  ✓ Will this work for ALL inputs? (not just test cases)
  ✓ Am I adding allocations? (check with -benchmem)
  ✓ Is this readable? (maintainability matters)

□ Keep running tests
  → go test ./vm -v (specific package)
  → Fix any failures immediately
```

### **Validation: Testing Phase** (30 min)

```bash
□ Correctness (MANDATORY)
  → go test ./...
  → go test ./... -race
  → go test ./vm -v -count=5

□ Performance (MANDATORY)
  → go test -bench=BenchmarkVM_X -cpuprofile=after.prof -count=10 > after.txt
  → benchstat before.txt after.txt
  → Check: p-value < 0.05, improvement ≥ 5%

□ Profile Analysis (MANDATORY)
  → go tool pprof -base=before.prof after.prof
  → Verify bottleneck reduced by >50%

□ Allocation Check (MANDATORY)
  → go test -bench=BenchmarkVM_X -benchmem
  → Verify: 0 B/op, 0 allocs/op

□ Cross-Validation
  → Test with different expression types
  → Ensure no regressions in other benchmarks
```

### **Documentation: Recording Phase** (15 min)

```bash
□ Update optimization-journey.md
  → Add new section for this optimization
  → Include before/after results
  → Add CPU profile analysis
  → Document lessons learned

□ Update optimization-rollout-plan.md
  → Check off completed phase
  → Note actual vs expected results
  → Update if timeline changes

□ Commit changes
  → git add .
  → git commit -m "optimize: [component] - [description] (Xns → Yns, Z% gain)"
  → Example: "optimize: pipe handlers - scope reuse (1500ns → 1000ns, 33% gain)"
```

---

## 🎯 Performance Targets

### **Current State**
- Boolean expressions: **62 ns/op** (0 allocs)
- Arithmetic: ~80 ns/op
- String concat: ~100 ns/op
- Map pipe (10 items): ~1500 ns/op
- Filter pipe (10 items): ~1800 ns/op

### **Targets**

| Tier | Target | Improvement | Status |
|------|--------|-------------|--------|
| **Stretch Goal** | 20 ns/op | 68% faster | Ambitious but possible |
| **Realistic** | 30-35 ns/op | 45-50% faster | Very achievable |
| **Minimum** | 40 ns/op | 35% faster | Guaranteed |

**All tiers beat competitors:**
- expr: 105 ns/op
- cel-go: 127 ns/op

**Mandatory:** 0 B/op, 0 allocs/op (no allocations allowed)

---

## ✅ Validation Checklist (MANDATORY)

### **Before You Start**
- [ ] Profile baseline established
- [ ] Benchmark baseline saved (10+ runs)
- [ ] Bottleneck identified (>5% CPU in profile)
- [ ] Expected improvement documented

### **During Development**
- [ ] No hardcoded values for specific expressions
- [ ] No test-specific code paths
- [ ] Code follows best practices (readable, maintainable)
- [ ] Tests passing continuously

### **Before Committing**
- [ ] ✅ All tests pass: `go test ./...`
- [ ] ✅ Race detector clean: `go test ./... -race`
- [ ] ✅ Benchmark improved: `benchstat` shows p < 0.05, improvement ≥ 5%
- [ ] ✅ Profile improved: Bottleneck reduced >50%
- [ ] ✅ Zero allocations: `0 B/op, 0 allocs/op`
- [ ] ✅ No regressions: Other benchmarks not slower
- [ ] ✅ Cross-validated: Tested with varied inputs
- [ ] ✅ Documented: Updated optimization-journey.md

**If ANY checkbox is unchecked, DO NOT commit.**

---

## ❌ Forbidden Optimizations

These will be **REJECTED in code review:**

### 1. Hardcoded Results
```go
// ❌ FORBIDDEN
if vm.currentExpr == "a && b" {
    return true, nil  // Hardcoded for benchmark
}
```

### 2. Test-Specific Paths
```go
// ❌ FORBIDDEN
if testing.Testing() {
    return fastPath()  // Different behavior in tests
}
```

### 3. Expression Caching
```go
// ❌ FORBIDDEN (not general optimization)
var exprCache = map[string]any{}
if cached, ok := exprCache[expr]; ok {
    return cached
}
```

### 4. Skipping Validation
```go
// ❌ FORBIDDEN
// +build !test
func (vm *VM) Push(val any) {
    vm.stack[vm.sp] = val  // No bounds check in production
    vm.sp++
}
```

### 5. Benchmark Detection
```go
// ❌ FORBIDDEN
if len(bytecode.Instructions) == 7 {  // Specific to test
    return quickResult()
}
```

---

## ✅ Acceptable Optimizations

These are **ENCOURAGED:**

### 1. Type-Specific Dispatch
```go
// ✅ ACCEPTABLE - General for all types
func (vm *VM) executeNumberArithmetic(op code.Opcode, left, right float64) error {
    // No type assertions, works for all float64 operations
}
```

### 2. Pre-computation
```go
// ✅ ACCEPTABLE - Computed once, reused
func (vm *VM) setBaseInstructions(bc *ByteCode, ctx map[string]any) {
    vm.contextVarCache = buildCache(ctx)  // Pre-compute
}
```

### 3. Scope/Frame Reuse
```go
// ✅ ACCEPTABLE - Avoid allocations
vm.pushPipeScope()  // Once
for item := range arr {
    pipeScope["$item"] = item  // Reuse
}
vm.popPipeScope()  // Once
```

### 4. Pattern Detection (General)
```go
// ✅ ACCEPTABLE - Works for ANY matching pattern
if isSimpleArithmetic(instructions) {
    return vectorizedOperation()  // Not specific to one expression
}
```

### 5. Sentinel Values
```go
// ✅ ACCEPTABLE - Avoid allocations
var contextVarNotProvided = contextVarMissing{}  // Singleton
```

---

## 🔍 Measurement Protocol

### **Establishing Baseline**
```bash
# Clean cache
go clean -testcache

# Run baseline benchmarks (20 iterations for stability)
go test -bench=BenchmarkVM_Boolean_Current \
  -benchtime=20s \
  -count=20 \
  -cpuprofile=baseline.prof \
  > baseline.txt

# Profile baseline
go tool pprof -http=:8080 baseline.prof

# Save baseline
git add baseline.txt baseline.prof
git commit -m "baseline: establish performance reference"
```

### **After Each Optimization**
```bash
# Run optimized benchmarks
go test -bench=BenchmarkVM_Boolean_Current \
  -benchtime=20s \
  -count=20 \
  -cpuprofile=optimized.prof \
  > optimized.txt

# Statistical comparison (p < 0.05 required)
benchstat baseline.txt optimized.txt

# Must show:
# name                old time/op  new time/op  delta
# VM_Boolean_Current    62.0ns ± 2%  XX.Xns ± 2%  -YY.YY%  (p=0.000 n=20+20)
#                                                  ^^^^^^^^
#                                                  p-value < 0.05 = significant

# Profile comparison
go tool pprof -base=baseline.prof optimized.prof -top
```

### **Acceptance Criteria**
- ✅ p-value < 0.05 (statistically significant)
- ✅ Improvement ≥ 5% (meaningful gain)
- ✅ Variance ±2-3% (stable, reproducible)
- ✅ 0 allocs/op (no new allocations)
- ✅ All tests pass (no regressions)

---

## 📊 Implementation Priority

Follow this order (from [optimization-rollout-plan.md](optimization-rollout-plan.md)):

### **🔴 HIGH PRIORITY (Weeks 1-2)**

**Phase 3: Pipe Operations** ← **START HERE**
- FilterPipeHandler, ReducePipeHandler, etc.
- Apply scope/frame reuse pattern
- Expected: 15-25% improvement
- Highest user-visible impact

**Phase 1: Arithmetic Operations**
- Type-specific function signatures
- Expected: 5-8% improvement

**Phase 2: String Operations**
- Remove type assertions
- Expected: 3-5% improvement

### **🟡 MEDIUM PRIORITY (Week 3)**

**Phase 4: Array/Object Access**
- Type-specific indexing
- Expected: 5-7% improvement

### **🟢 LOW PRIORITY (Week 4+)**

**Phase 5: Unary Operations**
- Expected: 2-4% improvement

**Phase 6: Boolean Operations**
- Expected: 1-2% improvement

---

## 🛠️ Tools & Commands

### **Testing**
```bash
# All tests
go test ./...

# With race detector
go test ./... -race

# Specific package
go test ./vm -v

# Repeat 5 times
go test ./vm -count=5
```

### **Benchmarking**
```bash
# Run benchmark
go test -bench=BenchmarkVM_Boolean_Current

# With memory stats
go test -bench=BenchmarkVM_Boolean_Current -benchmem

# Multiple runs for stability
go test -bench=BenchmarkVM_Boolean_Current -count=10 -benchtime=10s > results.txt

# Statistical comparison
benchstat baseline.txt optimized.txt
```

### **Profiling**
```bash
# CPU profile
go test -bench=BenchmarkVM_Boolean_Current -cpuprofile=cpu.prof

# Analyze interactively
go tool pprof cpu.prof
# Then: top, list, web, etc.

# Web interface
go tool pprof -http=:8080 cpu.prof

# Compare profiles
go tool pprof -base=before.prof after.prof
```

---

## 🎓 Learning Path

### **Day 1: Understand Current State**
1. Read [README.md](README.md) - Overview
2. Read [optimization-journey.md](optimization-journey.md) - History
3. Run current benchmarks, establish baseline
4. Review [best-practices.md](best-practices.md) - Philosophy

### **Day 2: Learn Techniques**
1. Read [optimization-techniques.md](optimization-techniques.md) - Patterns
2. Study [dos-and-donts.md](dos-and-donts.md) - Examples
3. Review existing optimized code (comparison operators)
4. Practice profiling with [profiling-guide.md](profiling-guide.md)

### **Day 3+: Start Optimizing**
1. Follow [optimization-rollout-plan.md](optimization-rollout-plan.md)
2. Start with Phase 3 (Pipe Operations)
3. Use this document as daily checklist
4. Document results in [optimization-journey.md](optimization-journey.md)

---

## 🆘 Troubleshooting

### **Tests Failing After Optimization**
1. Revert changes: `git checkout .`
2. Review [dos-and-donts.md](dos-and-donts.md) - Common mistakes
3. Check if you violated any "Forbidden Optimizations"
4. Run tests in isolation: `go test ./vm -v -run=TestSpecificTest`

### **Benchmark Slower Than Baseline**
1. Profile both: `go tool pprof -base=baseline.prof current.prof`
2. Check for new allocations: `go test -bench=X -benchmem`
3. Look for accidental complexity added
4. Verify you're testing same code path

### **Results Unstable (High Variance)**
1. Close other applications
2. Run longer: `-benchtime=20s -count=20`
3. Check system isn't throttling (thermal/power)
4. Use `benchstat` for statistical validation

### **Not Sure What to Optimize**
1. Profile first: `go test -bench=X -cpuprofile=cpu.prof`
2. Look at `pprof -top` output
3. Only optimize functions showing >5% CPU time
4. Consult [optimization-rollout-plan.md](optimization-rollout-plan.md) for priorities

---

## 📞 Quick Reference

**Need to know what to work on?**
→ [optimization-rollout-plan.md](optimization-rollout-plan.md)

**Writing code right now?**
→ [dos-and-donts.md](dos-and-donts.md)

**Need to measure performance?**
→ [profiling-guide.md](profiling-guide.md) + [benchmarking-guide.md](benchmarking-guide.md)

**Looking for a specific technique?**
→ [optimization-techniques.md](optimization-techniques.md)

**Finished an optimization?**
→ Update [optimization-journey.md](optimization-journey.md)

**Planning long-term?**
→ [pending-optimizations.md](pending-optimizations.md) + [best-practices.md](best-practices.md)

---

## 🎯 Success Criteria

You've successfully optimized when:
- ✅ All tests pass (`go test ./...`)
- ✅ Benchmark improved ≥5% (p < 0.05)
- ✅ Profile shows bottleneck reduced >50%
- ✅ Zero allocations maintained
- ✅ No regressions in other benchmarks
- ✅ Code is readable and maintainable
- ✅ Documentation updated
- ✅ Code reviewed and approved

---

## 📝 Summary: The Essential 3

**Keep these 3 documents open during development:**

1. **[optimization-rollout-plan.md](optimization-rollout-plan.md)** - Main roadmap (LEFT PANEL)
2. **[dos-and-donts.md](dos-and-donts.md)** - Quick reference (RIGHT PANEL)
3. **[profiling-guide.md](profiling-guide.md)** / **[benchmarking-guide.md](benchmarking-guide.md)** - Measurement (BOTTOM PANEL)

**Update after each optimization:**
- **[optimization-journey.md](optimization-journey.md)** - Historical record

**Everything else is reference material - use as needed!**

---

**Ready to start?** → Open [optimization-rollout-plan.md](optimization-rollout-plan.md) and begin with Phase 3! 🚀
