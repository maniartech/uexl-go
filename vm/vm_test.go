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
	vm := vm.New(&compiler.ByteCode{}, vm.Builtins)
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
		vm := vm.New(comp.ByteCode(), vm.Builtins)
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
	case *parser.StringLiteral:
		if expectedStr, ok := expected.(string); ok {
			if actual.Value != expectedStr {
				return fmt.Errorf("expected %q, got %q", expectedStr, actual.Value)
			}
		} else {
			return fmt.Errorf("expected a string, got %T", expected)
		}
	case *parser.ArrayLiteral:
		if expectedArray, ok := expected.([]any); ok {
			if len(actual.Elements) != len(expectedArray) {
				return fmt.Errorf("expected array of length %d, got %d", len(expectedArray), len(actual.Elements))
			}
			for i, elem := range actual.Elements {
				if err := testExpectedObject(t, expectedArray[i], elem); err != nil {
					return fmt.Errorf("error at index %d: %s", i, err)
				}
			}
		} else {
			return fmt.Errorf("expected an array, got %T", expected)
		}
	case *parser.ObjectLiteral:
		if expectedMap, ok := expected.(map[string]any); ok {
			if len(actual.Properties) != len(expectedMap) {
				return fmt.Errorf("expected object with %d properties, got %d", len(expectedMap), len(actual.Properties))
			}
			for key, value := range actual.Properties {
				if expectedValue, exists := expectedMap[key]; exists {
					if err := testExpectedObject(t, expectedValue, value); err != nil {
						return fmt.Errorf("error for key %q: %s", key, err)
					}
				} else {
					return fmt.Errorf("unexpected key %q in object", key)
				}
			}
			// Also check that all expected keys are present
			for expectedKey := range expectedMap {
				if _, exists := actual.Properties[expectedKey]; !exists {
					return fmt.Errorf("missing expected key %q in object", expectedKey)
				}
			}
		} else {
			return fmt.Errorf("expected an object, got %T", expected)
		}
	default:
		return fmt.Errorf("unsupported type: %T (value: %v)", actual, actual)
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
		{"(5 - 2) < (2 * 2)", true},
		{"(10 / 2) >= 5", true},
		{"(10 / 2) <= 5", true},
	}
	runVmTests(t, tests)
}
func TestBitwiseOperations(t *testing.T) {
	tests := []vmTestCase{
		{"5 & 3", 1},             // 0101 & 0011 = 0001
		{"5 | 3", 7},             // 0101 | 0011 = 0111
		{"5 ^ 3", 6},             // 0101 ^ 0011 = 0110
		{"8 << 2", 32},           // 1000 << 2 = 100000
		{"32 >> 3", 4},           // 100000 >> 3 = 100
		{"15 & 7", 7},            // 1111 & 0111 = 0111
		{"15 | 7", 15},           // 1111 | 0111 = 1111
		{"15 ^ 7", 8},            // 1111 ^ 0111 = 1000
		{"1 << 4", 16},           // 0001 << 4 = 10000
		{"16 >> 2", 4},           // 10000 >> 2 = 100
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
func TestStringLiterals(t *testing.T) {
	tests := []vmTestCase{
		{`"hello"`, "hello"},
		{`"world"`, "world"},
		{`""`, ""},
	}
	runVmTests(t, tests)
}

func TestStringConcatenation(t *testing.T) {
	tests := []vmTestCase{
		{`"hello" + " " + "world"`, "hello world"},
		{`"foo" + "bar"`, "foobar"},
		{`"a" + "" + "b"`, "ab"},
	}
	runVmTests(t, tests)
}

func TestStringComparison(t *testing.T) {
	tests := []vmTestCase{
		{`"abc" == "abc"`, true},
		{`"abc" == "def"`, false},
		{`"abc" != "def"`, true},
		{`"abc" != "abc"`, false},
	}
	runVmTests(t, tests)
}

func TestStringLengthFunction(t *testing.T) {
	tests := []vmTestCase{
		{`len("hello")`, 5},
		{`len("")`, 0},
		{`len("abcde")`, 5},
	}
	runVmTests(t, tests)
}

func TestStringSubstringFunction(t *testing.T) {
	tests := []vmTestCase{
		{`substr("hello", 0, 2)`, "he"},
		{`substr("world", 1, 3)`, "orl"},
		{`substr("foobar", 3, 3)`, "bar"},
	}
	runVmTests(t, tests)
}

func TestStringContainsFunction(t *testing.T) {
	tests := []vmTestCase{
		{`contains("hello", "ll")`, true},
		{`contains("hello", "z")`, false},
		{`contains("foobar", "foo")`, true},
	}
	runVmTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []any{}},
		{"[1]", []any{1}},
		{"[1, 2, 3]", []any{1, 2, 3}},
		{"[true, false, 1, \"hello\"]", []any{true, false, 1, "hello"}},
		{"[1, 1 + 4, 3 * 4]", []any{1, 5, 12}},
	}
	runVmTests(t, tests)
}
func TestArrayIndexing(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"[10, 20, 30, 40][3]", 40},
		{"[true, false, true][1]", false},
		{`["a", "b", "c"][2]`, "c"},
		{"[1 + 2, 3 * 4, 5 - 1][1]", 12},
		{"[1, 2, 3][0] == 1", true},
		{"[1, 2, 3][1] + 5", 7},
		{"[1, 2, 3][2] * 2", 6},
		{"[1, 2, 3][0] + [4, 5, 6][2]", 7},
		{"[[1,2],[3,4]][1][0]", 3},
		{"[1, [2, 3+5], 4][1][1]", 8},
	}
	runVmTests(t, tests)
}
func TestObjectLiterals(t *testing.T) {
	tests := []vmTestCase{
		{`{}`, map[string]any{}},
		{`{"a": 1}`, map[string]any{"a": 1}},
		{`{"a": 1, "b": 2}`, map[string]any{"a": 1, "b": 2}},
		{`{"x": true, "y": false, "z": "hello"}`, map[string]any{"x": true, "y": false, "z": "hello"}},
		{`{"num": 42, "arr": [1,2,3], "obj": {"nested": "yes"}}`, map[string]any{
			"num": 42,
			"arr": []any{1, 2, 3},
			"obj": map[string]any{"nested": "yes"},
		}},
	}
	runVmTests(t, tests)
}

func TestObjectIndexing(t *testing.T) {
	tests := []vmTestCase{
		{`{"a": 1, "b": 2}["a"]`, 1},
		{`{"a": 1, "b": 2}["b"]`, 2},
		{`{"x": true, "y": false}["y"]`, false},
		{`{"foo": "bar"}["foo"]`, "bar"},
		{`{"arr": [1,2,3]}["arr"][1]`, 2},
		{`{"obj": {"nested": 99}}["obj"]["nested"]`, 99},
	}
	runVmTests(t, tests)
}
