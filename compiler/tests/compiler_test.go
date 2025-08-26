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
		{"1 ** 2", []any{1.0, 2.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpPow),
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

func TestLogicalShortCircuitCompilation(t *testing.T) {
	cases := []compilerTestCase{
		// OR: first truthy wins; all falsy -> last term
		{
			`false || 0 || ""`,
			[]any{0.0, ""}, // constants: 0, ""
			[]code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpJumpIfTruthy, 13),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpJumpIfTruthy, 13),
				code.Make(code.OpConstant, 1),
				// No normalization OpFalse at the end
			},
		},
		{
			`true || 0 || ""`,
			[]any{0.0, ""}, // constants: 0, ""
			[]code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpJumpIfTruthy, 13),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpJumpIfTruthy, 13),
				code.Make(code.OpConstant, 1),
			},
		},

		// AND: first falsy wins; all truthy -> last term
		{
			`false && 0 && ""`,
			[]any{0.0, ""}, // constants: 0, ""
			[]code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpJumpIfFalsy, 13),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpJumpIfFalsy, 13),
				code.Make(code.OpConstant, 1),
				// No jump-to-false block
			},
		},
		{
			`true && 1 && "hello"`,
			[]any{1.0, "hello"},
			[]code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpJumpIfFalsy, 13),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpJumpIfFalsy, 13),
				code.Make(code.OpConstant, 1),
			},
		},
		{
			`false || true`,
			[]any{},
			[]code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpJumpIfTruthy, 5),
				code.Make(code.OpTrue),
			},
		},
		{
			`true && false`,
			[]any{},
			[]code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpJumpIfFalsy, 5),
				code.Make(code.OpFalse),
			},
		},
	}
	runCompilerTestCases(t, cases)
}

func TestContextVariables(t *testing.T) {
	cases := []compilerTestCase{
		{"foo", []any{}, []code.Instructions{
			code.Make(code.OpContextVar, 0),
		}},
		{"foo + bar", []any{}, []code.Instructions{
			code.Make(code.OpContextVar, 0),
			code.Make(code.OpContextVar, 1),
			code.Make(code.OpAdd),
		}},
		{"foo - bar * baz", []any{}, []code.Instructions{
			code.Make(code.OpContextVar, 0),
			code.Make(code.OpContextVar, 1),
			code.Make(code.OpContextVar, 2),
			code.Make(code.OpMul),
			code.Make(code.OpSub),
		}},
		{"-foo", []any{}, []code.Instructions{
			code.Make(code.OpContextVar, 0),
			code.Make(code.OpMinus),
		}},
		{"foo && bar", []any{}, []code.Instructions{
			code.Make(code.OpContextVar, 0),
			code.Make(code.OpContextVar, 1),
			code.Make(code.OpLogicalAnd),
		}},
		{"foo == bar", []any{}, []code.Instructions{
			code.Make(code.OpContextVar, 0),
			code.Make(code.OpContextVar, 1),
			code.Make(code.OpEqual),
		}},
		{"foo != bar", []any{}, []code.Instructions{
			code.Make(code.OpContextVar, 0),
			code.Make(code.OpContextVar, 1),
			code.Make(code.OpNotEqual),
		}},
		{"foo < bar", []any{}, []code.Instructions{
			code.Make(code.OpContextVar, 0),
			code.Make(code.OpContextVar, 1),
			code.Make(code.OpGreaterThan),
		}},
		{"foo <= bar", []any{}, []code.Instructions{
			code.Make(code.OpContextVar, 0),
			code.Make(code.OpContextVar, 1),
			code.Make(code.OpGreaterThanOrEqual),
		}},
	}
	runCompilerTestCases(t, cases)
}
func TestArrayLiterals(t *testing.T) {
	cases := []compilerTestCase{
		{"[]", []any{}, []code.Instructions{
			code.Make(code.OpArray, 0),
		}},
		{"[1]", []any{1.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpArray, 1),
		}},
		{"[1, 2]", []any{1.0, 2.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpArray, 2),
		}},
		{"[1 + 2, 3]", []any{1.0, 2.0, 3.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpAdd),
			code.Make(code.OpConstant, 2),
			code.Make(code.OpArray, 2),
		}},
		{"[foo, bar]", []any{}, []code.Instructions{
			code.Make(code.OpContextVar, 0),
			code.Make(code.OpContextVar, 1),
			code.Make(code.OpArray, 2),
		}},
		{"[1, [2, 3]]", []any{1.0, 2.0, 3.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpConstant, 2),
			code.Make(code.OpArray, 2),
			code.Make(code.OpArray, 2),
		}},
	}
	runCompilerTestCases(t, cases)
}

func TestArrayBracketIndexAccess(t *testing.T) {
	cases := []compilerTestCase{
		{"[1, 2, 3][0]", []any{1.0, 2.0, 3.0, 0.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpConstant, 2),
			code.Make(code.OpArray, 3),
			code.Make(code.OpConstant, 3),
			code.Make(code.OpIndex),
		}},
		{"arr[1]", []any{1.0}, []code.Instructions{
			code.Make(code.OpContextVar, 0),
			code.Make(code.OpConstant, 0),
			code.Make(code.OpIndex),
		}},
		{"[1, 2, 3][2]", []any{1.0, 2.0, 3.0, 2.0}, []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpConstant, 2),
			code.Make(code.OpArray, 3),
			code.Make(code.OpConstant, 3),
			code.Make(code.OpIndex),
		}},
		{"arr[0] + arr[1]", []any{0.0, 1.0}, []code.Instructions{
			code.Make(code.OpContextVar, 0),
			code.Make(code.OpConstant, 0),
			code.Make(code.OpIndex),
			code.Make(code.OpContextVar, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpIndex),
			code.Make(code.OpAdd),
		}},
		{"[foo, bar][0]", []any{0.0}, []code.Instructions{
			code.Make(code.OpContextVar, 0),
			code.Make(code.OpContextVar, 1),
			code.Make(code.OpArray, 2),
			code.Make(code.OpConstant, 0),
			code.Make(code.OpIndex),
		}},
	}
	runCompilerTestCases(t, cases)
}
func TestObjectLiterals(t *testing.T) {
	cases := []compilerTestCase{
		// Empty object
		{"{}", []any{}, []code.Instructions{
			code.Make(code.OpObject, 0),
		}},
		// Single property with number value
		{`{"a": 1}`, []any{"a", 1.0}, []code.Instructions{
			code.Make(code.OpConstant, 0), // "a"
			code.Make(code.OpConstant, 1), // 1.0
			code.Make(code.OpObject, 2),
		}},
		// Multiple properties with number values
		{`{"a": 1, "b": 2}`, []any{"a", 1.0, "b", 2.0}, []code.Instructions{
			code.Make(code.OpConstant, 0), // "a"
			code.Make(code.OpConstant, 1), // 1.0
			code.Make(code.OpConstant, 2), // "b"
			code.Make(code.OpConstant, 3), // 2.0
			code.Make(code.OpObject, 4),
		}},
		// Property with boolean value
		{`{"flag": true}`, []any{"flag"}, []code.Instructions{
			code.Make(code.OpConstant, 0), // "flag"
			code.Make(code.OpTrue),
			code.Make(code.OpObject, 2),
		}},
		// Property with string value
		{`{"name": "bob"}`, []any{"name", "bob"}, []code.Instructions{
			code.Make(code.OpConstant, 0), // "name"
			code.Make(code.OpConstant, 1), // "bob"
			code.Make(code.OpObject, 2),
		}},
		// Nested object
		{`{"outer": {"inner": 42}}`, []any{"outer", "inner", 42.0}, []code.Instructions{
			code.Make(code.OpConstant, 0), // "outer"
			code.Make(code.OpConstant, 1), // "inner"
			code.Make(code.OpConstant, 2), // 42.0
			code.Make(code.OpObject, 2),
			code.Make(code.OpObject, 2),
		}},
		// Object with array property
		{`{"arr": [1, 2]}`, []any{"arr", 1.0, 2.0}, []code.Instructions{
			code.Make(code.OpConstant, 0), // "arr"
			code.Make(code.OpConstant, 1), // 1.0
			code.Make(code.OpConstant, 2), // 2.0
			code.Make(code.OpArray, 2),
			code.Make(code.OpObject, 2),
		}},
		// Object with context variable property
		{`{"foo": bar}`, []any{"foo"}, []code.Instructions{
			code.Make(code.OpConstant, 0),   // "foo"
			code.Make(code.OpContextVar, 0), // bar
			code.Make(code.OpObject, 2),
		}},
		// Multiple property types
		{`{"x": 1, "y": true, "z": "hi"}`, []any{"x", 1.0, "y", "z", "hi"}, []code.Instructions{
			code.Make(code.OpConstant, 0), // "x"
			code.Make(code.OpConstant, 1), // 1.0
			code.Make(code.OpConstant, 2), // "y"
			code.Make(code.OpTrue),
			code.Make(code.OpConstant, 3), // "z"
			code.Make(code.OpConstant, 4), // "hi"
			code.Make(code.OpObject, 6),
		}},
		{ // Object with mixed types and nested structure
			`{"a": 1, "b": {"c": 2}}`, []any{"a", 1.0, "b", "c", 2.0}, []code.Instructions{
				code.Make(code.OpConstant, 0), // "a"
				code.Make(code.OpConstant, 1), // 1.0
				code.Make(code.OpConstant, 2), // "b"
				code.Make(code.OpConstant, 3), // "c"
				code.Make(code.OpConstant, 4), // 2.0
				code.Make(code.OpObject, 2),   // {"c": 2}
				code.Make(code.OpObject, 4),
			}},
	}
	runCompilerTestCases(t, cases)
}

func TestObjectBracketIndexAccess(t *testing.T) {
	cases := []compilerTestCase{
		// Access property by string literal
		{`{"a": 1}["a"]`, []any{"a", 1.0, "a"}, []code.Instructions{
			code.Make(code.OpConstant, 0), // "a"
			code.Make(code.OpConstant, 1), // 1.0
			code.Make(code.OpObject, 2),
			code.Make(code.OpConstant, 2), // "a"
			code.Make(code.OpIndex),
		}},
		// Access property by context variable
		{`obj["key"]`, []any{"key"}, []code.Instructions{
			code.Make(code.OpContextVar, 0), // obj
			code.Make(code.OpConstant, 0),   // "key"
			code.Make(code.OpIndex),
		}},
		// Nested object property access
		{`{"a": {"b": 2}}["a"]["b"]`, []any{"a", "b", 2.0, "a", "b"}, []code.Instructions{
			code.Make(code.OpConstant, 0), // "a"
			code.Make(code.OpConstant, 1), // "b"
			code.Make(code.OpConstant, 2), // 2.0
			code.Make(code.OpObject, 2),
			code.Make(code.OpObject, 2),
			code.Make(code.OpConstant, 3), // "a"
			code.Make(code.OpIndex),
			code.Make(code.OpConstant, 4), // "b"
			code.Make(code.OpIndex),
		}},
	}
	runCompilerTestCases(t, cases)
}
func TestFunctionCalls(t *testing.T) {
	cases := []compilerTestCase{
		// Simple function call with no arguments
		{`foo()`, []any{"foo"}, []code.Instructions{
			code.Make(code.OpCallFunction, 0, 0),
		}},
		// Function call with one argument
		{`bar(42)`, []any{42.0, "bar"}, []code.Instructions{
			code.Make(code.OpConstant, 0), // 42.0
			code.Make(code.OpCallFunction, 1, 1),
		}},
		// Function call with multiple arguments
		{`baz(1, 2, 3)`, []any{1.0, 2.0, 3.0, "baz"}, []code.Instructions{
			code.Make(code.OpConstant, 0), // 1.0
			code.Make(code.OpConstant, 1), // 2.0
			code.Make(code.OpConstant, 2), // 3.0
			code.Make(code.OpCallFunction, 3, 3),
		}},
		// Function call with context variables
		{`doSomething(foo, bar)`, []any{"doSomething"},
			[]code.Instructions{
				code.Make(code.OpContextVar, 0), // foo
				code.Make(code.OpContextVar, 1), // bar
				code.Make(code.OpCallFunction, 0, 2),
			}},
		// Function call with mixed arguments
		{`calculate(1, bar, 3)`, []any{1.0, 3.0, "calculate"}, []code.Instructions{
			code.Make(code.OpConstant, 0),   // 1.0
			code.Make(code.OpContextVar, 0), // bar
			code.Make(code.OpConstant, 1),   // 3.0
			code.Make(code.OpCallFunction, 2, 3),
		}},
	}
	runCompilerTestCases(t, cases)
}

func TestPipeExpression(t *testing.T) {
	cases := []compilerTestCase{
		// Access property by string literal
		{`[1,2] |map: $item * 2`, []any{1, 2, "pipe", 2, "map"}, []code.Instructions{
			code.Make(code.OpConstant, 0), // "a"
			code.Make(code.OpConstant, 1), // 1.0
			code.Make(code.OpArray, 2),
			code.Make(code.OpPipe, 2, 3),
			code.Make(code.OpIdentifier, 0), // $item
			code.Make(code.OpConstant, 4),   // 2
			code.Make(code.OpMul),
			code.Make(code.OpPipe, 5, 6),
		}},
		// Multiple pipes with different operators
		{`[1,2,3] |filter: $item > 1 |map: $item * 2`, []any{1.0, 2.0, 3.0, 1.0, 2.0, 1, "filter", "", 2, "map", ""}, []code.Instructions{
			code.Make(code.OpConstant, 0), // 1.0
			code.Make(code.OpConstant, 1), // 2.0
			code.Make(code.OpConstant, 2), // 3.0
			code.Make(code.OpArray, 3),
			code.Make(code.OpPipe, 3, 4),
			code.Make(code.OpIdentifier, 0), // $item
			code.Make(code.OpConstant, 5),   // 1.0
			code.Make(code.OpGreaterThan),
			code.Make(code.OpPipe, 6, 7),
			code.Make(code.OpIdentifier, 0), // $item
			code.Make(code.OpConstant, 8),   // 2.0
			code.Make(code.OpMul),
			code.Make(code.OpPipe, 9, 10),
		}},
	}
	runCompilerTestCases(t, cases)
}

func TestTernaryOperatorCompilation(t *testing.T) {
	cases := []compilerTestCase{
		{
			`true ? 1 : 2`,
			[]any{1.0, 2.0},
			[]code.Instructions{
				code.Make(code.OpTrue),            // 0
				code.Make(code.OpJumpIfFalsy, 10), // 1..3 (elsePos=10)
				code.Make(code.OpConstant, 0),     // 4..6 (consequent 1)
				code.Make(code.OpJump, 14),        // 7..9 (endPos=14)
				code.Make(code.OpPop),             // 10 (discard falsy cond when jumping)
				code.Make(code.OpConstant, 1),     // 11..13 (alternate 2)
			},
		},
		{
			`false ? 1 : 2`,
			[]any{1.0, 2.0},
			[]code.Instructions{
				code.Make(code.OpFalse),           // 0
				code.Make(code.OpJumpIfFalsy, 10), // 1..3
				code.Make(code.OpConstant, 0),     // 4..6
				code.Make(code.OpJump, 14),        // 7..9
				code.Make(code.OpPop),             // 10
				code.Make(code.OpConstant, 1),     // 11..13
			},
		},
		{
			`(true ? 5 : 10) + 3`,
			[]any{5.0, 10.0, 3.0},
			[]code.Instructions{
				code.Make(code.OpTrue),            // 0
				code.Make(code.OpJumpIfFalsy, 10), // 1..3
				code.Make(code.OpConstant, 0),     // 4..6 (5)
				code.Make(code.OpJump, 14),        // 7..9
				code.Make(code.OpPop),             // 10
				code.Make(code.OpConstant, 1),     // 11..13 (10)
				code.Make(code.OpConstant, 2),     // 14..16 (3)
				code.Make(code.OpAdd),             // 17
			},
		},
		{
			`false ? 1 : true ? 2 : 3`,
			[]any{1.0, 2.0, 3.0},
			[]code.Instructions{
				// Outer condition
				code.Make(code.OpFalse),           // 0
				code.Make(code.OpJumpIfFalsy, 10), // 1..3 (outer elsePos=10)
				code.Make(code.OpConstant, 0),     // 4..6 (outer consequent 1)
				code.Make(code.OpJump, 25),        // 7..9 (outer endPos=25)
				code.Make(code.OpPop),             // 10 (discard falsy outer cond)

				// Inner condition (starts at 11)
				code.Make(code.OpTrue),            // 11
				code.Make(code.OpJumpIfFalsy, 21), // 12..14 (inner elsePos=21)
				code.Make(code.OpConstant, 1),     // 15..17 (inner consequent 2)
				code.Make(code.OpJump, 25),        // 18..20 (inner endPos=25)
				code.Make(code.OpPop),             // 21 (discard falsy inner cond)
				code.Make(code.OpConstant, 2),     // 22..24 (inner alternate 3)
			},
		},
	}
	runCompilerTestCases(t, cases)
}
