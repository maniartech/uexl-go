# Accessing Graphemes within Strings (v2)

Working with human text requires care: what users perceive as a single “character” can be composed of multiple Unicode code points (e.g., a base letter plus one or more combining marks). Indexing by raw bytes or even by code points can split these. Grapheme clusters model user‑perceived characters.

This note describes the recommended v2 semantics and proposed built‑ins for grapheme‑aware operations in UExL.

## Background: Bytes vs Code Points vs Graphemes
- Byte: a raw unit of storage. Not meaningful for text indexing in Unicode.
- Code point (rune): a single Unicode scalar value. Many scripts require multiple code points for one displayed “character.”
- Grapheme cluster: one user‑perceived character, possibly composed of multiple code points (as defined by Unicode Text Segmentation, UAX #29).

Examples
- "café": the "é" can be a single code point U+00E9.
- "éclair": the initial "é" can be two code points: U+0065 ("e") + U+0301 (COMBINING ACUTE ACCENT).

## Default indexing in UExL
- Bracket indexing on strings remains position‑based and returns a single character string.
- This is suitable for simple ASCII but may not align with a full grapheme for combined characters.
- Out‑of‑bounds indexing returns `null`. Strings are immutable.

For full user‑perceived character handling, use the grapheme‑aware helpers below.

## Proposed grapheme built‑ins
To work with grapheme clusters explicitly, v2 proposes the following built‑ins:

- graphemeAt(s, i)
	- Returns the i‑th grapheme cluster of string `s` as a string.
	- Index is zero‑based. If `i` is out of range, returns `null`.
- lenGraphemes(s)
	- Returns the number of grapheme clusters in `s` as a number.
- graphemeSlice(s, start, end)
	- Returns the substring consisting of graphemes in the half‑open range `[start, end)`.
	- Indices are zero‑based and clamped to the valid range. If `start >= end`, returns an empty string.

Optional convenience:
- graphemes(s)
	- Returns an array of grapheme cluster strings for `s`.

Note: These functions focus on stable behavior for end‑user text. They avoid splitting combined characters.

## Behavior and edge cases
- Index origin: zero‑based for all indexing functions.
- Out of range: `graphemeAt` returns `null`. `graphemeSlice` clamps indices; empty range yields `""`.
- Non‑string inputs: type error.
- Immutability: all functions return new values; strings aren’t mutated.

## Examples
Basic indexing differences
```
"café"[3]            // "é" (works for precomposed character)
"éclair"[1]          // "́" (combining mark only; not a full grapheme)

graphemeAt("éclair", 0)  // "é" (full user‑perceived character)
lenGraphemes("éclair")   // 6 (counts user‑perceived characters)
```

Slicing by graphemes
```
// Take the first 2 graphemes
graphemeSlice("नमस्ते", 0, 2)  // language‑dependent clusterization; returns the first two graphemes

// Map words to their initial grapheme (safe for accents)
words |map: graphemeAt($1, 0)
```

Producing arrays of graphemes
```
// If graphemes() is available
letters = graphemes("été")   // ["é", "t", "é"]
letters[0]                     // "é"
```

## Performance considerations
- Grapheme segmentation is more expensive than simple indexing. Prefer it when working with end‑user text or when correctness matters (names, UI labels, search, truncation).
- For ASCII‑only data or purely programmatic tokens, default indexing may be sufficient.

## Migration tips
- If your code relies on visual characters (e.g., initials, previews, truncation), migrate to `graphemeAt`, `lenGraphemes`, or `graphemeSlice`.
- Keep the general rule: arrays/objects/strings use zero‑based indexing; out‑of‑bounds returns `null` for single‑element access.

## Implementation guidance (host)
While UExL specifies behavior, host implementations should use a Unicode segmentation library that follows UAX #29 to implement these built‑ins. In Go, this is commonly achieved with a grapheme iterator.

---
By separating default string indexing (fast, position‑based) from grapheme‑aware helpers (correct for user‑perceived characters), UExL provides both performance and correctness where it matters.
