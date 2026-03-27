# Design philosophy: explicit nullish and boolish semantics

UExL is deliberately explicit about how values are accessed and how defaults are applied. This keeps expressions predictable, safe, and easy to reason about across diverse data sources.

## Why explicitness?

- Avoid surprises: Real data often contains valid “falsy” values like 0, empty string, or false. We don’t silently replace those.
- Clear error boundaries: Accessing something that doesn’t exist should fail loudly unless you’ve asked for a softer behavior.
- Composability: Pure, explicit building blocks are easier to combine and optimize.

## Core principles

1) Strict access by default
- Property access (`a.b`) and index access (`a[i]`) throw when the key/index doesn’t exist or is out of bounds.
- The optional forms (`a?.b`, `a?.[i]`) guard the base from being nullish **and** soften a missing key/index. If the base is nullish, or the member/index is absent, the result is null and the rest of the chain is short-circuited.

2) Nullish means null or absent
- A value is **nullish** if it is `null` or absent (a context variable not provided by the caller).
- Absent context variables are treated as `null` — they are interchangeable in nullish checks.
- This mirrors the `null | undefined` convention found in other expression languages, making guards intuitive: `user?.name` works whether `user` is `null` or simply not in the context.
- Nullish coalescing (`a ?? b`) falls back when `a` is null or absent. It does not treat valid "falsy" values (0, "", false, empty array/object) as nullish.
- Logical operators (`||`, `&&`, `!`) work on truthiness. They’re for control flow, not for data defaulting.

3) Safe mode with nullish coalescing
- Only the immediate property/index access on the left of `??` is evaluated in “safe” mode.
- Example: `x.a.b ?? c` softens only the final access of `b` in `a`.
  - If `b` exists in `a` but is null → use `c`.
  - If `b` doesn’t exist in `a` at all → still an error.
  - If `a` doesn’t exist in `x` or is nullish (or any prior link fails) → still an error. Earlier links remain strict.

4) Short‑circuiting without side effects

- `a && b` evaluates `b` only if `a` is truthy; `a || b` evaluates `b` only if `a` is falsy.
- `a?.b` and `a?.[i]` do not evaluate the property or index expression when `a` is nullish.

5) Unicode-level explicitness

- String operations default to grapheme clusters (user-perceived characters) for correctness.
- When different Unicode levels are needed, use explicit view functions: `char()` for code points, `utf8()` for bytes, `utf16()` for code units.
- All operations work consistently across all views, maintaining composability.

## Explicitness policies (nullish and boolish)

- Strict access by default: property (`. …`) and index (`[i]`) access throw on missing keys or out-of-bounds indices.
- Optional access guards both a null base and a missing key/index: `a?.b`, `a?.[i]` return null when `a` is nullish **or** the member/index doesn't exist.
- Nullish means null or absent:
  - Absent context variables (not provided by the caller) are treated as `null`.
  - Use `a ?? b` to default only when `a` is null or absent; it preserves valid falsy like `0`, `""`, `false`.
  - Use `||`, `&&`, and `!` for truthiness-based control flow, not for data defaulting.
- Safe mode with nullish coalescing: `??` provides safety for the immediate property/index access only
  - `x.a.b ?? c` is safe for accessing `b` in `a`, but not for accessing `a` in `x`
- Short-circuit evaluation without side effects: `a?.[i]` does not evaluate `i` when `a` is nullish; logical ops only evaluate the right side when needed.

## Practical guidance and examples

- Defaulting with nullish coalescing keeps valid falsy values:
  - `count ?? 0` → use a default only when `count` is null or absent.
- Avoid using `||` for defaults when falsy values are meaningful:
  - `count || 0` would replace 0 with 0 again (fine) but also replace "" or false, and can be wrong in other contexts.
- Safe property access with nullish coalescing (immediate access only):
  - `user.name ?? "Anonymous"` → provides fallback when `name` property is missing or null in `user`
  - `data.user.name ?? "Anonymous"` → safe for `name` in `user`, but `user` must exist in `data`
- Optional chaining on present but incomplete objects:
  - `user?.address?.city ?? "unknown"` → `"unknown"` when `user` exists but has no `address` key; `?.` softens both the null base and the absent key.
  - `(user?.address).city` can still error if `user` exists but `address` is missing — `.city` is a strict access on the null result.
- Unicode-level operations are explicit about their target level:
  - `len("éclair")` → 6 graphemes (safe for user display)
  - `len(char("éclair"))` → explicit code point count when needed
  - `len(utf8("éclair"))` → explicit byte count for protocols/storage

## Precedence notes

- Access operators (`.`, `[ ]`, `?.`, `?[ ]`) bind tighter than `??`, `||`, and `&&`.
- `??` binds tighter than `||` and `&&`. Parenthesize for readability when mixing.

## Related reading

- Null chaining operator: `book/operators/null-chaining.md`
- Nullish coalescing operator: `book/operators/nullish-coalescing.md`
- Strings and Unicode: `book/strings-unicode.md`
- Mutability and purity: `book/mutability.md`
