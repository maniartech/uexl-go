# UExL Go - AI Coding Agent Instructions

## Project Overview

**UExL (Universal Expression Language)** is a bytecode-compiled, embedded expression evaluation engine for Go with a unique **three-pass architecture**: Parser → Compiler → VM. The project emphasizes **zero-panic robustness**, **explicit nullish/boolish semantics**, and **pipe-based data transformations**.

## Core Architecture (3-Stage Pipeline)

```
Expression String → Parser (AST) → Compiler (Bytecode) → VM (Execution)
```

### 1. Parser (`parser/`)
- **Entry**: `parser.ParseString(expr)` or `parser.NewParser(input).Parse()`
- Produces AST nodes (`parser.Node` interface) from tokenized input
- **Key files**: `parser.go`, `tokenizer.go`, `types.go`
- **Error handling**: Returns structured `errors.ParserError`, NEVER panics
- **Sub-packages**: `constants/`, `errors/`, `specs/`, `tests/`

### 2. Compiler (`compiler/`)
- **Entry**: `compiler.New()` then `Compile(ast)`
- Transforms AST into `*compiler.ByteCode` (constants + instructions)
- **Key concepts**:
  - `CompilationScope` for nested scopes (pipes, function blocks)
  - `InstructionBlock` for pipe predicates (compiled as bytecode constants)
  - **Backpatching** for short-circuit jumps (`OpJumpIfFalsy`, `OpJumpIfTruthy`)
- **Key files**: `compiler.go`, `compiler_utils.go`, `bytecode.go`
- **Pipe compilation**: See `pipe_compilation_and_evaluation.md` for design

### 3. VM (`vm/`)
- **Entry**: `vm.New(LibContext{...})` then `Run(bytecode, contextVars)`
- Stack-based bytecode interpreter with frame management
- **Key files**: `vm.go`, `vm_handlers.go`, `pipes.go`, `builtins.go`
- **VM components**:
  - Stack (1024 slots), Frame stack (1024 frames)
  - Pipe scope stack for `$item`, `$index`, `$acc`, `$last`, etc.
  - Built-in functions registry (`VMFunctions`)
  - Pipe handlers registry (`PipeHandlers`)

## Critical Developer Workflows

### Running Tests
```bash
# All tests with race detection
go test ./... -race

# Specific package
go test ./parser/tests -v

# Benchmarks (compare against expr/cel-go)
go test -bench=. -benchtime=20s
```

### Using Tasks (VSCode)
- **`go test all`**: Runs `go test ./...`
- **`go test after enabling specials`**: Runs `go test ./... -v`

### Typical Expression Evaluation Flow
```go
// 1. Parse
node, err := parser.ParseString("arr |map: $item * 2")
// 2. Compile
comp := compiler.New()
comp.Compile(node)
bytecode := comp.ByteCode()
// 3. Execute
machine := vm.New(vm.LibContext{
    Functions: vm.Builtins,
    PipeHandlers: vm.DefaultPipeHandlers,
})
result, err := machine.Run(bytecode, map[string]any{"arr": []any{1.0, 2.0}})
```

## Project-Specific Conventions

### 1. Pipe Operator (`|:` or `|type:`) - THE KEY DIFFERENTIATOR
- Pipes transform data through **chained stages** with **ephemeral scope variables**
- **Syntax**: `expr |map: $item * 2 |filter: $item > 5`
- **Built-in pipes**: `map`, `filter`, `reduce`, `find`, `some`, `every`, `unique`, `sort`, `groupBy`, `window`, `chunk`, `:` (default/passthrough)
- **Context vars**: `$item`, `$index`, `$acc`, `$window`, `$chunk`, `$last`
- **Compilation**: Each pipe predicate becomes an `InstructionBlock` constant
- **VM execution**: `OpPipe(pipeTypeIdx, aliasIdx, blockIdx)` → handler lookup → predicate execution in isolated frame
- **Extensibility**: Register custom pipes via `PipeHandlers` map

### 2. Explicit Nullish/Boolish Semantics (NOT JavaScript-like for defaults)
- **Strict access**: `a.b` and `a[i]` throw on missing keys/indices
- **Optional chaining** (`?.`, `?.[`) ONLY guards nullish base, NOT missing members
- **Nullish coalescing** (`??`) falls back ONLY on `null`, preserves `0`, `""`, `false`
- **Safe mode**: `x.a.b ?? c` makes ONLY `b` access safe, NOT earlier links
- See `book/design-philosophy.md` for rationale

### 3. Testing Patterns
- **Parser tests**: Use `parser_test` package, check AST structure + errors
- **Compiler tests**: Verify bytecode instructions + constants (see `compiler/tests/help_test.go`)
- **VM tests**: Use `vmTestCase{input, expected}` + `runVmTests()`
- **White-box tests**: Suffix `_wb_test.go` for internal package tests
- **Benchmark naming**: `Benchmark_<component>_<scenario>` (e.g., `BenchmarkVM_Boolean_Current`)

### 4. Error Handling (Zero-Panic Policy)
- **NEVER use panic** in production code paths
- Return structured errors: `errors.ParserError`, `fmt.Errorf()`, or custom types
- VM propagates errors from opcodes without wrapping (preserve stack traces)
- See `compiler-vm-upgrade.md` for robustness guidelines

### 5. Performance Benchmarks (Context: UExL is ~10x slower than expr/cel-go)
- **Current target**: Sub-100 ns/op for simple boolean expressions
- **Comparison project**: `../golang-expression-evaluation-comparison/`
- **Profiling**: CPU profiles in `*.prof` files (use `go tool pprof`)
- **Optimization docs**: `wip-notes/phase1-implementation-plan.md`, `MILITARY_GRADE_PERFORMANCE_FRAMEWORK.md`

## Key Design Decisions & Gotchas

### Instruction Blocks for Pipes
- Pipe predicates are NOT interpreted AST—they're compiled to bytecode
- `InstructionBlock` is a `parser.Node` implementation holding `code.Instructions`
- Stored in constants pool, executed in isolated VM frames

### Short-Circuit Compilation (Backpatching)
- `&&` and `||` use `OpJumpIfFalsy`/`OpJumpIfTruthy` with placeholder offsets
- Compiler backpatches jump targets after compiling right-hand side
- See `compiler/compiler.go:compileShortCircuitChain()`

### VM Frame Management
- Each pipe predicate execution pushes a new frame (`vm.pushFrame()`)
- Frames are NOT automatically popped on error—handlers must clean up
- Stack pointer (`sp`) and frame index (`framesIdx`) managed manually

### Context Variables vs System Variables
- **Context vars**: User-provided data (`map[string]any` in `Run()`)
- **System vars**: Internal bindings like `$item`, `$index` (managed by VM)
- **Alias vars**: Pipe aliases (`|map as $result:`)

### Unicode Handling (Future: Grapheme Clusters)
- Current: Rune-based indexing/slicing (UTF-32 code points)
- Planned: Grapheme-aware operations (see `wip-notes/graphemes-upgrade-design.md`)
- View functions: `char()`, `utf8()`, `utf16()` for explicit levels

## File Organization Patterns

```
uexl-go/
├── parser/           # Tokenizer + Parser + AST
│   ├── constants/    # Token types, pipe types
│   ├── errors/       # ParserError types
│   ├── specs/        # Design docs (roadmap.md)
│   └── tests/        # Parser + tokenizer tests
├── compiler/         # AST → Bytecode
│   └── tests/        # Bytecode verification tests
├── vm/               # Bytecode → Result
│   ├── builtins.go   # len(), substr(), contains(), etc.
│   ├── pipes.go      # Pipe handlers (map, filter, reduce, ...)
│   └── *_wb_test.go  # White-box VM tests
├── code/             # Opcode definitions
├── book/             # User-facing docs (GitBook structure)
│   ├── v2/           # Future features (slicing, graphemes)
│   └── internals/    # (empty - architecture docs in wip-notes)
├── wip-notes/        # Design docs, upgrade plans
└── performance_benchmark_test.go  # Performance tracking
```

## Integration Points & Dependencies

- **No external dependencies** except `github.com/stretchr/testify` for tests
- **Vendored**: All dependencies in `vendor/` (Go modules enabled)
- **Comparison benchmarks**: Separate workspace at `../golang-expression-evaluation-comparison/`

## When Making Changes

1. **Parser changes**: Update AST types in `parser/types.go`, add tests in `parser/tests/`
2. **New operators**: Add to tokenizer, parser, compiler (emit opcodes), VM (opcode handlers)
3. **New pipes**: Add handler in `vm/pipes.go`, register in `DefaultPipeHandlers`
4. **New builtins**: Add in `vm/builtins.go`, register in `Builtins` map
5. **Performance work**: Run benchmarks BEFORE and AFTER, document in commit message
6. **Breaking changes**: Follow roadmap in `parser/specs/roadmap.md`

## Common Pitfalls

- **Don't mix `parser.Node` and `any` in constants pool** (compiler stores both)
- **Always check `ok` on type assertions** in VM handlers (no panics!)
- **Pipe handlers must clean up scopes/frames on error** (use defer or explicit cleanup)
- **Backpatch jump offsets are 2-byte operands** (use `code.ReadUint16`/`WriteUint16`)
- **VM stack operations use `sp` (stack pointer)** - watch for off-by-one errors

## Quick Reference Links

- **Language design**: `book/design-philosophy.md`, `book/syntax.md`
- **Pipe design**: `compiler/pipe_compilation_and_evaluation.md`
- **Performance plans**: `wip-notes/phase1-implementation-plan.md`
- **Upgrade roadmap**: `compiler-vm-upgrade.md`, `parser/specs/roadmap.md`
- **Progress tracker**: `progress.md`
