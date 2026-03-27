# Chapter 9: Functions

> "UExL functions come in two kinds: the small, deliberate set built into the runtime, and the unlimited set your application registers. Know which is which."

---

## 9.1 What a Function Call Looks Like

Functions in UExL use the familiar call syntax:

```uexl
len(name)                  // one argument
substr(text, 0, 10)        // three arguments
join(tags, ", ")           // two arguments
```

Functions are not first-class values — you cannot store them in variables or pass them as arguments. They are called by name, resolved against the registered `VMFunctions` map at runtime.

---

## 9.2 The Built-In Function Library

UExL ships with a fixed set of built-in functions registered in `vm.Builtins`. Every embedding has these available. They divide cleanly into four groups:

### 9.2.1 Length and Measurement

| Function | Signature | Returns | Level |
|----------|-----------|---------|-------|
| `len(v)` | `len(string\|array)` | `number` | Byte / element count |
| `runeLen(s)` | `runeLen(string)` | `number` | Unicode code points |
| `graphemeLen(s)` | `graphemeLen(string)` | `number` | Grapheme clusters (UAX #29) |

```uexl
len("hello")            // => 5   (bytes)
len("naïve")           // => 6   (ï is 2 UTF-8 bytes)
len([1, 2, 3])         // => 3   (array elements)
runeLen("naïve")       // => 5   (code points)
graphemeLen("café")    // => 4   (visible characters)
```

### 9.2.2 Substring Extraction

| Function | Signature | Level |
|----------|-----------|-------|
| `substr(s, start, length)` | byte-level | Byte index |
| `runeSubstr(s, start, length)` | rune-level | Code point index |
| `graphemeSubstr(s, start, length)` | grapheme-level | Visible character |

All three follow the same convention: `start` is zero-based, `length` is the count of units.

```uexl
substr("hello world", 6, 5)          // => "world"   (bytes 6–10)
runeSubstr("naïve", 2, 3)            // => "ïve"     (runes 2–4)
graphemeSubstr("café\u0301", 2, 2)   // => "fé"      (graphemes 2–3)
```

An out-of-range `start` or the combination `start + length > len` does not panic — it returns the available portion or an empty string.

### 9.2.3 Search

| Function | Signature | Returns |
|----------|-----------|---------|
| `contains(s, sub)` | `contains(string, string)` | `bool` |

```uexl
contains("hello world", "world")    // => true
contains("hello world", "xyz")      // => false
contains("", "")                    // => true  (empty string contains empty string)
```

`contains` is byte-level, matching Go's `strings.Contains`.

### 9.2.4 String Explosion and Assembly

These functions break strings apart and put them back together:

| Function | Signature | Returns |
|----------|-----------|---------|
| `runes(s)` | `runes(string)` | `[]any` of single-rune strings |
| `graphemes(s)` | `graphemes(string)` | `[]any` of grapheme-cluster strings |
| `bytes(s)` | `bytes(string)` | `[]any` of byte values as `float64` |
| `join(arr)` | `join(array)` | `string` (no separator) |
| `join(arr, sep)` | `join(array, string)` | `string` (with separator) |

```uexl
runes("hello")              // => ["h", "e", "l", "l", "o"]
graphemes("café")           // => ["c", "a", "f", "é"]
bytes("hi")                 // => [104, 105]  (ASCII values as float64)

join(["a", "b", "c"])       // => "abc"
join(["a", "b", "c"], "-")  // => "a-b-c"
join(["a", "b", "c"], ", ") // => "a, b, c"
```

`join` requires every element to be a string. A non-string element is a runtime error.

### 9.2.5 String Conversion

| Function | Signature | Returns |
|----------|-----------|---------|
| `str(v)` | `str(any)` | `string` |

```uexl
str(42)        // => "42"
str(3.14)      // => "3.14"
str(true)      // => "true"
str(null)      // => "<nil>"
str([1, 2, 3]) // => "[1 2 3]"  (Go fmt default — suitable for display, not parsing)
```

`str` uses `fmt.Sprintf("%v", value)` internally. It is the only built-in any-to-string converter.

> **NOTE:** The output format for arrays and objects is Go's default `%v` format, not JSON. Use it for display labels, not for interoperability. For JSON output, register a host `toJSON()` function.

### 9.2.6 Object Mutation

| Function | Signature | Returns |
|----------|-----------|---------|
| `set(obj, key, value)` | `set(object, string, any)` | `object` (same reference) |

```uexl
set({name: "Alice"}, "age", 30)    // => {name: "Alice", age: 30}
```

**Important:** `set` mutates the original object in place and also returns it. It does not create a copy. In a pure expression context this is usually fine — the object lives only for the expression lifetime.

```uexl
// Typical usage: build up an enriched product record
set(product, "displayPrice", str(product.basePrice) + " USD")
```

---

## 9.3 Complete Built-In Reference

| Function | Category | Args | Return |
|----------|----------|------|--------|
| `len(v)` | Measure | string or array | number |
| `runeLen(s)` | Measure | string | number |
| `graphemeLen(s)` | Measure | string | number |
| `substr(s, start, len)` | Cut | string, number, number | string |
| `runeSubstr(s, start, len)` | Cut | string, number, number | string |
| `graphemeSubstr(s, start, len)` | Cut | string, number, number | string |
| `contains(s, sub)` | Search | string, string | bool |
| `runes(s)` | Explode | string | array |
| `graphemes(s)` | Explode | string | array |
| `bytes(s)` | Explode | string | array |
| `join(arr)` | Assemble | array | string |
| `join(arr, sep)` | Assemble | array, string | string |
| `str(v)` | Convert | any | string |
| `set(obj, key, val)` | Mutate | object, string, any | object |

These fourteen functions are always available regardless of what host functions are registered.

---

## 9.4 Host Functions — Extending the Runtime

Everything beyond the fourteen built-ins must be registered by the embedding application. This is not a limitation — it is the design. UExL's built-in set deliberately stops at the boundary of "universally correct" functions. Anything with domain-specific semantics (number parsing, locale-aware string transforms, math utilities) belongs in the host.

### Registering a host function

```go
import "github.com/maniartech/uexl/vm"

myFunctions := vm.VMFunctions{
    // Math utilities
    "floor": func(args ...any) (any, error) {
        if len(args) != 1 {
            return nil, fmt.Errorf("floor expects 1 argument")
        }
        n, ok := args[0].(float64)
        if !ok {
            return nil, fmt.Errorf("floor expects a number")
        }
        return math.Floor(n), nil
    },

    "min": func(args ...any) (any, error) {
        if len(args) != 2 {
            return nil, fmt.Errorf("min expects 2 arguments")
        }
        a, aOk := args[0].(float64)
        b, bOk := args[1].(float64)
        if !aOk || !bOk {
            return nil, fmt.Errorf("min expects two numbers")
        }
        if a <= b {
            return a, nil
        }
        return b, nil
    },

    // String utilities
    "upper": func(args ...any) (any, error) {
        if len(args) != 1 {
            return nil, fmt.Errorf("upper expects 1 argument")
        }
        s, ok := args[0].(string)
        if !ok {
            return nil, fmt.Errorf("upper expects a string")
        }
        return strings.ToUpper(s), nil
    },

    "lower": func(args ...any) (any, error) {
        if len(args) != 1 {
            return nil, fmt.Errorf("lower expects 1 argument")
        }
        s, ok := args[0].(string)
        if !ok {
            return nil, fmt.Errorf("lower expects a string")
        }
        return strings.ToLower(s), nil
    },
}

// Merge with (or replace) built-ins:
allFunctions := vm.Builtins
for name, fn := range myFunctions {
    allFunctions[name] = fn
}

machine := vm.New(vm.LibContext{
    Functions:    allFunctions,
    PipeHandlers: vm.DefaultPipeHandlers,
})
```

Now expressions can call `floor()`, `min()`, `upper()`, and `lower()` in addition to all built-ins.

### Function registry is per-VM instance

Each `*vm.VM` carries its own function registry. You can create specialized VM instances for different tenants, permission levels, or domains:

```go
// Public expressions — only allow safe string/math functions
publicMachine := vm.New(vm.LibContext{
    Functions:    safeSubset,
    PipeHandlers: vm.DefaultPipeHandlers,
})

// Admin expressions — allow database lookup functions
adminMachine := vm.New(vm.LibContext{
    Functions:    fullFunctionSet,
    PipeHandlers: vm.DefaultPipeHandlers,
})
```

Chapter 14 covers host function design patterns in full, including error conventions, argument validation, and thread safety.

---

## 9.5 Common Host Functions for ShopLogic

Throughout the book's ShopLogic examples, we assume the following host functions are registered. Chapters that use them note them explicitly.

| Function | Behaviour |
|----------|-----------|
| `min(a, b)` | `a <= b ? a : b` |
| `max(a, b)` | `a >= b ? a : b` |
| `floor(n)` | `math.Floor(n)` |
| `ceil(n)` | `math.Ceil(n)` |
| `round(n)` | `math.Round(n)` |
| `abs(n)` | `math.Abs(n)` |
| `upper(s)` | `strings.ToUpper(s)` |
| `lower(s)` | `strings.ToLower(s)` |
| `trim(s)` | `strings.TrimSpace(s)` |
| `replace(s, old, new)` | `strings.ReplaceAll(s, old, new)` |
| `split(s, sep)` | `strings.Split(s, sep)` → `[]any` |
| `number(v)` | Parse string to float64; return nil on failure |

The canonical registration code for ShopLogic is in Chapter 14 (Section 14.2).

---

## 9.6 Functions vs. the `!!` and `str()` Idioms

Three common conversion needs arise constantly. Here is the right tool for each:

| Need | Built-in answer | Why |
|------|----------------|-----|
| Any value → string | `str(v)` | Built-in, always available |
| Any value → boolean | `!!v` | Operator idiom, no function needed |
| String → number | Host `number(v)` | Semantics vary (locale? error value?) |
| Number → bool | `!!v` or `v != 0` | Expression idiom, no function needed |
| Bool → number | `v ? 1 : 0` | Ternary — no function needed |

Prefer built-in idioms over host functions when they exist. They have zero registration overhead and are guaranteed available in every context.

---

## 9.7 Error Handling in Function Calls

A function that returns an error propagates it immediately — UExL does not catch or wrap function errors:

```uexl
join(["a", "b", 42], "-")   // RuntimeError: join: element 2 must be a string, got float64
```

Function errors surface as expression-level errors returned from `vm.Run()`. The VM never panics on a function error — it returns cleanly.

From the host perspective, your registered functions should:
1. Validate argument count and types first
2. Return `(nil, fmt.Errorf(...))` for invalid inputs
3. Never panic — panics escape the VM's error boundary

```go
// Good: explicit validation
"myFunc": func(args ...any) (any, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf("myFunc expects 2 arguments, got %d", len(args))
    }
    n, ok := args[0].(float64)
    if !ok {
        return nil, fmt.Errorf("myFunc: first argument must be a number")
    }
    // ...
},
```

---

## 9.8 ShopLogic: Functions in Practice

Assuming the ShopLogic host functions from Section 9.5 are registered:

**Pricing with capped discount:**

```uexl
product.basePrice * (1 - min(
    customer.tier == 'platinum' ? 0.20 :
    customer.tier == 'gold'     ? 0.15 :
    customer.tier == 'silver'   ? 0.08 : 0.00,
    maxDiscount
))
```

**Formatted product display string:**

```uexl
upper(substr(product.category, 0, 3)) + " | " + product.name + " — " + str(product.basePrice) + " USD"
```

**Count orders above threshold:**

```uexl
orders |filter: $item.total > threshold |: len($last)
```

**Highest-value order total (no host `max` needed — use `|reduce:`):**

```uexl
orders |map: $item.total |reduce: $acc == null || $item > $acc ? $item : $acc
```

---

## 9.9 Summary

- UExL has **14 built-in functions** covering measurement, substring, search, string explosion/assembly, string conversion, and object mutation.
- The built-ins are always available via `vm.Builtins`; register them in `LibContext.Functions`.
- Host functions extend the runtime — register anything domain-specific: math utilities, string transforms, lookup functions, validators.
- `str(v)` and `!!v` cover the two most common conversion idioms without host functions.
- Functions never panic — they return `(result, error)`. Validate arguments and return errors explicitly.
- The compile-once/run-many pattern applies fully to function calls — function resolution happens at runtime, not compile time.

---

## Exercises

**9.1 — Recall.** List all 14 built-in UExL functions by name. Which three categories have multiple members (length at three levels, substring at three levels, and string explosion in three forms)?

**9.2 — Apply.** Without any host functions, write UExL expressions that:
1. Check if the array `tags` contains more than 3 elements and returns the first 3 joined by `", "`.
2. Produce a string summary of a product: `"[PROD-001] Widget (3 in stock)"` using only `str()`, `+`, and context variables `product.id`, `product.name`, `product.stock`.
3. Find the first order whose total exceeds 500 using `|find:`.

**9.3 — Extend.** Register a host function called `clamp(value, min, max)` that keeps a number within a range. Write the Go implementation, register it in `LibContext.Functions`, and write a UExL expression that clamps `product.rating` to the range `[1, 5]`. What happens if `value` is not a number?
