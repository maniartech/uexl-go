package vm_test

import (
	"fmt"
	"testing"

	"github.com/maniartech/uexl_go/compiler"
	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/vm"
)

// parse is a helper function to parse the input string
// this runs the lexer and parser to create an AST
func parse(input string) parser.Node {
	p := parser.NewParser(input)
	node, err := p.Parse()
	if err != nil {
		fmt.Printf("Parse error: %s\n", err)
		return nil
	}
	if node == nil {
		fmt.Println("Parse error: no node returned")
		return nil
	}
	return node
}

type vmTestCases struct {
	input    string
	expected any
}

func TestVM(t *testing.T) {
	vm := vm.New(&compiler.ByteCode{})
	if vm == nil {
		t.Fatal("failed to create VM")
	}
}

type vmTestCase struct {
	input    string
	expected interface{}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()
	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}
		vm := vm.New(comp.ByteCode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}
		stackElem := vm.LastPoppedStackElem()
		err = testExpectedObject(t, tt.expected, stackElem)
		if err != nil {
			t.Fatalf("testExpectedObject error: %s", err)
		}

	}
}

// testExpectedObject compares the expected value with the actual value from the VM stack.
func testExpectedObject(t *testing.T, expected any, actual parser.Node) error {
	switch actual := actual.(type) {
	case *parser.NumberLiteral:
		// NumberLiteral can be integer or float
		if expectedInt, ok := expected.(int); ok {
			if actual.Value != float64(expectedInt) {
				return fmt.Errorf("expected %d, got %f", expectedInt, actual.Value)
			}
		} else if expectedFloat, ok := expected.(float64); ok {
			if actual.Value != expectedFloat {
				return fmt.Errorf("expected %f, got %f", expectedFloat, actual.Value)
			}
		} else {
			return fmt.Errorf("expected a number, got %T", expected)
		}
	case *parser.BooleanLiteral:
		if expectedBool, ok := expected.(bool); ok {
			if actual.Value != expectedBool {
				return fmt.Errorf("expected %t, got %t", expectedBool, actual.Value)
			}
		} else {
			return fmt.Errorf("expected a boolean, got %T", expected)
		}
	}
	return nil
}

func TestNumberArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"-(-20 + 10)", 10},
		{"--10", 10},

		// Floating point tests
		{"1.5", 1.5},
		{"2.5 + 3.5", 6.0},
		{"5.5 - 2.2", 3.3},
		{"2.0 * 3.5", 7.0},
		{"7.5 / 2.5", 3.0},
		{"1.1 + 2.2 + 3.3", 6.6},
		{"10.0 - 5.5", 4.5},
		{"2.5 * 2.0 * 2.0", 10.0},
		{"5.0 * (2.0 + 3.0)", 25.0},
		{"-5.5", -5.5},
		{"-10.1", -10.1},
		{"-50.5 + 100.5 + -50.0", 0.0},
		{"(5.5 + 10.5 * 2.0 + 15.0 / 3.0) * 2.0 + -10.0", 53.0},
		{"-(-20.5 + 10.5)", 10.0},
		{"--10.5", 10.5},
	}
	runVmTests(t, tests)
}

func TestNumberComparison(t *testing.T) {
	// Contains tests for number comparison operations (float and integer)
	// starts with simple comparisons and includes more complex expressions
	tests := []vmTestCase{
		{"1 == 1", true},
		{"1 == 2", false},
		{"1 != 2", true},
		{"2 != 2", false},
		{"2 > 1", true},
		{"1 > 2", false},
		{"2 >= 2", true},
		{"2 >= 3", false},
		{"1 < 2", true},
		{"2 < 1", false},
		{"2 <= 2", true},
		{"3 <= 2", false},
		{"1.5 == 1.5", true},
		{"1.5 != 2.5", true},
		{"2.5 > 1.5", true},
		{"1.5 < 2.5", true},
		{"(1 + 2) == 3", true},
		{"(2 * 2) > (3 + 1)", false},
		{"(5 - 2) < (2 * 2)", false},
		{"(10 / 2) >= 5", true},
		{"(10 / 2) <= 5", true},
	}
	runVmTests(t, tests)
}
func TestBitwiseOperations(t *testing.T) {
	tests := []vmTestCase{
		{"5 & 3", 1},      // 0101 & 0011 = 0001
		{"5 | 3", 7},      // 0101 | 0011 = 0111
		{"5 ^ 3", 6},      // 0101 ^ 0011 = 0110
		{"8 << 2", 32},    // 1000 << 2 = 100000
		{"32 >> 3", 4},    // 100000 >> 3 = 100
		{"15 & 7", 7},     // 1111 & 0111 = 0111
		{"15 | 7", 15},    // 1111 | 0111 = 1111
		{"15 ^ 7", 8},     // 1111 ^ 0111 = 1000
		{"1 << 4", 16},    // 0001 << 4 = 10000
		{"16 >> 2", 4},    // 10000 >> 2 = 100
		{"(5 & 3) | (2 ^ 1)", 3}, // (1) | (3) = 3
	}
	runVmTests(t, tests)
}
func TestBooleanLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
	}
	runVmTests(t, tests)
}

func TestBooleanComparisons(t *testing.T) {
	tests := []vmTestCase{
		{"true == true", true},
		{"true == false", false},
		{"false == false", true},
		{"true != false", true},
		{"false != true", true},
		{"true != true", false},
		{"false != false", false},
	}
	runVmTests(t, tests)
}

func TestBooleanLogic(t *testing.T) {
	tests := []vmTestCase{
		{"true && true", true},
		{"true && false", false},
		{"false && false", false},
		{"true || false", true},
		{"false || false", false},
		{"!true", false},
		{"!false", true},
		{"!(true == false)", true},
		{"!(1 == 1)", false},
		{"(1 == 1) && (2 == 2)", true},
		{"(1 == 2) || (2 == 2)", true},
	}
	runVmTests(t, tests)
}
