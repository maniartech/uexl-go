# Glossary

A quick reference for UExL terminology and concepts.

## Operators and evaluation

- Null chaining (?.)
  - Also known as optional chaining (JS/TS), null‑aware access (Dart), safe navigation.
  - Short-circuits property/index access when the left-hand side is nullish; returns null.
  - See: v2/null-chaining-operator.md
- Nullish coalescing (??)
  - Returns the left value unless it is nullish; otherwise returns the right.
  - Keeps valid falsy (0, "", false).
  - See: v2/nullish-coalescing-operator.md
- Short-circuiting
  - Logical ops (||, &&) and null/optional access may stop evaluating later parts when the result is already determined.
  - Example: a && b only evaluates b if a is truthy; a?.[i] only evaluates i if a is non‑nullish.
- Chaining (member/index)
  - Applying successive property/index accesses: a.b[c].d
  - With null chaining: a?.b?.[i]?.d
- Precedence and associativity
  - Defines how expressions group without parentheses. Access (., [], ?., ?[]) binds tighter than ??, ||, &&, and ?:
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
  - Zero values by type: number 0; string ""; boolean false; empty array []; empty object {}.
  - Null and unavailable/missing are nullish and therefore falsy.
  - Logical operators (||, &&) rely on this notion of truthiness.
- Nullish
  - A value is nullish if it is `null` or unavailable/missing/undefined (e.g., absent key, out‑of‑bounds index, unresolved identifier).
  - Operators (??) and (?. / ?[ ]) treat nullish as “absent” for fallback and short‑circuiting.
  - See: v2/nullish-coalescing-operator.md and v2/null-chaining-operator.md.

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
