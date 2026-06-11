# UExL Competitive Analysis — Expression Language Landscape

Status: Living document — update as competitors evolve and as UExL ships
Last updated: 2026-06-05
Audience: Maintainers, contributors, adopters evaluating UExL
Companion artifacts: `datetime-spec.md`, `standard-library-roadmap.md`,
`builtins-implementation-plan.md`, `language-spec-requirements.md`

This document compares UExL against the embedded expression-language field — primarily **CEL**
(Google's Common Expression Language) and **Expr** (expr-lang), plus the wider set (govaluate, gval,
goja/otto, JsonLogic, JMESPath, starlark). It records where UExL leads, where it lags, and what must
ship for each claim to be defensible. Keep it honest: this doc is only useful if it survives contact
with a skeptical adopter.

> Maintenance: re-verify the performance numbers (§5) after significant VM changes, and re-check
> competitor capabilities roughly twice a year (CEL and Expr both move).

---

## 1. The field

| Engine | Positioning | Spec? | Implementations | Notes |
|--------|-------------|-------|-----------------|-------|
| **UExL** | Cross-platform expression spec ("RegEx for expressions") | Draft (requirements + datetime complete) | Go (JS/Rust/Dart planned) | Pipes, NITES, grapheme-aware Unicode |
| **CEL** | Policy/config expressions (K8s, Envoy, Firebase) | **Mature + conformance suite** | **Go, C++, Java (official)** | Protobuf-typed, non-Turing-complete |
| **Expr** | Go-embedded expressions | No — implementation is the spec | Go only | Fast, pragmatic, rich Go interop |
| govaluate | Go-embedded, minimal | No | Go only | Largely dormant |
| gval | Go-embedded, extensible | No | Go only | |
| goja / otto | Full JS engines | (ECMAScript) | Go | Heavyweight for expressions |
| starlark | Config language (Bazel) | Yes | Go, Java, Rust | A language, not an expression DSL |
| JsonLogic | JSON-encoded portable logic | Thin | Many | Tiny feature set |
| JMESPath | JSON query | Yes + compliance tests | Many | Query, not general expressions |

**CEL is the only true benchmark** — the one competitor sharing UExL's exact ambition: spec-first,
multi-implementation, deterministic, portable. Most of this analysis is therefore UExL vs CEL, with
Expr as the Go-ecosystem pragmatist baseline.

---

## 2. Capability matrix

| Capability | **UExL** | **CEL** | **Expr** |
|------------|----------|---------|----------|
| Datetime type | ✅ first-class, spec'd (`d"…"` literal) | ✅ protobuf `Timestamp` (no literal) | ⚠️ host `time.Time` via env |
| Duration type | ✅ first-class, exact-only, `7d` literals | ✅ protobuf `Duration`, **max unit: hour** | ⚠️ Go `time.Duration` |
| Calendar arithmetic (add/diff months, years) | ✅ **spec'd with end-of-month clamp** | ❌ none | ❌ none |
| Date/time format pattern language | ✅ **NITES** (own standard, ref impl) | ❌ RFC3339 string conversion only | ❌ none |
| Deterministic `now` | ✅ injected clock, stable per evaluation | ⚠️ env convention | ❌ wall clock |
| Timezone determinism | ✅ fixed-offset-only core; named zones excluded | ⚠️ IANA names → host tzdb dependence | ⚠️ host `Location` |
| Unicode model | ✅ explicit byte/rune/grapheme views | ⚠️ code points only | ⚠️ Go strings (bytes) |
| Iteration | ✅ pipes (`\|map:`, `\|filter:`, params) | ⚠️ comprehension macros | ✅ builtins + closures |
| Static type checking | ❌ dynamic | ✅ **protobuf-typed, compile-time** | ⚠️ partial |
| Conformance suite | ⚠️ seeded (datetime), harness pending | ✅ **full, run by 3 impls** | ❌ |
| Library consistency tooling | ✅ **manifest + drift checker + symmetry audit** | ❌ | ❌ |
| Implemented stdlib breadth | ❌ 13 functions today (~80 specified) | ✅ core + extension libs | ✅ dozens of builtins |
| Non-Turing-complete / sandboxable | ✅ | ✅ | ✅ |
| Purity guarantees (no mutation) | ✅ documented policy | ✅ | ⚠️ host-dependent |

---

## 3. Datetime deep-dive (UExL's strongest design lead)

The datetime/duration specification (`datetime-spec.md`) is, by design comparison, the strongest
temporal model in this class:

1. **Calendar arithmetic exists and is pinned.** CEL and Expr have *no* month/year arithmetic at all —
   "invoice date + 1 month" is inexpressible. UExL specifies `addMonths`/`addYears`/`diffMonths`/
   `diffYears` with a normative end-of-month clamp and an inverse-of-add difference rule, each with
   conformance cases.
2. **Durations cover days and weeks.** CEL's duration inherits Go's hour-capped units — `7d` must be
   written `duration("168h")`. UExL: `7d`, `1.5h`, `1w` literals, with the exact-vs-calendar boundary
   keeping months/years out of durations (same refusal java.time and Go make — then UExL adds the
   calendar functions they stop short of).
3. **Parse-time temporal literals.** `d"2024-12-01"` validates at parse time and compiles to a
   canonical constant. CEL's `timestamp("…")` is a runtime conversion; Expr has no literal.
4. **Stronger determinism.** Zoneless UTC instants; fixed numeric offsets only in core; named-IANA
   zones excluded (host resolves them) so no tzdb-version nondeterminism — the exact leak CEL's
   timezone-aware accessors permit. `now()` is context-injected and stable within one evaluation
   (SQL `CURRENT_TIMESTAMP` semantics).
5. **A formatting standard.** NITES (case-insensitive, intuitive specifiers, named layouts) with
   `gotime` as a 100%-coverage reference implementation. Neither CEL nor Expr has any pattern language.
6. **Footgun avoidance.** CEL ships `getMonth()` → 0–11 and both `getDate()` (1-based) and
   `getDayOfMonth()` (0-based). UExL: `datePart(d, "month")` → 1–12, singular unit names only, and
   `date + number` is a type error (the unit must be visible: `+ 7d` or `duration(n, unit)`).
7. **Independent validation.** CEL pins the same portable range (years 1–9999) UExL chose;
   java.time/Go made the same no-months-in-durations call. UExL's choices sit on the same ground as
   the most rigorous prior art, then extend it.

Caveat for honesty: the UExL datetime functions are **specified, not yet implemented** (status
`planned` in the manifest). The design lead is real; the shipping lead is CEL's until Phase 4 of the
builtins plan lands.

---

## 4. Standard library

- **Today:** 13 implemented functions (strings/Unicode-heavy) — behind both CEL and Expr in shipped
  breadth. The symmetry audit (`standard-library.md`) records ~50 missing siblings, including the
  entire math and type-introspection families.
- **Specified:** ~80 functions across 9 families with signatures, failure modes (`to*` lenient-null /
  `parse*` strict-error / `format*`), and priority tiers (`standard-library-roadmap.md`).
- **Unique process advantage:** no competitor has a machine-readable manifest + drift checker +
  family-symmetry audit (`cmd/builtins-check`). CEL's conformance suite tests behavior; nothing in the
  field catches "documented-but-unimplemented" or "join-without-split" automatically.

---

## 5. Performance (verified against repo benchmarks)

Source: `wip-notes/FINAL_PERFORMANCE_RESULTS.md`, `wip-notes/BENCHMARK_COMPARISON.md`
(AMD Ryzen 7 5700G, Windows, Go 1.21+, 10–15s runs; expression:
`(Origin == "MOW" || Country == "RU") && (Value >= 100 || Adults == 1)`).

### Headline micro-benchmark (boolean policy expression)

| Framework | ns/op | B/op | allocs/op |
|-----------|-------|------|-----------|
| expr | 132.1 | 32 | 1 |
| cel-go | 174.1 | 16 | 1 |
| **UExL** | **227.1** | **0** | **0** ✅ |
| govaluate | ~315 | — | — |
| goja | ~324 | — | — |
| otto | ~650 | — | — |
| gval | ~894 | — | — |
| starlark | ~7,248 | — | — |

**#3 of 10 on raw dispatch; the only framework at zero allocations.** Zero allocs means no GC
pressure and predictable latency under sustained load — expr spends ~16% of its time in `mallocgc`.

### Where UExL is outright fastest

| Workload | UExL | expr | cel-go |
|----------|------|------|--------|
| Array map (100 elems, `\|map: $item * 2`) | **3,428 ns** | 10,588 ns (3.1× slower) | 44,339 ns (13× slower) |
| String pattern matching | **256 ns** | 285 ns | 347 ns |
| Function calls | **186 ns** | 198 ns | 234 ns |

The pipe engine is the standout: collection processing — the realistic workload for an expression
language fed arrays of data — is ~3× faster than expr and ~13× faster than cel-go.

### Honest framing

- On the raw boolean micro-benchmark, expr and cel-go are faster (~95 ns gap, attributed to the
  48-byte Value struct vs 16-byte interface stacks — a deliberate architecture-over-microspeed call).
- On allocations (0 vs 1) and on collection/pipe workloads, **UExL leads the field**.
- Historic trajectory: 9,388 ns → 227 ns (≈41× improvement) across the optimization phases.

**Performance claim that survives scrutiny:** *top-3 raw speed among 10 Go expression engines, #1 in
allocation efficiency, and #1 — by 3–13× — in collection processing.*

---

## 6. Where UExL is honestly behind (and what closes each gap)

| Gap | Leader | What closes it |
|-----|--------|----------------|
| Implemented stdlib (13 fns) | Expr / CEL | Builtins plan Phases 0–3 (T0 reconcile + T1 core families) |
| Datetime functions unimplemented | CEL (shipping Timestamp/Duration) | Builtins plan Phase 4 (wrap `gotime`) |
| Conformance suite is a seed | CEL (full suite, 3 impls) | Verify §12 numbers; build the harness per requirements doc §8.9 |
| Single implementation | CEL (Go/C++/Java) | One non-Go port (JS first) passing the suite |
| No static type checking | CEL | Possibly out of scope (dynamic-by-design for the Excel-user audience) — decide and document |
| Raw dispatch speed | expr | Optional micro-work (stack-value compaction); architecture call already made |

---

## 7. Verdict

**Design: best in class.** The datetime/duration specification solves problems CEL declined to touch
(calendar arithmetic, temporal literals, formatting), avoids footguns CEL shipped (0-based months,
hour-capped durations, tzdb leakage), and grounds every rule in a conformance case. The pipe syntax,
grapheme-aware Unicode model, and the library-consistency tooling are likewise ahead of the field.

**Performance: verified excellent.** Top-3 of 10 on raw micro-benchmarks, uniquely zero-allocation,
and the outright leader (3–13×) on collection workloads — the workload that matters most in practice.

**Overall, today: best-in-class trajectory, not yet best-in-class fact.** CEL still wins on shipped
maturity — conformance suite, three implementations, static typing, production mileage. The honest
one-liner: **UExL has out-designed CEL and out-engineered the Go field on the workloads that matter;
it now has to out-ship CEL.** The path is written down: execute the builtins plan (T0–T1), implement
the datetime spec, verify the conformance numbers, stand up one non-Go port. Each converts a design
win into a defensible public claim.

## 8. Claim checklist (update as items land)

- [x] Best-in-class temporal *specification* (calendar math, literals, NITES, determinism)
- [x] Zero-allocation evaluation (verified)
- [x] Fastest collection/pipe processing among Go engines (verified, 3–13×)
- [x] Library consistency tooling (manifest + drift checker + symmetry audit) — unique in field
- [ ] Core stdlib parity with Expr/CEL (builtins plan T0–T1)
- [ ] Datetime/duration implemented and passing conformance cases
- [ ] Machine-readable conformance suite + harness
- [ ] Second implementation (non-Go) passing the suite
- [ ] Spec 1.0 published (authority hierarchy flips from impl to spec)

When every box is checked, "best in class" is no longer a claim — it's a description.
