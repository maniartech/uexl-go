# UExL — Universal Expression Language for Go

**An embeddable, platform-independent expression evaluation engine with zero allocations on the hot path.**

UExL turns runtime strings into evaluated results — no codegen, no `eval()`, no panics. It is perfect for applications where expressions are not known at compile time, or where you need flexibility through configuration files or databases. Pre-compile once and re-evaluate faster than cel-go and expr — with zero GC pressure on the hot path. With comprehensive pipe operations and Excel-compatible syntax, UExL makes dynamic expression evaluation both powerful and production-ready. It ships as a single import with a clean, environment-first API designed for production Go code.

```go
// One line for scripts and REPLs.
result, err := uexl.Eval("price * qty * (1 - discount)", vars)

// Pre-compile for hot paths — goroutine-safe, pool-backed, zero extra allocs.
expr := uexl.MustCompile("price * qty * (1 - discount)")
result, err := expr.Eval(ctx, vars)

// Transform collections with chainable pipes.
result, err := uexl.Eval("[1,2,3,4,5] |filter: $item > 2 |map: $item * 10", nil)
// → [30, 40, 50]
```

---

## Table of Contents

- [Why UExL?](#why-uexl)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Core Concepts](#core-concepts)
  - [The Three-Stage Pipeline](#the-three-stage-pipeline)
  - [Environments](#environments)
  - [Compiled Expressions](#compiled-expressions)
  - [Pipe Operators](#pipe-operators)
  - [Result Helpers](#result-helpers)
- [API Reference](#api-reference)
  - [Package-Level Shortcuts](#package-level-shortcuts)
  - [Env — Evaluation Environment](#env--evaluation-environment)
  - [CompiledExpr — Pre-compiled Expression](#compiledexpr--pre-compiled-expression)
  - [Extensibility — Functions, Pipes, and Libs](#extensibility--functions-pipes-and-libs)
  - [Result Type Helpers](#result-type-helpers)
- [Performance](#performance)
- [Operator Reference](#operator-reference)
- [Language Examples](#language-examples)

---

## Why UExL?

| Concern | UExL’s answer |
|---------|---------------|
| **Zero panics** | Every error path returns an `error`; the VM never panics |
| **Goroutine safety** | `Env` and `CompiledExpr` are immutable after construction; VMs are pool-borrowed |
| **Zero allocations** | 0 allocs/op for primitive evaluation — no GC pressure in hot paths |
| **Explicit semantics** | `??` (nullish coalescing), `?.` (optional chaining) — no JavaScript-style surprises |
| **Extensible** | Register custom functions, pipe handlers, and global vars per environment |
| **One import** | All public types re-exported from the root package; no sub-package imports needed |

---

## Installation

> UExL is not yet stable-released. The API below reflects the current `main` branch.

```bash
go get github.com/maniartech/uexl
```

```go
import "github.com/maniartech/uexl"
```

---

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/maniartech/uexl"
)

func main() {
    // ── One-liner evaluation ───────────────────────────────────────────────────
    result, err := uexl.Eval("10 + 20", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result) // 30

    // ── With variables ───────────────────────────────────────────────────────────
    result, err = uexl.Eval("price * qty", map[string]any{
        "price": 4.99,
        "qty":   3.0,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result) // 14.97

    // ── Pre-compiled for repeated evaluation (goroutine-safe) ──────────────────────
    expr := uexl.MustCompile("price * qty * (1 - discount)")

    for _, order := range orders {
        result, err := expr.Eval(context.Background(), map[string]any{
            "price":    order.Price,
            "qty":      order.Qty,
            "discount": order.Discount,
        })
        if err != nil {
            log.Println(err)
            continue
        }
        fmt.Println(result)
    }

    // ── Pipe transformations ──────────────────────────────────────────────────
    result, err = uexl.Eval(
        `scores |filter: $item >= 60 |map: $item * 1.1`,
        map[string]any{"scores": []any{45.0, 62.0, 78.0, 91.0}},
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result) // [68.2, 85.8, 100.1]
}
```

---

## Core Concepts

### The Three-Stage Pipeline

Every expression travels through three stages:

```
Expression string  ►►  Parser (AST)  ►►  Compiler (Bytecode)  ►►  VM (Result)
```

`uexl.Eval` runs all three stages for you. `Env.Compile` stops after stage two and hands you a `*CompiledExpr` you can evaluate thousands of times without re-parsing or re-compiling.

### Environments

An `Env` bundles together the functions, pipe handlers, and global variables available during expression evaluation. It is **immutable after construction** and **safe to share across goroutines**.

```go
// Blank environment — no stdlib, no pipes, no globals.
env := uexl.NewEnv()

// Default environment — stdlib functions + built-in pipe handlers.
env := uexl.Default()

// Extend Default with extra functions; Default() itself is unchanged.
env := uexl.DefaultWith(
    uexl.WithFunctions(uexl.Functions{
        "discount": func(args ...any) (any, error) {
            price, _ := args[0].(float64)
            rate, _  := args[1].(float64)
            return price * (1 - rate), nil
        },
    }),
)

// Build a fully custom environment.
env := uexl.NewEnv(
    uexl.WithFunctions(myFunctions),
    uexl.WithPipeHandlers(myPipeHandlers),
    uexl.WithGlobals(map[string]any{"version": "2.0", "maxItems": 100.0}),
)

// Extend any env — copy-on-write, original is untouched.
child := env.Extend(uexl.WithGlobals(map[string]any{"tenantID": "acme"}))
```

### Compiled Expressions

`*CompiledExpr` is the workhorse for production use. Compile once, evaluate anywhere, concurrently:

```go
expr, err := uexl.Default().Compile("base ** 2 + offset")
if err != nil {
    // parse + compile + function validation errors reported here
    log.Fatal(err)
}

// Introspect the variables the expression needs.
fmt.Println(expr.Variables()) // ["base", "offset"]

// Evaluate concurrently — each call borrows a pooled VM.
var wg sync.WaitGroup
for i := 0; i < 1000; i++ {
    wg.Add(1)
    go func(n float64) {
        defer wg.Done()
        result, _ := expr.Eval(ctx, map[string]any{"base": n, "offset": 5.0})
        _ = result
    }(float64(i))
}
wg.Wait()
```

Function validation happens at **compile time**, not eval time — if you reference `unknownFn()` and it is not registered in the env, `Compile` returns an error immediately.

### Pipe Operators

Pipes transform data through chainable stages using ephemeral scope variables (`$item`, `$index`, `$acc`, etc.):

```
input  |pipeName: predicate  |pipeName: predicate  ...
```

**Built-in pipes:**

| Pipe | Scope vars | Description |
|------|-----------|-------------|
| `\|map:` | `$item`, `$index` | Transform each element |
| `\|filter:` | `$item`, `$index` | Keep elements where predicate is truthy |
| `\|reduce:` | `$acc`, `$item`, `$index`, `$last` | Fold to a single value |
| `\|find:` | `$item`, `$index` | First matching element or null |
| `\|some:` | `$item`, `$index` | True if any element matches |
| `\|every:` | `$item`, `$index` | True if all elements match |
| `\|sort:` | `$item`, `$index` | Sort by predicate result |
| `\|unique:` | `$item`, `$index` | Deduplicate by predicate result |
| `\|groupBy:` | `$item`, `$index` | Group into `map[string][]any` |
| `\|chunk:` | `$chunk`, `$index` | Split into fixed-size sub-arrays |
| `\|window:` | `$window`, `$index` | Sliding window sub-arrays |
| `\|:` | `$item` | Passthrough / single-value transform |

```go
// Chained pipes — each stage feeds the next.
result, _ := uexl.Eval(
    `orders
        |filter:  $item.status == "shipped"
        |map:     $item.total
        |reduce:  ($acc || 0) + $item`,
    map[string]any{"orders": orders},
)
```

### Result Helpers

UExL evals return `any`. Use the typed helpers to convert without surprises — no truthy coercion, no `fmt.Sprint` fallbacks:

```go
result, err := uexl.Eval("price * qty", vars)

total, err := uexl.AsFloat64(result)  // float64, int, int64, float32 → float64
ok,    err := uexl.AsBool(result)     // bool only; AsBool(1) → error
label, err := uexl.AsString(result)   // string only; no fmt.Sprint
items, err := uexl.AsSlice(result)    // []any only
meta,  err := uexl.AsMap(result)      // map[string]any only
```

---

## API Reference

### Package-Level Shortcuts

These convenience functions use the `Default()` singleton env with `context.Background()`:

```go
// Evaluate expr with optional variables. vars may be nil.
result, err := uexl.Eval(expr string, vars map[string]any) (any, error)

// Parse, compile, and validate expr against the default stdlib environment.
err := uexl.Validate(expr string) error

// Compile expr against Default(); panics on failure.
// Use at package level for known-valid expressions only.
ce := uexl.MustCompile(expr string) *CompiledExpr

// Return a new Env extending Default() with additional options.
env := uexl.DefaultWith(opts ...Option) *Env
```

### Env — Evaluation Environment

```go
// Default returns the stdlib singleton (sync.Once, goroutine-safe).
env := uexl.Default() *Env

// NewEnv builds a blank environment, then applies opts left-to-right.
env := uexl.NewEnv(opts ...Option) *Env

// Extend creates a child environment inheriting all symbols, then applying opts.
// The receiver is never mutated.
child := env.Extend(opts ...Option) *Env

// Compile parses, compiles, and validates expr against this env.
// Function calls are validated at compile time.
ce, err := env.Compile(expr string) (*CompiledExpr, error)

// MustCompile is a panic-on-error variant of Compile.
ce := env.MustCompile(expr string) *CompiledExpr

// Validate is a thin wrapper: compile, discard the result, return any error.
err := env.Validate(expr string) error

// Eval is a convenience: Compile + CompiledExpr.Eval in one call.
result, err := env.Eval(ctx context.Context, expr string, vars map[string]any) (any, error)

// Info returns a sorted, copy-independent snapshot of registered symbols.
info := env.Info() EnvInfo  // Info.String() produces a human-readable summary

// Introspect individual registries.
env.HasFunction("len")  bool
env.HasPipe("map")      bool
env.HasGlobal("rate")   bool
```

### CompiledExpr — Pre-compiled Expression

```go
// Eval executes the bytecode, honoring ctx for cancellation/deadline.
// Borrows a VM from the env's pool — goroutine-safe, near-zero allocation.
result, err := ce.Eval(ctx context.Context, vars map[string]any) (any, error)

// Variables returns sorted variable names the expression references.
// Returns []string{} (never nil) for constant expressions. Copy-safe.
names := ce.Variables() []string

// Env returns the *Env the expression was compiled against.
env := ce.Env() *Env
```

### Extensibility — Functions, Pipes, and Libs

**Options** (all panic on nil input to surface misconfigurations early):

```go
uexl.WithFunctions(fns Functions) Option
uexl.WithPipeHandlers(pipes PipeHandlers) Option
uexl.WithGlobals(vars map[string]any) Option
uexl.WithLib(lib Lib) Option  // Lib bundles all three into one Apply() call
```

**Writing a custom function:**

```go
tax := func(args ...any) (any, error) {
    price, ok := args[0].(float64)
    if !ok {
        return nil, fmt.Errorf("tax: expected float64, got %T", args[0])
    }
    return price * 1.08, nil
}

env := uexl.NewEnv(uexl.WithFunctions(uexl.Functions{"tax": tax}))
result, _ := env.Eval(ctx, "tax(price)", map[string]any{"price": 100.0})
// → 108.0
```

**Writing a custom pipe:**

```go
double := func(ctx uexl.PipeContext, input any) (any, error) {
    arr, ok := input.([]any)
    if !ok {
        return nil, fmt.Errorf("double: expected array")
    }
    out := make([]any, len(arr))
    for i, v := range arr {
        out[i] = v.(float64) * 2
    }
    return out, nil
}

env := uexl.NewEnv(uexl.WithPipeHandlers(uexl.PipeHandlers{"double": double}))
result, _ := env.Eval(ctx, "[1,2,3] |double:", nil)
// → [2, 4, 6]
```

**Shipping a reusable library (the `Lib` interface):**

```go
type FinanceLib struct{}

func (FinanceLib) Apply(cfg *uexl.EnvConfig) {
    cfg.AddFunctions(uexl.Functions{
        "pv":  presentValue,
        "fv":  futureValue,
        "irr": internalRateOfReturn,
    })
    cfg.AddGlobals(map[string]any{
        "defaultRate": 0.05,
    })
}

env := uexl.NewEnv(uexl.WithLib(FinanceLib{}))
```

### Result Type Helpers

```go
uexl.AsFloat64(v any) (float64, error)       // float64 | int | int64 | float32 → float64
uexl.AsBool(v any)    (bool, error)           // bool only — no truthy coercion
uexl.AsString(v any)  (string, error)         // string only — no fmt.Sprint
uexl.AsSlice(v any)   ([]any, error)          // []any only
uexl.AsMap(v any)     (map[string]any, error) // map[string]any only
```

All helpers return a descriptive error for `nil` input and for mismatched types.

---

## Performance

UExL uses a zero-allocation `Value` type for all runtime data. Benchmarked head-to-head against [expr](https://github.com/antonmedv/expr) and [cel-go](https://github.com/google/cel-go) on the same hardware:

### Benchmark Comparison

All results measured with `-benchtime=10s -benchmem` on AMD Ryzen 7 5700G (Windows/amd64).
UExL uses a pre-compiled expression with a pooled VM (same as `CompiledExpr.Eval()`). Competitors use their equivalent pre-compiled hot paths.
**Reproduce it:** clone [golang-expression-evaluation-comparison](https://github.com/antonmedv/golang-expression-evaluation-comparison), add the `uexl_test.go` from this repo, and run `go test -bench=. -benchmem -benchtime=10s`.
| Scenario | expr | cel-go | **UExL** |
|----------|:----:|:------:|:--------:|
| Boolean expression | 262 ns \| 1 alloc | 335 ns \| 1 alloc | **223 ns \| 0 allocs** |
| String match | 532 ns \| 4 allocs | 535 ns \| 4 allocs | **119 ns \| 0 allocs** |
| Custom function call | 367 ns \| 4 allocs | 397 ns \| 4 allocs | **197 ns \| 2 allocs** |
| Map over 100 items | 21,071 ns \| 111 allocs | 80,574 ns \| 621 allocs | **16,425 ns \| 104 allocs** |

In these benchmarks, UExL outperforms both frameworks across every scenario and is the only one with zero allocations on the boolean and string-matching paths.

> **Note:** `Eval()` (parse + compile + run in one call) costs ~10,200 ns/op and 64 allocs. Always pre-compile with `MustCompile()` or `Compile()` for repeated use.

### Pipe Performance

Pipe operations are measured at the VM layer. A `|map:` over 100 numeric items runs at approximately **16,400 ns/op** — 28% faster than expr and 5× faster than cel-go — with 104 allocs vs 111 (expr) and 621 (cel-go).

`CompiledExpr.Eval` borrows a `*vm.VM` from a `sync.Pool` on each call and returns it via `defer`. In a steady-state concurrent workload the pool stays warm, meaning most evaluations pay zero VM allocation cost.

### Performance Deep-Dives

- [Value type architecture](designdocs/value-migration/README.md) — how zero allocations were achieved
- [Optimization journey](designdocs/value-migration/performance-optimization.md) — measurements and decisions
- [Planned improvements](designdocs/performance/README.md) — targeting 20–35 ns/op

---

## Operator Reference

| Operators | Type | Associativity | Notes |
|-----------|------|---------------|-------|
| `(` `)` | Grouping | — | |
| `.` `[]` `?.` `?.[` | Member access | Left | Optional variants guard against null base |
| `**` `^` | Power | Right | `^` Excel-style, `**` Python/JS-style |
| `-` `!` `~` *(unary)* | Unary | Right | Negate, logical NOT, bitwise NOT |
| `%` | Modulus | Left | |
| `*` `/` | Multiplicative | Left | |
| `+` `-` | Additive | Left | `+` concatenates strings |
| `<<` `>>` | Bitwise shift | Left | |
| `??` | Nullish coalescing | Left | Falls back only on `null`; preserves `0`, `""`, `false` |
| `<` `>` `<=` `>=` | Relational | Left | |
| `==` `!=` `<>` | Equality | Left | `<>` Excel-style not-equals |
| `&` | Bitwise AND | Left | |
| `~` *(binary)* | Bitwise XOR | Left | Lua-style |
| `\|` *(binary)* | Bitwise OR | Left | |
| `&&` | Logical AND | Left | Short-circuits |
| `\|\|` | Logical OR | Left | Short-circuits |
| `?:` | Ternary | Right | `condition ? then : else` |
| `\|name:` | Pipe | Left | `\|map:`, `\|filter:`, `\|reduce:`, … |

---

## Language Examples

```go
// ── Arithmetic ───────────────────────────────────────────────────────────────────────
uexl.Eval("10 + 20", nil)           // 30
uexl.Eval("2 ** 10", nil)           // 1024  (Python/JS style)
uexl.Eval("2 ^ 10", nil)            // 1024  (Excel style)
uexl.Eval("17 % 5", nil)            // 2

// ── Strings ──────────────────────────────────────────────────────────────────────────
uexl.Eval(`"Hello, " + name + "!"`, map[string]any{"name": "World"})
uexl.Eval(`len("hello")`, nil)      // 5
uexl.Eval(`upper("uexl")`, nil)     // "UEXL"

// ── Nullish coalescing ───────────────────────────────────────────────────────────
uexl.Eval("config ?? 'default'", map[string]any{"config": nil}) // "default"
uexl.Eval("count  ?? 99",        map[string]any{"count": 0.0})  // 0   (not 99!)

// ── Optional chaining ────────────────────────────────────────────────────────────
uexl.Eval("user?.address?.city", map[string]any{"user": nil})   // null, no error

// ── Conditionals ─────────────────────────────────────────────────────────────────
uexl.Eval("score >= 60 ? 'pass' : 'fail'", map[string]any{"score": 72.0})

// ── Pipe: map ────────────────────────────────────────────────────────────────────────
uexl.Eval("[1,2,3] |map: $item * $item", nil)  // [1, 4, 9]

// ── Pipe: filter ─────────────────────────────────────────────────────────────────
uexl.Eval("[1,2,3,4,5] |filter: $item % 2 == 0", nil)  // [2, 4]

// ── Pipe: reduce ─────────────────────────────────────────────────────────────────
uexl.Eval("[1,2,3,4,5] |reduce: ($acc || 0) + $item", nil)  // 15

// ── Pipe: chained ─────────────────────────────────────────────────────────────────
uexl.Eval("[10,20,30,40] |filter: $item > 15 |map: $item / 10", nil)  // [2, 3, 4]

// ── Custom environment ────────────────────────────────────────────────────────────
env := uexl.NewEnv(
    uexl.WithFunctions(uexl.Functions{
        "celsius": func(args ...any) (any, error) {
            f, _ := args[0].(float64)
            return (f - 32) * 5 / 9, nil
        },
    }),
)
env.Eval(ctx, "celsius(temp)", map[string]any{"temp": 212.0})  // 100

// ── Introspection ───────────────────────────────────────────────────────────────────
expr, _ := uexl.Default().Compile("price * qty - discount")
fmt.Println(expr.Variables())  // [discount price qty]
fmt.Println(uexl.Default().Info())
// Env:
//   Functions (42): abs, ceil, contains, ...
//   PipeHandlers (13): chunk, every, filter, find, groupBy, map, ...
//   Globals (0):
```
