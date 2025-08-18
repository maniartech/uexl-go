# Opcodes with Pipe Expressions

This document shows a proposed (KISS) bytecode layout for a pipeline expression and presents it in two complementary formats:

1. Byte offset table (existing style)
2. Linear instruction index format (requested style similar to other VM listings)

Target expression:

```text
[1, 2, 3, 4, 5] |filter: $item > 2 |map: $item * 2 |reduce: $acc + $item
```

## Assumptions (Proposed Design)

- New `OpPipe` uses **3 operands**: `(pipeTypeIdx, aliasIdx, blockIdx)`.
- Predicate bodies are compiled into `InstructionBlock` constants (not inlined in the main stream).
- Constant pool deduplicates identical values (e.g. number `2`, empty alias string).
- Pipe local identifiers (`$item`, `$acc`) compile via `OpIdentifier` referencing a local/pipe-scope index (shown symbolically as `idItem`, `idAcc`).
- No initial accumulator argument for `reduce`; first element is used (future enhancement could support optional seed).

> If the current engine still has a 2‑operand `OpPipe`, the third operand (`blockIdx`) addition would require: (a) updating the opcode definition, (b) incrementing operand width metadata, (c) adjusting the VM dispatch. A fallback layout for 2‑operand form is included later.

## Constant Pool (Example Ordering)

Idx | Kind              | Value / Detail
----|-------------------|-----------------------------
0   | NumberLiteral     | 1
1   | NumberLiteral     | 2
2   | NumberLiteral     | 3
3   | NumberLiteral     | 4
4   | NumberLiteral     | 5
5   | String            | "filter"
6   | String            | "" (empty alias shared)
7   | InstructionBlock  | filter predicate `$item > 2`
8   | String            | "map"
9   | InstructionBlock  | map predicate `$item * 2`
10  | String            | "reduce"
11  | InstructionBlock  | reduce predicate `$acc + $item`

### Predicate Block Internal Instructions

Block 7 (filter: `$item > 2`):

```text
Idx  Opcode        Operands  Notes
0    OpIdentifier  idItem    push $item
1    OpConstant    1         push 2
2    OpGreaterThan           compare
```

Block 9 (map: `$item * 2`):

```text
Idx  Opcode        Operands  Notes
0    OpIdentifier  idItem    push $item
1    OpConstant    1         push 2
2    OpMul                   multiply
```

Block 11 (reduce: `$acc + $item`):

```text
Idx  Opcode        Operands  Notes
0    OpIdentifier  idAcc     push $acc
1    OpIdentifier  idItem    push $item
2    OpAdd                   sum
```

## Main Instruction Stream – Byte Offset View

(Offsets reflect current encoding: 1 byte opcode + 2 bytes per operand for opcodes with operands.)

Addr | Opcode      | Operands (decoded) | Comment / Stack Effect
-----|-------------|--------------------|------------------------
0000 | OpConstant  | 0                  | push 1
0003 | OpConstant  | 1                  | push 2
0006 | OpConstant  | 2                  | push 3
0009 | OpConstant  | 3                  | push 4
0012 | OpConstant  | 4                  | push 5
0015 | OpArray     | 5                  | build array [1..5]
0018 | OpPipe      | 5, 6, 7            | filter stage (type="filter", alias="", blockIdx=7)
0025 | OpPipe      | 8, 6, 9            | map stage (type="map", alias="", blockIdx=9)
0032 | OpPipe      | 10, 6, 11          | reduce stage (type="reduce", alias="", blockIdx=11)
0039 | (End)       |                    | final result on stack

## Main Instruction Stream – Linear Instruction Index Format

(Each row = one instruction irrespective of operand byte width.)

Idx | Opcode     | Operands      | Annotation
----|------------|---------------|-----------
0   | OpConstant | <0>           | 1
1   | OpConstant | <1>           | 2
2   | OpConstant | <2>           | 3
3   | OpConstant | <3>           | 4
4   | OpConstant | <4>           | 5
5   | OpArray    | <5>           | len=5 → [1..5]
6   | OpPipe     | <5,6,7>       | filter (block 7)
7   | OpPipe     | <8,6,9>       | map (block 9)
8   | OpPipe     | <10,6,11>     | reduce (block 11)

Legend:

- Angle brackets `<a,b,c>` show operand constant indices in order.
- An empty alias uses constant index 6 (the shared empty string literal).

## 2‑Operand Fallback (Current Engine Form)

If `OpPipe` currently has only two operands `(pipeTypeIdx, aliasIdx)` and predicates are compiled *inline before* the pipe (legacy style), the layout would differ:

1. Emit predicate code directly before each `OpPipe`.
2. `OpPipe` consumes (and/or peeks) top-of-stack predicate result along with prior stage value.
3. No `InstructionBlock` constants (indices 7,9,11 would disappear; constant indices shift).

This fallback is **not** recommended for long-term extensibility because:

- Makes `$item`, `$acc` binding order delicate.
- Harder to defer / re-run predicate or support short‑circuit semantics cleanly.

## Validation / Sanity Checks

- Only one array construction before first pipe.
- Each pipe instruction appears exactly once per stage.
- Reuse of constant `2` (index 1) across both filter and map predicate blocks.
- All stages share identical alias constant (empty string) – reduces constant pool churn.

## Potential Optimizations (Future)

- Inline small (single opcode) predicate blocks by flagging empty / trivial blocks.
- Cache resolved handler function pointers keyed by pipeTypeIdx to skip registry map lookup.
- Optional constant folding inside predicate blocks prior to sealing `InstructionBlock`.
- Potential single pass to compute needed pipe locals mask to speed identifier resolution.

## Error Scenarios (Runtime)

Condition | Detection Point | Action
--------- | ---------------- | ------
Non-array input to map/filter/reduce | Handler | Raise type error
Empty array reduce (no seed) | Handler | Return Nil (or configurable)
Unknown pipe type | VM dispatch | Error & abort pipeline
Missing predicate where required (e.g. map w/out block) | Compile or handler | Error

## Summary

The 3‑operand `OpPipe` plus `InstructionBlock` constants keeps the main instruction stream flat and defers predicate execution to handler-controlled scopes, enabling clean user-defined pipe extensions while remaining minimal and deterministic.

If/when implementation diverges from these assumptions (e.g. operand widths, alias semantics), update this document to maintain accuracy.