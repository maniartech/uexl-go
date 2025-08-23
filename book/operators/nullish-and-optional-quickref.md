# Quick Reference: Optional Chaining (?.) and Nullish Coalescing (??)

This guide summarizes UExL’s handling of `?.` and `??` in one place.

## Core rules

- Optional chaining (?.)
  - Guards only the base being nullish; if the base is nullish, the whole access returns null.
  - Missing members on a non-nullish base still cause an error.
- Nullish coalescing (??)
  - When the left side is a member or index access, only the final access is evaluated in “soft” mode (missing → null).
  - Preceding chain links remain strict and can still error.
  - If the left value is non-null, the right side is not evaluated.

## Common patterns

1) Provide a default for a leaf property
- `user.name ?? "Anonymous"`
  - If `user` is missing or not an object → error
  - If `name` missing/null → "Anonymous"
  - Else → value of `user.name`

2) Provide a default for a leaf index
- `arr[i] ?? 0`
  - If `arr` missing/non-indexable → error
  - If index OOB → 0

3) Mix with optional chaining (guard an earlier link)
- `(order?.customer).address.city ?? "N/A"`
  - `?.` protects only `order` being nullish
  - Missing `address` or `city` still errors unless you place `??` at that exact access

4) Strict chain until the last step
- `x.y.a.b ?? 7`
  - Missing `x|y|a` → error
  - Missing/null `b` → 7

5) Falsy vs nullish
- `config.count ?? 0` preserves `0` if present; only null/missing triggers the default.

6) Computed index/property stays strict except the final step
- `table[computeIndex()] ?? ""` → computeIndex() must succeed; if final index OOB → ""

## Tips

- Place `??` at the exact leaf you want to provide a default for.
- Use `?.` to guard a base that may be nullish; it does not soften missing members.
- Parenthesize for clarity when mixing logical operators.
