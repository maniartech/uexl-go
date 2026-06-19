package uexl

import "github.com/maniartech/uexl/builtins/datetime"

// WithDatetime returns an Option that registers the datetime standard-library functions (datetime-spec
// §11): date/datetime/time, parseDate/tryParseDate, duration, the epoch conversions
// (to/fromEpochMillis, to/fromEpochSeconds), the calendar functions addMonths/addYears/diffMonths/
// diffYears, datePart, durationIn, and formatDate/formatDur/parseDur/tryParseDur.
//
// The datetime/duration value types, their literals (d"…", 7d), and the temporal operators are part of
// the always-present core; this option adds only the function surface. Clock-dependent now()/today() are
// registered separately. Attach with uexl.DefaultWith(uexl.WithDatetime()).
func WithDatetime() Option {
	fns := make(Functions, len(datetime.Builtins))
	for name, f := range datetime.Builtins {
		fns[name] = Function(f)
	}
	return WithFunctions(fns)
}

// WithClock returns an Option that fixes the instant returned by now()/today() to nowMillis (milliseconds
// since the Unix epoch) for every evaluation in the environment — convenient for deterministic tests and
// a per-Env clock. For production, inject the instant per evaluation instead by passing "$now" in the
// eval vars (which shadows this global). now()/today() are read once and stable within a single
// evaluation either way (datetime-spec §9.1).
func WithClock(nowMillis int64) Option {
	return WithGlobals(map[string]any{"$now": nowMillis})
}

// isReservedClockFunc reports whether name is a built-in clock function (now/today). These are resolved
// from the injected instant ("$now") at evaluation time rather than the function registry, so the
// compiler accepts them without registration. A host may still register its own function of that name to
// override the default. See datetime-spec §9.1.
func isReservedClockFunc(name string) bool {
	return name == "now" || name == "today"
}
