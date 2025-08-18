# Nullish coalescing operator (??)

The nullish coalescing operator returns the left-hand value if it is not nullish, otherwise it returns the right-hand value. Use it to provide safe defaults without treating valid “falsy” values (like `0`, `""`, or `false`) as missing.

In UExL, “nullish” means:
- `null`, or
- a lookup that resolves to no value (an undefined identifier resolves to `null` in UExL).

Short-circuiting: the right-hand side is only evaluated if the left-hand side is nullish.

## Why not use || for defaults?

`||` is a truthiness operator: it replaces any falsy value. That’s often wrong for real data where `0`, `false`, and `""` are meaningful.

- `0 || 10` → `10` (replaced, maybe wrong)
- `0 ?? 10` → `0` (preserved, correct)

Prefer `??` when you only want to fall back for truly missing values.

## Common applications and use cases

1) Default for possibly missing variables
- `userName ?? "Anonymous"`
- `pageSize ?? 25`

2) Property or index with a safe default
- `user ? (user.nickname ?? user.name) : "Anonymous"`
- `arr[10] ?? 0`  // out-of-bounds yields `null` → default to 0

3) Configuration fallbacks (layered defaults)
- `config.timeout ?? defaults.timeout ?? 5000`
- `options.locale ?? env.DEFAULT_LOCALE ?? "en"`

4) Normalizing external/API data
- `payload.meta.version ?? 1`
- `payload.address.city ?? "Unknown"`

5) Function arguments with defaults
- `len(name ?? "")`
- `substr(text ?? "", 0, 10)`

6) Pipes: provide a default for a stage result
- `data |: $last.items ?? [] |map: $item.name`
- `users |find: $item.id == targetId ?? null`  // if not found → `null`

7) Reduce: initialize accumulator safely (see notes below)
- `[1,2,3,4] |reduce: ($acc ?? 0) + $item`
- `[1,2,3,4] |reduce: set($acc ?? {}, $index, $item)`

Tip: When available in your version, prefer the explicit initializer form for reduce:
- `[1,2,3,4] |reduce(0): $acc + $item`
- `[1,2,3,4] |reduce({}): set($acc, $index, $item)`

8) Chained fallbacks (progressively relax)
- `user.nickname ?? user.name ?? "Anonymous"`
- `request.query.limit ?? request.body.limit ?? 50`

## Patterns and recipes

- Guarded property access with a single default:
	- `user ? (user.email ?? "") : ""`

- Default-after-pipeline:
	- `orders |filter: $item.status == 'paid' ?? []`
	- `orders |reduce: ($acc ?? 0) + $item.total ?? 0`

- Final fallback for an entire chain:
	- `(source |: $last.items |map: $item.value) ?? []`

## Edge cases and best practices

- Property access on `null` will fail before `??` runs. Guard the access:
	- Bad: `user.name ?? "N/A"` when `user` may be `null`
	- Good: `user ? (user.name ?? "N/A") : "N/A"`

- Distinguish falsy vs nullish:
	- `false ?? true` → `false`
	- `"" ?? "(empty)"` → `""`
	- `0 ?? 1` → `0`

- Reduce and empty inputs:
	- If you don’t provide an explicit initializer, a reduce over an empty input returns `null`.
	- Add a default afterward or use an initializer:
		- `values |reduce: ($acc ?? 0) + $item ?? 0`
		- `values |reduce(0): $acc + $item`

- Mixing with logical operators:
	- When combining `??` with `||` or `&&`, use parentheses for clarity and to avoid ambiguity:
		- `(a ?? b) || c`
		- `a && (b ?? d)`

- Performance: `??` short-circuits. If the left side is not nullish, the right side isn’t evaluated.

## Quick comparison: ?? vs ||

- Use `??` to fill in only when a value is missing (nullish).
- Use `||` to replace when a value is falsy (including `0`, `false`, `""`, empty array/object).

## Related reading

- Operators overview and precedence
- Pipes overview (using `$last`, `$item`, `$index`, `$acc`)
- Reduce notes (initializer and nullish patterns)
- Data types (null)
- Context (undefined identifiers resolve to null)

