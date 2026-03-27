# V2 Feature Status

This document tracks the implementation status of all features planned or mentioned in the v2 design docs and `pending-things.md`. Last updated: 2026-03-27.

## Legend

| Symbol | Meaning |
|--------|---------|
| ✅ Done | Fully implemented and tested |
| ⚠️ Partial | Workaround exists; full syntactic/API form not yet implemented |
| ❌ Not Started | Designed/documented but no implementation |
| 🚫 Deferred | Rejected or out of scope by design decision |

---

## Already-Implemented (listed as pending in older docs)

These items appear in `pending-things.md` or earlier planning notes but are **fully implemented**.

| Feature | Notes |
|---------|-------|
| ✅ `??` nullish coalescing | Core language; full book chapter |
| ✅ `?.` optional chaining | Core language; full book chapter |
| ✅ `arr[start:end]` slicing | `OpSlice` in VM + compiler; parser `SliceExpression` |
| ✅ `arr[start:end:step]` slicing with step | Same implementation; step handled in `sliceArray`/`sliceString` |
| ✅ `arr[-1]` negative index access | `vm/indexing.go`: `if intIdx < 0 { intIdx = max + intIdx }` |
| ✅ `arr[-2:]` negative-index slicing | Handled via `adjustSliceIndex` in `vm/slicing.go` |
| ✅ `?[` optional slicing | `OpSlice` carries an `optional` flag |
| ✅ `NaN`, `Inf`, `-Inf` literals | Tokenizer opt-in; full IEEE-754 semantics doc; VM tests |
| ✅ `obj["key"]` dynamic key access | `vm/indexing.go` `executeObjectKey` |
| ✅ `obj["computed" + key]` computed key | Works naturally — any expression evaluates to a key |
| ✅ `obj.prop` dot access | Core — has always been implemented |
| ✅ Function calls (with/without args) | Core |
| ✅ Context variables and identifiers | Core |
| ✅ `flatMap` pipe | Registered in `vm/pipes.go`; `FlatMapPipeHandler` |
| ✅ All 11 core pipes | `map`, `filter`, `reduce`, `find`, `some`, `every`, `sort`, `unique`, `groupBy`, `chunk`, `window` |
| ✅ Unicode: `runeLen`, `runeSubstr` | Registered in `vm/builtins.go` |
| ✅ Unicode: `graphemeLen`, `graphemeSubstr` | Registered in `vm/builtins.go` |
| ✅ Unicode: `runes`, `graphemes`, `bytes` | Registered in `vm/builtins.go` |
| ✅ `join` builtin | Registered in `vm/builtins.go` |
| ✅ `($acc ?? 0) + $item` — safe reduce init | The correct and recommended pattern; `??` preserves valid falsy accumulators (`0`, `""`, `false`) |

---

## Not Started

### Pipe Parameters
> Designed in `pipe-parameters.md`. Requirements finalized 2026-03-27.

| Feature | Notes |
|---------|-------|
| ❌ `\|window(n):` — parametric window size | Requires tokenizer + parser + compiler + VM changes; `OpPipe` gains 4th operand |
| ❌ `\|chunk(n):` — parametric chunk size | Same infrastructure as window |
| ❌ `PipeContext.Args() []any` method | New interface method; `pipeContextImpl.args` field + constant pool entry |

**Finalized decisions:**
- Args are **compile-time literals only** (number, string, bool, null). No expressions or variables.
- Multiple args are supported: `|someHandler(3, "asc", true):`
- `|window:` (no args) defaults to size `2` — fully backward compatible.
- `|reduce(n):` is **not** implemented — `($acc ?? 0) + $item` is the canonical pattern.
- Sentinel `0xFFFF` in the 4th operand of `OpPipe` means "no args provided".
- See `pipe-parameters.md` for full specification and implementation checklist.

---

### Dynamic Pipe Expressions
> Designed in `dynamic-pipe-expressions.md` and `migration-and-design.md`.

| Feature | Notes |
|---------|-------|
| ❌ `RegisterPipeExpression(name, template)` API | No implementation in `uexl.go` or `env.go` |
| ❌ `UnregisterPipeExpression` / `HasPipeExpression` | — |
| ❌ `$PRED` placeholder substitution at compile-time | Requires compiler macro registry + AST splice |
| ❌ Compile-time error for missing/extra predicate | — |
| ❌ Macro composition (macros calling other macros) | — |

> **Note:** Registering a **Go** function as a pipe handler via `WithPipeHandlers` is fully implemented. `RegisterPipeExpression` is a separate feature — it allows authoring new pipes entirely in UExL syntax, with no Go code.

### Dynamic Function Expressions
> Designed in `dynamic-function-expressions.md`.

| Feature | Notes |
|---------|-------|
| ❌ `RegisterFunctionExpression(name, template)` API | No implementation |
| ❌ `UnregisterFunctionExpression` / `HasFunctionExpression` | — |
| ❌ `$1`, `$2`, ... positional argument substitution | Requires compiler-time template expansion |

> **Note:** Registering a **Go** function via `WithFunctions` is fully implemented. `RegisterFunctionExpression` is the pure-UExL-template variant.

### Cross-type Operator Polymorphism
> Designed in `cross-type-operator-polymorphism.md`. Currently the VM requires both operands to be the same type.

| Feature | Notes |
|---------|-------|
| ❌ `"ab" * 3` → string repetition | |
| ❌ `3 * "ab"` → string repetition (commutative) | |
| ❌ `"x:" + 3` → string + number concatenation | |
| ❌ `3 + "x"` → number + string concatenation | |
| ❌ `[1,2] + [3,4]` → array concatenation | |
| ❌ `[1,2] + 10` → array append | |
| ❌ `10 + [1,2]` → array prepend | |
| ❌ `[1,2] * 3` → array repetition | |
| ❌ `["a","b"] * ","` → array join (optional candidate) | |
| ❌ `arr1 \| arr2` → array union (optional candidate) | |
| ❌ `arr1 & arr2` → array intersection (optional candidate) | |
| ❌ `arr1 - arr2` → array difference (optional candidate) | |
| ❌ `obj1 \| obj2` → object shallow merge (optional candidate) | |

### Additional Pending Items (from `pending-things.md`)

| Feature | Notes |
|---------|-------|
| ❌ `arr[0, 1]` — 2D index access | |
| ❌ Array creation with ranges | |
| ❌ `in` operator | |
| ❌ `typeof`, `is`, `as` operators | |
| ❌ `typeof()`, `isNaN()`, `isFinite()`, `isNullish()`, `isTruthy()`, `isFalsy()` builtins | |
| ❌ Date and time values and functions | |
| ❌ Regular expressions | |
| ❌ Raw string literals | |
| ❌ `zip` pipe | |
| ❌ `partition` pipe | |

---

## Deferred / Discontinued

| Feature | Reason |
|---------|--------|
| 🚫 `obj.method()` — object method calls | Rejected by design; parser tests explicitly assert this is a parse error |
| 🚫 `BigInt`, `Symbol` types | Out of scope for a Go-embedded engine; no active discussion |
| 🚫 `\|reduce(0): expr` initializer syntax | Discontinued. Use `(\$acc ?? 0) + \$item` with nullish coalescing — the correct and idiomatic approach. |
