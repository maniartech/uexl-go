# Chapter 14: Custom Functions and Pipes

> "UExL's built-in set is intentionally small. The real power is in what you add. This chapter shows you how to extend UExL safely and idiomatically."

---

## 14.1 The Extension Model

UExL has two extension mechanisms:

1. **Custom functions** — callable by name in expressions: `upper(name)`, `clamp(x, 0, 100)`
2. **Custom pipe handlers** — new pipe types callable as `arr |take: 5`, `items |normalize: $item`

Both are registered in `LibContext` (directly) or via `Env` options (`WithFunctions`, `WithPipeHandlers`). Both are per-environment — multiple envs can have different function/pipe sets.

---

## 14.2 Writing Custom Functions

A custom function has the signature:

```go
type VMFunction func(args ...any) (any, error)
```

Every argument arrives as `any`. Your function is responsible for type-asserting and validating them. The function **must never panic** — return `(nil, error)` for invalid inputs.

### Anatomy of a well-written function

```go
func clampFn(args ...any) (any, error) {
    // 1. Arity check
    if len(args) != 3 {
        return nil, fmt.Errorf("clamp expects 3 arguments (value, lo, hi), got %d", len(args))
    }
    // 2. Type check each argument
    value, ok := args[0].(float64)
    if !ok {
        return nil, fmt.Errorf("clamp: first argument must be a number, got %T", args[0])
    }
    lo, ok := args[1].(float64)
    if !ok {
        return nil, fmt.Errorf("clamp: second argument (lo) must be a number, got %T", args[1])
    }
    hi, ok := args[2].(float64)
    if !ok {
        return nil, fmt.Errorf("clamp: third argument (hi) must be a number, got %T", args[2])
    }
    // 3. Validate logical constraints
    if lo > hi {
        return nil, fmt.Errorf("clamp: lo (%v) must be <= hi (%v)", lo, hi)
    }
    // 4. Compute and return
    if value < lo { return lo, nil }
    if value > hi { return hi, nil }
    return value, nil
}
```

### Registering

```go
myEnv := uexl.DefaultWith(
    uexl.WithFunctions(uexl.Functions{
        "clamp": clampFn,
    }),
)
```

Or alongside the full ShopLogic function set:

```go
shopEnv := uexl.DefaultWith(
    uexl.WithFunctions(shopLogicFunctions()),
)

func shopLogicFunctions() uexl.Functions {
    return uexl.Functions{
        "min":    minFn,
        "max":    maxFn,
        "floor":  floorFn,
        "ceil":   ceilFn,
        "round":  roundFn,
        "abs":    absFn,
        "clamp":  clampFn,
        "upper":  upperFn,
        "lower":  lowerFn,
        "trim":   trimFn,
        "split":  splitFn,
        "number": numberFn,
    }
}
```

---

## 14.3 The ShopLogic Standard Function Library

Here is the complete production-ready implementation of the functions used throughout this book's ShopLogic examples:

```go
package shoplogic

import (
    "fmt"
    "math"
    "strconv"
    "strings"
)

func minFn(args ...any) (any, error) {
    a, b, err := twoFloats("min", args)
    if err != nil { return nil, err }
    return math.Min(a, b), nil
}

func maxFn(args ...any) (any, error) {
    a, b, err := twoFloats("max", args)
    if err != nil { return nil, err }
    return math.Max(a, b), nil
}

func floorFn(args ...any) (any, error) {
    n, err := oneFloat("floor", args)
    if err != nil { return nil, err }
    return math.Floor(n), nil
}

func ceilFn(args ...any) (any, error) {
    n, err := oneFloat("ceil", args)
    if err != nil { return nil, err }
    return math.Ceil(n), nil
}

func roundFn(args ...any) (any, error) {
    n, err := oneFloat("round", args)
    if err != nil { return nil, err }
    return math.Round(n), nil
}

func absFn(args ...any) (any, error) {
    n, err := oneFloat("abs", args)
    if err != nil { return nil, err }
    return math.Abs(n), nil
}

func clampFn(args ...any) (any, error) {
    if len(args) != 3 {
        return nil, fmt.Errorf("clamp expects 3 arguments, got %d", len(args))
    }
    v, vOk := args[0].(float64)
    lo, loOk := args[1].(float64)
    hi, hiOk := args[2].(float64)
    if !vOk || !loOk || !hiOk {
        return nil, fmt.Errorf("clamp expects numbers")
    }
    if lo > hi {
        return nil, fmt.Errorf("clamp: lo must be <= hi")
    }
    if v < lo { return lo, nil }
    if v > hi { return hi, nil }
    return v, nil
}

func upperFn(args ...any) (any, error) {
    s, err := oneString("upper", args)
    if err != nil { return nil, err }
    return strings.ToUpper(s), nil
}

func lowerFn(args ...any) (any, error) {
    s, err := oneString("lower", args)
    if err != nil { return nil, err }
    return strings.ToLower(s), nil
}

func trimFn(args ...any) (any, error) {
    s, err := oneString("trim", args)
    if err != nil { return nil, err }
    return strings.TrimSpace(s), nil
}

func replaceFn(args ...any) (any, error) {
    if len(args) != 3 {
        return nil, fmt.Errorf("replace expects 3 arguments (str, old, new)")
    }
    s, sOk := args[0].(string)
    old, oldOk := args[1].(string)
    newStr, newOk := args[2].(string)
    if !sOk || !oldOk || !newOk {
        return nil, fmt.Errorf("replace expects strings")
    }
    return strings.ReplaceAll(s, old, newStr), nil
}

func splitFn(args ...any) (any, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf("split expects 2 arguments (str, sep)")
    }
    s, sOk := args[0].(string)
    sep, sepOk := args[1].(string)
    if !sOk || !sepOk {
        return nil, fmt.Errorf("split expects strings")
    }
    parts := strings.Split(s, sep)
    result := make([]any, len(parts))
    for i, p := range parts {
        result[i] = p
    }
    return result, nil
}

// numberFn parses a string to float64; returns nil on failure (not an error).
func numberFn(args ...any) (any, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("number expects 1 argument")
    }
    switch v := args[0].(type) {
    case float64:
        return v, nil
    case string:
        n, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
        if err != nil {
            return nil, nil  // unparseable → null (not an error)
        }
        return n, nil
    case bool:
        if v { return 1.0, nil }
        return 0.0, nil
    case nil:
        return nil, nil
    default:
        return nil, nil
    }
}

// --- helpers ---

func oneFloat(name string, args []any) (float64, error) {
    if len(args) != 1 {
        return 0, fmt.Errorf("%s expects 1 argument, got %d", name, len(args))
    }
    n, ok := args[0].(float64)
    if !ok {
        return 0, fmt.Errorf("%s expects a number, got %T", name, args[0])
    }
    return n, nil
}

func twoFloats(name string, args []any) (float64, float64, error) {
    if len(args) != 2 {
        return 0, 0, fmt.Errorf("%s expects 2 arguments, got %d", name, len(args))
    }
    a, aOk := args[0].(float64)
    b, bOk := args[1].(float64)
    if !aOk || !bOk {
        return 0, 0, fmt.Errorf("%s expects numbers", name)
    }
    return a, b, nil
}

func oneString(name string, args []any) (string, error) {
    if len(args) != 1 {
        return "", fmt.Errorf("%s expects 1 argument, got %d", name, len(args))
    }
    s, ok := args[0].(string)
    if !ok {
        return "", fmt.Errorf("%s expects a string, got %T", name, args[0])
    }
    return s, nil
}
```

---

## 14.4 Thread Safety

Custom functions are shared across all concurrent evaluations in an environment. **They must be goroutine-safe.** The rules:

- **No shared mutable state** in the function closure. If you need shared state, protect it with a mutex or use atomic values.
- **No side effects on UExL values** — the `any` arguments passed to functions should not be mutated; UExL may reuse them internally.
- **Closures that capture Go variables** — the captured variables must be safe for concurrent reads (e.g., immutable after env construction).

```go
// SAFE: stateless pure function
"upper": func(args ...any) (any, error) {
    s, ok := args[0].(string)
    if !ok { return nil, fmt.Errorf("upper: string expected") }
    return strings.ToUpper(s), nil
},

// SAFE: read-only lookup into immutable map
rateTable := map[string]float64{"A": 0.05, "B": 0.10}  // frozen at startup
"rate": func(args ...any) (any, error) {
    tier, ok := args[0].(string)
    if !ok { return nil, fmt.Errorf("rate: string expected") }
    r, exists := rateTable[tier]
    if !exists { return nil, nil }
    return r, nil
},

// UNSAFE: mutates a shared counter without synchronization — DON'T DO THIS
counter := 0
"increment": func(args ...any) (any, error) {
    counter++  // DATA RACE — concurrent evals see non-atomic mutation
    return float64(counter), nil
},

// SAFE alternative: use atomic
var counter int64
"increment": func(args ...any) (any, error) {
    return float64(atomic.AddInt64(&counter, 1)), nil
},
```

---

## 14.5 Writing Custom Pipe Handlers

A custom pipe handler implements:

```go
type PipeHandler func(ctx uexl.PipeContext, input any) (any, error)
```

`PipeContext` provides two methods:
- `EvalItem(item any, index int) (any, error)` — sets `$item` and `$index`, runs the predicate
- `EvalWith(scope map[string]any) (any, error)` — sets arbitrary scope variables, runs the predicate

### Example: `take` pipe — first N elements

```go
func takePipeHandler(ctx uexl.PipeContext, input any) (any, error) {
    arr, ok := input.([]any)
    if !ok {
        return nil, fmt.Errorf("take: expected array input, got %T", input)
    }

    // EvalWith({}) runs the predicate with no extra scope vars —
    // the predicate is expected to be a literal number like |take: 5
    raw, err := ctx.EvalWith(map[string]any{})
    if err != nil {
        return nil, fmt.Errorf("take: cannot evaluate count: %w", err)
    }
    n, ok := raw.(float64)
    if !ok || n < 0 {
        return nil, fmt.Errorf("take: count must be a non-negative number, got %v", raw)
    }

    limit := int(n)
    if limit > len(arr) {
        limit = len(arr)
    }
    // Return a slice of the original array — no copy needed
    result := make([]any, limit)
    copy(result, arr[:limit])
    return result, nil
}
```

### Example: `skip` pipe — skip first N elements

```go
func skipPipeHandler(ctx uexl.PipeContext, input any) (any, error) {
    arr, ok := input.([]any)
    if !ok {
        return nil, fmt.Errorf("skip: expected array input, got %T", input)
    }
    raw, err := ctx.EvalWith(map[string]any{})
    if err != nil {
        return nil, fmt.Errorf("skip: cannot evaluate count: %w", err)
    }
    n, ok := raw.(float64)
    if !ok || n < 0 {
        return nil, fmt.Errorf("skip: count must be a non-negative number, got %v", raw)
    }
    offset := int(n)
    if offset >= len(arr) {
        return []any{}, nil
    }
    result := make([]any, len(arr)-offset)
    copy(result, arr[offset:])
    return result, nil
}
```

### Registering custom pipes

```go
allPipes := vm.DefaultPipeHandlers
allPipes["take"] = takePipeHandler
allPipes["skip"] = skipPipeHandler

myEnv := uexl.NewEnv(
    uexl.WithFunctions(vm.Builtins),
    uexl.WithPipeHandlers(allPipes),
)
```

> **NOTE:** `vm.DefaultPipeHandlers` is a map — assigning to it directly mutates the default. Prefer copying it first:
> ```go
> allPipes := make(uexl.PipeHandlers)
> for k, v := range vm.DefaultPipeHandlers {
>     allPipes[k] = v
> }
> allPipes["take"] = takePipeHandler
> ```

---

## 14.6 The `Lib` Interface for Bundled Extensions

For reusable function+pipe bundles, implement the `Lib` interface:

```go
type Lib interface {
    Apply(cfg *EnvConfig)
}
```

This lets you package a cohesive library and register it in one call:

```go
type MathLib struct{}

func (l MathLib) Apply(cfg *uexl.EnvConfig) {
    cfg.RegisterFunctions(uexl.Functions{
        "min":   minFn,
        "max":   maxFn,
        "floor": floorFn,
        "ceil":  ceilFn,
        "round": roundFn,
        "abs":   absFn,
        "clamp": clampFn,
    })
}

type StringLib struct{}

func (l StringLib) Apply(cfg *uexl.EnvConfig) {
    cfg.RegisterFunctions(uexl.Functions{
        "upper":   upperFn,
        "lower":   lowerFn,
        "trim":    trimFn,
        "replace": replaceFn,
        "split":   splitFn,
    })
}

// Usage:
shopEnv := uexl.DefaultWith(
    uexl.WithLib(MathLib{}),
    uexl.WithLib(StringLib{}),
)
```

---

## 14.7 Function Naming Conventions

| Pattern | Example | Use for |
|---------|---------|---------|
| Verb | `trim`, `upper`, `round` | Pure transforms |
| Verb + noun | `parseInt`, `formatDate` | Conversions with domain context |
| camelCase | `startsWith`, `isValid` | Multi-word predicates |
| No spaces, no special chars | required | Tokenizer constraint |

Do not register functions whose names shadow built-ins (`len`, `str`, `contains`, etc.) — doing so overrides them silently. The compiler resolves function names from the registry; the last-registered name wins.

---

## 14.8 Summary

- Custom functions implement `func(args ...any) (any, error)`. Validate arity and types; never panic.
- Custom pipes implement `func(ctx PipeContext, input any) (any, error)`. Use `EvalItem` for per-element or `EvalWith` for arbitrary scope.
- Both functions and pipes are goroutine-safe by contract — use only immutable closures or atomic state.
- Register via `WithFunctions` / `WithPipeHandlers` on `Env`, or bundle in a `Lib` for reuse.
- Copy `vm.DefaultPipeHandlers` before adding custom pipes — do not mutate the package-level default.
- Do not shadow built-in function names.

---

## Exercises

**14.1 — Recall.** A custom function receives a `nil` argument (the user passed `null` from the expression). What Go type does `args[0]` have? How should you handle it?

**14.2 — Apply.** Implement a `formatCurrency(amount, currencyCode)` function that returns a string like `"USD 99.99"`. Register it in `shopEnv` and write a UExL expression that formats `product.basePrice` with `currencyCode`.

**14.3 — Extend.** Implement a custom `|page:` pipe handler that accepts two parameters (`offset` and `size`) expressed as `arr |page: {offset: 1, size: 10}`. The pipe should skip `offset * size` elements and return the next `size` elements. Show the Go implementation, the registration, and a UExL expression that pages through `products`.
