# UExL Date and Time Handling

Status: Draft for future version
Audience: Maintainers, language designers, implementers of UExL ports
Scope: Proposed future handling for date/time values, durations, arithmetic, parsing, and capability boundaries. This document does not make date/time part of the current mandatory core standard library.

## 1. Purpose

Date and time handling is essential for many real expression-hosting use cases: scheduling, expiry checks, SLAs, reporting windows, retention rules, pricing windows, and audit logic.

UExL does not currently include date/time helpers in the mandatory core built-in function set. This document records the intended future direction so the feature can be added in a controlled, portable way rather than as host-specific one-off helpers.

## 2. Current status

- Date and time helpers are not part of the current mandatory core built-in set.
- Hosts MAY inject custom date/time functions today, but those are host extensions and are not portability-safe.
- Any future standardization should happen as a named profile, not as an unscoped expansion of the current core.

## 3. Why this is deferred from the current core

Date/time support is important, but it is also one of the easiest places to create accidental incompatibility.

Key risks:

- Parsing formats vary widely across hosts and standard libraries.
- Time zone databases and naming support differ by platform.
- Local times can be ambiguous across daylight-saving transitions.
- Formatting often pulls in locale behavior, calendar assumptions, and host-specific layout syntax.
- "Current time" introduces external state into what is otherwise a deterministic expression model.

For those reasons, date/time should not enter the mandatory core until the profile boundary is explicit.

## 4. Design goals

Any future date/time profile SHOULD satisfy all of the following:

- Portable across conforming implementations.
- Locale-independent by default.
- Explicit about time zone handling.
- Small enough for non-Go ports to implement without depending on host-specific formatting quirks.
- Honest about nondeterministic behavior such as `now()`.
- Structured in layers so a minimal profile can ship before full formatting and time-zone support.

## 5. Recommended profile structure

Date/time support SHOULD be standardized as layered named profiles.

### 5.1 `datetime-core`

This should be the first profile and should stay intentionally small.

It SHOULD define:

- a first-class `datetime` value kind representing an instant in time
- a first-class `duration` value kind representing a signed span of time
- comparison of datetime values
- subtraction of datetime minus datetime resulting in duration
- addition and subtraction of duration with datetime
- deterministic parsing of a narrow, fixed input set
- basic extraction helpers such as year/month/day/hour/minute/second/weekday

It SHOULD NOT initially include:

- locale-aware formatting
- arbitrary host format strings
- named time-zone database lookups
- calendar systems other than the default civil calendar
- "parse anything the host library accepts" behavior

### 5.2 `datetime-zones`

This optional follow-on profile MAY add explicit time-zone identifiers, zone conversion, and named-zone support.

This profile would need to define:

- whether IANA time-zone names are required
- what happens when a zone is unknown on a host
- how ambiguous local times are resolved
- whether zone offsets are preserved in formatting/parsing round-trips

### 5.3 `datetime-format`

This optional follow-on profile MAY add explicit formatting and custom parsing.

This profile MUST avoid host-native layout syntax unless that syntax is standardized by UExL itself. Otherwise portability breaks immediately.

## 6. Proposed `datetime-core` surface

The following is the recommended first-pass surface for future consideration.

### 6.1 Value kinds

- `datetime`
- `duration`

These are conceptual UExL value categories. Implementations MAY map them to host-native types internally, but the observable language behavior must be standardized.

### 6.2 Proposed built-ins

#### Parsing and creation

- `datetime(str)`
  - parse a datetime string using a fixed, standardized accepted set
- `duration(str)`
  - parse a duration string such as `"1h30m"`

#### Field extraction

- `year(dt)`
- `month(dt)`
- `day(dt)`
- `hour(dt)`
- `minute(dt)`
- `second(dt)`
- `weekday(dt)`

#### Optional nondeterministic helper

- `now()`

`now()` SHOULD be capability-gated even within `datetime-core`, because it introduces host clock dependence.

## 7. Proposed operator behavior

If `datetime-core` is adopted, the following operator semantics SHOULD be standardized:

- `dt1 == dt2`, `!=`, `<`, `<=`, `>`, `>=` on two datetime values
- `dt + dur` -> datetime
- `dt - dur` -> datetime
- `dt1 - dt2` -> duration

The following SHOULD remain errors unless a later profile says otherwise:

- `datetime + datetime`
- `duration * duration`
- string arithmetic involving datetime values

## 8. Parsing scope for `datetime-core`

To remain portable, the initial parser SHOULD accept a small fixed set only.

Recommended first-pass accepted inputs:

- RFC3339 / ISO-8601 timestamp with explicit offset or `Z`
- date-only form `YYYY-MM-DD`
- time-only form only if the language defines exactly how it is anchored

Recommended first-pass rejected inputs:

- host-specific shorthand formats
- locale-dependent month/day names
- platform-native layout strings
- zone names like `Europe/Zurich` in `datetime-core`

## 9. Time-zone model

The initial profile SHOULD prefer explicit offsets over named zones.

That means:

- offset-bearing timestamps are portable
- named-zone conversions belong in `datetime-zones`
- local-machine default timezone behavior MUST NOT be implicit

If a datetime lacks explicit zone information, the profile MUST define exactly how it is interpreted. Silent host-local interpretation SHOULD be avoided.

## 10. Error model

Date/time support SHOULD define stable error categories for at least:

- invalid datetime literal or string input
- invalid duration input
- unsupported profile or capability
- invalid timezone identifier
- ambiguous or nonexistent local time, if local-zone parsing is ever standardized
- wrong arity
- wrong type

Human-readable messages may vary, but the error categories and required triggering conditions should be stable.

## 11. Capability declaration guidance

Hosts that support future date/time features SHOULD declare support explicitly.

Recommended capability structure:

- language version
- supported profile list, including `datetime-core` if implemented
- optional capability flag for `now()` if host clock access is enabled
- optional support for `datetime-zones`
- optional support for `datetime-format`

## 12. Non-goals for the first date/time profile

The first standardized version SHOULD NOT attempt to solve all of these at once:

- localized formatting
- arbitrary format tokens
- business calendars
- recurring schedules and cron semantics
- leap-second-aware semantics beyond the chosen host abstraction
- time-zone database distribution and update policy

## 13. Recommended rollout order

1. `datetime-core`
2. optional `now()` capability within `datetime-core`
3. `datetime-zones`
4. `datetime-format`

This preserves momentum while keeping portability honest.

## 14. Relationship to current core builtin finalization

Date/time helpers are intentionally excluded from the current mandatory core built-in set. Their future standardization should happen through the profile model described here, not by expanding the current core ad hoc.