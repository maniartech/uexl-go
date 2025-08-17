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

func TestVM(t *testing.T) {
	vm := vm.New(vm.LibContext{
		Functions:    nil,
		PipeHandlers: nil,
	})
	if vm == nil {
		t.Fatal("failed to create VM")
	}
}

type vmTestCase struct {
	input    string
	expected any
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()
	for i, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("[case %d] compiler error: %s", i+1, err)
		}
		vm := vm.New(vm.LibContext{
			Functions:    vm.Builtins,
			PipeHandlers: vm.DefaultPipeHandlers,
		})
		bytecode := comp.ByteCode()
		output, err := vm.Run(bytecode)
		if err != nil {
			t.Fatalf("[case %d] vm error: %s", i+1, err)
		}
		err = testExpectedObject(t, tt.expected, output)
		if err != nil {
			t.Fatalf("[case %d] testExpectedObject error: %s", i+1, err)
		}
	}
}

// testExpectedObject compares the expected value with the actual value from the VM stack.
func testExpectedObject(t *testing.T, expected any, actual any) error {
	switch v := actual.(type) {
	case float64:
	case int:
		// converting the value to float64 for comparison
		if float64(v) != actual.(float64) {
			return fmt.Errorf("expected %f, got %f", float64(v), actual.(float64))
		}
	case bool:
		if a, ok := actual.(bool); ok {
			if v != a {
				return fmt.Errorf("expected %t, got %t", v, a)
			}
		} else {
			return fmt.Errorf("expected a boolean, got %T", actual)
		}
	case string:
		if a, ok := actual.(string); ok {
			if v != a {
				return fmt.Errorf("expected %q, got %q", v, a)
			}
		} else {
			return fmt.Errorf("expected a string, got %T", actual)
		}
	case []any:
		e, ok := expected.([]any)
		if !ok {
			return fmt.Errorf("expected an array, got %T", expected)
		}
		if len(v) != len(e) {
			return fmt.Errorf("expected array of length %d, got %d", len(e), len(v))
		}
		for i := range v {
			if err := testExpectedObject(t, e[i], v[i]); err != nil {
				return fmt.Errorf("error at index %d: %s", i, err)
			}
		}
	case map[string]any:
		e, ok := expected.(map[string]any)
		if !ok {
			return fmt.Errorf("expected an object, got %T", expected)
		}
		if len(v) != len(e) {
			return fmt.Errorf("expected object with %d properties, got %d", len(e), len(v))
		}
		for key, val := range v {
			expVal, exists := e[key]
			if !exists {
				return fmt.Errorf("unexpected key %q in object", key)
			}
			if err := testExpectedObject(t, expVal, val); err != nil {
				return fmt.Errorf("error for key %q: %s", key, err)
			}
		}
		for key := range e {
			if _, exists := v[key]; !exists {
				return fmt.Errorf("missing expected key %q in object", key)
			}
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
		{"2 ** 3", 8}, // Power operator test

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

func TestLogicalShortCircuit(t *testing.T) {
	tests := []vmTestCase{
		// || returns first truthy value
		{`false || true`, true},
		{`false || false || true`, true},
		{`false || 0 || "hello"`, "hello"},
		{`false || 0 || ""`, ""}, // all falsy, returns last

		// && returns first falsy value, otherwise last value
		{`true && false`, false},
		{`true && true && false`, false},
		{`true && 1 && "hello"`, "hello"},
		{`true && 1 && ""`, ""},     // "" is falsy
		{`false && 0 && ""`, false}, // all falsy, returns false

		// Chained with numbers and strings
		{`0 || 42`, 42},
		{`"foo" || "bar"`, "foo"},
		{`"" || "bar"`, "bar"},
		{`1 && 2 && 3`, 3},
		{`1 && 0 && 3`, 0},
		{`0 && 1`, 0},

		// Nested expressions
		{`(false || 0) && (true || "baz")`, 0},
		{`(true && 1) || (false && "baz")`, 1},
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
		{`contains("f"+"oobar", "foo")`, true},
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

func TestPipeFunction(t *testing.T) {
	tests := []vmTestCase{
		{`"foo" as $foo |pipe: $foo + "bar"`, "foobar"},
		{"[1,2] |map: $item * 2", []any{2, 4}},
		{"[1,2] |map: $item * $index", []any{0, 2}},
		{"[1,2] |map: $item * 2 |map: $item + 1", []any{3, 5}},
		{"[1,2,3,4,5,6] |filter: $item > 2", []any{3, 4, 5, 6}},

		// Reduce: sum all items
		{"[1,2,3,4] |reduce: ($acc || 0) + $item", 10},

		// TODO: Reduce to an object
		{"[1,2,3,4] |reduce: set($acc || {}, $index, $item)", map[string]any{
			"0": 1,
			"1": 2,
			"2": 3,
			"3": 4,
		}},

		// Find: first item greater than 2
		{"[1,2,3,4] |find: $item > 2", 3},

		// Some: any item is even
		{"[1,2,3,4] |some: $item % 2 == 0", true},
		{"[1,3,5] |some: $item % 2 == 0", false},

		// Every: all items are positive
		{"[1,2,3,4] |every: $item > 0", true},
		{"[1,2,3,0] |every: $item > 0", false},

		// Unique: remove duplicates
		{"[1,2,2,3,1,4] |unique: $item", []any{1, 2, 3, 4}},

		// Sort: sort by value
		{"[3,1,2] |sort: $item", []any{1, 2, 3}},
		// Sort: sort by computed value
		{"[3,1,2] |sort: $item * -1", []any{3, 2, 1}},

		// GroupBy: group by even/odd
		// {`[1,2,3,4] |groupBy: $item % 2`, map[string]any{
		// 	"1": []any{1, 3},
		// 	"0": []any{2, 4},
		// }},

		// // Window: window size 2, sum each window
		{"[1,2,3,4] |window: $window[0] + $window[1]", []any{3, 5, 7}},

		// Chunk: chunk size 2, sum each chunk
		// {"[1,2,3,4,5] |chunk: $chunk[0] + ($chunk[1] ?? 0)", []any{3, 7, 5}},
	}
	runVmTests(t, tests)
}
