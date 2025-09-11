# V2: Cross-type operator polymorphism (proposal)

While the current VM is strict about cross-type operations, for v2 we propose a small, predictable set of operator overloads to improve ergonomics without compromising type clarity. These do not change numeric IEEE-754 semantics; they only define behavior when one operand is a string or array. All other cross-type cases remain errors unless explicitly listed here.

## Design goals
- No implicit stringification or boolean-to-number coercions in `+`, except when at least one operand is a string (then `+` is concatenation and numbers are stringified).
- Keep behaviors symmetric where possible and intuitive (inspired by Python for `*`).
- Guard with clear, finite-integer requirements; reject NaN/±Inf in counts.

## Proposed rules

### 1) String repetition
- `string * n` and `n * string` repeat the string `n` times.
- Constraints on `n`:
  - Must be a finite integer value (float with zero fractional part is allowed, e.g., `3.0`).
  - Must be ≥ 0.
  - NaN or ±Inf → error: "repeat count must be a finite integer".
- Examples:
  - `"*" * 3` → `"***"`
  - `3 * "ab"` → `"ababab"`
  - `"x" * 0` → `""`
  - `"x" * NaN` → error
  - `"x" * 2.7` → error (non-integer count)

### 2) String concatenation with numbers
- If either operand is a string and the other is a number (int or float64), `+` performs string concatenation after converting the number to a string.
- Numeric-to-string formatting:
  - Finite numbers: use shortest round‑trip formatting (Go `strconv.FormatFloat(v, 'g', -1, 64)` semantics). This yields `"3"`, `"3.14"`, or scientific notation when appropriate.
  - Integers represented as floats (e.g., `3.0`) format without a trailing `.0` → `"3"`.
  - NaN/±Inf stringify to `"NaN"`, `"+Inf"`, `"-Inf"`.
- Examples:
  - `"x:" + 3` → `"x:3"`
  - `3 + "x"` → `"3x"`
  - `"pi=" + 3.14159` → `"pi=3.14159"`
  - `"n=" + NaN` → `"n=NaN"`
  - `"lim=" + (+Inf)` → `"lim=+Inf"`

### 3) Array concatenation, append, prepend
- If either operand is an array, `+` becomes array concatenation with the non-array (if any) treated as a single-element array.
- Forms:
  - `array + array` → concatenation (left then right).
  - `array + value` → append `value` to a copy of `array`.
  - `value + array` → prepend `value` to a copy of `array`.
- Notes:
  - Values can be of any type (including NaN/±Inf), since they are just elements.
  - No implicit deep copy of nested structures beyond the top-level array object itself (same semantics as current array literal behavior).
- Examples:
  - `[1, 2, 3] + 10` → `[1, 2, 3, 10]`
  - `10 + [1, 2, 3]` → `[10, 1, 2, 3]`
  - `[1, 2] + [3, 4]` → `[1, 2, 3, 4]`
  - `[1] + NaN` → `[1, NaN]`

### 4) Array repetition
- `array * n` and `n * array` repeat/concatenate the array `n` times.
- Constraints on `n` mirror string repetition:
  - Must be a finite integer value, ≥ 0; NaN/±Inf → error.
- Examples:
  - `[1, 2] * 3` → `[1, 2, 1, 2, 1, 2]`
  - `3 * ["a"]` → `["a", "a", "a"]`
  - `[1] * 0` → `[]`
  - `[1] * 2.5` → error

### 5) Non-goals for now (remain errors)
- Object/Map merging via `+` is not defined here.
- Arithmetic with booleans (e.g., `true + 1`) remains an error.
- `string + boolean` remains an error (no implicit boolean-to-string in binary `+`).

## Interactions with NaN/±Inf
- Repetition counts: NaN or ±Inf → error. Finite non-integer → error.
- Array append/prepend/concat: if the element is a non-finite number, it is included as-is (e.g., `[1] + (+Inf)` → `[1, +Inf]`).
- String concatenation: NaN/±Inf stringify exactly as `"NaN"`, `"+Inf"`, `"-Inf"`.
- These cross-type rules do not alter numeric arithmetic, comparison, or truthiness rules elsewhere in this document.

## Operator resolution order (disambiguation)
To avoid ambiguity, binary `+` dispatch resolves as follows:
1. If either operand is a string → string concatenation (with numeric operand stringified as above; booleans still error).
2. Else if either operand is an array → array concat/append/prepend as defined above.
3. Else if both operands are numbers → numeric addition (IEEE‑754 semantics).
4. Otherwise → type error.

Binary `*` dispatch resolves as follows:
1. If one operand is a string and the other is a number → string repetition (count must be a finite integer ≥ 0).
2. Else if one operand is an array and the other is a number → array repetition (count must be a finite integer ≥ 0).
3. Else if both operands are numbers → numeric multiplication (IEEE‑754 semantics).
4. Otherwise → type error.

## Commutativity and operand order
Not all cross-type operations are commutative. We define commutativity per-operator and per-type-pair:

- Commutative:
  - `string * number` ⇄ `number * string` (same result)
  - `array * number` ⇄ `number * array` (same result)
  - `number + number`, `number * number` (numeric)

- Order-sensitive (asymmetric):
  - `string + number` vs `number + string` → both valid, results differ by order (concatenation order).
  - `array + array` → concatenates left then right (order matters).
  - `array + value` (append) vs `value + array` (prepend) → intentionally different results.

This mirrors common practice (e.g., Python’s `"s" * n` and list repetition are commutative; concatenations preserve left-to-right order; JS defines `string + number` with order-sensitive concatenation).

## Implementation notes (VM)
To achieve speed, clarity, and robustness, implement operator dispatch using a commutativity-aware multi-method table:

- Tag operands into a small set: number, string, array, boolean, null, object.
- Maintain an O(1) dispatch table indexed by (operator, leftTag, rightTag) → handler + commutative flag.
- For symmetric cases (e.g., `string * number`), either register both directions or register one and allow the dispatcher to swap operands when marked commutative.
- For asymmetric cases (e.g., `array + value` vs `value + array`), register distinct handlers so order is explicit.
- Keep helpers fast and focused:
  - Numeric detection/normalization (int/float → float64) without reflection.
  - Shortest-round‑trip number→string formatting; NaN/±Inf map to `"NaN"`, `"+Inf"`, `"-Inf"`.
  - Size-guarded string/array repetition using doubling to minimize allocations.
- Guardrails: reject NaN/±Inf and non-integers for repetition counts; optional size caps to prevent runaway allocations; consistent error messages.

This pattern is widely recognized (akin to Julia’s multiple dispatch) and delivers predictable semantics with high performance in Go (no reflect, small constant-time tables).
