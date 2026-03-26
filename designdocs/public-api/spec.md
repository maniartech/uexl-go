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

### 2.2 Functions

```
Default()  *Env
NewEnv(opts ...Option) *Env
Eval(expr string, vars map[string]any) (any, error)
WithFunctions(fns vm.VMFunctions) Option
WithPipeHandlers(pipes vm.PipeHandlers) Option
WithGlobals(vars map[string]any) Option
WithLib(lib Lib) Option
```

### 2.3 Methods on `*Env`

```
(e *Env) Extend(opts ...Option) *Env
(e *Env) Compile(expr string) (*CompiledExpr, error)
(e *Env) Eval(expr string, vars map[string]any) (any, error)

// Introspection
(e *Env) Info() EnvInfo
(e *Env) HasFunction(name string) bool
(e *Env) HasPipe(name string) bool
(e *Env) HasGlobal(name string) bool
```

### 2.4 `Lib` interface

```go
type Lib interface {
    // Apply registers the lib's functions, pipe handlers, and globals
    // into the supplied config. Called once during Env construction.
    Apply(cfg *envConfig)
}
```

### 2.5 `EnvInfo` (value type, read-only snapshot)

```
type EnvInfo struct {
    Functions    []string  // sorted function names
    PipeHandlers []string  // sorted pipe handler names
    Globals      []string  // sorted global variable names
}
```

### 2.6 Methods on `*CompiledExpr`

```
(c *CompiledExpr) Eval(vars map[string]any) (any, error)
```

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

### 3.5 `WithFunctions(fns vm.VMFunctions) Option`

Returns an `Option` that merges `fns` into the in-progress `envConfig.functions` map. Existing keys from a prior option are overwritten if `fns` contains the same key.

**Panics if `fns` is nil.** (Validated at option application time, not deferred to use, so misuse surfaces immediately at app startup.)

---

### 3.6 `WithPipeHandlers(pipes vm.PipeHandlers) Option`

Returns an `Option` that merges `pipes` into `envConfig.pipeHandlers`. Same semantics as `WithFunctions`.

**Panics if `pipes` is nil.**

---

### 3.7 `WithGlobals(vars map[string]any) Option`

Returns an `Option` that merges `vars` into `envConfig.globals`. These become env-level context variables, resolved **before** per-call vars in the resolution order (see §4).

**Panics if `vars` is nil.**

---

### 3.11 `WithLib(lib Lib) Option`

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

func (c *EnvConfig) AddFunctions(fns vm.VMFunctions)
func (c *EnvConfig) AddPipeHandlers(pipes vm.PipeHandlers)
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
    return Default().Eval(expr, vars)
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
3. Wrap `*compiler.ByteCode` and a reference to the parent `*Env` into a `*CompiledExpr`.
4. Return `*CompiledExpr, nil`.

`*CompiledExpr` is goroutine-safe. The `*compiler.ByteCode` it holds is immutable after compilation. Multiple goroutines may call `compiledExpr.Eval(...)` concurrently.

---

### 3.14 `(*Env).Eval(expr string, vars map[string]any) (any, error)`

One-shot parse + compile + run within the environment. Convenience wrapper; equivalent to:

```go
func (e *Env) Eval(expr string, vars map[string]any) (any, error) {
    c, err := e.Compile(expr)
    if err != nil {
        return nil, err
    }
    return c.Eval(vars)
}
```

Use this for expressions evaluated once. Use `Compile` + `CompiledExpr.Eval` for expressions evaluated many times.

---

### 3.15 `(*CompiledExpr).Eval(vars map[string]any) (any, error)`

Evaluates the pre-compiled bytecode against `vars`. This is the **hot path**.

Steps:
1. Borrow a `*vm.VM` from `env.pool` — no allocation if pool has a spare.
2. Build the merged variable map (see §4 for resolution order).
3. Call `machine.Run(bytecode, mergedVars)`.
4. Return the VM to pool via `defer pool.Put(machine)`.
5. Return the result or error.

**Variable merging (step 2):** Produce a single `map[string]any` where eval-level vars shadow env globals. For performance, if `env.globals` is empty (no globals registered), pass `vars` directly without creating a new map. If globals are non-empty, allocate a merged map: start from a copy of `env.globals`, then apply `vars` on top.

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
| `*CompiledExpr` | **Yes** | `ByteCode` immutable; `pool` thread-safe |
| `vm.VMFunctions` | **Yes** | Read-only map after registration |
| `vm.PipeHandlers` | **Yes** | Read-only map after registration |
| `*vm.VM` (internal) | **No** | Mutable execution state; managed by pool, never exposed |

---

## 6. File Layout

```
uexl-go/
├── uexl.go        — Eval(), Default(), NewEnv(), Option, WithFunctions,
│                    WithPipeHandlers, WithGlobals, WithLib (replaces current stub)
├── env.go         — Env struct, NewEnv impl, Extend, Compile, Eval methods,
│                    HasFunction, HasPipe, HasGlobal, Info
├── env_config.go  — envConfig (unexported), EnvConfig (public projection for Lib.Apply),
│                    Lib interface, Option type alias
├── env_info.go    — EnvInfo struct and String() method
├── compiled.go    — CompiledExpr struct and Eval method
└── doc.go         — Package-level godoc
```

`envConfig` is unexported and lives in `env_config.go`. `EnvConfig` is its public projection, only accessible inside `Lib.Apply`; its pointer is never stored or returned externally.

---

## 7. Exact Signatures

```go
// uexl.go

package uexl

import (
    "github.com/maniartech/uexl/vm"
)

type Option func(*envConfig)

func WithFunctions(fns vm.VMFunctions) Option
func WithPipeHandlers(pipes vm.PipeHandlers) Option
func WithGlobals(vars map[string]any) Option
func WithLib(lib Lib) Option

func Default() *Env
func NewEnv(opts ...Option) *Env
func Eval(expr string, vars map[string]any) (any, error)
```

```go
// env_config.go

// Lib is implemented by packages that ship reusable bundles of UExL extensions.
type Lib interface {
    Apply(cfg *EnvConfig)
}

// EnvConfig is the public projection of envConfig, accessible only inside Lib.Apply.
// It exposes additive operations only — no replacement, no read-back.
type EnvConfig struct {
    cfg *envConfig // unexported
}

func (c *EnvConfig) AddFunctions(fns vm.VMFunctions)
func (c *EnvConfig) AddPipeHandlers(pipes vm.PipeHandlers)
func (c *EnvConfig) AddGlobals(vars map[string]any)

// unexported — internal accumulation state
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
func (e *Env) Compile(expr string) (*CompiledExpr, error)
func (e *Env) Eval(expr string, vars map[string]any) (any, error)

// Introspection — all read-only, goroutine-safe
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

import "github.com/maniartech/uexl/compiler"

type CompiledExpr struct {
    bytecode *compiler.ByteCode
    env      *Env
}

func (c *CompiledExpr) Eval(vars map[string]any) (any, error)
```

```go
// doc.go

// Package uexl provides a bytecode-compiled expression evaluation engine.
// ...
package uexl
```

---

## 8. Error Behavior

| Situation | Returned error text (prefix) |
|---|---|
| Parse failure | propagated from `parser.ParseString` (structured `errors.ParserError`) |
| Compile failure | propagated from `compiler.Compile` |
| Undefined variable at runtime | propagated from VM opcode handler |
| Type mismatch at runtime | propagated from VM opcode handler |
| `WithFunctions(nil)` | `panic("uexl: WithFunctions: fns must not be nil")` |
| `WithPipeHandlers(nil)` | `panic("uexl: WithPipeHandlers: pipes must not be nil")` |
| `WithGlobals(nil)` | `panic("uexl: WithGlobals: vars must not be nil")` |
| `WithLib(nil)` | `panic("uexl: WithLib: lib must not be nil")` |
| `EnvConfig.AddFunctions(nil)` | `panic("uexl: EnvConfig.AddFunctions: fns must not be nil")` |
| `EnvConfig.AddPipeHandlers(nil)` | `panic("uexl: EnvConfig.AddPipeHandlers: pipes must not be nil")` |
| `EnvConfig.AddGlobals(nil)` | `panic("uexl: EnvConfig.AddGlobals: vars must not be nil")` |

The panics are **programmer errors** (wrong API usage at startup), not runtime data errors. All other error conditions return `error` values.

---

## 9. Deprecated Symbol

`EvalExpr(expr string)` in the current `uexl.go` is the only existing export. After this API is implemented:

- `EvalExpr` is **removed** (it has no callers outside the repo — confirmed by grep).
- It is **not** preserved with a deprecation comment since the module has no tagged public release yet; no compatibility obligation exists.

---

## 10. Implementation Checklist

### Phase 1 — Core types and zero-option path

- [ ] Create `env_config.go`: define `envConfig` (unexported), `EnvConfig` (public projection), `Lib` interface, `Option` type alias
- [ ] Implement `EnvConfig.AddFunctions`, `AddPipeHandlers`, `AddGlobals` with nil-guards
- [ ] Create `env.go`: define `Env` struct
- [ ] Implement `NewEnv(opts ...Option) *Env` (applies options to blank `envConfig`)
- [ ] Implement `Default()` with `sync.Once` singleton
- [ ] Create `env_info.go`: define `EnvInfo` struct, implement `String() string` with sorted, stable output
- [ ] Implement `(*Env).Info() EnvInfo` — collect sorted keys from all three maps
- [ ] Implement `(*Env).HasFunction(name string) bool`
- [ ] Implement `(*Env).HasPipe(name string) bool`
- [ ] Implement `(*Env).HasGlobal(name string) bool`
- [ ] Create `compiled.go`: define `CompiledExpr`
- [ ] Implement `(*Env).Compile(expr string) (*CompiledExpr, error)` (parse + compile only)
- [ ] Implement `(*CompiledExpr).Eval(vars map[string]any) (any, error)` with `sync.Pool`
- [ ] Implement `(*Env).Eval(expr string, vars map[string]any) (any, error)` as thin wrapper
- [ ] Rewrite `uexl.go`: add `Eval`, `WithFunctions`, `WithPipeHandlers`, `WithGlobals`, `WithLib`; remove `EvalExpr`
- [ ] Create `doc.go` with package godoc

### Phase 2 — Extend, WithLib, and globals

- [ ] Implement `(*Env).Extend(opts ...Option) *Env` (copy-on-extend)
- [ ] Implement `WithLib(lib Lib) Option` — calls `lib.Apply` during option application
- [ ] Implement globals variable merging in `(*CompiledExpr).Eval` (fast path: skip allocation if no globals)
- [ ] Implement nil-guard panics in `WithLib` and all `EnvConfig.Add*` methods

### Phase 3 — Tests (file: `uexl_test.go`)

- [ ] `TestEval_basic` — package-level `Eval` with vars, no config
- [ ] `TestEval_noVars` — `vars` is `nil`
- [ ] `TestEval_parseError` — malformed expression returns error, not panic
- [ ] `TestNewEnv_blankSlate` — bare `NewEnv()` has no builtins; calling a builtin returns error
- [ ] `TestDefault_hasBuiltins` — `Default().Eval("len('hi')", nil)` returns `2`
- [ ] `TestDefault_singleton` — two calls to `Default()` return the same pointer
- [ ] `TestEnv_Compile_andEval` — compile once, eval with two different var maps
- [ ] `TestEnv_Compile_parseError` — returns error
- [ ] `TestEnv_Extend_inherits` — child env has parent functions
- [ ] `TestEnv_Extend_override` — child overrides a function, parent unaffected
- [ ] `TestEnv_Extend_additionalPipes` — child adds a custom pipe handler
- [ ] `TestEnv_Extend_multiLevel` — three-level chain (stdlib → domain → tenant); each level only sees its own + inherited symbols
- [ ] `TestWithGlobals_shadowedByVars` — eval var shadows global with same name
- [ ] `TestWithGlobals_nilShadow` — eval var of `nil` correctly shadows global (not skipped)
- [ ] `TestWithGlobals_usedWhenNoVar` — global is used when eval vars do not contain the key
- [ ] `TestWithLib_appliesFunctionsAndPipes` — `WithLib` registers functions and pipes correctly
- [ ] `TestWithLib_overriddenByLaterOption` — `WithFunctions` after `WithLib` wins for same name
- [ ] `TestWithLib_nil_panics` — `WithLib(nil)` panics with correct message
- [ ] `TestEnvConfig_AddFunctions_nil_panics`
- [ ] `TestEnvConfig_AddPipeHandlers_nil_panics`
- [ ] `TestEnvConfig_AddGlobals_nil_panics`
- [ ] `TestEnv_HasFunction_true` — registered function returns true
- [ ] `TestEnv_HasFunction_false` — unregistered name returns false
- [ ] `TestEnv_HasFunction_emptyString` — empty string always returns false
- [ ] `TestEnv_HasPipe_true`, `_false`, `_emptyString`
- [ ] `TestEnv_HasGlobal_true`, `_false`, `_emptyString`
- [ ] `TestEnv_HasFunction_inheritedFromParent` — child env sees parent functions via `HasFunction`
- [ ] `TestEnvInfo_sorted` — `Info()` returns all three slices in sorted order
- [ ] `TestEnvInfo_stable` — two calls to `Info()` on same Env return equal slices
- [ ] `TestEnvInfo_independent` — mutating returned slices does not affect the Env
- [ ] `TestEnvInfo_String_format` — output matches the documented multiline format
- [ ] `TestCompiledExpr_concurrentEval` — 50 goroutines call `Eval` on one `*CompiledExpr`; no races
- [ ] `TestEnv_concurrentEval` — 50 goroutines call env `Eval` concurrently; no races
- [ ] `TestEnv_concurrentInfo` — 50 goroutines call `Info` on one `*Env` concurrently; no races
- [ ] `TestWithFunctions_nil_panics`
- [ ] `TestWithPipeHandlers_nil_panics`
- [ ] `TestWithGlobals_nil_panics`

### Phase 4 — Documentation update

- [ ] Update `book/golang/overview.md` with the real API surface (replacing stubs)
- [ ] Update `README.md` usage examples to use `uexl.Eval`, `uexl.Default().Extend(...)`, and `WithLib`
- [ ] Update Launch Checklist §2 to mark items complete

---

## 11. Non-Goals (explicitly out of scope for this spec)

- **Expression caching / memoization** inside `Env.Eval` — callers who want this should call `Compile` explicitly.
- **Thread-local VM pools** — `sync.Pool` is sufficient; NUMA-aware pooling is premature.
- **Streaming / async evaluation** — not in scope for v0.1.0.
- **Serialization of `CompiledExpr`** — bytecode is not designed for wire transmission in this release.
- **Hot-reload / mutation of `Env` after construction** — explicitly excluded to preserve goroutine safety.
- **Introspection of function signatures** — `HasFunction` reports presence only, not arity or parameter types. Full reflection is a v0.x+ concern.
- **`EnvInfo` diffing helpers** — callers can diff two `[]string` slices themselves; no built-in `Diff(a, b EnvInfo)` is provided.
- **Dynamic lib loading** (plugins, `.so`) — `Lib` is a static Go interface; dynamic loading is out of scope.
