package types

import "testing"

func TestFormatNITES(t *testing.T) {
	dt := ms(t, "2026-12-21T15:04:05Z")
	cases := []struct {
		pattern string
		off     int
		want    string
	}{
		{"yyyy-mm-dd", 0, "2026-12-21"},   // dt-fmt-001
		{"iso", 0, "2026-12-21T15:04:05"}, // dt-fmt-002 (default layout)
		{"yyyy-mm-ddThhh:ii:ss", 0, "2026-12-21T15:04:05"},
		{"mmm", 0, "Dec"},
		{"mmmm", 0, "December"},
		{"hh:ii aa", 0, "03:04 PM"},       // 12-hour + AM/PM
		{"h:i:s", 0, "3:4:5"},             // no-pad
		{"yyyy-mm-dd", 330, "2026-12-21"}, // offset shifts the wall clock but date holds here
		{"hhh:ii", 330, "20:34"},          // 15:04 UTC at +05:30 -> 20:34
		{"yy", 0, "26"},                   // 2-digit year
		{"m/d", 0, "12/21"},               // no-pad month/day
		{"a", 0, "pm"},                    // lowercase am/pm
		{"z", 0, "Z"},                     // zone, UTC
		{"z", 330, "+05:30"},              // zone with offset
		{"o", 330, "+0530"},               // offset, no colon
		{"ooo", 330, "+05:30"},            // offset with colon
	}
	for _, c := range cases {
		if got := FormatNITES(dt, c.off, c.pattern); got != c.want {
			t.Errorf("FormatNITES(%q, off=%d) = %q, want %q", c.pattern, c.off, got, c.want)
		}
	}
	// Weekday names (2024-12-02 is a Monday).
	if got := FormatNITES(ms(t, "2024-12-02"), 0, "www"); got != "Mon" {
		t.Errorf("www = %q, want Mon", got)
	}
	if got := FormatNITES(ms(t, "2024-12-02"), 0, "wwww"); got != "Monday" {
		t.Errorf("wwww = %q, want Monday", got)
	}
}

func TestFormatISODuration(t *testing.T) {
	cases := map[int64]string{
		5400000:   "PT1H30M", // dt-fmt-003
		604800000: "P7D",
		0:         "PT0S",
		500:       "PT0.5S",
		1000:      "PT1S",
		-5400000:  "-PT1H30M",
		90061000:  "P1DT1H1M1S",
		1500:      "PT1.5S",
	}
	for ms, want := range cases {
		if got := FormatISODuration(ms); got != want {
			t.Errorf("FormatISODuration(%d) = %q, want %q", ms, got, want)
		}
	}
}

func TestParseISODuration(t *testing.T) {
	ok := map[string]int64{
		"PT1H30M":    5400000, // dt-fmt-004
		"P7D":        604800000,
		"PT0.5S":     500,
		"P1W":        604800000,
		"PT0S":       0,
		"-PT1H":      -3600000,
		"+PT1H":      3600000,
		"P1DT1H1M1S": 90061000,
	}
	for s, want := range ok {
		if got, err := ParseISODuration(s); err != nil || got != want {
			t.Errorf("ParseISODuration(%q) = (%d,%v), want %d", s, got, err, want)
		}
	}
	bad := []string{
		"P1Y",           // dt-err-010: year not exact
		"P3M",           // dt-err-011: date-part month not exact
		"P1Y2M",         // year + month
		"1H",            // missing P
		"P",             // empty
		"PT",            // empty time
		"PT1X",          // bad designator
		"PT1.5",         // number without designator
		"PTH",           // designator without number
		"PT9999999999H", // out of range
		"P99999999999D", // out of range (date part)
	}
	for _, s := range bad {
		if got, err := ParseISODuration(s); err == nil {
			t.Errorf("ParseISODuration(%q): expected error, got %d", s, got)
		}
	}
}
