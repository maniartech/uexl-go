package uexl_test

import (
	"testing"

	"github.com/maniartech/uexl"
	"github.com/maniartech/uexl/types"
)

// End-to-end (tokenize -> parse -> compile -> evaluate) checks that datetime/duration literals produce
// the correct typed constant Value at the public API boundary.

func TestEval_DateTimeLiteral(t *testing.T) {
	cases := map[string]int64{
		`d"2024-12-01"`:           1733011200000,
		`d"1970-01-01T00:00:00Z"`: 0,
		`d"1969-12-31T00:00:00Z"`: -86400000,
	}
	for expr, want := range cases {
		got, err := uexl.Eval(expr, nil)
		if err != nil {
			t.Errorf("%s: unexpected error %v", expr, err)
			continue
		}
		dt, ok := got.(types.DateTime)
		if !ok {
			t.Errorf("%s: expected types.DateTime, got %T (%v)", expr, got, got)
			continue
		}
		if dt.Millis != want {
			t.Errorf("%s: got %d, want %d", expr, dt.Millis, want)
		}
	}
}

func TestEval_DurationLiteral(t *testing.T) {
	cases := map[string]int64{
		`7d`:    604800000,
		`30ms`:  30,
		`1.5h`:  5400000,
		`500ms`: 500,
	}
	for expr, want := range cases {
		got, err := uexl.Eval(expr, nil)
		if err != nil {
			t.Errorf("%s: unexpected error %v", expr, err)
			continue
		}
		d, ok := got.(types.Duration)
		if !ok {
			t.Errorf("%s: expected types.Duration, got %T (%v)", expr, got, got)
			continue
		}
		if d.Millis != want {
			t.Errorf("%s: got %d, want %d", expr, d.Millis, want)
		}
	}
}

func TestEval_DateTimeLiteralErrors(t *testing.T) {
	for _, expr := range []string{`d"2024-13-01"`, `d"2024-02-30"`, `d"0000-01-01"`, `7d"2024-12-01"`} {
		if _, err := uexl.Eval(expr, nil); err == nil {
			t.Errorf("%s: expected a parse-time error", expr)
		}
	}
}

// TestEval_TemporalInArray confirms temporal literals also work inside a larger expression (an array),
// each materializing to the correct typed value at the boundary.
func TestEval_TemporalInArray(t *testing.T) {
	got, err := uexl.Eval(`[d"2024-12-01", 7d]`, nil)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	arr, ok := got.([]any)
	if !ok || len(arr) != 2 {
		t.Fatalf("expected a 2-element array, got %T %v", got, got)
	}
	if dt, ok := arr[0].(types.DateTime); !ok || dt.Millis != 1733011200000 {
		t.Errorf("arr[0]: expected DateTime{1733011200000}, got %T %v", arr[0], arr[0])
	}
	if dur, ok := arr[1].(types.Duration); !ok || dur.Millis != 604800000 {
		t.Errorf("arr[1]: expected Duration{604800000}, got %T %v", arr[1], arr[1])
	}
}
