# Glossary

A quick reference for UExL terminology and concepts.

## Operators and evaluation

- Null chaining (?.)
  - Also known as optional chaining (JS/TS), null‑aware access (Dart), safe navigation.
  - Short-circuits property/index access when the left-hand side is nullish; returns null.
  - See: operators/null-chaining.md
- Nullish coalescing (??)
  - Returns the left value unless it is nullish; otherwise returns the right.
  - Keeps valid falsy (0, "", false, the epoch datetime, a zero duration) — these are falsy but not nullish.
  - See: operators/nullish-coalescing.md
- Equality (==, !=)
  - Exact equality for primitives; no cross‑type coercion (e.g., "1" == 1 is false).
  - Deep equality for arrays and objects in both == and !=.
  - There are no ===/!== operators; use ==/!= for equality and rely on ?./?? for nullish flow.
  - Convert to boolean explicitly with !! when needed (double NOT via truthiness).
- Short-circuiting
  - Logical ops (||, &&) and null/optional access may stop evaluating later parts when the result is already determined.
  - Example: a && b only evaluates b if a is truthy; a?.[i] only evaluates i if a is non‑nullish.
- Chaining (member/index)
  - Applying successive property/index accesses: a.b[c].d
  - With null chaining: a?.b?.[i]?.d
- Precedence and associativity
  - Defines how expressions group without parentheses. Access (., [], ?., ?[]) binds tighter than `??`, `||`, `&&`, and `?:`. Also, `??` binds tighter than `||`/`&&` in UExL.
  - See: operators/precedence.md

## Pipes

- Pipe stage
  - A single transformation step: |:, |map:, |filter:, |reduce:, etc.
- Pipe chaining
  - Linking multiple stages; each stage receives the previous stage result.
- Emitted context variables
  - $last: previous stage value
  - $item: current element (map/filter/find/every/some/unique/sort/groupBy)
  - $index: current index
  - $acc: accumulator (reduce)
  - $window: current window (window)
  - $chunk: current chunk (chunk)
  - See: pipes/overview.md and pipes/types.md
- Predicate
  - The expression evaluated within a stage to decide or compute per element.
- Accumulator
  - The running value in reduce; updated each iteration by the predicate expression.

## Values and truthiness

- Truthiness (aka Boolish)
  - A value is truthy if it is non‑nullish and not the zero value of its type.
  - Zero values by type: number 0; string ""; boolean false; empty array []; empty object {}; datetime at the epoch (1970‑01‑01T00:00:00.000Z, i.e. 0 ms); duration of 0 ms.
  - Null and unavailable/missing are nullish and therefore falsy.
  - Logical operators (||, &&) rely on this notion of truthiness.
- Boolean conversion
  - Use !!x to coerce any value to a boolean via truthiness; !x yields the opposite.
- Nullish
  - A value is nullish if it is `null` or unavailable/missing/undefined (e.g., absent key, out‑of‑bounds index, unresolved identifier).
  - Operators (??) and (?. / ?[ ]) treat nullish as “absent” for fallback and short‑circuiting.
  - See: operators/nullish-coalescing.md and operators/null-chaining.md.

## Mutability and purity

- Pure by default
  - Expressions read and compute; they do not mutate ambient data.
- Update helpers return copies
  - set(obj, key, val) returns a new object with the key set; input is not mutated.
- No assignment or ++/-- operators
  - Mutation is not expressed via operators in UExL.
  - See: mutability.md

## Context and identifiers

- Context
  - The environment (map/object) providing values for identifiers during evaluation.
- Identifiers
  - Names that resolve in the current context; missing names are treated as null.

## Error handling

- Short-circuit prevents some errors
  - a?.b avoids errors when a is nullish; returns null instead.
- Normal access errors still apply when base is non‑nullish
  - Accessing properties on incompatible types may error per language rules.

For deeper dives, follow the links in each section to the full documentation pages.


## Special numeric values (enabled by default)

- NaN (Not‑a‑Number)
  - An IEEE‑754 special value representing an undefined or unrepresentable number.
  - Available as a literal by default. By IEEE‑754 rules, `NaN != NaN` is true, and comparisons with `NaN` are false.
  - Truthiness: numbers are truthy except zero; `NaN` is treated as a number and therefore truthy in truthiness checks.
- Inf (Infinity)
  - An IEEE‑754 special value representing positive infinity; write as `Inf`. Negative infinity is expressed as `-Inf` (unary minus).
  - Enabled by default (configurable); there is no `Infinity` literal.
  - Ordering: `-Inf < any finite number < +Inf`.
  - See `vm/ieee754-semantics.md` for operator semantics with infinities and NaN.

## Temporal values (datetime, duration)

- datetime
  - A value type representing an instant in time, stored as milliseconds since the Unix epoch (1970‑01‑01T00:00:00.000Z, UTC). Zoneless: it identifies an absolute instant, not a wall‑clock‑in‑a‑zone.
  - Falsy only at the epoch (0 ms); every other instant is truthy. Never nullish.
- duration
  - A value type representing an exact span of time in milliseconds (may be negative, e.g. from subtracting a later instant from an earlier one).
  - Falsy only when 0 ms; every other span is truthy. Never nullish.
- Status
  - The datetime/duration types, their **literals** (`d"…"`, and `7d`/`30ms`/`1.5h`), and their truthiness are implemented and work today. Temporal **operators** (`date − date`, `date ± duration`, comparisons) and the datetime **functions** (`parseDate`, `formatDate`, `now`, …) are being introduced incrementally. See the [DateTime Specification](../specs/datetime-spec.md).
