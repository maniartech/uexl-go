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
		{"1 + 2", []any{1.0, 2.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpAdd),
		}},
		{"1 - 2", []any{1.0, 2.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpSub),
		}},
		{"1 * 2", []any{1.0, 2.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpMul),
		}},
		{"1 / 2", []any{1.0, 2.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpDiv),
		}},
		{"1 + 2 * 3", []any{1.0, 2.0, 3.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpConstant, 2),
			code.Make(code.OpMul),
			code.Make(code.OpAdd),
		}},
		{"-1", []any{1.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpMinus),
		}},
		{"5 % 2", []any{5.0, 2.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpMod),
		}},
		{"(1 + 2) * 3", []any{1.0, 2.0, 3.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpAdd),
			code.Make(code.OpConstant, 2),
			code.Make(code.OpMul),
		}},
		{"1 + 2 * 3 - 4 / 2", []any{1.0, 2.0, 3.0, 4.0, 2.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpConstant, 2),
			code.Make(code.OpMul),
			code.Make(code.OpAdd),
			code.Make(code.OpConstant, 3),
			code.Make(code.OpConstant, 4),
			code.Make(code.OpDiv),
			code.Make(code.OpSub),
		}},
		{"-(1 + 2) * (3 - 4)", []any{1.0, 2.0, 3.0, 4.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpAdd),
			code.Make(code.OpMinus),
			code.Make(code.OpConstant, 2),
			code.Make(code.OpConstant, 3),
			code.Make(code.OpSub),
			code.Make(code.OpMul),
		}},
	}
	runCompilerTestCases(t, cases)
}
func TestBooleanLiteralsAndOperations(t *testing.T) {
	cases := []compilerTestCase{
		{"true", []any{}, []code.Instructions{
			code.Make(code.OpTrue),
		}},
		{"false", []any{}, []code.Instructions{
			code.Make(code.OpFalse),
		}},
		{"!true", []any{}, []code.Instructions{
			code.Make(code.OpTrue),
			code.Make(code.OpBang),
		}},
		{"!false", []any{}, []code.Instructions{
			code.Make(code.OpFalse),
			code.Make(code.OpBang),
		}},
		{"true && false", []any{}, []code.Instructions{
			code.Make(code.OpTrue),
			code.Make(code.OpFalse),
			code.Make(code.OpLogicalAnd),
		}},
		{"true || false", []any{}, []code.Instructions{
			code.Make(code.OpTrue),
			code.Make(code.OpFalse),
			code.Make(code.OpLogicalOr),
		}},
		{"true == false", []any{}, []code.Instructions{
			code.Make(code.OpTrue),
			code.Make(code.OpFalse),
			code.Make(code.OpEqual),
		}},
		{"true != false", []any{}, []code.Instructions{
			code.Make(code.OpTrue),
			code.Make(code.OpFalse),
			code.Make(code.OpNotEqual),
		}},
		{"!true == false", []any{}, []code.Instructions{
			code.Make(code.OpTrue),
			code.Make(code.OpBang),
			code.Make(code.OpFalse),
			code.Make(code.OpEqual),
		}},
		{"!(true && false)", []any{}, []code.Instructions{
			code.Make(code.OpTrue),
			code.Make(code.OpFalse),
			code.Make(code.OpLogicalAnd),
			code.Make(code.OpBang),
		}},
		{"1 < 2", []any{2.0, 1.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpGreaterThan),
		}},
		{"2 > 1", []any{2.0, 1.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpGreaterThan),
		}},
		{"1 <= 2", []any{2.0, 1.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpGreaterThanOrEqual),
		}},
		{"2 >= 1", []any{2.0, 1.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpGreaterThanOrEqual),
		}},
	}
	runCompilerTestCases(t, cases)
}
func TestBitwiseOperators(t *testing.T) {
	cases := []compilerTestCase{
		{"1 & 2", []any{1.0, 2.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpBitwiseAnd),
		}},
		{"3 | 4", []any{3.0, 4.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpBitwiseOr),
		}},
		{"5 ^ 6", []any{5.0, 6.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpBitwiseXor),
		}},
		// {"~7", []any{7.0}, []code.Instructions{
		// 	code.Make(code.OpConstant, 0),
		// 	code.Make(code.OpBitwiseNot),
		// }},
		{"8 << 2", []any{8.0, 2.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpShiftLeft),
		}},
		{"16 >> 3", []any{16.0, 3.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpShiftRight),
		}},
		// Complex equations
		{"(1 & 2) | (3 ^ 4)", []any{1.0, 2.0, 3.0, 4.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpBitwiseAnd),
			code.Make(code.OpConstant, 2),
			code.Make(code.OpConstant, 3),
			code.Make(code.OpBitwiseXor),
			code.Make(code.OpBitwiseOr),
		}},
		// {"~(5 | 6) ^ (7 & 8)", []any{5.0, 6.0, 7.0, 8.0}, []code.Instructions{
		// 	code.Make(code.OpConstant, 0),
		// 	code.Make(code.OpConstant, 1),
		// 	code.Make(code.OpBitwiseOr),
		// 	code.Make(code.OpBitwiseNot),
		// 	code.Make(code.OpConstant, 2),
		// 	code.Make(code.OpConstant, 3),
		// 	code.Make(code.OpBitwiseAnd),
		// 	code.Make(code.OpBitwiseXor),
		// }},
		{"(2 << 3) & (32 >> 2)", []any{2.0, 3.0, 32.0, 2.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpShiftLeft),
			code.Make(code.OpConstant, 2),
			code.Make(code.OpConstant, 3),
			code.Make(code.OpShiftRight),
			code.Make(code.OpBitwiseAnd),
		}},
		{"((1 | 2) & 3) ^ (4 << 1)", []any{1.0, 2.0, 3.0, 4.0, 1.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpBitwiseOr),
			code.Make(code.OpConstant, 2),
			code.Make(code.OpBitwiseAnd),
			code.Make(code.OpConstant, 3),
			code.Make(code.OpConstant, 4),
			code.Make(code.OpShiftLeft),
			code.Make(code.OpBitwiseXor),
		}},
	}
	runCompilerTestCases(t, cases)
}
