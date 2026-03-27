package vm_test

import (
	"testing"
)

// TestPipeParams_Window exercises |window(n): and the backward-compatible |window: form.
func TestPipeParams_Window(t *testing.T) {
	tests := []vmTestCase{
		// explicit size — result count
		{"len([1,2,3,4,5] |window(3): $window)", 3.0},
		{"len([1,2,3,4,5] |window(4): $window)", 2.0},
		{"len([1,2,3,4,5] |window(5): $window)", 1.0},
		// size > len → 0 windows
		{"len([1,2,3,4,5] |window(6): $window)", 0.0},

		// backward compat: no args → default size 2
		{"len([1,2,3,4] |window: $window)", 3.0},

		// nested array shape
		{"[1,2,3,4,5] |window(3): $window", []any{
			[]any{1.0, 2.0, 3.0},
			[]any{2.0, 3.0, 4.0},
			[]any{3.0, 4.0, 5.0},
		}},
		{"[1,2,3,4] |window: $window", []any{
			[]any{1.0, 2.0},
			[]any{2.0, 3.0},
			[]any{3.0, 4.0},
		}},

		// predicate over window elements
		{"[1,2,3,4,5] |window(3): $window[0] + $window[1] + $window[2]", []any{6.0, 9.0, 12.0}},

		// empty input → empty result
		{"[] |window(3): $window", []any{}},

		// whitespace-insensitive args syntax
		{"len([1,2,3,4,5] |window( 3 ): $window)", 3.0},
	}
	runVmTests(t, tests)
}

// TestPipeParams_Chunk exercises |chunk(n): and the backward-compatible |chunk: form.
func TestPipeParams_Chunk(t *testing.T) {
	tests := []vmTestCase{
		// explicit size — result count
		{"len([1,2,3,4,5] |chunk(4): $chunk)", 2.0}, // [1,2,3,4] + [5]
		{"len([1,2,3,4,5] |chunk(3): $chunk)", 2.0}, // [1,2,3] + [4,5]
		{"len([1,2,3,4,5,6] |chunk(3): $chunk)", 2.0},
		// chunk size equals array length → single chunk
		{"len([1,2,3] |chunk(3): $chunk)", 1.0},

		// backward compat: no args → default size 2
		{"len([1,2,3,4] |chunk: $chunk)", 2.0},

		// nested array shape
		{"[1,2,3,4,5] |chunk(4): $chunk", []any{
			[]any{1.0, 2.0, 3.0, 4.0},
			[]any{5.0},
		}},
		{"[1,2,3,4] |chunk: $chunk", []any{
			[]any{1.0, 2.0},
			[]any{3.0, 4.0},
		}},

		// empty input → empty result
		{"[] |chunk(3): $chunk", []any{}},

		// whitespace-insensitive args syntax
		{"len([1,2,3,4,5] |chunk( 4 ): $chunk)", 2.0},
	}
	runVmTests(t, tests)
}

// TestPipeParams_Chain verifies that pipes with args compose correctly with other pipes.
func TestPipeParams_Chain(t *testing.T) {
	tests := []vmTestCase{
		// window(3) then map to extract first element of each window
		{"[1,2,3,4,5] |window(3): $window |map: $item[0]", []any{1.0, 2.0, 3.0}},

		// chunk(3) then keep only full chunks
		{"len([1,2,3,4,5,6] |chunk(3): $chunk |filter: len($chunk) == 3)", 2.0},

		// window(2) default then map to sum of pair
		{"[1,2,3,4] |window: $window |map: $item[0] + $item[1]", []any{3.0, 5.0, 7.0}},
	}
	runVmTests(t, tests)
}
