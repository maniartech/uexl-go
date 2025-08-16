# Pipe Compilation and Evaluation (Proposed Robust / KISS Design)

Goal: Treat each pipe stage like a lightweight inline function (lambda) whose body is the predicate expression. Compilation produces a self‑contained instruction block per stage; evaluation executes those blocks with well-defined ephemeral scope variables ($last, $item, etc.). Keeps code simple, predictable, and extensible.

## 1. Design Objectives

- Deterministic, linear compilation.
- No interpretation at runtime (all bytecode).
- Minimal new concepts: reuse existing Compiler scope + VM frame logic.
- Explicit per-stage scope; no hidden mutations.
- Easy to extend with new pipe types.

## 2. New / Adjusted Concepts

1. Instruction Block Constant
   A constant representing a compiled predicate body:
   ```
   type InstructionBlock struct {
     Instructions code.Instructions
     // (No nested constants copy; shares parent constant pool.)
   }
   ```
2. Unified Pipe Stage Emission
   Each stage (after the initial left expression) emits:
   - pipe type string constant (e.g. "map", "filter", "pipe")
   - alias string constant (empty string if none)
   - predicate block constant (may be empty for pass‑through)
   - OpPipe(pipeTypeIdx, aliasIdx, blockIdx)

3. Runtime Scope Stack
   Replace specialized pipeScopes with generic:
   ```
   scopes []map[string]parser.Node
   ```
   Helper methods: pushScope, popScope, setVar, getVar.

4. Standard Emitted Variables
   - Always: $last (incoming value to the stage)
   - map / flatMap / filter / unique / sort / groupBy / find / some / every: $item, $index
   - reduce: $acc (accumulator), $item, $index
   - window: $window, $index
   - chunk: $chunk, $index

## 3. Opcode Changes

OpPipe always has 3 operands `(pipeTypeIdx, aliasIdx, blockIdx)` in this system (no legacy 2‑operand form retained).

VM sequence:

1. Pop input value (result of prior stage or initial expression).
2. Lookup block constant (InstructionBlock).
3. Dispatch by pipe type.

(If backward compatibility needed: version flag or unused alias slot.)

## 4. Compilation Algorithm (Pseudo)

```
compilePipeChain(root):
  compile(leftMostExpression)  // leaves initial value on stack
  prevResultOnStack = true
  for each stage in order:
     pipeTypeIdx = addConstant(StringLiteral(pipeType))
     aliasIdx    = addConstant(StringLiteral(aliasOrEmpty))
     blockIdx    = compilePredicateBlock(stage.predicate)
     emit OpPipe(pipeTypeIdx, aliasIdx, blockIdx)
```

compilePredicateBlock(expr):
```
enterScope()
  compile(expr)         // result left on stack
  blockInstr = currentScope.instructions
exitScope()
return addConstant(&InstructionBlock{Instructions: blockInstr})
```

NOTE: No implicit Pop; predicate’s result consumed by pipe handler logic.

## 5. VM Execution (Per OpPipe)

General steps:
1. Read operands (pipeType, alias, block).
2. Pop inputValue (nil only for malformed chain – treat as error).
3. Set base “incoming” = inputValue.
4. For pipeline semantics:
   - Stage-specific iteration (map/filter/etc.).
   - Each iteration:
     - pushScope()
     - setVar("$last", incoming)
     - set stage vars ($item/$index/$acc/…)
     - executeBlock(blockIdx) → result (top of stack)
     - popScope()
   - Produce stage result (array, accumulator, boolean, or transformed value).
5. Push stage result (becomes input for next OpPipe / final program result).

executeBlock(blockIdx):
```
blk := constants[blockIdx].(*InstructionBlock)
pushFrame(blk.Instructions)
run until frame ends
result := Pop()   // predicate result
return result
```

## 6. Pipe Type Semantics (Core Set)

- pipe / ":" (default):
  - If block empty: pass-through input (still re-emit as $last)
  - Else: single execution with $last set; result replaces input.
- map:
  - Input must be ArrayLiteral; iterate; collect predicate result array.
- filter:
  - Keep element if predicate truthy.
- reduce:
  - First element (or explicit initial syntax future) seeds $acc.
  - Each iteration sets $acc (prev accumulated), $item, $index.
  - Predicate result becomes new $acc.
  - Final $acc pushed.
- find / some / every: short-circuit on condition; still update $last with final result.
- sort / groupBy / unique: collect key via predicate; perform operation after iteration.
- window / chunk: slice segmentation; predicate runs per segment with $window or $chunk.

(Handlers beyond map can be added incrementally following same template.)

## 7. Error Handling

At compile time:
- Missing predicate where required (e.g. map without predicate) → compile error.
- Unsupported pipe type string (optional early validation or defer to runtime).

At runtime:
- Type mismatches (e.g. map on non-array).
- Undefined stage variable lookup -> “undefined identifier: $item” (consistent OpIdentifier error).
- Division / other runtime errors propagate directly.

## 8. Simplicity / KISS Justification

- Single new constant type (InstructionBlock).
- Reuses existing frame loop; no special interpreter path.
- Uniform OpPipe shape regardless of pipe type.
- Stage handlers remain small, pure functions.

## 9. Example Bytecode Sketch

Expression:
```
[1,2] |map: $item * 2 |filter: $item > 3
```
Constants (order illustrative):
0: 1
1: 2
2: "map"
3: ""          (alias)
4: InstructionBlock ( $item * 2 )
5: "filter"
6: ""          (alias)
7: InstructionBlock ( $item > 3 )

Instructions:
OpConstant 0
OpConstant 1
OpArray 2
OpPipe 2 3 4
OpPipe 5 6 7

## 10. Migration Steps

Phase 1:
- Introduce InstructionBlock + new OpPipe signature.
- Update compiler to emit new form.
- Adjust VM OpPipe implementation (map + default only).
Phase 2:
- Add remaining standard pipes (filter, reduce, etc.).
- Replace pipeScopes with generic scopes stack; unify OpIdentifier resolution.
Phase 3:
- Optimize (optional): micro-cache for small blocks or JIT later.

## 11. Minimal Code Changes (Outline)

(Expanded with concrete implementation details.)

### 11.1 File / Package Placement

- Add `compiler/instruction_block.go` for `InstructionBlock` (compiler owns generation).
- Keep execution types (`PipeContext`, handler registry) in `vm` (e.g. `vm/pipes_runtime.go`).
- Expose `InstructionBlock` via the constants pool (already `[]parser.Node`).

### 11.2 Type Definition (compiler/instruction_block.go)

```go
package compiler

type InstructionBlock struct { // implements parser.Node
    Instructions code.Instructions
}
func (b *InstructionBlock) Type() string { return "INSTRUCTION_BLOCK" }
func (b *InstructionBlock) Line() int    { return 0 }
func (b *InstructionBlock) Column() int  { return 0 }
```

(If `parser.Node` requires different methods, adjust accordingly.)

### 11.3 Opcode Update (code/code.go)

Set (single authoritative form):

```go
OpPipe: {"OpPipe", []int{2, 2, 2}}, // pipeTypeIdx, aliasIdx, blockIdx
```

(No legacy variant kept.)

### 11.4 Compiler Scope Helpers (if not present)

```go
func (c *Compiler) enterScope() {
    c.scopes = append(c.scopes, CompilationScope{})
    c.scopeIndex++
}
func (c *Compiler) leaveScope() code.Instructions {
    scope := c.scopes[c.scopeIndex]
    c.scopes = c.scopes[:c.scopeIndex]
    c.scopeIndex--
    return scope.instructions
}
func (c *Compiler) currentInstructions() code.Instructions {
    return c.scopes[c.scopeIndex].instructions
}
```

(Adjust if these already exist.)

### 11.5 Predicate Block Compilation

```go
func (c *Compiler) compilePredicateBlock(expr parser.Expression) int {
    if expr == nil { // pass-through stage
        return c.addConstant(&InstructionBlock{Instructions: nil})
    }
    c.enterScope()
    _ = c.Compile(expr) // errors bubble up externally if desired
    instr := c.leaveScope()
    return c.addConstant(&InstructionBlock{Instructions: instr})
}
```

### 11.6 Pipe Chain Compilation

Assuming AST now: `ProgramNode{ Initial parser.Expression, Stages []*parser.PipeExpression }` where each stage holds `PipeType string`, `Alias string`, `Predicate parser.Expression`.

```go
case *parser.ProgramNode:
    // compile initial
    if err := c.Compile(node.Initial); err != nil { return err }
    for _, st := range node.Stages {
        ptIdx := c.addConstant(&parser.StringLiteral{Value: st.PipeType})
        alIdx := c.addConstant(&parser.StringLiteral{Value: st.Alias}) // empty if none
        blkIdx := c.compilePredicateBlock(st.Predicate)
        c.emit(code.OpPipe, ptIdx, alIdx, blkIdx)
    }
```

### 11.7 VM Constant Accessors

```go
func (vm *VM) constantString(i uint16) string {
    return vm.constants[i].(*parser.StringLiteral).Value
}
func (vm *VM) constantInstructionBlock(i uint16) *compiler.InstructionBlock {
    if vm.constants[i] == nil { return nil }
    return vm.constants[i].(*compiler.InstructionBlock)
}
```

### 11.8 VM OpPipe Execution (New Form)

```go
case code.OpPipe:
    pipeTypeIdx := code.ReadUint16(ins[ip+1:ip+3])
    aliasIdx    := code.ReadUint16(ins[ip+3:ip+5])
    blockIdx    := code.ReadUint16(ins[ip+5:ip+7])

    input := vm.Pop()

    pipeType := vm.constantString(pipeTypeIdx)
    alias := vm.constantString(aliasIdx)
    block := vm.constantInstructionBlock(blockIdx)

    handler, ok := vm.pipeHandlers[pipeType]
    if !ok { return fmt.Errorf("unknown pipe type: %s", pipeType) }

    result, err := handler(input, block, alias, vm)
    if err != nil { return err }

    vm.Push(result)
    ip += 7
```

## 12. Rationale vs Alternatives

| Approach | Rejected Because |
|----------|------------------|
| Inline predicate before OpPipe (current) | Wrong execution order for $item; cannot defer binding cleanly |
| Interpret AST nodes per element | Slower; duplicates evaluator logic |
| Jump table replay of inline code | More IP bookkeeping; harder to maintain |
| Full closure objects now | Overkill without user functions |

Chosen method balances clarity + extensibility.

## 13. Extension Path

Later user-defined pipes can be registered mapping pipeType → handler using same block execution pattern.

---

This plan gives a “function-like” model for pipes with minimal surface area and predictable evaluation behavior.

## 14. User-Provided Pipe Handlers (Extensibility Spec)


The core runtime ships only with a minimal built-in set (e.g. `:`, `map`, `filter`, `reduce`, `find`, `some`, `every`, `unique`, `sort`, `groupBy`, `window`, `chunk`). All additional pipe types are injected by users at init/boot via a registry.

### 14.1 Goals

- Zero coupling: adding / removing a pipe handler does not modify compiler.
- Uniform invocation path (`OpPipe` → lookup → handler).
- Sandboxed predicate execution (no direct frame manipulation in handlers).
- Consistent variable emission (`$last` always present + stage vars).

### 14.2 Public Types (Proposed)

```go
// PipeInstruction is a compiled predicate block (may be nil for pass-through pipes).
type PipeInstruction = InstructionBlock

// PipeContext carries execution-time data for a single OpPipe stage.
type PipeContext struct {
    VM          *VM            // reference for executing predicate blocks
    PipeType    string         // e.g. "map"
    Alias       string         // alias after stage (if any)
    Input       Value          // $last (incoming value)
    Block       *PipeInstruction // compiled predicate (nil or empty allowed)

    // Emit/lookup helpers (scoped to predicate execution frames)
    ExecPredicate func(inject map[string]Value) (Value, error)
}

// PipeHandler executes a pipe stage and returns the stage result.
type PipeHandler func(ctx *PipeContext) (Value, error)
```

Rationale: Handlers never directly manipulate scopes; instead they call `ExecPredicate`, passing the stage variables to inject for that predicate execution. This keeps scope mechanics centralized inside the VM.

### 14.3 Registration API

```go
var pipeRegistry = map[string]PipeHandler{}

func RegisterPipe(name string, handler PipeHandler) {
    pipeRegistry[name] = handler
}

func getPipeHandler(name string) (PipeHandler, bool) {
    h, ok := pipeRegistry[name]
    return h, ok
}
```

Registration is typically performed in an `init()` of an extensions package.

### 14.4 VM Integration (OpPipe Path)

Pseudo:

```go
func (vm *VM) execOpPipe(pipeTypeIdx, aliasIdx, blockIdx int) error {
    pipeType := vm.constantString(pipeTypeIdx)
    alias    := vm.constantString(aliasIdx)
    block    := vm.constantInstructionBlock(blockIdx) // may be nil / empty

    input := vm.pop() // $last

    handler, ok := getPipeHandler(pipeType)
    if !ok { return fmt.Errorf("unknown pipe type: %s", pipeType) }

    ctx := &PipeContext{
        VM: vm,
        PipeType: pipeType,
        Alias: alias,
        Input: input,
        Block: block,
        ExecPredicate: func(inject map[string]Value) (Value, error) {
            return vm.execPredicateBlock(block, inject)
        },
    }

    result, err := handler(ctx)
    if err != nil { return err }

    if alias != "" { vm.bindAlias(alias, result) }

    vm.push(result)
    return nil
}
```

### 14.5 Predicate Execution Helper

```go
func (vm *VM) execPredicateBlock(block *PipeInstruction, inject map[string]Value) (Value, error) {
    if block == nil || len(block.Instructions) == 0 { return vm.nilValue(), nil }

    // push new scope frame for predicate
    vm.pushScope()
    // Always emit $last first (caller must have included Input in inject if needed)
    for k, v := range inject { vm.setVar(k, v) }

    // Run the instruction sequence in an isolated frame
    result, err := vm.runBlock(block.Instructions)

    vm.popScope()
    return result, err
}
```

### 14.6 Built-in Handler Patterns

Example: `map`
```go
func mapHandler(ctx *PipeContext) (Value, error) {
    arr, ok := ctx.Input.(ArrayValue)
    if !ok { return Nil, fmt.Errorf("map expects array input") }
    out := make([]Value, 0, len(arr.Items))
    for idx, item := range arr.Items {
        val, err := ctx.ExecPredicate(map[string]Value{
            "$last": ctx.Input,
            "$item": item,
            "$index": Int(idx),
        })
        if err != nil { return Nil, err }
        out = append(out, val)
    }
    return Array(out), nil
}
```

Example: default (`:`)
```go
func passthroughHandler(ctx *PipeContext) (Value, error) {
    if ctx.Block == nil || len(ctx.Block.Instructions) == 0 {
        return ctx.Input, nil
    }
    val, err := ctx.ExecPredicate(map[string]Value{"$last": ctx.Input})
    if err != nil { return Nil, err }
    return val, nil
}
```

Example: `reduce`
```go
func reduceHandler(ctx *PipeContext) (Value, error) {
    arr, ok := ctx.Input.(ArrayValue)
    if !ok { return Nil, fmt.Errorf("reduce expects array input") }
    if len(arr.Items) == 0 { return Nil, nil }
    acc := arr.Items[0]
    for idx := 1; idx < len(arr.Items); idx++ {
        item := arr.Items[idx]
        next, err := ctx.ExecPredicate(map[string]Value{
            "$last": ctx.Input,
            "$acc": acc,
            "$item": item,
            "$index": Int(idx),
        })
        if err != nil { return Nil, err }
        acc = next
    }
    return acc, nil
}
```

### 14.7 Alias Binding Rule

If alias is non-empty (e.g. `as $temp`), after handler returns successfully the resulting value is inserted into the current top-level pipeline scope under that identifier. Aliases should not overwrite existing locals unless explicitly allowed (simple rule: last wins, or reject duplicates—choose policy early).

### 14.8 Handler Responsibilities

| Responsibility | Handler | VM |
| -------------- | ------- | -- |
| Validate input type | ✓ | |
| Iterate / segment data | ✓ | |
| Construct stage vars map | ✓ | |
| Execute predicate | → via ExecPredicate | ✓ runs block |
| Manage scopes / frames | | ✓ |
| Push final result | | ✓ (after handler returns) |
| Bind alias | | ✓ |

### 14.9 Error Semantics

- Handler returns error → pipeline aborts; error propagates.
- Predicate runtime error bubbles up unchanged.
- Unknown variable inside predicate -> standard identifier resolution error.

### 14.10 Performance Notes

- Reuse allocated maps for inject (optional micro-optimization later; start simple).
- Avoid capturing closures with large outer references.
- Potential future optimization: precompute a fast path for handlers that only need `$item` & `$index`.

### 14.11 Testing Strategy

- Unit test each built-in handler with: empty input, single element, multi element, predicate error, type mismatch.
- Integration test mixed chains including user-registered custom pipe.
- Fuzz test `ExecPredicate` isolation (no leakage of transient vars).

### 14.12 Minimal User-Defined Pipe Example

```go
func init() { RegisterPipe("first", firstHandler) }

func firstHandler(ctx *PipeContext) (Value, error) {
    arr, ok := ctx.Input.(ArrayValue)
    if !ok { return Nil, fmt.Errorf("first expects array input") }
    if len(arr.Items) == 0 { return Nil, nil }
    // Optional predicate: acts as a filter condition
    if ctx.Block == nil || len(ctx.Block.Instructions) == 0 {
        return arr.Items[0], nil
    }
    for idx, item := range arr.Items {
        passed, err := ctx.ExecPredicate(map[string]Value{
            "$last": ctx.Input,
            "$item": item,
            "$index": Int(idx),
        })
        if err != nil { return Nil, err }
        if truthy(passed) { return item, nil }
    }
    return Nil, nil
}
```

### 14.13 KISS Recap

- One opcode (`OpPipe`) for all pipe types.
- One registry.
- One execution helper for predicates.
- Handlers are pure functions: Input + (optional) predicate → Output.

---

This extended section defines the stable extension surface for user pipe handlers while retaining the simple core execution model.