# Chapter 13: The Three-Stage API

> "Understanding UExL's internal pipeline gives you the vocabulary to optimize your integration, structure your code, and diagnose failures at the right level."

---

## 13.1 The Three Stages Revisited

Every UExL evaluation follows three stages:

```
Expression String
       │ parser.ParseString()
       ▼
   Abstract Syntax Tree (AST)
       │ compiler.New().Compile(ast)
       ▼
   ByteCode (constants + instructions)
       │ vm.New(libCtx).Run(bytecode, vars)
       ▼
   Result (any value)
```

In most production integrations, you want to decouple stage two from stage three:

- **Stage 1 (Parse)** — once, at startup or on expression change
- **Stage 2 (Compile)** — once, at startup or on expression change
- **Stage 3 (Run)** — many times per request, with different `vars`

The `uexl` package provides a high-level API that manages this lifecycle for you.

---

## 13.2 The `Env`-Centric API (Recommended)

The recommended integration path uses `*uexl.Env` and `*uexl.CompiledExpr`. These types encapsulate the three stages, pool VM instances for concurrency, and validate function calls at compile time.

### The Default environment

```go
import "github.com/maniartech/uexl"

// One-shot evaluation using the package-level Default env:
result, err := uexl.Eval("product.basePrice * 0.9", map[string]any{
    "product": map[string]any{"basePrice": 99.99},
})
```

`uexl.Eval` uses a singleton `*Env` pre-loaded with:
- `vm.Builtins` — all 14 built-in functions
- `vm.DefaultPipeHandlers` — all 13 default pipe handlers

This is the fastest path for scripts, CLIs, and low-volume evaluations.

### Compile once, run many

```go
// At startup or rule-change time — compile once:
compiled, err := uexl.Default().Compile("orders |filter: $item.total > threshold |: len($last)")
if err != nil {
    // Handle parse/compile error — expression is invalid
    log.Fatalf("bad expression: %v", err)
}

// Per request — run many times with different vars:
result, err := compiled.Eval(ctx, map[string]any{
    "orders":    getOrders(),
    "threshold": 100.0,
})
```

`CompiledExpr.Eval` is goroutine-safe — multiple goroutines can call it concurrently. The env maintains a `sync.Pool` of `*vm.VM` instances so each concurrent call borrows a VM, runs it, and returns it to the pool without allocating.

---

## 13.3 Creating Custom Environments

When your application needs custom functions or pipe handlers, create a named environment:

```go
import (
    "github.com/maniartech/uexl"
    "github.com/maniartech/uexl/vm"
    "math"
    "strings"
    "fmt"
)

// Build the ShopLogic environment once at startup:
shopEnv := uexl.DefaultWith(
    uexl.WithFunctions(uexl.Functions{
        "min": func(args ...any) (any, error) {
            if len(args) != 2 {
                return nil, fmt.Errorf("min expects 2 arguments")
            }
            a, aOk := args[0].(float64)
            b, bOk := args[1].(float64)
            if !aOk || !bOk {
                return nil, fmt.Errorf("min expects two numbers")
            }
            return math.Min(a, b), nil
        },
        "max": func(args ...any) (any, error) {
            if len(args) != 2 {
                return nil, fmt.Errorf("max expects 2 arguments")
            }
            a, aOk := args[0].(float64)
            b, bOk := args[1].(float64)
            if !aOk || !bOk {
                return nil, fmt.Errorf("max expects two numbers")
            }
            return math.Max(a, b), nil
        },
        "upper": func(args ...any) (any, error) {
            if len(args) != 1 {
                return nil, fmt.Errorf("upper expects 1 argument")
            }
            s, ok := args[0].(string)
            if !ok {
                return nil, fmt.Errorf("upper expects a string")
            }
            return strings.ToUpper(s), nil
        },
        "floor": func(args ...any) (any, error) {
            if len(args) != 1 {
                return nil, fmt.Errorf("floor expects 1 argument")
            }
            n, ok := args[0].(float64)
            if !ok {
                return nil, fmt.Errorf("floor expects a number")
            }
            return math.Floor(n), nil
        },
    }),
    uexl.WithGlobals(map[string]any{
        "TAX_RATE": 0.08,
    }),
)
```

`DefaultWith` creates a new env that includes all of `Default()`'s functions + pipes *plus* your additions. The default singleton is never mutated.

---

## 13.4 Compiling and Validating Expressions

### Compile

```go
compiled, err := shopEnv.Compile(expr)
```

`Compile` runs stages 1 and 2 and also **validates all function call sites** — if the expression calls an unknown function, you get an error here, at compile time, not at eval time:

```go
_, err := shopEnv.Compile("unknown_func(x)")
// err: compile error: unknown function "unknown_func" — not registered in this environment
```

This means you can validate user-submitted expression strings immediately when they are saved, without waiting for the first evaluation.

### Validate only

```go
err := shopEnv.Validate(expr)
// nil → expression is valid in this env
// error → parse or compile error
```

`Validate` does not return the `*CompiledExpr`. Use it to power expression editors, linters, or API validation endpoints.

### MustCompile for package-level constants

```go
// At package level — for expressions you own and know are valid:
var discountExpr = uexl.MustCompile(
    "product.basePrice * (1 - min(discount, maxDiscount))",
)
```

`MustCompile` panics on error. Use it only for constants that cannot be broken by user input. Never use it for user-submitted expressions.

---

## 13.5 Running Expressions

### With context and vars

```go
result, err := compiled.Eval(ctx, vars)
```

- `ctx` — a `context.Context` for cancellation and deadline enforcement. The VM checks the context at each opcode boundary. Pass `context.Background()` for unconstrained evaluation.
- `vars` — `map[string]any`. Keys become expression-visible variables. `nil` is treated as an empty map.

### Globals and per-call vars

An `Env` can carry **global variables** — values always present in every evaluation in that environment:

```go
shopEnv := uexl.DefaultWith(
    uexl.WithGlobals(map[string]any{
        "TAX_RATE": 0.08,
        "today":    time.Now().Format("2006-01-02"),
    }),
)

// Per-call: only product and customer — today and TAX_RATE come from globals
result, _ := compiled.Eval(ctx, map[string]any{
    "product":  product,
    "customer": customer,
})
```

Per-call vars shadow globals of the same name.

---

## 13.6 Introspection: What Variables Does an Expression Use?

```go
compiled, _ := shopEnv.Compile("product.basePrice * discount + TAX_RATE")
vars := compiled.Variables()
// => ["TAX_RATE", "discount", "product"]  (sorted, globals included)
```

Use `Variables()` to:
- Generate context schema documentation automatically
- Validate that required variables are present in the context before evaluation
- Build expression editors that show what variables an expression expects

---

## 13.7 Extending Environments

Environments are immutable. `Extend` creates a new env that inherits all functions, pipe handlers, and globals, then applies additional options on top:

```go
baseEnv := uexl.DefaultWith(
    uexl.WithFunctions(mathFunctions),
    uexl.WithGlobals(map[string]any{"TAX_RATE": 0.08}),
)

// Admin env adds privileged lookup functions:
adminEnv := baseEnv.Extend(
    uexl.WithFunctions(uexl.Functions{
        "lookupCustomerTier": lookupCustomerTierFn,
        "sendAlert":          sendAlertFn,
    }),
)
```

The base env is never mutated. Admin expressions run in `adminEnv`; public expressions run in `baseEnv`.

---

## 13.8 The Low-Level API (When You Need It)

For special cases — testing, tooling, or non-standard compilation — you can use the three packages directly:

```go
import (
    "github.com/maniartech/uexl/parser"
    "github.com/maniartech/uexl/compiler"
    "github.com/maniartech/uexl/vm"
)

// Stage 1: Parse
ast, err := parser.ParseString(expr)
if err != nil {
    // Returns *parser/errors.ParseErrors — structured error with line/column
    log.Printf("parse error: %v", err)
    return
}

// Stage 2: Compile
comp := compiler.New()
if err := comp.Compile(ast); err != nil {
    log.Printf("compile error: %v", err)
    return
}
bytecode := comp.ByteCode()

// Stage 3: Execute
machine := vm.New(vm.LibContext{
    Functions:    vm.Builtins,
    PipeHandlers: vm.DefaultPipeHandlers,
})
result, err := machine.Run(bytecode, contextVars)
```

> **NOTE:** The low-level API does not pool VM instances or validate function names at compile time. Prefer `Env.Compile` + `CompiledExpr.Eval` for production use.

---

## 13.9 Error Types

Errors propagate without wrapping:

| Stage | Error source | Check for |
|-------|-------------|-----------|
| Parse | `*parsererrors.ParseErrors` | `.Error()` gives all parser errors |
| Compile | `error` string from compiler | `strings.HasPrefix(err.Error(), "compile error:")` |
| Function validation | `error` from `Env.Compile` | `"unknown function"` in message |
| Runtime (VM) | `error` from `Run`/`Eval` | Pipe errors, type errors, division by zero |

```go
result, err := compiled.Eval(ctx, vars)
if err != nil {
    // err could be:
    // - context.DeadlineExceeded / context.Canceled  (timeout/cancel)
    // - a runtime error from the expression (type mismatch, nil deref, etc.)
    log.Printf("eval error: %v", err)
}
```

---

## 13.10 ShopLogic: The Full Integration Pattern

```go
package shoplogic

import (
    "context"
    "fmt"
    "math"
    "strings"

    "github.com/maniartech/uexl"
)

// ShopEnv is the global ShopLogic expression environment.
// Created once at startup; thread-safe for concurrent use.
var ShopEnv *uexl.Env

func init() {
    ShopEnv = uexl.DefaultWith(
        uexl.WithFunctions(shopFunctions()),
        uexl.WithGlobals(map[string]any{
            "TAX_RATE": 0.08,
        }),
    )
}

func shopFunctions() uexl.Functions {
    return uexl.Functions{
        "min":   minFn,
        "max":   maxFn,
        "floor": floorFn,
        "upper": upperFn,
        "lower": lowerFn,
    }
}

// PricingRule is a compiled pricing expression.
type PricingRule struct {
    expr *uexl.CompiledExpr
    name string
}

// NewPricingRule compiles and validates a pricing expression at creation time.
func NewPricingRule(name, expression string) (*PricingRule, error) {
    compiled, err := ShopEnv.Compile(expression)
    if err != nil {
        return nil, fmt.Errorf("pricing rule %q invalid: %w", name, err)
    }
    return &PricingRule{expr: compiled, name: name}, nil
}

// Apply evaluates the pricing rule for a product/customer context.
func (r *PricingRule) Apply(ctx context.Context, product, customer map[string]any, maxDiscount float64) (float64, error) {
    result, err := r.expr.Eval(ctx, map[string]any{
        "product":     product,
        "customer":    customer,
        "maxDiscount": maxDiscount,
    })
    if err != nil {
        return 0, fmt.Errorf("pricing rule %q failed: %w", r.name, err)
    }
    price, ok := result.(float64)
    if !ok {
        return 0, fmt.Errorf("pricing rule %q: expected number result, got %T", r.name, result)
    }
    return price, nil
}

// --- helpers ---

func minFn(args ...any) (any, error) {
    if len(args) != 2 { return nil, fmt.Errorf("min: want 2 args") }
    a, aOk := args[0].(float64)
    b, bOk := args[1].(float64)
    if !aOk || !bOk { return nil, fmt.Errorf("min: want numbers") }
    return math.Min(a, b), nil
}

func maxFn(args ...any) (any, error) {
    if len(args) != 2 { return nil, fmt.Errorf("max: want 2 args") }
    a, aOk := args[0].(float64)
    b, bOk := args[1].(float64)
    if !aOk || !bOk { return nil, fmt.Errorf("max: want numbers") }
    return math.Max(a, b), nil
}

func floorFn(args ...any) (any, error) {
    if len(args) != 1 { return nil, fmt.Errorf("floor: want 1 arg") }
    n, ok := args[0].(float64)
    if !ok { return nil, fmt.Errorf("floor: want number") }
    return math.Floor(n), nil
}

func upperFn(args ...any) (any, error) {
    if len(args) != 1 { return nil, fmt.Errorf("upper: want 1 arg") }
    s, ok := args[0].(string)
    if !ok { return nil, fmt.Errorf("upper: want string") }
    return strings.ToUpper(s), nil
}

func lowerFn(args ...any) (any, error) {
    if len(args) != 1 { return nil, fmt.Errorf("lower: want 1 arg") }
    s, ok := args[0].(string)
    if !ok { return nil, fmt.Errorf("lower: want string") }
    return strings.ToLower(s), nil
}
```

---

## 13.11 Summary

- The `Env` API is the recommended integration path: immutable, goroutine-safe, VM-pooled.
- `uexl.Eval()` is the one-liner for simple use; `Env.Compile()` + `CompiledExpr.Eval()` is the production pattern.
- `Compile` validates function names at compile time — catch unknown functions before the first evaluation.
- `DefaultWith` extends the singleton default env without mutating it.
- `Env.Extend` creates tenant- or role-scoped environments from a base.
- `WithGlobals` adds per-environment constants; per-call vars shadow globals.
- The low-level `parser` / `compiler` / `vm` packages are available for tooling but rare in production.

---

## Exercises

**13.1 — Recall.** Which statement is false?
a) `CompiledExpr.Eval` is goroutine-safe
b) `Env.Compile` validates unknown function names
c) `uexl.Default()` can be mutated with `WithFunctions`
d) Per-call vars shadow env globals of the same name

**13.2 — Apply.** Write a Go program that:
1. Creates an env with a custom `clamp(value, lo, hi)` function
2. Compiles the expression `clamp(rating, 1, 5)`
3. Runs it for `rating = -3` and `rating = 7.5`

**13.3 — Extend.** Design a multi-tenant system where each tenant has a different `TAX_RATE` global. Use `Env.Extend` to create per-tenant envs from a shared base and show how a compiled pricing rule runs correctly per tenant.
