# Chapter 17: Performance Tuning

> "The biggest performance win in UExL costs nothing to implement: compile once, run many times. Everything else is a fine-tuning exercise."

---

## 17.1 The Performance Hierarchy

UExL performance optimizations, in order of impact:

1. **Compile once, evaluate many times** — eliminates parse and compile overhead on every call
2. **Use env globals for constants** — eliminates per-call context copying for invariants
3. **Design flat, minimal context maps** — reduces property lookup depth and map allocation cost
4. **Use pipes for bulk operations** — batch processing in a single expression pass
5. **Move heavy computation to Go** — expressions are decision logic, not number crunching
6. **Profile before optimizing** — measure actual bottlenecks, not assumed ones

---

## 17.2 Compile Once, Evaluate Many Times

Parsing and compilation are relatively expensive operations that scan, tokenize, parse, and compile the expression to bytecode. Evaluation against pre-compiled bytecode is the fast path.

**Anti-pattern (one-shot eval on every request):**

```go
// SLOW: re-parses and re-compiles on every rule evaluation
func applyDiscount(product, customer map[string]any) (any, error) {
    return uexl.Eval(
        "product.basePrice * (customer.tier == 'platinum' ? 0.85 : 1.0)",
        map[string]any{"product": product, "customer": customer},
    )
}
```

**Correct pattern (compile at startup):**

```go
var discountExpr *uexl.CompiledExpr

func init() {
    var err error
    discountExpr, err = shopEnv.Compile(
        "product.basePrice * (customer.tier == 'platinum' ? 0.85 : 1.0)",
    )
    if err != nil {
        panic(err) // startup failure — correct place to fail fast
    }
}

// FAST: only VM execution on each call
func applyDiscount(ctx context.Context, product, customer map[string]any) (any, error) {
    return discountExpr.Eval(ctx, map[string]any{
        "product":  product,
        "customer": customer,
    })
}
```

`CompiledExpr.Eval` is goroutine-safe — you can call it from thousands of goroutines simultaneously without locks. The `Env` internally pools `*vm.VM` instances so each goroutine gets a clean VM without allocation on most calls.

---

## 17.3 Goroutine Safety: What You Get For Free

`CompiledExpr` is safe for concurrent use without any locking on your side:

- The compiled bytecode (constants + instructions) is **read-only** at evaluation time
- The VM needed to run bytecode is provided from a `sync.Pool` per `Env`
- The context vars map you pass is **not** modified by UExL

This means you can share a single `*uexl.CompiledExpr` across all goroutines and HTTP handlers:

```go
// Shared across all goroutines — correct
var (
    priceRule    = mustCompile(shopEnv, "product.basePrice * discount")
    eligibility  = mustCompile(shopEnv, "customer.active && customer.totalSpent > 500")
)

func mustCompile(env *uexl.Env, expr string) *uexl.CompiledExpr {
    c, err := env.Compile(expr)
    if err != nil {
        panic(fmt.Sprintf("invalid expression %q: %v", expr, err))
    }
    return c
}
```

---

## 17.4 Using Globals for Constants

Constants that do not change between calls belong in the env, not the context map. Each key in the context map costs a lookup at evaluation time; globals cached in the VM are looked up via a pre-resolved index.

```go
// Less efficient: pass constants as context vars every call
vars := map[string]any{
    "TAX_RATE":     0.08,
    "product":      productMap,
    "customer":     customerMap,
}

// More efficient: constants as env globals, dynamic data as context
shopEnv := uexl.DefaultWith(
    uexl.WithGlobals(map[string]any{
        "TAX_RATE": 0.08,
    }),
)
// Then per call:
vars := map[string]any{
    "product":  productMap,
    "customer": customerMap,
}
```

Globals are resolved at compile time — the compiler bakes their indices into the bytecode. Context vars are resolved at runtime by map lookup, which has a small but measurable cost when the map is large.

---

## 17.5 Profiling Expression Performance

When you need to identify the bottleneck in your expression pipeline, use Go's built-in profiling tools.

### Running a CPU profile

```bash
go test ./... -bench=BenchmarkVM_Boolean_Current -benchtime=20s -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

In the pprof shell:
```
(pprof) top 10
(pprof) list uexl
```

### Writing your own micro-benchmark

```go
func BenchmarkMyRule(b *testing.B) {
    env := uexl.Default()
    compiled, _ := env.Compile("customer.totalSpent > 500 && customer.active")
    vars := map[string]any{
        "customer": map[string]any{
            "totalSpent": 750.0,
            "active":     true,
        },
    }
    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = compiled.Eval(ctx, vars)
    }
}
```

Run with:
```bash
go test -bench=BenchmarkMyRule -benchmem -benchtime=10s
```

Look at both `ns/op` (latency) and `allocs/op` (allocation pressure). Zero allocations is achievable for simple boolean and arithmetic expressions.

---

## 17.6 Context Map Optimization

### Keep context maps flat where possible

Nested map lookups (`product.basePrice`) require two map accesses under the hood. For extremely tight loops, flatten the structure:

```go
// Nested (two lookups for product.basePrice)
vars := map[string]any{
    "product": map[string]any{
        "basePrice": p.BasePrice,
    },
}

// Flat (one lookup — use for hot paths when worth the API cost)
vars := map[string]any{
    "basePrice": p.BasePrice,
    "category":  p.Category,
}
// Expression: basePrice * (category == 'electronics' ? 1.15 : 1.0)
```

### Reuse map allocations in tight loops

For benchmark-grade workloads, reuse the outer map and only update values that change:

```go
vars := map[string]any{
    "product":  nil,
    "customer": nil,
}

for _, pair := range pairs {
    vars["product"] = pair.ProductMap
    vars["customer"] = pair.CustomerMap
    result, _ := compiled.Eval(ctx, vars)
    _ = result
}
```

This avoids the map allocation on every iteration.

---

## 17.7 When to Move Logic to Go

Expressions excel at decision logic: comparisons, boolean combinations, property access, and simple arithmetic. They are **not** optimized for:

- Heavy numeric computation (FFT, matrix math, statistical aggregations)
- Large dataset processing (sort 100,000 items in an expression)
- Recursive algorithms
- I/O or side effects

When an expression becomes the bottleneck because it is doing too much work, move that work to either a pre-computed context variable or a registered host function.

**Before (expression does heavy work):**
```uexl
orders |filter: $item.status == 'complete'
       |map: $item.total
       |reduce: ($acc ?? 0) + $item
```

**After (pre-compute in Go — especially for large datasets):**
```go
// In Go
totalCompleted := sumCompleteOrders(orders)

// In expression (now trivial)
vars["totalCompleted"] = totalCompleted
// Expression: totalCompleted > threshold
```

The `|map: |reduce:` chain works correctly and efficiently for small arrays (under a few hundred items). For thousands of items per evaluation, prefer Go.

---

## 17.8 Summary

- **Compile once at startup, evaluate at runtime** — this is the single most impactful optimization.
- `CompiledExpr.Eval` is goroutine-safe; share `*CompiledExpr` across all goroutines.
- Use `WithGlobals` for constants that do not change between calls.
- Keep context maps small and reasonably flat.
- Reuse context map allocations if garbage collector pressure appears in profiles.
- Profile with `-cpuprofile` before making micro-optimization decisions.
- Pipe operators are fast for small arrays but not a replacement for Go-level bulk data processing.
- Expressions are decision logic — keep them that way.

---

## Exercises

**17.1 — Recall.** Why is `uexl.Eval(expr, vars)` inappropriate for production use inside an HTTP handler? What should be used instead?

**17.2 — Apply.** You have a benchmark showing `1200 ns/op` for a rule that calls `uexl.Eval` on every invocation. Rewrite the benchmark and the calling code to use a `*uexl.CompiledExpr`. What `ns/op` range would you expect for a simple boolean expression after this change?

**17.3 — Extend.** A ShopLogic pricing rule runs 10,000 times per second during peak sales. It has 12 context variables, 3 of which are global constants (TAX_RATE, MAX_DISCOUNT, FREE_SHIPPING_THRESHOLD). Describe the full optimization strategy — what goes into env globals, what goes into per-call context, and how you structure the compilation lifecycle.
