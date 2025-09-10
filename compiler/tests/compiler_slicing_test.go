package compiler_test

import (
	"testing"

	"github.com/maniartech/uexl_go/code"
)

func TestCompileSliceExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "arr[1:5]",
			expectedConstants: []any{
				1.0,
				5.0,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpContextVar, 0),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNull),
				code.Make(code.OpSlice, 0),
			},
		},
		{
			input: "arr[:5]",
			expectedConstants: []any{
				5.0,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpContextVar, 0),
				code.Make(code.OpNull),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpNull),
				code.Make(code.OpSlice, 0),
			},
		},
		{
			input: "arr[1:]",
			expectedConstants: []any{
				1.0,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpContextVar, 0),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpNull),
				code.Make(code.OpNull),
				code.Make(code.OpSlice, 0),
			},
		},
		{
			input:             "arr[:]",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpContextVar, 0),
				code.Make(code.OpNull),
				code.Make(code.OpNull),
				code.Make(code.OpNull),
				code.Make(code.OpSlice, 0),
			},
		},
		{
			input: "arr[1:5:2]",
			expectedConstants: []any{
				1.0,
				5.0,
				2.0,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpContextVar, 0),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpSlice, 0),
			},
		},
		{
			input: "arr?[1:5]",
			expectedConstants: []any{
				1.0,
				5.0,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpContextVar, 0),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNull),
				code.Make(code.OpSlice, 1),
			},
		},
		{
			input:             "arr[a:b]",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpContextVar, 0), // arr
				code.Make(code.OpContextVar, 1), // a
				code.Make(code.OpContextVar, 2), // b
				code.Make(code.OpNull),
				code.Make(code.OpSlice, 0),
			},
		},
		{
			input: "arr[1+2:5*2]",
			expectedConstants: []any{
				1.0,
				2.0,
				5.0,
				2.0,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpContextVar, 0),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpMul),
				code.Make(code.OpNull),
				code.Make(code.OpSlice, 0),
			},
		},
		{
			input: "arr[1:5][0]",
			expectedConstants: []any{
				1.0,
				5.0,
				0.0,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpContextVar, 0),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNull),
				code.Make(code.OpSlice, 0),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpIndex, 0),
			},
		},
	}

	runCompilerTestCases(t, tests)
}
