# Chapter 6: Property and Index Access

> "How you access data is just as important as how you compute with it. UExL's access operators let you navigate any data structure — safely or strictly, depending on what you need."

---

## 6.1 The Access Operators

UExL provides four ways to access properties and array elements:

| Operator | Syntax | Behavior |
|----------|--------|----------|
| Dot access | `obj.key` | Strict — throws if key missing |
| Bracket access | `obj[expr]` | Strict — throws if key/index missing |
| Optional dot | `obj?.key` | Safe — returns null if obj is nullish OR key missing |
| Optional bracket | `obj?.[expr]` | Safe — returns null if obj is nullish OR index/key missing |

The first two are the **strict** operators. The last two are the **optional (safe)** operators. This chapter covers both in detail; the next chapter covers the nullish semantics that underpin them.

---

## 6.2 Dot Access

Dot access reads a named property from an object (or a named key from a map):

```uexl
user.name            // reads "name" from user
order.total          // reads "total" from order
product.category     // reads "category" from product
```

Dot access is strict by default — if the key doesn't exist, a ReferenceError is thrown:

```uexl
// Context: { "user": {"name": "Alice"} }
user.email    // ReferenceError — "email" key does not exist in user
```

This strictness is intentional. In production data pipelines, silently returning null for a missing field can mask data quality problems. UExL wants you to be explicit about when missing data is acceptable.

---

## 6.3 Bracket Access

Bracket access accepts any expression that evaluates to a string (for object keys) or a number (for array indices):

```uexl
// Static string key — same as dot access
user["name"]          // => same as user.name

// Numeric index — for arrays
items[0]              // first element
items[2]              // third element
items[len(items) - 1] // last element

// Dynamic key — computed at runtime
obj[fieldName]        // fieldName is a context variable
items[currentIndex]   // currentIndex is a context variable
```

### Array access

Arrays are zero-indexed. Out-of-bounds returns `null` (not an error):

```uexl
// Context: { "items": [10, 20, 30] }
items[0]      // => 10
items[2]      // => 30
items[5]      // => null  (out of bounds — safe)
items[-1]     // => null  (negative index — safe)
```

> **NOTE:** The different behavior between object key access (missing = error) and array index access (out of bounds = null) is intentional. Missing object properties usually indicate a schema mismatch (programming error), while out-of-bounds array access often indicates "no more items" (expected behavior in loops and filtering).

### String character access

Strings can be indexed like arrays (byte-level):

```uexl
"hello"[0]     // => "h"
"hello"[4]     // => "o"
"hello"[10]    // => null  (out of bounds)
```

---

## 6.4 Chained Access

Property access operators compose naturally — the result of one access becomes the base for the next:

```uexl
order.customer.address.city        // three levels deep
order.items[0].sku                 // array access then property access
order.items[currentIndex].price    // dynamic index then property access
```

Chains evaluate left-to-right. Each step is evaluated before the next:

```uexl
a.b.c.d    // parsed as ((a.b).c).d
```

If any step fails (key missing or null base for strict access), the entire chain throws an error.

---

## 6.5 The Optional Chaining Operators: `?.` and `?.[]`

Optional chaining changes one rule: **if the base is nullish, the access short-circuits to `null`** instead of throwing an error.

```uexl
user?.name          // null if user is null/absent; user.name otherwise
user?.address?.city // null if user or address is null/absent; city otherwise
arr?.[0]            // null if arr is null; arr[0] otherwise
```

### The short-circuit scope

A common misconception is that optional chaining "turns off" errors for the entire chain to the right. It does not. **Each `?.` or `?.[` only guards its own step.**

```uexl
a?.b.c    // if a is null → null (chain stops)
           // if a is non-null → evaluates a.b.c (strict .c)
```

So `a?.b.c` is equivalent to: "if a is null, return null; otherwise evaluate `a.b` strictly and then `.c` strictly." If `a` is non-null but `a.b` is null, then `a.b.c` will throw.

```uexl
// Context: { "user": { "name": "Alice" } }  // no "address" key

user?.address?.city     // => null (address is absent, ?. returns null, chain stops)
user?.address.city      // => ReferenceError! address is absent; .city is strict
```

The second expression has a subtle bug: `user?.address` is safe (returns null if address is absent), but then `.city` is a strict dot access on that null result, which throws.

> **TIP:** Use `?.` consistently throughout a chain when any intermediate value might be missing. Mix strict `.` only when you're certain the property exists.

### Optional chaining on any expression

`?.` can follow any expression, not just identifiers. Use parentheses to make the base explicit:

```uexl
(getValue())?.name        // method call result might be null
(a ?? b)?.property        // nullish coalescing result, then safe access
(items[0])?.price         // array element might be null
```

### `?.` with dynamic keys

```uexl
user?.[fieldName]         // safe dynamic key access
items?.[0]?.price         // safe index then safe property
```

---

## 6.6 When to Use Which Operator

Use this decision tree when writing access expressions:

```
Is the base guaranteed to be non-null at runtime?
│
├── YES — use strict access: .key or [expr]
│         Structural errors fail loudly (good for development)
│
└── NO ──► Is a null base an expected, normal case?
           │
           ├── YES — use optional: ?.key or ?.[expr]
           │         Returns null instead of throwing
           │
           └── NO — guard explicitly before the access:
                     user != null ? user.key : fallback
                     OR:
                     user?.key ?? fallback
```

### Practical guidance for ShopLogic

In the ShopLogic context, most values come from the context map that we control. For these, use strict access — they should always be present:

```uexl
product.basePrice    // we always provide this
customer.tier        // we always provide this
```

For optional enrichment data that might not be present for all customers:

```uexl
customer?.loyaltyPoints ?? 0      // external loyalty system might not have data
product?.promotionalPrice ?? product.basePrice  // might not have a promo price
```

---

## 6.7 Slicing

Beyond single-element access, UExL supports **slice notation** for extracting sub-sequences from arrays and strings:

```uexl
arr[start:end]    // elements from index start (inclusive) to end (exclusive)
arr[start:]       // from start to the end
arr[:end]         // from the beginning to end
arr[:]            // a full copy
```

Examples:

```uexl
[10, 20, 30, 40, 50][1:4]    // => [20, 30, 40]
[10, 20, 30, 40, 50][2:]     // => [30, 40, 50]
[10, 20, 30, 40, 50][:3]     // => [10, 20, 30]
"hello world"[6:]            // => "world"
"hello world"[0:5]           // => "hello"
```

Slice indices out of range are clamped — they don't throw:

```uexl
[1, 2, 3][0:100]    // => [1, 2, 3]  (end clamped to len)
[1, 2, 3][10:20]    // => []          (both beyond end → empty)
```

---

## 6.8 Putting It Together: Navigating Complex Data

Here is a realistic ShopLogic scenario: calculate whether to apply a "bundle discount" when the order contains at least one product from the featured list, and the order has more than 3 items.

```uexl
order.items[0]?.category == 'electronics'
  && len(order.items) > 3
  && featured?.[0] != null
```

More practically — finding the most expensive item in an order:

```uexl
order.items
  |sort: $item.price
  |: $last[len($last) - 1]   // last element after ascending sort = most expensive
```

Accessing the first item's product name safely when items might be empty:

```uexl
order.items?.[0]?.name ?? "No items"
```

---

## 6.9 Summary

- Dot access (`.key`) and bracket access (`[expr]`) are strict — they throw if the key is missing or the base is null.
- Optional dot (`?.key`) and optional bracket (`?.[expr]`) are safe — they return null if the base is nullish or the key/index is absent.
- Each `?.` guards exactly its own step; the rest of the chain retains its own strictness.
- Arrays return `null` on out-of-bounds access (not an error); missing object keys throw.
- Slice notation (`arr[start:end]`) extracts sub-sequences; indices out of range are clamped.
- Use strict access when data is guaranteed present; use `?.` consistently throughout a chain when any step might be missing.

---

## Exercises

**6.1 — Recall.** What is the difference in behavior between `user.email` and `user?.email` when `email` is absent from the `user` object? What about when `user` itself is null?

**6.2 — Apply.** Given the following context:
```go
ctx := map[string]any{
    "order": map[string]any{
        "customer": map[string]any{
            "name": "Bob",
        },
        "items": []any{
            map[string]any{"sku": "A1", "price": 10.0},
            map[string]any{"sku": "B2", "price": 25.0},
        },
    },
}
```
Write expressions for:
1. The customer's name
2. The price of the second item
3. The SKU of the third item (which doesn't exist) — return `"N/A"` instead
4. Safe access to `order.customer.address.city` — return `"unknown"` if absent

**6.3 — Extend.** For ShopLogic, the product object may or may not have a `promotionalPrice` field. Write an expression that returns the effective price: `promotionalPrice` if present and less than `basePrice`, otherwise `basePrice`. Then extend it to apply the customer tier discount on top of the effective price (reuse the tier discount ternary from Chapter 5).
