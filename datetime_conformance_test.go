package uexl_test

import (
	"context"
	"testing"

	"github.com/maniartech/uexl"
	"github.com/maniartech/uexl/types"
)

// This is the authoritative executable form of datetime-spec.md §12. Each row is one conformance case,
// traceable by ID. The env has the datetime library attached (WithDatetime); the core types, literals,
// and operators are always present.

// conformCase is a §12 row. want is matched type-aware: an int64 matches a datetime/duration ms or a
// numeric result; bool/string match directly; nil means the result must be null.
type conformCase struct {
	id   string
	expr string
	want any
}

func conformMatch(got, want any) bool {
	switch w := want.(type) {
	case int64:
		switch g := got.(type) {
		case types.DateTime:
			return g.Millis == w
		case types.Duration:
			return g.Millis == w
		case float64:
			return int64(g) == w
		}
		return false
	case bool:
		return got == w
	case string:
		return got == w
	case nil:
		return got == nil
	}
	return false
}

func TestConformance_Section12(t *testing.T) {
	env := uexl.DefaultWith(uexl.WithDatetime())
	cases := []conformCase{
		{"dt-lit-001", `d"1970-01-01T00:00:00Z"`, int64(0)},
		{"dt-lit-002", `d"2024-12-01"`, int64(1733011200000)},
		{"dt-lit-003", `d"2024-12-01T10:30:00+05:30" == d"2024-12-01T05:00:00Z"`, true},
		{"dt-lit-004", `d"1969-12-31T00:00:00Z"`, int64(-86400000)},
		{"dt-lit-005", `d"0001-01-01T00:00:00Z"`, int64(-62135596800000)},
		{"dt-op-001", `d"2024-12-02" - d"2024-12-01"`, int64(86400000)},
		{"dt-op-002", `d"2024-12-01" + duration(2, "day")`, int64(1733184000000)},
		{"dt-op-003", `durationIn(d"2024-12-02" - d"2024-12-01", "hour")`, int64(24)},
		{"dt-op-004", `d"2024-12-01" < d"2024-12-02"`, true},
		{"dt-lit-006", `d"2024-12-01" + 7d`, int64(1733616000000)},
		{"dt-lit-007", `30ms == duration(30, "millisecond")`, true},
		{"dt-lit-008", `1.5h == duration(90, "minute")`, true},
		{"dt-lit-009", `-7d == duration(-7, "day")`, true},
		{"dt-cal-001", `addMonths(d"2024-01-31", 1) == d"2024-02-29"`, true},
		{"dt-cal-002", `addYears(d"2024-02-29", 1) == d"2025-02-28"`, true},
		{"dt-cal-003", `addMonths(d"2024-03-15", -2) == d"2024-01-15"`, true},
		{"dt-cal-004", `datePart(d"2026-12-21", "month")`, int64(12)},
		{"dt-cal-005", `diffMonths(d"2024-01-15", d"2024-03-10")`, int64(1)},
		{"dt-cal-006", `diffYears(d"2020-06-01", d"2024-05-01")`, int64(3)},
		{"dt-div-001", `duration(1, "hour") / duration(15, "minute")`, int64(4)},
		{"dt-div-002", `duration(1, "hour") / 2 == duration(30, "minute")`, true},
		{"dt-fmt-001", `formatDate(d"2026-12-21T15:04:05Z", "yyyy-mm-dd")`, "2026-12-21"},
		{"dt-fmt-002", `formatDate(d"2026-12-21T15:04:05Z")`, "2026-12-21T15:04:05"},
		{"dt-fmt-003", `formatDur(duration(90, "minute"))`, "PT1H30M"},
		{"dt-fmt-004", `parseDur("PT1H30M") == duration(90, "minute")`, true},
		{"dt-ctor-001", `date(2024, 12, 1) == d"2024-12-01"`, true},
		{"dt-ctor-002", `datetime(2024, 12, 1, 10, 30) == d"2024-12-01T10:30:00Z"`, true},
		{"dt-ctor-003", `time(10, 30) == d"1970-01-01T10:30:00Z"`, true},
		{"dt-epoch-001", `toEpochMillis(d"2024-12-01")`, int64(1733011200000)},
		{"dt-epoch-002", `fromEpochMillis(0) == d"1970-01-01T00:00:00Z"`, true},
		{"dt-epoch-003", `toEpochSeconds(d"1970-01-01T00:00:01Z")`, int64(1)},
		{"dt-off-001", `formatDate(d"2024-12-01T00:30:00Z", "yyyy-mm-dd", "+05:30")`, "2024-12-01"},
		{"dt-off-002", `datePart(d"2024-12-01T00:30:00Z", "hour", "+05:30")`, int64(6)},
		{"dt-off-003", `datePart(d"2024-12-02", "weekday")`, int64(1)},
		{"dt-clamp-001", `addYears(d"9999-06-01", 5) == d"9999-12-31T23:59:59.999Z"`, true},
		{"dt-clamp-002", `date(12000, 1, 1) == d"9999-12-31T23:59:59.999Z"`, true},
		{"dt-parse-001", `tryParseDate("not-a-date")`, nil},
		{"dt-parse-002", `tryParseDur("nonsense")`, nil},
		{"dt-parse-003", `tryParseDate("2024-12-01") == d"2024-12-01"`, true},
	}
	for _, c := range cases {
		got, err := env.Eval(context.Background(), c.expr, nil)
		if err != nil {
			t.Errorf("[%s] %s: unexpected error: %v", c.id, c.expr, err)
			continue
		}
		if !conformMatch(got, c.want) {
			t.Errorf("[%s] %s: got %v (%T), want %v", c.id, c.expr, got, got, c.want)
		}
	}

	// §12 error cases (dt-err-001..016).
	errCases := map[string]string{
		"dt-err-001": `d"2024-13-01"`,
		"dt-err-002": `d"2024-02-30"`,
		"dt-err-003": `d"2024-12-01T10:30:60Z"`,
		"dt-err-004": `d"0000-01-01"`,
		"dt-err-005": `d"10000-01-01"`,
		"dt-err-006": `duration(2, "month")`,
		"dt-err-007": `d"2024-12-01" + 7`,
		"dt-err-008": `d"2024-12-01" + d"2024-12-02"`,
		"dt-err-009": `7d"2024-12-01"`,
		"dt-err-010": `parseDur("P1Y")`,
		"dt-err-011": `parseDur("P3M")`,
		"dt-err-012": `duration(2, "months")`,
		"dt-err-013": `date(2023, 2, 29)`,
		"dt-err-014": `date(2024, 13, 1)`,
		"dt-err-015": `datetime(2024, 1, 1, 25, 0)`,
		"dt-err-016": `parseDate("not-a-date")`,
	}
	for id, expr := range errCases {
		if _, err := env.Eval(context.Background(), expr, nil); err == nil {
			t.Errorf("[%s] %s: expected an error", id, expr)
		}
	}
}
