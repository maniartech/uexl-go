# UExL Standard Library — Inventory, Drift & Roadmap

Status: Draft
Audience: Maintainers, spec authors, port implementers
Companion artifacts:
- `standard-library.json` — the machine-readable manifest (source of truth for the tooling)
- `standard-library.md` — the generated consistency matrix + family-completeness report
- `cmd/builtins-check` — the drift checker / matrix generator
- `datetime-spec.md` — normative datetime/duration functions
- `language-spec-requirements.md` — overall spec process

This document is the human-readable companion to the generated matrix. It lists **everything**: what
exists, what has drifted, what is committed (specified), and what is missing by family symmetry — with a
proposed signature and a priority tier for each. It is the working backlog for the core standard library.

> Naming/behavior conventions used throughout (see also `datetime-spec.md`):
> - **`to*`** = lenient coercion; returns `null` on failure (e.g. `toNumber`).
> - **`parse*`** = strict parser; **raises an error** on failure (e.g. `parseDate`).
> - **`format*`** = value → string.
> - Unicode operations come in **byte / rune / grapheme / utf16** views (e.g. `len`/`runeLen`/
>   `graphemeLen`/`utf16Len`).
> - All builtins are **pure** (no mutation of inputs, no ambient state; "current time" is injected).
> - Errors are raised for invalid *types/arguments*; specific functions document `null` vs error.

---

## 1. Current state — implemented builtins (13)

These exist today in `vm.Builtins`.

| Function | Family | Signature | Returns | Notes |
|----------|--------|-----------|---------|-------|
| `len` | length | `(string\|array)` | number | byte length for strings; element count for arrays |
| `substr` | substring | `(string, start, length)` | string | byte-indexed |
| `contains` | string-ops | `(string, string)` | bool | |
| `set` | collection | `(object\|array, key, value)` | object\|array | copy-returning (immutable) update |
| `str` | conversion | `(any)` | string | stringify; book calls this `toString` |
| `runeLen` | length | `(string)` | number | code-point count |
| `runeSubstr` | substring | `(string, start, length)` | string | code-point indexed |
| `graphemeLen` | length | `(string)` | number | grapheme-cluster count (UAX #29) |
| `graphemeSubstr` | substring | `(string, start, length)` | string | grapheme indexed |
| `runes` | explode | `(string)` | array | explode to code points |
| `graphemes` | explode | `(string)` | array | explode to grapheme clusters |
| `bytes` | explode | `(string)` | array | explode to bytes |
| `join` | reassemble | `(array, sep?)` | string | |

---

## 2. Drift to reconcile (Tier 0 — do first)

Inconsistencies between docs, implementation, and naming. These are correctness/credibility issues, not
new features.

| Item | Problem | Action |
|------|---------|--------|
| `toNumber` | Documented in `book/type-conversion.md`, **not implemented** | Implement (lenient, `null` on fail) **or** remove from book |
| `toBoolean` | Documented, **not implemented** | Implement **or** remove from book |
| `toString` vs `str` | Book documents `toString`; impl provides `str` | Pick one canonical name; alias the other (recommend `str` canonical, `toString` alias) |
| `min`/`max`/`sum`/`concat` | Used in `book/functions/overview.md` examples, **not implemented** | Implement (see Math, §4.1) or fix the examples |

---

## 3. Committed — specified but not yet implemented (datetime/duration)

Fully specified in `datetime-spec.md` (status `planned`). Listed here for completeness.

**Construction:** `now`, `today`, `date`, `datetime`, `time`, `parseDate`
**Formatting/parse:** `formatDate`, `parseDuration`, `formatDuration`
**Epoch:** `toEpochMillis`, `fromEpochMillis`, `toEpochSeconds`, `fromEpochSeconds`
**Calendar:** `addMonths`, `addYears`, `diffMonths`, `diffYears`, `datePart`
**Duration:** `duration`, `durationIn`

(20 functions; semantics are normative in the datetime spec. Go reference: `gotime`.)

---

## 4. Missing by family (symmetry audit)

The "exists on some, not on others" gaps. Each function has a proposed signature and a priority tier.
Tiers: **T1** = core essential (target v1), **T2** = core nice-to-have, **T3** = optional/extension.

### 4.1 Math (entire family absent) — T1

No numeric functions exist; `min`/`max`/`sum` are already referenced in book examples.

| Function | Signature | Returns | onFail | Tier | Notes |
|----------|-----------|---------|--------|------|-------|
| `abs` | `(number)` | number | error | T1 | |
| `sign` | `(number)` | number | error | T1 | -1 / 0 / 1 |
| `round` | `(number, digits?)` | number | error | T1 | half-away-from-zero; `digits` default 0 |
| `floor` | `(number)` | number | error | T1 | |
| `ceil` | `(number)` | number | error | T1 | |
| `trunc` | `(number)` | number | error | T1 | toward zero |
| `min` | `(number, …)` / `(array)` | number | error | T1 | variadic + array form |
| `max` | `(number, …)` / `(array)` | number | error | T1 | |
| `sum` | `(array)` / `(number, …)` | number | error | T1 | |
| `avg` | `(array)` / `(number, …)` | number | error | T2 | mean |
| `mod` | `(number, number)` | number | error | T1 | result sign convention to pin (recommend Euclidean or match `%`) |
| `pow` | `(number, number)` | number | error | T1 | |
| `sqrt` | `(number)` | number | error | T1 | |
| `clamp` | `(number, lo, hi)` | number | error | T2 | |

> Decisions: rounding mode (recommend half-away-from-zero), `mod` sign, NaN/Inf propagation (defer to
> `numeric-semantics.md` IEEE-754 rules).

### 4.2 Type introspection (entire family absent) — T1

Especially relevant now that `datetime`/`duration` exist.

| Function | Signature | Returns | Tier | Notes |
|----------|-----------|---------|------|-------|
| `typeOf` | `(any)` | string | T1 | `"number"`/`"string"`/`"bool"`/`"array"`/`"object"`/`"null"`/`"datetime"`/`"duration"` |
| `isNull` | `(any)` | bool | T1 | true for null/absent |
| `isNumber` | `(any)` | bool | T1 | |
| `isString` | `(any)` | bool | T1 | |
| `isBool` | `(any)` | bool | T1 | |
| `isArray` | `(any)` | bool | T1 | |
| `isObject` | `(any)` | bool | T1 | |
| `isDate` | `(any)` | bool | T1 | datetime type |
| `isDuration` | `(any)` | bool | T2 | |
| `isEmpty` | `(any)` | bool | T2 | empty string/array/object; null → true |

> Decision: the exact `typeOf` string set is the canonical type-name vocabulary — pin it in the data
> model spec so all ports agree.

### 4.3 String operations — T1/T2

Only `contains`/`substr` exist; the common toolkit is missing.

| Function | Signature | Returns | onFail | Tier | Notes |
|----------|-----------|---------|--------|------|-------|
| `indexOf` | `(string, sub, from?)` | number | error | T1 | -1 if not found; index unit must match `substr` (byte) |
| `startsWith` | `(string, prefix)` | bool | error | T1 | |
| `endsWith` | `(string, suffix)` | bool | error | T1 | |
| `replace` | `(string, old, new)` | string | error | T1 | all occurrences (decide first-vs-all; recommend all) |
| `split` | `(string, sep)` | array | error | T1 | inverse of `join` (also listed in §4.4) |
| `trim` | `(string)` | string | error | T1 | both ends; whitespace definition to pin |
| `trimStart` | `(string)` | string | error | T2 | |
| `trimEnd` | `(string)` | string | error | T2 | |
| `upper` | `(string)` | string | error | T1 | invariant case in core; locale-aware later |
| `lower` | `(string)` | string | error | T1 | invariant case in core |
| `padStart` | `(string, len, pad?)` | string | error | T2 | |
| `padEnd` | `(string, len, pad?)` | string | error | T2 | |
| `repeat` | `(string, count)` | string | error | T2 | |

> Decisions: index/length unit (byte vs rune vs grapheme — recommend byte for `indexOf` to match
> `substr`, with rune/grapheme variants later); `replace` first-vs-all; whitespace set for `trim`;
> case-fold locale policy (invariant in core).

### 4.4 Explode / reassemble — T1

| Function | Signature | Returns | onFail | Tier | Notes |
|----------|-----------|---------|--------|------|-------|
| `split` | `(string, sep)` | array | error | T1 | **the missing inverse of `join`** |

### 4.5 Collection accessors — T1/T2

`set` exists but its companions don't. (Iteration — map/filter/reduce/sort — is **pipe-based by design**
and intentionally not duplicated as functions.)

| Function | Signature | Returns | onFail | Tier | Notes |
|----------|-----------|---------|--------|------|-------|
| `get` | `(object\|array, key, default?)` | any | error/null | T1 | safe read; `default` returned when absent |
| `has` | `(object\|array, key)` | bool | error | T1 | |
| `keys` | `(object)` | array | error | T1 | |
| `values` | `(object)` | array | error | T1 | |
| `remove` | `(object\|array, key)` | object\|array | error | T2 | copy-returning (immutable), mirrors `set` |
| `merge` | `(object, object)` | object | error | T2 | copy-returning shallow merge |

> Decision: `get`/`set`/`has`/`remove` are the immutable accessor set; confirm `get` default-vs-error
> behavior (recommend: `default` provided → return it; omitted → error on missing, consistent with strict
> access philosophy).

### 4.6 Conversion completion — T1/T2

Make the `to*` / `parse*` / `format*` triple consistent across types (datetime/duration already have it).

| Function | Signature | Returns | onFail | Tier | Notes |
|----------|-----------|---------|--------|------|-------|
| `toNumber` | `(any)` | number\|null | null | T0 | reconcile drift (§2) |
| `parseNumber` | `(string)` | number | **error** | T1 | strict counterpart to `toNumber` |
| `formatNumber` | `(number, pattern?)` | string | error | T2 | decimals/grouping; pattern model TBD |
| `toBoolean` | `(any)` | bool | null | T0 | reconcile drift (§2) |
| `parseBoolean` | `(string)` | bool | error | T2 | strict `"true"`/`"false"` |

### 4.7 JSON — T3 (optional / extension)

Only needed when expressions receive JSON *as a string* (object/array literals already cover JSON-shaped
data). Drags in number-precision and parsing-semantics decisions → keep out of core.

| Function | Signature | Returns | onFail | Tier | Notes |
|----------|-----------|---------|--------|------|-------|
| `parseJson` | `(string)` | any | error | T3 | strict |
| `toJson` | `(any, pretty?)` | string | error | T3 | serialization rules must be pinned for portability |

### 4.8 Unicode utf16 views & view functions — T2/T3

Round out the byte/rune/grapheme/utf16 symmetry and the view functions promised in
`design-philosophy.md`.

| Function | Signature | Returns | Tier | Notes |
|----------|-----------|---------|------|-------|
| `utf16Len` | `(string)` | number | T2 | JS/UTF-16 host parity |
| `utf16Substr` | `(string, start, length)` | string | T2 | |
| `char` | `(string)` | view | T3 | code-point view (per design-philosophy) |
| `utf8` | `(string)` | view | T3 | byte view |
| `utf16` | `(string)` | view | T3 | code-unit view |

### 4.9 Datetime conveniences — T2

Beyond the committed datetime set; map cleanly onto `gotime`.

| Function | Signature | Returns | Tier | gotime |
|----------|-----------|---------|------|--------|
| `startOfDay`/`startOfMonth`/`startOfYear` | `(datetime)` | datetime | T2 | `DayStart`/`MonthStart`/`YearStart` |
| `endOfDay`/`endOfMonth`/`endOfYear` | `(datetime)` | datetime | T2 | `MonthEnd`/`YearEnd` |
| `isLeapYear` | `(number)` | bool | T2 | `IsLeapYear` |
| `daysInMonth` | `(year, month)` | number | T2 | `DaysInMonth` |
| `datePart "dayOfYear"/"weekOfYear"/"quarter"` | component | number | T2 | extends `datePart` (datetime-spec §11 Open) |

---

## 5. Prioritized roadmap (summary)

| Tier | Theme | Functions (count) |
|------|-------|-------------------|
| **T0 — Reconcile** | Fix doc/impl drift | `toNumber`, `toBoolean`, `toString`/`str`, example fns (4 items) |
| **T1 — Core essentials** | Math, type-introspection, core string ops, `split`, collection accessors, `parseNumber` + the committed datetime/duration set | ~50 |
| **T2 — Core nice-to-have** | `avg`/`clamp`, `trimStart/End`, `pad*`, `repeat`, `remove`/`merge`, `formatNumber`, utf16 views, datetime conveniences | ~20 |
| **T3 — Optional / extension** | JSON, Unicode view functions, locale-aware case/format, named timezones | ~10 |

Recommended order: **T0 → T1 (math + introspection + split/get/has first, they unblock the most
expressions) → T2 → T3.**

---

## 6. Cross-cutting decisions to pin (block precise specs)

1. **Numeric:** rounding mode, `mod` sign convention, NaN/Inf propagation (link `numeric-semantics.md`).
2. **String indexing unit** for `indexOf`/`replace`/`split` (recommend byte to match `substr`, with
   rune/grapheme variants later).
3. **Case-fold & whitespace** policy: invariant in core; locale-aware behind a profile.
4. **`typeOf` vocabulary:** the canonical type-name strings (shared with the data-model spec).
5. **`get`/`replace`** behavior: default-vs-error, first-vs-all.
6. **JSON serialization** rules (if/when JSON becomes an extension profile).

---

## 7. Maintenance

- The manifest `standard-library.json` is the source of truth for the tooling.
- Regenerate the matrix + completeness report: `go run ./cmd/builtins-check`.
- CI gate (drift only): `go run ./cmd/builtins-check -check` (exit 1 on hard drift).
- When a function is implemented, update its manifest `status` to `implemented`; when specified, add a
  `specRef`. The checker will flag any divergence between this roadmap, the manifest, and the registry.
