# Appendix G: Migration Guide

This appendix helps teams that are replacing hand-coded evaluators, JSON Logic, or other expression engines with UExL.

---

## G.1 Migrating from Hand-Coded Go Conditionals

The most common starting point: logic buried in Go functions that stakeholders cannot modify without a deployment.

### Before (Go code)

```go
func computeDiscount(customer Customer, product Product) float64 {
    if customer.Tier == "platinum" {
        if product.BasePrice > 100 {
            return 0.25
        }
        return 0.15
    }
    if customer.TotalSpent > 2000 {
        return 0.10
    }
    return 0.0
}
```

### After (UExL expression stored in database)

```uexl
customer.tier == 'platinum'
  ? (product.basePrice > 100 ? 0.25 : 0.15)
  : (customer.totalSpent > 2000 ? 0.10 : 0.0)
```

**Migration steps:**
1. Add `uexl.Env` and a rule store to your service
2. Translate each branch tree into a ternary expression
3. Load the expression from config or database
4. Compile at startup, evaluate per request
5. Verify with a test suite comparing old and new outputs

---

## G.2 Migrating from `text/template` or `html/template`

Go templates are for rendering text, not computing values. If you are using templates for decision logic (e.g., `{{if gt .Score 90}}A{{end}}`), UExL is a more appropriate tool.

### Before (template-based logic)

```go
const tpl = `{{if eq .Customer.Tier "platinum"}}0.25{{else if gt .Customer.TotalSpent 2000.0}}0.10{{else}}0.0{{end}}`
```

### After (UExL expression)

```uexl
customer.tier == 'platinum' ? 0.25 : (customer.totalSpent > 2000 ? 0.10 : 0.0)
```

**Key differences:**
- UExL returns a typed value, not a rendered string
- Parse errors are structured (line/column), not template error strings
- UExL validates function names at compile time

---

## G.3 Migrating from `govaluate`

`govaluate` is a popular expression evaluator. Key differences:

| Feature | govaluate | UExL |
|---------|-----------|------|
| Type system | Weak (auto-coerce) | No coercion — type mismatches error |
| Null handling | Limited | First-class `null`, `??`, `?.` |
| Pipe operators | Not supported | Built-in |
| Array operations | Limited | Full pipe system |
| Unicode support | Byte-level | Byte, rune, and grapheme levels |
| Compile-time checks | None | Function name validation |
| Goroutine safe | Evaluator per-goroutine | `CompiledExpr` goroutine-safe via pool |

### `govaluate` to UExL mapping

| govaluate | UExL equivalent |
|-----------|-----------------|
| `customer.Tier == "platinum"` | `customer.tier == 'platinum'` (case-sensitive keys) |
| `IN` operator | `arr \|some: $item == val` |
| Parameter functions | `WithFunctions` registered functions |

---

## G.4 Migrating from `expr-lang/expr`

`expr` is a type-safe expression language with struct reflection. UExL takes a different approach:

| Feature | expr | UExL |
|---------|------|------|
| Context type | Go struct (reflected) | `map[string]any` |
| Type safety | Compile-time (struct types) | Runtime (map values) |
| Pipe operators | Not built-in | First-class (`\|map:`, `\|filter:`, etc.) |
| Custom functions | `env.Function(...)` | `WithFunctions(map)` |
| Performance | Very fast (compiled to Go) | Fast (bytecode VM) |

### Key migration consideration

In `expr`, context is a Go struct and the compiler reflects struct fields. In UExL, context is `map[string]any`. You must:
1. Convert all context structs to `map[string]any` (see Chapter 15)
2. Ensure all numbers are `float64` (not `int`)
3. Ensure all slices are `[]any` (not `[]string`, `[]int`)

---

## G.5 Migrating from JSON Logic

JSON Logic represents rules as JSON trees: `{">":[{"var":"price"},100]}`. UExL represents the same as `price > 100`.

### JSON Logic to UExL mapping

| JSON Logic | UExL |
|-----------|------|
| `{"var": "x"}` | `x` |
| `{"==": [a, b]}` | `a == b` |
| `{"!=": [a, b]}` | `a != b` |
| `{"!": [a]}` | `!a` |
| `{"and": [a, b]}` | `a && b` |
| `{"or": [a, b]}` | `a \|\| b` |
| `{"if": [cond, a, b]}` | `cond ? a : b` |
| `{"in": [val, arr]}` | `arr \|some: $item == val` |
| `{"map": [arr, fn]}` | `arr \|map: expr` |
| `{"filter": [arr, fn]}` | `arr \|filter: expr` |
| `{"reduce": [arr, fn, init]}` | `arr \|reduce: ($acc ?? init) + $item` |
| `{"missing": ["x"]}` | `x == null` |
| `{"cat": [a, b, c]}` | `a + b + c` |

---

## G.6 Migrating from Environment Variable-Based Feature Flags

Teams using `os.Getenv` for feature toggles can move to expression-based flags for more sophisticated targeting.

### Before

```go
if os.Getenv("ENABLE_PLATINUM_DISCOUNT") == "true" && customer.Tier == "platinum" {
    price *= 0.75
}
```

### After

```go
// Rule stored in DB: "ENABLE_PLATINUM_DISCOUNT && customer.tier == 'platinum'"
env := uexl.DefaultWith(
    uexl.WithGlobals(map[string]any{
        "ENABLE_PLATINUM_DISCOUNT": featureFlags["enable_platinum_discount"],
    }),
)
```

---

## G.7 Type Coercion Differences

Many expression engines auto-coerce types (e.g., `"5" + 3 = 8`). UExL does not. Be prepared to:

1. Convert strings to numbers in Go before putting them in context
2. Convert integers to `float64` (UExL's only numeric type)
3. Replace `typeof(x) == 'number'` with `x != x` (NaN check) or Go-level type validation
4. Replace `bool(x)`, `number(x)`, `string(x)` with either Go-level conversion or registered host functions

---

## G.8 Migration Testing Strategy

1. **Inventory**: List all expressions currently evaluated (in templates, hand-coded logic, other engines)
2. **Categorize**: Group by expression type (boolean, numeric, string)
3. **Create a reference test suite**: For each expression, record input + expected output pairs
4. **Translate**: Convert each expression to UExL syntax
5. **Compile-check**: Run `env.Validate()` on all translated expressions
6. **Cross-check**: Run both old and new evaluators against the same inputs; compare outputs
7. **Performance test**: Benchmark the hot paths with the compile-once pattern
8. **Rollout**: Shadow-run UExL alongside the old evaluator for one release cycle; compare results in production logs
