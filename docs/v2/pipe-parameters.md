# Pipe Parameters (v2)

Compile-time literal arguments passed directly to a pipe handler via `|pipeName(arg1, arg2, ...):` syntax.

> Status: **Requirements finalized — not yet implemented.**
> See `status.md` for implementation progress.

---

## Motivation

Several built-in pipe handlers (currently `window` and `chunk`) have a single hardcoded configuration value — the window/chunk size — which defaults to `2`. There is no current way to override this from an expression.

```uexl
arr |window:   # always size 2 — cannot be changed
arr |chunk:    # always size 2 — cannot be changed
```

Pipe parameters solve this by letting the expression author pass literal values directly in the pipe header:

```uexl
arr |window(3):                   # sliding window of size 3
arr |chunk(4):                    # fixed chunks of size 4
arr |window(5) as $w: $w.sum()   # window of 5 with alias
```

---

## Non-Goals

The following are explicitly **out of scope** for this feature:

- **`|reduce(initialValue):`** — The `($acc ?? 0) + $item` idiom covers this cleanly. Adding a parameter would create two equivalent ways to express the same thing. See the relevant entry in `status.md`.
- **Runtime/dynamic arguments** — Args must be known at compile time. Variable references, function calls, and arithmetic expressions are not allowed as pipe args.
- **Custom Go pipe handlers with required args** — Existing handlers registered via `WithPipeHandlers` receive args through the `PipeContext.Args()` method and use them at their own discretion. There is no enforcement mechanism for required vs optional args at the API level.

---

## Syntax

```
|pipeName(arg1, arg2, ...):    predicate
|pipeName(arg1, arg2, ...) as $alias:    predicate
```

Only **compile-time literals** are allowed as arguments:

| Literal type | Example |
|---|---|
| Number (float64) | `|window(3):` |
| String | `|sort("desc"):` |
| Boolean | `|someHandler(true):` |
| Null | `|someHandler(null):` |

Multiple arguments are supported: `|someHandler(3, "asc", true):`

Trailing commas are **not** allowed.

### Zero-args form (backward compat)

Calling a pipe with no args (`|window:`) continues to work exactly as before. Handlers that support optional args must define their own defaults (e.g., `window` defaults to size `2`).

---

## Parsing Rules

1. The tokenizer recognizes a pipe name token when it sees `|` followed by letters, and the sequence is terminated by either `:` or `(`.
   - `|window:` → `TokenPipe("window")` + `:` consumed (existing behavior)
   - `|window(3):` → `TokenPipe("window")` + `(` left unconsumed for the parser

2. The parser's `processPipeSegment` function:
   - If the next token after the pipe name is `(`, switches into arg-parsing mode.
   - Consumes `(`, reads zero or more comma-separated literals, consumes `)`, then expects and consumes `:`.
   - If any token between `(` and `)` is not a literal, it is a **compile-time parse error**.
   - If `)` is missing, it is a **parse error**.
   - If `:` is missing after `)`, it is a **parse error**.

3. Parsed args are stored as `[]any` on the `PipeExpression` AST node (`Args` field). Empty/nil slice means no args were provided.

---

## AST

The `PipeExpression` struct gains an `Args []any` field:

```go
type PipeExpression struct {
    Expression Expression
    PipeType   string
    Alias      string
    Args       []any  // nil = no args; compile-time literals only
    Index      int
    Line       int
    Column     int
}
```

No other AST node is affected.

---

## Bytecode

### OpPipe operand layout

`OpPipe` gains a 4th 2-byte operand — `argsIdx`:

| Operand | Width | Description |
|---|---|---|
| `pipeTypeIdx` | 2 bytes | Index into constants: pipe type string |
| `aliasIdx` | 2 bytes | Index into constants: alias string (empty string = no alias) |
| `blockIdx` | 2 bytes | Index into constants: `InstructionBlock` for predicate |
| `argsIdx` | 2 bytes | Index into constants: `[]any` args slice, **or `0xFFFF` = no args** |

Total instruction size: `1 (opcode) + 8 (operands) = 9 bytes` (was 7 bytes).

### Sentinel value

`0xFFFF` is reserved as the "no args" sentinel. This means a valid constants index can never be `65535`. In practice this is not a constraint — a constants pool of 65535+ entries would require extremely pathological expressions.

### Compiler behavior

```
if len(pipeExpr.Args) == 0 {
    emit OpPipe(pipeTypeIdx, aliasIdx, blockIdx, 0xFFFF)
} else {
    argsIdx = addConstant(pipeExpr.Args)  // []any
    emit OpPipe(pipeTypeIdx, aliasIdx, blockIdx, argsIdx)
}
```

---

## VM / PipeContext

### PipeContext interface

Add `Args() []any` as a new method:

```go
type PipeContext interface {
    EvalItem(item any, index int) (any, error)
    EvalWith(scopeVars map[string]any) (any, error)
    Alias() string
    Args() []any          // NEW — nil if no args were provided
    Context() context.Context
}
```

### pipeContextImpl struct

```go
type pipeContextImpl struct {
    vm    *VM
    block *compiler.InstructionBlock
    alias string
    args  []any           // NEW — nil if argsIdx == 0xFFFF
    frame *Frame
}

func (p *pipeContextImpl) Args() []any { return p.args }
```

### OpPipe handler (vm.go)

```go
case code.OpPipe:
    pipeTypeIdx := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
    aliasIdx    := code.ReadUint16(frame.instructions[frame.ip+3 : frame.ip+5])
    blockIdx    := code.ReadUint16(frame.instructions[frame.ip+5 : frame.ip+7])
    argsIdx     := code.ReadUint16(frame.instructions[frame.ip+7 : frame.ip+9])
    frame.ip += 9  // was 7

    var pipeArgs []any
    if argsIdx != 0xFFFF {
        pipeArgs, _ = vm.constants[argsIdx].([]any)
    }
    // ... build pipeContextImpl with args: pipeArgs
```

---

## Built-in Handler Changes

### WindowPipeHandler

Current behavior: `windowSize` is hardcoded to `2`.

New behavior:

```go
windowSize := 2
if args := ctx.Args(); len(args) > 0 {
    if n, ok := args[0].(float64); ok && n >= 2 {
        windowSize = int(n)
    }
    // non-float64 or n < 2: silently use default (2)
}
```

### ChunkPipeHandler

Same pattern as `WindowPipeHandler` — reads `args[0]` as `float64` chunk size, defaults to `2` if not provided or invalid.

### All other built-in handlers

No changes. They simply ignore `ctx.Args()`.

---

## Error Handling

| Scenario | Error |
|---|---|
| Non-literal token inside `(...)` | Parse error: `"pipe arguments must be compile-time literals"` |
| Missing `)` | Parse error: `"expected ')' after pipe arguments"` |
| Missing `:` after `)` | Parse error: `"expected ':' after pipe arguments"` |
| `argsIdx != 0xFFFF` but constant is not `[]any` | VM: silently treat as no args (shouldn't happen in valid bytecode) |
| `args[0]` for window/chunk is not a number or is `< 2` | Handler: silently use default `2` |

---

## Examples

```uexl
# default window (size 2 — backward compat)
[1,2,3,4,5] |window: $window

# explicit window size
[1,2,3,4,5] |window(3): $window

# explicit chunk size
[1,2,3,4,5] |chunk(4): $chunk

# window with alias
data |window(5) as $w: avg($w)

# hypothetical multi-arg handler (not a built-in; for extensibility)
data |myPipe(10, "desc", true): $item
```

---

## Performance Requirements

UExL currently leads benchmarks against expr and cel-go across every measured scenario (see README). This feature must not regress any of those numbers. The architecture below is designed so that **expressions without pipe args pay exactly zero additional cost at eval time**.

### Zero-cost no-args path (common case)

Expressions that never use pipe parameters — which is the entire existing corpus — must pay no additional cost:

- The `0xFFFF` sentinel means the VM skips the constants lookup entirely with a single `uint16` comparison (`argsIdx != 0xFFFF`). This branch is always predictable because all existing pipes have `argsIdx = 0xFFFF`.
- Reading the 4th 2-byte operand from the instruction stream costs nothing extra: it is fetched from the same cache line as the other 3 operands. Modern CPUs prefetch the full 9-byte instruction as a single cache operation.
- `pipeContextImpl.args` is `nil` for all no-args calls. `ctx.Args()` returns `nil` in a single instruction. Handlers that don't use args never touch this field.

### Compile-time allocation only (no eval-time allocs)

Args are stored once in the constants pool during `Compile()`. At eval time:

- The VM reads a pointer from `vm.constants[argsIdx]` — a slice header, no allocation.
- `pipeContextImpl.args` is set to that pointer — a single pointer assignment, no copy, no allocation.
- `ctx.Args()` returns the slice reference directly — no allocation.

**All args-related allocation is confined to `Compile()`, never to `Eval()`.**

### Tokenizer: no new allocations

The `(` check in `readPipeOrBitwiseOr()` is one additional comparison in an existing byte scan loop. No new buffers, no additional token objects for the common `|pipe:` path.

### Parser: lazy allocation

The `args []any` slice in `processPipeSegment` starts as `nil`. It is only constructed when `TokenLeftParen` is the current token. Expressions without args pay zero parse-time allocation for this feature.

### Window/chunk handler overhead

`WindowPipeHandler` and `ChunkPipeHandler` add:
1. One `len(args) > 0` check — a register comparison.
2. One `float64` type assertion — only when args are present.

Both handlers are already O(n) in the input length. This overhead is unmeasurable relative to the loop cost.

### pipeContextImpl struct size

`args []any` adds a 24-byte slice header (pointer + len + cap) to `pipeContextImpl`. Verify this does not push the struct into a higher Go size class or cause pooled copies to exceed 128 bytes. If it does, the `args` field can be changed to `args unsafe.Pointer` with length stored separately, but that optimization should only be applied if benchmarks show a regression.

### Benchmark gate (required before merge)

Run the full benchmark suite before and after implementation. No existing benchmark may regress beyond measurement noise (~5%):

```bash
# Baseline (before changes)
go test -bench=. -benchmem -benchtime=20s ./... > before.txt

# After implementation
go test -bench=. -benchmem -benchtime=20s ./... > after.txt

# Compare
benchstat before.txt after.txt
```

Key benchmarks to monitor:
- `BenchmarkVM_Boolean_Current` — must stay ≤ 223 ns/op, 0 allocs
- `BenchmarkVM_String_Current` — must stay ≤ 119 ns/op, 0 allocs
- `BenchmarkVM_Function_Current` — must stay ≤ 197 ns/op, 2 allocs
- `BenchmarkVM_Pipe_Map100_Current` — must stay ≤ 16,425 ns/op
- Any `|window:` or `|chunk:` benchmark — should show **identical** ns/op vs pre-change (since no args path is unchanged)

New benchmarks to add:

```go
// Verify parametric pipes don't regress standard pipe performance
BenchmarkVM_Pipe_Window3_Parametric   // |window(3): $window over 20-item array
BenchmarkVM_Pipe_Chunk4_Parametric    // |chunk(4): $chunk over 20-item array
BenchmarkVM_Pipe_Window_NoArgs_Compat // |window: $window — must match pre-change baseline
```

---

## Backward Compatibility

All existing expressions continue to work unchanged:
- `|window:` and `|chunk:` without args default to size `2` as before.
- All other pipes ignore `ctx.Args()` and behave identically.
- The bytecode format changes (9-byte `OpPipe` instead of 7-byte), so any serialized/cached bytecode from a previous version is **not** compatible. This is acceptable since UExL does not currently offer a bytecode serialization API.

---

## Implementation Checklist

The following files must be modified, in this order:

1. `code/code.go` — `OpPipe` 4th operand (`[]int{2,2,2}` → `[]int{2,2,2,2}`)
2. `parser/types.go` — add `Args []any` to `PipeExpression`
3. `parser/tokenizer.go` — `readPipeOrBitwiseOr()`: stop at `(` in addition to `:`
4. `parser/parser.go` — `processPipeSegment()`: parse literal args; thread `[][]any` parallel slice; update `ProgramNode` construction
5. `compiler/compiler.go` — emit 4th operand (`argsIdx` or `0xFFFF`)
6. `vm/vm_utils.go` — add `Args() []any` to `PipeContext` interface; add `args []any` to `pipeContextImpl`; implement method
7. `vm/vm.go` — read 4th operand in `case code.OpPipe:`; set `pipeContextImpl.args`; change `ip += 7` to `ip += 9`
8. `vm/pipes.go` — update `WindowPipeHandler` and `ChunkPipeHandler` to read `ctx.Args()`

Performance gate (required before merge):
- Run `benchstat before.txt after.txt` as described in the Performance Requirements section
- All existing benchmarks must be within ~5% of baseline
- Add `BenchmarkVM_Pipe_Window3_Parametric`, `BenchmarkVM_Pipe_Chunk4_Parametric`, `BenchmarkVM_Pipe_Window_NoArgs_Compat`
- Confirm 0 allocs/op is preserved for boolean and string benchmarks

---

## Test Coverage

Test coverage must be comprehensive — every code path introduced by this feature must be exercised, including all error branches, boundary values, and edge cases. The goal is A+ grade quality with no dead untested code.

### Tokenizer tests (`parser/tests/`)

Test file: `pipe_params_tokenizer_test.go` (package `parser_test`)

All tests call `parser.NewParser(input).Parse()` and inspect the resulting AST or error.

**Happy path — token recognition:**

| Input | Expected token sequence |
|---|---|
| `arr \|window(3): $window` | Pipe token `"window"`, args `float64(3)` |
| `arr \|chunk(4): $chunk` | Pipe token `"chunk"`, args `float64(4)` |
| `arr \|window("desc"): $item` | Pipe token `"window"`, args `string("desc")` |
| `arr \|foo(true): $item` | Pipe token `"foo"`, args `bool(true)` |
| `arr \|foo(null): $item` | Pipe token `"foo"`, args `nil` |
| `arr \|foo(3, "asc", true): $item` | Pipe token `"foo"`, args `[3.0, "asc", true]` |
| `arr \|window: $window` | Pipe token `"window"`, no args (backward compat) |
| `arr \|map: $item * 2` | Pipe token `"map"`, no args (all other pipes unchanged) |

**Tokenizer does NOT swallow `(` — parser sees it directly:**

Verify that after a pipe token with `(` following, the parser's next current token is `TokenLeftParen`. This is the contract between tokenizer and parser.

**Pipe name followed by space then `(` — must NOT be treated as args:**

`arr |map (x): $item` — the space breaks the `(` from the pipe name, so this should NOT parse as pipe args (it should produce a parse error or be treated as bitwise OR based on existing grammar rules). Document the exact expected behavior.

---

### Parser tests (`parser/tests/`)

Test file: `pipe_params_parser_test.go` (package `parser_test`)

**Happy path — AST structure:**

```go
// Single number arg
{"[1,2,3] |window(3): $window", PipeExpression{PipeType: "window", Args: []any{float64(3)}}},

// Single string arg
{"arr |sort(\"desc\"): $item", PipeExpression{PipeType: "sort", Args: []any{"desc"}}},

// Single bool arg
{"arr |myPipe(true): $item", PipeExpression{PipeType: "myPipe", Args: []any{true}}},

// Null arg
{"arr |myPipe(null): $item", PipeExpression{PipeType: "myPipe", Args: []any{nil}}},

// Multiple args
{"arr |myPipe(3, \"asc\", true): $item", PipeExpression{PipeType: "myPipe", Args: []any{float64(3), "asc", true}}},

// Empty parens — zero args provided explicitly: |pipe(): ...
// Decision needed: treat as nil args (same as no parens) or parse error.
// Recommended: treat as nil args (args = nil) for simplicity.
{"arr |window(): $window", PipeExpression{PipeType: "window", Args: nil}},

// No args — backward compat
{"[1,2,3] |window: $window", PipeExpression{PipeType: "window", Args: nil}},

// Args combined with alias
{"arr |window(3) as $w: $w", PipeExpression{PipeType: "window", Args: []any{float64(3)}, Alias: "$w"}},

// Negative number arg
{"arr |window(-3): $window", PipeExpression{PipeType: "window", Args: []any{float64(-3)}}},

// Float arg
{"arr |myPipe(3.14): $item", PipeExpression{PipeType: "myPipe", Args: []any{float64(3.14)}}},

// Chained pipes — only one has args
{"arr |filter: $item > 0 |window(3): $window", two PipeExpressions; second has Args: []any{float64(3)}},

// Both chained pipes have args
{"arr |window(3): $window |chunk(2): $chunk", two PipeExpressions; Args verified on each},
```

**Parse errors:**

| Input | Expected error code | Reason |
|---|---|---|
| `arr \|window($x): $window` | `ErrInvalidArgument` | Variable reference not allowed |
| `arr \|window(1+2): $window` | `ErrInvalidArgument` | Expression not allowed |
| `arr \|window(len(x)): $window` | `ErrInvalidArgument` | Function call not allowed |
| `arr \|window(3: $window` | `ErrUnclosedFunction` | Missing `)` |
| `arr \|window(3) $window` | `ErrExpectedToken` | Missing `:` after `)` |
| `arr \|window(3,): $window` | `ErrInvalidArgument` | Trailing comma |
| `arr \|window(,3): $window` | `ErrInvalidArgument` | Leading comma |

**Chained pipe parse errors propagate to top-level:**

`arr |map: $item * 2 |window($x): $window` — the second pipe's error must be reported and parsing must not partially succeed.

---

### Compiler tests (`compiler/tests/`)

Test file: `pipe_params_compiler_test.go` (package `compiler_test`)

Use the existing `compilerTestCase` / `testInstructions` / `testConstants` helpers.

**No-args path — 4th operand is `0xFFFF`:**

```go
// |map: with no args — 4th operand must be 0xFFFF
{"[1] |map: $item", expectedInstructions: []code.Instructions{
    // ... constant loads ...
    code.Make(code.OpPipe, pipeTypeIdx, aliasIdx, blockIdx, 0xFFFF),
}},

// |filter: with no args
{"[1,2] |filter: $item > 1", 4th operand == 0xFFFF},

// |window: with no args (backward compat)
{"[1,2,3] |window: $window", 4th operand == 0xFFFF},
```

**Args path — 4th operand is constants index:**

```go
// |window(3): — args slice stored in constants, index != 0xFFFF
{"[1,2,3] |window(3): $window",
    // constants pool must contain []any{float64(3)} at argsIdx
    // 4th operand of OpPipe must == argsIdx
},

// |chunk(4):
{"[1,2,3,4,5] |chunk(4): $chunk",
    // constants pool must contain []any{float64(4)} at argsIdx
},

// Multi-arg
{"arr |myPipe(3, \"desc\"): $item",
    // constants pool must contain []any{float64(3), "desc"} at argsIdx
},
```

**Chained pipes — each gets independent argsIdx:**

```go
{"[1,2,3,4,5] |window(3): $window |chunk(2): $chunk",
    // two OpPipe instructions
    // first:  4th operand → []any{float64(3)}
    // second: 4th operand → []any{float64(2)}
},
```

**No double-allocation — same args slice not duplicated in constants:**

(Informative test — verify the constants pool length is as expected, not inflated with duplicates from sharing the same literal value. May be deferred to a follow-up.)

---

### VM white-box tests (`vm/pipes_wb_test.go`, package `vm`)

These test the handler functions directly, bypassing the full pipeline. They verify that `pipeContextImpl.args` is correctly read by `WindowPipeHandler` and `ChunkPipeHandler`.

**`WindowPipeHandler` — args control:**

```go
// Args nil → default size 2
pctx := &pipeContextImpl{vm: machine, block: someBlock, args: nil}
res, err := WindowPipeHandler(pctx, []any{1.0, 2.0, 3.0, 4.0, 5.0})
// → [[1,2],[2,3],[3,4],[4,5]]

// args[0] = 3 → size 3
pctx.args = []any{float64(3)}
res, err = WindowPipeHandler(pctx, []any{1.0, 2.0, 3.0, 4.0, 5.0})
// → [[1,2,3],[2,3,4],[3,4,5]]

// args[0] = 5 → size 5 (exact fit)
// args[0] = 6 → size 6 > len(input), results in empty slice

// args[0] = 1 → invalid (< 2), silently fall back to size 2
// args[0] = 0 → invalid, fall back to size 2
// args[0] = -1 → invalid (negative), fall back to size 2
// args[0] = "3" (string, not float64) → invalid type, fall back to size 2
// args[0] = true → invalid type, fall back to size 2
// args is empty slice ([]any{}) → fall back to size 2
// args[0] = 2.0 → size 2 (explicit, same as default)
// args[0] = math.NaN() → fall back to size 2
// args[0] = math.Inf(1) → fall back to size 2 (very large window makes no sense)
```

**`ChunkPipeHandler` — args control (same pattern):**

```go
// args nil → default size 2
// args[0] = 4 → size 4
// args[0] = 1 → invalid, fall back to size 2
// args[0] = 0 → invalid, fall back to size 2
// args[0] = -5 → invalid, fall back to size 2
// args[0] = "4" → invalid type, fall back to size 2
// args[0] = 3.7 (non-integer float) → truncate to int(3), use size 3
// args[0] = len(input) → single chunk containing all elements
// args[0] = len(input)+1 → single chunk containing all elements (no panic)
```

**`pipeContextImpl.Args()` method:**

```go
// nil args
pctx := &pipeContextImpl{args: nil}
assert.Nil(t, pctx.Args())

// non-nil args
pctx.args = []any{float64(3), "desc"}
assert.Equal(t, []any{float64(3), "desc"}, pctx.Args())

// empty slice
pctx.args = []any{}
assert.Equal(t, []any{}, pctx.Args())
assert.NotNil(t, pctx.Args())  // empty != nil
```

**`PipeContext` interface satisfaction:**

Verify `*pipeContextImpl` still satisfies `PipeContext` after adding `Args()`:

```go
var _ PipeContext = (*pipeContextImpl)(nil)  // compile-time check in _wb_test.go
```

---

### VM integration tests (`vm/vm_test.go`, package `vm_test`)

Use `runVmTests` / `vmTestCase`.

**`|window(n):` — correct results:**

```go
{"[1,2,3,4,5] |window(3): $window", []any{
    []any{1.0, 2.0, 3.0},
    []any{2.0, 3.0, 4.0},
    []any{3.0, 4.0, 5.0},
}},
{"[1,2,3,4,5] |window(5): $window", []any{
    []any{1.0, 2.0, 3.0, 4.0, 5.0},
}},
// Size larger than array → empty result
{"[1,2,3] |window(10): $window", []any{}},
// Size == 2 explicit → same as no args
{"[1,2,3,4] |window(2): $window", []any{
    []any{1.0, 2.0}, []any{2.0, 3.0}, []any{3.0, 4.0},
}},
```

**`|window:` backward compat (no args → default 2, MUST equal pre-change output):**

```go
{"[1,2,3,4] |window: $window", []any{
    []any{1.0, 2.0}, []any{2.0, 3.0}, []any{3.0, 4.0},
}},
// Single-element array → empty (window can't form)
{"[1] |window: $window", []any{}},
// Empty array → empty
{"[] |window: $window", []any{}},
```

**`|chunk(n):` — correct results:**

```go
{"[1,2,3,4,5] |chunk(3): $chunk", []any{
    []any{1.0, 2.0, 3.0},
    []any{4.0, 5.0},
}},
{"[1,2,3,4,5] |chunk(5): $chunk", []any{
    []any{1.0, 2.0, 3.0, 4.0, 5.0},
}},
// Size larger than array → single chunk
{"[1,2,3] |chunk(10): $chunk", []any{
    []any{1.0, 2.0, 3.0},
}},
// Evenly divisible
{"[1,2,3,4] |chunk(2): $chunk", []any{
    []any{1.0, 2.0}, []any{3.0, 4.0},
}},
```

**`|chunk:` backward compat:**

```go
{"[1,2,3,4,5] |chunk: $chunk", []any{
    []any{1.0, 2.0}, []any{3.0, 4.0}, []any{5.0},
}},
{"[] |chunk: $chunk", []any{}},
{"[1] |chunk: $chunk", []any{[]any{1.0}}},
```

**Predicate access within window/chunk:**

```go
// Sum each window using predicate
{"[1,2,3,4,5] |window(3): $window[0] + $window[1] + $window[2]", []any{6.0, 9.0, 12.0}},

// Access $index inside chunk predicate
{"[10,20,30,40] |chunk(2): $index", []any{0.0, 1.0}},
```

**Args with alias:**

```go
{"[1,2,3,4,5] |window(3) as $w: $w", []any{
    []any{1.0, 2.0, 3.0},
    []any{2.0, 3.0, 4.0},
    []any{3.0, 4.0, 5.0},
}},
```

**Chained: parametric pipe followed by another pipe:**

```go
// window(3) then map to sum
{"[1,2,3,4,5] |window(3): $window |map: $item[0]", []any{1.0, 2.0, 3.0}},

// chunk(3) then filter non-full chunks
{"[1,2,3,4,5] |chunk(3): $chunk |filter: len($item) == 3", []any{
    []any{1.0, 2.0, 3.0},
}},
```

**Custom pipe handler receives args via `ctx.Args()`:**

```go
// Register a test handler that echoes its args, verify Args() returns them.
// This tests the full PipeContext.Args() plumbing end-to-end.
argsReceived := []any(nil)
testHandler := func(ctx uexl.PipeContext, input any) (any, error) {
    argsReceived = ctx.Args()
    return input, nil
}
env := uexl.NewEnv(uexl.WithPipeHandlers(uexl.PipeHandlers{"echo": testHandler}))
env.Eval(ctx, `[1] |echo(42, "hello"): $item`, nil)
// assert argsReceived == []any{float64(42), "hello"}
```

**All existing (non-parametric) pipe tests continue to pass:**

Run the full existing `vm_test.go` pipe suite without modification. Zero regressions permitted.

---

### Parser error tests (`parser/tests/error_test.go` or new file)

Extend `TestNewErrorSystem` or add `TestPipeParamsParseErrors`:

```go
// All must return a parse error with the specified code.
{`arr |window($x): $window`,    errors.ErrInvalidArgument},
{`arr |window(1+2): $window`,   errors.ErrInvalidArgument},
{`arr |window(len(x)): $window`,errors.ErrInvalidArgument},
{`arr |window(3: $window`,      errors.ErrUnclosedFunction},
{`arr |window(3) $window`,      errors.ErrExpectedToken},
{`arr |window(3,): $window`,    errors.ErrInvalidArgument},
{`arr |window(,3): $window`,    errors.ErrInvalidArgument},
```

---

### Integration / end-to-end tests (`uexl_test.go` or `vm/vm_test.go`)

These test the full `uexl.Eval` / `uexl.Default().Compile()` path:

```go
// Eval shortcut
result, err := uexl.Eval("[1,2,3,4,5] |window(3): $window", nil)
assert.NoError(t, err)
assert.Equal(t, []any{...}, result)

// Compile + re-eval (goroutine-safety smoke test)
ce, err := uexl.Default().Compile("[1,2,3,4,5] |chunk(3): $chunk")
assert.NoError(t, err)
for i := 0; i < 100; i++ {
    go func() {
        res, err := ce.Eval(ctx, nil)
        assert.NoError(t, err)
        _ = res
    }()
}

// Variables() still works (args are literals, not variables)
ce, _ = uexl.Default().Compile("arr |window(3): $window")
assert.Equal(t, []string{"arr"}, ce.Variables())
```

---

### Boundary and edge cases (required for full coverage)

These must be explicitly tested, not assumed:

| Scenario | Expected behavior |
|---|---|
| `|window(2):` — explicit size equal to default | Identical output to `\|window:` |
| `|window(1):` — size below minimum | Silently use default size 2 |
| `|window(0):` — zero size | Silently use default size 2 |
| `|window(-1):` — negative size | Silently use default size 2 |
| `|window(1000000):` — size much larger than array | Empty result, no panic |
| `|chunk(1):` — chunk of one | Each element in its own `[]any{v}` |
| `|window(3):` on empty array | Empty result `[]any{}`, no panic |
| `|chunk(3):` on empty array | Empty result `[]any{}`, no panic |
| `|window(3):` on single-element array | Empty result (window can't fill) |
| `|chunk(3):` on single-element array | `[[v]]` — one partial chunk |
| `|window(3):` on exactly 3 elements | One window `[[a,b,c]]` |
| `|chunk(3):` on exactly 3 elements | One chunk `[[a,b,c]]` |
| Pipe with `null` arg: `\|foo(null):` | `args[0] == nil` in handler |
| Pipe with `false` arg: `\|foo(false):` | `args[0] == false` in handler |
| Pipe with `0` arg: `\|foo(0):` | `args[0] == float64(0)` in handler |
| Pipe with empty string arg: `\|foo(""):` | `args[0] == ""` in handler |
| Multiple pipes, only middle has args | Other pipes still get `args == nil` |
| Non-array input to `\|window(3):` | VM error returned, no panic |
| Non-array input to `\|chunk(3):` | VM error returned, no panic |
