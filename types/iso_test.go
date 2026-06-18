package types

import "testing"

func TestParseISODateTime_Valid(t *testing.T) {
	cases := []struct {
		in   string
		want int64
	}{
		{"1970-01-01T00:00:00Z", 0},
		{"2024-12-01", 1733011200000},                   // date-only -> midnight UTC
		{"2024-12-01T05:00:00Z", 1733029200000},         // explicit Z
		{"2024-12-01T05:00:00", 1733029200000},          // missing offset -> UTC
		{"2024-12-01T10:30:00Z", 1733049000000},         // HH:MM:SS
		{"2024-12-01T10:30", 1733049000000},             // HH:MM only, seconds default 0
		{"1969-12-31T00:00:00Z", -86400000},             // pre-epoch (negative)
		{"0001-01-01T00:00:00Z", MinDateTimeMillis},     // MIN boundary
		{"9999-12-31T23:59:59.999Z", MaxDateTimeMillis}, // MAX boundary
		{"2024-12-01T00:00:00.250Z", 1733011200250},     // fractional ms
		{"2024-02-29T00:00:00Z", 1709164800000},         // valid leap day
	}
	for _, c := range cases {
		got, err := ParseISODateTime(c.in)
		if err != nil {
			t.Errorf("%s: unexpected error %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("%s: got %d, want %d", c.in, got, c.want)
		}
	}
}

// TestParseISODateTime_OffsetEquivalence covers §4.2: an offset is applied then discarded, so the two
// spellings denote the SAME instant.
func TestParseISODateTime_OffsetEquivalence(t *testing.T) {
	a, err1 := ParseISODateTime("2024-12-01T10:30:00+05:30")
	b, err2 := ParseISODateTime("2024-12-01T05:00:00Z")
	if err1 != nil || err2 != nil {
		t.Fatalf("unexpected errors: %v %v", err1, err2)
	}
	if a != b {
		t.Errorf("+05:30 form = %d, Z form = %d; want equal", a, b)
	}
	// And a negative offset shifts the other way.
	c, _ := ParseISODateTime("2024-12-01T00:00:00-05:00")
	d, _ := ParseISODateTime("2024-12-01T05:00:00Z")
	if c != d {
		t.Errorf("-05:00 form = %d, want %d", c, d)
	}
}

// TestParseISODateTime_FractionTruncation covers §4.3: sub-millisecond digits truncate toward zero.
func TestParseISODateTime_FractionTruncation(t *testing.T) {
	base := int64(1733011200000)
	cases := map[string]int64{
		"2024-12-01T00:00:00.2Z":      base + 200,
		"2024-12-01T00:00:00.25Z":     base + 250,
		"2024-12-01T00:00:00.250Z":    base + 250,
		"2024-12-01T00:00:00.2509Z":   base + 250, // 4th digit discarded
		"2024-12-01T00:00:00.250999Z": base + 250,
	}
	for in, want := range cases {
		got, err := ParseISODateTime(in)
		if err != nil {
			t.Errorf("%s: unexpected error %v", in, err)
			continue
		}
		if got != want {
			t.Errorf("%s: got %d, want %d", in, got, want)
		}
	}
}

func TestParseISODateTime_Errors(t *testing.T) {
	bad := []string{
		"2024-13-01",                     // month out of range
		"2024-02-30",                     // day out of range
		"2023-02-29",                     // not a leap year
		"2024-12-01T25:00:00Z",           // hour out of range
		"2024-12-01T10:60:00Z",           // minute out of range
		"2024-12-01T10:30:60Z",           // leap second rejected
		"0000-01-01",                     // year 0000
		"10000-01-01",                    // 5-digit year -> format error
		"2024-12-01 10:30:00",            // space separator (only 'T' accepted)
		"2024-12-01T",                    // 'T' with no time
		"2024-12-01Tgarbage",             // malformed time
		"2024-12-01T10:30:00+05:30extra", // trailing characters
		"2024-12-01T10:30:00.Z",          // empty fraction
		"24-12-01",                       // 2-digit year
		"2024/12/01",                     // wrong separators
		"",                               // empty
		"9999-12-31T23:59:59.999-01:00",  // in-range local, out-of-range UTC after offset
		"0001-01-01T00:00:00+05:00",      // underflow below MIN after offset
	}
	for _, in := range bad {
		if got, err := ParseISODateTime(in); err == nil {
			t.Errorf("%q: expected error, got %d", in, got)
		}
	}
}

func TestParseISODateTime_ZeroAllocOnSuccess(t *testing.T) {
	if allocs := testing.AllocsPerRun(1000, func() {
		_, _ = ParseISODateTime("2024-12-01T10:30:00.250Z")
	}); allocs != 0 {
		t.Errorf("ParseISODateTime allocated %v times on success, want 0", allocs)
	}
}

func BenchmarkParseISODateTime(b *testing.B) {
	var sink int64
	for i := 0; i < b.N; i++ {
		sink, _ = ParseISODateTime("2024-12-01T10:30:00.250Z")
	}
	_ = sink
}

func TestParseISODateTime_MoreErrors(t *testing.T) {
	bad := []string{
		"2024-04-31",                // April has 30 days
		"2024-06-31",                // June has 30 days
		"2024-12-01T10:30:00+99:00", // offset hour out of range
		"2024-12-01T10:30:00+05:60", // offset minute out of range
		"2024-12-01T10:30:00X",      // unrecognized offset token
		"2024-12-01T10:30:00+5:00",  // malformed offset (one-digit hour)
	}
	for _, in := range bad {
		if got, err := ParseISODateTime(in); err == nil {
			t.Errorf("%q: expected error, got %d", in, got)
		}
	}
}

// TestParseISODateTime_ThirtyDayMonth exercises the 30-day branch of daysInMonth on the success path.
func TestParseISODateTime_ThirtyDayMonth(t *testing.T) {
	for _, in := range []string{"2024-04-30", "2024-06-30", "2024-09-30", "2024-11-30"} {
		if _, err := ParseISODateTime(in); err != nil {
			t.Errorf("%q: unexpected error %v", in, err)
		}
	}
}

// TestISO_DefensiveBranches covers the internal guards that the public callers never trigger.
func TestISO_DefensiveBranches(t *testing.T) {
	if _, ok := parseDigitsN(""); ok {
		t.Error(`parseDigitsN("") should return ok=false`)
	}
	if daysInMonth(2024, 0) != 0 || daysInMonth(2024, 13) != 0 {
		t.Error("daysInMonth with an invalid month should return 0")
	}
}

func TestParseISODateTime_SecondsFieldEdges(t *testing.T) {
	for _, in := range []string{
		"2024-12-01T10:30:0",  // incomplete seconds field
		"2024-12-01T10:30:ab", // non-numeric seconds field
		"2024-12-01T10:3",     // incomplete minute field (< HH:MM)
	} {
		if _, err := ParseISODateTime(in); err == nil {
			t.Errorf("%q: expected error", in)
		}
	}
}

// TestDaysFromCivil_NegativeYears covers the proleptic-Gregorian algorithm for years <= 0. UExL clamps
// datetimes to >= 0001, so this branch is unreachable from the public parser, but the algorithm must
// still be correct (the datetime library may reuse it).
func TestDaysFromCivil_NegativeYears(t *testing.T) {
	if daysFromCivil(-1, 12, 31) >= daysFromCivil(1, 1, 1) {
		t.Error("a year -1 date should precede 0001-01-01")
	}
	if daysFromCivil(0, 1, 1) >= daysFromCivil(1, 1, 1) {
		t.Error("year 0 should precede year 1")
	}
}
