# Chapter 8: Strings and Unicode

> "Strings are deceptively simple — until your application runs in a language other than English. UExL's approach: be explicit about which level of Unicode you're working at, and never silently do the wrong thing."

---

## 8.1 The Three Unicode Levels

A string in UExL is a sequence of bytes. But "a character" can mean three different things depending on what you're trying to do:

| Level | Unit | What it counts | Use for |
|-------|------|---------------|---------|
| **Bytes** | Raw storage byte | Storage size, network protocols | Byte-level operations (default) |
| **Runes** (code points) | Unicode scalar value | Every Unicode character slot | Code point operations |
| **Grapheme clusters** | User-perceived character | What a person sees as one "letter" | Display, truncation, UI |

These three levels diverge for multi-byte Unicode text. Consider `"café"` with decomposed `é` (the letter `e` followed by a combining acute accent `U+0301`):

```
String: c  a  f  e  ́
Bytes:  63 61 66 65 CC 81  →  6 bytes
Runes:  c  a  f  e  ́      →  5 runes (code points)
Graphs: c  a  f  é         →  4 grapheme clusters (user-visible characters)
```

And a family emoji `"👨‍👩‍👧‍👦"`:
```
Bytes:  25 bytes (UTF-8 encoding of 7 code points joined by ZWJ)
Runes:  7 runes
Graphs: 1 grapheme cluster (one visible glyph)
```

UExL makes this explicit rather than silently choosing one level for you.

---

## 8.2 Default: Byte Level (Go-Compatible)

All primitive string operations in UExL are **byte-based**, matching Go's native string semantics. This is intentional:

- Expressions that call Go UDFs (user-defined functions) receive byte positions back — `strings.Index`, regexp matches, etc.
- O(1) `len()` is possible at byte level.
- No Unicode decoding overhead for ASCII strings (the common case).

```uexl
len("hello")        // => 5  (bytes — same as runes and graphemes for ASCII)
"hello"[1]          // => "e"  (byte index 1)
"hello"[1:4]        // => "ell"  (bytes 1–3)
```

For pure-ASCII strings, all three levels agree. The distinction only matters for multi-byte Unicode:

```uexl
len("naïve")        // => 6  (ï encodes as 2 UTF-8 bytes)
"naïve"[0:3]        // => "naï"  (bytes 0–2, valid UTF-8 boundary)
"naïve"[0:2]        // => "na"  (stops before ï)

// Splitting mid-codepoint is allowed syntactically but produces
// invalid UTF-8 — be careful with non-ASCII byte slicing:
"naïve"[2:3]        // ⚠️ one raw byte — part of ï, not a stand-alone character
```

> **WARNING:** Byte-level slicing on multi-byte strings can split a code point and produce invalid UTF-8. When you need safe character-level sub-sequences, use `runeSubstr()` or `graphemeSubstr()`.

---

## 8.3 Rune Level (Code Point Operations)

Use the rune-level functions when you need to work with Unicode scalar values — individual code points.

### Counting

```uexl
runeLen("naïve")         // => 5  (5 code points: n, a, ï, v, e)
runeLen("café\u0301")    // => 5  (c, a, f, e, combining-acute)
runeLen("hello")         // => 5  (ASCII: runes = bytes)
```

### Substring by rune index

```uexl
runeSubstr("naïve", 0, 3)    // => "naï"  (code points 0–2)
runeSubstr("naïve", 2, 3)    // => "ïve"  (starting at code point 2)
```

`runeSubstr(s, start, length)` — `start` is the zero-based rune index, `length` is the number of runes.

### Explode to rune array

```uexl
runes("hello")     // => ["h", "e", "l", "l", "o"]
runes("naïve")     // => ["n", "a", "ï", "v", "e"]
len(runes("naïve")) // => 5  (rune count via array)
```

`runes(s)` returns an array of single-code-point strings. This lets you use pipe operators on them:

```uexl
runes("hello world") |filter: $item != " " |: join($last, "")
// => "helloworld"  (remove spaces at rune level)
```

---

## 8.4 Grapheme Level (User-Visible Characters)

Use grapheme-level functions when you need to match what users see — for UI display, safe truncation, or user-facing string operations.

### Counting visible characters

```uexl
graphemeLen("naïve")        // => 5  (5 visible characters)
graphemeLen("café\u0301")   // => 4  (c, a, f, é — é is decomposed but looks like 1)
graphemeLen("👨‍👩‍👧‍👦")          // => 1  (one visible glyph)
```

### Substring by grapheme position

```uexl
graphemeSubstr("café\u0301", 0, 3)   // => "caf"   (first 3 graphemes)
graphemeSubstr("café\u0301", 2, 2)   // => "fé"    (graphemes 2 and 3)
```

### Explode to grapheme array

```uexl
graphemes("café")     // => ["c", "a", "f", "é"]
graphemes("👨‍👩‍👧‍👦")    // => ["👨‍👩‍👧‍👦"]   (one cluster)
```

---

## 8.5 Choosing the Right Level

| Operation | What to use | Why |
|-----------|-------------|-----|
| Byte count for HTTP `Content-Length` | `len(s)` | Protocol uses bytes |
| Column width for CSV/TSV | `len(utf8Bytes)` or `len(s)` | Byte-based formats |
| Counting characters in a CSV field | `runeLen(s)` | Consistent Unicode scalar count |
| Truncating display text to 20 chars | `graphemeSubstr(s, 0, 20)` | Matches what users see |
| Safe split for processing | `graphemes(s)` | Character boundaries |
| Reversing a string letter by letter | `graphemes(s) \|: join($last[..] reversed)` | Reverse by visible chars |
| Position from a Go regexp match | `s[matchStart:matchEnd]` | Go returns byte positions |

---

## 8.6 String Functions Reference

### Measure

```uexl
len(s)              // byte count  — O(1)
runeLen(s)          // rune (code point) count
graphemeLen(s)      // grapheme cluster count (UAX #29)
```

### Cut

```uexl
substr(s, start, length)           // byte-level  (default)
runeSubstr(s, start, length)       // rune-level
graphemeSubstr(s, start, length)   // grapheme-level
```

### Search and test

```uexl
contains(s, sub)        // => true/false  (byte-level substring check)  [built-in]
```

> **NOTE:** `startsWith()` and `endsWith()` are **not built-in**. Use `substr()` instead:
> ```uexl
> substr(s, 0, 4) == "PRD-"          // starts with "PRD-"
> substr(s, len(s) - 3, 3) == "jpg"  // ends with "jpg"
> ```
> Or register host functions in `LibContext.Functions` (Chapter 14).

### Transform

> **NOTE:** `upper()`, `lower()`, `trim()`, and `replace()` are **not built-in**. Register them as host functions in `LibContext.Functions` if your application needs them (Chapter 14). UExL is designed for extensibility — many common utilities are provided by the embedding application.

### Split, join, and explode

```uexl
join(arr)                // => string  (join array of strings, no separator)  [built-in]
join(arr, separator)     // => string  (join with separator)                   [built-in]
runes(s)                 // => array of single code point strings              [built-in]
graphemes(s)             // => array of grapheme cluster strings               [built-in]
bytes(s)                 // => array of byte values (float64)                  [built-in]
```

> **NOTE:** `split(s, delimiter)` is **not built-in**. Register it as a host function if needed.

### Concatenation

```uexl
"Hello, " + name          // concatenate with + operator  [always available]
"Qty: " + str(qty)        // number to string, then concatenate
```

> **NOTE:** There is no built-in multi-argument `concat()` function. Use `+` chaining.

---

## 8.7 Common String Patterns

### Safe display truncation (30 visible characters)

```uexl
graphemeLen(description) > 30
  ? graphemeSubstr(description, 0, 30) + "…"
  : description
```

### Check if string starts with a prefix (no host function needed)

```uexl
substr(productCode, 0, 4) == "PRD-"    // starts with "PRD-"
```

### Check if string contains any digit

```uexl
runes(code) |some: $item >= '0' && $item <= '9'
```

### Count words (using host-provided `split`, or pipe-based approach)

```uexl
// With host-provided split():
len(split(text, ' '))

// Built-in only — count space characters and add 1:
runes(text) |filter: $item == ' ' |: len($last) + 1
```

### Normalize a product SKU (requires host-provided `upper` and `trim`)

```uexl
// With host-registered upper() and trim():
upper(trim(product.sku))
```

### Build a display label for ShopLogic

```uexl
product.name + " — " + str(product.basePrice) + " USD"
```

---

## 8.8 Type Conversion

UExL has no implicit type coercion — operations on mismatched types produce a TypeError. The built-in conversion tools are intentionally minimal:

### Built-in: any value → string

```uexl
str(42)        // => "42"
str(3.14)      // => "3.14"
str(true)      // => "true"
str(null)      // => "<nil>"
str([1, 2])    // => "[1 2]"   (Go fmt default — for display only)
```

`str()` uses Go's `fmt.Sprintf("%v", value)` internally. It is the only built-in type-to-string converter.

### Built-in idiom: any value → boolean

```uexl
!!value    // convert to boolean via truthiness
!!1        // => true
!!0        // => false
!!""       // => false
!!"hello"  // => true
!!null     // => false
```

The `!!` double-not is the canonical UExL idiom for explicit boolean coercion.

### Host-provided: string → number, number → bool, etc.

`number()`, `bool()`, `parseInt()`, `parseFloat()` and similar are **not built-in**. They are registered by the embedding application in `LibContext.Functions`:

```go
// In your Go embedding (Chapter 14):
functions := vm.VMFunctions{
    "number": func(args ...any) (any, error) {
        // parse string or return nil on failure
    },
    "bool": func(args ...any) (any, error) {
        // coerce to boolean
    },
}
machine := vm.New(vm.LibContext{Functions: functions, PipeHandlers: vm.DefaultPipeHandlers})
```

This design is deliberate: conversion semantics vary by domain (should `"true"` → `true`? should `"1"` → `true`?). The embedding application defines the rules.

---

## 8.9 ShopLogic: String Operations in Practice

**Product display name with category badge (requires host-provided `upper`):**

```uexl
// With host-registered upper():
"[" + upper(substr(product.category, 0, 3)) + "] " + product.name
```

For a product in category `"electronics"` named `"AirPods Pro"`:
```
=> "[ele] AirPods Pro"   (without upper())
=> "[ELE] AirPods Pro"   (with host upper())
```

**Validate that a customer email looks reasonable (built-ins only):**

```uexl
contains(customer.email, "@") && contains(customer.email, ".") && len(customer.email) > 5
```

**Generate a short product code from grapheme-safe name prefix (requires host-provided `upper`):**

```uexl
// With host-registered upper():
upper(graphemeSubstr(product.name, 0, 3)) + "-" + product.id

// Without upper(), using built-ins only:
graphemeSubstr(product.name, 0, 3) + "-" + product.id
```

---

## 8.10 Summary

- UExL strings are byte sequences by default — `len`, `[]`, and `[:]` are byte-level, like Go.
- Use `runeLen`/`runeSubstr`/`runes()` when you need Unicode code-point (scalar) semantics.
- Use `graphemeLen`/`graphemeSubstr`/`graphemes()` when you need user-visible character semantics.
- For ASCII-only strings, all three levels agree — no overhead.
- UExL has no implicit type coercion. Built-in conversions: `str(v)` (any → string), `!!v` (any → boolean). For `number()`, `bool()`, `upper()`, `lower()`, `trim()`, `split()` and similar, register host functions in `LibContext.Functions`.
- Built-in string functions: `len`, `substr`, `contains`, `runeLen`, `runeSubstr`, `graphemeLen`, `graphemeSubstr`, `runes`, `graphemes`, `bytes`, `join`, `str`. Functions like `upper`, `lower`, `trim`, `split`, `replace`, `startsWith`, `endsWith` are host-provided.

---

## Exercises

**8.1 — Recall.** For the string `"café"` (with precomposed `é` = U+00E9), what do `len()`, `runeLen()`, and `graphemeLen()` each return? What does `"café"[0:3]` return?

**8.2 — Apply.** Write UExL expressions using only built-in functions for:
1. Truncate `title` to 50 grapheme clusters, appending `"..."` if truncated.
2. Check if `productCode` (a string) starts with `"PRD-"` (hint: use `substr`).
3. Count the number of spaces in a `description` string (hint: use `runes` and `|filter:`).

**8.3 — Extend.** For ShopLogic, product tags are an array of lowercase strings (e.g., `["sale", "new", "featured"]`). Write a UExL expression that produces a formatted badge string like `"#sale • #new • #featured"` from the tags array using only built-in functions and pipe operators. (If you have a host-provided `upper()` function, describe how the expression changes to produce `"#SALE • #NEW • #FEATURED"`.)
