package types

import "testing"

func ms(t *testing.T, s string) int64 {
	t.Helper()
	v, err := ParseISODateTime(s)
	if err != nil {
		t.Fatalf("ParseISODateTime(%q): %v", s, err)
	}
	return v
}

func TestComponentsOf_RoundTrip(t *testing.T) {
	for _, s := range []string{
		"1970-01-01T00:00:00.000Z", "2024-12-01T10:30:45.250Z", "0001-01-01T00:00:00.000Z",
		"9999-12-31T23:59:59.999Z", "1969-12-31T23:59:59.999Z", "2024-02-29T00:00:00.000Z",
	} {
		want := ms(t, s)
		y, mo, d, h, mi, se, ml, _ := ComponentsOf(want)
		got, err := MillisFromComponents(y, mo, d, h, mi, se, ml)
		if err != nil {
			t.Errorf("%s: recompose error %v", s, err)
			continue
		}
		if got != want {
			t.Errorf("%s: round-trip %d != %d (components %d-%d-%d %d:%d:%d.%d)", s, got, want, y, mo, d, h, mi, se, ml)
		}
	}
}

func TestComponentsOf_Weekday(t *testing.T) {
	// 2024-12-01 is a Sunday (ISO 7), 2024-12-02 a Monday (ISO 1).
	if _, _, _, _, _, _, _, wd := ComponentsOf(ms(t, "2024-12-01")); wd != 7 {
		t.Errorf("2024-12-01 weekday = %d, want 7", wd)
	}
	if _, _, _, _, _, _, _, wd := ComponentsOf(ms(t, "2024-12-02")); wd != 1 {
		t.Errorf("2024-12-02 weekday = %d, want 1", wd)
	}
}

func TestMillisFromComponents(t *testing.T) {
	if got, err := MillisFromComponents(2024, 12, 1, 0, 0, 0, 0); err != nil || got != 1733011200000 {
		t.Errorf("date(2024,12,1) = (%d,%v), want 1733011200000", got, err)
	}
	// Year out of range clamps (datetime-spec §5.4).
	if got, err := MillisFromComponents(12000, 1, 1, 0, 0, 0, 0); err != nil || got != MaxDateTimeMillis {
		t.Errorf("date(12000,1,1) = (%d,%v), want MAX", got, err)
	}
	if got, err := MillisFromComponents(0, 6, 1, 0, 0, 0, 0); err != nil || got != MinDateTimeMillis {
		t.Errorf("date(0,6,1) = (%d,%v), want MIN", got, err)
	}
	// Structurally invalid components error.
	for _, c := range []struct {
		y, mo, d, h, mi, s, ml int
	}{
		{2023, 2, 29, 0, 0, 0, 0}, // not a leap year
		{2024, 13, 1, 0, 0, 0, 0}, // month
		{2024, 4, 31, 0, 0, 0, 0}, // April has 30 days
		{2024, 1, 1, 25, 0, 0, 0}, // hour
		{2024, 1, 1, 0, 60, 0, 0}, // minute
		{2024, 1, 1, 0, 0, 60, 0}, // second
		{2024, 1, 1, 0, 0, 0, 1000},
	} {
		if got, err := MillisFromComponents(c.y, c.mo, c.d, c.h, c.mi, c.s, c.ml); err == nil {
			t.Errorf("components %+v: expected error, got %d", c, got)
		}
	}
}

func TestAddMonthsYears(t *testing.T) {
	cases := []struct {
		in   string
		n    int
		want string
		year bool
	}{
		{"2024-01-31", 1, "2024-02-29", false}, // end-of-month clamp (leap)
		{"2023-01-31", 1, "2023-02-28", false},
		{"2024-03-31", 1, "2024-04-30", false},
		{"2024-03-15", -2, "2024-01-15", false}, // dt-cal-003
		{"2024-12-15", 1, "2025-01-15", false},  // year rollover
		{"2024-02-29", 1, "2025-02-28", true},   // addYears, leap -> non-leap
	}
	for _, c := range cases {
		var got int64
		if c.year {
			got = AddYears(ms(t, c.in), c.n)
		} else {
			got = AddMonths(ms(t, c.in), c.n)
		}
		if got != ms(t, c.want) {
			t.Errorf("add(%s, %d) = %d, want %s", c.in, c.n, got, c.want)
		}
	}
}

func TestDiffMonthsYears(t *testing.T) {
	mcases := []struct {
		a, b string
		want int
	}{
		{"2024-01-15", "2024-03-10", 1},  // dt-cal-005
		{"2024-01-31", "2024-02-29", 1},  // clamped target reached
		{"2024-01-31", "2024-02-28", 0},  // clamped target not reached
		{"2024-03-10", "2024-01-15", -1}, // negative (truncated toward zero)
		{"2024-06-01", "2024-06-01", 0},
	}
	for _, c := range mcases {
		if got := DiffMonths(ms(t, c.a), ms(t, c.b)); got != c.want {
			t.Errorf("diffMonths(%s,%s) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
	if got := DiffYears(ms(t, "2020-06-01"), ms(t, "2024-05-01")); got != 3 { // dt-cal-006
		t.Errorf("diffYears = %d, want 3", got)
	}
	if got := DiffYears(ms(t, "2024-05-01"), ms(t, "2020-06-01")); got != -3 {
		t.Errorf("diffYears (negative) = %d, want -3", got)
	}
}

func TestFloorHelpers(t *testing.T) {
	if floorDivInt(-7, 3) != -3 || floorModInt(-7, 3) != 2 {
		t.Errorf("floor(-7,3) = (%d,%d), want (-3,2)", floorDivInt(-7, 3), floorModInt(-7, 3))
	}
	if floorDivInt(7, 3) != 2 || floorModInt(7, 3) != 1 {
		t.Errorf("floor(7,3) = (%d,%d), want (2,1)", floorDivInt(7, 3), floorModInt(7, 3))
	}
	if q, r := floorDivMod(-1, msPerDay); q != -1 || r != msPerDay-1 {
		t.Errorf("floorDivMod(-1) = (%d,%d)", q, r)
	}
}

func TestAddMonths_Clamp(t *testing.T) {
	// Shifting past the upper bound clamps to MAX.
	if got := AddMonths(ms(t, "9999-12-01"), 12); got != MaxDateTimeMillis {
		t.Errorf("AddMonths upper clamp = %d, want MAX", got)
	}
	if got := AddYears(ms(t, "0001-06-01"), -5); got != MinDateTimeMillis {
		t.Errorf("AddYears lower clamp = %d, want MIN", got)
	}
}

func TestParseFixedOffset(t *testing.T) {
	cases := map[string]struct {
		mins int
		ok   bool
	}{
		"Z":      {0, true},
		"+05:30": {330, true},
		"-08:00": {-480, true},
		"+00:00": {0, true},
		"+99:00": {0, false},
		"5:30":   {0, false},
	}
	for in, want := range cases {
		got, ok := ParseFixedOffset(in)
		if ok != want.ok || (ok && got != want.mins) {
			t.Errorf("ParseFixedOffset(%q) = (%d,%v), want (%d,%v)", in, got, ok, want.mins, want.ok)
		}
	}
}
