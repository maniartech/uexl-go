package vm

import "testing"

// TestIsTruthyValue_Temporal covers datetime-spec §6.2: temporal values follow the zero-is-falsy rule —
// the epoch instant (ms=0) and a zero-length duration are falsy (like 0/""/false), every other value is
// truthy. Also re-checks the primitive cases since the hot/cold split was reorganized to keep
// isTruthyValue inlinable.
func TestIsTruthyValue_Temporal(t *testing.T) {
	cases := []struct {
		name string
		v    Value
		want bool
	}{
		{"datetime epoch (zero ms)", newDateTimeValue(0), false},
		{"datetime nonzero", newDateTimeValue(1733011200000), true},
		{"datetime negative", newDateTimeValue(-86400000), true},
		{"duration zero", newDurationValue(0), false},
		{"duration positive", newDurationValue(604800000), true},
		{"duration negative", newDurationValue(-604800000), true},
		{"null", newNullValue(), false},
		{"bool false", newBoolValue(false), false},
		{"bool true", newBoolValue(true), true},
		{"number zero", newFloatValue(0), false},
		{"number nonzero", newFloatValue(1), true},
		{"string empty", newStringValue(""), false},
		{"string nonempty", newStringValue("x"), true},
	}
	for _, c := range cases {
		if got := isTruthyValue(c.v); got != c.want {
			t.Errorf("%s: isTruthyValue=%v, want %v", c.name, got, c.want)
		}
	}
}

// TestIsTruthy_TemporalAny covers the any-based truthiness path (isTruthy) for boxed temporal wrappers,
// which must agree with the Value-based path (datetime-spec §6.2).
func TestIsTruthy_TemporalAny(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want bool
	}{
		{"DateTime epoch", DateTime{Millis: 0}, false},
		{"DateTime nonzero", DateTime{Millis: 1733011200000}, true},
		{"Duration zero", Duration{Millis: 0}, false},
		{"Duration nonzero", Duration{Millis: -604800000}, true},
	}
	for _, c := range cases {
		if got := isTruthy(c.v); got != c.want {
			t.Errorf("%s: isTruthy=%v, want %v", c.name, got, c.want)
		}
	}
}
