# Chapter 7: Nullish and Optional Semantics

> "The difference between 'no value' and 'false' is not academic — it's the difference between a correct system and one that silently drops valid data."

---

## 7.1 Why This Chapter Matters

If there is one chapter in this book to read carefully and internalize, it is this one. The nullish/optional semantics of UExL are its most deliberate and most differentiating design choices. They determine how expressions handle missing, undefined, or zero-value data — and getting this wrong is the source of the most common production bugs in expression-driven systems.

---

## 7.2 The Three Categories UExL Separates

Most expression languages conflate categories that UExL treats distinctly:

| Category | What it means | UExL mechanism |
|----------|--------------|----------------|
| **Nullish** | Null or absent | `??` coalescing, `?.` chaining |
| **Falsy** | Zero, empty, false | `\|\|` / `&&` logical operators |
| **Missing (structural error)** | Wrong data shape | Strict access throws |

Understanding which category applies to a given situation determines which operator to use.

---

## 7.3 Nullish: `null` and Absent Are the Same

In UExL, a value is nullish if it is:
- The literal `null`
- A context variable not provided by the caller (absent → resolves to `null`)

These two cases are indistinguishable. There is no separate "undefined" — absent context variables are null. This mirrors the `null | undefined` convention in TypeScript/JavaScript optional chains, but UExL codifies it into the type system directly.

```uexl
// Context: { "user": { "name": "Alice" } }
user.email        // ReferenceError (key missing, strict access)
user?.email       // => null (optional access: absent = null)
user?.email == null  // => true (absent is null)
```

---

## 7.4 The Nullish Coalescing Operator: `??`

`??` returns the left-hand value if it is not nullish; otherwise it returns the right-hand value.

```uexl
null ?? "default"      // => "default"
"hello" ?? "default"  // => "hello"
0 ?? 42               // => 0       (0 is not nullish!)
false ?? true         // => false   (false is not nullish!)
"" ?? "empty"         // => ""      (empty string is not nullish!)
```

The key insight: **`??` only falls back on `null` and absent. It preserves all falsy-but-valid values.**

### When `||` is wrong and `??` is right

This distinction matters enormously for real data:

```uexl
// Inventory count — 0 is a meaningful value (out of stock)
stock ?? 0       // Good: falls back only when stock is absent
stock || 0       // BAD: 0 stock would be replaced by 0 (same result here, but...)
stock || 10      // VERY BAD: "out of stock" becomes "10 in stock"!

// Score — 0 is a valid score
score ?? 0       // Good: provides default only when score is absent
score || 0       // BAD: a score of 0 would silently become 0 (looks OK but wrong intent)

// Premium status — false is meaningful
isPremium ?? false   // Good: defaults to false, preserves true
isPremium || false   // BAD: always returns false when isPremium is false
```

> **WARNING:** Using `||` as a null-coalescing substitute is one of the most common expression bugs. Use `??` whenever your semantic is "provide a default for missing data."

### Chained nullish coalescing

```uexl
// Try each source in order until you find a non-null value
user.nickname ?? user.name ?? "Anonymous"
config.timeout ?? env.DEFAULT_TIMEOUT ?? 5000
```

### Short-circuit evaluation

The right-hand side of `??` is only evaluated if the left side is nullish:

```uexl
cheapProduct ?? expensiveComputation()  // expensiveComputation only runs if left is null
```

---

## 7.5 Safe Mode: How `??` Interacts with Strict Access

`??` not only provides defaults — it also provides "safe mode" for the immediate preceding access.

When you write `obj.key ?? default`, UExL softens the `obj.key` access:
- If `obj.key` results in a missing key (what would normally be a ReferenceError), returns `default` instead.
- If `obj.key` is present and `null`, returns `default`.
- If `obj.key` is present and non-null, returns its value.

```uexl
// Context: { "user": { "name": "Alice" } }  // no "email" key

user.email ?? "no-email"     // => "no-email"  (email absent → fallback, no ReferenceError)
user.name ?? "no-name"       // => "Alice"     (present and non-null)
```

### The scope of safe mode: one step only

Critical: `??` only softens the **immediately preceding** access. It does not protect earlier steps in the chain.

```uexl
// Context: { "order": null }

order.customer.name ?? "Unknown"
// Still a ReferenceError: order is null, order.customer throws BEFORE ?? can rescue it
```

To protect the full chain, use `?.` for intermediate steps:

```uexl
order?.customer?.name ?? "Unknown"
// Now: if order is null → null, if customer is null → null, name falls back to "Unknown"
```

### The rule to remember

> `a.b.c ?? default` — only `c` is softened (the final `??`-adjacent access).
> All earlier steps (`a`, `.b`) remain strict.

---

## 7.6 The Optional Chaining Operators: `?.` and `?.[`

While `??` provides defaults for null values, `?.` prevents errors when *navigating through* null values.

```uexl
user?.name              // null if user is null; user.name otherwise
user?.address?.city     // null at any null in the chain
arr?.[0]                // null if arr is null; arr[0] otherwise
```

### The interaction between `?.` and `??`

These operators are designed to compose:

```uexl
user?.address?.city ?? "Unknown"
```

Reading this left to right: navigate safely through `user` and `address` (either may be null), then provide an "Unknown" default if the final result is null.

This is the idiomatic pattern for accessing optional deeply nested data with a fallback.

### What `?.` does NOT do

`?.` guards the *base* of the access from being nullish. If the base is non-nullish:
- `?.key` behaves identically to `.key` — the key must exist, or it throws
- `?.[i]` behaves identically to `[i]` — out-of-bounds returns null

Wait — if `?.[i]` with an out-of-bounds index returns null, isn't that safe already? Yes, but `?.` also protects against the base being null:

```uexl
// Difference:
arr[0]     // TypeError if arr is null
arr?.[0]   // null if arr is null, otherwise arr[0]
```

For **object keys**, `?.key` also softens missing keys (like `?? `'s safe mode):

```uexl
user?.email    // null if user is null OR if email key is missing in user
```

This is the key difference from bare dot access:

```uexl
user.email     // throws ReferenceError if email is missing
user?.email    // returns null if email is missing (or if user is null)
```

---

## 7.7 Combining `??` and `?.`: The Complete Patterns

### Pattern 1: Safe leaf access with fallback

Access a property that may not exist on an object that is guaranteed non-null:

```uexl
user.nickname ?? user.name    // nickname might be absent; name is guaranteed
```

### Pattern 2: Safe navigation with fallback

Traverse a chain where any level might be null:

```uexl
user?.company?.address?.city ?? "Unknown"
```

### Pattern 3: Safe navigation with complex fallback

```uexl
order?.discount?.rate > 0 ? order.discount.rate : 0
```

Or more cleanly:

```uexl
(order?.discount?.rate ?? 0) > 0 ? order.discount.rate : 0
```

### Pattern 4: Layered defaults for configuration

```uexl
config.timeout ?? environment.DEFAULT_TIMEOUT ?? 3000
```

### Pattern 5: Pipe with null-safe result

```uexl
users |find: $item.id == targetId
  |: $last?.name ?? "Not found"
```

---

## 7.8 Logical Operators: Truthiness Guards

`||` and `&&` work on truthiness, not nullishness. Use them for control flow, not data defaults.

### `||` for fallback to any truthy value

```uexl
message || "Default message"    // falls back if message is falsy (null, "", 0, false)
```

OK when you genuinely want to treat all falsy values as missing. Wrong when falsy values are meaningful.

### `&&` for conditional evaluation

```uexl
user.active && user.score > 80    // only check score if user is active
items.length > 0 && items[0].price    // only access price if items is non-empty
```

### The `!!` conversion pattern revisited

```uexl
!!customer.email           // true if email is present and non-empty (truthy)
!!customer?.email          // true if email exists and is truthy, false if customer is null
```

---

## 7.9 ShopLogic: Null-Safe Pricing

Let's apply everything from this chapter to the ShopLogic pricing expression. The rules:

1. Use the product's `promotionalPrice` if it exists and is lower than `basePrice`
2. Apply the customer's tier discount (Platinum 20%, Gold 15%, Silver 8%, Standard 0%)
3. Apply an extra 2% for customers with `loyaltyPoints > 500` (loyalty may not exist for all customers)
4. Cap total discount at `maxDiscount`
5. Express the final price

```uexl
// Step 1: Effective base price (promotional or standard)
(product?.promotionalPrice ?? product.basePrice) < product.basePrice
  ? product.promotionalPrice
  : product.basePrice
```

Let's use `effectiveBase` as a mental alias (since UExL doesn't have variables):

> **NOTE:** `min(a, b)` below is a **host-provided function** — register it in `LibContext.Functions` (Chapter 14). Without it, use `a <= b ? a : b`.

```uexl
// Full pricing expression
(product?.promotionalPrice != null && product.promotionalPrice < product.basePrice
    ? product.promotionalPrice
    : product.basePrice)
  * (1 - min(
      (customer.tier == 'platinum' ? 0.20 :
       customer.tier == 'gold'     ? 0.15 :
       customer.tier == 'silver'   ? 0.08 : 0.00)
      + ((customer?.loyaltyPoints ?? 0) > 500 ? 0.02 : 0),
      maxDiscount
    ))
```

Breaking it down:
- `product?.promotionalPrice != null` — safe check for promotion existence
- `customer?.loyaltyPoints ?? 0` — zero if loyalty program not available for this customer
- `min(..., maxDiscount)` — caps total discount (host-provided function)

---

## 7.10 Quick Reference: Choose the Right Operator

| Situation | Use | Why |
|-----------|-----|-----|
| Default when value is null or absent | `value ?? default` | Preserves valid falsy |
| Navigate through possibly-null objects | `obj?.key` | Short-circuits on null |
| Navigate AND provide default | `obj?.key ?? default` | Safe nav + fallback |
| Default when any falsy value | `value \|\| default` | Falsy-based |
| Guard before using value | `value && doSomething(value)` | Truthiness check |
| Force boolean result | `!!value` | Explicit coercion |
| Check if absent | `value == null` | Treats null and absent together |

---

## 7.11 Summary

- "Nullish" in UExL means `null` or absent (unset context variable). These are equivalent.
- `??` falls back only on nullish values — it preserves `0`, `""`, `false`, `[]`, `{}`.
- Use `??` for data defaults; use `||` only for truthiness-based fallbacks where falsy = missing.
- `??` applies safe mode to the immediately preceding access only — earlier chain steps remain strict.
- `?.` guards against a null base at each step where it appears; use it consistently throughout optional chains.
- `?.` on property access also softens missing keys.
- Compose `?.` and `??` for the complete safe-navigation-with-fallback pattern: `a?.b?.c ?? default`.

---

## Exercises

**7.1 — Recall.** What does `0 ?? 42` evaluate to? What does `0 || 42` evaluate to? Why are they different? What is the "scope" of safe mode applied by `??`?

**7.2 — Apply.** Given the context below, predict the result of each expression:
```go
ctx := map[string]any{
    "user": map[string]any{
        "name":  "Alice",
        "score": 0,      // valid zero score
        // "email" is absent
    },
    // "config" is absent
}
```
Expressions:
1. `user.score ?? 100`
2. `user.score || 100`
3. `user.email ?? "no email"`
4. `user.email || "no email"`
5. `config?.maxItems ?? 50`
6. `config.maxItems ?? 50`

**7.3 — Extend.** For ShopLogic, the product object may have an optional `bundleDiscount` field — an additional percentage off to apply when the customer buys at least 3 items. The `order.itemCount` field tells you how many items are in the order. Write an expression that:
1. Applies the tier discount (reuse from previous chapters)
2. Adds `bundleDiscount` (if present) when `itemCount >= 3`
3. Caps at `maxDiscount`
4. Returns the final price

Use `??` appropriately to handle the case where `bundleDiscount` is absent.
