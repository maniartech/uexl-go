# Chapter 4: Identifiers, Variables, and Context

> "An expression without data is just a constant. Context is what makes expressions powerful — it's the bridge between your data and your logic."

---

## 4.1 What Is an Identifier?

An identifier is a name that refers to a value. In a UExL expression, identifiers are resolved against the *context* — the data you provide at runtime.

```uexl
price               // identifier: looks up 'price' in the context
customer.tier       // chained identifiers: customer object, then its 'tier' property
```

### Identifier naming rules

Valid identifier characters:
- Start with a letter (`a`–`z`, `A`–`Z`) or underscore (`_`)
- Continue with letters, digits (`0`–`9`), or underscores
- Unicode letters are allowed (e.g., `café`, `名前`)

```uexl
price        // valid
_total       // valid (leading underscore)
item2        // valid
2item        // INVALID — cannot start with a digit
my-price     // INVALID — hyphens are not allowed
```

The special identifier prefix `$` is **reserved for system variables** (covered in Section 4.4):

```uexl
$item        // system variable inside a pipe
$index       // system variable inside a pipe
$last        // system variable — previous pipe stage result
```

Do not use `$` as a prefix for your own context variable names.

### Case sensitivity

Identifiers are case-sensitive:

```uexl
price    // resolves to context["price"]
Price    // resolves to context["Price"] — a different key!
PRICE    // resolves to context["PRICE"] — yet another key!
```

Use consistent casing (we recommend `camelCase` for context keys) and document your context schema so expression authors know exactly what names are available.

---

## 4.2 The Context Map

When you call `machine.Run(bytecode, ctx)` or `uexl.Eval(expr, ctx)`, the second argument is the **context** — a `map[string]any` that provides all the data the expression can read.

```go
ctx := map[string]any{
    "price":    99.99,
    "customer": map[string]any{
        "name": "Alice",
        "tier": "gold",
    },
    "tags": []any{"sale", "featured"},
}
```

Inside an expression:

```uexl
price                  // => 99.99
customer.name          // => "Alice"
customer.tier          // => "gold"
tags[0]                // => "sale"
len(tags)              // => 2
```

The context map is the **only** way data enters a UExL expression. Expressions cannot read environment variables, files, network resources, or any global state. This is deliberate — it defines the security boundary.

---

## 4.3 Absent Variables

When an identifier is not present in the context map, it resolves to `null` — it does not throw an error:

```uexl
// Context: {"price": 99.99}
discount    // => null (not in context)
```

This behavior makes expressions more resilient when contexts are partially populated. However, it can also hide typos:

```uexl
prce * 2    // typo: 'prce' resolves to null, null * 2 is a TypeError
```

> **TIP:** Enable strict mode (in host configuration) to make absent variable references throw a ReferenceError instead of resolving to null. This is strongly recommended for development and testing environments.

---

## 4.4 System Variables

UExL defines a set of **system variables** that are automatically created by the runtime in specific contexts. These always start with `$` and are never provided by the host application.

### Inside pipe stages

System variables are the mechanism through which pipes communicate between stages. Each pipe type injects its own set of variables:

| Variable | Available in | Meaning |
|----------|-------------|---------|
| `$last` | All pipes | Result from the previous stage (or left-hand expression) |
| `$item` | `map`, `filter`, `find`, `some`, `every`, `sort`, `unique`, `groupBy`, `flatMap` | Current element being processed |
| `$index` | Same as `$item` | Zero-based position of `$item` in the collection |
| `$acc` | `reduce` | Accumulator value built up across iterations |
| `$window` | `window` | Current sliding window sub-array |
| `$chunk` | `chunk` | Current fixed-size chunk sub-array |

These are covered in full in Chapters 10–12. For now, the key point: **`$`-prefixed names are reserved for the runtime** and must not be used as host context variable names.

### How `$last` threads through a pipeline

`$last` is the simplest system variable to understand — it always holds the output of the previous stage:

```uexl
10 + 20         // produces 30
  |: $last * 2  // $last = 30, result = 60
```

---

## 4.5 Nested Context and Property Access

Context values can be nested to any depth. Access nested properties using dot notation or bracket notation:

```go
ctx := map[string]any{
    "order": map[string]any{
        "id": "ORD-001",
        "customer": map[string]any{
            "name":    "Alice",
            "address": map[string]any{
                "city":    "London",
                "country": "UK",
            },
        },
        "items": []any{
            map[string]any{"sku": "A1", "qty": 2, "price": 15.00},
            map[string]any{"sku": "B2", "qty": 1, "price": 45.00},
        },
    },
}
```

```uexl
order.id                          // => "ORD-001"
order.customer.name               // => "Alice"
order.customer.address.city       // => "London"
order.items[0].sku                // => "A1"
order.items[1].price              // => 45.00
len(order.items)                  // => 2
```

> **NOTE:** Strict property access (`.key`) throws a ReferenceError if the key does not exist at any point in the chain. If `order.customer.address` might be missing from some orders, you need optional chaining: `order.customer?.address?.city`. Chapter 6 covers access operators in full, and Chapter 7 covers nullish/optional semantics.

---

## 4.6 Dynamic Key Access

When the property name is determined at runtime (rather than being a literal), use bracket notation with an expression:

```uexl
obj["name"]              // static key — equivalent to obj.name
obj[keyVar]              // dynamic key — key comes from a variable
obj[prefix + "_total"]  // dynamic key — computed from an expression
```

This is particularly useful in pipe expressions where the key varies per item:

```uexl
// Get the value of a dynamically specified field from each product
products |map: $item[fieldName]
```

---

## 4.7 Context as a Public API

Think of the context map as a **public API** for your expressions. The expression authors (potentially business users or other developers) will write against this API. The same principles that apply to REST API design apply here:

**Stability.** Once you publish a context variable name, expressions in production depend on it. Renaming `customer.tier` to `customer.level` will break every expression that references the old name.

**Completeness.** Provide everything the expression author needs. If they need the current date, add it. If they need a formatted version of a value, consider adding the pre-formatted version rather than requiring expressions to format it.

**Clarity.** Use self-documenting names. `customer.totalSpent` is better than `customer.ts`. `order.lineItems` is better than `order.li`.

**Separation.** Don't expose raw database models. Create purpose-built context objects that expose only what expressions should see. A `customer` context object might contain `id`, `name`, `tier`, `totalSpent`, and `active` — but not internal fields like `passwordHash`, `createdAt`, or `stripeCustomerId`.

We'll return to context design principles with concrete patterns and anti-patterns in Chapter 15.

---

## 4.8 ShopLogic: Finalizing the Context Schema

Based on the ShopLogic requirements we've outlined so far, here is the finalized context schema for version 1. We'll use this throughout the rest of the book:

```go
// ShopLogicContext builds the expression context for a pricing/filtering evaluation.
func ShopLogicContext(product Product, customer Customer, env Environment) map[string]any {
    return map[string]any{
        "product": map[string]any{
            "id":        product.ID,
            "name":      product.Name,
            "basePrice": product.BasePrice,
            "category":  product.Category,
            "rating":    product.Rating,
            "stock":     float64(product.Stock),
            "tags":      toAnySlice(product.Tags),
        },
        "customer": map[string]any{
            "id":          customer.ID,
            "name":        customer.Name,
            "tier":        customer.Tier,       // "platinum", "gold", "silver", "standard"
            "totalSpent":  customer.TotalSpent,
            "active":      customer.Active,
            "memberSince": customer.MemberSince.Format("2006-01-02"),
        },
        "today":       env.Today.Format("2006-01-02"),
        "maxDiscount": env.MaxDiscount, // e.g., 0.30 (30% cap)
    }
}

func toAnySlice(ss []string) []any {
    result := make([]any, len(ss))
    for i, s := range ss {
        result[i] = s
    }
    return result
}
```

This function is the **gateway** between Go types and UExL expressions. It converts Go structs into the flat `map[string]any` structure that UExL can read, using only the fields that expression authors need.

A sample expression using this context:

> **NOTE:** `min(a, b)` below is a **host-provided function** — register it in `LibContext.Functions` (Chapter 14). Without it, use `a <= b ? a : b`.

```uexl
product.basePrice
    * (1 - min(
        customer.tier == 'platinum' ? 0.20 :
        customer.tier == 'gold'     ? 0.15 :
        customer.tier == 'silver'   ? 0.08 : 0,
        maxDiscount
    ))
```

This expression is readable, business-friendly, and fully testable with different context values.

---

## 4.9 Summary

- Identifiers resolve names against the context map; they are case-sensitive.
- Absent identifiers resolve to `null` (not an error); enable strict mode in development to catch typos as errors.
- The context map is the only data entry point into expressions — it defines the security boundary.
- System variables (`$item`, `$index`, `$acc`, `$last`, etc.) are created by the runtime inside pipe stages. Never use `$`-prefixed names for your own context variables.
- Nested context objects are accessed with dot notation (`customer.tier`) or bracket notation (`customer["tier"]`).
- Treat context as a public API: stable, self-documenting, and purpose-built for expression consumption.
- The ShopLogic context schema (`product`, `customer`, `today`, `maxDiscount`) is our working example for the rest of the book.

---

## Exercises

**4.1 — Recall.** What does UExL return when an expression references an identifier that is not in the context map? In which environment would you enable strict mode, and why?

**4.2 — Apply.** Given the ShopLogic context schema from Section 4.8, write a UExL expression that:
1. Returns `true` if the customer has been a member for at least one year (use string comparison of ISO dates: `customer.memberSince <= oneYearAgo`)
2. Returns `"loyal_gold"` if they are a gold customer AND have been a member for at least a year, otherwise `"standard"`

Note: Add a `oneYearAgo` variable to the context in Go (a string in `"2006-01-02"` format).

**4.3 — Extend.** ShopLogic needs to support a new concept: "featured products" — a list of product IDs that are promoted this week. Add a `featured` field to the context (as an array of product ID strings) and write a UExL expression that returns `true` if `product.id` appears in the `featured` list. (Hint: use `|some:` — preview of Chapter 10.)
