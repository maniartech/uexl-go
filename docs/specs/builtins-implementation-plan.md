# UExL Standard Library — Implementation & Directory Plan

Status: Draft
Audience: Maintainers / implementers
Companion artifacts:
- `standard-library-roadmap.md` — what functions exist, drifted, or are missing (the backlog)
- `standard-library.json` — machine-readable manifest (tooling source of truth)
- `standard-library.md` — generated consistency matrix + family-completeness report
- `cmd/builtins-check` — drift checker / matrix generator
- `datetime-spec.md` — normative datetime/duration functions

> Naming: the root package/directory is **`builtins`** (plural) — consistent with `vm.Builtins`,
> `cmd/builtins-check`, and the "built-in functions" vocabulary used throughout the docs. (Singular
> `builtin` is avoided because Go has a predeclared `builtin` pseudo-package.)

This document defines **how** the standard library is organized in code and the **phased plan** to migrate
the existing builtins and add the missing ones. It complements the roadmap (the *what*) with the
*structure* and *sequence*.

---

## 1. Goals

1. House all built-in functions in a root `builtins/` directory, organized **by type/family**.
2. Keep the runtime call mechanism unchanged: a builtin is still
   `func(args ...any) (any, error)`, assembled into the VM's function registry.
3. Make the code the **source of truth**: each function carries a descriptor (name, family, arity) so
   the manifest and the consistency checker can validate against real code, and the manifest can
   eventually be generated from it.
4. No import cycles; nothing imports the `vm` package.
5. Migrate incrementally — the build and tests stay green at every phase.

---

## 2. Current architecture (what we're refactoring)

- A builtin has signature `vm.VMFunction = func(args ...any) (any, error)`.
- `vm.VMFunctions = map[string]VMFunction` is the registry type.
- `vm.Builtins` (a literal map in `vm/builtins.go`) holds the 13 current builtins.
- The VM merges `vm.Builtins` with host-registered functions (`LibContext.Functions`).
- Public API exposes `uexl.Functions = vm.VMFunctions`.

The refactor changes **where functions live and how the default set is assembled** — not how they are
called.

---

## 3. Target directory structure

```
builtins/
  fn/
    fn.go                 # Func, Descriptor, Registry types — LEAF package, no deps
    args.go               # arg helpers: Num(args,i), Str(args,i), Bool, Arity(name,args,min,max)
    fn_test.go
  builtins.go               # Default() Registry — aggregates every family
  builtins_test.go          # no duplicate names; cross-checks against the manifest

  numbers/                # math family
    numbers.go            # abs, sign, round, floor, ceil, trunc, min, max, sum, avg, mod, pow, sqrt, clamp
    numbers_test.go
  strings/                # string + unicode-view operations
    basic.go              # substr, len(string), indexOf, startsWith, endsWith, contains, repeat
    case.go               # upper, lower
    trim.go               # trim, trimStart, trimEnd, padStart, padEnd
    edit.go               # replace
    split_join.go         # split, join
    unicode_views.go      # runeLen/runeSubstr, graphemeLen/graphemeSubstr, utf16Len/utf16Substr,
                          #   runes, graphemes, bytes, (char/utf8/utf16 views — T3)
    strings_test.go
  collections/            # array/object operations
    collections.go        # len, get, set, has, keys, values, remove, merge
    collections_test.go
  conversion/             # to* (lenient/null) · parse* (strict/error) · format*
    conversion.go         # str/toString, toNumber, parseNumber, formatNumber, toBoolean, parseBoolean
    conversion_test.go
  introspection/          # type queries
    introspection.go      # typeOf, isNull, isNumber, isString, isBool, isArray, isObject,
                          #   isDate, isDuration, isEmpty
    introspection_test.go
  datetime/               # wraps gotime; implements datetime-spec.md
    construct.go          # now, today, date, datetime, time
    arithmetic.go         # addMonths, addYears, diffMonths, diffYears, datePart
    epoch.go              # toEpochMillis/fromEpochMillis, toEpochSeconds/fromEpochSeconds
    format.go             # parseDate, formatDate (NITES)
    duration.go           # duration, durationIn, parseDuration, formatDuration
    datetime_test.go
  json/                   # T3 — optional / extension profile
    json.go               # parseJson, toJson
    json_test.go
```

Notes:
- `len` is polymorphic (string byte-length and collection size). It lives in `collections` and dispatches
  on argument type; the string view-lengths (`runeLen`, etc.) live in `strings/unicode_views.go`.
- `split`/`join` live together in `strings/split_join.go` (string-centric explode/reassemble).
- `datetime` is the only family with an external dependency (`gotime`).

---

## 4. Layering (no import cycles)

```
builtins/fn          ← leaf; imports nothing
builtins/<family>    → builtins/fn          (datetime additionally → gotime)
builtins             → builtins/fn + every family
vm                 → builtins             (builds vm.Builtins from builtins.Default())
```

Nothing imports `vm`. The shared types live in the leaf `builtins/fn` so families and the aggregator can
both use them without a cycle.

---

## 5. Shared types & registration

```go
// builtins/fn/fn.go
package fn

// Func is the universal builtin signature (identical to vm.VMFunction).
type Func func(args ...any) (any, error)

// Descriptor is the in-code declaration of a builtin; mirrors the manifest fields.
type Descriptor struct {
    Name    string
    Family  string
    MinArgs int
    MaxArgs int // -1 = variadic
    Fn      Func
}

type Registry map[string]Descriptor
```

```go
// builtins/fn/args.go — kill per-function boilerplate
func Arity(name string, args []any, min, max int) error
func Num(args []any, i int) (float64, error)
func Str(args []any, i int) (string, error)
func Bool(args []any, i int) (bool, error)
```

Each family exports its descriptors:

```go
// builtins/numbers/numbers.go
package numbers

import "github.com/maniartech/uexl/builtins/fn"

var Funcs = []fn.Descriptor{
    {Name: "abs", Family: "math", MinArgs: 1, MaxArgs: 1, Fn: abs},
    // ...
}

func abs(args ...any) (any, error) {
    if err := fn.Arity("abs", args, 1, 1); err != nil { return nil, err }
    x, err := fn.Num(args, 0); if err != nil { return nil, err }
    return math.Abs(x), nil
}
```

The aggregator assembles the default core library:

```go
// builtins/builtins.go
package builtins

func Default() fn.Registry {
    reg := fn.Registry{}
    for _, set := range [][]fn.Descriptor{
        numbers.Funcs, strings.Funcs, collections.Funcs,
        conversion.Funcs, introspection.Funcs, datetime.Funcs,
    } {
        for _, d := range set {
            if _, dup := reg[d.Name]; dup { panic("duplicate builtin: " + d.Name) }
            reg[d.Name] = d
        }
    }
    return reg
}
```

`vm` builds its registry from the descriptors (identical signature → trivial conversion):

```go
// vm: replace the vm/builtins.go literal
var Builtins = func() VMFunctions {
    out := VMFunctions{}
    for name, d := range builtins.Default() {
        out[name] = VMFunction(d.Fn)
    }
    return out
}()
```

---

## 6. Phased plan (build & tests stay green)

| Phase | Work | Manifest status change |
|-------|------|------------------------|
| **0 — Scaffold** | Add `builtins/fn` + `builtins.Default()`; point `vm.Builtins` at it (initially wrapping the existing 13) | — |
| **1 — Migrate** | Move the 13 current builtins into `strings`/`collections`/`conversion`; remove the `vm/builtins.go` literal; update manifest `impl` refs | impl refs updated |
| **2 — Reconcile (T0)** | `conversion`: implement `toNumber`/`toBoolean`; make `toString` an alias of `str`; implement (or remove) example fns | `documented-only → implemented` |
| **3 — Core families (T1)** | `numbers`, `introspection`, string ops, `split`, collection accessors, `parseNumber` | `gap → implemented` |
| **4 — Datetime (T1)** | `builtins/datetime` wrapping `gotime`; implements the datetime spec | `planned → implemented` |
| **5 — Nice-to-have (T2)** | `avg`/`clamp`, `trim*`, `pad*`, `merge`/`remove`, `formatNumber`, utf16 views, datetime conveniences | `gap → implemented` |
| **6 — Optional (T3)** | `json`, unicode view functions, locale/case profiles | extension entries |

Per-phase checklist:
1. Implement functions + unit tests in the family package.
2. Flip each function's `status` in `standard-library.json`.
3. `go run ./cmd/builtins-check` — matrix/completeness regenerates clean.
4. `go test ./...` — green.

Recommended order within T1: **numbers + introspection + `split`/`get`/`has` first** — they unblock the
most expressions.

---

## 7. Tooling tie-in (code as source of truth)

Once functions carry descriptors, extend `cmd/builtins-check`:

1. **Validate** the manifest against the live descriptors — compare `name`, `family`, and `arity`
   (`MinArgs`/`MaxArgs`), not just presence. Mismatches become drift.
2. **`-emit` mode** — generate `standard-library.json` from `builtins.Default()` descriptors, so the
   manifest becomes a build artifact rather than a hand-maintained file. The roadmap doc remains the
   human narrative; the manifest is derived; the matrix is generated.

This closes the loop: **code (descriptors) → manifest (emitted) → checker (validates) → matrix
(generated)**, with the roadmap as the prose companion.

---

## 8. Conventions (recap)

- `to*` = lenient coercion, `null` on failure; `parse*` = strict, **error** on failure; `format*` =
  value → string.
- Unicode ops come in byte / rune / grapheme / utf16 views.
- All builtins are pure; "current time" is context-injected (see datetime-spec §9.1).
- Cross-cutting semantics still to pin (rounding mode, `mod` sign, string-index unit, `typeOf`
  vocabulary, case/whitespace policy, JSON serialization) are tracked in
  `standard-library-roadmap.md` §6.
