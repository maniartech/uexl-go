# UExL v2 Overview

UExL v2 introduces a set of forward-looking enhancements to improve reusability, composability, and ergonomics without requiring host-language code.

The centerpiece of v2 is Dynamic Pipe Expressions — authoring new pipes purely in UExL syntax (pipe macros), similar in spirit to dynamic function expressions, but tailored for pipeline stages.

Status: Experimental (subject to change). Contributions and feedback are welcome.

## Goals

- Author reusable pipe logic entirely in UExL (no Go code required)
- Keep pipeline syntax consistent (`|name:`) and ergonomic
- Allow predicate expressions to be passed unevaluated and used inside the macro
- Preserve current pipe semantics (`$last`, `$item`, `$index`, `$acc`, etc.)

## Highlights

- Register dynamic pipe expressions by name with a UExL pipeline fragment
- Predicate capture via `$PRED` placeholder in templates
- Strict usage rules: macros that require `$PRED` must be called with a predicate
- Full composition: macros can chain multiple stages and even invoke other macros

See Dynamic Pipe Expressions for details and examples.
