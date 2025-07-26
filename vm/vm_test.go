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
		testExpectedObject(t, tt.expected, stackElem)
	}
}

// testExpectedObject compares the expected value with the actual value from the VM stack.
func testExpectedObject(t *testing.T, expected any, actual parser.Node) {
	switch actual := actual.(type) {
	case *parser.NumberLiteral:
		// NumberLiteral can be integer or float
		if expectedInt, ok := expected.(int); ok {
			if actual.Value != float64(expectedInt) {
				t.Errorf("expected %d, got %f", expectedInt, actual.Value)
			}
		} else if expectedFloat, ok := expected.(float64); ok {
			if actual.Value != expectedFloat {
				t.Errorf("expected %f, got %f", expectedFloat, actual.Value)
			}
		} else {
			t.Errorf("expected a number, got %T", expected)
		}
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		// {"1", 1},
		// {"2", 2},
		// {"1 + 2", 3},
		// {"1 - 2", -1},
		// {"1 * 2", 2},
		// {"4 / 2", 2},
		// {"50 / 2 * 2 + 10 - 5", 55},
		// {"5 + 5 + 5 + 5 - 10", 10},
		// {"2 * 2 * 2 * 2 * 2", 32},
		// {"5 * 2 + 10", 20},
		// {"5 + 2 * 10", 25},
		// {"5 * (2 + 10)", 60},
		// {"-5", -5},
		// {"-10", -10},
		// {"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		// {"-(-20 + 10)", 10},
		// {"--10", 10},
	}
	runVmTests(t, tests)
}

