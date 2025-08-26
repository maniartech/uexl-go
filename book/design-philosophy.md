# Design philosophy: explicit nullish and boolish semantics

UExL is deliberately explicit about how values are accessed and how defaults are applied. This keeps expressions predictable, safe, and easy to reason about across diverse data sources.

## Why explicitness?

- Avoid surprises: Real data often contains valid “falsy” values like 0, empty string, or false. We don’t silently replace those.
- Clear error boundaries: Accessing something that doesn’t exist should fail loudly unless you’ve asked for a softer behavior.
- Composability: Pure, explicit building blocks are easier to combine and optimize.

## Core principles

1) Strict access by default
- Property access (`a.b`) and index access (`a[i]`) throw when the key/index doesn’t exist or is out of bounds.
- The optional forms (`a?.b`, `a?.[i]`) only guard the base from being nullish; they do not soften missing keys/indices. If the base is non‑nullish but the member/index is not present, it’s still an error.

2) Separate nullish from falsy
- Nullish coalescing (`a ?? b`) only falls back when `a` is null. It does not treat valid “falsy” values (0, "", false, empty array/object) as missing.
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

## Explicitness policies (nullish and boolish)

- Strict access by default: property (`. …`) and index (`[i]`) access throw on missing keys or out-of-bounds indices.
- Optional access only guards the base: `a?.b`, `a?.[i]` prevent errors only when `a` is nullish; they do not soften missing members/indices when `a` exists.
- Separate nullish from falsy:
  - Use `a ?? b` to default only when `a` is null; it preserves valid falsy like `0`, `""`, `false`.
  - Use `||`, `&&`, and `!` for truthiness-based control flow, not for data defaulting.
- Safe mode with nullish coalescing: `??` provides safety for the immediate property/index access only
  - `x.a.b ?? c` is safe for accessing `b` in `a`, but not for accessing `a` in `x`
- Short-circuit evaluation without side effects: `a?.[i]` does not evaluate `i` when `a` is nullish; logical ops only evaluate the right side when needed.

## Practical guidance and examples

- Defaulting with nullish coalescing keeps valid falsy values:
  - `count ?? 0` → use a default only when `count` is truly null.
- Avoid using `||` for defaults when falsy values are meaningful:
  - `count || 0` would replace 0 with 0 again (fine) but also replace "" or false, and can be wrong in other contexts.
- Safe property access with nullish coalescing (immediate access only):
  - `user.name ?? "Anonymous"` → provides fallback when `name` property is missing or null in `user`
  - `data.user.name ?? "Anonymous"` → safe for `name` in `user`, but `user` must exist in `data`
- Optional access guards only the base being nullish, not missing members:
  - `(user?.address).city` can still error if `user` exists but `address` is missing.
  - Use explicit checks or host helpers for existence if you want to treat missing as acceptable.

## Precedence notes

- Access operators (`.`, `[ ]`, `?.`, `?[ ]`) bind tighter than `??`, `||`, and `&&`.
- `??` binds tighter than `||` and `&&`. Parenthesize for readability when mixing.

## Related reading

- v2/Null chaining operator: `book/v2/null-chaining-operator.md`
- v2/Nullish coalescing operator: `book/v2/nullish-coalescing-operator.md`
- Mutability and purity: `book/mutability.md`
