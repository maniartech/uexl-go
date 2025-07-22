package compiler_test

import (
	"testing"

	"github.com/maniartech/uexl_go/code"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []any
	expectedInstructions []code.Instructions
}

func TestNumberArithmetic(t *testing.T) {
	cases := []compilerTestCase{
		{"1 + 2", []any{"1", "2"}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpAdd),
		}},
		{"1 - 2", []any{"1", "2"}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpSub),
		}},
		{"1 * 2", []any{"1", "2"}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpMul),
		}},
		{"1 / 2", []any{"1", "2"}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpDiv),
		}},
	}
	runCompilerTestCases(t, cases)
}
