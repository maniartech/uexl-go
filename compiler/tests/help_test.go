package compiler_test

import (
	"fmt"
	"testing"

	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/compiler"
	"github.com/maniartech/uexl_go/parser"
)

func parse(input string) parser.Node {
	p := parser.NewParser(input)
	node, err := p.Parse()
	if err != nil {
		panic(err)
	}
	return node
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}
	for _, ins := range s {
		out = append(out, ins...)
	}
	return out
}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	concatted := concatInstructions(expected)
	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length.\nwant=%q\ngot =%q",
			concatted, actual)
	}
	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot =%q",
				i, concatted, actual)
		}
	}
	return nil
}

func testConstants(
	t *testing.T,
	expected []any,
	actual []any,
) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, want=%d",
			len(actual), len(expected))
	}
	for i, constant := range expected {
		if constant == nil { // sentinel: skip validation for this constant (e.g., InstructionBlock)
			continue
		}
		switch constant := constant.(type) {
		case float64:
		case int:
			if err := testNumerLiteral(float64(constant), actual[i]); err != nil {
				return fmt.Errorf("test case %d: %s", i, err)
			}
		case string:
			if err := testStringLiteral(constant, actual[i]); err != nil {
				return fmt.Errorf("test case %d: %s", i, err)
			}
		default:
			return fmt.Errorf("unknown constant type %T", constant)
		}
	}
	return nil
}

func testNumerLiteral(expected float64, actual any) error {
	// Accept either a *parser.NumberLiteral (old representation) or a raw float64 (current representation)
	switch v := actual.(type) {
	case *parser.NumberLiteral:
		if v.Value != expected {
			return fmt.Errorf("wrong number literal. got=%v, want=%v", v.Value, expected)
		}
		return nil
	case float64:
		if v != expected {
			return fmt.Errorf("wrong number literal. got=%v, want=%v", v, expected)
		}
		return nil
	case int:
		if float64(v) != expected {
			return fmt.Errorf("wrong number literal. got=%v, want=%v", v, expected)
		}
		return nil
	}
	return fmt.Errorf("expected a number literal, got %T", actual)
}

func testStringLiteral(expected string, actual any) error {
	// Accept either a *parser.StringLiteral (old) or raw string (current)
	switch v := actual.(type) {
	case *parser.StringLiteral:
		if v.Value != expected {
			return fmt.Errorf("wrong string literal. got=%q, want=%q", v.Value, expected)
		}
		return nil
	case string:
		if v != expected {
			return fmt.Errorf("wrong string literal. got=%q, want=%q", v, expected)
		}
		return nil
	}
	return fmt.Errorf("expected a string literal, got %T", actual)
}

func runCompilerTestCases(t *testing.T, cases []compilerTestCase) {
	for i, tc := range cases {
		program := parse(tc.input)
		compiler := compiler.New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("Test case %d: Compile error: %s", i, err)
		}
		byteCode := compiler.ByteCode()
		err = testConstants(t, tc.expectedConstants, byteCode.Constants)
		if err != nil {
			t.Fatalf("Test case %d: %s", i, err)
		}
		err = testInstructions(tc.expectedInstructions, byteCode.Instructions)
		if err != nil {
			t.Fatalf("Test case %d: %s", i, err)
		}
		t.Logf("Test case %d passed", i)
	}
}
