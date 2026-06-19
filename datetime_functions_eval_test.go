package uexl_test

import (
	"context"
	"testing"

	"github.com/maniartech/uexl"
)

func dtEnv() *uexl.Env {
	return uexl.DefaultWith(uexl.WithDatetime())
}

func evalDT(t *testing.T, expr string) any {
	t.Helper()
	v, err := dtEnv().Eval(context.Background(), expr, nil)
	if err != nil {
		t.Fatalf("%s: unexpected error: %v", expr, err)
	}
	return v
}

// TestDatetimeFunctions_True covers the §12 conformance cases that evaluate to true (mostly `X == Y`).
func TestDatetimeFunctions_True(t *testing.T) {
	trueExprs := []string{
		// construction
		`date(2024, 12, 1) == d"2024-12-01"`,                       // dt-ctor-001
		`datetime(2024, 12, 1, 10, 30) == d"2024-12-01T10:30:00Z"`, // dt-ctor-002
		`time(10, 30) == d"1970-01-01T10:30:00Z"`,                  // dt-ctor-003
		`parseDate("2024-12-01") == d"2024-12-01"`,
		`tryParseDate("not-a-date") == null`,
		`duration(2, "day") == 2d`,
		// epoch
		`toEpochMillis(d"2024-12-01") == 1733011200000`, // dt-epoch-001
		`fromEpochMillis(0) == d"1970-01-01T00:00:00Z"`, // dt-epoch-002
		`toEpochSeconds(d"1970-01-01T00:00:01Z") == 1`,  // dt-epoch-003
		`fromEpochSeconds(0) == d"1970-01-01T00:00:00Z"`,
		// calendar
		`addMonths(d"2024-01-31", 1) == d"2024-02-29"`,             // dt-cal-001
		`addYears(d"2024-02-29", 1) == d"2025-02-28"`,              // dt-cal-002
		`addMonths(d"2024-03-15", -2) == d"2024-01-15"`,            // dt-cal-003
		`datePart(d"2026-12-21", "month") == 12`,                   // dt-cal-004
		`diffMonths(d"2024-01-15", d"2024-03-10") == 1`,            // dt-cal-005
		`diffYears(d"2020-06-01", d"2024-05-01") == 3`,             // dt-cal-006
		`datePart(d"2024-12-02", "weekday") == 1`,                  // dt-off-003 (Monday)
		`datePart(d"2024-12-01T00:30:00Z", "hour", "+05:30") == 6`, // dt-off-002
		// division / durationIn
		`durationIn(d"2024-12-02" - d"2024-12-01", "hour") == 24`, // dt-op-003
		`duration(1, "hour") / duration(15, "minute") == 4`,       // dt-div-001
		`duration(1, "hour") / 2 == duration(30, "minute")`,       // dt-div-002
		// formatting
		`formatDate(d"2026-12-21T15:04:05Z", "yyyy-mm-dd") == "2026-12-21"`, // dt-fmt-001
		`formatDate(d"2026-12-21T15:04:05Z") == "2026-12-21T15:04:05"`,      // dt-fmt-002
		`formatDur(duration(90, "minute")) == "PT1H30M"`,                    // dt-fmt-003
		`parseDur("PT1H30M") == duration(90, "minute")`,                     // dt-fmt-004
		`tryParseDur("nonsense") == null`,
		`tryParseDate("2024-12-01") == d"2024-12-01"`,
		`tryParseDur("PT1H") == duration(1, "hour")`,
		`tryParseDate(5) == null`, // non-string -> null, not error
		`tryParseDur(5) == null`,
		`formatDate(d"2024-12-01T00:30:00Z", "yyyy-mm-dd", "+05:30") == "2024-12-01"`, // dt-off-001
		// datePart over every component
		`datePart(d"2024-12-01T10:30:45.250Z", "year") == 2024`,
		`datePart(d"2024-12-01T10:30:45.250Z", "day") == 1`,
		`datePart(d"2024-12-01T10:30:45.250Z", "hour") == 10`,
		`datePart(d"2024-12-01T10:30:45.250Z", "minute") == 30`,
		`datePart(d"2024-12-01T10:30:45.250Z", "second") == 45`,
		`datePart(d"2024-12-01T10:30:45.250Z", "millisecond") == 250`,
		// clamp / saturation
		`date(12000, 1, 1) == d"9999-12-31T23:59:59.999Z"`,                     // dt-clamp-002
		`addYears(d"9999-06-01", 5) == d"9999-12-31T23:59:59.999Z"`,            // dt-clamp-001
		`durationIn(duration(1e300, "day"), "millisecond") == 315537897599999`, // boundDuration saturates
		`fromEpochMillis(1e300) == d"9999-12-31T23:59:59.999Z"`,                // clamp MAX
		`fromEpochMillis(-1e300) == d"0001-01-01T00:00:00Z"`,                   // clamp MIN
	}
	for _, expr := range trueExprs {
		if got := evalDT(t, expr); got != true {
			t.Errorf("%s: got %v (%T), want true", expr, got, got)
		}
	}
}

// TestDatetimeFunctions_Errors covers the §12 conformance error cases.
func TestDatetimeFunctions_Errors(t *testing.T) {
	for _, expr := range []string{
		`duration(2, "month")`,                 // dt-err-006: calendar unit on duration
		`duration(1, "year")`,                  // calendar unit
		`duration(1, "fortnight")`,             // unknown unit
		`date(2023, 2, 29)`,                    // dt-err-013: invalid day (non-leap)
		`date(2024, 13, 1)`,                    // dt-err-014: invalid month
		`datetime(2024, 1, 1, 25, 0)`,          // dt-err-015: invalid hour
		`parseDate("not-a-date")`,              // dt-err-016: strict parse errors
		`parseDur("P1Y")`,                      // dt-err-010: year not exact
		`parseDur("P3M")`,                      // dt-err-011: date-part month not exact
		`durationIn(7d, "month")`,              // calendar unit on durationIn
		`datePart(d"2024-12-01", "fortnight")`, // unknown component
		`datePart(d"2024-12-01", "hour", "bad-offset")`,
		`date(2024, 12)`, // wrong arity
		// argument type mismatches (cover the arg-helper error branches)
		`date("x", 1, 1)`,
		`datetime(2024, 1, 1, "x")`,
		`time("x", 0)`,
		`toEpochMillis(5)`,
		`toEpochSeconds("x")`,
		`fromEpochMillis("x")`,
		`fromEpochMillis(NaN)`,
		`fromEpochSeconds("x")`,
		`fromEpochSeconds(NaN)`,
		`addMonths(5, 1)`,
		`addMonths(d"2024-01-01", 1.5)`, // non-integer
		`diffMonths(5, d"2024-01-01")`,
		`durationIn(5, "hour")`,
		`duration("x", "day")`,
		`formatDate(5)`,
		`formatDate(d"2024-01-01", 5)`,
		`formatDur(5)`,
		`parseDur(5)`,
		`parseDate(5)`,
	} {
		if _, err := dtEnv().Eval(context.Background(), expr, nil); err == nil {
			t.Errorf("%s: expected an error", expr)
		}
	}
}
