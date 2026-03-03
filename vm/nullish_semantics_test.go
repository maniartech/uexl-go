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

// TestOptionalChainingAbsentKey verifies that ?. softens both a null base AND a
// missing key/index on a non-null object — consistent with the "absent = null" principle.
func TestOptionalChainingAbsentKey(t *testing.T) {
	tests := []vmTestCase{
		// key absent in non-null object -> null (safe access)
		{`{}?.foo`, nil},
		{`{"a": 1}?.b`, nil},

		// chain: key absent at intermediate step -> null short-circuits the rest
		{`{}?.u?.name`, nil},
		{`{"user": {}}?.user?.address?.city`, nil},

		// combined with ?? -> fallback
		{`{}?.u?.name ?? "anon"`, "anon"},
		{`{"user": {"name": "Alice"}}?.user?.address?.city ?? "unknown"`, "unknown"},
		{`{"user": {}}?.user?.name ?? "Guest"`, "Guest"},

		// key present -> normal value, no fallback
		{`{"user": {"name": "Alice"}}?.user?.name ?? "Guest"`, "Alice"},
		{`{"user": {"address": {"city": "Paris"}}}?.user?.address?.city ?? "unknown"`, "Paris"},

		// optional index access on array: index absent -> null
		{`[1, 2, 3]?.[10] ?? 99`, 99.0},
		{`[1, 2, 3]?.[1] ?? 99`, 2.0},

		// safe step returns null -> next safe step short-circuits
		{`{"a": null}?.a?.b?.c ?? "end"`, "end"},
	}
	runVmTests(t, tests)

	// With an actual user context var present but lacking address
	userCtx := map[string]any{"user": map[string]any{"name": "Alice"}}
	ctxTests := []vmTestCase{
		{`user?.address?.city ?? "unknown"`, "unknown"},
		{`user?.name ?? "anon"`, "Alice"},
		{`user?.address?.zip ?? 0`, 0.0},
	}
	runVmTests(t, ctxTests, userCtx)
}

// TestAbsentContextVarIsNullish verifies that a context variable not provided by the
// caller is treated as null (nullish), allowing ?., ?? and ?? chains to guard it.
func TestAbsentContextVarIsNullish(t *testing.T) {
	// No context values provided — all variables are absent.
	tests := []vmTestCase{
		// Direct ?? fallback on absent variable
		{`x ?? "default"`, "default"},
		{`count ?? 0`, 0.0},

		// Optional chaining on absent variable
		{`user?.name ?? "anon"`, "anon"},
		{`user?.address?.city ?? "unknown"`, "unknown"},

		// Optional index access on absent variable
		{`items?.[0] ?? "none"`, "none"},

		// Deep ?? chain: multiple absent vars -> last fallback wins
		{`a ?? b ?? "last"`, "last"},
		{`a ?? b ?? c ?? 42`, 42.0},
	}
	runVmTests(t, tests)

	// Falsy values in context must NOT trigger ?? fallback
	falsyCtx := map[string]any{"n": 0.0, "s": "", "b": false}
	falsyTests := []vmTestCase{
		{`n ?? 99`, 0.0},
		{`s ?? "fallback"`, ""},
		{`b ?? true`, false},
	}
	runVmTests(t, falsyTests, falsyCtx)
}

// TestNullLiteralIsNullish verifies that an explicit null literal is nullish, and that
// falsy inline literals (0, "", false, [], {}) are NOT nullish.
func TestNullLiteralIsNullish(t *testing.T) {
	tests := []vmTestCase{
		// null literal triggers ?? fallback
		{`null ?? "default"`, "default"},
		{`null ?? 0`, 0.0},
		{`null ?? false`, false},

		// null literal with optional chaining
		{`null?.foo ?? "x"`, "x"},
		{`null?.[0] ?? "x"`, "x"},

		// falsy inline literals are NOT nullish
		{`0 ?? 99`, 0.0},
		{`"" ?? "fallback"`, ""},
		{`false ?? true`, false},
		{`[] ?? "x"`, []any{}},
		{`{} ?? "x"`, map[string]any{}},
	}
	runVmTests(t, tests)
}

// TestNullContextVarIsNullish verifies that a context var explicitly set to null
// behaves identically to an absent context var — both are nullish.
func TestNullContextVarIsNullish(t *testing.T) {
	nullCtx := map[string]any{"a": nil, "user": nil}
	tests := []vmTestCase{
		// Explicit null in context triggers ?? fallback
		{`a ?? "default"`, "default"},
		{`a ?? 0`, 0.0},

		// Optional chaining on null-valued context var
		{`user?.name ?? "anon"`, "anon"},
		{`user?.address?.city ?? "unknown"`, "unknown"},

		// null and absent are interchangeable for ??
		{`a ?? b ?? "last"`, "last"}, // a=null, b absent -> "last"
	}
	runVmTests(t, tests, nullCtx)
}
