# Benchmarking Guide

## Overview

This guide covers how to write, run, and interpret benchmarks for UExL VM performance testing. Includes best practices for reliable measurements and statistical analysis.

---

## Quick Start

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=.

# Run specific benchmark
go test -bench=BenchmarkVM_Boolean_Current

# Run with memory stats
go test -bench=. -benchmem

# Extended run for stability
go test -bench=. -benchtime=10s

# Statistical analysis (10 iterations)
go test -bench=. -count=10 | tee results.txt
benchstat results.txt
```

---

## Writing Benchmarks

### Basic Benchmark Structure

```go
func BenchmarkVM_Boolean_Current(b *testing.B) {
    // Setup (not timed)
    node, err := parser.ParseString("a && b || c")
    if err != nil {
        b.Fatal(err)
    }
    
    comp := compiler.New()
    if err := comp.Compile(node); err != nil {
        b.Fatal(err)
    }
    bytecode := comp.ByteCode()
    
    machine := vm.New(vm.LibContext{
        Functions:    vm.Builtins,
        PipeHandlers: vm.DefaultPipeHandlers,
    })
    
    ctx := map[string]any{
        "a": true,
        "b": false,
        "c": true,
    }
    
    // Reset timer before measurement
    b.ResetTimer()
    
    // Benchmark loop (timed)
    for i := 0; i < b.N; i++ {
        result, err := machine.Run(bytecode, ctx)
        if err != nil {
            b.Fatal(err)
        }
        _ = result  // Prevent optimization
    }
}
```

### Key Elements

1. **Setup Outside Loop:** Parse/compile once, execute many times
2. **b.ResetTimer():** Exclude setup from timing
3. **Error Handling:** Use `b.Fatal()` not `t.Fatal()`
4. **Prevent Optimization:** Use results (`_ = result`)
5. **Realistic Data:** Use production-like inputs

---

## Benchmark Patterns

### Pattern 1: Isolated Component

Test single component in isolation:

```go
func BenchmarkVM_StackOperations(b *testing.B) {
    machine := vm.New(vm.LibContext{})
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        machine.Push(42.0)
        machine.Push(10.0)
        machine.Pop()
        machine.Pop()
    }
}
```

### Pattern 2: End-to-End

Test complete execution path:

```go
func BenchmarkVM_CompleteExpression(b *testing.B) {
    expr := "users |filter: $item.age > 18 |map: $item.name"
    
    // Setup
    node, _ := parser.ParseString(expr)
    comp := compiler.New()
    comp.Compile(node)
    bytecode := comp.ByteCode()
    
    machine := vm.New(vm.LibContext{
        PipeHandlers: vm.DefaultPipeHandlers,
    })
    
    ctx := map[string]any{
        "users": []any{
            map[string]any{"name": "Alice", "age": 25.0},
            map[string]any{"name": "Bob", "age": 17.0},
            map[string]any{"name": "Carol", "age": 30.0},
        },
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        machine.Run(bytecode, ctx)
    }
}
```

### Pattern 3: Comparative

Compare different implementations:

```go
func BenchmarkVM_ContextAccess_Map(b *testing.B) {
    // Old: Map lookup
    ctx := map[string]any{"a": 42.0}
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = ctx["a"]
    }
}

func BenchmarkVM_ContextAccess_Cache(b *testing.B) {
    // New: Array access
    cache := []any{42.0}
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = cache[0]
    }
}
```

### Pattern 4: Realistic Workload

Simulate production usage:

```go
func BenchmarkVM_ProductionLike(b *testing.B) {
    expressions := []string{
        "user.age > 18 && user.verified",
        "items |filter: $item.price > 100 |map: $item.name",
        "total * 1.15",  // Tax calculation
    }
    
    contexts := []map[string]any{
        {"user": map[string]any{"age": 25.0, "verified": true}},
        {"items": []any{...}},
        {"total": 100.0},
    }
    
    // Pre-compile all expressions
    bytecodes := make([]*compiler.ByteCode, len(expressions))
    for i, expr := range expressions {
        node, _ := parser.ParseString(expr)
        comp := compiler.New()
        comp.Compile(node)
        bytecodes[i] = comp.ByteCode()
    }
    
    machine := vm.New(vm.LibContext{
        Functions:    vm.Builtins,
        PipeHandlers: vm.DefaultPipeHandlers,
    })
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        idx := i % len(expressions)
        machine.Run(bytecodes[idx], contexts[idx])
    }
}
```

---

## Interpreting Results

### Basic Metrics

```bash
$ go test -bench=BenchmarkVM_Boolean -benchmem

BenchmarkVM_Boolean-8   19354838   62.08 ns/op   0 B/op   0 allocs/op
                        ^^^^^^^^   ^^^^^^^^^^^   ^^^^^^   ^^^^^^^^^^^
                        iterations   time/op     bytes    allocations
```

**Metrics:**
- **Iterations:** Number of times loop ran (higher = more stable)
- **ns/op:** Nanoseconds per operation (lower = faster)
- **B/op:** Bytes allocated per operation (lower = better)
- **allocs/op:** Allocations per operation (0 = optimal)

### Performance Targets

Based on UExL goals:

| Operation Type | Target | Good | Acceptable |
|----------------|--------|------|------------|
| Simple boolean | < 50ns | 50-100ns | 100-200ns |
| Arithmetic | < 100ns | 100-200ns | 200-500ns |
| String ops | < 200ns | 200-500ns | 500-1000ns |
| Pipe (small) | < 1µs | 1-5µs | 5-10µs |
| Pipe (large) | < 10µs | 10-50µs | 50-100µs |

### Stability Indicators

```bash
$ go test -bench=. -count=10 | tee results.txt
$ benchstat results.txt

name                  time/op
VM_Boolean_Current    62.1ns ± 2%
                      ^^^^^^^^^  ^
                      mean       variation
```

**Interpretation:**
- **± < 5%:** Stable, reliable measurements
- **± 5-10%:** Acceptable, some noise
- **± > 10%:** Unstable, re-run on quiet machine

---

## Statistical Analysis with benchstat

### Installation

```bash
go install golang.org/x/perf/cmd/benchstat@latest
```

### Comparing Implementations

```bash
# Baseline
git checkout main
go test -bench=BenchmarkVM_Boolean -count=10 > old.txt

# After optimization
git checkout feature
go test -bench=BenchmarkVM_Boolean -count=10 > new.txt

# Compare
benchstat old.txt new.txt
```

**Output:**
```
name                  old time/op  new time/op  delta
VM_Boolean_Current    106ns ± 3%    62ns ± 2%  -41.51%  (p=0.000 n=10+10)
                      ^^^^^^^^^     ^^^^^^^^   ^^^^^^^^  ^^^^^^^^^^^^^^^^^^
                      before        after      change    statistical significance
```

**Interpretation:**
- **Delta:** -41.51% = 41.51% faster
- **p-value:** < 0.05 = statistically significant
- **n=10+10:** 10 samples each side

### Multiple Benchmarks

```bash
$ benchstat old.txt new.txt

name                     old time/op    new time/op    delta
VM_Boolean_Current         106ns ± 3%      62ns ± 2%  -41.51%  (p=0.000 n=10+10)
VM_Arithmetic_Current      125ns ± 5%     120ns ± 3%   -4.00%  (p=0.023 n=10+10)
VM_String_Current          250ns ± 7%     245ns ± 4%      ~     (p=0.089 n=10+10)
                                                        ^
                                                        ~ = no significant change
```

---

## Advanced Benchmarking

### Sub-Benchmarks

Group related benchmarks:

```go
func BenchmarkVM_Comparisons(b *testing.B) {
    testCases := []struct {
        name string
        expr string
        ctx  map[string]any
    }{
        {"Numbers", "a > b", map[string]any{"a": 10.0, "b": 5.0}},
        {"Strings", "name == 'Alice'", map[string]any{"name": "Alice"}},
        {"Booleans", "active && verified", map[string]any{"active": true, "verified": true}},
    }
    
    for _, tc := range testCases {
        b.Run(tc.name, func(b *testing.B) {
            // Setup
            node, _ := parser.ParseString(tc.expr)
            comp := compiler.New()
            comp.Compile(node)
            bytecode := comp.ByteCode()
            machine := vm.New(vm.LibContext{})
            
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                machine.Run(bytecode, tc.ctx)
            }
        })
    }
}
```

**Run:**
```bash
$ go test -bench=BenchmarkVM_Comparisons

BenchmarkVM_Comparisons/Numbers-8     20000000    65.3 ns/op
BenchmarkVM_Comparisons/Strings-8     15000000    78.2 ns/op
BenchmarkVM_Comparisons/Booleans-8    19000000    62.1 ns/op
```

### Benchmark Table

Test multiple scenarios systematically:

```go
func BenchmarkVM_ScaleTest(b *testing.B) {
    sizes := []int{10, 100, 1000, 10000}
    
    for _, size := range sizes {
        b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
            // Create array of specified size
            items := make([]any, size)
            for i := range items {
                items[i] = float64(i)
            }
            
            node, _ := parser.ParseString("arr |map: $item * 2")
            comp := compiler.New()
            comp.Compile(node)
            bytecode := comp.ByteCode()
            
            machine := vm.New(vm.LibContext{
                PipeHandlers: vm.DefaultPipeHandlers,
            })
            
            ctx := map[string]any{"arr": items}
            
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                machine.Run(bytecode, ctx)
            }
        })
    }
}
```

### Parallel Benchmarks

Test concurrent performance:

```go
func BenchmarkVM_Parallel(b *testing.B) {
    node, _ := parser.ParseString("a && b")
    comp := compiler.New()
    comp.Compile(node)
    bytecode := comp.ByteCode()
    
    ctx := map[string]any{"a": true, "b": false}
    
    b.RunParallel(func(pb *testing.PB) {
        // Each goroutine gets its own VM
        machine := vm.New(vm.LibContext{})
        
        for pb.Next() {
            machine.Run(bytecode, ctx)
        }
    })
}
```

---

## Cross-Library Comparison

### Benchmark Against Competitors

```go
// UExL
func BenchmarkUExL_Boolean(b *testing.B) {
    node, _ := parser.ParseString("a && b || c")
    comp := compiler.New()
    comp.Compile(node)
    bytecode := comp.ByteCode()
    machine := vm.New(vm.LibContext{})
    ctx := map[string]any{"a": true, "b": false, "c": true}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        machine.Run(bytecode, ctx)
    }
}

// expr (competitor)
func BenchmarkExpr_Boolean(b *testing.B) {
    program, _ := expr.Compile("a && b || c")
    env := map[string]any{"a": true, "b": false, "c": true}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        expr.Run(program, env)
    }
}

// cel-go (competitor)
func BenchmarkCelGo_Boolean(b *testing.B) {
    env, _ := cel.NewEnv(
        cel.Variable("a", cel.BoolType),
        cel.Variable("b", cel.BoolType),
        cel.Variable("c", cel.BoolType),
    )
    ast, _ := env.Compile("a && b || c")
    prg, _ := env.Program(ast)
    
    ctx := map[string]any{"a": true, "b": false, "c": true}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        prg.Eval(ctx)
    }
}
```

**Compare:**
```bash
$ go test -bench=Boolean

BenchmarkUExL_Boolean-8     19354838    62.08 ns/op
BenchmarkExpr_Boolean-8      9511806   105.3 ns/op
BenchmarkCelGo_Boolean-8     7456231   127.8 ns/op
```

**Analysis:**
- UExL: 62ns (baseline)
- expr: 105ns (1.7x slower)
- cel-go: 127ns (2.1x slower)

---

## Best Practices

### 1. Realistic Benchmarks

```go
// ❌ Bad: Unrealistic (new context every time)
func BenchmarkVM_Unrealistic(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ctx := map[string]any{"a": 42.0}  // Allocation overhead
        machine.Run(bytecode, ctx)
    }
}

// ✅ Good: Realistic (reuse context like production)
func BenchmarkVM_Realistic(b *testing.B) {
    ctx := map[string]any{"a": 42.0}  // Reused
    for i := 0; i < b.N; i++ {
        machine.Run(bytecode, ctx)
    }
}
```

### 2. Sufficient Runtime

```bash
# ❌ Too short - unstable results
go test -bench=. -benchtime=100ms

# ✅ Good - stable results
go test -bench=. -benchtime=10s

# ✅ Best - very stable with stats
go test -bench=. -benchtime=10s -count=10 | benchstat
```

### 3. Isolated Testing

```go
// ❌ Bad: Mixed concerns
func BenchmarkEverything(b *testing.B) {
    for i := 0; i < b.N; i++ {
        node, _ := parser.ParseString("a + b")  // Parser overhead
        comp := compiler.New()
        comp.Compile(node)                       // Compiler overhead
        bytecode := comp.ByteCode()
        machine.Run(bytecode, ctx)               // VM overhead
    }
}

// ✅ Good: Focused on VM
func BenchmarkVM_Only(b *testing.B) {
    // Parse & compile once (outside timer)
    node, _ := parser.ParseString("a + b")
    comp := compiler.New()
    comp.Compile(node)
    bytecode := comp.ByteCode()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        machine.Run(bytecode, ctx)  // Only VM timed
    }
}
```

### 4. Prevent Compiler Optimizations

```go
// ❌ Bad: Result unused, may be optimized away
func BenchmarkVM_Optimized(b *testing.B) {
    for i := 0; i < b.N; i++ {
        machine.Run(bytecode, ctx)  // Compiler may skip
    }
}

// ✅ Good: Use result to prevent optimization
var result any  // Package-level to prevent inlining

func BenchmarkVM_Correct(b *testing.B) {
    var r any
    for i := 0; i < b.N; i++ {
        r, _ = machine.Run(bytecode, ctx)
    }
    result = r  // Assign to prevent optimization
}
```

### 5. Multiple Iterations

```bash
# ❌ Single run - may be unrepresentative
go test -bench=.

# ✅ Multiple runs with statistical analysis
go test -bench=. -count=10 > results.txt
benchstat results.txt
```

---

## Common Pitfalls

### Pitfall 1: Setup Inside Loop

```go
// ❌ Wrong: Setup included in timing
func BenchmarkVM_Wrong(b *testing.B) {
    for i := 0; i < b.N; i++ {
        node, _ := parser.ParseString("a + b")  // TIMED!
        // ... rest of benchmark
    }
}

// ✅ Correct: Setup excluded
func BenchmarkVM_Correct(b *testing.B) {
    node, _ := parser.ParseString("a + b")  // NOT TIMED
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // ... benchmark
    }
}
```

### Pitfall 2: Hidden Allocations

```go
// ❌ Hidden allocation (interface conversion)
func BenchmarkVM_Hidden(b *testing.B) {
    for i := 0; i < b.N; i++ {
        vm.Push(42)  // int → any = allocation!
    }
}

// ✅ Explicit type (no conversion)
func BenchmarkVM_NoAlloc(b *testing.B) {
    for i := 0; i < b.N; i++ {
        vm.Push(42.0)  // float64 → any (still allocates, but expected)
    }
}
```

Verify with:
```bash
go test -bench=. -benchmem
# Check B/op and allocs/op
```

### Pitfall 3: Shared State

```go
// ❌ Race condition in parallel benchmark
var counter int
func BenchmarkVM_Race(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            counter++  // RACE!
        }
    })
}

// ✅ Thread-local state
func BenchmarkVM_Safe(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        localCounter := 0
        for pb.Next() {
            localCounter++  // Safe
        }
    })
}
```

---

## Automation Scripts

### Continuous Benchmark Tracking

```bash
#!/bin/bash
# scripts/track_performance.sh

BENCHMARK="BenchmarkVM_Boolean_Current"
OUTPUT="benchmarks/history/$(date +%Y%m%d_%H%M%S).txt"

mkdir -p benchmarks/history

# Run benchmark
go test -bench=$BENCHMARK -count=10 > $OUTPUT

# Store commit info
git rev-parse HEAD >> $OUTPUT
git log -1 --oneline >> $OUTPUT

echo "Results saved to $OUTPUT"

# Compare with previous
PREV=$(ls -t benchmarks/history/*.txt | sed -n '2p')
if [ -f "$PREV" ]; then
    echo "Comparing with previous run:"
    benchstat $PREV $OUTPUT
fi
```

### Pre-Commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

# Run benchmarks before commit
go test -bench=. -benchtime=5s > /tmp/bench_results.txt

# Check for regressions
if grep -q "slower" /tmp/bench_results.txt; then
    echo "WARNING: Performance regression detected!"
    cat /tmp/bench_results.txt
    read -p "Commit anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi
```

---

## References

- [Go Blog: Benchmarking](https://go.dev/blog/benchmarks)
- [benchstat Documentation](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)
- [optimization-journey.md](optimization-journey.md) - Real optimization examples
- [profiling-guide.md](profiling-guide.md) - CPU profiling

**Last Updated:** October 17, 2025
