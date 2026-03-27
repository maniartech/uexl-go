# Mutability and purity in UExL

UExL is a declarative expression language. By design, expressions read values and produce new values; they do not mutate existing data implicitly. This makes expressions predictable, composable, and easier to optimize.

This page summarizes the mutability policy and the practical patterns to work with data safely.

## Policy at a glance

- No assignment or increment/decrement operators: there is no `=`, `+=`, `-=`, `++`, `--`, etc.
- Pure by default: operators like `+ - * /`, property/index access (`.`, `[i]`), and their optional variants (`?.`, `?[i]`) do not mutate.
- Explicit updates via functions return copies: update helpers like `set(obj, key, val)` return a new object/array with the change applied; the input is not mutated.
- Deterministic evaluation: left-to-right evaluation within an expression, with short-circuiting for logical ops and null chaining; avoid re-evaluation of already computed subexpressions.

## Why this design?

- Predictability: no hidden writes, fewer surprising interactions between parts of an expression.
- Safety: expressions are safe to run in different contexts without risking accidental state changes.
- Composability and caching: pure expressions are easier to memoize and reorder when safe.
- Host responsibility: stateful changes are performed by the embedding application, or through explicit update functions that make intent clear.

## Reading vs updating data

Reading is always non-mutating:
- Member access: `user.name`, `order.items[0]`
- Null chaining (null-aware access): `user?.address?.city`, `list?.[0]`

Updating is explicit and returns a new value:
- `set(obj, key, value)` → returns a copy of `obj` with `key` set to `value`.
  - Numeric keys are coerced to strings when targeting objects.
  - When used in pipelines (e.g., `reduce`), you can accumulate a new object without mutating prior values.

Examples:
- Build an index map:
  - `[1,2,3,4] |reduce: set($acc ?? {}, $index, $item)`
  - Result: `{ "0": 1, "1": 2, "2": 3, "3": 4 }`
- Update by dynamic key without mutation:
  - `set(user ?? {}, fieldName, value)`

## Null chaining and side effects

Null chaining does not evaluate the right side when the base is nullish:
- `(obj)?.[computeIndex()]` → `computeIndex()` only runs if `obj` is non-nullish

This prevents accidental side effects when data is missing. See v2/optional-chaining-operator.md for full semantics.

## Pipes and purity

Pipes like `map`, `filter`, `find`, `every`, `some`, `sort` are expected to be pure with respect to their input; they compute derived values.

Pipes that aggregate, such as `reduce`, can use update helpers to build new structures:
- `[1,2,3,4] |reduce: set($acc ?? {}, $index, $item)`

Avoid mutating external state inside pipes; if needed, prefer explicit return values and let the host apply changes after evaluation.

## What’s intentionally not supported

- Assignment operators and in-place mutation in expressions.
- Increment/decrement (`++`, `--`).
- Hidden mutation via property/index access.

These choices keep UExL a safe, declarative layer atop your application’s state.

## Implementation notes (for readers of the engine)

- Builtins like `set` should not mutate their first argument; they must return a new value. If performance is a concern, consider copy-on-write or persistent data structures in the future.
- Null chaining must short-circuit without evaluating index/key expressions.
- Logical operators should short-circuit and return the conventional values (e.g., `||` returns first truthy, `&&` returns first falsy or last value).

## Future directions

- Add non-mutating variants for common operations (e.g., `assoc`, `dissoc`, `push`, `merge`) that always return new values.
- Consider namespacing side-effect-capable helpers (e.g., `update.set`) if you ever introduce mutating forms, to make effects obvious. For now, all provided helpers are non-mutating.
