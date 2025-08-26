package vm_test

import (
	"testing"
)

// These tests exercise nullish coalescing around immediate property/index access
// The current semantics should soft-fail only the last access on the left of ??.
func TestNullishCoalescingSemantics(t *testing.T) {
	tests := []vmTestCase{
		// user exists? If not, earlier step should throw (but our harness expects a value),
		// therefore we structure tests where base exists to observe immediate softening.

		// Immediate property missing -> fallback
		{`{"user": {}}.user.name ?? "Anonymous"`, "Anonymous"},
		// Immediate property null -> fallback
		{`{"user": {"name": null}}.user.name ?? "Anonymous"`, "Anonymous"},
		// Immediate property present -> no fallback
		{`{"user": {"name": "Alice"}}.user.name ?? "Anonymous"`, "Alice"},

		// Indexing cases
		{`{"arr": [1,2,3]}.arr[5] ?? 99`, 99.0},
		{`{"arr": [1,2,3]}.arr[1] ?? 99`, 2.0},

		// Optional chaining combined with ?? should still work
		{`{"u": null}?.u?.name ?? "anon"`, "anon"},
		{`{"u": {"name": null}}?.u?.name ?? "anon"`, "anon"},
		{`{"u": {"name": "Bob"}}?["u"]?.name ?? "anon"`, "Bob"},
	}

	runVmTests(t, tests)
}
