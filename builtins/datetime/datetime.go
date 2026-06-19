// Package datetime is the UExL "datetime" standard-library bundle: thin wrappers that expose the pure
// datetime engine in package types (calendar math, extraction, NITES formatting, ISO 8601 duration) as
// UExL builtin functions. Attach via uexl.WithDatetime(). Implements datetime-spec §11. now()/today()
// (clock-injected) are registered separately.
package datetime

import (
	"fmt"
	"math"
	"strings"
	"time"

	gotime "github.com/maniartech/gotime/v2"
	"github.com/maniartech/uexl/types"
)

// namedLayouts are the NITES named layouts recognized by the datetime library, expanded to gotime's
// NITES dialect (24-hour is hhhh; see datetime-spec §10.1). gotime itself does not know these names, so
// formatDate/parseDate resolve them before delegating. "iso" is the default.
var namedLayouts = map[string]string{
	"iso":     "yyyy-mm-ddThhhh:ii:ss",
	"isodate": "yyyy-mm-dd",
	"date":    "yyyy-mm-dd",
	"time":    "hhhh:ii:ss",
	"sql":     "yyyy-mm-dd hhhh:ii:ss",
	"rfc":     "www, dd mmm yyyy hhhh:ii:ss",
}

// defaultLayout is the layout used by formatDate when no pattern is supplied.
const defaultLayout = "iso"

// resolveLayout expands a NITES named layout (case-insensitive) to its gotime pattern, or returns the
// pattern unchanged when it is already a raw NITES pattern.
func resolveLayout(pattern string) string {
	if exp, ok := namedLayouts[strings.ToLower(strings.TrimSpace(pattern))]; ok {
		return exp
	}
	return pattern
}

// localTime materializes a canonical (zoneless, UTC-ms) instant as a time.Time whose wall clock is the
// instant as seen offMin minutes east of UTC, so gotime renders the intended local components.
func localTime(msUTC int64, offMin int) time.Time {
	if offMin == 0 {
		return time.UnixMilli(msUTC).UTC()
	}
	return time.UnixMilli(msUTC).In(time.FixedZone("", offMin*60))
}

// Builtins maps UExL function names to their implementations (the datetime standard library, minus the
// clock-dependent now()/today() which are added by the host with an injected instant).
var Builtins = map[string]func(args ...any) (any, error){
	"date":             builtinDate,
	"datetime":         builtinDatetime,
	"time":             builtinTime,
	"parseDate":        builtinParseDate,
	"tryParseDate":     builtinTryParseDate,
	"duration":         builtinDuration,
	"toEpochMillis":    builtinToEpochMillis,
	"fromEpochMillis":  builtinFromEpochMillis,
	"toEpochSeconds":   builtinToEpochSeconds,
	"fromEpochSeconds": builtinFromEpochSeconds,
	"addMonths":        builtinAddMonths,
	"addYears":         builtinAddYears,
	"diffMonths":       builtinDiffMonths,
	"diffYears":        builtinDiffYears,
	"datePart":         builtinDatePart,
	"durationIn":       builtinDurationIn,
	"formatDate":       builtinFormatDate,
	"formatDur":        builtinFormatDur,
	"parseDur":         builtinParseDur,
	"tryParseDur":      builtinTryParseDur,
}

// exactUnits maps the EXACT duration unit names to their millisecond weight (datetime-spec §11.1). There
// is deliberately no month/year (variable length).
var exactUnits = map[string]int64{
	"millisecond": 1,
	"second":      1000,
	"minute":      60000,
	"hour":        3600000,
	"day":         86400000,
	"week":        604800000,
}

// ---- argument helpers ----

func arity(name string, args []any, min, max int) error {
	if len(args) < min || len(args) > max {
		if min == max {
			return fmt.Errorf("%s expects %d argument(s), got %d", name, min, len(args))
		}
		return fmt.Errorf("%s expects %d..%d arguments, got %d", name, min, max, len(args))
	}
	return nil
}

func argDateTime(name string, args []any, i int) (types.DateTime, error) {
	if d, ok := args[i].(types.DateTime); ok {
		return d, nil
	}
	return types.DateTime{}, fmt.Errorf("%s: argument %d must be a datetime, got %T", name, i+1, args[i])
}

func argDuration(name string, args []any, i int) (types.Duration, error) {
	if d, ok := args[i].(types.Duration); ok {
		return d, nil
	}
	return types.Duration{}, fmt.Errorf("%s: argument %d must be a duration, got %T", name, i+1, args[i])
}

func argInt(name string, args []any, i int) (int, error) {
	f, ok := args[i].(float64)
	if !ok {
		return 0, fmt.Errorf("%s: argument %d must be a number, got %T", name, i+1, args[i])
	}
	if math.IsNaN(f) || math.IsInf(f, 0) || f != math.Trunc(f) {
		return 0, fmt.Errorf("%s: argument %d must be an integer, got %v", name, i+1, f)
	}
	return int(f), nil
}

func argFloat(name string, args []any, i int) (float64, error) {
	f, ok := args[i].(float64)
	if !ok {
		return 0, fmt.Errorf("%s: argument %d must be a number, got %T", name, i+1, args[i])
	}
	return f, nil
}

func argString(name string, args []any, i int) (string, error) {
	s, ok := args[i].(string)
	if !ok {
		return "", fmt.Errorf("%s: argument %d must be a string, got %T", name, i+1, args[i])
	}
	return s, nil
}

// optInt returns the integer at index i, or def if the argument is absent.
func optInt(name string, args []any, i, def int) (int, error) {
	if i >= len(args) {
		return def, nil
	}
	return argInt(name, args, i)
}

func boundDuration(ms float64) (any, error) {
	if math.IsNaN(ms) {
		return nil, fmt.Errorf("duration is not a number")
	}
	if ms >= float64(types.MaxDurationMillis) {
		return types.Duration{Millis: types.MaxDurationMillis}, nil
	}
	if ms <= float64(-types.MaxDurationMillis) {
		return types.Duration{Millis: -types.MaxDurationMillis}, nil
	}
	return types.Duration{Millis: int64(ms)}, nil
}

func clampDateFromFloat(ms float64) (any, error) {
	if math.IsNaN(ms) {
		return nil, fmt.Errorf("datetime is not a number")
	}
	if ms >= float64(types.MaxDateTimeMillis) {
		return types.DateTime{Millis: types.MaxDateTimeMillis}, nil
	}
	if ms <= float64(types.MinDateTimeMillis) {
		return types.DateTime{Millis: types.MinDateTimeMillis}, nil
	}
	return types.DateTime{Millis: int64(ms)}, nil
}

// ---- construction ----

func builtinDate(args ...any) (any, error) {
	if err := arity("date", args, 3, 3); err != nil {
		return nil, err
	}
	y, e1 := argInt("date", args, 0)
	m, e2 := argInt("date", args, 1)
	d, e3 := argInt("date", args, 2)
	if err := firstErr(e1, e2, e3); err != nil {
		return nil, err
	}
	ms, err := types.MillisFromComponents(y, m, d, 0, 0, 0, 0)
	if err != nil {
		return nil, err
	}
	return types.DateTime{Millis: ms}, nil
}

func builtinDatetime(args ...any) (any, error) {
	if err := arity("datetime", args, 3, 7); err != nil {
		return nil, err
	}
	y, e1 := argInt("datetime", args, 0)
	mo, e2 := argInt("datetime", args, 1)
	d, e3 := argInt("datetime", args, 2)
	h, e4 := optInt("datetime", args, 3, 0)
	mi, e5 := optInt("datetime", args, 4, 0)
	s, e6 := optInt("datetime", args, 5, 0)
	ml, e7 := optInt("datetime", args, 6, 0)
	if err := firstErr(e1, e2, e3, e4, e5, e6, e7); err != nil {
		return nil, err
	}
	ms, err := types.MillisFromComponents(y, mo, d, h, mi, s, ml)
	if err != nil {
		return nil, err
	}
	return types.DateTime{Millis: ms}, nil
}

func builtinTime(args ...any) (any, error) {
	if err := arity("time", args, 2, 4); err != nil {
		return nil, err
	}
	h, e1 := argInt("time", args, 0)
	mi, e2 := argInt("time", args, 1)
	s, e3 := optInt("time", args, 2, 0)
	ml, e4 := optInt("time", args, 3, 0)
	if err := firstErr(e1, e2, e3, e4); err != nil {
		return nil, err
	}
	ms, err := types.MillisFromComponents(1970, 1, 1, h, mi, s, ml)
	if err != nil {
		return nil, err
	}
	return types.DateTime{Millis: ms}, nil
}

func builtinParseDate(args ...any) (any, error) {
	if err := arity("parseDate", args, 1, 2); err != nil {
		return nil, err
	}
	s, err := argString("parseDate", args, 0)
	if err != nil {
		return nil, err
	}
	if len(args) == 2 {
		pattern, e := argString("parseDate", args, 1)
		if e != nil {
			return nil, e
		}
		return parseWithPattern("parseDate", pattern, s)
	}
	ms, perr := types.ParseISODateTime(s)
	if perr != nil {
		return nil, perr
	}
	return types.DateTime{Millis: ms}, nil
}

func builtinTryParseDate(args ...any) (any, error) {
	if err := arity("tryParseDate", args, 1, 2); err != nil {
		return nil, err
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, nil
	}
	if len(args) == 2 {
		pattern, ok := args[1].(string)
		if !ok {
			return nil, nil
		}
		v, perr := parseWithPattern("tryParseDate", pattern, s)
		if perr != nil {
			return nil, nil
		}
		return v, nil
	}
	ms, perr := types.ParseISODateTime(s)
	if perr != nil {
		return nil, nil
	}
	return types.DateTime{Millis: ms}, nil
}

// parseWithPattern parses s against a NITES pattern (or named layout) via the gotime engine, returning a
// canonical zoneless datetime. A pattern without an offset field is read as a UTC wall clock; one with an
// offset field yields the corresponding UTC instant. The result is bounded to the datetime range.
func parseWithPattern(name, pattern, s string) (any, error) {
	t, err := gotime.Parse(resolveLayout(pattern), s)
	if err != nil {
		return nil, fmt.Errorf("%s: cannot parse %q with pattern %q: %w", name, s, pattern, err)
	}
	ms := t.UnixMilli()
	if ms < types.MinDateTimeMillis || ms > types.MaxDateTimeMillis {
		return nil, fmt.Errorf("%s: parsed datetime %q is out of range", name, s)
	}
	return types.DateTime{Millis: ms}, nil
}

func builtinDuration(args ...any) (any, error) {
	if err := arity("duration", args, 2, 2); err != nil {
		return nil, err
	}
	amount, e1 := argFloat("duration", args, 0)
	unit, e2 := argString("duration", args, 1)
	if err := firstErr(e1, e2); err != nil {
		return nil, err
	}
	if unit == "month" || unit == "year" {
		return nil, fmt.Errorf("duration: %q is a calendar unit, not an exact duration; use addMonths/addYears", unit)
	}
	w, ok := exactUnits[unit]
	if !ok {
		return nil, fmt.Errorf("duration: unknown unit %q", unit)
	}
	return boundDuration(amount * float64(w))
}

// ---- epoch conversion ----

func builtinToEpochMillis(args ...any) (any, error) {
	if err := arity("toEpochMillis", args, 1, 1); err != nil {
		return nil, err
	}
	d, err := argDateTime("toEpochMillis", args, 0)
	if err != nil {
		return nil, err
	}
	return float64(d.Millis), nil
}

func builtinFromEpochMillis(args ...any) (any, error) {
	if err := arity("fromEpochMillis", args, 1, 1); err != nil {
		return nil, err
	}
	n, err := argFloat("fromEpochMillis", args, 0)
	if err != nil {
		return nil, err
	}
	return clampDateFromFloat(n)
}

func builtinToEpochSeconds(args ...any) (any, error) {
	if err := arity("toEpochSeconds", args, 1, 1); err != nil {
		return nil, err
	}
	d, err := argDateTime("toEpochSeconds", args, 0)
	if err != nil {
		return nil, err
	}
	return math.Floor(float64(d.Millis) / 1000), nil
}

func builtinFromEpochSeconds(args ...any) (any, error) {
	if err := arity("fromEpochSeconds", args, 1, 1); err != nil {
		return nil, err
	}
	n, err := argFloat("fromEpochSeconds", args, 0)
	if err != nil {
		return nil, err
	}
	return clampDateFromFloat(n * 1000)
}

// ---- calendar arithmetic / extraction ----

func builtinAddMonths(args ...any) (any, error) {
	return addCalendar("addMonths", args, false)
}

func builtinAddYears(args ...any) (any, error) {
	return addCalendar("addYears", args, true)
}

func addCalendar(name string, args []any, years bool) (any, error) {
	if err := arity(name, args, 2, 2); err != nil {
		return nil, err
	}
	d, e1 := argDateTime(name, args, 0)
	n, e2 := argInt(name, args, 1)
	if err := firstErr(e1, e2); err != nil {
		return nil, err
	}
	if years {
		return types.DateTime{Millis: types.AddYears(d.Millis, n)}, nil
	}
	return types.DateTime{Millis: types.AddMonths(d.Millis, n)}, nil
}

func builtinDiffMonths(args ...any) (any, error) {
	return diffCalendar("diffMonths", args, false)
}

func builtinDiffYears(args ...any) (any, error) {
	return diffCalendar("diffYears", args, true)
}

func diffCalendar(name string, args []any, years bool) (any, error) {
	if err := arity(name, args, 2, 2); err != nil {
		return nil, err
	}
	a, e1 := argDateTime(name, args, 0)
	b, e2 := argDateTime(name, args, 1)
	if err := firstErr(e1, e2); err != nil {
		return nil, err
	}
	if years {
		return float64(types.DiffYears(a.Millis, b.Millis)), nil
	}
	return float64(types.DiffMonths(a.Millis, b.Millis)), nil
}

func builtinDatePart(args ...any) (any, error) {
	if err := arity("datePart", args, 2, 3); err != nil {
		return nil, err
	}
	d, e1 := argDateTime("datePart", args, 0)
	comp, e2 := argString("datePart", args, 1)
	if err := firstErr(e1, e2); err != nil {
		return nil, err
	}
	offMin := 0
	if len(args) == 3 {
		off, e3 := argString("datePart", args, 2)
		if e3 != nil {
			return nil, e3
		}
		m, ok := types.ParseFixedOffset(off)
		if !ok {
			return nil, fmt.Errorf("datePart: invalid offset %q", off)
		}
		offMin = m
	}
	y, mo, day, h, mi, se, ml, wd := types.ComponentsOf(d.Millis + int64(offMin)*60000)
	switch comp {
	case "year":
		return float64(y), nil
	case "month":
		return float64(mo), nil
	case "day":
		return float64(day), nil
	case "hour":
		return float64(h), nil
	case "minute":
		return float64(mi), nil
	case "second":
		return float64(se), nil
	case "millisecond":
		return float64(ml), nil
	case "weekday":
		return float64(wd), nil
	default:
		return nil, fmt.Errorf("datePart: unknown component %q", comp)
	}
}

func builtinDurationIn(args ...any) (any, error) {
	if err := arity("durationIn", args, 2, 2); err != nil {
		return nil, err
	}
	d, e1 := argDuration("durationIn", args, 0)
	unit, e2 := argString("durationIn", args, 1)
	if err := firstErr(e1, e2); err != nil {
		return nil, err
	}
	if unit == "month" || unit == "year" {
		return nil, fmt.Errorf("durationIn: %q is a calendar unit, not an exact duration", unit)
	}
	w, ok := exactUnits[unit]
	if !ok {
		return nil, fmt.Errorf("durationIn: unknown unit %q", unit)
	}
	return float64(d.Millis) / float64(w), nil
}

// ---- formatting / parsing ----

func builtinFormatDate(args ...any) (any, error) {
	if err := arity("formatDate", args, 1, 3); err != nil {
		return nil, err
	}
	d, err := argDateTime("formatDate", args, 0)
	if err != nil {
		return nil, err
	}
	pattern := defaultLayout
	if len(args) >= 2 {
		p, e := argString("formatDate", args, 1)
		if e != nil {
			return nil, e
		}
		pattern = p
	}
	offMin := 0
	if len(args) == 3 {
		off, e := argString("formatDate", args, 2)
		if e != nil {
			return nil, e
		}
		m, ok := types.ParseFixedOffset(off)
		if !ok {
			return nil, fmt.Errorf("formatDate: invalid offset %q", off)
		}
		offMin = m
	}
	return gotime.Format(localTime(d.Millis, offMin), resolveLayout(pattern)), nil
}

func builtinFormatDur(args ...any) (any, error) {
	if err := arity("formatDur", args, 1, 1); err != nil {
		return nil, err
	}
	d, err := argDuration("formatDur", args, 0)
	if err != nil {
		return nil, err
	}
	return types.FormatISODuration(d.Millis), nil
}

func builtinParseDur(args ...any) (any, error) {
	if err := arity("parseDur", args, 1, 1); err != nil {
		return nil, err
	}
	s, err := argString("parseDur", args, 0)
	if err != nil {
		return nil, err
	}
	ms, perr := types.ParseISODuration(s)
	if perr != nil {
		return nil, perr
	}
	return types.Duration{Millis: ms}, nil
}

func builtinTryParseDur(args ...any) (any, error) {
	if err := arity("tryParseDur", args, 1, 1); err != nil {
		return nil, err
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, nil
	}
	ms, perr := types.ParseISODuration(s)
	if perr != nil {
		return nil, nil
	}
	return types.Duration{Millis: ms}, nil
}

func firstErr(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}
