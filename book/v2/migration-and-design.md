# v2 Migration and Design Notes

This document outlines the rationale, compatibility considerations, and the high-level design of v2 features, focusing on Dynamic Pipe Expressions.

## Motivation
- Reduce friction for extending the language by allowing users to author new pipe behavior in UExL itself.
- Promote reuse: common patterns (e.g., filtering, mapping, aggregating) can be named and shared.
- Keep pipelines declarative and readable.

## Compatibility
- v1 expressions continue to work unchanged.
- Dynamic Pipe Expressions are additive; they do not alter existing pipe semantics.
- Macro registration is opt-in. Without registration, names are resolved as built-in or user-implemented pipes as before.

## Semantics
- Macro templates expand at compile-time into one or more standard pipe stages.
- `$PRED` in the template denotes where the call-site predicate (if any) is inlined.
- If `$PRED` is present, the call-site must provide a predicate after the colon.
- If `$PRED` is absent, the macro must be invoked without a predicate.
- Expanded stages obey normal scoping for `$last`, `$item`, `$index`, `$acc`.

## Error Handling
- Helpful compile-time errors should be produced:
  - Unknown pipe: `unknown pipe 'name'` (existing behavior)
  - Predicate required: `pipe 'name' requires a predicate`
  - Predicate not allowed: `pipe 'name' does not accept a predicate`
  - Template parse error: show the underlying parse error location in the expanded template

## Implementation Sketch

1. Registry
   - Maintain a global/ctx registry mapping `name -> template`.
   - Optionally track a boolean `requiresPredicate` by scanning for `$PRED`.

2. Parser/Compiler Hook
   - When encountering a pipe stage `|name: <expr?>`:
     - If `name` is in registry:
       - Validate predicate presence/absence against `requiresPredicate`.
       - Construct the effective template string by substituting `$PRED` with the original predicate source (or AST splice if your compiler prefers AST-level integration).
       - Parse the resulting pipeline fragment (it must start with a `|` stage).
       - Splice the resulting AST stages into the current pipeline at that location.
     - Else fallback to built-in/registered pipe resolution.

3. VM/Runtime
   - No new runtime handlers are required if expansion occurs before codegen.
   - Bytecode will reflect the expanded stages exactly as if they were authored inline.

## Examples of Macro Templates

- `where`: `|filter: $PRED`
- `select`: `|map: $PRED`
- `sumBy`: `|reduce: ($acc || 0) + ($PRED)`
- `concatStr`: `|reduce: ($acc || '') + ($PRED) + ',' |: substr($last, 0, len($last)-1)`

## Limitations & Future Work
- Currently single-placeholder `$PRED`. We may extend to multiple placeholders (e.g., `$P1`, `$P2`) to support multi-argument pipes.
- Optional call-site aliases or parameters (e.g., `$ALIAS`, `$IN`) may be introduced to support additional authoring patterns.
- Namespacing or shadowing rules if macro names overlap with built-ins.

## Migration Tips
- Start by extracting frequently repeated pipeline fragments into macros.
- Keep the same names as your conceptual operations (`where`, `select`), and document their intended use.
- Add tests for macro expansion boundaries (error cases and scoping interactions) as you adopt v2.
