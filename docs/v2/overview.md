# Planned Features

This section documents language features that are designed but not yet implemented. Everything here is subject to change based on implementation experience and community feedback.

For the full implementation status of all features (including those that are done, deferred, or discontinued), see [status.md](status.md).

## Dynamic Pipe Expressions

The centerpiece of the planned work is **Dynamic Pipe Expressions** — authoring new pipe behaviors purely in UExL syntax (pipe macros), without writing any Go code. This is similar in spirit to dynamic function expressions, but tailored for pipeline stages.

### Goals

- Author reusable pipe logic entirely in UExL (no Go code required)
- Keep pipeline syntax consistent (`|name:`) and ergonomic
- Allow predicate expressions to be passed unevaluated and used inside the macro
- Preserve current pipe semantics (`$last`, `$item`, `$index`, `$acc`, etc.)

See [Dynamic Pipe Expressions](dynamic-pipe-expressions.md) for details and examples.

## Dynamic Function Expressions

Similar to dynamic pipe expressions but for callable functions — author new functions in UExL syntax without Go code.

See [Dynamic Function Expressions](dynamic-function-expressions.md) for details.

## Cross-type Operator Polymorphism

A small, predictable set of operator overloads for ergonomic cross-type operations (e.g. `string + number`, `string * count`, `array + array`). Currently a proposal; the semantics are not yet finalized.

See [Cross-type Operator Polymorphism](cross-type-operator-polymorphism.md) for the design.

## Other Pending Items

See [pending-things.md](pending-things.md) for the remaining backlog of smaller features.

---

> **Graduated to main book:** The following features were originally planned under "v2" and are now part of the main language:
> [??](../operators/nullish-coalescing.md), [?.](../operators/null-chaining.md), [Slicing](../operators/slicing.md) (including negative indices and step), [Strings and Unicode](../strings-unicode.md), [Numeric Semantics / NaN / Inf](../numeric-semantics.md), [flatMap pipe](../pipes/types.md), computed object key access, and all 12 built-in pipes.

