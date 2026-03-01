# Grapheme-Aware String Implementation Strategy

## Design Goal

Implement multi-level Unicode string operations in a way that is:
1. **Easy to implement**: Minimal changes to existing codebase — pure function additions, no VM state changes
2. **UDF-transparent**: User-defined functions receive plain `string` — no awareness of levels needed
3. **Performance-conscious**: ASCII fast-path throughout; grapheme segmentation only when explicitly requested
4. **Zero-panic robustness**: Maintains UExL's error handling philosophy
5. **Industry-standard naming**: Matches conventions familiar to Python/Ruby/Go developers

## Chosen Design: Explicit Function Approach

### Why Not the Tagged String / Metadata Approach

The earlier design stored a `StringLevel` tag on `Value` and propagated it through operations. This was rejected for three reasons:

1. **Allocation cost**: `stringWithMeta` cannot go into `Value.StrVal` — it boxes into `Value.AnyVal` as `any`, costing a heap allocation on every `char()` / `utf8()` call. Undoes Phase 8 inlining gains.
2. **UDF boundary**: Level silently resets at every host function call (host sees plain `string`). Invisible, untestable. Users can't reason about propagation.
3. **Complexity**: Mixed-level arguments require propagation rules, dominant-level heuristics, and documented edge cases. Every VM operation needs an extra type-switch arm.

### The Chosen Approach: Explicit Pure Functions

No new VM state. No tagged values. No propagation rules. Level-aware work is done by explicit built-in functions that return arrays or strings. UDFs always receive plain strings.

```
Expression                     VM operation                  Host function
──────────────────────────────────────────────────────────────────────────
graphemes(name)            →   segmentGraphemes() → []any   (none)
  |map: myUDF($item)       →   pipe each element           myUDF(string) → string
  |join: ""                →   join array                  (none)
```

## Architecture

### What Changes vs v1

| Component | v1 actual | v2 Change |
|-----------|-----------|----------|
| `len(s)` | byte count (`builtinLen`) | **unchanged** |
| `s[i]` | rune-based (`[]rune`) | **→ byte-based** (standardized to match `len`/`substr`) |
| `s[i:j]` | rune-based (`[]rune`) | **→ byte-based** (standardized to match `len`/`substr`) |
| `substr(s,i,n)` | byte-based (`str[start:end]`) | **unchanged** |
| `Value` struct | no level field | **unchanged** |
| `VMFunctions` signature | `func(...any)(any,error)` | **unchanged** |
| New: grapheme measure | — | `graphemeLen(s)` |
| New: grapheme cut | — | `graphemeSubstr(s,i,n)` |
| New: rune measure | — | `runeLen(s)` |
| New: rune cut | — | `runeSubstr(s,i,n)` |
| New: explode functions | — | `graphemes`, `runes`, `bytes` |
| New: reassemble function | pipe `\|join:` only | `join(arr, sep?)` as function |

### Function Inventory

| Category | Function | Returns | Level | Dependency |
|----------|----------|---------|-------|-----------|
| Measure | `len(s)` | float64 | byte | none — unchanged |
| Measure | `runeLen(s)` | float64 | rune | none — new |
| Measure | `graphemeLen(s)` | float64 | grapheme | UAX #29 |
| Cut | `substr(s, i, n)` | string | byte | none — unchanged |
| Cut | `runeSubstr(s, i, n)` | string | rune | none — new |
| Cut | `graphemeSubstr(s, i, n)` | string | grapheme | UAX #29 |
| Explode | `runes(s)` | []any | rune | none |
| Explode | `graphemes(s)` | []any | grapheme | UAX #29 |
| Explode | `bytes(s)` | []any | byte | none |
| Reassemble | `join(arr)` | string | — | none |
| Reassemble | `join(arr, sep)` | string | — | none |

## Implementation Plan

### Phase 1: Foundation — `internal/utils/unicode.go`

All Unicode processing logic lives in a single utility file, isolated from the VM.

```go
package utils

import "unicode/utf8"

// StringLevel is not needed in the VM — only used internally in this package.

// GraphemeLength returns grapheme cluster count with an ASCII fast-path.
func GraphemeLength(s string) int {
    if isASCII(s) {
        return len(s)
    }
    return segmentGraphemes(s)  // UAX #29 via github.com/rivo/uniseg
}

// GraphemeSlice returns graphemes [start : start+length].
func GraphemeSlice(s string, start, length int) (string, error) { ... }

// CollectRunes returns all runes as a slice of single-rune strings.
func CollectRunes(s string) []string {
    runes := []rune(s)
    result := make([]string, len(runes))
    for i, r := range runes {
        result[i] = string(r)
    }
    return result
}

// CollectGraphemes returns all grapheme clusters as strings.
func CollectGraphemes(s string) []string {
    if isASCII(s) {
        result := make([]string, len(s))
        for i := range s { result[i] = s[i : i+1] }
        return result
    }
    return collectGraphemesUAX29(s)  // via uniseg
}

// CollectBytes returns all UTF-8 bytes directly as []any of float64 — single pass, no intermediate slice.
func CollectBytes(s string) []any {
    result := make([]any, len(s))
    for i := range s { result[i] = float64(s[i]) }
    return result
}

func isASCII(s string) bool {
    for i := 0; i < len(s); i++ {
        if s[i] >= 128 { return false }
    }
    return true
}
```

### Phase 2: Built-in Functions — `vm/builtins.go`

Add all new functions as standard VMFunctions. They are identical in structure to existing built-ins.

```go
// Measure (grapheme-level only — len() is unchanged, byte-based)
"graphemeLen": func(args ...any) (any, error) {
    s, err := requireString(args, 1, "graphemeLen")
    if err != nil { return nil, err }
    return float64(utils.GraphemeLength(s)), nil
},

// Cut (grapheme-level only — substr() is unchanged, byte-based)
"graphemeSubstr": func(args ...any) (any, error) {
    s, start, length, err := requireStringIntInt(args, "graphemeSubstr")
    if err != nil { return nil, err }
    return utils.GraphemeSlice(s, start, length)
},

// Explode to array
"runes": func(args ...any) (any, error) {
    s, err := requireString(args, 1, "runes")
    if err != nil { return nil, err }
    strs := utils.CollectRunes(s)
    result := make([]any, len(strs))
    for i, r := range strs { result[i] = r }
    return result, nil
},

"graphemes": func(args ...any) (any, error) {
    s, err := requireString(args, 1, "graphemes")
    if err != nil { return nil, err }
    clusters := utils.CollectGraphemes(s)
    result := make([]any, len(clusters))
    for i, c := range clusters { result[i] = c }
    return result, nil
},

"bytes": func(args ...any) (any, error) {
    s, err := requireString(args, 1, "bytes")
    if err != nil { return nil, err }
    return utils.CollectBytes(s), nil  // returns []any directly — single pass
},

// Reassemble
"join": func(args ...any) (any, error) {
    if len(args) < 1 || len(args) > 2 {
        return nil, fmt.Errorf("join expects 1 or 2 arguments")
    }
    arr, ok := args[0].([]any)
    if !ok {
        return nil, fmt.Errorf("join: first argument must be an array, got %T", args[0])
    }
    sep := ""
    if len(args) == 2 {
        sep, ok = args[1].(string)
        if !ok {
            return nil, fmt.Errorf("join: separator must be a string, got %T", args[1])
        }
    }
    var sb strings.Builder
    sb.Grow(len(arr) * 4)  // rough pre-size: avoids realloc for typical short strings
    for i, v := range arr {
        s, ok := v.(string)
        if !ok {
            return nil, fmt.Errorf("join: element %d must be a string, got %T", i, v)
        }
        if i > 0 { sb.WriteString(sep) }
        sb.WriteString(s)
    }
    return sb.String(), nil
},
```

### Phase 3: Standardize `s[i]` and `s[i:j]` to Byte-Based

`executeStringIndex` and `sliceString` are changed from `[]rune`-based to byte-based, so all four primitive string operations are consistent:

| Operation | Implementation | Level |
|-----------|---------------|-------|
| `len(s)` | Go `len(s)` — O(1) | byte |
| `substr(s,i,n)` | `s[i:i+n]` — unchanged | byte |
| `s[i]` | `s[i:i+1]` — was `[]rune` | **byte** |
| `s[i:j]` | `[]byte(s)` accumulation — was `[]rune` | **byte** |

This matches Go's native convention and UDF expectations (e.g. `strings.Index` returns byte offsets). No changes to `vm_handlers.go`, `vm_utils.go`, hot-path stack operations, or Phase 8 inlining gains.

### Phase 4: Dependency — `github.com/rivo/uniseg`

Add to `go.mod`. This is the only new dependency for the entire feature. It is:
- Used exclusively by `grapheme*` functions
- Not imported by any VM hot-path code
- Well-maintained, UAX #29 compliant, used by many Go projects

```bash
go get github.com/rivo/uniseg
go mod vendor
```

## UDF Integration

No changes to `VMFunctions` signature. The host contract is unchanged:

```go
// Existing host functions work as-is:
functions["upper"] = func(args ...any) (any, error) {
    return strings.ToUpper(args[0].(string)), nil
}
```

UDFs receive `string` elements from `graphemes()`/`runes()` arrays via pipes:

```javascript
// myTransform receives plain strings "c", "a", "f", "é" — no changes needed:
graphemes(name) |map: myTransform($item) |join: ""
```

## Performance Characteristics

| Function | Time | Allocation | Notes |
|----------|------|-----------|-------|
| `len(s)` | O(1) | 0 | byte count — unchanged, Go-compatible |
| `graphemeLen(s)` ASCII | O(n) scan | 0 | fast-path |
| `graphemeLen(s)` Unicode | O(n) UAX #29 | small | uniseg |
| `runes(s)` | O(n) | O(n) []any | unavoidable |
| `graphemes(s)` ASCII | O(n) | O(n) []any | fast-path |
| `graphemes(s)` Unicode | O(n) UAX #29 | O(n) []any | uniseg |
| `bytes(s)` | O(n) | O(n) []any | unavoidable |
| `join(arr, sep)` | O(n) | 1 strings.Builder | same as \|join: |

## Implementation Checklist

### Phase 1: Unicode Utilities
- [ ] Create `internal/utils/unicode.go`
- [ ] Implement `isASCII()`, `GraphemeLength()`, `CollectGraphemes()` with ASCII fast-path
- [ ] Implement `CollectRunes()`, `CollectBytes()` (single-pass `[]any` output)
- [ ] Implement `GraphemeSlice(s, start, length int)`
- [ ] Add `github.com/rivo/uniseg` to go.mod + vendor

### Phase 2: Built-in Functions
- [ ] Add `graphemeLen` to `vm/builtins.go`
- [ ] Add `graphemeSubstr` to `vm/builtins.go`
- [ ] Add `graphemes`, `runes`, `bytes` to `vm/builtins.go`
- [ ] Add `join(arr, sep?)` to `vm/builtins.go` — strings only, type error for non-string elements
- [ ] Register all new functions in `Builtins` map
- [ ] No changes to `builtinLen`, `builtinSubstr`, or `executeIndexValue` — byte-based behavior retained

### Phase 3: Tests
- [ ] Unit tests for `internal/utils/unicode.go` — ASCII fast-path, decomposed strings, emoji
- [ ] Confirm `len("naïve")` = 6 (bytes, `builtinLen` unchanged)
- [ ] Confirm `"naïve"[2]` = `"\xAF"` (3rd UTF-8 byte — raw byte, `executeStringIndex` now byte-based)
- [ ] Confirm `"naïve"[0:3]` = first 3 bytes = `"na\xC3"` (splits ï, `sliceString` now byte-based)
- [ ] Confirm `"naïve"[0:4]` = `"naï"` (first 4 bytes — valid boundary, ï is bytes 3-4)
- [ ] Confirm `substr("naïve", 0, 3)` = first 3 bytes (`builtinSubstr` byte-based — unchanged)
- [ ] Confirm `runes("naïve")[2]` = `"ï"` (rune-level access via explicit explode)
- [ ] Confirm `len(runes("naïve"))` = 5 (rune count via explicit explode)
- [ ] VM integration tests for each new builtin
- [ ] Edge cases: empty string, out-of-bounds (clamped for substr, error for single index)
- [ ] `byteSubstr` explicitly produces invalid UTF-8 — document and test
- [ ] `join` type error: non-string element returns error (not silent conversion)
- [ ] `join` with empty array = ""
- [ ] `graphemeLen`/`graphemes` with decomposed + precomposed + emoji sequences
- [ ] Negative indices: `substr("hello", -2, 2)` = "lo"
- [ ] Pipe integration: `graphemes(s) |map: upper($item) |join: ""`
- [ ] UDF integration: `graphemes(s) |map: myUDF($item)` (UDF receives plain string)

### Phase 4: Documentation
- [ ] `book/v2/accessing-graphemes.md` — user-facing (done)
- [ ] `book/functions/` — per-function reference entries
- [ ] `wip-notes/grapheme-implementation-strategy.md` — this file (done)
