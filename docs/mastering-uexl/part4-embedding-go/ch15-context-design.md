# Chapter 15: Context Design

> "The context map is the contract between Go and expressions. Design it well and your expressions stay simple. Design it poorly and every expression becomes a workaround."

---

## 15.1 What the Context Map Is

Every UExL evaluation receives a `map[string]any` called the *context* (or *vars*). This map is the sole data channel between Go and expressions:

```go
result, err := compiled.Eval(ctx, map[string]any{
    "product":  productMap,
    "customer": customerMap,
    "today":    "2025-01-15",
})
```

Inside the expression:
```uexl
product.basePrice * (customer.tier == 'platinum' ? 0.80 : 1.0)
```

The context map determines what is visible. Keys not in the map resolve to `null`. Keys in the map that are not referenced by the expression are silently ignored.

---

## 15.2 Mapping Go Types to UExL Types

UExL's type system is `map[string]any`-native. The VM reads values via type assertions. Ensure your context values use the types UExL expects:

| Go type | UExL type | Notes |
|---------|-----------|-------|
| `float64` | Number | ✅ Direct |
| `int`, `int32`, `int64` | Number | ⚠️ Must convert to `float64` |
| `uint`, `uint64` | Number | ⚠️ Must convert to `float64` |
| `string` | String | ✅ Direct |
| `bool` | Boolean | ✅ Direct |
| `nil` | Null | ✅ Direct |
| `[]any` | Array | ✅ Direct |
| `[]string`, `[]int` | ❌ Typed slices | ⚠️ Must convert to `[]any` |
| `map[string]any` | Object | ✅ Direct |
| Go struct | ❌ No support | Must convert to map |

> **WARNING:** Passing typed slices like `[]string` or `[]int` will cause pipe operators to fail — they expect `[]any`. Always convert to `[]any` before passing to context.

### Safe conversion helpers

```go
// Convert []string → []any
func stringSliceToAny(ss []string) []any {
    result := make([]any, len(ss))
    for i, s := range ss {
        result[i] = s
    }
    return result
}

// Convert []int → []any (and to float64)
func intSliceToAny(ns []int) []any {
    result := make([]any, len(ns))
    for i, n := range ns {
        result[i] = float64(n)
    }
    return result
}
```

---

## 15.3 Flattening Structs

Go structs do not map transparently to UExL objects. The cleanest conversion approach depends on your requirements:

### JSON round-trip (simplest, for complex structs)

```go
type Product struct {
    ID        string   `json:"id"`
    BasePrice float64  `json:"basePrice"`
    Category  string   `json:"category"`
    Tags      []string `json:"tags"`
    Stock     int      `json:"stock"`
}

func productToMap(p Product) (map[string]any, error) {
    data, err := json.Marshal(p)
    if err != nil {
        return nil, err
    }
    var m map[string]any
    if err := json.Unmarshal(data, &m); err != nil {
        return nil, err
    }
    return m, nil
}
```

JSON round-trip handles nested structs, slices, and field name mapping (`json:"..."` tags) automatically. Numbers become `float64`, string slices become `[]any` of strings via the JSON decoder. The trade-off is CPU cost — around 1–3 µs for typical structs.

### Direct map conversion (fastest, for hot paths)

```go
func productToMap(p Product) map[string]any {
    return map[string]any{
        "id":        p.ID,
        "basePrice": p.BasePrice,        // float64 → direct
        "category":  p.Category,
        "tags":      stringSliceToAny(p.Tags),  // []string → []any
        "stock":     float64(p.Stock),    // int → float64
    }
}
```

This is explicit — you control exactly which fields appear in the context and how they are named. It is roughly 10× faster than JSON round-trip but requires manual maintenance as the struct evolves.

---

## 15.4 What Should Be in the Context?

**Put in context:**
- All data the expression needs to make a decision
- Derived values that would be expensive to compute in the expression
- Configuration values that vary per call (thresholds, rates, dates)

**Do not put in context:**
- Internal Go types that cannot be converted to UExL values
- Large blobs that are not accessed by the expression (use env globals for constants)
- Sensitive data that expressions should not be able to leak (the log vulnerability: `str(secretKey)`)

**Put in env globals:**
- Values constant across all evaluations (tax rates, config flags)
- Values that change only at startup (holidays, rate tables)

```go
// These never change per-request — globals
shopEnv := uexl.DefaultWith(
    uexl.WithGlobals(map[string]any{
        "TAX_RATE":     0.08,
        "MAX_DISCOUNT": 0.40,
    }),
)

// These change per-request — context vars
result, _ := compiled.Eval(ctx, map[string]any{
    "product":  productMap,
    "customer": customerMap,
    "today":    todayStr,
})
```

---

## 15.5 Context as the Security Boundary

The context map is the only way data enters expressions. This is intentional — it is the **security boundary** between Go and user-submitted expressions.

**Principle: Whitelist, not blacklist.** Build the context by explicitly including what is needed, not by dumping your entire data model and hoping expressions don't misuse it.

```go
// RISKY: pass the entire customer struct
ctx := map[string]any{
    "customer": entireCustomerObject,  // includes internal fields, credit card hash, etc.
}

// SAFE: explicitly project only needed fields
ctx := map[string]any{
    "customer": map[string]any{
        "tier":        customer.Tier,
        "totalSpent":  customer.TotalSpent,
        "memberSince": customer.MemberSince,
        // No: customer.PasswordHash, customer.PaymentMethodID, etc.
    },
}
```

For expressions authored by end-users (not just developers), consider running them in a restricted env that omits functions with side effects or access to sensitive data.

---

## 15.6 Naming Context Variables

Context variable names are case-sensitive and become identifiers in expressions. Follow these conventions:

| Convention | Reason |
|------------|--------|
| `camelCase` for objects | `product.basePrice`, `customer.tier` |
| Short, purpose-named keys | `today`, `threshold` — not `currentDateInISO8601Format` |
| Avoid `$`-prefixed names | System variables (`$item`, `$acc`) use `$` — never use `$` for context vars |
| Plural for arrays | `orders`, `products` — not `order`, `productList` |

---

## 15.7 The `today` Variable Pattern

Date handling in expressions is a common challenge. Since UExL has no date type, dates are ISO 8601 strings. Use lexicographic comparison (it works for `YYYY-MM-DD` format):

```go
// In Go — inject today at evaluation time
vars["today"] = time.Now().Format("2006-01-02")
```

```uexl
// In expressions — lexicographic string comparison works for ISO dates
customer.memberSince < today   // true if member joined before today
order.dueDate < today          // true if order is overdue
```

For date arithmetic (add 30 days, check if within a window), compute the threshold date in Go and inject it:

```go
vars["thirtyDaysAgo"] = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
```

```uexl
customer.memberSince > thirtyDaysAgo   // joined within last 30 days
```

---

## 15.8 Caching and Reuse

The most important performance principle: **compile once, evaluate many times**.

```go
// Startup: compile once per rule
rules := map[string]*uexl.CompiledExpr{}
for name, expr := range ruleConfig.Expressions {
    compiled, err := shopEnv.Compile(expr)
    if err != nil {
        return fmt.Errorf("rule %q invalid: %w", name, err)
    }
    rules[name] = compiled
}

// Per-request: evaluate with fresh context
func ApplyRule(name string, ctx context.Context, vars map[string]any) (any, error) {
    rule, ok := rules[name]
    if !ok {
        return nil, fmt.Errorf("unknown rule %q", name)
    }
    return rule.Eval(ctx, vars)
}
```

`CompiledExpr.Eval` pools the `*vm.VM` internally — no lock contention for concurrent calls. The context map is the only per-call allocation when the expr references far fewer vars than are provided.

---

## 15.9 ShopLogic: The Complete Context Function

```go
package shoplogic

import (
    "time"
)

// Product represents a product in the ShopLogic system.
type Product struct {
    ID        string
    BasePrice float64
    Category  string
    SKU       string
    Name      string
    Tags      []string
    Stock     int
    Rating    float64
}

// Customer represents a customer.
type Customer struct {
    ID             string
    Tier           string   // "platinum", "gold", "silver", "standard"
    TotalSpent     float64
    LoyaltyPoints  int
    MemberSince    string   // ISO date "YYYY-MM-DD"
    Active         bool
}

// EvalContext builds the expression context for a product+customer pair.
// This is the primary data gateway — only fields that expressions need appear here.
func EvalContext(product Product, customer Customer) map[string]any {
    return map[string]any{
        "product": map[string]any{
            "id":        product.ID,
            "basePrice": product.BasePrice,
            "category":  product.Category,
            "sku":       product.SKU,
            "name":      product.Name,
            "tags":      stringSliceToAny(product.Tags),
            "stock":     float64(product.Stock),
            "rating":    product.Rating,
        },
        "customer": map[string]any{
            "id":            customer.ID,
            "tier":          customer.Tier,
            "totalSpent":    customer.TotalSpent,
            "loyaltyPoints": float64(customer.LoyaltyPoints),
            "memberSince":   customer.MemberSince,
            "active":        customer.Active,
        },
        "today": time.Now().Format("2006-01-02"),
    }
}

func stringSliceToAny(ss []string) []any {
    result := make([]any, len(ss))
    for i, s := range ss {
        result[i] = s
    }
    return result
}
```

---

## 15.10 Summary

- The context map is the sole data entry point into expressions — it is the security boundary.
- Use `float64` for all numbers; convert `int`, `int64`, `uint`, etc. at the boundary.
- Use `[]any` for all arrays; never pass `[]string`, `[]int`, or other typed slices.
- Convert structs to `map[string]any` either via JSON round-trip (simple) or direct field mapping (fast).
- Expose only the fields expressions need — whitelist, never dump the full domain model.
- Put constants in env globals; put per-request values in per-call vars.
- Inject `today` as an ISO date string; compute date arithmetic in Go before passing the result.
- Context variable names should be `camelCase`, plural for arrays, never `$`-prefixed.

---

## Exercises

**15.1 — Recall.** Why should you never pass `[]string` directly as a context variable? What error would occur at runtime?

**15.2 — Apply.** A `Product` struct has an `Attributes map[string]string` field containing arbitrary key-value metadata. Write a Go function that converts this struct to the UExL context map, including `Attributes` so it is accessible in expressions as `product.attributes["color"]`.

**15.3 — Extend.** ShopLogic needs to support "flash sale" windows: a sale is active if the current UTC time is between `sale.Start` (RFC3339) and `sale.End` (RFC3339). Design the Go context injection so that a UExL expression can check `today >= sale.start && today <= sale.end`. What format for the date strings ensures correct lexicographic comparison?
