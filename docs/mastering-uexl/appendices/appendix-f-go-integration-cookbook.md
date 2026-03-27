# Appendix F: Go Integration Cookbook

This appendix collects ready-to-use Go code patterns for embedding UExL.

---

## F.1 Creating the Default Environment

```go
import "github.com/shoplogic/uexl-go"

// The simplest setup: built-ins + all default pipes
env := uexl.Default()
```

---

## F.2 Adding Custom Functions

```go
env := uexl.DefaultWith(
    uexl.WithFunctions(map[string]vm.Function{
        "upper": func(args []any) (any, error) {
            if len(args) != 1 {
                return nil, fmt.Errorf("upper expects 1 argument")
            }
            s, ok := args[0].(string)
            if !ok {
                return nil, fmt.Errorf("upper: argument must be a string, got %T", args[0])
            }
            return strings.ToUpper(s), nil
        },
        "lower": func(args []any) (any, error) {
            if len(args) != 1 {
                return nil, fmt.Errorf("lower expects 1 argument")
            }
            s, ok := args[0].(string)
            if !ok {
                return nil, fmt.Errorf("lower: argument must be a string, got %T", args[0])
            }
            return strings.ToLower(s), nil
        },
    }),
)
```

---

## F.3 Adding Global Constants

```go
env := uexl.DefaultWith(
    uexl.WithGlobals(map[string]any{
        "TAX_RATE":     0.08,
        "MAX_DISCOUNT": 0.40,
        "CURRENCY":     "USD",
    }),
)
```

---

## F.4 Compile Once, Run Many

```go
// At startup
compiled, err := env.Compile("product.basePrice * (1 - discount ?? 0)")
if err != nil {
    log.Fatalf("compile failed: %v", err)
}

// Per request — goroutine-safe
result, err := compiled.Eval(ctx, map[string]any{
    "product":  productMap,
    "discount": 0.15,
})
```

---

## F.5 Validating Expressions Without Compiling

```go
if err := env.Validate(userExpr); err != nil {
    return fmt.Errorf("expression is invalid: %w", err)
}
```

---

## F.6 Structured Parse Error Handling

```go
import parsererrors "github.com/uexl-go/parser/errors"

compiled, err := env.Compile(expr)
if err != nil {
    var pe parsererrors.ParseErrors
    if errors.As(err, &pe) {
        for _, e := range pe.Errors {
            fmt.Printf("line %d col %d [%s]: %s\n",
                e.Line, e.Column, e.Code, e.Message)
        }
    } else {
        // Compile error (unknown function, etc.)
        fmt.Printf("compile error: %v\n", err)
    }
}
```

---

## F.7 Discovering Expression Variables

```go
compiled, _ := env.Compile("customer.tier == 'platinum' && product.stock < threshold")
vars := compiled.Variables()
// vars: ["customer", "product", "threshold"]
// Note: globals and scope vars ($item, etc.) are NOT included
```

---

## F.8 Per-Tenant Environment Extension

```go
baseEnv := uexl.NewEnv(
    uexl.WithFunctions(sharedFunctions),
)

func tenantEnv(cfg TenantConfig) *uexl.Env {
    return baseEnv.Extend(
        uexl.WithGlobals(map[string]any{
            "MAX_DISCOUNT": cfg.MaxDiscount,
            "CURRENCY":     cfg.Currency,
        }),
    )
}
```

---

## F.9 Using the Lib Interface

```go
// Implement vm.Lib
type ShopLib struct{}

func (s ShopLib) Apply(cfg *uexl.EnvConfig) {
    cfg.Functions["upper"] = upperFn
    cfg.Functions["lower"] = lowerFn
    cfg.Functions["trim"]  = trimFn
}

// Register via WithLib
env := uexl.DefaultWith(uexl.WithLib(ShopLib{}))
```

---

## F.10 Hot Reload of Rules

```go
type RuleStore struct {
    mu      sync.RWMutex
    rules   map[string]*uexl.CompiledExpr
}

func (rs *RuleStore) Reload(env *uexl.Env, defs map[string]string) error {
    compiled := make(map[string]*uexl.CompiledExpr, len(defs))
    for name, expr := range defs {
        c, err := env.Compile(expr)
        if err != nil {
            return fmt.Errorf("rule %q: %w", name, err)
        }
        compiled[name] = c
    }
    rs.mu.Lock()
    rs.rules = compiled
    rs.mu.Unlock()
    return nil
}

func (rs *RuleStore) Eval(ctx context.Context, name string, vars map[string]any) (any, error) {
    rs.mu.RLock()
    rule, ok := rs.rules[name]
    rs.mu.RUnlock()
    if !ok {
        return nil, fmt.Errorf("unknown rule %q", name)
    }
    return rule.Eval(ctx, vars)
}
```

---

## F.11 Context Cancellation with Timeout

```go
ctx, cancel := context.WithTimeout(r.Context(), 50*time.Millisecond)
defer cancel()

result, err := compiled.Eval(ctx, vars)
if errors.Is(err, context.DeadlineExceeded) {
    http.Error(w, "rule evaluation timed out", http.StatusGatewayTimeout)
    return
}
```

---

## F.12 Struct to Context Map (JSON Round-Trip)

```go
func structToMap(v any) (map[string]any, error) {
    data, err := json.Marshal(v)
    if err != nil {
        return nil, err
    }
    var m map[string]any
    if err := json.Unmarshal(data, &m); err != nil {
        return nil, err
    }
    return m, nil
}

// Usage
productMap, err := structToMap(product)
if err != nil {
    return nil, err
}
result, err := compiled.Eval(ctx, map[string]any{"product": productMap})
```

---

## F.13 Extracting a Boolean Result

```go
result, err := compiled.Eval(ctx, vars)
if err != nil {
    return false, err
}
b, ok := result.(bool)
if !ok {
    return false, fmt.Errorf("expression must return bool, got %T", result)
}
return b, nil
```

---

## F.14 Extracting a Float64 Result

```go
result, err := compiled.Eval(ctx, vars)
if err != nil {
    return 0, err
}
f, ok := result.(float64)
if !ok {
    return 0, fmt.Errorf("expression must return number, got %T", result)
}
return f, nil
```

---

## F.15 Writing a Custom Pipe Handler (take N)

```go
import "github.com/shoplogic/uexl-go/vm"

func makeTakePipe(n int) vm.PipeHandler {
    return func(pc vm.PipeContext) (any, error) {
        input, ok := pc.Input().([]any)
        if !ok {
            return nil, fmt.Errorf("take: input must be an array")
        }
        if n > len(input) {
            n = len(input)
        }
        return input[:n], nil
    }
}

env := uexl.DefaultWith(
    uexl.WithPipeHandlers(map[string]vm.PipeHandler{
        "take3": makeTakePipe(3),
        "take5": makeTakePipe(5),
    }),
)
```

---

## F.16 The MustCompile Pattern for Package-Level Vars

```go
var (
    discountRule  = mustCompile(shopEnv, "customer.tier == 'platinum' ? 0.25 : 0.0")
    eligibleRule  = mustCompile(shopEnv, "customer.active && customer.totalSpent > 0")
)

func mustCompile(env *uexl.Env, expr string) *uexl.CompiledExpr {
    c, err := env.Compile(expr)
    if err != nil {
        panic(fmt.Sprintf("invalid expression %q: %v", expr, err))
    }
    return c
}
```

> Only use `mustCompile` (or `env.MustCompile`) for expressions that are part of your application code — not for user-supplied expressions.

---

## F.17 Batch Rule Evaluation

```go
func evalRules(
    ctx context.Context,
    rules map[string]*uexl.CompiledExpr,
    vars map[string]any,
) (map[string]any, error) {
    results := make(map[string]any, len(rules))
    for name, rule := range rules {
        result, err := rule.Eval(ctx, vars)
        if err != nil {
            return nil, fmt.Errorf("rule %q: %w", name, err)
        }
        results[name] = result
    }
    return results, nil
}
```

---

## F.18 Safe Type Conversion Helpers

```go
// All numbers from UExL are float64
func toFloat(v any) (float64, bool) {
    f, ok := v.(float64)
    return f, ok
}

func toInt(v any) (int, bool) {
    f, ok := v.(float64)
    if !ok {
        return 0, false
    }
    i := int(f)
    if float64(i) != f {
        return 0, false // not a whole number
    }
    return i, true
}

func toBool(v any) (bool, bool) {
    b, ok := v.(bool)
    return b, ok
}

func toString(v any) (string, bool) {
    s, ok := v.(string)
    return s, ok
}
```
