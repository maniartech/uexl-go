# Reduce Accumulator Syntax

This note captures the recommended way to initialize and use the accumulator in the `|reduce:` pipe, plus a small syntax improvement proposal.

## Motivation

- We want a simple, safe way to provide an initial accumulator value.
- Current examples often use truthiness to initialize, which can be error‑prone for valid falsy values like `0`, `""`, or `false`.

## Current pattern (works, but risky)

Examples seen in tests and docs:

```uexl
"[1,2,3,4] |reduce: ($acc || 0) + $item"
"[1,2,3,4] |reduce: set($acc || {}, $index, $item)"
```

Issue: `||` is a truthiness operator. If `$acc` ever becomes a falsy-but-valid value, it will be replaced by the default on the next iteration (e.g., product after hitting `0`, or concatenation with an empty string).

## Recommended pattern (use nullish coalescing)

Prefer the nullish coalescing operator `??` to initialize only when the accumulator is not set (null/undefined), not when it is merely falsy:

```uexl
"[1,2,3,4] |reduce: ($acc ?? 0) + $item"
"[1,2,3,4] |reduce: set($acc ?? {}, $index, $item)"
```

This keeps valid falsy accumulator values intact.

## Proposal: explicit initializer syntax (recommended)

Add an explicit initializer to `reduce` to make the intent obvious and remove the need for inline `??` in most cases.

- Preferred form:

	```uexl
	"[1,2,3,4] |reduce(0): $acc + $item"
	"[1,2,3,4] |reduce({}): set($acc, $index, $item)"
	```

- Alternative forms (if parser preferences differ):
	- Trailing argument: `|reduce: $acc + $item, 0`
	- With-clause: `|reduce: $acc + $item with 0`

Back-compat: keep supporting the current `|reduce: <expr>` form. If no initializer is provided, use the “no-initializer” semantics defined below.

## Semantics

- With initializer (e.g., `|reduce(seed): expr`):
	- `$acc` starts as `seed`.
	- `$index` starts at `0` (first element), `$item` is the element at that index.
	- For an empty input, the result is `seed` and the body is not executed.

- Without initializer (current behavior when only `: expr` is provided):
	- `$acc` starts as the first element of the input sequence.
	- `$index` starts at `1` (second element).
	- For an empty input, the result is `null` (or consider an error—choose and document; returning `null` is friendlier for expressions).

Note: Using `??` inside the body remains valid and useful when you intentionally want a dynamic initializer.

## Examples

- Sum:
	```uexl
	"[1,2,3,4] |reduce(0): $acc + $item"      // preferred
	"[1,2,3,4] |reduce: ($acc ?? 0) + $item"  // current-compatible
	```

- Product:
	```uexl
	"[1,2,3,4] |reduce(1): $acc * $item"
	// Avoid: ($acc || 1) * $item — breaks if $acc becomes 0
	```

- Build object by index → value:
	```uexl
	"[1,2,3,4] |reduce({}): set($acc, $index, $item)"
	"[1,2,3,4] |reduce: set($acc ?? {}, $index, $item)"
	```

- Concatenate strings (safe for empty strings):
	```uexl
	"[\"a\", \"\", \"c\"] |reduce(\"\"): $acc + $item"
	"[\"a\", \"\", \"c\"] |reduce: ($acc ?? \"\") + $item"
	```

## Migration notes

- Replace truthy initialization with nullish initialization:
	- `$acc || 0` → `$acc ?? 0`
	- `$acc || {}` → `$acc ?? {}`
- Prefer the explicit initializer form for clarity going forward:
	- `|reduce(0): ...` or `|reduce({}): ...`

## Optional sugar (future)

Provide shorthands built on top of `reduce(seed):`:

- `|sum` → `|reduce(0): $acc + $item`
- `|product` → `|reduce(1): $acc * $item`
- `|toObject: keyExpr, valExpr` → internally uses `reduce({})` and `set($acc, keyExpr, valExpr)`
- `|keyBy: keyExpr` → last value per key (also `reduce({})`)
- `|groupBy: keyExpr` → array values per key (existing or planned)

These keep pipelines terse for common folds while remaining consistent with reduce semantics.

