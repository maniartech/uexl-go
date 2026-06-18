package types

import "testing"

// Portable-range boundaries (datetime-spec §5.2).
const (
	minMillis = int64(-62135596800000) // 0001-01-01T00:00:00.000Z
	maxMillis = int64(253402300799999) // 9999-12-31T23:59:59.999Z
)

func TestDateTimeValue_RoundTrip(t *testing.T) {
	cases := []int64{0, 1, -1, minMillis, maxMillis, 1733011200000, -86400000}
	for _, ms := range cases {
		v := NewDateTimeValue(ms)
		if v.Typ != TypeDateTime {
			t.Fatalf("ms=%d: Typ=%d, want TypeDateTime", ms, v.Typ)
		}
		if !v.IsDateTime() || v.IsDuration() {
			t.Errorf("ms=%d: IsDateTime=%v IsDuration=%v, want true/false", ms, v.IsDateTime(), v.IsDuration())
		}
		if got, ok := v.AsDateTimeMillis(); !ok || got != ms {
			t.Errorf("ms=%d: AsDateTimeMillis=(%d,%v), want (%d,true)", ms, got, ok, ms)
		}
	}
}

func TestDurationValue_RoundTrip(t *testing.T) {
	cases := []int64{0, 1, -1, 604800000, -604800000, maxMillis - minMillis}
	for _, ms := range cases {
		v := NewDurationValue(ms)
		if !v.IsDuration() || v.IsDateTime() {
			t.Errorf("ms=%d: IsDuration=%v IsDateTime=%v, want true/false", ms, v.IsDuration(), v.IsDateTime())
		}
		if got, ok := v.AsDurationMillis(); !ok || got != ms {
			t.Errorf("ms=%d: AsDurationMillis=(%d,%v), want (%d,true)", ms, got, ok, ms)
		}
	}
}

// TestTemporal_ExactFloat64RoundTrip is the load-bearing invariant behind storing ms in FloatVal:
// every value in (and a bit beyond) the portable range is < 2^53, so int64<->float64 is lossless.
func TestTemporal_ExactFloat64RoundTrip(t *testing.T) {
	for _, ms := range []int64{minMillis, maxMillis, 0, 1 << 52, -(1 << 52), maxMillis - minMillis} {
		if int64(float64(ms)) != ms {
			t.Errorf("float64 round-trip lost precision for %d", ms)
		}
	}
}

func TestTemporal_BoundaryConversions(t *testing.T) {
	if got := NewDateTimeValue(maxMillis).ToAny(); got != (DateTime{Millis: maxMillis}) {
		t.Errorf("datetime ToAny=%v, want DateTime{%d}", got, maxMillis)
	}
	if got := NewDurationValue(-1).ToAny(); got != (Duration{Millis: -1}) {
		t.Errorf("duration ToAny=%v, want Duration{-1}", got)
	}
	// Host boundary round-trip: external wrapper -> Value -> external wrapper.
	dt := NewAnyValue(DateTime{Millis: 1733011200000})
	if !dt.IsDateTime() {
		t.Fatal("NewAnyValue(DateTime) is not a datetime")
	}
	if ms, _ := dt.AsDateTimeMillis(); ms != 1733011200000 {
		t.Errorf("round-trip datetime ms=%d, want 1733011200000", ms)
	}
	du := NewAnyValue(Duration{Millis: 5400000})
	if !du.IsDuration() {
		t.Fatal("NewAnyValue(Duration) is not a duration")
	}
	if ms, _ := du.AsDurationMillis(); ms != 5400000 {
		t.Errorf("round-trip duration ms=%d, want 5400000", ms)
	}
}

// TestTemporal_TypeSafety verifies temporal values are NOT interchangeable with number or each other,
// even though they share the FloatVal storage slot.
func TestTemporal_TypeSafety(t *testing.T) {
	dt := NewDateTimeValue(1000)
	if _, ok := dt.AsFloat(); ok {
		t.Error("datetime AsFloat ok=true, want false (not a number)")
	}
	if _, ok := dt.AsDurationMillis(); ok {
		t.Error("datetime AsDurationMillis ok=true, want false")
	}
	if _, ok := NewFloatValue(1000).AsDateTimeMillis(); ok {
		t.Error("number AsDateTimeMillis ok=true, want false")
	}
	if _, ok := NewDurationValue(1000).AsDateTimeMillis(); ok {
		t.Error("duration AsDateTimeMillis ok=true, want false")
	}
}

func TestTemporal_ZeroAllocations(t *testing.T) {
	if allocs := testing.AllocsPerRun(1000, func() {
		v := NewDateTimeValue(1733011200000)
		d := NewDurationValue(604800000)
		_, _ = v.AsDateTimeMillis()
		_, _ = d.AsDurationMillis()
	}); allocs != 0 {
		t.Errorf("temporal construct/extract allocated %v times, want 0", allocs)
	}
}

func BenchmarkNewDateTimeValue(b *testing.B) {
	var sink Value
	for i := 0; i < b.N; i++ {
		sink = NewDateTimeValue(int64(i))
	}
	_ = sink
}
