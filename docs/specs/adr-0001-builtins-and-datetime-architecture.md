# ADR-0001 — Standard-Library Architecture, Conversion Naming, and Clock Injection

Status: Proposed (awaiting ratification)
Date: 2026-06-14
Related: [datetime-spec.md](datetime-spec.md), [builtins-implementation-plan.md](builtins-implementation-plan.md),
[standard-library-roadmap.md](standard-library-roadmap.md), [language-spec-requirements.md](language-spec-requirements.md)

## Context

UExL is a cross-platform expression language (Go reference implementation; planned JS/Rust/Dart ports).
Three coupled decisions were open before datetime implementation could begin:

- **A.** Is datetime mandatory *core*, or an optional module — and *how* is the optional standard
  library physically attached (separate libraries? separate binaries? build tags?)?
- **B.** How are the value-conversion functions named, and what is their failure contract?
- **C.** How do `now()` / `today()` obtain the current instant without breaking evaluation purity?

These were grounded in the actual engine, not chosen abstractly. The key code facts (verified):

- **Literal prefixes are hard-coded in the tokenizer** — `parser/tokenizer.go` dispatches the raw-string
  prefix with `if t.input[t.pos] == 'r'`. Adding `d"..."` / `7d` literals means editing the tokenizer.
- **Operators are a hard-coded opcode switch** — `vm/vm_handlers.go` `executeBinaryExpression` (and the
  dispatch in `vm/vm.go`) switch on `code.OpAdd`/`OpSub`/… with type-specific handlers. New operator
  semantics (date − date → duration) mean editing the VM.
- **A library-composition mechanism already exists** — `type Lib interface { Apply(*EnvConfig) }`
  (`env_config.go`), applied via the `WithLib` / `WithFunctions` options (`uexl.go`) at `Env`
  construction. `EnvConfig` exposes **additive operations only**: `AddFunctions`, `AddPipeHandlers`,
  `AddGlobals`. It **cannot** register new literal syntax or operators.

## Decision A — Layered core + runtime-composed libraries

**The rule:** *if the parser or VM must understand it, it is **core**; if it is just a function call, it
is an attachable **library**.* The code enforces this line — a `Lib` can add functions/pipes/globals but
not syntax or operators.

**Core (always compiled in):**

- `datetime` / `duration` value types
- the `d"..."` literal, **including its fixed ISO 8601 parser** (see sub-decision below)
- duration suffix literals (`7d`, `30ms`)
- the temporal operators (`date − date`, `date ± duration`, duration arithmetic, comparisons)
- `typeof` and a small set of irreducible builtins

**Attachable libraries (opt-in `Lib` bundles, composed via `WithLib`):**

- **`datetime`** — every datetime/duration *function* (parse/format/calendar/construct/extract/epoch)
- **`math`**, **`conversion`**, **`collections`/complex-logic**, etc.

```go
// general-purpose host
env := uexl.DefaultWith(
    uexl.WithLib(datetimex.Lib()),
    uexl.WithLib(mathx.Lib()),
)
// minimal host
env := uexl.Default()   // none of the above attached
```

**Packaging answer: not separate binaries, not build tags.** Attachable libraries register at runtime
through the existing `Lib` / `WithLib` API. A library MAY be a separate Go *package* (per
[builtins-implementation-plan.md](builtins-implementation-plan.md)) — or a separate Go *module* for
independent versioning — but composition happens at `Env` construction, in one binary.

**Consequence (worth ratifying explicitly):** because the *type*, *literal*, and *operators* are core,
the datetime/duration **type is effectively always present** — every host can write `d"2026-12-31"`,
compare, subtract, and add durations out of the box. Only the **function library** is optional (though
*recommended* for general-purpose hosts). This is a *stronger* portability guarantee than a fully
optional profile, and is safe to make core because UExL's core datetime layer has **zero host-dependent
behavior** (zoneless UTC milliseconds; no timezone database). That is a deliberate improvement over CEL,
whose timezone-aware accessors reintroduce host-database dependence.

**Sub-decision (Option 1):** the `d"..."` literal parses the **fixed** ISO 8601 grammar in core; **all**
runtime `parseDate` / `formatDate` / NITES-format functions live in the `datetime` library. The ISO
parser code is shared internally. Net: bare core can *write* date literals; parsing dynamic strings,
formatting, and calendar math require attaching the `datetime` library.

## Decision B — Conversion naming: `parseX` (strict) / `tryParseX` (safe) + `str`

Functions that turn a string into a typed value come as a pair distinguished by one prefix:

- **`parseX`** — strict; **raises an error** on bad input.
- **`tryParseX`** — safe; **returns `null`** on bad input.
- **`str(x)`** — value → string; never fails, so it needs no twin.

Type tokens are the short canonical forms: **`Num`, `Bool`, `Date`, `Dur`**.

| Intent | Strict (error) | Safe (null) |
|--------|----------------|-------------|
| string → number | `parseNum` | `tryParseNum` |
| string → boolean | `parseBool` | `tryParseBool` |
| string → datetime | `parseDate` | `tryParseDate` |
| string → duration | `parseDur` | `tryParseDur` |
| any → string | `str` | — |

This **supersedes** the earlier `to*` (lenient) / `parse*` (strict) split. `tryParse*` states the
safe-failure contract more explicitly and shares the `parse` root, so a strict/lenient pair reads as two
halves of one operation. Dropped: `toNumber`, `toBoolean`, `toString`, `parseNumber`, `parseBoolean`.
Renamed for token consistency: `parseDuration` → `parseDur`, `formatDuration` → `formatDur`.

## Decision C — `now()` / `today()` read a single frozen instant

The current instant is captured **once at the start of evaluation** (the host snaps the system clock) and
injected as context state (canonical milliseconds). `now()` reads that value; `today()` is the same
instant truncated to UTC midnight. Every `now()` / `today()` call anywhere in the expression — **including
across pipe stages** — returns the identical instant, so `now() == now()` and `now() - now() == 0`.

- **Production:** host snaps the real system clock once at eval start.
- **Tests:** host injects a fixed instant → fully reproducible.
- **No clock injected:** behavior is host-defined (error or wall clock) and outside the deterministic
  conformance guarantee.

**Mechanism:** reuse the existing context-variable plumbing — no change to the builtin call ABI. The
normative statement lives in [datetime-spec.md](datetime-spec.md) §9.1.

## Alternatives considered

- **A:** mandatory core (CEL's choice) — rejected: burdens minimal/edge hosts with a calendar engine.
  Host-native delegation (Expr's choice) — rejected: cannot be specified once and ported. Build tags /
  separate binaries — rejected: unnecessary given the runtime `Lib` mechanism already in the engine.
- **B:** `str`/`num`/`bool` — rejected: only one failure mode per function, and clashes with the
  already-committed `parseDate`. `to*`/`parse*` — rejected in favor of `parse*`/`tryParse*` (clearer
  failure contract, shared root).
- **C:** change the builtin call ABI to pass a clock — rejected: large refactor for a one-function need.
  VM special-case of `now()` — rejected: hard-coded magic. Host-supplied `now()` — rejected: breaks
  determinism and portability.

## Consequences

- [datetime-spec.md](datetime-spec.md): scope header updated; new core-vs-`datetime`-library boundary
  (§1.2); §9.1 extended for `today()` stability and the context-injection mechanism; conversion-naming
  references updated to `parseX`/`tryParseX`; duration functions renamed `parseDur`/`formatDur`.
- [builtins-implementation-plan.md](builtins-implementation-plan.md): families expose `Lib` bundles; the
  core-vs-attachable split is documented; conversion-family naming updated.
- The datetime implementation track is unblocked once this ADR is ratified.
