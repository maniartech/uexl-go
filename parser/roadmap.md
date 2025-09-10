# Parser Roadmap (Breaking Changes + Migration Guide)

This roadmap proposes targeted, standards-driven improvements to the `parser` package that will introduce breaking changes in `compiler` and `vm`. It includes clear rationale, concrete examples, and step-by-step mitigation strategies to make the upgrade safe and predictable.

## Why change? (Strong reasons to proceed)

- Industry-standard APIs: Aligns with Go stdlib patterns (explicit `(T, error)` returns, typed values, clear semantics).
- Performance: Lower allocations and faster hot paths established by recent tokenizer improvements; the same rigor applied to AST and parser will compound system-wide.
- Maintainability and Safety: Replace `any` in tokens/AST with narrow types to prevent runtime type assertions and class of bugs.
- Better tooling: Clearer contracts make compilers and VMs simpler to reason about, test, and optimize.
- Future features: A stable, typed AST and options open the door to optimization passes (folding, DCE) and richer editor tooling.

## Summary of planned changes

1) Token typing overhaul (BREAKING)
- Today: `Token.Value any`
- Proposed: strongly typed token payloads via a tagged union struct

```go
// New in parser
// Replaces Token.Value any
// Provides exact type and payload without type assertions

type TokenValueKind uint8
const (
    TVKNone TokenValueKind = iota
    TVKNumber
    TVKString
    TVKBoolean
    TVKNull
    TVKIdentifier
    TVKOperator
)

type TokenValue struct {
    Kind TokenValueKind
    Num  float64
    Str  string
    Bool bool
}

// Token will embed TokenValue
// type Token struct { Type constants.TokenType; Value TokenValue; ... }
```

Impact:
- Compiler and VM locations that assert `token.Value.(float64|string|bool)` will fail to compile.
- Benefit: compile-time guidance and faster paths (no interface{} boxing/unboxing, fewer allocs).

Mitigation:
- Provide helper accessors and conversion shims during transition:

```go
// Temporary helpers (parser/internal/compat)
func AsFloat(t Token) (float64, bool)
func AsString(t Token) (string, bool)
func AsBool(t Token) (bool, bool)
```

2) AST node typing and constructors (BREAKING)
- Today: Many nodes use `any` for fields like MemberAccess.Property.
- Proposed: use a small sum type:

```go
// For properties that can be string or int

type PropertyKind uint8
const (
    PropString PropertyKind = iota
    PropInt
)

type Property struct {
    Kind PropertyKind
    S    string
    I    int
}

// MemberAccess.Property: Property (instead of any)
```

Impact:
- Compiler pattern-matches property kinds; existing code doing type switches will need updates.

Mitigation:
- Add small constructor helpers and matchers:

```go
func PropS(s string) Property { return Property{Kind: PropString, S: s} }
func PropI(i int) Property   { return Property{Kind: PropInt, I: i} }
func (p Property) IsString() bool { return p.Kind == PropString }
func (p Property) IsInt() bool    { return p.Kind == PropInt }
```

3) Parser options instead of boolean state (BREAKING)
- Today: `Parser` tracks `subExpressionActive`, `inParenthesis` as booleans.
- Proposed: `Options` struct to make behavior explicit and testable.

```go
// Examples only; details may evolve

type Options struct {
    // Language feature toggles
    EnableNullish bool
    EnableOptionalChaining bool
    EnablePipes bool

    // Limits & safety
    MaxDepth int // 0 => unlimited
}

func NewParserWithOptions(input string, opt Options) *Parser
```

Impact:
- Parser constructors and call sites change.
- VM and compiler tests that assume all features always-on may need to pass Options or use defaults.

Mitigation:
- Keep `NewParser(input string) *Parser` as a convenience that uses `DefaultOptions()`.
- Introduce `DefaultOptions()` with today’s behavior.

4) Position handling and node interfaces (MINOR BREAKING)
- Normalize AST nodes to implement a small `Pos()` interface for uniform position access.

```go
// New common interface

type Node interface {
    Pos() (line, col int)
}
```

Impact:
- Compiler/VM walkers can be simplified; custom access may need small refactors.

Mitigation:
- Provide default implementations on all nodes and ensure tests cover position propagation.

5) Optional: AST arena/pooling (OPT-IN)
- Introduce an internal arena or sync.Pool for frequently created nodes, guarded by options and benchmarks.
- No public API break if kept internal; can be iterated safely.

## Concrete examples and migrations

### Example A: Numbers and strings (compiler token consumers)

Before:
```go
switch tok.Type {
case TokenNumber:
    v := tok.Value.(float64)
case TokenString:
    s := tok.Value.(string)
}
```
After:
```go
switch tok.Type {
case TokenNumber:
    v := tok.Value.Num // no type assertion
case TokenString:
    s := tok.Value.Str
}
```

Mitigation (intermediate):
```go
v, _ := compat.AsFloat(tok)
s, _ := compat.AsString(tok)
```

### Example B: Member access property handling

Before:
```go
switch v := ma.Property.(type) {
case string:
    // field access
case int:
    // index-like property
}
```
After:
```go
if ma.Property.IsString() {
    field := ma.Property.S
    // field access
} else {
    idx := ma.Property.I
    // index-like property
}
```

### Example C: Parser construction in compiler tests

Before:
```go
p := parser.NewParser(src)
expr, err := p.Parse()
```
After (explicit options):
```go
p := parser.NewParserWithOptions(src, parser.DefaultOptions())
expr, err := p.Parse()
```
Or keep convenience:
```go
p := parser.NewParser(src) // uses DefaultOptions under the hood
expr, err := p.Parse()
```

## Breaking change surface for compiler and vm

- Token consumption paths: replace interface assertions with typed fields.
- AST property handling: update switches to the new `Property` type.
- Parser construction: optionally pass `Options` or rely on `DefaultOptions()`.
- Node interface: if you leverage `Node`/`Pos()`, you can simplify walkers.

Expected effort: mechanical refactors across a limited set of choke points (token handling, member access, parser creation). Tests will guide changes at compile time.

## Step-by-step migration plan

1. Introduce compatibility shims
   - Add `parser/internal/compat` helpers (`AsFloat`, `AsString`, `AsBool`).
   - Keep old behavior working while incrementally refactoring consumers.

2. Land typed TokenValue
   - Change `Token.Value` to `TokenValue`.
   - Update parser/tokenizer to populate TokenValue.
   - Update a small set of compiler/vm hot spots using compat helpers first.

3. Update AST property typing
   - Replace `any` with `Property` in AST nodes.
   - Update compiler property handling (minimally invasive, switch → methods).

4. Parser options
   - Add `Options`, `DefaultOptions()`, `NewParserWithOptions`.
   - Keep `NewParser` delegating to defaults to minimize churn.

5. Clean up
   - Remove compat helpers once consumers stop using them.
   - Re-run full test and benchmark suites; compare with baselines.

## Rollout strategy

- Branching: perform changes on a feature branch `feat/parser-typed-api`.
- CI gates: require `go test ./...`, race, and benchmarks for parser/compiler/vm to pass.
- Versioning: bump major version in go.mod to signal breaking changes.
- Changelog: document migration snippets similar to examples above.

## Risks and mitigations

- Risk: Wide impact of Token.Value change.
  - Mitigate with compatibility helpers and staged refactors.
- Risk: Hidden runtime assumptions in vm/compiler.
  - Mitigate with compile-time breaks, exhaustive tests, and targeted greps for `.(type)`.
- Risk: Performance regressions due to new structs.
  - Mitigate with benchmarks; `TokenValue` is stack-friendly and should reduce interface costs.

## Acceptance criteria

- All packages build with new typed TokenValue and Property.
- All existing tests pass unchanged or with minimal, mechanical updates.
- Benchmarks show neutral-to-better performance across parser/compiler/vm critical paths.
- Documentation updated: this roadmap, package docs, and migration notes.

## Timeline (suggested)

- Week 1: Introduce TokenValue + compat; refactor compiler hot paths.
- Week 2: Property typing in AST; adjust vm/compiler; stabilize tests.
- Week 3: Parser options + cleanups; remove compat; final benchmarks and docs.

## Appendix: What we’ve already improved

- Tokenizer correctness and speed-ups (UTF-8-aware peek, fewer line/col inconsistencies, retained ASCII fast paths) with double-digit improvements on larger inputs.
- Package docs and review notes updated with measurable wins and quality gates.

---

## Exact impact points in compiler (files and anchors)

These references use current repository structure and common code paths observed; adjust line numbers as code evolves.

- `compiler/compiler.go`
    - access step representation (near top of file):
        - Today uses `property any` to carry either a string (member) or `parser.Node` (index expression).
        - Action: replace with typed fields:
            - `propertyStr string` for member names
            - `propertyExpr parser.Node` for bracket/index expressions
    - Case handling of `*parser.MemberAccess` and `*parser.IndexAccess` (around emit/compile phases):
        - Member: read `v.Property` via new `Property` type (`if v.Property.IsString() { v.Property.S }`).
        - Index: use `v.Index` as before (expression), assign to `propertyExpr`.

- `compiler/compiler_utils.go`
    - Building access chains (functions that linearize member/index access):
        - Previously: `property: v.Property // any`
        - After: `propertyStr: v.Property.S` or `propertyExpr: v.Index`.
    - Emitting constants for member names:
        - Previously: `propIdx := c.addConstant(step.property)` (where `property` was any/string)
        - After: `propIdx := c.addConstant(step.propertyStr)`
    - Removing type assertions:
        - Previously: `idxExpr, ok := step.property.(parser.Node)` → After: use `step.propertyExpr` directly.

- `compiler/tests/help_test.go`
    - Parser construction:
        - If `NewParser` remains, no change.
        - Otherwise use: `parser.NewParserWithOptions(input, parser.DefaultOptions())`.

Notes:
- The compiler compiles from AST, not tokens; the `TokenValue` change does not directly affect compiler code.

### Compiler migration snippet (illustrative)

```go
// Before
type accessStep struct {
        safe     bool
        property any // string or parser.Node
}

// After
type accessStep struct {
        safe         bool
        propertyStr  string      // member name
        propertyExpr parser.Node // index expression
}

// When building steps
switch v := node.(type) {
case *parser.MemberAccess:
        if v.Property.IsString() {
                steps = append(steps, accessStep{safe: v.Optional, propertyStr: v.Property.S})
        } else {
                steps = append(steps, accessStep{safe: v.Optional, propertyExpr: &parser.NumberLiteral{Value: float64(v.Property.I)}})
        }
case *parser.IndexAccess:
        steps = append(steps, accessStep{safe: v.Optional, propertyExpr: v.Index})
}

// When emitting
if step.propertyStr != "" {
        propIdx := c.addConstant(step.propertyStr)
        c.emit(code.OpConstant, propIdx)
        c.emit(code.OpMemberAccess)
} else if step.propertyExpr != nil {
        if err := c.Compile(step.propertyExpr); err != nil { return err }
        c.emit(code.OpIndexAccess)
}
```

---

## Exact impact points in VM (files and anchors)

VM predominantly consumes bytecode and runtime values; it does not depend on parser tokens or AST internals.

- `vm/vm_test.go`, `vm/nullish_semantics_errors_test.go`
    - Parser construction in tests:
        - If `NewParser` remains: no changes.
        - If not: switch to `NewParserWithOptions(input, parser.DefaultOptions())`.

- `vm/vm.go`, `vm/vm_handlers.go`
    - Opcode semantics unaffected:
        - `OpMemberAccess` expects member name as a string constant on the stack (pushed by compiler).
        - `OpIndexAccess` evaluates the index expression (emitted by compiler) and uses existing handlers.
    - No changes required to handler signatures or behavior.

---

## Grep checklist to drive migration

Use these to locate and update all relevant sites:

```bash
# Compiler: removal of any-based property
grep -RIn "property\s\+any" compiler

# Compiler: places that handle MemberAccess / IndexAccess
grep -RIn "MemberAccess\|IndexAccess" compiler

# Compiler: pushing property constants
grep -RIn "addConstant\(step\.property" compiler

# Parser construction in tests (compiler + vm)
grep -RIn "NewParser\(" compiler vm

# VM handlers for awareness (no changes expected)
grep -RIn "OpMemberAccess\|executeMemberAccess\|OpIndexAccess" vm
```
