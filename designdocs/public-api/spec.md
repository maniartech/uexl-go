# UExL Public API — Design Specification

**Status:** Draft
**Scope:** `github.com/maniartech/uexl` (root package, `uexl.go` + `doc.go`)
**Implements:** Launch Checklist §2 — Public API

---

## 1. Goals

| Goal | Rationale |
|---|---|
| One-shot convenience (`Eval`) | Lowest possible barrier for simple use cases |
| Compile-once / run-many (`CompiledExpr`) | Avoid re-parsing and re-compiling hot expressions |
| Layered, immutable environments (`Env`) | Stdlib built-ins → lib env → app env → per-call vars without mutation or locking |
| Multi-level user libs (`Lib`) | Package authors ship reusable bundles of functions, pipes, and globals as a single composable unit |
| Goroutine-safety without external synchronization | All public types safe to share across goroutines after construction |
| Zero-panic surface | All error paths return `error`; no internal panics escape to caller |
| Introspection / debugging | Operators can enumerate registered functions and pipes and check presence at any time |
| Validate without compile artifact | `Validate` checks syntax and compilation correctness without allocating a `*CompiledExpr` |
| Result type ergonomics | Typed coercion helpers eliminate per-call type assertions on `Eval` results |

---

## 2. Package-Level Symbols (exported from `uexl.go`)

### 2.1 Types

```
Option          — opaque functional option applied to an Env at construction
Lib             — shareable, composable bundle of functions + pipes + globals
Env             — immutable environment (functions + pipe handlers + global vars)
EnvInfo         — snapshot of an Env's registered symbols, used for introspection
CompiledExpr    — immutable pre-compiled expression, produced by Env.Compile
```

> **Note:** `Function`, `PipeContext`, `Functions`, `PipeHandler`, `PipeHandlers`, `ParserError`, and `ParseErrors` are type aliases, not new types — see §2.2.

### 2.2 Type Re-exports (avoid forcing users to import internal packages)

Type aliases are defined at the `uexl` package level so callers never need to import internal sub-packages directly.

**From `vm` package:**
```
Function     = vm.VMFunction      — func(args ...any) (any, error)
Functions    = vm.VMFunctions     — map[string]Function
PipeHandler  = vm.PipeHandler     — func(ctx PipeContext, input any) (any, error)
PipeHandlers = vm.PipeHandlers    — map[string]PipeHandler
PipeContext  = vm.PipeContext     — interface for pipe predicate evaluation (see §3.25)
```

**From `parser/errors` package:**
```
ParserError  = parsererrors.ParserError   — single structured parse error (value type)
ParseErrors  = parsererrors.ParseErrors   — collection of parse errors (value type)
```

All of these are Go type aliases (`=`), not new types — existing values from the originating packages are directly assignable without conversion.

> **Why only parser errors?** Compile errors and runtime errors are returned as plain `error` because there is no useful concrete type to expose — they carry only a message string. Only the parser produces structured, field-rich error values (`Line`, `Column`, `Code`, `Message`) worth surfacing directly in the public API.

### 2.3 Functions

```
// Environment construction
Default()                                            *Env
DefaultWith(opts ...Option)                          *Env
NewEnv(opts ...Option)                               *Env

// One-shot evaluation and validation (use Default env + context.Background())
Eval(expr string, vars map[string]any)               (any, error)
Validate(expr string)                                error
MustCompile(expr string)                             *CompiledExpr

// Options
WithFunctions(fns Functions)                         Option
WithPipeHandlers(pipes PipeHandlers)                 Option
WithGlobals(vars map[string]any)                     Option
WithLib(lib Lib)                                     Option

// Result coercion helpers (no dependency on Env)
AsFloat64(v any)                                     (float64, error)
AsBool(v any)                                        (bool, error)
AsString(v any)                                      (string, error)
AsSlice(v any)                                       ([]any, error)
AsMap(v any)                                         (map[string]any, error)
```

### 2.4 Methods on `*Env`

```
// Extension and evaluation
(e *Env) Extend(opts ...Option)                                      *Env
(e *Env) Compile(expr string)                                        (*CompiledExpr, error)
(e *Env) MustCompile(expr string)                                    *CompiledExpr
(e *Env) Validate(expr string)                                       error
(e *Env) Eval(ctx context.Context, expr string, vars map[string]any) (any, error)

// Introspection — read-only, goroutine-safe
(e *Env) Info()                                      EnvInfo
(e *Env) HasFunction(name string)                    bool
(e *Env) HasPipe(name string)                        bool
(e *Env) HasGlobal(name string)                      bool
```

### 2.5 `Lib` interface

```go
type Lib interface {
    // Apply registers the lib's functions, pipe handlers, and globals
    // into the supplied config. Called once during Env construction.
    Apply(cfg *EnvConfig)
}
```

### 2.6 `EnvInfo` (value type, read-only snapshot)

```
type EnvInfo struct {
    Functions    []string  // sorted function names
    PipeHandlers []string  // sorted pipe handler names
    Globals      []string  // sorted global variable names
}
```

### 2.7 Methods on `*CompiledExpr`

```
(c *CompiledExpr) Eval(ctx context.Context, vars map[string]any) (any, error)
(c *CompiledExpr) Variables()                                   []string
(c *CompiledExpr) Env()                                         *Env
```

`Eval` executes the pre-compiled bytecode against `vars`, respecting `ctx` for cancellation and deadline. `Variables()` returns the sorted list of variable names (without `$` prefix) that the expression references. `Env()` returns the `*Env` the expression was compiled against, useful for introspection and logging.

---

### 2.8 Complete API Quick Reference

A consolidated, scannable view of every exported symbol. Use this as the go/no-go checklist before implementation begins.

#### Concrete Types

| Type | Kind | Goroutine-safe | Notes |
|---|---|---|---|
| `Option` | `func(*envConfig)` | n/a | Opaque; create only via `With*` helpers |
| `Env` | struct (pointer) | ✅ after construction | Immutable fields; own `sync.Pool` |
| `EnvInfo` | struct (value) | ✅ | Read-only snapshot; safe to copy and pass around |
| `EnvConfig` | struct (pointer) | scoped | Only accessible inside `Lib.Apply`; never stored |
| `CompiledExpr` | struct (pointer) | ✅ | Immutable bytecode + parent env; pool-based eval |
| `Lib` | interface | n/a | Implemented by third-party library packages |

#### Type Aliases (re-exported — no sub-package import needed)

| Alias | Resolves to | Underlying shape |
|---|---|---|
| `Function` | `vm.VMFunction` | `func(args ...any) (any, error)` |
| `Functions` | `vm.VMFunctions` | `map[string]Function` |
| `PipeHandler` | `vm.PipeHandler` | `func(ctx PipeContext, input any) (any, error)` |
| `PipeHandlers` | `vm.PipeHandlers` | `map[string]PipeHandler` |
| `PipeContext` | `vm.PipeContext` | interface (`EvalItem`, `EvalWith`, `Alias`, `Context`) |
| `ParserError` | `parsererrors.ParserError` | value struct (`Line`, `Column`, `Code`, `Message`) |
| `ParseErrors` | `parsererrors.ParseErrors` | value struct (`Errors []ParserError`) |

#### Package-level Functions — Environment Construction

| Signature | Returns | Notes |
|---|---|---|
| `Default()` | `*Env` | Singleton; stdlib builtins + default pipes; no globals |
| `DefaultWith(opts ...Option)` | `*Env` | `Default().Extend(opts...)`; does not mutate singleton |
| `NewEnv(opts ...Option)` | `*Env` | Blank slate — no builtins, no pipes, no globals |

#### Package-level Functions — Options

| Signature | Panics on nil? | Notes |
|---|---|---|
| `WithFunctions(fns Functions) Option` | ✅ | Merges into env; later call wins for same key |
| `WithPipeHandlers(pipes PipeHandlers) Option` | ✅ | Same merge semantics as `WithFunctions` |
| `WithGlobals(vars map[string]any) Option` | ✅ | Env-level vars; shadowed by per-call vars |
| `WithLib(lib Lib) Option` | ✅ | Calls `lib.Apply(cfg)` during construction |

#### Package-level Functions — One-Shot Evaluation

| Signature | Returns | Notes |
|---|---|---|
| `Eval(expr string, vars map[string]any)` | `(any, error)` | Uses `Default()` + `context.Background()` internally |
| `Validate(expr string)` | `error` | Uses `Default()`; no `*CompiledExpr` artifact |
| `MustCompile(expr string)` | `*CompiledExpr` | Uses `Default()`; panics — for `var` declarations only |

#### Package-level Functions — Result Coercion Helpers

| Signature | Widening | Truthy coercion |
|---|---|---|
| `AsFloat64(v any) (float64, error)` | `int`, `int64`, `float32` → `float64` | — |
| `AsBool(v any) (bool, error)` | none | ❌ `AsBool(1)` → error |
| `AsString(v any) (string, error)` | none | — |
| `AsSlice(v any) ([]any, error)` | none | — |
| `AsMap(v any) (map[string]any, error)` | none | — |

#### Methods on `*Env`

| Signature | Returns | Notes |
|---|---|---|
| `Extend(opts ...Option)` | `*Env` | New child env; parent unchanged; own pool |
| `Compile(expr string)` | `(*CompiledExpr, error)` | Parse + compile + fn-name validation; no VM |
| `MustCompile(expr string)` | `*CompiledExpr` | Panics — for startup `var` declarations only |
| `Validate(expr string)` | `error` | Thin wrapper: `_, err := e.Compile(expr)` |
| `Eval(ctx, expr string, vars map[string]any)` | `(any, error)` | One-shot: Compile + Eval; borrows VM from pool |
| `Info()` | `EnvInfo` | Sorted snapshot of all symbols; goroutine-safe |
| `HasFunction(name string)` | `bool` | `false` for empty string |
| `HasPipe(name string)` | `bool` | `false` for empty string |
| `HasGlobal(name string)` | `bool` | `false` for empty string |

#### Methods on `*CompiledExpr`

| Signature | Returns | Notes |
|---|---|---|
| `Eval(ctx context.Context, vars map[string]any)` | `(any, error)` | Hot path; borrows `*vm.VM` from env pool |
| `Variables()` | `[]string` | Sorted; derived from `bytecode.ContextVars`; copy |
| `Env()` | `*Env` | Allocation-free pointer return |

#### Methods on `EnvInfo`

| Signature | Returns | Notes |
|---|---|---|
| `String()` | `string` | Implements `fmt.Stringer`; multiline; stable sort |

#### Methods on `*EnvConfig` (only accessible inside `Lib.Apply`)

| Signature | Panics on nil? | Notes |
|---|---|---|
| `AddFunctions(fns Functions)` | ✅ | Merges into in-progress config |
| `AddPipeHandlers(pipes PipeHandlers)` | ✅ | Same merge semantics |
| `AddGlobals(vars map[string]any)` | ✅ | Same merge semantics |

---

## 3. Detailed Specification

### 3.1 `Option`

```go
type Option func(*envConfig)
```

An opaque type. Users never construct an `Option` directly; they call one of the three constructor helpers (`WithFunctions`, `WithPipeHandlers`, `WithGlobals`). The alias to `func(*envConfig)` is intentional — it avoids an extra interface indirection and keeps the implementation trivial.

`envConfig` is an **unexported** struct used only during `NewEnv` / `Extend` to accumulate option state before freezing it into an `Env`.

```go
// unexported — never exposed
type envConfig struct {
    functions    vm.VMFunctions
    pipeHandlers vm.PipeHandlers
    globals      map[string]any
}
```

---

### 3.2 `Env`

```go
type Env struct {
    functions    vm.VMFunctions   // merged, read-only after construction
    pipeHandlers vm.PipeHandlers  // merged, read-only after construction
    globals      map[string]any   // merged, read-only after construction
    pool         sync.Pool        // pool of *vm.VM for this env's config
}
```

**All fields are unexported.** `Env` is goroutine-safe to use after construction.

#### Immutability contract
- No method on `*Env` mutates its fields.
- `Extend` always produces a **new** `Env` by copying parent fields and applying options on top. The parent is unchanged.
- The `pool` field is **not** copied on `Extend`; each `Env` maintains its own pool keyed to its own function/pipe configuration.

#### `pool` details
- `pool.New` creates a `*vm.VM` via `vm.New(vm.LibContext{Functions: e.functions, PipeHandlers: e.pipeHandlers})`.
- VMs are returned to the pool via `defer` in every `Eval` path.
- Pool is per-`Env` so each env's VM is pre-configured with the correct functions and pipe handlers — no reconfiguration on borrow.

---

### 3.3 `Default() *Env`

Returns an `*Env` pre-loaded with:
- `vm.Builtins` as the function set (all built-in functions: `len`, `substr`, `contains`, `str`, `runeLen`, `runeSubstr`, `graphemeLen`, `graphemeSubstr`, `runes`, `graphemes`, `bytes`, etc.)
- `vm.DefaultPipeHandlers` as the pipe handler set (`map`, `filter`, `reduce`, `find`, `some`, `every`, `unique`, `sort`, `groupBy`, `window`, `chunk`, `flatMap`, `pipe`)
- No globals.

`Default()` returns the **same singleton** on every call — it is initialized once via `sync.Once` at package init time. Callers MUST NOT mutate the returned value (impossible because all fields are unexported, but noted for documentation clarity).

```go
var (
    defaultEnvOnce sync.Once
    defaultEnv     *Env
)

func Default() *Env {
    defaultEnvOnce.Do(func() {
        defaultEnv = NewEnv(
            WithFunctions(vm.Builtins),
            WithPipeHandlers(vm.DefaultPipeHandlers),
        )
    })
    return defaultEnv
}
```

---

### 3.4 `NewEnv(opts ...Option) *Env`

Creates an `Env` from a blank slate — **no built-ins, no pipes, no globals** — then applies all provided options in order.

Steps:
1. Initialize an empty `envConfig`.
2. Apply each `Option` to `envConfig` in order (left to right). Later options for the same key win.
3. Shallow-copy the resulting maps into a frozen `Env`.
4. Initialize `pool.New` to create `*vm.VM` using the frozen config.
5. Return `&Env`.

**Merging behavior within a single `NewEnv` call:** If `WithFunctions` is called twice, the second call's entries are merged on top of the first. No entries are deleted.

---

### 3.5 `WithFunctions(fns Functions) Option`

Returns an `Option` that merges `fns` into the in-progress `envConfig.functions` map. Existing keys from a prior option are overwritten if `fns` contains the same key.

`Functions` is a type alias for `vm.VMFunctions` defined at the `uexl` package level — no import of the `vm` package required.

**Panics if `fns` is nil.** (Validated at option application time, not deferred to use, so misuse surfaces immediately at app startup.)

---

### 3.6 `WithPipeHandlers(pipes PipeHandlers) Option`

Returns an `Option` that merges `pipes` into `envConfig.pipeHandlers`. Same semantics as `WithFunctions`.

`PipeHandlers` is a type alias for `vm.PipeHandlers` defined at the `uexl` package level.

**Panics if `pipes` is nil.**

---

### 3.7 `WithGlobals(vars map[string]any) Option`

Returns an `Option` that merges `vars` into `envConfig.globals`. These become env-level context variables, resolved **before** per-call vars in the resolution order (see §4).

**Panics if `vars` is nil.**

---

### 3.8 `WithLib(lib Lib) Option`

Returns an `Option` that calls `lib.Apply(cfg)` during env construction, allowing the lib to register any combination of functions, pipe handlers, and globals in a single step.

`WithLib` may be combined freely with `WithFunctions`, `WithPipeHandlers`, and `WithGlobals` in the same `NewEnv` or `Extend` call. Options are applied left-to-right, so a later `WithFunctions` or a later `WithLib` can override a key that an earlier lib registered.

**Panics if `lib` is nil.**

Package authors implement the `Lib` interface to publish a distributable library:

```go
// In a third-party package, e.g., github.com/acme/uexlfinance:

type FinanceLib struct{}

func (FinanceLib) Apply(cfg *uexl.EnvConfig) {
    cfg.AddFunctions(uexl.VMFunctions{
        "pv":  presentValue,
        "fv":  futureValue,
        "npv": netPresentValue,
    })
    cfg.AddPipeHandlers(uexl.PipeHandlers{
        "amortize": amortizePipeHandler,
    })
}
```

> **Note:** `EnvConfig` (the public face of `envConfig`) exposes only `AddFunctions`, `AddPipeHandlers`, and `AddGlobals` methods — no direct field access. This keeps the internal merging logic centralized and prevents libraries from accidentally replacing rather than merging.

---

### 3.9 `Lib` interface

```go
type Lib interface {
    Apply(cfg *EnvConfig)
}
```

`EnvConfig` is the **public projection** of the internal `envConfig`. It is only accessible inside an `Apply` call — the pointer is never stored or returned anywhere else.

```go
// Public — only constructors and Apply see this, never stored long-term
type EnvConfig struct {
    cfg *envConfig  // unexported backing store
}

func (c *EnvConfig) AddFunctions(fns Functions)
func (c *EnvConfig) AddPipeHandlers(pipes PipeHandlers)
func (c *EnvConfig) AddGlobals(vars map[string]any)
```

Each method panics on nil input (same contract as `WithFunctions` etc.).

---

### 3.10 Multi-Level Extension Model

`Extend` is arbitrarily chainable. Each call produces a new, independent, immutable `Env` that inherits from its parent. There is no depth limit.

```
Default()                              ← stdlib layer (builtins + default pipes)
  └── Extend(WithLib(FinanceLib{}))    ← domain lib layer
        └── Extend(WithLib(TenantLib{}), WithGlobals(tenantConfig))  ← tenant layer
              └── Extend(WithFunctions(requestFns))                  ← per-request layer
```

At each level:
- All symbols from the parent are inherited.
- New symbols are added; conflict resolution is "last writer wins" (child over parent, left-to-right among options in a single call).
- No level can mutate any ancestor.
- Each level has its own `sync.Pool` of VMs configured for its exact symbol set.

This model supports the following reuse patterns without any special mechanism:

| Pattern | How |
|---|---|
| Shared stdlib | Call `Default()` once; all envs in the process share it |
| Domain lib (e.g., finance, HR) | Implement `Lib`, ship as a package, apply via `WithLib` |
| Multi-tenant customization | One app env; per tenant call `Extend(WithLib(tenantLib))` |
| Per-request context injection | One compiled expr; inject via `Eval(requestVars)` |
| Feature flags as globals | `appEnv.Extend(WithGlobals(featureFlags))` |

---

### 3.11 `Eval(expr string, vars map[string]any) (any, error)`  *(package-level)*

Package-level convenience function. Equivalent to:

```go
func Eval(expr string, vars map[string]any) (any, error) {
    return Default().Eval(context.Background(), expr, vars)
}
```

This function always uses the default stdlib environment. No configuration is possible. For custom functions or pipes, use `NewEnv` or `Default().Extend(...)`.

`vars` may be `nil` — treated as an empty map.

---

### 3.12 `(*Env).Extend(opts ...Option) *Env`

Creates a new `*Env` inheriting all functions, pipe handlers, and globals from the receiver, then applies `opts` on top.

Steps:
1. Copy receiver's `functions`, `pipeHandlers`, and `globals` into a new `envConfig` (shallow copy of each map).
2. Apply each option in `opts` over the copy (same merge logic as `NewEnv`).
3. Freeze into a new `Env` with its own `pool`.
4. Return the new `Env`. Receiver is unchanged.

**Key invariant:** If the caller does not override a key, the child inherits it from the parent exactly. If the caller overrides a key (function name, pipe name, global var name), the child's version takes precedence for all `Eval`/`Compile` calls on the child. The parent retains its original value.

---

### 3.13 `(*Env).Compile(expr string) (*CompiledExpr, error)`

Parses and compiles `expr` into a `*CompiledExpr`. Compile is a pure transformation — it allocates no VM.

Steps:
1. `parser.ParseString(expr)` → AST; return error on parse failure.
2. `compiler.New().Compile(ast)` → error on compile failure.
3. Validate all function call sites against `e.functions`: for each `OpCall "name"` in the bytecode, check whether `"name"` is a key in `e.functions`. If not, return a compile error: `compile error: unknown function "<name>" — not registered in this environment`.
4. Wrap `*compiler.ByteCode` and a reference to the parent `*Env` into a `*CompiledExpr`.
5. Return `*CompiledExpr, nil`.

**Compile-time function existence checking:** Since `Compile` is called on `*Env` (which holds the function registry), function names referenced in the expression can be validated at compile time. This means errors like calling `discount(price)` in an env that has no `discount` function are caught immediately at `Compile`/`Validate` time rather than at eval time. This is a significant improvement over the previous design and closes the gap with cel-go.

**Note on `NewEnv()` (blank env):** A blank env has no functions registered. Compiling `"len('hi')"` against a blank env will fail with `compile error: unknown function "len"`. This is correct and expected — use `Default()` or `DefaultWith(...)` for access to the stdlib.

`*CompiledExpr` is goroutine-safe. The `*compiler.ByteCode` it holds is immutable after compilation. Multiple goroutines may call `compiledExpr.Eval(...)` concurrently.

---

### 3.14 `(*Env).Eval(ctx context.Context, expr string, vars map[string]any) (any, error)`

One-shot parse + compile + run within the environment. Context is forwarded to the VM for cancellation and deadline enforcement.

```go
func (e *Env) Eval(ctx context.Context, expr string, vars map[string]any) (any, error) {
    ce, err := e.Compile(expr)
    if err != nil {
        return nil, err
    }
    return ce.Eval(ctx, vars)
}
```

Use this for expressions evaluated once. Use `Compile` + `CompiledExpr.Eval` for expressions evaluated many times.

**Package-level `Eval`** is the convenience path for novices — it uses `context.Background()` internally and omits the ctx parameter entirely:

```go
func Eval(expr string, vars map[string]any) (any, error) {
    return Default().Eval(context.Background(), expr, vars)
}
```

---

### 3.15 `(*CompiledExpr).Eval(ctx context.Context, vars map[string]any) (any, error)`

Executes the pre-compiled bytecode against `vars`, honoring `ctx` for cancellation and deadline. This is the **hot path**.

Steps:
1. Check `ctx.Err()` before starting — return immediately if already cancelled.
2. Borrow a `*vm.VM` from `env.pool` — no allocation if pool has a spare.
3. Set the context on the borrowed VM (for pipe handlers and cooperative cancellation).
4. Build the merged variable map (see §4 for resolution order).
5. Call `machine.Run(bytecode, mergedVars)`.
6. Return the VM to pool via `defer pool.Put(machine)`.
7. Return the result or error.

**Variable merging (step 4):** Produce a single `map[string]any` where eval-level vars shadow env globals. For performance, if `env.globals` is empty (no globals registered), pass `vars` directly without creating a new map. If globals are non-empty, allocate a merged map: start from a copy of `env.globals`, then apply `vars` on top.

**Context cancellation:** The VM checks `ctx.Done()` at loop boundaries (between opcode executions) and returns `ctx.Err()` if the context is cancelled. Long-running pipe iterations (e.g., `|reduce:` over a large array) will be cancelled within one iteration of detecting cancellation.

`vars` may be `nil` — treated as empty.

---

### 3.16 Introspection API

All three methods are **read-only** and goroutine-safe. They never allocate VM instances or perform compilation.

#### `(*Env).Info() EnvInfo`

Returns a **value-type snapshot** of everything registered in the environment. The snapshot is independent of the `Env` — mutating the returned slices has no effect.

```go
type EnvInfo struct {
    Functions    []string // sorted, all registered function names
    PipeHandlers []string // sorted, all registered pipe handler names
    Globals      []string // sorted, all registered global variable names
}
```

Example usage:

```go
info := appEnv.Info()
fmt.Println("Functions:",    info.Functions)
fmt.Println("Pipe handlers:", info.PipeHandlers)
fmt.Println("Globals:",       info.Globals)
```

Sorted order is guaranteed so that output is stable and diffable across runs.

---

#### `(*Env).HasFunction(name string) bool`

Returns `true` if a function with the given name is registered in this env (including inherited). Returns `false` for empty string.

```go
if !appEnv.HasFunction("discount") {
    log.Fatal("discount function not registered")
}
```

---

#### `(*Env).HasPipe(name string) bool`

Returns `true` if a pipe handler with the given name is registered. Returns `false` for empty string.

```go
if !appEnv.HasPipe("amortize") {
    log.Fatal("amortize pipe not registered")
}
```

---

#### `(*Env).HasGlobal(name string) bool`

Returns `true` if a global variable with the given name is registered. Returns `false` for empty string.

```go
if !tenantEnv.HasGlobal("currency") {
    log.Fatal("currency global not set")
}
```

---

#### `EnvInfo.String() string`

`EnvInfo` implements `fmt.Stringer`. Output format is fixed and human-readable for use in logs, tests, and debug output:

```
Env:
  Functions (12): contains, discount, fv, graphemeLen, graphemeSubstr, graphemes, len, npv, pv, runeLen, runeSubstr, runes
  PipeHandlers (14): amortize, chunk, every, filter, find, flatMap, groupBy, map, pipe, reduce, some, sort, unique, window
  Globals (2): appVersion, currency
```

This output is produced by `fmt.Println(info)` with no additional formatting code by the caller.

---

### 3.17 `DefaultWith(opts ...Option) *Env`

Shorthand for `Default().Extend(opts...)`. Creates a new `Env` with all stdlib built-ins and default pipe handlers, plus any additional options applied on top.

```go
func DefaultWith(opts ...Option) *Env {
    return Default().Extend(opts...)
}
```

This is the single most common entry point for apps that want the standard library plus custom functions:

```go
// Without DefaultWith (verbose):
env := uexl.Default().Extend(uexl.WithFunctions(myFns))

// With DefaultWith (idiomatic):
env := uexl.DefaultWith(uexl.WithFunctions(myFns))
```

`DefaultWith` does **not** modify the stdlib singleton; it always produces a new child `Env`.

---

### 3.18 `MustCompile(expr string) *CompiledExpr`  *(package-level)*

Compiles `expr` using `Default()` and panics on failure. Intended exclusively for **package-level `var` declarations** where expressions are known-valid at write time.

```go
func MustCompile(expr string) *CompiledExpr {
    c, err := Default().Compile(expr)
    if err != nil {
        panic(fmt.Sprintf("uexl: MustCompile: %v", err))
    }
    return c
}
```

Typical use:

```go
var totalRule   = uexl.MustCompile("price * qty * (1 - discount)")
var filterAdult = uexl.MustCompile("age >= 18")
```

**Do NOT use in request-handling paths.** For dynamic or user-supplied expressions, always use `Compile` and handle the error.

---

### 3.19 `(*Env).MustCompile(expr string) *CompiledExpr`

Same contract as the package-level `MustCompile` but compiles within a specific `Env` (its functions, pipes, and globals are available).

```go
func (e *Env) MustCompile(expr string) *CompiledExpr {
    c, err := e.Compile(expr)
    if err != nil {
        panic(fmt.Sprintf("uexl: Env.MustCompile: %v", err))
    }
    return c
}
```

Use when the env is constructed at startup and expressions are known-valid:

```go
appEnv := uexl.DefaultWith(uexl.WithFunctions(domainFns))

var discountRule = appEnv.MustCompile("price |filter: discount($item) > 0")
```

---

### 3.20 `(*CompiledExpr).Variables() []string`

Returns the **sorted list of variable names** (without `$` prefix) that the expression references. The list is derived directly from `compiler.ByteCode.ContextVars`, which is populated at compile time. Returns an empty slice (never `nil`) when the expression has no variable references.

```go
func (c *CompiledExpr) Variables() []string
```

The returned slice is a **copy** — mutating it has no effect.

**Use cases:**

```go
prog, _ := appEnv.Compile("price * qty - discount")
vars := prog.Variables()
// vars == ["discount", "price", "qty"]

// Preflight: ensure all needed vars are present before eval
for _, v := range prog.Variables() {
    if _, ok := inputRecord[v]; !ok {
        return fmt.Errorf("missing required field: %s", v)
    }
}

// Form schema: tell the UI which fields a rule reads
schema := prog.Variables() // drive which columns to fetch
```

---

### 3.21 `(*CompiledExpr).Env() *Env`

Returns the `*Env` the expression was compiled against. Useful for logging, debugging, and asserting that an expression has access to the functions it needs.

```go
func (c *CompiledExpr) Env() *Env
```

```go
// Note: with compile-time function validation (§3.13), this check is now redundant
// — the Compile call would have already failed if "discount" was not registered.
// Env() is still useful for introspection in logging and diagnostics.
prog, _ := appEnv.Compile("discount(price, tier)")
log.Printf("compiled against env with %d functions", len(prog.Env().Info().Functions))
```

The returned pointer is the same `*Env` passed at compile time; calling `Env()` is allocation-free.

---

### 3.22 `Validate(expr string) error`  *(package-level)*

Parses, compiles, and validates `expr` against the default stdlib environment. Returns `nil` if the expression is syntactically valid, all functions are registered, and bytecode compiles successfully. No `*CompiledExpr` artifact is allocated in the success case.

```go
func Validate(expr string) error {
    return Default().Validate(expr)
}
```

---

### 3.23 `(*Env).Validate(expr string) error`

Parses, compiles, and validates `expr` within this environment's context. Returns the first error, or `nil`. Because `Validate` calls `Compile` internally (and `Compile` now validates function names — see §3.13), this detects both syntax errors and calls to unregistered functions.

```go
func (e *Env) Validate(expr string) error {
    _, err := e.Compile(expr)
    return err
}
```

Example uses: REST endpoint validation, CI lint step for rule libraries, pre-storage sanity check.

---

### 3.24 Result Coercion Helpers

Five package-level functions provide safe, typed extraction from `any` values returned by `Eval`. These have no dependency on `Env` or `CompiledExpr` and may be called on any `any` value.

```go
func AsFloat64(v any) (float64, error)
func AsBool(v any)    (bool, error)
func AsString(v any)  (string, error)
func AsSlice(v any)   ([]any, error)
func AsMap(v any)     (map[string]any, error)
```

Each function attempts a direct type assertion, then falls back to numeric widening or conversion where reasonable. They return an error (never panic) if the value cannot be converted.

**Conversion rules:**

| Target | Direct | Widening / conversion |
|---|---|---|
| `float64` | `float64` | `int`, `int64`, `float32` via numeric cast |
| `bool` | `bool` | none — no truthy coercion |
| `string` | `string` | none — no `fmt.Sprint` fallback |
| `[]any` | `[]any` | none |
| `map[string]any` | `map[string]any` | none |

**No truthy coercion:** `AsBool(0)` returns an error, not `false`. UExL's explicit nullish/boolish semantics (see design-philosophy.md) apply here too.

Example:

```go
result, err := env.Eval("price * qty", vars)
if err != nil { ... }

total, err := uexl.AsFloat64(result)
if err != nil {
    return fmt.Errorf("expected numeric result, got %T", result)
}
```

These helpers are defined in a new file `result.go` (see §6).

---

### 3.25 `PipeContext` interface

The `PipeContext` interface is the sole parameter passed to custom pipe handlers. It completely replaces the previous `*vm.VM` parameter, eliminating the need to import internal VM packages when writing custom pipes.

```go
// Defined in vm package, re-exported as uexl.PipeContext via type alias.
type PipeContext interface {
    // EvalItem runs the pipe predicate with $item=item and $index=index set in scope.
    // Optimized for iteration: the underlying VM frame is reused across calls.
    EvalItem(item any, index int) (any, error)

    // EvalWith runs the pipe predicate with arbitrary scope variables.
    // Use for accumulation patterns (reduce: $acc), window ($window), chunk ($chunk), etc.
    EvalWith(scopeVars map[string]any) (any, error)

    // Alias returns the user-defined alias from the pipe expression (e.g., "$x" from
    // "|map as $x:"). Returns empty string if no alias was specified.
    Alias() string

    // Context returns the evaluation context, enabling cancellation and deadline checks.
    Context() context.Context
}
```

**Updated `PipeHandler` signature:**

```go
// Old (required importing vm package):
type PipeHandler func(input any, block any, alias string, vm *VM) (any, error)

// New (self-contained — only import "github.com/maniartech/uexl"):
type PipeHandler func(ctx PipeContext, input any) (any, error)
```

**Writing a custom pipe with the new API:**

```go
// A pipe that doubles all numeric items: |double:
func doublePipe(ctx uexl.PipeContext, input any) (any, error) {
    items, ok := input.([]any)
    if !ok {
        return nil, fmt.Errorf("double pipe expects array input")
    }
    result := make([]any, len(items))
    for i, item := range items {
        val, err := ctx.EvalItem(item, i)  // runs predicate with $item, $index set
        if err != nil {
            return nil, err
        }
        result[i] = val
    }
    return result, nil
}

// Usage:
env := uexl.DefaultWith(uexl.WithPipeHandlers(uexl.PipeHandlers{
    "double": doublePipe,
}))
```

A custom `reduce`-style pipe using `EvalWith`:

```go
func runningTotalPipe(ctx uexl.PipeContext, input any) (any, error) {
    items, ok := input.([]any)
    if !ok {
        return nil, fmt.Errorf("runningTotal expects array")
    }
    var acc any
    for i, item := range items {
        if ctx.Context().Err() != nil {  // respect cancellation
            return nil, ctx.Context().Err()
        }
        var err error
        acc, err = ctx.EvalWith(map[string]any{
            "$acc":   acc,
            "$item":  item,
            "$index": i,
        })
        if err != nil {
            return nil, err
        }
    }
    return acc, nil
}
```

**Implementation note:** The VM creates an unexported `pipeContextImpl` struct that implements this interface for each pipe invocation. Frame reuse and scope management happen inside `EvalItem`. Alias is captured at OpPipe dispatch time.

---

## 4. Variable Resolution Order

Within a single `Eval(vars)` call, variable lookup proceeds:

```
1. eval-level vars   (map[string]any passed to CompiledExpr.Eval or Env.Eval)
2. env globals       (registered via WithGlobals on this Env or any ancestor)
3. <undefined>       → runtime error: "undefined variable: $<name>"
```

Eval-level vars always shadow env globals with the same name. This is by design — per-call data is more specific than env-wide defaults.

Globals from child `Env` (via `Extend`) shadow globals from the parent, applying the same merge logic as functions and pipe handlers.

---

## 5. Goroutine Safety

| Type | Safe to share across goroutines? | Notes |
|---|---|---|
| `*Env` | **Yes** | All fields frozen after construction |
| `*CompiledExpr` | **Yes** | `ByteCode` immutable; `pool` thread-safe; `Eval` borrows VM per call |
| `vm.VMFunctions` | **Yes** | Read-only map after registration |
| `vm.PipeHandlers` | **Yes** | Read-only map after registration |
| `*vm.VM` (internal) | **No** | Mutable execution state; managed by pool, never exposed |

---

## 6. File Layout

```
uexl-go/
├── uexl.go        — Eval(), Validate(), MustCompile(), Default(), DefaultWith(),
│                    NewEnv(), Option, type aliases (Functions, PipeHandler, PipeHandlers,
│                    PipeContext, ParserError, ParseErrors),
│                    WithFunctions, WithPipeHandlers, WithGlobals, WithLib
├── env.go         — Env struct, NewEnv impl, Extend, Compile, MustCompile,
│                    Validate, Eval, HasFunction, HasPipe, HasGlobal, Info
├── env_config.go  — envConfig (unexported), EnvConfig (public projection), Lib interface
├── env_info.go    — EnvInfo struct and String() method
├── compiled.go    — CompiledExpr struct, Eval, Variables, Env methods
├── result.go      — AsFloat64, AsBool, AsString, AsSlice, AsMap helpers
└── doc.go         — Package-level godoc
```

---

## 7. Exact Signatures

```go
// uexl.go

package uexl

import (
    "github.com/maniartech/uexl/vm"
    parsererrors "github.com/maniartech/uexl/parser/errors"
)

// Type aliases — users import only "github.com/maniartech/uexl"
type Function     = vm.VMFunction      // func(args ...any) (any, error)
type Functions    = vm.VMFunctions     // map[string]Function
type PipeHandler  = vm.PipeHandler     // func(ctx PipeContext, input any) (any, error)
type PipeHandlers = vm.PipeHandlers
type PipeContext  = vm.PipeContext     // see §3.25

// Parser error types re-exported so callers never import parser/errors
type ParserError = parsererrors.ParserError   // single structured parse error (value type)
type ParseErrors = parsererrors.ParseErrors   // collection of parse errors (value type)

type Option func(*envConfig)

// Options
func WithFunctions(fns Functions) Option
func WithPipeHandlers(pipes PipeHandlers) Option
func WithGlobals(vars map[string]any) Option
func WithLib(lib Lib) Option

// Environment constructors
func Default() *Env
func DefaultWith(opts ...Option) *Env
func NewEnv(opts ...Option) *Env

// One-shot helpers (use Default env + context.Background())
func Eval(expr string, vars map[string]any) (any, error)
func Validate(expr string) error
func MustCompile(expr string) *CompiledExpr
```

```go
// result.go

package uexl

func AsFloat64(v any) (float64, error)
func AsBool(v any)    (bool, error)
func AsString(v any)  (string, error)
func AsSlice(v any)   ([]any, error)
func AsMap(v any)     (map[string]any, error)
```

```go
// env_config.go

// Lib is implemented by packages that ship reusable bundles of UExL extensions.
type Lib interface {
    Apply(cfg *EnvConfig)
}

// EnvConfig is the public projection of envConfig, accessible only inside Lib.Apply.
// Exposes additive operations only — no replacement, no read-back.
type EnvConfig struct {
    cfg *envConfig // unexported
}

func (c *EnvConfig) AddFunctions(fns Functions)
func (c *EnvConfig) AddPipeHandlers(pipes PipeHandlers)
func (c *EnvConfig) AddGlobals(vars map[string]any)

// unexported — internal accumulation state only
type envConfig struct {
    functions    vm.VMFunctions
    pipeHandlers vm.PipeHandlers
    globals      map[string]any
}
```

```go
// env.go

import "sync"

type Env struct {
    functions    vm.VMFunctions   // frozen after construction
    pipeHandlers vm.PipeHandlers  // frozen after construction
    globals      map[string]any   // frozen after construction
    pool         sync.Pool        // per-Env VM pool
}

func (e *Env) Extend(opts ...Option) *Env
func (e *Env) Compile(expr string) (*CompiledExpr, error)  // validates fn names at compile time
func (e *Env) MustCompile(expr string) *CompiledExpr
func (e *Env) Validate(expr string) error
func (e *Env) Eval(ctx context.Context, expr string, vars map[string]any) (any, error)

// Introspection — all read-only, goroutine-safe, allocation-free
func (e *Env) Info() EnvInfo
func (e *Env) HasFunction(name string) bool
func (e *Env) HasPipe(name string) bool
func (e *Env) HasGlobal(name string) bool
```

```go
// env_info.go

type EnvInfo struct {
    Functions    []string // sorted function names
    PipeHandlers []string // sorted pipe handler names
    Globals      []string // sorted global variable names
}

// String implements fmt.Stringer for human-readable debug output.
func (i EnvInfo) String() string
```

```go
// compiled.go

import (
    "context"
    "github.com/maniartech/uexl/compiler"
)

type CompiledExpr struct {
    bytecode *compiler.ByteCode
    env      *Env
}

func (c *CompiledExpr) Eval(ctx context.Context, vars map[string]any) (any, error)
func (c *CompiledExpr) Variables() []string   // derived from bytecode.ContextVars; sorted copy
func (c *CompiledExpr) Env() *Env             // returns the Env used at compile time; no allocation
```

```go
// doc.go

// Package uexl provides a bytecode-compiled expression evaluation engine.
// ...
package uexl
```

---

## 8. Error Behavior

### 8.1 Error Type Taxonomy

Users can distinguish error categories using `errors.As`. The parser can return two types:

| Category | Concrete type | `errors.As` target | Import |
|---|---|---|---|
| Single parse error | `uexl.ParserError` (value) | `var pe uexl.ParserError` | `github.com/maniartech/uexl` |
| Multiple parse errors | `uexl.ParseErrors` (value) | `var pe uexl.ParseErrors` | `github.com/maniartech/uexl` |
| Compile error | `error` (plain) | n/a | — |
| Runtime error | `error` (plain) | n/a | — |

> **Note:** `ParserError` and `ParseErrors` implement `error` via **value receivers**, so `errors.As` targets must be value types, not pointers.

```go
result, err := env.Eval(ctx, expr, vars)
if err != nil {
    var single uexl.ParserError
    var multi  uexl.ParseErrors
    switch {
    case errors.As(err, &single):
        // single syntax error: single.Line, single.Column, single.Message
    case errors.As(err, &multi):
        // multiple syntax errors: multi.Errors []uexl.ParserError
    default:
        // compile-time or runtime error (plain error, check message)
    }
}
```

The `uexl` package **does not wrap errors** — all errors are propagated as-is from the subsystem that produced them, preserving `errors.As` / `errors.Is` chains.

### 8.2 Nil Argument Panics

| Call | Panic message |
|---|---|
| `WithFunctions(nil)` | `"uexl: WithFunctions: fns must not be nil"` |
| `WithPipeHandlers(nil)` | `"uexl: WithPipeHandlers: pipes must not be nil"` |
| `WithGlobals(nil)` | `"uexl: WithGlobals: vars must not be nil"` |
| `WithLib(nil)` | `"uexl: WithLib: lib must not be nil"` |
| `EnvConfig.AddFunctions(nil)` | `"uexl: EnvConfig.AddFunctions: fns must not be nil"` |
| `EnvConfig.AddPipeHandlers(nil)` | `"uexl: EnvConfig.AddPipeHandlers: pipes must not be nil"` |
| `EnvConfig.AddGlobals(nil)` | `"uexl: EnvConfig.AddGlobals: vars must not be nil"` |

These are **programmer errors** (wrong API usage detectable at startup), not runtime data errors. All other error conditions return `error` values.

---

## 9. Deprecated Symbol

`EvalExpr(expr string)` in the current `uexl.go` is the only existing export. After this API is implemented:

- `EvalExpr` is **removed** (it has no callers outside the repo — confirmed by grep).
- It is **not** preserved with a deprecation comment since the module has no tagged public release yet; no compatibility obligation exists.

---

## 10. Implementation Checklist

### Phase 1 — Core types, type aliases, zero-option path

- [ ] Create `env_config.go`: define `envConfig` (unexported), `EnvConfig` (public projection), `Lib` interface
- [ ] Implement `EnvConfig.AddFunctions`, `AddPipeHandlers`, `AddGlobals` with nil-guards
- [ ] Create `env.go`: define `Env` struct with unexported fields
- [ ] Implement `NewEnv(opts ...Option) *Env` (applies options to blank `envConfig`)
- [ ] Implement `Default()` with `sync.Once` singleton
- [ ] Create `env_info.go`: define `EnvInfo` struct, implement `String() string` (sorted, stable, multiline format)
- [ ] Implement `(*Env).Info() EnvInfo` — collect sorted keys from all three maps into independent slices
- [ ] Implement `(*Env).HasFunction(name string) bool`
- [ ] Implement `(*Env).HasPipe(name string) bool`
- [ ] Implement `(*Env).HasGlobal(name string) bool`
- [ ] Create `compiled.go`: define `CompiledExpr` struct
- [ ] Implement `(*Env).Compile(expr string) (*CompiledExpr, error)` (parse + compile only, no VM)
- [ ] Implement `(*CompiledExpr).Eval(ctx context.Context, vars map[string]any) (any, error)` with `sync.Pool`
- [ ] Implement `(*CompiledExpr).Variables() []string` — sorted copy of `bytecode.ContextVars`
- [ ] Implement `(*CompiledExpr).Env() *Env` — returns stored `env` pointer, no allocation
- [ ] Implement `(*Env).Eval(ctx context.Context, expr string, vars map[string]any) (any, error)` as thin wrapper over Compile+Eval
- [ ] Implement `(*Env).Validate(expr string) error` as thin wrapper (`_, err := e.Compile(expr); return err`)
- [ ] Implement package-level `Validate(expr string) error` using `Default()`
- [ ] Create `result.go`: implement `AsFloat64`, `AsBool`, `AsString`, `AsSlice`, `AsMap` with nil-guards and numeric widening; no panics
- [ ] Rewrite `uexl.go`: declare `Function`, `Functions`, `PipeHandler`, `PipeHandlers`, `PipeContext`, `ParserError`, `ParseErrors` type aliases; add `Eval`, `Validate`, `WithFunctions`, `WithPipeHandlers`, `WithGlobals`, `WithLib`; remove `EvalExpr`
- [ ] Create `doc.go` with package-level godoc

### Phase 2 — Extend, MustCompile, DefaultWith, WithLib, globals

- [ ] Implement `(*Env).Extend(opts ...Option) *Env` (copy-on-extend; new pool per child)
- [ ] Implement `(*Env).MustCompile(expr string) *CompiledExpr` (panic with message `"uexl: Env.MustCompile: <err>"`)
- [ ] Implement `DefaultWith(opts ...Option) *Env` as `Default().Extend(opts...)`
- [ ] Implement package-level `MustCompile(expr string) *CompiledExpr` using `Default()`
- [ ] Implement `WithLib(lib Lib) Option` — calls `lib.Apply(&EnvConfig{cfg})` during option application
- [ ] Implement globals variable merging in `(*CompiledExpr).Eval` (fast path: skip allocation when `env.globals` is empty)
- [ ] Implement nil-guard panics in `WithLib` and all `EnvConfig.Add*` methods

### Phase 3 — Tests (file: `uexl_test.go`)

**Basic evaluation**
- [ ] `TestEval_basic` — package-level `Eval` with vars
- [ ] `TestEval_noVars` — `vars` is `nil`
- [ ] `TestEval_parseError` — malformed expression returns error, not panic
- [ ] `TestEval_errorIsParserError` — `errors.As(err, &pe)` with `var pe uexl.ParserError` succeeds on parse failure

**Default and NewEnv**
- [ ] `TestDefault_hasBuiltins` — `Default().Eval("len('hi')", nil)` returns `2.0`
- [ ] `TestDefault_singleton` — two calls to `Default()` return the same pointer
- [ ] `TestNewEnv_blankSlate` — bare `NewEnv()` has no builtins; calling a builtin returns error
- [ ] `TestDefaultWith_extendsDefault` — `DefaultWith(WithFunctions(...))` has stdlib AND custom fns
- [ ] `TestDefaultWith_doesNotMutateDefault` — `Default()` before and after `DefaultWith` have same Info

**Compile and CompiledExpr**
- [ ] `TestEnv_Compile_andEval` — compile once, eval with two different var maps
- [ ] `TestEnv_Compile_parseError` — returns error
- [ ] `TestMustCompile_packageLevel` — valid expr returns `*CompiledExpr`
- [ ] `TestMustCompile_panicsOnBadExpr` — invalid expr panics with `"uexl: MustCompile:"` prefix
- [ ] `TestEnv_MustCompile_panicsOnBadExpr` — same for `(*Env).MustCompile`
- [ ] `TestCompiledExpr_Variables_basic` — `"price * qty"` returns `["price", "qty"]` (sorted)
- [ ] `TestCompiledExpr_Variables_noVars` — `"1 + 2"` returns `[]string{}`
- [ ] `TestCompiledExpr_Variables_independent` — mutating returned slice does not affect next call
- [ ] `TestCompiledExpr_Env` — `Env()` returns the env it was compiled with

**Extend and multi-level**
- [ ] `TestEnv_Extend_inherits` — child has parent functions
- [ ] `TestEnv_Extend_override` — child overrides a function; parent unaffected
- [ ] `TestEnv_Extend_additionalPipes` — child adds a custom pipe handler
- [ ] `TestEnv_Extend_multiLevel` — three-level chain; each level has correct symbol set
- [ ] `TestEnv_Extend_parentUnchanged` — child `Extend` never modifies parent's functions map

**WithLib**
- [ ] `TestWithLib_appliesFunctionsAndPipes` — lib registers correctly
- [ ] `TestWithLib_overriddenByLaterOption` — `WithFunctions` after `WithLib` for same name wins
- [ ] `TestWithLib_nil_panics`

**Globals**
- [ ] `TestWithGlobals_shadowedByVars` — eval var shadows global with same name
- [ ] `TestWithGlobals_nilValueShadows` — eval var `nil` correctly shadows global, not skipped
- [ ] `TestWithGlobals_usedWhenNoVar` — global used when eval vars don't contain the key
- [ ] `TestWithGlobals_inheritedByExtend` — child env inherits parent globals

**Introspection**
- [ ] `TestEnv_HasFunction_true`, `_false`, `_emptyString`
- [ ] `TestEnv_HasPipe_true`, `_false`, `_emptyString`
- [ ] `TestEnv_HasGlobal_true`, `_false`, `_emptyString`
- [ ] `TestEnv_HasFunction_inheritedFromParent`
- [ ] `TestEnvInfo_sorted` — all three slices are sorted
- [ ] `TestEnvInfo_stable` — two calls return equal slices
- [ ] `TestEnvInfo_independent` — mutating slices does not affect Env
- [ ] `TestEnvInfo_String_format` — output matches documented multiline format exactly

**Nil-guard panics**
- [ ] `TestWithFunctions_nil_panics`
- [ ] `TestWithPipeHandlers_nil_panics`
- [ ] `TestWithGlobals_nil_panics`
- [ ] `TestWithLib_nil_panics`
- [ ] `TestEnvConfig_AddFunctions_nil_panics`
- [ ] `TestEnvConfig_AddPipeHandlers_nil_panics`
- [ ] `TestEnvConfig_AddGlobals_nil_panics`

**Concurrency (run with `-race`)**
- [ ] `TestCompiledExpr_concurrentEval` — 50 goroutines, one `*CompiledExpr`
- [ ] `TestEnv_concurrentEval` — 50 goroutines calling `env.Eval`
- [ ] `TestEnv_concurrentInfo` — 50 goroutines calling `env.Info()`
- [ ] `TestEnv_concurrentExtend` — 50 goroutines each calling `Extend` on same parent
- [ ] `TestEnv_concurrentValidate` — 50 goroutines calling `env.Validate` on same Env (valid and invalid exprs)

**Validate**
- [ ] `TestValidate_packageLevel_valid` — `Validate("1 + 2")` returns `nil`
- [ ] `TestValidate_packageLevel_parseError` — `Validate("1 +")` returns non-nil error
- [ ] `TestEnv_Validate_valid` — `env.Validate("price * qty")` returns `nil` (no VM, vars unused)
- [ ] `TestEnv_Validate_parseError` — returns error, not panic
- [ ] `TestEnv_Validate_noArtifact` — confirm no `*CompiledExpr` is returned (API-level: signature is `error` only)

**Result coercion helpers (file: `result_test.go`)**
- [ ] `TestAsFloat64_float64` — `AsFloat64(3.14)` → `3.14, nil`
- [ ] `TestAsFloat64_int` — `AsFloat64(7)` → `7.0, nil`
- [ ] `TestAsFloat64_int64` — `AsFloat64(int64(7))` → `7.0, nil`
- [ ] `TestAsFloat64_float32` — `AsFloat64(float32(1.5))` → `1.5, nil`
- [ ] `TestAsFloat64_string` — `AsFloat64("x")` → error
- [ ] `TestAsFloat64_nil` — `AsFloat64(nil)` → error (nil has no numeric meaning)
- [ ] `TestAsBool_true` — `AsBool(true)` → `true, nil`
- [ ] `TestAsBool_false` — `AsBool(false)` → `false, nil`
- [ ] `TestAsBool_int` — `AsBool(1)` → error (no truthy coercion)
- [ ] `TestAsBool_nil` — `AsBool(nil)` → error
- [ ] `TestAsString_string` — `AsString("hello")` → `"hello", nil`
- [ ] `TestAsString_int` — `AsString(42)` → error (no fmt.Sprint fallback)
- [ ] `TestAsString_nil` — `AsString(nil)` → error
- [ ] `TestAsSlice_slice` — `AsSlice([]any{1, 2})` → `[]any{1,2}, nil`
- [ ] `TestAsSlice_nil` — `AsSlice(nil)` → error
- [ ] `TestAsMap_map` — `AsMap(map[string]any{"k":1})` → `map, nil`
- [ ] `TestAsMap_nil` — `AsMap(nil)` → error
- [ ] `TestAsFloat64_roundtrip` — `Eval("price * qty", vars)` → `AsFloat64` → expected value

### Phase 4 — Documentation update

- [ ] Update `book/golang/overview.md` with the real API surface (replacing stubs)
- [ ] Update `README.md` usage examples to use `uexl.Eval`, `uexl.DefaultWith(...)`, and `WithLib`
- [ ] Update Launch Checklist §2 to mark items complete

---

## 11. Non-Goals (explicitly out of scope for this spec)

- **Full compile-time type inference** — function argument types and return types are not checked at compile time. UExL is dynamically typed by design. Only function *existence* is validated at `Compile` time (see §3.13).
- **Expression caching / memoization** inside `Env.Eval` — callers who want this should call `Compile` explicitly.
- **`context.Context` in `Compile`** — `Compile` is a pure CPU-bound operation with no I/O; it does not accept a context. Only `Eval` (which involves the VM loop) accept context.
- **Thread-local VM pools** — `sync.Pool` is sufficient; NUMA-aware pooling is premature.
- **Streaming / async evaluation** — not in scope for v0.1.0.
- **Serialization of `*CompiledExpr`** — bytecode is not designed for wire transmission in this release.
- **Hot-reload / mutation of `Env` after construction** — explicitly excluded to preserve goroutine safety.
- **Introspection of function signatures** — `HasFunction` reports presence only, not arity or parameter types.
- **`EnvInfo` diffing helpers** — callers can diff two `[]string` slices themselves.
- **Dynamic lib loading** (plugins, `.so`) — `Lib` is a static Go interface; dynamic loading is out of scope.
- **`PipeContext.EvalItem` arity checking** — `EvalItem` sets `$item` and `$index` but not custom scope vars (use `EvalWith` for those). The interface is intentionally narrow to keep handler code readable.

---

## 12. Competitive Analysis

This section benchmarks the UExL public API design against three major Go expression-evaluation frameworks: **cel-go**, **expr (expr-lang)**, and **gval**. UExL sits in the same niche as these frameworks — safe embedded expression evaluation — but has a distinct identity through its pipe operator and explicit nullish semantics.

### 12.1 Framework Profiles

| | **UExL** | **cel-go** | **expr** | **gval** |
|---|---|---|---|---|
| **Module** | `github.com/maniartech/uexl` | `github.com/google/cel-go` | `github.com/expr-lang/expr` | `github.com/PaesslerAG/gval` |
| **Backend** | Bytecode VM | Bytecode (native CEL) | Bytecode VM | DSL-based evaluable |
| **Type system** | Dynamic (runtime) | Static (compile-time) | Semi-static (struct env) | Dynamic (runtime) |
| **Primary use case** | Business rules + data transformation | Policy evaluation (Kubernetes) | Template / CEL-lite | Formula evaluation |

### 12.2 Feature Comparison

| Feature | UExL | cel-go | expr | gval |
|---|---|---|---|---|
| **Compile-once / run-many** | ✅ `*CompiledExpr` | ✅ `cel.Program` | ✅ `*vm.Program` | ✅ `Evaluable` (function) |
| **Layered env / inheritance** | ✅ `Extend()` chain | ✅ `cel.Library` | ❌ global only | ✅ `NewLanguage(pieces...)` |
| **Lib bundling interface** | ✅ `Lib.Apply(*EnvConfig)` | ✅ `cel.Library` (2 methods) | ❌ | ❌ |
| **Env-level globals** | ✅ `WithGlobals` | ❌ | ❌ | ❌ |
| **`MustCompile` variant** | ✅ | ❌ | ❌ | ❌ |
| **`DefaultWith` shorthand** | ✅ | ❌ | ❌ | ✅ `gval.Base()` |
| **`Validate` (no artifact)** | ✅ (syntax+compile only) | ✅ `env.Parse()`+`env.Check()` | ❌ (Compile IS validate) | ❌ |
| **`Variables()` extraction** | ✅ sorted `[]string` | ❌ | ✅ (via `expr.Compile` options) | ❌ |
| **Introspection (`HasFunction`)** | ✅ | ❌ | ❌ | ❌ |
| **`EnvInfo` listing** | ✅ `Info() EnvInfo` | ❌ | ❌ | ❌ |
| **Result type helpers** | ✅ `AsFloat64` etc. | ✅ `ref.Val` interface | ❌ native Go values | ❌ native Go values |
| **Pipe / transform operators** | ✅ `\|map:` `\|filter:` etc. | ❌ | ❌ | ❌ |
| **`context.Context` in eval** | ✅ `Eval(ctx, vars)` | ❌ | ❌ | ✅ `evaluable(ctx, vars)` |
| **Compile-time fn existence check** | ✅ (§3.13) | ✅ | ✅ | ❌ |
| **Compile-time type checking** | ❌ (dynamic by design) | ✅ | ✅ (with typed env) | ❌ |
| **No import of internal packages** | ✅ type aliases | ✅ | ✅ | ✅ |

### 12.3 UExL Strengths

**Pipe operators — unique in the Go ecosystem.** No other Go expression evaluator provides `|map:`, `|filter:`, `|reduce:`, `|sort:`, `|groupBy:` etc. as first-class language syntax. This makes UExL the only framework suitable for expression-driven ETL/data transformation and spreadsheet-like formulas without bolting on external libraries.

**Richer environment ergonomics than peers.**
- `WithGlobals` — zero peers provide env-level constant bindings without a custom `Eval` wrapper.
- `Extend()` chain — simpler and more composable than cel-go's `cel.Library` (which mandates two interface methods: `CompileOptions()` + `ProgramOptions()`).
- `Lib.Apply(*EnvConfig)` — single-method interface, lower burden than cel-go.
- `DefaultWith(opts...)` — avoids the boilerplate of creating a blank env just to add one function.

**Best-in-class introspection.** `HasFunction`, `HasPipe`, `HasGlobal`, and `EnvInfo.String()` are absent from all three peers. This is essential for diagnostic tooling, REPL tab-completion, and runtime library auditing.

**`MustCompile` and `Validate`.** No peer provides `MustCompile`; it allows clean top-level variable initialisation. `Validate` (syntax + compile + fn-existence check, no artifact) fills the gap between "just try to compile" and full type-checked static analysis.

**`context.Context` in `Run` with novice-friendly escape hatch.** No peer (cel-go, expr) supports ctx in eval. UExL matches gval's cancellation support while keeping the package-level `Eval(expr, vars)` ctx-free for novices.

**Result coercion helpers.** cel-go provides `ref.Val` (a typed wrapper returned instead of `any`); expr and gval return native Go values directly (reducing the need for helpers but requiring callers to write their own type assertions). UExL returns `any` like expr/gval but provides opt-in helpers, giving callers a safety net without mandating a wrapper type.

### 12.4 UExL Limitations vs Peers

~~**No compile-time function existence check.**~~ **Resolved in this spec (§3.13).** `Compile` now validates function names against the env's registry. Unknown function calls are caught at compile time, not runtime.

**No compile-time type inference.** UExL is dynamically typed (like gval). cel-go and expr with typed environments provide full static type checking over expressions. This is deliberate for UExL: dynamic expressions are a feature, not a bug, in rule-engine and configuration contexts.

~~**No `context.Context`.**~~ **Resolved in this spec (§3.15).** `CompiledExpr.Eval(ctx, vars)` and `Env.Eval(ctx, ...)` accept a `context.Context`. The VM checks for cancellation between opcode executions. The package-level `Eval(expr, vars)` convenience wrapper uses `context.Background()` for novice ergonomics.

~~**`PipeHandler` leaks `*vm.VM`.**~~ **Resolved in this spec (§3.25).** The new `PipeContext` interface replaces `*vm.VM`. Custom pipe authors import only `github.com/maniartech/uexl`; no internal packages needed.

### 12.5 Design Decision Rationale

| Decision | Rationale |
|---|---|
| `Env` as value receiver pattern (not functional options on a builder) | Simpler API footprint; aligns with `sync.Pool` ownership; matches `Extend` composability |
| Type aliases for `Functions`, `PipeHandler`, `PipeHandlers`, `PipeContext`, `ParserError`, `ParseErrors` | Users import only `uexl`; internal refactoring doesn't break import paths |
| `Default()` singleton via `sync.Once` | Zero-boilerplate one-liner evals; peers (expr, gval) have similar one-shot helpers |
| No truthy coercion in `AsBool` | Preserves UExL's explicit boolish semantics (see design-philosophy.md); prevents silent bugs in boolean-heavy rule engines |
| `Variables()` returns sorted `[]string` | Deterministic output for tests and documentation; sorted copy prevents aliasing |
| `MustCompile` panics with `"uexl: MustCompile: <err>"` prefix | Distinguishes init-time bugs from runtime errors; consistent prefix aids `grep` in post-mortem logs |
