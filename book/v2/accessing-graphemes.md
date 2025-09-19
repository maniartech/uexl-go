# Accessing Graphemes within Strings (v2)

Working with human text requires care: what users perceive as a single "character" can be composed of multiple Unicode code points (e.g., a base letter plus one or more combining marks). Indexing by raw bytes or even by code points can split these. Grapheme clusters model user‑perceived characters.

This note describes the recommended v2 semantics for grapheme‑aware operations in UExL using string views and grapheme-aware defaults.

## Background: Bytes vs Code Points vs Graphemes

- **Byte**: A raw unit of storage. Not meaningful for text indexing in Unicode.
- **Code point (rune)**: A single Unicode scalar value. Many scripts require multiple code points for one displayed "character."
- **Grapheme cluster**: One user‑perceived character, possibly composed of multiple code points (as defined by Unicode Text Segmentation, UAX #29).

**Examples:**

- `"café"`: The "é" can be a single code point U+00E9.
- `"éclair"`: The initial "é" can be two code points: U+0065 ("e") + U+0301 (COMBINING ACUTE ACCENT).

## Advantages of Grapheme-Aware Defaults

### International Text Safety

- **Prevents text corruption**: Never splits user-perceived characters
- **Consistent behavior**: Same results across all Unicode scripts and languages
- **Future-proof**: Works with new Unicode additions (emoji, combining marks, etc.)

### Predictable Operations

- **Intuitive length**: `len("👨‍👩‍👧‍👦")` returns 1 (one family emoji), not 7 code points
- **Safe truncation**: `"naïve"[0:3]` gives "naï", not "naï" + partial character
- **Reliable indexing**: `"éclair"[0]` always returns complete "é", never broken pieces

### Developer Experience

- **Correct by default**: No need to remember Unicode complexities for common operations
- **Explicit optimization**: Use views only when performance requires it
- **Composable**: All existing UExL operations work seamlessly with any view
- **No rewrite needed**: Existing code handles international text correctly automatically

### Cross-Platform Consistency

- **Same behavior everywhere**: Works identically on Go, JavaScript, Python, Java hosts
- **No platform surprises**: Eliminates differences in how platforms handle Unicode
- **Portable expressions**: UExL code behaves the same regardless of host language

### Comparison with Other Approaches

| Operation | UExL v2 (Grapheme-aware) | Traditional (Code points/Bytes) |
|-----------|-------------------------|--------------------------------|
| `len("naïve")` | 5 (always correct) | 5 or 6 (depends on encoding) |
| `"éclair"[0]` | "é" (complete) | "e" or broken char (unsafe) |
| `"👨‍👩‍👧‍👦".length` | 1 (family emoji) | 7 (fragmented) |
| Performance | Fast for ASCII, safe for Unicode | Fast but breaks international text |
| Migration | Zero code changes needed | Requires Unicode library adoption |

## UExL v2 String Processing Approach

UExL v2 provides a view-based approach to string processing that works with all existing operations:

### Default Behavior: Grapheme-Aware

By default, UExL string operations work at the grapheme (user-perceived character) level:

- `len("éclair")` - Returns 6 (grapheme count)
- `"éclair"[0:3]` - Returns first 3 graphemes as string
- `"éclair" |map: upper($item)` - Maps over graphemes
- `contains("naïve", "aï")` - Grapheme-aware search

### String Views for Different Levels

When you need to work at different Unicode levels, use view functions:

- `char("text")` - Code point/rune view
- `utf8("text")` - UTF-8 byte view
- `utf16("text")` - UTF-16 code unit view

### Universal Operations on Views

**All existing UExL operations work with any view:**

```javascript
// Length operations
len("éclair")           // 6 (graphemes)
len(char("éclair"))     // 6 (code points)
len(utf8("éclair"))     // 7 (UTF-8 bytes)

// Indexing operations
"éclair"[0:3]           // First 3 graphemes: "écl"
char("éclair")[0:3]     // First 3 code points
utf8("éclair")[0:3]     // First 3 bytes

// Pipe operations
"éclair" |map: upper($item)         // Over graphemes
char("éclair") |filter: $item != "é" // Over code points
utf8("éclair") |filter: $item < 128  // Over bytes

// String functions
contains("éclair", "é")             // Grapheme search
contains(char("éclair"), "é")       // Code point search
```

## Examples

### Default Grapheme-Aware Behavior

```javascript
// Length and indexing (grapheme-aware by default)
len("éclair")           // 6 (user-perceived characters)
"éclair"[0]             // "é" (complete grapheme)
"éclair"[0:3]           // "écl" (first 3 graphemes)

// Comparison with code-point view
len(char("éclair"))     // 6 (code points - same for this example)
char("éclair")[1]       // "c" (second code point)
char("café\u0301")[3]   // "́" (combining mark only)
```

### Practical Use Cases

```javascript
// User interface - safe text processing
userNames
  |filter: len($item) >= 2        // At least 2 visual characters
  |map: $item[0] + "."           // Safe initial extraction

// Data processing - code point level when needed
fileNames
  |map: char($item)              // Switch to code point view
  |filter: len($last) <= 255     // Technical filename limits
  |join: ""                      // Back to string

// Protocol/encoding - byte level
httpHeaders
  |map: utf8($item)              // Switch to UTF-8 byte view
  |filter: all($last, $item < 128) // ASCII-only headers
  |join: ""                      // Back to string
```

### Mixed Operations

```javascript
// Process user text safely, then apply technical constraints
userName
  |: $last[0:20]                 // Truncate to 20 graphemes (safe)
  |: char($last)                 // Switch to code point view
  |filter: isalnum($item)        // Keep only alphanumeric
  |join: ""                      // Result: cleaned username
```

## Behavior and Edge Cases

- **Index origin**: Zero‑based for all indexing operations
- **Out of range**: Returns `null` for single-element access; slicing clamps indices
- **Empty ranges**: Return empty string `""`
- **Non‑string inputs**: Type error
- **Immutability**: All operations return new values; strings aren't mutated
- **Performance**: ASCII-only strings can use fast-path optimizations automatically

## Migration and Design Guidelines

### For User-Facing Text

- **Default operations work correctly**: UExL's grapheme-aware defaults handle international text properly
- **No special syntax needed**: `len()`, indexing, and pipes work with user-perceived characters
- **Use for**: names, content, UI labels, search, truncation

### For Technical/Performance-Critical Code

- **Use view functions when needed**: `char()` for code points, `utf8()` for bytes
- **Use for**: identifiers, protocols, parsing, ASCII-only data, performance optimization

### Design Principles

- **Correctness by default**: Default behavior serves end users and prevents Unicode bugs
- **Performance when needed**: View functions provide explicit control without sacrificing safety
- **Composable**: All operations work with all views, maintaining UExL's pipe-friendly design
- **No function explosion**: One set of operations, multiple views - simple and consistent
- **International-first**: Built for global applications from day one, not as an afterthought
- **Swift-compatible**: Proven architecture used successfully in production systems

## Implementation Guidance (Host)

### Core Requirements

- **Default string operations**: Implement grapheme-aware behavior using Unicode UAX #29 segmentation
- **View functions**: `char()`, `utf8()`, `utf16()` should return string-like objects that work with all operations
- **Performance optimization**: Detect ASCII-only strings for fast-path operations when possible

### Implementation Notes

- Use established Unicode libraries (ICU, Go's runes + segmentation library, etc.)
- View objects should be lightweight wrappers that maintain the original string data
- Consider lazy evaluation for view conversions
- Ensure all existing UExL operations (indexing, slicing, pipes, functions) work transparently with views

---

This unified approach provides both correctness (grapheme-aware defaults) and performance (explicit views) while maintaining a clean, composable API that works with all existing UExL operations.
