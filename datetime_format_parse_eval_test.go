package uexl_test

import (
	"context"
	"testing"
)

// TestFormatDate_GotimeDialect exercises formatDate across the NITES specifiers as rendered by the
// gotime engine (datetime-spec §10.1): 24-hour is hhhh, 12-hour is hh/h, minute is ii, AM/PM is aa/a,
// weekday is www/wwww, and numeric offsets are o/oo/ooo. It replaces the unit coverage that previously
// lived in types.FormatNITES (now retired in favour of gotime).
func TestFormatDate_GotimeDialect(t *testing.T) {
	trueExprs := []string{
		`formatDate(d"2026-12-21T15:04:05Z", "hhhh:ii:ss") == "15:04:05"`, // 24-hour
		`formatDate(d"2026-12-21T15:04:05Z", "hh:ii aa") == "03:04 PM"`,   // 12-hour + AM/PM
		`formatDate(d"2026-12-21T15:04:05Z", "h:i:s a") == "3:4:5 pm"`,    // no-pad + lowercase am/pm
		`formatDate(d"2026-12-21T15:04:05Z", "mmm") == "Dec"`,
		`formatDate(d"2026-12-21T15:04:05Z", "mmmm") == "December"`,
		`formatDate(d"2024-12-02", "www") == "Mon"`, // 2024-12-02 is a Monday
		`formatDate(d"2024-12-02", "wwww") == "Monday"`,
		`formatDate(d"2026-12-21T15:04:05Z", "yy") == "26"`, // 2-digit year
		`formatDate(d"2026-12-21T15:04:05Z", "sql") == "2026-12-21 15:04:05"`,
		// offset shifts the wall clock: 15:04 UTC seen at +05:30 is 20:34
		`formatDate(d"2026-12-21T15:04:05Z", "hhhh:ii", "+05:30") == "20:34"`,
		`formatDate(d"2026-12-21T15:04:05Z", "ooo", "+05:30") == "+05:30"`, // offset with colon
		`formatDate(d"2026-12-21T15:04:05Z", "oo", "+05:30") == "+0530"`,   // offset, no colon
	}
	env := dtEnv()
	for _, expr := range trueExprs {
		v, err := env.Eval(context.Background(), expr, nil)
		if err != nil || v != true {
			t.Errorf("%s: got %v, err %v", expr, v, err)
		}
	}
	for _, expr := range []string{
		`formatDate(d"2024-01-01", "yyyy", "not-an-offset")`, // invalid offset string
		`formatDate(d"2024-01-01", "yyyy", 5)`,               // offset must be a string
	} {
		if _, err := env.Eval(context.Background(), expr, nil); err == nil {
			t.Errorf("%s: expected an error", expr)
		}
	}
}

// TestParseDate_Pattern covers the gotime-backed pattern-directed parseDate(s, pattern) / tryParseDate
// (datetime-spec §10.2) — the capability that motivated wiring gotime. A pattern without an offset field
// reads the value as a UTC wall clock.
func TestParseDate_Pattern(t *testing.T) {
	trueExprs := []string{
		`parseDate("2026-12-21", "yyyy-mm-dd") == d"2026-12-21"`,
		`parseDate("21/12/2026", "dd/mm/yyyy") == d"2026-12-21"`,
		`parseDate("Dec 21, 2026", "mmm dd, yyyy") == d"2026-12-21"`,
		`parseDate("2026-12-21 15:04", "yyyy-mm-dd hhhh:ii") == d"2026-12-21T15:04:00Z"`,
		`parseDate("2026-12-21", "isodate") == d"2026-12-21"`, // named layout
		// round-trip through the gotime engine
		`parseDate(formatDate(d"2026-12-21T15:04:05Z", "yyyy-mm-dd hhhh:ii:ss"), "yyyy-mm-dd hhhh:ii:ss") == d"2026-12-21T15:04:05Z"`,
		`tryParseDate("2026-12-21", "yyyy-mm-dd") == d"2026-12-21"`,
		`tryParseDate("not-a-date", "yyyy-mm-dd") == null`,
		`tryParseDate(5, "yyyy-mm-dd") == null`, // non-string value -> null, not error
		`tryParseDate("2024-01-01", 5) == null`, // non-string pattern -> null, not error
	}
	env := dtEnv()
	for _, expr := range trueExprs {
		v, err := env.Eval(context.Background(), expr, nil)
		if err != nil || v != true {
			t.Errorf("%s: got %v, err %v", expr, v, err)
		}
	}
	for _, expr := range []string{
		`parseDate("not-a-date", "yyyy-mm-dd")`, // strict: error
		`parseDate("2026-12-21", 5)`,            // pattern must be a string
	} {
		if _, err := env.Eval(context.Background(), expr, nil); err == nil {
			t.Errorf("%s: expected an error", expr)
		}
	}
}
