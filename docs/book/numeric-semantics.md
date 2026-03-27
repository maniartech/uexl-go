# VM numeric semantics for NaN and Inf (IEEE-754)

This document specifies how the VM evaluates arithmetic, comparisons, logicals, and other operators when operands include IEEE-754 special values: NaN and +Inf/-Inf. We follow widely accepted, industry-standard semantics (IEEE-754 and mainstream languages like Go/JS) while preserving UExL's existing error policy.

Scope: runtime behavior in the VM only. Tokenization/parsing of `NaN` and `Inf` is defined in parser docs.

## Quick principles
- Base type for numbers is IEEE-754 double (Go float64).
- NaN payload is not observable; all NaNs are treated the same in the VM.
- Unless explicitly noted, operations use Go's math semantics for float64 (propagate NaN, compute with +Inf/-Inf) and push the resulting float on the stack.
- Deliberate deviation: division by zero is an error (fail-fast), even though IEEE-754 would yield +Inf/-Inf/NaN in some cases. This matches current VM behavior and keeps consistency with existing error expectations.
- Bitwise/shifts are integer-only operations; with non-finite operands they are errors (see below).

## Truthiness and nullish
- Truthiness: numbers are truthy except 0. Therefore, `NaN` and `+Inf`/`-Inf` are truthy.
- Nullish: only `nil` is nullish in the VM. `NaN`/`+Inf`/`-Inf` are not nullish.

## Arithmetic operators
Let x and y be float64 values; operations push a float64 result unless error is stated.

- Addition `+`, Subtraction `-`, Multiplication `*`:
  - If either operand is NaN → result NaN.
  - `(+Inf) + (+Inf) = +Inf`, `(-Inf) + (-Inf) = -Inf`, `(+Inf) + (-Inf) = NaN` (and same for subtraction by algebraic equivalence).
  - `x * (+Inf)` yields sign-appropriate +Inf/-Inf if x ≠ 0 and finite; `0 * (+Inf/-Inf) = NaN`.

- Division `/`:
  - If right operand is 0 → error "division by zero" (intentional deviation from IEEE-754); no value is pushed.
  - Otherwise, IEEE-754 semantics apply:
    - `(+Inf or -Inf) / finite ≠ 0 = +Inf or -Inf` (sign per operands).
    - `finite / (+Inf or -Inf) = 0` (with sign per operands).
    - `(+Inf or -Inf) / (+Inf or -Inf) = NaN`.
    - Any operand NaN → NaN.

- Modulo `%` (math.Mod):
  - Uses `math.Mod(x, y)` for finite operands. Edge cases:
    - `math.Mod(x, 0)` yields NaN (no runtime error is raised by the VM).
    - If either operand is NaN → NaN.
    - `math.Mod(+Inf/-Inf, y)` and `math.Mod(x, +Inf/-Inf)` follow Go/IEEE behavior (commonly NaN for the first, and `x` for the second when finite). The VM forwards `math.Mod` results.

- Power `**` (math.Pow):
  - Follows `math.Pow(x, y)` semantics.
  - Any operand NaN → NaN (except cases where result is defined by `math.Pow`, e.g., `Pow(NaN, 0) = 1`).
  - `Pow(0, negative) = +Inf`; `Pow(0, positive) = 0`.
  - `Pow(+Inf, positive) = +Inf`; `Pow(+Inf, negative) = 0`; similar with `-Inf` per `math.Pow` rules.

- Unary minus `-x`:
  - If x is NaN → NaN; if x is +Inf → -Inf; if x is -Inf → +Inf; else `-x`.

## Comparison and ordering
Number comparisons use IEEE-754 ordering semantics:
- Equality `==`:
  - `NaN == any` (including `NaN`) → false.
  - `+Inf == +Inf` and `-Inf == -Inf` → true; `+Inf == -Inf` → false.
- Inequality `!=`:
  - `NaN != any` (including `NaN`) → true.
- Relational `<, <=, >, >=`:
  - Any comparison with NaN → false.
  - `-Inf < finite < +Inf` holds for all finite numbers.

The VM currently implements number equality/inequality and greater/greater-equal; less/less-equal follow the same rules where present.

## Logical operators (&&, ||, !)
- `!x` (unary logical NOT): Uses truthiness. With numbers, `!0` → true; `!NaN`, `!+Inf`, `!-Inf`, `!nonzero` → false.
- Binary logical ops (&&, ||) operate on booleans in the VM. If user code produces booleans from numbers (e.g., via `!!x`), NaN/Inf participate via truthiness in that conversion.

## Bitwise and shifts (&, |, ^, <<, >>)
- Defined only for integer operands (the VM currently casts from float64 to int internally).
- For non-finite numbers (NaN, +Inf, -Inf): behavior is an error. The VM should reject bitwise operations when either operand is non-finite, with a message like "bitwise requires finite integers".
  - Rationale: Avoid unsafe or surprising casts from non-finite floats to integers.
- For finite non-integral floats: the VM continues to truncate toward zero via `int(...)` as today (implementation detail unchanged). Future tightening to require integer-valued floats could be considered.

## String operations
- Unchanged. `+` on strings is concatenation and does not mix with numeric NaN/Inf.

## Mixed-type operations (NaN/Inf with non-number)
> For proposed v2 cross-type operator semantics (e.g., string/array with numeric operators), see the dedicated page: [Cross-type operator polymorphism](./cross-type-operator-polymorphism.md).
Currently, the VM does not perform implicit type coercions. Numeric operators require numeric operands; string operators require strings; boolean operators require booleans. When NaN or ±Inf are combined with non-number operands, the result is a type error, not coercion.

- Arithmetic and numeric operators (+, -, *, /, %, **, bitwise, shifts):
  - Both operands must be numbers (float64 or int). If either is a non-number (string, bool, array, object, nil), the VM raises a type error like "expected number, got <type>".
  - Examples:
    - `NaN * "hello"` → error: expected number, got string
    - `(+Inf) | true` → error: expected number, got bool
    - `5 % nil` → error: expected number, got <nil>
  - For bitwise and shift operators, see also the earlier rule that non-finite numbers (NaN/±Inf) are errors even when both operands are numeric.

- String concatenation (+):
  - Only defined for two strings. Mixing with numbers (including NaN/±Inf) is an error; there is no auto-stringification in binary `+`.
  - Examples:
    - `"hello" + NaN` → error: string addition requires string operands
    - `NaN + "hello"` → error: expected number, got string

- Comparisons (==, !=, <, <=, >, >=):
  - No cross-type comparisons. Number-vs-number follows IEEE-754 rules above; string-vs-string and bool-vs-bool are supported separately.
  - Mixing types yields a type error like "number comparison requires float64 operands" or "string comparison requires string operands" depending on the left operand.
  - Examples:
    - `NaN == "1"` → error (no implicit conversion)
    - `+Inf > false` → error (no implicit conversion)

- Logical operators (&&, ||, !):
  - Binary logical operators require booleans. Using numbers directly (including NaN/±Inf) is a type error. Convert explicitly if desired using double-negation.
  - `!x` (unary) uses truthiness, so it accepts any type; see Truthiness above.
  - Examples:
    - `NaN && true` → error (operands must be booleans)
    - `!!NaN && true` → true (since NaN is truthy, `!!NaN` is true)

In short: there is no implicit coercion between numbers, strings, and booleans in binary operators. NaN/±Inf behave like ordinary numbers with respect to type requirements—they do not trigger special conversions.

## Error interoperability and propagation
- Where an operation returns a float64 NaN/Inf, that value is pushed and propagates through subsequent numeric operations per rules above.
- The division by zero rule remains an error (no NaN/Inf produced by `/ 0`).
- Other errors (type mismatches, out-of-range indexing, etc.) are unaffected by NaN/Inf.

## Summary table
- NaN propagation: any arithmetic with NaN → NaN; comparisons with NaN → false except `!=` → true; truthy.
- Infinity arithmetic: follows IEEE-754; opposite infinities in addition/subtraction → NaN; scaling by finite values yields +Inf/-Inf unless multiplied by 0 → NaN.
- Division by zero: ERROR (deliberate deviation).
- Modulo/Power: `math.Mod`/`math.Pow` semantics (NaN where applicable); no special error.
- Bitwise/Shifts: ERROR if either operand is NaN or +Inf/-Inf.
- Logical: `!` uses truthiness; `&&`/`||` expect booleans.

## Implementation notes
- Current VM code already follows math semantics for +, -, *, Pow, Mod; adds an error for `/ 0`.
- To align bitwise with this spec, add non-finiteness checks before casting to int in `executeBinaryArithmeticOperation` when handling bitwise and shifts:
  - If `math.IsNaN(leftValue)` or `math.IsInf(leftValue, 0)` or `math.IsNaN(rightValue)` or `math.IsInf(rightValue, 0)`: return an error.
- Tests to consider:
  - `NaN + 1 → NaN`, `Inf - Inf → NaN`, `1 * Inf → Inf`, `0 * Inf → NaN`.
  - `/ 0` still errors; `1 / Inf → 0`.
  - `Mod(5, 0) → NaN` (no error), `Pow(0, -1) → +Inf`.
  - Comparisons with NaN all false except `!=`.
  - Bitwise with NaN/Inf → error; unary minus flips infinity sign.
