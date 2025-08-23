# Handling Nullish Coalescing (??) in UExL: Parser, Compiler, and Runtime

This note documents how UExL implements the desired semantics for expressions like:

- `x.y.a.b ?? c`

Goal: Only the final member/index access on the left of `??` should be softened (missing → null), while all preceding accesses remain strict. As a result:

- If `x` has no `y`, or `y` has no `a` → error (not guarded)
- If `a` has no `b` (or `a.b` is null) → treat LHS as null → evaluate `c`

No precedence change is needed. The parse remains `(x.y.a.b) ?? c`.

## Contract (what users get)

- `obj.prop ?? def`
  - If `obj` is missing or not an object → error
  - If `prop` is missing in `obj` → use `def`
  - If `obj.prop` is null → use `def`
  - Else → return `obj.prop`
- `arr[i] ?? def`
  - If `arr` missing or non-indexable → error
  - If `i` out of bounds → use `def`
  - If `arr[i]` is null → use `def`
  - Else → return `arr[i]`
- Chained example: `x.y.a.b ?? c`
  - `x`, `x.y`, `x.y.a` must exist and be of correct types (strict)
  - Only the final `b` access is softened

## Parser/AST: unchanged

- Keep existing precedence: member/index bind tighter than `??`
- AST for `a.b ?? c` stays a BinaryExpression with operator `??` and left=Member(Index)Expression
- No new tokens, no optional chaining involved

## Compiler changes

Introduce a compile-time option to "soften only the final access" when compiling the left operand of `??`.

- When compiling `(LHS) ?? RHS`:
  1) Compile LHS with `softenFinalAccess=true`
  2) Emit a jump that skips RHS if the LHS result is not null
  3) If null: pop LHS and compile RHS

Pseudo-outline (names adapted to your codebase):

- Add an option struct:
  - `type compileOptions struct { softenFinalAccess bool }`
- Entry for `??` in compiler:
  - Compile left with options `{ softenFinalAccess: true }`
  - Emit `OpJumpIfNotNull` (new) with placeholder target
  - Emit `OpPop`
  - Compile right normally
  - Patch jump target
- Member/index compilation helpers accept the option and apply soft opcodes only at the last step of the access chain.

### Member access

If your AST represents `a.b.c` as nested member expressions, the outermost compile call is naturally the last step. Compile base (object) strictly, then:

- Strict last step: `OpGetProp name`
- Soft last step (when `softenFinalAccess` is true): `OpGetPropSoft name`

### Index access

Similarly for indexing:

- Strict last step: `OpIndex`
- Soft last step: `OpIndexSoft`

### Mixed chains

- `x.y[0].z ?? d` → soften only `.z`
- `x.y[maybeIdx] ?? d` → soften only `[maybeIdx]`
- `x?.y.z ?? d` → `?.` guards only the base being nullish; the final `.z` softening still applies only because it is left of `??`

## VM changes

Add three new instructions and handlers:

1) `OpGetPropSoft keyIdx`
   - Pop base
   - If base is object/map:
     - If key present → push value
     - Else → push null
   - Else → type error (non-object)

2) `OpIndexSoft`
   - Pop index, pop base
   - If base is array/slice/string:
     - If index OOB → push null
     - Else → push element
   - Else → type error (non-indexable)

3) `OpJumpIfNotNull target`
   - Peek top of stack
   - If value is not null → jump to target (keeping value on stack)
   - Else → fall through (RHS will be evaluated)

Note: If your runtime represents null with a dedicated type (e.g., `types.Null{}`), the checks should use that sentinel instead of `nil`.

## Edge cases and rules

- Softening does not swallow type errors:
  - Accessing property on a number or boolean should still error
  - Indexing into a non-indexable value should still error
- Only the final access softens. Preceding steps stay strict and may error.
- RHS is only evaluated when LHS becomes null (after softening).
- This rule composes with logical ops and pipes; short-circuit applies as usual.

## Examples

- `user.name ?? "Anonymous"` → "Anonymous" if `name` missing/null; error if `user` missing
- `data.items[10] ?? 0` → 0 if OOB; error if `items` missing or not an array
- `x.y.a.b ?? c` → `c` if `b` missing/null; error if `x`, `y`, or `a` missing
- `order.customer?.address.city ?? "N/A"` → optional chaining guards `customer` only; final `.city` still softened by `??`

## Testing checklist

- Missing final property → coalesces to default
- Null final property → coalesces to default
- Missing earlier link → error
- Type mismatch on base → error
- Index OOB on final link → coalesces to default
- Index OOB on earlier link → error

## Why this design fits UExL

- Keeps explicit, strict-by-default semantics
- Minimizes surprises while making real-world data access concise
- Does not alter parsing or precedence; changes are localized to compiler+VM

