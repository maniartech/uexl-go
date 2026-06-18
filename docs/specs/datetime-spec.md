# UExL DateTime Specification

Status: Draft
Audience: Language designers, implementers of UExL ports, conformance authors
Scope: Normative definition of the UExL `datetime` and `duration` value types, their literals,
operators, and core semantics
Layering: The `datetime`/`duration` **types**, the `d"..."` and duration suffix **literals**, and the
temporal **operators** (§3, §8) are part of the always-present language **core**. The temporal
**functions** (§10–§11) ship as the recommended, attachable **`datetime` standard library** (registered
via `WithLib`). See §1.2 and [ADR-0001](adr-0001-builtins-and-datetime-architecture.md).

> Normative keywords `MUST`, `MUST NOT`, `SHOULD`, `SHOULD NOT`, and `MAY` are used as defined in the
> [Language Specification Requirements](language-spec-requirements.md), §7.
>
> Sections marked **(Open)** are not yet finalized and are placeholders for forthcoming design
> discussion. They are non-normative until promoted.

## 1. Overview

UExL defines two temporal value types:

- **`datetime`** — an absolute *instant* in time. It combines date and time-of-day into one value,
  mirroring the single combined datetime type provided by most host platforms (Go `time.Time`,
  JavaScript `Date`, Dart `DateTime`, Python `datetime`, Java `Instant`, .NET `DateTime`).
- **`duration`** — an *exact elapsed amount* of time (e.g. the result of subtracting two `datetime`
  values). A `duration` represents a fixed number of milliseconds.

Both are distinct value types and both are represented as a signed integer count of milliseconds
(see §2). A `datetime` is **zoneless**: it identifies an instant and carries no timezone, offset, or
calendar-precision metadata.

This design is chosen so that temporal values are:

- **Portable** — represented identically on every conforming implementation.
- **Deterministic** — every core operation is reproducible with fixed input and fixed output, and is
  therefore conformance-testable.
- **Host-friendly** — losslessly convertible to and from the host platform's native datetime/duration
  types at the integration boundary.

### 1.1 The fixed-vs-calendar boundary (Informative)

A single rule explains most of this specification:

> **Exact** time amounts (millisecond, second, minute, hour, day, week) have a fixed length. They are
> modeled by `duration` and the `+` / `-` operators.
>
> **Calendar** units (month, year) have *variable* length — a month is 28–31 days, a year is 365 or
> 366 days — and have no agreed fixed conversion. They are **never** `duration` values and are **never**
> used with `+` / `-`. Calendar-relative arithmetic is performed only by the dedicated functions
> `addMonths`, `addYears`, `diffMonths`, and `diffYears` (and components are read with `datePart`),
> which resolve the variability against a specific date by explicit rules (§5.3, §11).

Because of this, UExL defines **no** "month = N days" constant. There is no industry-standard fixed
month or year length (conventions range from 30 days to 30.436875 days to calendar-correct variable
lengths), so the core refuses the conversion rather than pick an arbitrary one. The rigorous host
libraries `java.time.Duration` and Go `time.Duration` make the same choice: their exact-duration types
do not support months or years at all.

### 1.2 Core vs `datetime` library boundary (Normative)

UExL is layered: a small always-present **core**, plus opt-in standard-library bundles a host attaches at
construction (`WithLib`). Temporal support straddles that line by a single rule:

> If the **parser or evaluator** must understand it, it is **core**. If it is **just a function call**,
> it lives in the attachable **`datetime` library**.

This boundary is enforced by the engine: a library may register functions, pipe handlers, and globals,
but it cannot add literal syntax or operators (those are compiled into the tokenizer and VM).

| Always in core | In the attachable `datetime` library |
|----------------|---------------------------------------|
| `datetime` / `duration` value **types** (§2) | all temporal **functions** (§10–§11): construction, calendar arithmetic, extraction, epoch conversion |
| the `d"..."` **literal** + its **fixed ISO 8601 parser** (§3.1, §4) | format-directed and string **parsing**: `parseDate`/`tryParseDate`, `parseDur`/`tryParseDur` (§10.2–10.3) |
| duration **suffix literals** `7d`, `30ms` (§3.3) | **formatting**: `formatDate` (NITES), `formatDur` (§10.1, §10.3) |
| the temporal **operators** (§8) | `now`/`today` (clock-injected, §9.1) |
| `typeof` recognising both types | — |

Consequences:

- A conforming **core** MUST accept the `datetime`/`duration` types, the `d"..."` literal (parsing the
  fixed §4 ISO grammar), duration suffix literals, and all §8 operators — so every host can author and
  compare temporal values without attaching anything.
- A conforming **`datetime` library** MUST provide the §10–§11 functions. It is *recommended* for any
  general-purpose host; a minimal host MAY omit it. When omitted, the §10–§11 functions are simply not
  registered and calling one is an "unknown function" error — never a silently wrong result.
- The fixed ISO parser that backs the `d"..."` literal is core; the same parsing logic MAY be shared with
  the library's `parseDate`, but the runtime `parseDate`/`formatDate` *functions* are library surface.

## 2. Canonical Representation

A conforming implementation MUST represent both `datetime` and `duration` values as an integer count of
**milliseconds**, using a signed integer with at least 64 bits of range:

- A `datetime` is milliseconds since the Unix epoch `1970-01-01T00:00:00.000Z` (UTC); values before the
  epoch are negative.
- A `duration` is a count of elapsed milliseconds; it MAY be negative (e.g. `d1 - d2` where `d1 < d2`).
- The canonical unit MUST be **milliseconds**. Implementations MUST NOT expose sub-millisecond precision
  through core behavior. (Rationale: milliseconds is the lowest common denominator across host platforms,
  notably JavaScript `Date`.)

`datetime` and `duration` MUST be **distinct types** — distinguishable at runtime and not interchangeable
with each other or with `number`. (This distinction is what allows the language to permit
`datetime - datetime` while rejecting `datetime + number`; see §12.)

The canonical `datetime` value carries **no** timezone, offset, or "this was date-only" granularity flag.
Any such information present in source text is consumed during parsing (§4) and is not retained.

> Informative: `1970-01-01T00:00:00Z` → `0`; `1969-12-31T00:00:00Z` → `-86400000`;
> `2024-12-01T00:00:00Z` → `1733011200000`.

## 3. The `datetime` Literal

### 3.1 Syntax

A `datetime` literal is written with a `d` prefix immediately followed by a quoted string containing an
ISO 8601 / RFC 3339 date-time in the accepted subset (§4):

```
d"2024-12-01"
d"2024-12-01T10:30:00Z"
d"2024-12-01T10:30:00+05:30"
d'2024-12-01T10:30:00Z'
```

- The `d` prefix MUST be immediately adjacent to the opening quote, with no intervening whitespace.
- Both single (`'`) and double (`"`) quotes MUST be accepted, consistent with ordinary string literals.
- A bare identifier `d`, or any identifier beginning with `d` that is not immediately followed by a
  quote (e.g. `day`, `d + 1`), MUST continue to lex as an identifier. Only the `d"`/`d'` sequence
  introduces a `datetime` literal.
- The contents of a `datetime` literal are **not** subject to string escape processing.

> Implementation note (Informative): this mirrors the existing raw-string prefix `r"..."` in the Go
> tokenizer, which uses the same prefix-plus-quote-lookahead mechanism.

Durations may be written as **suffix literals** (`7d`, `3h`, §3.3), produced by the `duration()`
constructor for dynamic amounts (§11), or obtained by subtracting two `datetime` values (§12).

### 3.2 Compile-Time Constant

A `datetime` literal is a compile-time constant. A conforming implementation:

- MUST validate the literal against §4 at parse time and raise a parse-time error (§8) for any
  non-conforming literal.
- MUST compile the literal to the **canonical representation** (§2), not to a host-native datetime
  object. Host-native materialization MUST occur only at the integration boundary (§9).

### 3.3 Duration Literals

A `duration` literal is a numeric value immediately followed by a **fixed-unit suffix**, with no
intervening whitespace. It evaluates to a `duration` value (§2) equal to the number times the unit's
millisecond weight.

```
7d            // 7 days     → duration of 604800000 ms
3h            // 3 hours    → duration of 10800000 ms
30m           // 30 minutes → duration of 1800000 ms
500ms         // 500 milliseconds
1.5h          // fractional permitted → duration of 5400000 ms
d"2026-12-21" + 7d
```

**Suffixes.** The suffix set maps one-to-one onto the EXACT unit set (§11.1). It MUST NOT include a
month or year suffix, because those are not fixed-length durations (§1.1):

| Suffix | Unit | ms weight |
|--------|------|-----------|
| `ms` | millisecond | 1 |
| `s` | second | 1000 |
| `m` | minute | 60000 |
| `h` | hour | 3600000 |
| `d` | day | 86400000 |
| `w` | week | 604800000 |

> The absence of a month/year suffix is what lets `m` unambiguously mean **minute**, consistent with
> universal developer convention (`30m` = 30 minutes). This is a separate micro-syntax from NITES format
> strings, where the minute specifier is `i` and `m` is month (§10.1).

Suffix literals are the compact **authoring** form. For **interchange/serialization**, durations also
have an ISO 8601 string form (`PT1H30M`) via `parseDur` / `formatDur` (§10.3), restricted to
the same exact units (no year/month).

**Lexical rules.** A conforming tokenizer:

1. MUST match the **longest** valid suffix: `30ms` is 30 milliseconds, not `30m` followed by `s`.
2. MUST distinguish a duration literal `<number><suffix>` from the datetime prefix `d"…"`: a digit
   immediately before `d` begins a duration (`7d`), whereas `d` immediately before a quote begins a
   `datetime` literal (`d"…"`). The sequence `7d"…"` (no operator between) MUST be a parse error.
3. MUST accept a fractional numeric part (e.g. `1.5h`); the resulting millisecond value MUST be
   truncated toward zero to a whole millisecond (§4.3).
4. MUST treat a unary minus before a duration literal as negation of the resulting `duration`
   (`-7d` → a `duration` of `-604800000` ms; see §8).

**Constant vs dynamic.** Suffix literals require a literal numeric amount and therefore cover only
*constant* durations. For a computed amount, use the `duration(amount, unit)` constructor (§11), e.g.
`duration(x, "day")`, since `x d` is not valid syntax.

> **Open:** whether compound literals such as `1h30m` are accepted as a single token. UExL v1 requires
> operator composition instead (`1h + 30m`); compound single-token form MAY be added later as sugar.

## 4. Accepted Input Format

`datetime` literals (and the `parseDate` function, §10) MUST accept the following subset of ISO 8601 /
RFC 3339, and MUST reject all other forms with a parse-time error (for literals) or the function's
defined error/null result (for `parseDate`).

### 4.1 Accepted Forms

| Form | Example | Meaning |
|------|---------|---------|
| Date only | `2024-12-01` | Date at `00:00:00.000`, UTC |
| Date + time | `2024-12-01T10:30:00` | Given wall time, interpreted as UTC |
| Date + time + `Z` | `2024-12-01T10:30:00Z` | Exact UTC instant |
| Date + time + offset | `2024-12-01T10:30:00+05:30` | Instant computed from the offset |
| Fractional seconds | `2024-12-01T10:30:00.250Z` | Milliseconds component |

Grammar (informative sketch; the normative grammar will live in the grammar artifact):

```
datetime   = date [ "T" time [ offset ] ]
date       = year "-" month "-" day            ; year = 4 digits
time       = hour ":" minute [ ":" second [ "." frac ] ]
offset     = "Z" | ( ("+" | "-") hour ":" minute )
```

### 4.2 Defaulting Rules

When a component is absent, the following defaults MUST be applied:

1. **Missing time-of-day** defaults to `00:00:00.000` (midnight).
2. **Missing date** is not permitted; a `datetime` literal MUST include a date.
3. **Missing offset** MUST be interpreted as UTC.
4. **Present offset** MUST be used to compute the UTC instant and then MUST be discarded; it MUST NOT be
   retained on the value.

Consequently, `d"2024-12-01T10:30:00+05:30"` and `d"2024-12-01T05:00:00Z"` denote the **same** value.

### 4.3 Precision

Fractional seconds finer than milliseconds MUST be **truncated toward the epoch-millisecond boundary**
(truncated, not rounded) to produce the canonical value.

## 5. Range and Calendar

### 5.1 Calendar

All dates MUST be interpreted in the **proleptic Gregorian calendar** — the Gregorian calendar extended
backward through all earlier dates, with no Julian-calendar switchover at any historical date.

> Informative: this matches the internal calendar used by Java `java.time`, Python `datetime`, and
> JavaScript `Date`. A pre-1582 date therefore differs from the historically-used Julian date.

### 5.2 Guaranteed Portable Range

The **guaranteed portable range** for `datetime` values is:

```
0001-01-01T00:00:00.000Z  …  9999-12-31T23:59:59.999Z
```

- A `datetime` literal whose value falls outside this range MUST raise a parse-time error.
- This bound is driven by the narrowest mainstream host type (Python `datetime`, .NET `DateTime`, both
  limited to years `0001`–`9999`).
- Years outside `0001`–`9999`, the ISO year `0000`, negative years, and ISO expanded/extended year
  representations (e.g. `±YYYYYY`) MUST be rejected with a parse-time error.

> Informative: the canonical millisecond representation (§2) spans roughly ±292 million years and handles
> pre-epoch dates as negative values; the `0001`–`9999` clamp is a portability bound, not a
> representational one.

### 5.3 Field Validation and Calendar Rules

- Month, day, hour, minute, and second fields MUST be range-checked against the proleptic Gregorian
  calendar (e.g. `2024-13-01`, `2024-02-30`, `2024-12-01T25:00:00` MUST be parse-time errors).
- Leap seconds (a seconds field of `60`) MUST be rejected with a parse-time error.
- **End-of-month clamp (calendar shifts).** When `addMonths` or `addYears` (§11) shifts a date and the
  original day-of-month does not exist in the target month, the result MUST be clamped to the last valid
  day of the target month. Examples:

  | Expression | Result |
  |------------|--------|
  | `addMonths(d"2024-01-31", 1)` | `2024-02-29` |
  | `addMonths(d"2023-01-31", 1)` | `2023-02-28` |
  | `addMonths(d"2024-03-31", 1)` | `2024-04-30` |
  | `addYears(d"2024-02-29", 1)`  | `2025-02-28` |

- **Whole-unit calendar difference.** `diffMonths` and `diffYears` (§11) MUST return the signed count of
  **complete** calendar units from `a` to `b`, **truncated toward zero**. The count is defined as the
  largest `k ≥ 0` such that `addMonths(a, k) ≤ b` (respectively `addYears`), negated when `b` precedes
  `a`. This makes the difference the inverse of the add operation, so the end-of-month clamp (above)
  applies consistently. Examples:

  | Expression | Result |
  |------------|--------|
  | `diffMonths(d"2024-01-15", d"2024-03-10")` | `1` (second month incomplete) |
  | `diffMonths(d"2024-01-31", d"2024-02-29")` | `1` (clamped target reached) |
  | `diffMonths(d"2024-01-31", d"2024-02-28")` | `0` (clamped target not reached) |
  | `diffYears(d"2020-06-01", d"2024-05-01")`  | `3` (fourth year incomplete) |

### 5.4 Construction Validity and Range Clamping

Two distinct failure modes are handled differently — *invalid values* fail loudly, *out-of-range
instants* clamp:

1. **Invalid component value → error.** A structurally invalid field MUST raise an error (parse-time for
   literals, runtime for constructors): month ∉ `1..12`; day not valid for the given month/year
   (e.g. `date(2023, 2, 29)`); hour ∉ `0..23`; minute or second ∉ `0..59`; millisecond ∉ `0..999`.
   Leap seconds (second = `60`) are likewise invalid (§5.3).

2. **Out-of-range instant → clamp.** When a *computed* instant falls outside the guaranteed portable
   range, it MUST be clamped to the nearest boundary instant:

   ```
   MIN = 0001-01-01T00:00:00.000Z   (-62135596800000 ms)
   MAX = 9999-12-31T23:59:59.999Z   ( 253402300799999 ms)
   ```

   Clamping applies to results of arithmetic operators (§8), the calendar functions `addMonths`/
   `addYears`, and numeric constructor inputs whose year falls outside `0001..9999`
   (e.g. `date(12000, 1, 1)` clamps to `MAX`; `fromEpochMillis(n)` with out-of-range `n` clamps).
   Clamping is silent (no error).

3. **Literals.** A `datetime` literal can only express a 4-digit year grammatically, so an out-of-range
   year such as `d"10000-01-01"` is a **format** parse error (§5.2), not a clamp. Invalid component
   values in a literal are parse errors per rule 1.

> Rationale: a wrong *value* (month 13, Feb 30) is almost always a mistake and should fail; drifting past
> the artificial `0001..9999` portability boundary during computation is not a correctness error, so
> graceful clamping keeps pipelines running. This matches the user-facing intent "clamp the range, throw
> on invalid dates."

## 6. Equality, Ordering, and Truthiness

### 6.1 Comparison Operators

The comparison operators `==`, `!=`, `<`, `<=`, `>`, `>=` MUST operate on temporal values:

- For two `datetime` values, and for two `duration` values, comparison MUST be by canonical millisecond
  value (§2).
- `==` / `!=` between values of different types (e.g. `datetime` vs `duration`, or `datetime` vs
  `number`) MUST be defined by the language's general cross-type equality rules; temporal values MUST NOT
  be considered equal to a `number` with the same millisecond value.
- Ordering operators (`<`, `<=`, `>`, `>=`) between a temporal value and a non-comparable type are
  governed by the general comparison rules in the semantic specification.

> Note: comparison is the only operator family that applies *between* two `datetime` values besides `-`.
> See §12 for the full operator table.

### 6.2 Truthiness

A temporal value is **falsy** when its canonical millisecond value is `0` — i.e. the epoch instant
`d"1970-01-01T00:00:00.000Z"` and a zero-length `duration` — and **truthy** for every other value (every
other instant, including any real-world `now()`, and every non-zero duration). This follows UExL's
general rule that the **zero value** of a type is falsy (`0`, `""`, `false`, empty list/map), independent
of whether the value is "present": an empty string is a real, non-null string yet falsy, and likewise the
zero instant and zero span are real values yet falsy.

A temporal value is **never nullish** — it is falsy-but-present, so the nullish-coalescing operator does
not treat it as null (`d"1970-01-01T00:00:00Z" ?? x` evaluates to the epoch instant, exactly as `0 ?? x`
evaluates to `0`). The truthiness test inspects the canonical millisecond value directly and performs
**no** implicit conversion to `number` (§8.2).

## 7. Timezones

UExL `datetime` values are **zoneless UTC instants**. For the core (mandatory) profile:

- Literals and core operations MUST accept **UTC** and **fixed numeric offsets** only (`Z`, `+HH:MM`,
  `-HH:MM`).
- A fixed offset MUST be used solely to compute the UTC instant and MUST then be discarded (§4.2).
- The host is responsible for resolving any **named** timezone (e.g. `America/New_York`) to a fixed
  numeric offset before passing it into UExL.

### 7.1 Out of Scope (Core)

The following MUST NOT be part of the core profile:

- Named IANA timezones.
- Daylight-saving-time-aware ("civil") arithmetic.
- Any operation whose result depends on a timezone database.

> Rationale (Informative): named-zone behavior depends on the IANA timezone database, which changes over
> time and varies by host. Such operations cannot be expressed as fixed input → fixed output conformance
> cases and would undermine the cross-platform determinism guarantee.

### 7.2 Forward Compatibility (Informative)

Because the canonical value is always a zoneless instant, named-zone support can be added later as an
optional `iana-timezones` capability/profile without changing the representation, the literal grammar, or
any existing expression.

## 8. Operators and Arithmetic

UExL v1 supports a deliberately small set of temporal operators. All are deterministic and exact.

### 8.1 Operator Table

| Left | Op | Right | Result | Notes |
|------|----|-------|--------|-------|
| `datetime` | `-` | `datetime` | `duration` | exact elapsed time (may be negative) |
| `datetime` | `+` | `duration` | `datetime` | shift forward |
| `duration` | `+` | `datetime` | `datetime` | shift forward (commutative form) |
| `datetime` | `-` | `duration` | `datetime` | shift backward |
| `duration` | `+` | `duration` | `duration` | |
| `duration` | `-` | `duration` | `duration` | |
| `duration` | `*` | `number` | `duration` | scale |
| `number` | `*` | `duration` | `duration` | scale |
| `duration` | `/` | `number` | `duration` | divide span (result truncated to whole ms) |
| `duration` | `/` | `duration` | `number` | ratio (e.g. how many fit) |
| (unary) | `-` | `duration` | `duration` | negation (e.g. `-7d`) |
| `datetime`/`duration` | `== != < <= > >=` | same type | `boolean` | §6 |

### 8.2 Rejected Operations

The following MUST raise a runtime type error (or, where statically determinable, a parse-time error):

- `datetime` `+` `datetime` — adding two instants is meaningless.
- `datetime` `+`/`-` `number`, and `duration` `+`/`-` `number` — bare numbers have no unit; a `duration`
  must be constructed explicitly via `duration()` (§11).
- Any `+` / `-` involving `month` or `year` — calendar shifts MUST use `addMonths` / `addYears` (§11),
  never operators.

> Rationale: forbidding `datetime + number` is the reason `datetime` and `duration` are distinct types
> (§2). It forces the unit to be explicit at the call site, eliminating the "+ 7 of what?" ambiguity.

### 8.3 Division

Both forms are supported (following `java.time.Duration`):

- `duration / number → duration` — divide a span (e.g. split into equal parts); the result MUST be
  truncated toward zero to a whole millisecond.
- `duration / duration → number` — a ratio (e.g. `duration(1,"hour") / duration(15,"minute")` → `4`).

Division by a zero `number` or a zero-length `duration` MUST raise a runtime error, consistent with the
language's general division-by-zero behavior.

### 8.4 Range Overflow

When an operator yields a `datetime` outside the guaranteed portable range, the result MUST be clamped to
the boundary instant per §5.4 (silent clamp). `duration` results are not range-bound beyond the signed
64-bit millisecond representation.

## 9. Host Interoperability

- At the integration boundary, a conforming implementation SHOULD materialize a `datetime` value as the
  host platform's native datetime type, constructed from the canonical millisecond instant in UTC
  (e.g. Go `time.UnixMilli(ms).UTC()`, JavaScript `new Date(ms)`, Dart
  `DateTime.fromMillisecondsSinceEpoch(ms, isUtc: true)`).
- A `duration` value SHOULD be materialized as the host's native duration type where one exists
  (e.g. Go `time.Duration` = `ms * time.Millisecond`, Dart `Duration(milliseconds: ms)`); hosts without
  a native duration type MAY represent it as an integer count of milliseconds.
- A host-native datetime passed **into** a UExL context MUST be normalized to the canonical representation
  (UTC instant, millisecond precision; sub-millisecond precision truncated per §4.3).
- The canonical representation — not any host-native object — MUST be used inside the evaluator and in
  any compiled/serialized constant form.

### 9.1 Current Instant (`now` / `today`)

The current instant MUST be supplied as **context-injected** state (a host-provided clock value), not
read from an ambient host clock during evaluation. The host captures the system clock **once, at the
start of the evaluation** (before any pipe stage runs) and injects it as the evaluation's current
instant. This follows the clock-injection pattern of `java.time.Clock` and keeps evaluation pure and
deterministic (conformance tests inject a fixed instant, equivalent to `Clock.fixed`).

A conforming implementation:

- MUST resolve `now()` (§11) from the injected clock instant for the current evaluation.
- MUST return the **same** instant for **every** call to `now()` within a single expression evaluation —
  **including across different pipe stages**. The instant is captured once at the start of evaluation and
  frozen for its duration. This mirrors SQL `CURRENT_TIMESTAMP` (stable within a statement) and
  guarantees that `now() == now()` and `now() - now()` evaluates to a zero `duration`, so a value cannot
  "age" between stages of the same expression.
- MUST resolve `today()` (§11) from the **same** captured instant, truncated to UTC midnight
  (`00:00:00.000Z`); `today()` is therefore equally stable across the whole evaluation.
- When no clock instant is injected, the behavior of `now()` / `today()` is host-defined (it MAY raise an
  error or use the host wall clock); such evaluations are outside the deterministic conformance guarantee.

> Implementation note (Informative): the captured instant is delivered through the host's existing
> context/variable channel — the same plumbing used for ordinary evaluation context — so `now`/`today`
> require **no** change to the builtin call signature. `now`/`today` are part of the attachable
> `datetime` library (§1.2); they read the core-provided injected instant.

## 10. Parsing and Formatting

Because a `datetime` carries no display metadata, all rendering is **explicit**: the caller chooses how
much of the value to show. UExL uses the **NITES** format syntax (Natural and Intuitive Time Expression
Syntax) as its date/time pattern language, and **ISO 8601** for durations.

### 10.1 DateTime Formatting — NITES

`formatDate(d: datetime, pattern: string)` → `string` MUST render a `datetime` using a NITES pattern or
NITES named layout. NITES is case-insensitive; the canonical reference is the NITES specification, and
`gotime` is a reference implementation.

Common specifiers (see the NITES specification for the full set):

| Specifier | Meaning | Example |
|-----------|---------|---------|
| `yyyy` / `yy` / `y` | year (4-digit / 2-digit / no-pad) | `2026` |
| `mm` / `m` | month number (padded / no-pad) | `12` |
| `mmm` / `mmmm` | month name (short / full) | `Dec` / `December` |
| `dd` / `d` | day of month | `21` |
| `hhh` | hour, 24-hour | `15` |
| `hh` / `h` | hour, 12-hour | `03` |
| `ii` / `i` | **minute** (padded / no-pad) | `04` |
| `ss` / `s` | second | `05` |
| `aa` / `a` | AM/PM | `PM` / `pm` |
| `www` / `wwww` | weekday name (short / full) | `Mon` / `Monday` |

> Note: in NITES, the minute specifier is `i` (and month is `m`). This is independent of the duration
> literal suffixes (§3.3), where `m` means minute by the universal convention (`30m` = 30 minutes). The
> two are different micro-syntaxes (a format *string* vs a numeric *suffix*) and do not interact.

**Named layouts.** Implementations MUST support at least the named layout `iso`
(`yyyy-mm-ddThhh:ii:ss`) and SHOULD support the other NITES named layouts (`rfc`, `sql`, `date`,
`isodate`, etc.). `formatDate(d)` called with no pattern MUST default to the `iso` layout.

**Locale and zone constraints (core).** In the core (deterministic) profile:

- Name-based specifiers (`mmm`, `mmmm`, `www`, `wwww`, `aa`) MUST render English/invariant names.
  Locale-aware rendering is reserved for a future locale profile.

  > Informative: the locale profile's data path is already designed (gotime
  > `dev-guidelines/localization/`): a slim CLDR-derived JSON is the single source of truth, and the
  > UExL host injects it as a context map (UExL reads no files). That contract embeds CLDR plural-rule
  > conditions **authored in UExL expression syntax** (`&&`/`||`/`!`, `==`, `%`) evaluated with the
  > CLDR operands (`n, i, v, w, f, t`) as context — making UExL the cross-platform rule notation for
  > the locale data itself.
- Since `datetime` is a zoneless UTC instant, components and the zone specifiers render in **UTC**
  (`z` → `Z`; `o`/`oo`/`ooo` → `+00`/`+0000`/`+00:00`). The timezone-abbreviation specifier (`zz`,
  e.g. `MST`) requires a timezone database and is therefore **not** part of the core profile.

### 10.2 DateTime Parsing — NITES and ISO

String parsing follows the language-wide **`parseX` (strict) / `tryParseX` (safe)** convention
([ADR-0001](adr-0001-builtins-and-datetime-architecture.md) §B):

- `parseDate(s: string)` → `datetime` MUST parse the §4 ISO 8601 subset (the canonical interchange form)
  and MUST **raise an error** on invalid input (never `null`). `parseDur` (§10.3) behaves the same way.
- `tryParseDate(s: string)` → `datetime | null` MUST behave identically to `parseDate` on valid input and
  MUST **return `null`** (not raise) on invalid input. `tryParseDur` (§10.3) is its duration counterpart.
  This is the same strict/safe split as `parseNum`/`tryParseNum` for numbers.
- `parseDate(s: string, pattern: string)` → `datetime` (pattern-directed parsing using a NITES pattern)
  is **Open** for v1 and MAY be deferred; ordinal and name specifiers that are documented as
  format-only in NITES are not valid for parsing.

### 10.3 Duration Interchange — ISO 8601

For serialization and interchange, durations use the **ISO 8601 duration** format, restricted to the
**exact** components (consistent with §1.1 — no calendar units):

- `formatDur(d: duration)` → `string` MUST produce an ISO 8601 duration using only week/day/hour/
  minute/second components (e.g. `PT1H30M`, `P7D`, `PT0.5S`).
- `parseDur(s: string)` → `duration` MUST accept ISO 8601 durations restricted to those exact
  components. The year designator (`Y`) and the **date-part** month designator (`M` before `T`) MUST be
  rejected, because months and years are not exact durations. The minute designator (`M` after `T`) is
  accepted. (ISO 8601 disambiguates `M` by position relative to `T`.) Invalid input MUST raise an error.
- `tryParseDur(s: string)` → `duration | null` is the safe counterpart: it returns `null` on invalid
  input instead of raising (per the `parseX`/`tryParseX` convention, §10.2).

> Rationale: this makes both temporal types standards-based on the wire — `datetime` via ISO 8601 / NITES
> `iso`, and `duration` via ISO 8601 duration — while the exact-only restriction prevents the variable
> calendar units from re-entering the `duration` type. Suffix literals (§3.3) remain the compact in-source
> authoring form; ISO 8601 strings are the interchange form.

### 10.4 Fixed-Offset Rendering and Extraction

Because a `datetime` is a zoneless UTC instant, formatting and component extraction default to UTC. To
render or read components *as seen at a fixed offset*, functions accept an optional `offset` argument:

- `formatDate(d, pattern, offset)` and `datePart(d, component, offset)` apply the offset only at
  render/extract time; the underlying instant is unchanged and no zoned value is created.
- `offset` is a fixed ISO 8601 offset string: `"Z"` / `"+00:00"` (UTC), `"+05:30"`, `"-08:00"`.
- When omitted, the offset is UTC.

```
formatDate(d"2024-12-01T00:30:00Z", "yyyy-mm-dd", "+05:30")   // "2024-12-01" (06:00 local)
datePart(d"2024-12-01T00:30:00Z", "day", "+05:30")            // 1
datePart(d"2024-12-01T00:30:00Z", "hour", "-01:00")           // 23 (prev day in UTC terms)
```

> The host resolves any **named** timezone to a fixed offset (§7) before passing it in. Named-zone /
> DST-aware rendering remains out of the core profile.

## 11. Standard-Library Functions

The following functions form the attachable **`datetime` standard library** (§1.2). They are *recommended*
for general-purpose hosts and registered via `WithLib`; a minimal host MAY omit them. Names are
provisional pending the standard-library artifact; semantics are normative.

Most temporal arithmetic is performed with the **core** operators (§8) and the `d"..."` / suffix literals,
so the library itself is small — it adds construction, calendar arithmetic, extraction, epoch conversion,
parsing, and formatting on top of the core types.

### 11.1 Unit and Component Names

Two closed string-name sets are defined, both lowercase singular:

```
EXACT      = "millisecond" | "second" | "minute" | "hour" | "day" | "week"
COMPONENT  = "year" | "month" | "day" | "hour" | "minute" | "second" | "millisecond" | "weekday"
```

`EXACT` units carry these millisecond weights: `millisecond` = 1, `second` = 1000, `minute` = 60000,
`hour` = 3600000, `day` = 86400000, `week` = 604800000. They are used by `duration` / `durationIn` and
correspond one-to-one to the duration suffixes (§3.3). Note `EXACT` contains **no** `month` or `year`,
because those have no fixed length (§1.1).

`COMPONENT` names are the fields readable by `datePart`. `weekday` returns the ISO-8601 day of week
(`1` = Monday … `7` = Sunday). (Day-of-year, week-of-year, and quarter are **Open** — see §13.)

When a function receives a unit/component argument that is a constant string, an unknown or
out-of-domain value MUST be a **parse-time** error; when the argument is dynamic, it MUST be a runtime
error.

Unit and component names are **singular only** in v1; plural forms (`"months"`) MUST be rejected. (A
single canonical spelling keeps conformance unambiguous; accepting aliases later would be a
non-breaking relaxation.)

### 11.2 Functions

Function names are **plural** for the calendar add/diff family (matching `.NET` `AddMonths`/`AddYears`).

**Construction**

| Function | Signature | Returns | Notes |
|----------|-----------|---------|-------|
| `now` | `()` | `datetime` | current instant; context-injected clock; stable within one evaluation (§9.1) |
| `today` | `()` | `datetime` | **UTC** midnight of the current date (injected clock); local "today" needs an offset |
| `date` | `(year, month, day)` | `datetime` | instant at `00:00:00.000Z`; validity/clamp per §5.4 |
| `datetime` | `(year, month, day, hour?, minute?, second?, millisecond?)` | `datetime` | trailing args default to `0`; §5.4 |
| `time` | `(hour, minute, second?, millisecond?)` | `datetime` | on the epoch date `1970-01-01`; trailing args default to `0`; §5.4 |
| `parseDate` | `(s: string)` | `datetime` | strict; ISO 8601 subset (§10.2); **error** on invalid |
| `tryParseDate` | `(s: string)` | `datetime \| null` | safe; `null` on invalid (§10.2) |
| `duration` | `(amount: number, unit)` | `duration` | EXACT units only; `month`/`year` MUST error |

**Epoch conversion** (the `datetime`↔`number` bridge; no implicit coercion)

| Function | Signature | Returns | Notes |
|----------|-----------|---------|-------|
| `toEpochMillis` | `(d: datetime)` | `number` | ms since `1970-01-01Z` (exact; canonical, §2) |
| `fromEpochMillis` | `(n: number)` | `datetime` | inverse; out-of-range `n` clamps (§5.4) |
| `toEpochSeconds` | `(d: datetime)` | `number` | Unix seconds (`floor(ms/1000)`) |
| `fromEpochSeconds` | `(n: number)` | `datetime` | inverse; out-of-range `n` clamps (§5.4) |

**Calendar arithmetic / extraction**

| Function | Signature | Returns | Accepted names | Notes |
|----------|-----------|---------|----------------|-------|
| `addMonths` | `(d: datetime, amount: number)` | `datetime` | — | end-of-month clamp (§5.3); negative `amount` subtracts |
| `addYears` | `(d: datetime, amount: number)` | `datetime` | — | end-of-month clamp (§5.3); negative `amount` subtracts |
| `diffMonths` | `(a: datetime, b: datetime)` | `number` | — | whole calendar months, truncated toward zero (§5.3) |
| `diffYears` | `(a: datetime, b: datetime)` | `number` | — | whole calendar years, truncated toward zero (§5.3) |
| `datePart` | `(d: datetime, component, offset?: string)` | `number` | **COMPONENT** | extraction; UTC unless `offset` given (§10.4) |
| `durationIn` | `(d: duration, unit)` | `number` | **EXACT** | may be fractional; `month`/`year` MUST error |

**Formatting**

| Function | Signature | Returns | Notes |
|----------|-----------|---------|-------|
| `formatDate` | `(d: datetime, pattern?: string, offset?: string)` | `string` | NITES pattern/layout; default `iso`; UTC unless `offset` given (§10.1, §10.4) |
| `formatDur` | `(d: duration)` | `string` | ISO 8601 duration (§10.3) |
| `parseDur` | `(s: string)` | `duration` | strict; ISO 8601 duration, exact components only (§10.3); **error** on invalid |
| `tryParseDur` | `(s: string)` | `duration \| null` | safe; `null` on invalid (§10.3) |

Notes:

- There is **no** general `addDate`/`diffDate` function. Exact arithmetic uses operators: shift with
  `date + duration` (§8) and exact difference with `date - date` → `duration`, decomposed via
  `durationIn`. Calendar *difference* in whole months/years uses `diffMonths` / `diffYears`.
- `date`/`datetime`/`time` are the dynamic (computed-argument) constructors; the `d"..."` literal (§3)
  is the constant form. They validate fields and clamp range per §5.4.
- `duration()` and `durationIn()` MUST reject `month` and `year` (the core defines no fixed month/year
  length; §1.1).
- `addMonths` / `addYears` / `diffMonths` / `diffYears` are the only calendar-aware functions, and the
  only place month/year variability is resolved.
- `datePart` and `formatDate` default to UTC; an explicit fixed `offset` (§10.4) renders/extracts as
  seen at that offset without changing the instant.

## 12. Conformance Cases (Starter)

Results are canonical millisecond values (§2) unless otherwise indicated. (IDs are provisional; all
numeric expectations MUST be verified against the reference implementation before promotion.)

| ID | Expression | Expected |
|----|------------|----------|
| dt-lit-001 | `d"1970-01-01T00:00:00Z"` | `0` |
| dt-lit-002 | `d"2024-12-01"` | `1733011200000` |
| dt-lit-003 | `d"2024-12-01T10:30:00+05:30" == d"2024-12-01T05:00:00Z"` | `true` |
| dt-lit-004 | `d"1969-12-31T00:00:00Z"` | `-86400000` |
| dt-lit-005 | `d"0001-01-01T00:00:00Z"` | `-62135596800000` |
| dt-op-001  | `d"2024-12-02" - d"2024-12-01"` | duration of `86400000` |
| dt-op-002  | `d"2024-12-01" + duration(2, "day")` | `1733184000000` |
| dt-op-003  | `durationIn(d"2024-12-02" - d"2024-12-01", "hour")` | `24` |
| dt-op-004  | `d"2024-12-01" < d"2024-12-02"` | `true` |
| dt-lit-006 | `d"2024-12-01" + 7d` | `1733616000000` |
| dt-lit-007 | `30ms == duration(30, "millisecond")` | `true` |
| dt-lit-008 | `1.5h == duration(90, "minute")` | `true` |
| dt-lit-009 | `-7d == duration(-7, "day")` | `true` |
| dt-cal-001 | `addMonths(d"2024-01-31", 1) == d"2024-02-29"` | `true` |
| dt-cal-002 | `addYears(d"2024-02-29", 1) == d"2025-02-28"` | `true` |
| dt-cal-003 | `addMonths(d"2024-03-15", -2) == d"2024-01-15"` | `true` |
| dt-cal-004 | `datePart(d"2026-12-21", "month")` | `12` |
| dt-cal-005 | `diffMonths(d"2024-01-15", d"2024-03-10")` | `1` |
| dt-cal-006 | `diffYears(d"2020-06-01", d"2024-05-01")` | `3` |
| dt-div-001 | `duration(1, "hour") / duration(15, "minute")` | `4` |
| dt-div-002 | `duration(1, "hour") / 2 == duration(30, "minute")` | `true` |
| dt-fmt-001 | `formatDate(d"2026-12-21T15:04:05Z", "yyyy-mm-dd")` | `"2026-12-21"` |
| dt-fmt-002 | `formatDate(d"2026-12-21T15:04:05Z")` | `"2026-12-21T15:04:05"` |
| dt-fmt-003 | `formatDur(duration(90, "minute"))` | `"PT1H30M"` |
| dt-fmt-004 | `parseDur("PT1H30M") == duration(90, "minute")` | `true` |
| dt-ctor-001 | `date(2024, 12, 1) == d"2024-12-01"` | `true` |
| dt-ctor-002 | `datetime(2024, 12, 1, 10, 30) == d"2024-12-01T10:30:00Z"` | `true` |
| dt-ctor-003 | `time(10, 30) == d"1970-01-01T10:30:00Z"` | `true` |
| dt-epoch-001 | `toEpochMillis(d"2024-12-01")` | `1733011200000` |
| dt-epoch-002 | `fromEpochMillis(0) == d"1970-01-01T00:00:00Z"` | `true` |
| dt-epoch-003 | `toEpochSeconds(d"1970-01-01T00:00:01Z")` | `1` |
| dt-off-001 | `formatDate(d"2024-12-01T00:30:00Z", "yyyy-mm-dd", "+05:30")` | `"2024-12-01"` |
| dt-off-002 | `datePart(d"2024-12-01T00:30:00Z", "hour", "+05:30")` | `6` |
| dt-off-003 | `datePart(d"2024-12-02", "weekday")` | `1` (Monday) |
| dt-clamp-001 | `addYears(d"9999-06-01", 5) == d"9999-12-31T23:59:59.999Z"` | `true` |
| dt-clamp-002 | `date(12000, 1, 1) == d"9999-12-31T23:59:59.999Z"` | `true` |
| dt-err-001 | `d"2024-13-01"` | parse error (field out of range) |
| dt-err-002 | `d"2024-02-30"` | parse error (field out of range) |
| dt-err-003 | `d"2024-12-01T10:30:60Z"` | parse error (leap second) |
| dt-err-004 | `d"0000-01-01"` | parse error (year out of range) |
| dt-err-005 | `d"10000-01-01"` | parse error (year out of range) |
| dt-err-006 | `duration(2, "month")` | error (calendar unit on duration) |
| dt-err-007 | `d"2024-12-01" + 7` | error (bare number) |
| dt-err-008 | `d"2024-12-01" + d"2024-12-02"` | error (datetime + datetime) |
| dt-err-009 | `7d"2024-12-01"` | parse error (no operator between duration and datetime) |
| dt-err-010 | `parseDur("P1Y")` | error (year not an exact duration) |
| dt-err-011 | `parseDur("P3M")` | error (date-part month not an exact duration) |
| dt-err-012 | `duration(2, "months")` | error (plural unit name rejected, §11.1) |
| dt-err-013 | `date(2023, 2, 29)` | error (invalid day for month, §5.4) |
| dt-err-014 | `date(2024, 13, 1)` | error (invalid month, §5.4) |
| dt-err-015 | `datetime(2024, 1, 1, 25, 0)` | error (invalid hour, §5.4) |
| dt-err-016 | `parseDate("not-a-date")` | error (strict parse; not `null`, §10.2) |
| dt-parse-001 | `tryParseDate("not-a-date")` | `null` (safe parse, §10.2) |
| dt-parse-002 | `tryParseDur("nonsense")` | `null` (safe parse, §10.3) |
| dt-parse-003 | `tryParseDate("2024-12-01") == d"2024-12-01"` | `true` |

## 13. Open Decisions Summary

Resolved in this revision:

- **Core vs library layering** — the `datetime`/`duration` **types**, the `d"..."` / suffix **literals**,
  and the **operators** are always-present **core**; the **functions** (§10–§11) ship as the recommended,
  attachable **`datetime` library** (`WithLib`). The `d"..."` literal's fixed-ISO parse is core; all
  runtime parse/format functions are library surface (§1.2,
  [ADR-0001](adr-0001-builtins-and-datetime-architecture.md) §A).
- Canonical representation; distinct `datetime`/`duration` types.
- DateTime and duration literal syntax (including suffix literals); duration ISO 8601 interchange.
- Accepted format and defaulting; range and calendar rules; timezone policy.
- Equality/ordering, truthiness, the full operator table including division (§8.3).
- Truthiness — temporal values follow the **zero-is-falsy** rule: the epoch instant and a zero duration
  are falsy (like `0`/`""`/`false`), every other value is truthy, and temporal values are never nullish
  (§6.2).
- Calendar difference — `diffMonths`/`diffYears` **added** (truncated toward zero, §5.3, §11).
- Formatting/parsing — **NITES** for `datetime` (default layout `iso`), **ISO 8601** for `duration`
  (§10), with `gotime` as the reference implementation.
- `now()` / `today()` — context-injected clock captured **once** at evaluation start; the same instant is
  returned by every call, **stable across all pipe stages** (§9.1).
- Unit/component names — **singular only** in v1 (§11.1).
- Function names — **plural** calendar family (`addMonths`/`addYears`/`diffMonths`/`diffYears`).
- Dynamic constructors — `today`, `date`, `datetime`, `time` (§11.2).
- Validity vs range — **invalid components throw, out-of-range instants clamp** to `MIN`/`MAX` (§5.4, §8.4).
- Fixed-offset rendering/extraction — optional `offset` argument on `formatDate`/`datePart` (§10.4).
- Epoch conversion — explicit `toEpochMillis`/`fromEpochMillis`/`toEpochSeconds`/`fromEpochSeconds` (§11.2).
- `datePart` gains `weekday` (ISO 1=Mon … 7=Sun) (§11.1).
- Conversion naming — **`parseX` (strict, raises) / `tryParseX` (safe, `null`)** with short type tokens
  (`parseDate`/`tryParseDate`, `parseDur`/`tryParseDur`); supersedes the earlier `to*`/`parse*` split
  (§10.2, §10.3, [ADR-0001](adr-0001-builtins-and-datetime-architecture.md) §B).

Remaining open:

1. Pattern-directed **parsing** `parseDate(s, pattern)` — deferred for v1; ISO parsing only (§10.2).
2. Compound duration literals such as `1h30m` as a single token — deferred (§3.3).
3. Locale-aware name rendering (`mmmm`, `wwww`, `aa`) — reserved for a future locale profile (§10.1).
4. Additional `datePart` components — `dayOfYear`, `weekOfYear`, `quarter` — deferred (§11.1).
5. All numeric conformance expectations to be verified against the reference implementation (§12).
