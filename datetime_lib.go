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
