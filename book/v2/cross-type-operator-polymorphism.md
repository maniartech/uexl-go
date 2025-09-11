# V2: Cross-type operator polymorphism (proposal)

While the current VM is strict about cross-type operations, for v2 we propose a small, predictable set of operator overloads to improve ergonomics without compromising type clarity. These do not change numeric IEEE-754 semantics; they only define behavior when one operand is a string or array. All other cross-type cases remain errors unless explicitly listed here.

## Design goals
- Zero implicit stringification or boolean-to-number coercions in `+`, except for string concatenation (string + string) already supported.
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

### 2) Array concatenation, append, prepend
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

### 3) Array repetition
- `array * n` and `n * array` repeat/concatenate the array `n` times.
- Constraints on `n` mirror string repetition:
  - Must be a finite integer value, ≥ 0; NaN/±Inf → error.
- Examples:
  - `[1, 2] * 3` → `[1, 2, 1, 2, 1, 2]`
  - `3 * ["a"]` → `["a", "a", "a"]`
  - `[1] * 0` → `[]`
  - `[1] * 2.5` → error

### 4) Non-goals for now (remain errors)
- `string + number` stays an error (no implicit number-to-string in binary `+`). Use functions or pipes to convert explicitly.
- Object/Map merging via `+` is not defined here.
- Arithmetic with booleans (e.g., `true + 1`) remains an error.

## Interactions with NaN/±Inf
- Repetition counts: NaN or ±Inf → error. Finite non-integer → error.
- Array append/prepend/concat: if the element is a non-finite number, it is included as-is (e.g., `[1] + (+Inf)` → `[1, +Inf]`).
- These cross-type rules do not alter numeric arithmetic, comparison, or truthiness rules elsewhere in this document.
