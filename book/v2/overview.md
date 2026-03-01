# Planned Features

This section documents language features that are designed but not yet implemented. Everything here is subject to change based on implementation experience and community feedback.

## Dynamic Pipe Expressions

The centerpiece of the planned work is **Dynamic Pipe Expressions** — authoring new pipe behaviors purely in UExL syntax (pipe macros), without writing any Go code. This is similar in spirit to dynamic function expressions, but tailored for pipeline stages.

### Goals

- Author reusable pipe logic entirely in UExL (no Go code required)
- Keep pipeline syntax consistent (`|name:`) and ergonomic
- Allow predicate expressions to be passed unevaluated and used inside the macro
- Preserve current pipe semantics (`$last`, `$item`, `$index`, `$acc`, etc.)

### Highlights

- Register dynamic pipe expressions by name with a UExL pipeline fragment
- Predicate capture via `$PRED` placeholder in templates
- Strict usage rules: macros that require `$PRED` must be called with a predicate
- Full composition: macros can chain multiple stages and even invoke other macros

See [Dynamic Pipe Expressions](dynamic-pipe-expressions.md) for details and examples.

## Dynamic Function Expressions

Similar to dynamic pipe expressions but for callable functions — author new functions in UExL syntax without Go code.

See [Dynamic Function Expressions](dynamic-function-expressions.md) for details.

## Cross-type Operator Polymorphism

A small, predictable set of operator overloads for ergonomic cross-type operations (e.g. `string + number`, `string * count`). Currently a proposal; the semantics are not yet finalized.

See [Cross-type Operator Polymorphism](cross-type-operator-polymorphism.md) for the design.

---

> **Already implemented:** Several features that were originally planned under "v2" are now part of the main language. See the main book for: [??](../operators/nullish-coalescing.md), [?.](../operators/null-chaining.md), [Slicing](../operators/slicing.md), [Strings and Unicode](../strings-unicode.md), [Numeric Semantics](../numeric-semantics.md), and [Performance](../performance.md).

