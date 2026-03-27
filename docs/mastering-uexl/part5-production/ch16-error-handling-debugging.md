# Chapter 16: Error Handling and Debugging

> "An expression evaluation system is only as reliable as its error reporting. Understanding where an error originates — parse, compile, or runtime — determines how you fix it."

---

## 16.1 The Three Error Origins

UExL processes expressions in three stages, and errors can originate at any of them:

| Stage | When detected | Who is responsible |
|-------|---------------|--------------------|
| **Parse** | During `env.Compile()` or `env.Validate()` | Bad expression syntax |
| **Compile** | During `env.Compile()` | Unknown function names |
| **Runtime** | During `compiled.Eval()` | Type mismatches, bad data, division by zero |

Understanding which stage raised an error tells you whether the problem is in the *expression* (parse/compile) or in the *context data* (runtime).

---

## 16.2 Parse Errors

Parse errors are returned as `*parsererrors.ParseErrors` (from the `parser/errors` package). They contain structured information — line, column, error code, and token — rather than flat strings.

```go
import parsererrors "github.com/uexl-go/parser/errors"

compiled, err := env.Compile("product.price *")  // incomplete expression
if err != nil {
    var pe parsererrors.ParseErrors
    if errors.As(err, &pe) {
        for _, e := range pe.Errors {
            fmt.Printf("Parse error [%s] at line %d, col %d: %s\n",
                e.Code, e.Line, e.Column, e.Message)
        }
    }
}
```

Output:
```
Parse error [missing-operand] at line 1, col 16: expected expression after '*'
```

### `ParseErrors` struct

```go
type ParseErrors struct {
    Errors []ParserError
}

type ParserError struct {
    Code     ErrorCode   // e.g. "unexpected-token", "missing-operand"
    Message  string      // human-readable description
    Line     int         // 1-based line number
    Column   int         // 1-based column number
    Token    string      // the offending token, if applicable
    Expected string      // what was expected, if known
    Context  string      // surrounding context snippet
}
```

### Common error codes

| Code | Meaning |
|------|---------|
| `unexpected-token` | A token appeared where it was not expected |
| `missing-operand` | An operator has no right-hand operand |
| `unexpected-eof` | Expression ended before it was complete |
| `unterminated-string` | String literal was never closed |
| `invalid-pipe-type` | Named pipe type is not a valid identifier |
| `empty-pipe` | Pipe `|:` has no predicate |
| `pipe-in-sub-expression` | Pipe used inside parentheses |
| `alias-in-sub-expression` | `as $alias` used inside parentheses |
| `unclosed-array` | Array `[` was never closed |
| `unclosed-object` | Object `{` was never closed |
| `unclosed-function` | Function call `)` was never closed |

### Checking for specific error codes

```go
if pe.HasErrorCode(parsererrors.ErrUnterminatedString) {
    // Give the user a hint to close their quotes
}
```

---

## 16.3 Compile Errors

Compile errors are plain Go `error` values returned by `env.Compile()`. The most common compile error is an unknown function name:

```go
compiled, err := env.Compile("product.price * round(product.price)")
// err: compile error: unknown function "round" — not registered in this environment
```

This error is a string, not a structured type. Use `strings.Contains` or `errors.Is`/`errors.As` as appropriate for your error handling strategy.

**Why compile-time function checking matters:** Catching unknown functions at compile time (when you load your rule) is far better than discovering them at runtime when a specific expression evaluates for the first time. Always use `env.Compile()` over `uexl.Eval()` in production.

### Validation-only check

If you want to validate an expression without compiling it:

```go
if err := env.Validate(expr); err != nil {
    return fmt.Errorf("expression invalid: %w", err)
}
```

`Validate` performs both parse and compile validation but does not return a `*CompiledExpr`.

---

## 16.4 Runtime Errors

Runtime errors are returned by `compiled.Eval()`. They are plain `error` values with descriptive messages. Common causes:

| Cause | Example expression | Error message |
|-------|--------------------|---------------|
| Wrong function arg count | `len(a, b)` | `len expects 1 argument` |
| Wrong argument type | `len(42)` | `len: unsupported type float64` |
| `substr` out of bounds | `substr(s, -1, 5)` | `substr: invalid start or length` |
| Array index out of bounds | `arr[99]` | runtime error with index |
| Key not found on strict access | `obj.missing` | property access error |
| Pipe input not array | `42 \|map: $item * 2` | pipe handler requires array |
| `reduce` on empty array | `[] \|reduce: ($acc ?? 0) + $item` | runtime error |
| Context cancellation | `compiled.Eval(cancelledCtx, vars)` | `context.DeadlineExceeded` or `context.Canceled` |

### Context cancellation

UExL respects Go context cancellation. Pass a context from your HTTP handler or job runner:

```go
ctx, cancel := context.WithTimeout(r.Context(), 50*time.Millisecond)
defer cancel()

result, err := compiled.Eval(ctx, vars)
if errors.Is(err, context.DeadlineExceeded) {
    return errors.New("expression timed out")
}
if errors.Is(err, context.Canceled) {
    return errors.New("request was canceled")
}
```

---

## 16.5 Structured Error Handling Pattern

For production services, consolidate error handling at the boundary so you log and respond consistently:

```go
type EvalError struct {
    Stage   string // "parse", "compile", "runtime"
    Expr    string
    Message string
    Cause   error
}

func (e *EvalError) Error() string {
    return fmt.Sprintf("[%s error] %s: %s", e.Stage, e.Expr, e.Message)
}

func (e *EvalError) Unwrap() error { return e.Cause }

func safeEval(env *uexl.Env, expr string, vars map[string]any) (any, error) {
    compiled, err := env.Compile(expr)
    if err != nil {
        var pe parsererrors.ParseErrors
        if errors.As(err, &pe) {
            return nil, &EvalError{
                Stage:   "parse",
                Expr:    expr,
                Message: pe.Error(),
                Cause:   err,
            }
        }
        return nil, &EvalError{
            Stage:   "compile",
            Expr:    expr,
            Message: err.Error(),
            Cause:   err,
        }
    }
    result, err := compiled.Eval(context.Background(), vars)
    if err != nil {
        return nil, &EvalError{
            Stage:   "runtime",
            Expr:    expr,
            Message: err.Error(),
            Cause:   err,
        }
    }
    return result, nil
}
```

---

## 16.6 Defensive Expression Patterns

Many runtime errors can be eliminated by using UExL's built-in safety operators in the expression itself.

### Nullish coalescing `??`

Provides a fallback when the left side is `null`. Does NOT fall back on `false`, `0`, or `""`.

```uexl
# Risky: throws if customer.discount is absent
customer.discount * product.basePrice

# Safe: falls back to 1.0 if discount is null
(customer.discount ?? 1.0) * product.basePrice
```

### Optional chaining `?.`

Guards against a null base object. If `customer` is null, short-circuits to null instead of throwing:

```uexl
# Risky: throws if customer is null
customer.discount

# Safe: returns null if customer is null
customer?.discount ?? 1.0
```

### Optional index access `?.[]`

```uexl
orders?.[0]?.total ?? 0.0
```

### Safe reduce with null-coalescing `$acc`

`$acc` is always `null` on the first iteration — always guard it:

```uexl
# CORRECT: guard $acc on first iteration
orders |reduce: ($acc ?? 0) + $item.total

# WRONG: will fail on first iteration
orders |reduce: $acc + $item.total
```

### Guard empty arrays before reduce

```uexl
# Risky: error if orders is empty
orders |reduce: ($acc ?? 0) + $item.total

# Safe: null if orders is empty; ?? catches null
(orders?.[0] != null ? orders |reduce: ($acc ?? 0) + $item.total : 0)
```

---

## 16.7 Validate Expressions at Load Time

In systems where expressions come from a database or configuration file, validate all expressions at startup — before any traffic hits your service:

```go
func loadRules(env *uexl.Env, ruleDB []RuleRow) (map[string]*uexl.CompiledExpr, error) {
    rules := make(map[string]*uexl.CompiledExpr, len(ruleDB))
    var errs []error

    for _, row := range ruleDB {
        compiled, err := env.Compile(row.Expression)
        if err != nil {
            errs = append(errs, fmt.Errorf("rule %q: %w", row.Name, err))
            continue
        }
        rules[row.Name] = compiled
    }

    if len(errs) > 0 {
        return nil, errors.Join(errs...)
    }
    return rules, nil
}
```

This pattern surfaces invalid expressions as startup errors rather than runtime panics.

---

## 16.8 The Validation Endpoint Pattern

For systems where end-users write expressions, expose a validation endpoint so users get immediate feedback:

```go
// POST /api/expressions/validate
// Body: { "expression": "product.price * 0.9" }
func validateHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Expression string `json:"expression"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }

    err := shopEnv.Validate(req.Expression)
    if err == nil {
        json.NewEncoder(w).Encode(map[string]any{"valid": true})
        return
    }

    resp := map[string]any{"valid": false, "errors": []any{}}

    var pe parsererrors.ParseErrors
    if errors.As(err, &pe) {
        errs := make([]any, len(pe.Errors))
        for i, e := range pe.Errors {
            errs[i] = map[string]any{
                "code":    string(e.Code),
                "message": e.Message,
                "line":    e.Line,
                "column":  e.Column,
            }
        }
        resp["errors"] = errs
    } else {
        resp["errors"] = []any{map[string]any{
            "code":    "compile-error",
            "message": err.Error(),
        }}
    }

    w.WriteHeader(http.StatusUnprocessableEntity)
    json.NewEncoder(w).Encode(resp)
}
```

---

## 16.9 Debugging Unknown Failures

When an expression fails at runtime and you are unsure why:

**Step 1: Isolate with `env.Validate`.**
```go
if err := env.Validate(expr); err != nil {
    log.Printf("syntax/compile error: %v", err)
}
```

**Step 2: Inspect the variables the expression references.**
```go
compiled, _ := env.Compile(expr)
log.Printf("expression references: %v", compiled.Variables())
```

This tells you which keys the expression expects. Cross-reference against your context map to find missing or misspelled keys.

**Step 3: Test sub-expressions.** If `product.discount.percentage * 100` fails, test `product.discount` alone to see if the nested access is where it breaks.

**Step 4: Add a `slog`/`zerolog` boundary.** Log the expression and variable list before evaluation during investigation:

```go
slog.Debug("evaluating expression",
    "expr", expressionName,
    "variables", compiled.Variables(),
    "context_keys", contextKeys(vars),
)
result, err := compiled.Eval(ctx, vars)
if err != nil {
    slog.Error("expression evaluation failed",
        "expr", expressionName,
        "error", err,
    )
}
```

---

## 16.10 Summary

- Errors come from three stages: **parse**, **compile**, and **runtime**. The stage tells you where to look.
- Parse errors are `parsererrors.ParseErrors` — structured with code, line, and column.
- Compile errors are plain `error`; most common is "unknown function".
- Runtime errors are plain `error`; most are type mismatches or missing data.
- Always `errors.As(err, &parsererrors.ParseErrors{})` before falling through to plain string handling.
- Use `??`, `?.`, and `?.[]` in expressions to defend against missing or null data.
- Always guard `$acc` in `|reduce:` with `($acc ?? initialValue)`.
- Validate expressions at load time, not at first evaluation.
- Use `compiled.Variables()` to verify which context keys an expression expects.

---

## Exercises

**16.1 — Recall.** At which stage does UExL detect an unknown function like `round()`? Why is this better than detecting it at runtime?

**16.2 — Apply.** An expression user submits `orders[0].total` but the `orders` array may sometimes be empty. Rewrite the expression to return `0` instead of throwing when the array is empty. (Hint: use `?.[]` and `??`.)

**16.3 — Extend.** Write a `loadExpressions` function that reads expressions from a `map[string]string` (name → expression), compiles all of them against a given `*uexl.Env`, and returns a `map[string]*uexl.CompiledExpr` plus a combined error listing all invalid expressions (not just the first). Use structured error types to include the expression name in each sub-error.
