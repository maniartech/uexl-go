package vm_test

import (
	"testing"
)

func TestBuiltins_RuneLen(t *testing.T) {
	tests := []vmTestCase{
		{`runeLen("hello")`, 5.0},
		{`runeLen("naïve")`, 5.0}, // ï = 1 rune, not 2 bytes
		{`runeLen("café")`, 4.0},  // precomposed é
		{`runeLen("")`, 0.0},
		{`runeLen("naïve") == len("naïve")`, false}, // rune ≠ byte for non-ASCII
	}
	runVmTests(t, tests)
}

func TestBuiltins_RuneSubstr(t *testing.T) {
	tests := []vmTestCase{
		{`runeSubstr("naïve", 0, 3)`, "naï"},  // first 3 runes
		{`runeSubstr("naïve", 2, 1)`, "ï"},    // rune at index 2
		{`runeSubstr("hello", 1, 3)`, "ell"},  // ASCII: runes == bytes
		{`runeSubstr("hello", 3, 100)`, "lo"}, // clamp
		{`runeSubstr("hello", 10, 3)`, ""},    // start beyond end
		{`runeSubstr("", 0, 3)`, ""},
	}
	runVmTests(t, tests)
}

func TestBuiltins_GraphemeLen(t *testing.T) {
	tests := []vmTestCase{
		{`graphemeLen("hello")`, 5.0},
		{`graphemeLen("naïve")`, 5.0},
		{`graphemeLen("café\u0301")`, 4.0}, // decomposed é = 1 grapheme
		{`graphemeLen("")`, 0.0},
		// grapheme ≤ rune count (combining sequences reduce count)
		{`graphemeLen("café\u0301") == runeLen("café\u0301")`, false},
	}
	runVmTests(t, tests)
}

func TestBuiltins_GraphemeSubstr(t *testing.T) {
	tests := []vmTestCase{
		{`graphemeSubstr("hello", 0, 3)`, "hel"},
		{`graphemeSubstr("naïve", 0, 3)`, "naï"},             // ï = 1 grapheme
		{`graphemeSubstr("café\u0301", 0, 3)`, "caf"},        // 3 graphemes
		{`graphemeSubstr("café\u0301", 0, 4)`, "café\u0301"}, // é cluster intact
		{`graphemeSubstr("hello", 10, 3)`, ""},               // start beyond end
		{`graphemeSubstr("", 0, 3)`, ""},
	}
	runVmTests(t, tests)
}

func TestBuiltins_Runes(t *testing.T) {
	tests := []vmTestCase{
		{`runes("hi")`, []any{"h", "i"}},
		{`runes("naïve")`, []any{"n", "a", "ï", "v", "e"}},
		{`runes("")`, []any{}},
		{`len(runes("naïve"))`, 5.0},
	}
	runVmTests(t, tests)
}

func TestBuiltins_Graphemes(t *testing.T) {
	tests := []vmTestCase{
		{`graphemes("hi")`, []any{"h", "i"}},
		{`graphemes("naïve")`, []any{"n", "a", "ï", "v", "e"}},
		{`graphemes("")`, []any{}},
		{`len(graphemes("café\u0301"))`, 4.0},                              // 4 clusters
		{`graphemes("café\u0301")[3]`, "é\u0301"},                          // index into grapheme array (byte-based on []any element is fine here)
		{`graphemes("naïve") |map: $item`, []any{"n", "a", "ï", "v", "e"}}, // pipe integration
	}
	runVmTests(t, tests)
}

func TestBuiltins_Bytes(t *testing.T) {
	tests := []vmTestCase{
		{`bytes("hi")`, []any{float64('h'), float64('i')}},
		// ï = 0xC3 0xAF — two bytes
		{`len(bytes("ï"))`, 2.0},
		{`len(bytes("naïve"))`, 6.0}, // same as len("naïve")
		// byte count == len(s)
		{`len(bytes("naïve")) == len("naïve")`, true},
		{`bytes("") `, []any{}},
		// All ASCII-only bytes are < 128
		{`bytes("hi") |every: $item < 128`, true},
	}
	runVmTests(t, tests)
}

func TestBuiltins_Join(t *testing.T) {
	tests := []vmTestCase{
		{`join(["a", "b", "c"])`, "abc"},
		{`join(["a", "b", "c"], "-")`, "a-b-c"},
		{`join([], ",")`, ""},
		{`join(["x"])`, "x"},
		// Round-trip: explode then join.
		{`join(runes("naïve"), "")`, "naïve"},
		{`join(graphemes("naïve"), "")`, "naïve"},
		// Pipe then join.
		{`join(runes("hello") |map: $item, "")`, "hello"},
		// Join with separator.
		{`join(["foo", "bar", "baz"], ", ")`, "foo, bar, baz"},
	}
	runVmTests(t, tests)
}

func TestBuiltins_NewStringErrors(t *testing.T) {
	tests := []vmTestCase{
		// Wrong arg count.
		{`runeLen()`, "error calling function runeLen: runeLen expects 1 argument"},
		{`runeLen("a", "b")`, "error calling function runeLen: runeLen expects 1 argument"},
		{`runeSubstr("a", 0)`, "error calling function runeSubstr: runeSubstr expects 3 arguments"},
		{`graphemeLen()`, "error calling function graphemeLen: graphemeLen expects 1 argument"},
		{`graphemeSubstr("a", 0)`, "error calling function graphemeSubstr: graphemeSubstr expects 3 arguments"},
		{`runes()`, "error calling function runes: runes expects 1 argument"},
		{`graphemes()`, "error calling function graphemes: graphemes expects 1 argument"},
		{`bytes()`, "error calling function bytes: bytes expects 1 argument"},
		{`join()`, "error calling function join: join expects 1 or 2 arguments"},
		// Wrong types.
		{`join(["a", 1], "")`, "error calling function join: join: element 1 must be a string, got float64"},
	}
	runVmErrorTests(t, tests)
}
