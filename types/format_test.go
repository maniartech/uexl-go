package types

import "testing"

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
