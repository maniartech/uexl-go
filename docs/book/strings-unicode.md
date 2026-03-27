# Strings and Unicode

Working with text requires choosing the right unit of processing: bytes for protocols and storage limits, code points (runes) for Unicode scalar operations, and grapheme clusters for user-visible characters.

UExL provides explicit, composable built-in functions that operate at clearly-named Unicode levels, with zero hidden state and full compatibility with user-defined functions (UDFs).

## Background: The Three Unicode Levels

| Level | Unit | Go term | Use for |
|-------|------|---------|----------|
| **Bytes** | Raw storage unit | `byte` | Protocols, encoding, storage limits |
| **Code points** | Unicode scalar value | `rune` | Identifiers, ASCII processing, character classification |
| **Grapheme clusters** | User-perceived character (1+ code points) | — (UAX #29) | Display, UI text, safe truncation |

**Why they differ — `"café\u0301"` (decomposed é = e + combining acute accent):**

```javascript
len("café\u0301")           // 8  (UTF-8 bytes — default, Go-compatible)
graphemeLen("café\u0301")   // 4  (4 visible characters: c, a, f, é)
len(runes("café\u0301"))    // 5  (rune count via explicit conversion)
```

**Examples of where levels diverge:**

- `"café"` with precomposed é (U+00E9): byte=5, rune=4, grapheme=4
- `"café\u0301"` with decomposed é: byte=8, rune=5, grapheme=4
- `"👨‍👩‍👧‍👦"` (family emoji): byte=25, rune=7, grapheme=1

## Default Behavior: Byte Level (Go-Compatible)

All four primitive string operations are **byte-based** — matching Go's native `len()`, `str[i]`, and `str[i:j]`. This keeps UExL predictable for Go developers and consistent with UDFs that return byte positions (e.g. `strings.Index`, regexp matches).

```javascript
// ASCII strings: bytes = runes = graphemes — all levels agree:
len("hello")           // 5
"hello"[1]             // "e"
"hello"[1:4]           // "ell"
substr("hello", 1, 3)  // "ell"

// Multi-byte Unicode: byte boundaries and rune boundaries diverge:
len("naïve")           // 6  (ï encodes as 2 UTF-8 bytes)
"naïve"[2]             // raw 3rd byte (0xAF) — part of ï, not a valid rune alone
"naïve"[0:4]           // "naï"   (bytes 0–3 happen to land on a valid boundary)
"naïve"[0:3]           // invalid UTF-8 ⚠️  (splits the 2-byte ï codepoint)
```

Use `runes(s)` or `graphemes(s)` to work at the character or display level safely.

## Function Reference

### Measure

```javascript
len(s)                  // UTF-8 byte count — O(1) (default)
runeLen(s)              // rune (code point) count (new)
graphemeLen(s)          // grapheme cluster count — UAX #29 segmentation (new)
len(runes(s))           // rune count via explicit array conversion (same as runeLen)
```

### Cut

```javascript
substr(s, start, length)             // byte-level substring (default)
s[start:end]                         // byte-level slice (default)
runeSubstr(s, start, length)         // rune-level — safe for Unicode strings (new)
graphemeSubstr(s, start, length)     // grapheme-level — safe for display (new)
```

### Explode to Array

```javascript
runes(s)        // []any of single code point strings
graphemes(s)    // []any of grapheme cluster strings
bytes(s)        // []any of byte values as float64
```

### Reassemble

```javascript
join(array)          // concatenate elements with "" separator
join(array, sep)     // concatenate elements with separator
```

Note: `|join:` pipe stage is still the cleanest choice for a final pipeline step. The `join()` function is needed when joining is required mid-expression or as a function argument.

## Examples

### Measure at Different Levels

```javascript
// Default — byte count (Go-compatible):
len("naïve")                  // 6  (ï = 2 UTF-8 bytes)
len("hello")                  // 5  (ASCII: bytes = runes = graphemes)
len("👨‍👩‍👧‍👦")                // 25 (25 UTF-8 bytes)

// For display/UI — visible character count:
graphemeLen("👨‍👩‍👧‍👦")          // 1  (one visible emoji)
graphemeLen("café\u0301")     // 4  (c, a, f, é as one grapheme)

// Rune count — explicit via array conversion:
len(runes("naïve"))           // 5
```

### Cut and Slice

```javascript
// Byte-based (default) — safe for ASCII, predictable for all:
"hello"[1:4]                           // "ell"  (bytes 1–3)
"hello"[2]                             // "l"    (byte 2)
substr("naïve", 0, 3)                   // first 3 bytes ⚠️ (may split multi-byte ï)

// Rune-level — explicit via explode + join or runeSubstr (new):
runes("naïve")[2]                       // "ï"  (3rd rune)
join(runes("naïve")[0:3])               // "naï"  (first 3 runes, rejoined)
runeSubstr("naïve", 0, 3)              // "naï"  (rune-level substr — new builtin)

// Safe display truncation — grapheme-level (new):
graphemeSubstr("café\u0301", 0, 3)      // "caf"   (3 graphemes)
graphemeSubstr("café\u0301", 0, 4)      // "café\u0301"  (é stays whole)
graphemeSubstr("👨‍👩‍👧‍👦 hello", 0, 2)  // "👨‍👩‍👧‍👦 "   (2 graphemes, emoji intact)
```

### Explode → Transform → Reassemble

```javascript
// Strip combining diacritics (normalize to base letters):
runes("café\u0301") |filter: $item != "\u0301" |join: ""
// → "cafe"

// Map upper over visible characters, then rejoin:
join(graphemes("café\u0301") |map: upper($item), "")
// → "CAFÉ"

// Mid-expression join (impossible with |join: pipe alone):
upper(join(runes("café\u0301") |filter: $item != "\u0301", ""))
// → "CAFE"

// Check all bytes are ASCII:
bytes(header) |every: $item < 128        // true/false

// byte count via array (same as byteLen):
len(bytes(str))
```

### Practical Use Cases

```javascript
// 1. Safe display name — max 20 visible characters:
graphemeSubstr(userName, 0, 20)

// 2. Username slug — lowercase, alpha only, grapheme-safe:
join(graphemes(lower(userName)) |filter: isAlpha($item), "")

// 3. Password validation — visual length for UX, byte length for storage:
len(password) >= 8 && len(password) <= 72

// 4. Per-name processing in a list:
names |map: join(graphemes($item) |filter: isAlpha($item), "")

// 5. ASCII-only protocol header check:
bytes(header) |every: $item < 128

// 6. Strip combining marks across a word list:
words |map: join(runes($item) |filter: $item < "\u0300" || $item > "\u036F", "")

// 7. Initials extraction (grapheme-safe):
name |split: " " |map: graphemes($item)[0] |join: ". "

// 8. Emoji-safe read-more truncation:
len(body) > 280 ? graphemeSubstr(body, 0, 279) + "…" : body

```

> **Note:** `bytes()` output (float64 values) is intended for analysis and protocol checks only, not for string reassembly in expressions. Use a UDF if you need to reconstruct a string from bytes.

## User-Defined Function (UDF) Integration

UDFs always receive **plain Go strings** — no level information, no special types. This is the complete contract:

```go
// Host developer writes normal Go — nothing changes:
functions["upper"] = func(args ...any) (any, error) {
    return strings.ToUpper(args[0].(string)), nil
}

functions["truncate"] = func(args ...any) (any, error) {
    s := args[0].(string)
    n := int(args[1].(float64))
    runes := []rune(s)          // host decides their own Unicode level
    if len(runes) > n {
        runes = runes[:n]
    }
    return string(runes), nil
}
```

Level-aware decomposition happens in the UExL expression **before** the UDF is called:

```javascript
// UDF gets individual grapheme strings — plain, no awareness needed:
graphemes(name) |map: myTransform($item)       // myTransform receives "c", "a", "f", "é"

// Level-aware work done in expression, UDF gets clean result:
myUDF(join(graphemes(name) |filter: isAlpha($item), ""))

// Built-in functions handle level-aware measuring; UDF handles content:
graphemeLen(userName)          // VM built-in — no UDF needed
upper(userName)                // UDF gets plain string ✅
```

**The rule:** Built-in `grapheme*` functions handle level-aware structural operations. UDFs handle content transformation. No mixing required, no new function signatures needed.

## Behavior and Edge Cases

- **Index origin**: Zero-based for all indexing and array operations
- **Out of range**: `substr`/`graphemeSubstr` clamp indices; single-element access returns `null`
- **Empty ranges**: Return `""` or `[]`
- **Non-string inputs**: Type error for `graphemeLen`, `graphemes`, `runes`, `bytes`
- **`graphemeSubstr` out of range**: indices are clamped to string length
- **`substr` out of range**: errors (Go-compatible behavior)
- **ASCII fast-path**: `graphemeLen` and `graphemes` detect ASCII-only strings and skip UAX #29 segmentation (O(n) scan, then O(1) per operation)
- **Immutability**: All operations return new values; strings are never mutated

## Design Principles

- **Byte default**: Matches Go — zero surprises for Go developers, consistent with standard library expectations
- **Explicit over magic**: `graphemeLen` is unambiguous; `runes(s)` makes the conversion explicit
- **No hidden state**: Pure functions, no tagged values, no level propagation rules
- **UDF-transparent**: Level awareness lives entirely in expression syntax, never in function signatures
- **Composable**: `graphemes()` / `runes()` / `bytes()` compose with every existing UExL pipe operation
- **Dependency-honest**: Default operations need nothing new; grapheme operations require one UAX #29 library, isolated to `grapheme*` functions
- **The rule:** Built-in `grapheme*` functions handle level-aware structural operations. UDFs handle content transformation.

## String Operation Reference

This table summarises the semantics of all string operations. The `s[i]` and `s[i:j]` operators are byte-based; the named functions provide explicit Unicode levels.

| Operation | Level | Notes |
|-----------|-------|------|
| `len("naïve")` | byte | 6 bytes |
| `"naïve"[2]` | byte | raw 3rd byte (may split a codepoint) |
| `"naïve"[0:3]` | byte | first 3 bytes (may split a codepoint) |
| `substr("naïve",0,3)` | byte | first 3 bytes |
| `runeLen("naïve")` | rune | 5 code points |
| `runeSubstr("naïve",0,3)` | rune | `"naï"` — safe for Unicode |
| `graphemeLen("👨\u200d👩\u200d👧\u200d👦")` | grapheme | 1 visible character |
| `graphemeSubstr(s,i,n)` | grapheme | display-safe substring |
| `runes("café\u0301")` | rune | `["c","a","f","e","́"]` |
| `graphemes("café")` | grapheme | `["c","a","f","é"]` |
| `bytes("hi")` | byte | `[104, 105]` as float64 |
| `join(arr, sep)` | — | reassemble any array into a string |
