# Chapter 3: Data Types and Literals

> "Before you can manipulate data, you need to understand what kinds of data exist. UExL's type system is small, deliberate, and predictable."

---

## 3.1 The Six Types

UExL has six data types. Every value in every expression is exactly one of them:

| Type | Examples | Go underlying type |
|------|----------|-------------------|
| **Number** | `42`, `3.14`, `1e6`, `NaN`, `Inf` | `float64` |
| **String** | `"hello"`, `'world'` | `string` |
| **Boolean** | `true`, `false` | `bool` |
| **Null** | `null` | `nil` |
| **Array** | `[1, 2, 3]`, `["a", true, null]` | `[]any` |
| **Object** | `{name: "Alice", score: 98}` | `map[string]any` |

There is no `integer` type separate from `float64`. There is no `undefined` (JavaScript's ghost) — missing values are `null`. There are no typed arrays or typed maps. This small type surface keeps expressions readable and reduces the mental overhead when authoring rules.

---

## 3.2 Numbers

Numbers in UExL are always 64-bit IEEE-754 floating-point (`float64`). There is one numeric type, not two.

```uexl
42          // => 42 (stored as float64: 42.0)
-42         // => -42
3.14        // => 3.14
1.5e3       // => 1500 (scientific notation)
0.001       // => 0.001
```

### Integer-like behavior

When a number has no fractional part, UExL displays and compares it as if it were an integer:

```uexl
10 / 2      // => 5 (not 5.0)
10 / 3      // => 3.3333333333333335
7 % 3       // => 1
```

Division always produces the mathematically correct floating-point result — there is no integer division. UExL has no built-in `floor()` function; register one as a host function (Chapter 14) if you need truncated integer division.

### Leading zeros are invalid

```uexl
01    // SyntaxError — leading zeros are not allowed
0.1   // valid — decimal, not octal
```

> **WARNING:** Unlike some languages, `08` and `09` in UExL are syntax errors (not silently treated as decimal). Always write numbers without leading zeros.

### IEEE-754 special values (enabled by default)

UExL supports three IEEE-754 special numeric literals, enabled by default:

| Literal | Meaning | Config flag |
|---------|---------|-------------|
| `NaN` | Not-a-Number | `EnableIeeeSpecials` |
| `Inf` | Positive infinity | `EnableIeeeSpecials` |
| `-Inf` | Negative infinity (unary minus + `Inf`) | `EnableIeeeSpecials` |

These can be disabled via parser options, in which case `NaN` and `Inf` are treated as plain identifiers.

```uexl
NaN == NaN    // => false  (IEEE-754 rule: NaN is never equal to itself)
NaN != NaN    // => true
NaN > 5       // => false  (all comparisons with NaN are false)
Inf > 1e308   // => true
-Inf < -1e308 // => true
1 / Inf       // => 0
```

> **WARNING:** `NaN != NaN` is counterintuitive but mathematically correct per IEEE-754. Use `x != x` to detect NaN — it is the only value not equal to itself. Alternatively, register a host `isNaN(x)` function in your `LibContext` (Chapter 14).

Division by zero remains a runtime error by design — UExL does not silently produce `Inf` from `1 / 0`:

```uexl
1 / 0    // RuntimeError: division by zero
```

---

## 3.3 Strings

Strings are immutable sequences of bytes. UExL strings use single or double quotes — they are interchangeable:

```uexl
"hello"           // => "hello"
'hello'           // => "hello"
"He said, 'hi'"  // Single quotes inside double quotes — valid
'She said, "hi"' // Double quotes inside single quotes — valid
```

The only escape sequences supported are `\\` (literal backslash) and the matching quote (`\"` inside `"..."`, `\'` inside `'...'`). There are no `\n`, `\t`, or `\uXXXX` escapes — embed literal characters directly.

### String indexing (byte-level by default)

```uexl
"hello"[0]    // => "h"     (byte index 0)
"hello"[4]    // => "o"     (byte index 4)
"hello"[10]   // => null    (out of bounds → null, no error)
```

> **NOTE:** Indexing is byte-level by default, matching Go's native string semantics. For multi-byte Unicode strings, use `runes()` or `graphemes()` for character-level access. Chapter 8 covers this in full.

### String slicing

```uexl
"hello world"[0:5]    // => "hello"   (bytes 0–4)
"hello world"[6:]     // => "world"   (bytes 6 to end)
"hello world"[:5]     // => "hello"   (bytes 0–4)
```

### String concatenation

Use the `+` operator to concatenate strings:

```uexl
"Hello, " + "world!"    // => "Hello, world!"
"Order #" + orderId     // => "Order #1042" (if orderId is a string)
```

> **WARNING:** `+` between a string and a number is a TypeError — UExL does not implicitly convert numbers to strings. Use `str(price)` for explicit conversion.

---

## 3.4 Booleans

Booleans have exactly two values: `true` and `false`.

```uexl
true     // => true
false    // => false
```

Booleans arise from comparison and logical operators:

```uexl
10 > 5          // => true
"abc" == "xyz"  // => false
!false          // => true
true && false   // => false
true || false   // => true
```

### Truthiness vs. booleanness

UExL distinguishes between a value *being* a boolean and being *truthy*. The logical operators (`&&`, `||`, `!`) work on *truthiness* — they consider any value that is not `false`, `null`, `0`, `""`, `[]`, or `{}` as truthy.

The truthiness table for common values:

| Value | Truthy? | Notes |
|-------|---------|-------|
| `true` | Yes | |
| `false` | No | |
| `1`, `3.14` | Yes | Any non-zero number |
| `0` | No | Zero is falsy |
| `-1` | Yes | Negative numbers are truthy |
| `"hello"` | Yes | Non-empty string |
| `""` | No | Empty string is falsy |
| `[1, 2]` | Yes | Non-empty array |
| `[]` | No | Empty array is falsy |
| `{a: 1}` | Yes | Non-empty object |
| `{}` | No | Empty object is falsy |
| `null` | No | Null is falsy |

> **TIP:** Use `!!value` as an idiom for explicit boolean conversion: `!!1` → `true`, `!!0` → `false`, `!!""` → `false`, `!!"x"` → `true`. This is more readable than `value != null && value != false && value != 0`.

---

## 3.5 Null

`null` is the explicit representation of "no value." It serves two roles:

1. **Absent context variable** — when a variable is not present in the context map, it resolves to `null`.
2. **Missing collection member** — out-of-bounds array access returns `null`.

```uexl
null                    // the null literal
user.middleName         // => null (if middleName key is absent in user object)
[1, 2, 3][10]           // => null (out of bounds)
```

UExL treats `null` and *absence* as interchangeable for nullish checks — the `??` and `?.` operators handle both the same way. This is intentional: the difference between "the key exists and is null" and "the key doesn't exist" should not matter for most expression logic.

> **NOTE:** Strict property access (`obj.key`) on an object where the key is **absent** throws a ReferenceError. But `??` and `?.` guard against both null and absent, making them the right tools for optional fields. Chapter 7 covers this in detail.

---

## 3.6 Arrays

Arrays are ordered, zero-indexed, heterogeneous collections. They can hold any type, including nested arrays and objects.

```uexl
[1, 2, 3]                    // => [1, 2, 3]
["a", "b", "c"]              // => ["a", "b", "c"]
[1, "two", true, null]       // heterogeneous — valid
[[1, 2], [3, 4]]             // nested arrays — valid
[]                           // empty array — valid
```

### Array access

```uexl
[10, 20, 30][0]    // => 10   (zero-based)
[10, 20, 30][2]    // => 30
[10, 20, 30][5]    // => null (out of bounds — not an error)
[10, 20, 30][-1]   // => null (negative index — not an error)
```

### Array length

```uexl
len([1, 2, 3])    // => 3
len([])           // => 0
```

> **NOTE:** Arrays are created as literals in expressions, but in practice most arrays come from context — e.g., `orders`, `products`, `users` — and are processed through pipe operators. Chapter 10 introduces pipes.

---

## 3.7 Objects

Objects are unordered collections of key-value pairs. Keys must be unique within an object.

```uexl
{"name": "Alice", "age": 30}             // quoted keys
{name: "Alice", age: 30}                 // unquoted identifier keys (equivalent)
{id: 1, tags: ["sale", "featured"]}      // nested array value
{user: {name: "Bob"}, active: true}      // nested object
{}                                        // empty object — valid
```

Both quoted (`"name"`) and unquoted (`name`) key forms are valid when the key is a valid identifier. Use quoted keys for keys containing spaces or special characters.

### Object access

```uexl
{name: "Alice", age: 30}.name    // => "Alice"
{name: "Alice", age: 30}["age"] // => 30
```

In practice, objects almost always come from the expression context rather than being constructed inline:

```uexl
customer.name       // dot notation
customer["name"]    // bracket notation (same result)
```

---

## 3.8 How Types Flow Through the Go Integration

When you pass context data to UExL from Go, the values must be representable as UExL types. Here is the mapping:

| Go type | UExL type | Notes |
|---------|-----------|-------|
| `float64` | Number | Direct |
| `int`, `int64`, `int32`, etc. | Number | Converted to float64 |
| `string` | String | Direct |
| `bool` | Boolean | Direct |
| `nil` | Null | Direct |
| `[]any` | Array | Elements recursively converted |
| `[]T` (typed slices) | Array | Elements converted to `any` |
| `map[string]any` | Object | Direct |
| Struct | Object | Use `encoding/json` marshal/unmarshal or manual conversion |

> **TIP:** The safest way to pass a Go struct to UExL context is via JSON round-trip:
> ```go
> jsonBytes, _ := json.Marshal(myStruct)
> var ctx map[string]any
> json.Unmarshal(jsonBytes, &ctx)
> ```
> This handles nested structs, slices, and field naming automatically. Chapter 15 discusses context design strategies in detail, including performance trade-offs of this approach.

---

## 3.9 No Implicit Type Coercion

A deliberate and important design choice: UExL **never silently converts one type to another**. Operations on mismatched types produce a TypeError.

```uexl
"5" + 3            // TypeError: cannot add string and number
"5" == 5           // false (different types, no coercion)
true + 1           // TypeError: cannot add boolean and number
```

This is unlike JavaScript, where `"5" + 3 === "53"` and `"5" == 5` is `true`. UExL's strictness means bugs in expression logic surface immediately as errors rather than silently producing unexpected results.

When you need type conversion, be explicit:

```uexl
str(42)        // => "42"  (number → string, built-in)
str(3.14)      // => "3.14"
str(true)      // => "true"
!!1            // => true   (any value → boolean via truthiness, built-in idiom)
!!""           // => false
```

For numeric parsing (`"5"` → `5`) or boolean coercion (`1` → `true`), register host functions in `LibContext.Functions` — these are not built-in. Chapter 9 shows the complete built-in function list; Chapter 14 shows how to register your own.

---

## 3.10 ShopLogic: Context Type Planning

For ShopLogic, we need to decide what types our context values will be. Here's an initial type plan:

```go
// ShopLogic expression context structure
ctx := map[string]any{
    // Product — Object
    "product": map[string]any{
        "id":        "PROD-001",           // String
        "basePrice": 99.99,                // Number (float64)
        "category":  "electronics",        // String
        "rating":    4.5,                  // Number
        "stock":     150,                  // Number (int→float64)
        "tags":      []any{"sale", "new"}, // Array of Strings
    },
    // Customer — Object
    "customer": map[string]any{
        "id":        "CUST-042",    // String
        "tier":      "gold",        // String
        "totalSpent": 2850.00,      // Number
        "active":    true,          // Boolean
    },
    // Scalar context values
    "today":      "2026-03-27",  // String (ISO date)
    "maxDiscount": 0.25,         // Number (cap for any discount)
}
```

Notice every value is a primitive or a map/slice — no raw Go structs. This ensures UExL can read them all as properties without any special handling.

---

## 3.11 Summary

- UExL has six types: Number (`float64`), String, Boolean, Null, Array, Object.
- Numbers are always `float64`; leading zeros are invalid; `NaN` and `Inf` are available when enabled.
- Strings use single or double quotes interchangeably; indexing is byte-level by default.
- `null` represents both explicit null values and absent context keys.
- Arrays are zero-indexed; out-of-bounds access returns `null`.
- Objects use dot or bracket notation; both quoted and unquoted keys are accepted.
- UExL performs **no implicit type coercion** — type mismatches are errors.
- Go types map cleanly to UExL types; structs should be converted via JSON or manual field mapping.

---

## Exercises

**3.1 — Recall.** What is the Go underlying type for UExL numbers? What happens when you access an out-of-bounds index on an array? What is the difference between `==` comparing a string and a number in UExL vs. JavaScript?

**3.2 — Apply.** Write a UExL expression that:
1. Takes a context variable `scores` (array of numbers)
2. Uses `len()` to check if the array is non-empty
3. Returns `true` if there are more than 5 scores, `false` otherwise

**3.3 — Extend.** For ShopLogic: write the Go code to construct a context map for a "Platinum" customer browsing a product in the "clothing" category with a base price of `45.00` and a stock of `0`. Then write a UExL expression that evaluates to `"out_of_stock"` when stock is 0, and `"available"` otherwise. Test it using `uexl.Eval()`.
