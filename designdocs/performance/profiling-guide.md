# CPU Profiling Guide

## Overview

This guide covers how to profile UExL VM performance, analyze results, and identify optimization opportunities using Go's built-in profiling tools.

---

## Quick Start

### Running a CPU Profile

```bash
# Profile a specific benchmark
go test -bench=BenchmarkVM_Boolean_Current -cpuprofile=cpu.prof

# Profile all benchmarks
go test -bench=. -cpuprofile=cpu_all.prof

# Profile with extended runtime for accuracy
go test -bench=BenchmarkVM_Boolean_Current -benchtime=20s -cpuprofile=cpu_long.prof
```

### Analyzing the Profile

```bash
# Interactive mode
go tool pprof cpu.prof

# Web interface (opens browser)
go tool pprof -http=:8080 cpu.prof

# Generate SVG graph
go tool pprof -svg cpu.prof > cpu_graph.svg
```

---

## Step-by-Step Profiling Workflow

### 1. Establish Baseline

Before making any changes:

```bash
# Run benchmark and save baseline
go test -bench=BenchmarkVM_Boolean_Current -benchtime=10s > baseline.txt

# Capture CPU profile
go test -bench=BenchmarkVM_Boolean_Current -benchtime=10s -cpuprofile=baseline.prof

# Save for comparison
cp baseline.prof profiles/baseline_$(date +%Y%m%d).prof
```

### 2. Profile Current State

```bash
# Run with CPU profiling
go test -bench=BenchmarkVM_Boolean_Current -cpuprofile=current.prof

# Verify test passes
go test ./... -v
```

### 3. Analyze Bottlenecks

#### Using Interactive Mode

```bash
go tool pprof current.prof

# Common commands:
(pprof) top           # Show top CPU consumers
(pprof) top10         # Top 10 functions
(pprof) list Run      # Show annotated source for vm.Run
(pprof) web           # Generate call graph (requires graphviz)
(pprof) peek Run      # Show callers and callees
(pprof) quit          # Exit
```

#### Using Web Interface

```bash
go tool pprof -http=localhost:8080 current.prof

# Navigate to:
# - Top: Function ranking by CPU time
# - Graph: Visual call graph
# - Flame Graph: Hierarchical time distribution
# - Source: Annotated source code
```

### 4. Interpret Results

Look for:

- **High % in single function** → Direct optimization target
- **Many small % across functions** → Systemic issue (e.g., allocations)
- **Unexpected functions in top 10** → Hidden bottlenecks

Example output:
```
(pprof) top10
      flat  flat%   sum%        cum   cum%
     2.84s 26.79% 26.79%      2.84s 26.79%  runtime.mapaccess2_faststr
     1.47s 13.87% 40.66%      1.47s 13.87%  github.com/maniartech/uexl-go/vm.(*VM).Pop
     0.84s  7.94% 48.58%      0.84s  7.94%  github.com/maniartech/uexl-go/vm.(*VM).Push
```

**Analysis:**
- 26.79% in map access → Context variable lookups are expensive
- 13.87% + 7.94% in Push/Pop → Stack operations need optimization

### 5. Make Changes

Based on profile analysis, implement optimizations (see [optimization-techniques.md](optimization-techniques.md))

### 6. Profile Again

```bash
# After changes
go test -bench=BenchmarkVM_Boolean_Current -cpuprofile=optimized.prof

# Compare profiles
go tool pprof -base=current.prof optimized.prof
```

### 7. Validate Improvement

```bash
# Run benchstat for statistical comparison
go test -bench=BenchmarkVM_Boolean_Current -count=10 > new.txt
benchstat baseline.txt new.txt
```

---

## Common Profiling Scenarios

### Scenario 1: Finding Memory Allocations

```bash
# Memory profile
go test -bench=. -memprofile=mem.prof -memprofilerate=1

# Analyze
go tool pprof -alloc_space mem.prof

(pprof) top
# Look for unexpected allocations
```

### Scenario 2: Comparing Before/After

```bash
# Before optimization
go test -bench=BenchmarkVM_Boolean_Current -cpuprofile=before.prof

# Make changes...

# After optimization
go test -bench=BenchmarkVM_Boolean_Current -cpuprofile=after.prof

# Diff profiles
go tool pprof -base=before.prof after.prof

(pprof) top
# Negative values = improvement
# Positive values = regression
```

### Scenario 3: Profiling Specific Code Path

```go
// Add targeted benchmark
func BenchmarkVM_ContextVarAccess(b *testing.B) {
    // Setup
    node, _ := parser.ParseString("a")
    comp := compiler.New()
    comp.Compile(node)
    bc := comp.ByteCode()
    
    machine := vm.New(vm.LibContext{})
    ctx := map[string]any{"a": 42.0}
    
    // Profile just context var access
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        machine.Run(bc, ctx)
    }
}
```

```bash
go test -bench=BenchmarkVM_ContextVarAccess -cpuprofile=context_access.prof
```

---

## Advanced Profiling Techniques

### 1. Continuous Profiling

Track performance over time:

```bash
#!/bin/bash
# scripts/continuous_profile.sh

DATE=$(date +%Y%m%d_%H%M%S)
PROFILE_DIR="profiles/history"

mkdir -p $PROFILE_DIR

# Run benchmark with profiling
go test -bench=. -cpuprofile=$PROFILE_DIR/cpu_$DATE.prof \
    -benchtime=20s > $PROFILE_DIR/bench_$DATE.txt

# Store commit hash
git rev-parse HEAD > $PROFILE_DIR/commit_$DATE.txt

echo "Profile saved to $PROFILE_DIR/cpu_$DATE.prof"
```

### 2. Differential Profiling

Compare two commits:

```bash
# Profile current commit
git checkout feature-branch
go test -bench=. -cpuprofile=feature.prof

# Profile base commit
git checkout main
go test -bench=. -cpuprofile=main.prof

# Compare
go tool pprof -base=main.prof feature.prof -http=:8080
```

### 3. Flame Graphs

Visualize time distribution:

```bash
# Generate flame graph
go tool pprof -http=:8080 cpu.prof

# Navigate to "Flame Graph" view in browser
# - Width = time spent
# - Height = call depth
# - Color = package/module
```

---

## Interpreting Profile Metrics

### Flat vs Cumulative Time

- **Flat:** Time spent in function itself (excluding callees)
- **Cumulative:** Time spent in function + all callees

```
(pprof) list Run
     flat    cum
     0.5s   10.2s   func (vm *VM) Run(bytecode, ctx) {
     0.1s    0.1s       vm.setBaseInstructions(bytecode, ctx)
     0.3s    9.8s       for vm.ip < len(vm.instructions) {
     0.1s    9.5s           op := vm.instructions[vm.ip]
                            vm.executeOp(op)  // <-- 9.4s cumulative
                        }
                    }
```

**Analysis:**
- Flat 0.5s in Run itself (loop overhead)
- Cumulative 10.2s includes all execution (9.7s in callees)
- Most time in executeOp → Investigate opcode handlers

### Sample Count

Profile samples at ~100Hz (every 10ms):

```
Total samples: 1000
Function X: 250 samples = 25% CPU time
```

**Low sample counts (<50):** May be noise  
**High sample counts (>100):** Reliable bottleneck

### Inlining Effects

Inlined functions don't appear in profiles:

```go
//go:noinline
func slowFunction() { ... }  // Will appear in profile

func fastFunction() { ... }  // May be inlined, disappears from profile
```

**Tip:** If function missing, check if inlined with:
```bash
go build -gcflags='-m' 2>&1 | grep inlined
```

---

## Common Bottleneck Patterns

### Pattern 1: Map Operations

```
(pprof) top
  26.79%  runtime.mapaccess2_faststr
  12.34%  runtime.mapassign_faststr
```

**Cause:** Frequent map lookups/assignments  
**Solution:** Cache values, use arrays/slices

### Pattern 2: Interface Conversions

```
(pprof) top
  15.23%  runtime.assertI2I
  10.45%  runtime.convT64
```

**Cause:** Type assertions, interface boxing  
**Solution:** Type-specific functions, avoid any

### Pattern 3: Allocations

```
(pprof) top -alloc_space
  1024MB  runtime.makeslice
   512MB  runtime.newobject
```

**Cause:** Unnecessary allocations  
**Solution:** Reuse buffers, object pooling

### Pattern 4: String Operations

```
(pprof) top
  18.67%  runtime.concatstrings
  12.34%  runtime.slicebytetostring
```

**Cause:** String concatenation, conversions  
**Solution:** strings.Builder, []byte operations

---

## Profiling Best Practices

### 1. Sufficient Benchmark Time

```bash
# Too short - noisy results
go test -bench=. -benchtime=1s

# Good - stable results
go test -bench=. -benchtime=10s

# Better - very stable
go test -bench=. -benchtime=20s
```

### 2. Isolated Benchmarks

```go
// Bad: Mixed concerns
func BenchmarkEverything(b *testing.B) {
    // Parse, compile, execute
}

// Good: Focused benchmarks
func BenchmarkParse(b *testing.B) { /* Only parsing */ }
func BenchmarkCompile(b *testing.B) { /* Only compilation */ }
func BenchmarkExecute(b *testing.B) { /* Only execution */ }
```

### 3. Realistic Workloads

```go
// Bad: Artificial benchmark
func BenchmarkVM_Simple(b *testing.B) {
    vm.Run(bytecode, map[string]any{"a": 1})
}

// Good: Production-like
func BenchmarkVM_Realistic(b *testing.B) {
    // Reuse context maps (like production)
    ctx := map[string]any{
        "user": map[string]any{"age": 25, "name": "John"},
        "items": []any{1.0, 2.0, 3.0},
    }
    for i := 0; i < b.N; i++ {
        vm.Run(bytecode, ctx)  // Reuse ctx
    }
}
```

### 4. Multiple Iterations

```bash
# Run multiple times for statistical validity
go test -bench=. -count=10 > results.txt
benchstat results.txt
```

---

## Tools and Resources

### Required Tools

```bash
# Install graphviz (for call graphs)
# Windows: choco install graphviz
# Mac: brew install graphviz
# Linux: apt-get install graphviz

# Install benchstat
go install golang.org/x/perf/cmd/benchstat@latest
```

### Useful Scripts

#### Compare Two Profiles

```bash
#!/bin/bash
# scripts/compare_profiles.sh

if [ $# -ne 2 ]; then
    echo "Usage: $0 <before.prof> <after.prof>"
    exit 1
fi

BEFORE=$1
AFTER=$2

echo "=== Top Functions ==="
go tool pprof -top -base=$BEFORE $AFTER

echo ""
echo "=== Detailed Diff ==="
go tool pprof -base=$BEFORE $AFTER -http=:8080
```

#### Profile All Benchmarks

```bash
#!/bin/bash
# scripts/profile_all.sh

BENCHMARKS=$(go test -bench=. -list='.*' | grep '^Benchmark')

for bench in $BENCHMARKS; do
    echo "Profiling $bench..."
    go test -bench=$bench -cpuprofile=profiles/${bench}.prof
done

echo "Profiles saved to profiles/"
ls -lh profiles/
```

---

## Troubleshooting

### Profile Shows Nothing

**Cause:** Benchmark too fast (< 10ms)  
**Solution:** Increase `-benchtime` or add more work

```bash
go test -bench=. -benchtime=20s -cpuprofile=cpu.prof
```

### Profile Shows Only Runtime Functions

**Cause:** Inlining or very fast execution  
**Solution:** Disable inlining during profiling

```bash
go test -bench=. -gcflags='-l' -cpuprofile=cpu.prof
```

### Inconsistent Results

**Cause:** Noise from other processes  
**Solution:** Run on dedicated machine or multiple iterations

```bash
# Close other applications
go test -bench=. -count=20 > results.txt
benchstat results.txt
```

### Web Interface Doesn't Open

**Cause:** Missing graphviz or port conflict  
**Solution:** Use different port or install graphviz

```bash
go tool pprof -http=localhost:9090 cpu.prof
```

---

## Real-World Example

### Problem: Slow Boolean Expression Evaluation

**Step 1: Baseline**
```bash
$ go test -bench=BenchmarkVM_Boolean -benchtime=10s
BenchmarkVM_Boolean-8   11309043   106.2 ns/op
```

**Step 2: Profile**
```bash
$ go test -bench=BenchmarkVM_Boolean -cpuprofile=slow.prof
$ go tool pprof -top slow.prof

(pprof) top
      flat  flat%   sum%        cum   cum%
     2.84s 26.79% 26.79%      2.84s 26.79%  runtime.mapaccess2_faststr
     1.47s 13.87% 40.66%      1.47s 13.87%  vm.(*VM).Pop
     0.84s  7.94% 48.58%      0.84s  7.94%  vm.(*VM).Push
```

**Analysis:**
- 26.79% in map access → Context variables
- 21.81% in stack operations → Push/Pop overhead

**Step 3: Optimize Context Access**

Implemented context variable caching (see [optimization-journey.md](optimization-journey.md#optimization-phase-2-context-variable-caching))

**Step 4: Profile Again**
```bash
$ go test -bench=BenchmarkVM_Boolean -cpuprofile=optimized.prof
$ go tool pprof -base=slow.prof optimized.prof

(pprof) top
      flat  flat%   sum%        cum   cum%
    -2.84s 100%   100%     -2.84s 100%   runtime.mapaccess2_faststr  # ELIMINATED
```

**Result:** 106ns → 95ns (10% improvement)

---

## References

- [Go Blog: Profiling Go Programs](https://go.dev/blog/pprof)
- [pprof Documentation](https://github.com/google/pprof/tree/main/doc)
- [optimization-journey.md](optimization-journey.md) - Optimization history
- [optimization-techniques.md](optimization-techniques.md) - Technique catalog

**Last Updated:** October 17, 2025
