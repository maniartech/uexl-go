package uexl_test

import (
	"testing"

	"github.com/maniartech/uexl"
	"github.com/maniartech/uexl/types"
)

// End-to-end (tokenize -> parse -> compile -> run) checks for Phase C temporal operators.

func mustEval(t *testing.T, expr string) any {
	t.Helper()
	v, err := uexl.Eval(expr, nil)
	if err != nil {
		t.Fatalf("%s: unexpected error: %v", expr, err)
	}
	return v
}

func evalDate(t *testing.T, expr string) int64 {
	t.Helper()
	v := mustEval(t, expr)
	d, ok := v.(types.DateTime)
	if !ok {
		t.Fatalf("%s: expected DateTime, got %T (%v)", expr, v, v)
	}
	return d.Millis
}

func evalDur(t *testing.T, expr string) int64 {
	t.Helper()
	v := mustEval(t, expr)
	d, ok := v.(types.Duration)
	if !ok {
		t.Fatalf("%s: expected Duration, got %T (%v)", expr, v, v)
	}
	return d.Millis
}

func evalBool(t *testing.T, expr string) bool {
	t.Helper()
	v := mustEval(t, expr)
	b, ok := v.(bool)
	if !ok {
		t.Fatalf("%s: expected bool, got %T (%v)", expr, v, v)
	}
	return b
}

func evalErr(t *testing.T, expr string) {
	t.Helper()
	if v, err := uexl.Eval(expr, nil); err == nil {
		t.Errorf("%s: expected an error, got %v", expr, v)
	}
}

func TestEval_DateArithmetic(t *testing.T) {
	if got := evalDur(t, `d"2024-12-02" - d"2024-12-01"`); got != 86400000 { // dt-op-001
		t.Errorf("date - date = %d, want 86400000", got)
	}
	if got := evalDate(t, `d"2024-12-01" + 1d`); got != 1733097600000 {
		t.Errorf("date + dur = %d", got)
	}
	if got := evalDate(t, `1d + d"2024-12-01"`); got != 1733097600000 { // commutes
		t.Errorf("dur + date = %d", got)
	}
	if got := evalDate(t, `d"2024-12-01" - 1d`); got != 1732924800000 {
		t.Errorf("date - dur = %d", got)
	}
}

func TestEval_DurationArithmetic(t *testing.T) {
	if got := evalDur(t, `7d + 1d`); got != 691200000 {
		t.Errorf("dur + dur = %d", got)
	}
	if got := evalDur(t, `7d - 1d`); got != 518400000 {
		t.Errorf("dur - dur = %d", got)
	}
	if got := evalDur(t, `2d * 3`); got != 518400000 {
		t.Errorf("dur * num = %d", got)
	}
	if got := evalDur(t, `3 * 2d`); got != 518400000 { // commutes
		t.Errorf("num * dur = %d", got)
	}
	if got := evalDur(t, `1h / 2`); got != 1800000 { // dt-div-002 shape
		t.Errorf("dur / num = %d", got)
	}
	if got := mustEval(t, `1h / 15m`).(float64); got != 4 { // dt-div-001: ratio
		t.Errorf("dur / dur = %v, want 4", got)
	}
	if got := evalDur(t, `-7d`); got != -604800000 { // unary minus negates a duration
		t.Errorf("-dur = %d", got)
	}
}

func TestEval_TemporalComparison(t *testing.T) {
	truthy := []string{
		`d"2024-12-01" < d"2024-12-02"`, // dt-op-004
		`d"2024-12-02" > d"2024-12-01"`,
		`d"2024-12-01" <= d"2024-12-01"`,
		`d"2024-12-01" >= d"2024-12-01"`,
		`d"2024-12-01" == d"2024-12-01"`,
		`d"2024-12-01" != d"2024-12-02"`,
		`7d > 1d`,
		`1h == 60m`,
		`d"1970-01-01T00:00:00Z" != 0`, // datetime is never equal to a number
	}
	for _, expr := range truthy {
		if !evalBool(t, expr) {
			t.Errorf("%s: expected true", expr)
		}
	}
	falsy := []string{
		`d"1970-01-01T00:00:00Z" == 0`, // distinct types, same ms -> not equal
		`d"2024-12-01" == 7d`,          // datetime vs duration -> not equal
		`d"2024-12-01" > d"2024-12-02"`,
	}
	for _, expr := range falsy {
		if evalBool(t, expr) {
			t.Errorf("%s: expected false", expr)
		}
	}
}

func TestEval_TemporalRejections(t *testing.T) {
	for _, expr := range []string{
		`d"2024-12-01" + 7`,             // dt-err-007: datetime + bare number
		`d"2024-12-01" - 7`,             // datetime - bare number
		`d"2024-12-01" + d"2024-12-02"`, // dt-err-008: datetime + datetime
		`7d + 7`,                        // duration + bare number
		`7d - 7`,                        // duration - bare number
		`7d * 1d`,                       // duration * duration
		`d"2024-12-01" * 2`,             // datetime * number
		`d"2024-12-01" % 2`,             // datetime % number
		`d"2024-12-01" * d"2024-12-02"`, // datetime * datetime
		`d"2024-12-01" * 7d`,            // datetime * duration
		`7d - d"2024-12-01"`,            // duration - datetime (only date - dur / date - date subtract)
		`2 / 7d`,                        // number / duration (only dur / num and dur / dur divide)
		`7d / 0`,                        // duration / zero number
		`7d / 0d`,                       // duration / zero duration
		`7d * NaN`,                      // duration scaled by a non-numeric value
		`7d / NaN`,                      // duration divided by a non-numeric value
		`NaN * 7d`,                      // non-numeric scale (commuted)
		`-d"2024-12-01"`,                // unary minus on a datetime
		`d"2024-12-01" < 7`,             // ordering datetime vs number
	} {
		evalErr(t, expr)
	}
}

func TestEval_TemporalClamp(t *testing.T) {
	// Adding past the upper bound clamps to MAX (datetime-spec §8.4), not int64-wrap.
	if got := evalDate(t, `d"9999-12-31T23:59:59.999Z" + 1d`); got != types.MaxDateTimeMillis {
		t.Errorf("upper clamp = %d, want %d", got, types.MaxDateTimeMillis)
	}
	// Subtracting past the lower bound clamps to MIN.
	if got := evalDate(t, `d"0001-01-01T00:00:00Z" - 1d`); got != types.MinDateTimeMillis {
		t.Errorf("lower clamp = %d, want %d", got, types.MinDateTimeMillis)
	}
	// Duration results are bounded to ±MaxDurationMillis. 3652058d is just under the bound; scaling it
	// past the bound clamps (keeping the int64<->float64 storage exact).
	if got := evalDur(t, `3652058d * 2`); got != types.MaxDurationMillis {
		t.Errorf("duration upper clamp = %d, want %d", got, types.MaxDurationMillis)
	}
	if got := evalDur(t, `3652058d * -2`); got != -types.MaxDurationMillis {
		t.Errorf("duration lower clamp = %d, want %d", got, -types.MaxDurationMillis)
	}
	// Extreme multipliers overflow the float64 product to ±Inf; the result must saturate to the bound
	// with the correct sign, not wrap via an implementation-defined int64(±Inf).
	if got := evalDur(t, `7d * 1e300`); got != types.MaxDurationMillis {
		t.Errorf("overflow *+big = %d, want %d", got, types.MaxDurationMillis)
	}
	if got := evalDur(t, `7d * -1e300`); got != -types.MaxDurationMillis {
		t.Errorf("overflow *-big = %d, want %d", got, -types.MaxDurationMillis)
	}
	// duration + duration past the bound clamps (int64 path).
	if got := evalDur(t, `3652058d + 3652058d`); got != types.MaxDurationMillis {
		t.Errorf("dur+dur upper clamp = %d, want %d", got, types.MaxDurationMillis)
	}
	if got := evalDur(t, `0d - 3652058d - 3652058d`); got != -types.MaxDurationMillis {
		t.Errorf("dur-dur lower clamp = %d, want %d", got, -types.MaxDurationMillis)
	}
}
