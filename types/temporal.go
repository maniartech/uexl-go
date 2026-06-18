package types

// DateTime and Duration are the boundary representations of UExL's temporal values — the form a
// datetime/duration Value materializes to when it leaves the evaluator (Value.ToAny) and the form a
// host may hand in (NewAnyValue). Inside the evaluator the value lives inline in Value.FloatVal with
// Typ TypeDateTime/TypeDuration and never allocates; these wrappers exist only so host code can
// distinguish a temporal value from a plain number at the integration boundary, and so the
// int64<->float64 conversion happens in exactly one place. See docs/specs/datetime-spec.md §2, §9.
//
// Millis is the canonical signed millisecond count:
//   - DateTime.Millis: milliseconds since the Unix epoch 1970-01-01T00:00:00.000Z (UTC).
//   - Duration.Millis: an exact elapsed span in milliseconds (may be negative).

// DateTime is the boundary form of a datetime instant.
type DateTime struct {
	Millis int64
}

// Duration is the boundary form of an exact duration span.
type Duration struct {
	Millis int64
}
