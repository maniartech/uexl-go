package vm_test

import (
	"testing"
)

func TestSlicingExpressions(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `[1, 2, 3, 4, 5][1:4]`,
			expected: []any{2.0, 3.0, 4.0},
		},
		{
			input:    `[1, 2, 3, 4, 5][:3]`,
			expected: []any{1.0, 2.0, 3.0},
		},
		{
			input:    `[1, 2, 3, 4, 5][2:]`,
			expected: []any{3.0, 4.0, 5.0},
		},
		{
			input:    `[1, 2, 3, 4, 5][:]`,
			expected: []any{1.0, 2.0, 3.0, 4.0, 5.0},
		},
		{
			input:    `[1, 2, 3, 4, 5][1:5:2]`,
			expected: []any{2.0, 4.0},
		},
		{
			input:    `[1, 2, 3, 4, 5][::2]`,
			expected: []any{1.0, 3.0, 5.0},
		},
		{
			input:    `[1, 2, 3, 4, 5][::-1]`,
			expected: []any{5.0, 4.0, 3.0, 2.0, 1.0},
		},
		{
			input:    `[1, 2, 3, 4, 5][-1:]`,
			expected: []any{5.0},
		},
		{
			input:    `[1, 2, 3, 4, 5][:-1]`,
			expected: []any{1.0, 2.0, 3.0, 4.0},
		},
		{
			input:    `[1, 2, 3, 4, 5][-3:-1]`,
			expected: []any{3.0, 4.0},
		},
		{
			input:    `"hello"[1:4]`,
			expected: "ell",
		},
		{
			input:    `"hello"[::-1]`,
			expected: "olleh",
		},
		{
			input:    `"hello"[:4]`,
			expected: "hell",
		},
		{
			input:    `"hello"[1:]`,
			expected: "ello",
		},
		{
			input:    `[1, 2, 3, 4, 5][1:4][1:2]`,
			expected: []any{3.0},
		},
		{
			input:    `[1, 2, 3, 4, 5][1:4][0]`,
			expected: 2.0,
		},
		{
			input:    `null?[1:4]`,
			expected: nil,
		},
		{
			input:    `[1, 2, 3]?[1:2]`,
			expected: []any{2.0},
		},
		{
			input:    `[1, 2, 3, 4, 5][10:20]`,
			expected: []any{},
		},
		{
			input:    `"hello"[10:20]`,
			expected: "",
		},
		{
			input:    `[1, 2, 3, 4, 5][1:1]`,
			expected: []any{},
		},
		{
			input:    `[1, 2, 3, 4, 5][4:1]`,
			expected: []any{},
		},
		{
			input:    `[1, 2, 3, 4, 5][4:1:-1]`,
			expected: []any{5.0, 4.0, 3.0},
		},
	}

	runVmTests(t, tests)
}

func TestSlicingErrors(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `1[1:2]`,
			expected: "invalid type for slice: float64",
		},
		{
			input:    `[1, 2, 3]["a":2]`,
			expected: "slice index must be a number, got string",
		},
		{
			input:    `[1, 2, 3][1:"b"]`,
			expected: "slice index must be a number, got string",
		},
		{
			input:    `[1, 2, 3][1:2:"c"]`,
			expected: "slice step must be a number, got string",
		},
		{
			input:    `[1, 2, 3][1:2:0]`,
			expected: "slice step cannot be zero",
		},
		{
			input:    `[1, 2, 3][1.5:2]`,
			expected: "slice index must be an integer, got 1.5",
		},
		{
			input:    `null[1:2]`,
			expected: "cannot slice a null value",
		},
	}

	runVmErrorTests(t, tests)
}
