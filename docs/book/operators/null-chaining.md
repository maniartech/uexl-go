# Null chaining operator (?.)

Null chaining provides null‑aware access to properties and indices. It returns `null` if the left-hand side is nullish, and otherwise performs the access normally. This prevents errors when traversing partially missing data.

Also known as: optional chaining (JS/TS), null‑aware access (Dart), safe navigation operator, null‑conditional operator.

UExL scope: identifiers, object member access, and index access. The call form (e.g., `fn?.(args)`) is not supported.

Note on equality and nullish: Equality operators (`==`, `!=`) are exact for primitives and deep for arrays/objects; they do not treat “missing” specially. Use `?.` to read safely and pair it with `??` to provide defaults. Use `!!x` when you need to coerce a value to boolean before applying logical ops.

## Summary

- Syntax: `expr?.prop` and `expr?.[index]`
- Semantics: short-circuit to `null` if `expr` is nullish; otherwise, behave like normal `.` or `[]` access
- Notation: no spaces between `?` and `.` or `[` for null chaining to be recognized
- Precedence: same as normal property/index access, and tighter than `??`, `||`, `&&`, and `?:`
- Short-circuiting: the right side (including the index expression) is not evaluated if the left side is nullish
- No call form: `fn?.(args)` is not part of this spec
- Scope of use: may appear after any expression result; parenthesize when chaining on the entire expression value (e.g., `(a ?? b)?.c`, `(cond ? a : b)?.[i]`, `(data |: $last.items)?.[0]`).

Nullish in UExL means:
- `null`, or
- a lookup that resolves to no value (an undefined identifier resolves to `null` in UExL)

## Motivation

Real-world data often has optional or partially missing fields. Null chaining makes access concise and safe without conflating “missing” with valid falsy values (`0`, `false`, `""`). It also composes naturally with the nullish coalescing operator `??` to provide layered fallbacks.

## Syntax and grammar

EBNF-style (informal) additions:

- Optional member access:
	- Postfix: `PostfixExpression := PostfixExpression '?.' Identifier`

- Optional index access:
	- Postfix: `PostfixExpression := PostfixExpression '?[' Expression ']'`

Notes:
- There must be no whitespace between `?` and `.` or `[` for the operator to be recognized.
- Normal access (`.` and `[expr]`) and optional access (`?.` and `?[expr]`) can be mixed in any order in a single chain.
- The call form is explicitly excluded.

### Chaining after any expression

Null chaining applies to the result of any expression, not just identifiers. Because `?.`/`?[ ]` are postfix and high‑precedence, add parentheses when you intend to chain off the whole expression:

- `(a ?? b)?.c`
- `(cond ? a : b)?.[i]`
- `(f(x))?.prop`  // allowed; the function call’s result may be object/array
- `(data |: $last.items)?.[0]`

## Semantics

For `A?.B` (member) or `A?.[I]` (index):
1) Evaluate `A`.
2) If `A` is nullish, the whole optional access expression evaluates to `null`. Do not evaluate `B` or `I`.
3) Otherwise, evaluate the access as if it were `A.B` or `A[I]` (using existing UExL rules). The result of that access is returned.

Short-circuiting examples:
- `null?.x` → `null` (no error)
- `user?.name` → `null` if `user` is `null`
- `arr?.[10]` → `null` if `arr` is `null`; otherwise normal indexing rules apply
- `obj?.[expensive()]` → `null` if `obj` is nullish; `expensive()` is not evaluated

Mixing with normal access:
- `a?.b.c` → `(a?.b).c` (left-to-right)
- `a.b?.c` → `a.b` evaluated first; then optional access to `.c`
- `a?.[i]?.j?.[k]` → each step short-circuits if the base at that step is nullish

Evaluation order and side effects:
- Evaluate the base exactly once at each optional step.
- If a step short-circuits (base is nullish), do not evaluate that step’s property/index part; later steps are not evaluated.
- For index access, the index expression is evaluated only if its base is non‑nullish.

### Relationship to nullish coalescing (??)

Null chaining reads safely; nullish coalescing provides defaults.

- `user?.nickname ?? "Anonymous"`
- `(source |: $last.items)?.[0]?.value ?? "N/A"`

Precedence examples (implicit parentheses):
- `a?.b ?? d` → `(a?.b) ?? d`
- `a && b?.c` → `a && (b?.c)`
- `a?.b.c` → `(a?.b).c`

## Common applications and examples

1) Deep property traversal
- `payload?.address?.city ?? "Unknown"`
- `user?.profile?.stats?.followers ?? 0`

2) Safe array access
- `users?.[0]?.email`
- `(data |: $last.items)?.[0] ?? null`

3) Layered defaults with ??
- `config?.timeout ?? defaults?.timeout ?? 5000`

4) Pipelines
- `orders |: $last.items |map: $item?.total ?? 0`
- `(source |: $last.items)?.[0]?.value ?? "N/A"`

5) Optional keys in heterogenous data
- `row?.[dynamicKey] ?? null`
- `record?.details?.[locale] ?? record?.details?.["en"] ?? null`

## Do and Don’t

Do
- Use `?.` and `?[ ]` to safely read values that may be missing
- Combine with `??` to supply a fallback only when missing
- Add parentheses when combining with `??`, `||`, `&&`, or `?:` for clarity
- Expect short-circuiting to skip expensive index computations

Don’t
- Don’t rely on `?.` to call functions (call form is not supported)
- Don’t write a standalone unary identifier form like `user?` (not supported)
- Don’t insert whitespace between `?` and `.`/`[` if you intend optional chaining
- Don’t expect `a?.b.c` to be fully safe; it equals `(a?.b).c` and may still error if `.c` is applied to `null`. Use `a?.b?.c` or add a final `??` fallback.

## Edge cases and details

- Index expression is not evaluated on short-circuit:
	- `(null)?.[expensive()]` → `null`, no side effects

- Missing keys and out-of-bounds indices normalize to `null` in UExL; null chaining preserves that behavior:
	- `arr?.[999]` → `null` (either from short-circuit on a nullish base or from missing index normalization)

- Access on non-object/non-array follows normal rules once the left is non-nullish:
	- If `A` is non-nullish but not indexable, `A?.[i]` behaves like `A[i]` and may error according to normal semantics

- Function-valued properties:
	- `obj?.fn` returns the function value or `null`. Call form `obj?.fn(args)` is not part of this spec.

- Mixing with ternary `?:` and `??`:
	- When in doubt, parenthesize: `(a?.b) ?? c`, `a ? (b?.c) : d`

- Identifiers and undefined:
	- UExL treats undefined identifiers as `null`. If any step in `a?.b?.c` resolves to `null`, the rest of the chain short-circuits.

## Operator precedence (relative)

From tighter to looser around optional access:
1) Member/index access (including `?.`, `?[ ]`)
2) Unary operators and power (`**` binds tighter than unary)
3) Multiplicative/additive and shifts
4) Nullish coalescing `??`
5) Comparisons and equality
6) Bitwise and logical `&&`, `||`
7) Ternary `?:`

Property/index access binds very tightly; `??` binds tighter than logical operators in UExL. Parenthesize for readability in complex expressions.

## Interactions with pipes

Null chaining is often used with the pipeline context variables:

- `data |: $last.items` may produce `null` if `items` is missing; follow with `?.[0]` safely:
	- `(data |: $last.items)?.[0] ?? null`

- In mappers/filters, especially when items may be partial:
	- `orders |map: $item?.total ?? 0`

## Implementation notes (engine)

Tokenization
- Add composite tokens for optional access to unambiguously distinguish from `?`, `?:`, and `??`:
	- `QUESTION_DOT` for `?.`
	- `QUESTION_LBRACKET` for `?[`
- Disambiguation order when seeing `?`:
	1) `??` (nullish coalescing)
	2) `?.` (optional member)
	3) `?[` (optional index)
	4) `?` (ternary)
- Require no whitespace between `?` and `.`/`[`.

Parsing
- Treat `?.` and `?[` as postfix operators with the same precedence and associativity as `.` and `[`.
- Allow chaining and mixing: `a?.b[c]?.d`, `a[b]?.c[d]`, etc.
- AST shape options:
	- Extend existing member/index nodes with `Optional: bool`
	- Or introduce `OptionalMemberExpression` and `OptionalIndexExpression` nodes

Compilation/VM
- Prefer explicit safe-get opcodes to keep semantics local and avoid runtime branching in general ops:
	- `SAFE_GET_PROP base, key` → push `null` if `base` is nullish; else same as `GET_PROP`
	- `SAFE_GET_INDEX base, idx` → push `null` if `base` is nullish; else same as `GET_INDEX`
- Ensure the index expression is compiled after the base only when needed (avoid evaluation if short-circuit occurs). In a bytecode VM, this is typically handled by emitting code to evaluate the base, then a conditional jump to skip evaluating the index expression and push `null` if base is nullish.

Error handling
- Optional access should not raise an error on nullish base; it yields `null`.
- Errors from normal access rules still apply when the base is non-nullish (e.g., indexing non-indexable types), unless your language specifies a softer behavior there as well.

## Migration and consistency notes

- Undefined identifiers are documented to resolve to `null`. Ensure identifier evaluation and error messaging are consistent with that rule. In contexts where an identifier currently raises, consider harmonizing with the documented behavior or at least ensuring it doesn’t raise when part of an optional chain.
- Keep `?` for ternary and `??` for nullish coalescing as-is. The `?.` operator does not conflict if tokenized with priority over `?`.

## Testing guide

Happy path
- `null?.x` → `null`
- `{x: 1}?.x` → `1`
- `[1,2]?.[0]` → `1`
- `a?.b?.c` with `a = null` → `null`

Chaining after arbitrary expressions
- `(a ?? b)?.c`
- `(cond ? a : b)?.[i]`
- `(f(x))?.prop`
- `(data |: $last.items)?.[0]`

Index side effects
- `(null)?.[boom()]` → does not call `boom`
- `(obj)?.[computeIndex()]` → calls `computeIndex()` only if `obj` is non-nullish

Mixing
- `a?.b.c` → `(a?.b).c`
- `a.b?.c` → optional only at the final step
- `a?.[i]?.j?.[k]` → consistent short-circuiting

With ??
- `user?.nickname ?? "Anonymous"`
- `(source |: $last.items)?.[0]?.value ?? "N/A"`

Errors propagate only when base is non-nullish
- `(42)?.x` → error or `null` according to normal rules for `42.x`; choose and document
- `(42)?.[0]` → same rationale as above

## Affected areas (codebase)

Parser and tokenizer
- `parser/tokenizer.go`: recognize `?.` and `?[` tokens
- `parser/constants/tokens.go`: define token kinds
- `parser/constants/operators.go`: add operator info (precedence/associativity)
- `parser/parser.go`: parse postfix optional access in the same place as `.` and `[]`

AST
- `ast/dot_expressions.go` and related property/index nodes: add `Optional bool` or new node types
- `ast/identifiers.go`: ensure identifier resolution of missing names aligns with nullish semantics

Compiler/VM
- `compiler/compiler.go`: emit SAFE_GET opcodes (or conditional jumps) for optional steps
- `compiler/bytecode.go`: add new opcodes if needed
- `vm/vm_handlers.go`: implement handlers for SAFE_GET_PROP and SAFE_GET_INDEX
- `operators/dot.go` (or equivalent): ensure underlying access helpers are reused

Tests and docs
- `parser/tests/...`: parsing and precedence cases
- `compiler/tests/...`: bytecode layout for short-circuiting
- `vm/vm_test.go`: runtime semantics and short-circuit behavior
- Docs: this page and cross-links from Operators overview and Pipes overview

## Quick comparison: ?., ?[ ] vs ??

- Use `?.` / `?[ ]` to read safely when the base might be missing.
- Use `??` to supply a fallback value for missing results.
- Combine them for the most expressive, predictable behavior in real data.



