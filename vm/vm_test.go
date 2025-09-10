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

func runVmTests(t *testing.T, tests []vmTestCase, contextValues ...map[string]any) {
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
		var context map[string]any
		if len(contextValues) > 0 {
			context = contextValues[0]
		} else {
			context = make(map[string]any)
		}
		// Set the context values in the VM
		output, err := vm.Run(bytecode, context)
		if err != nil {
			t.Fatalf("[case %d] vm error: %s", i+1, err)
		}
		err = testExpectedObject(t, tt.expected, output)
		if err != nil {
			t.Fatalf("[case %d][input: %q] testExpectedObject error: %s", i+1, tt.input, err)
		}
	}
}

// testExpectedObject compares the expected value with the actual value from the VM stack.
func testExpectedObject(t *testing.T, expected any, actual any) error {
	switch v := expected.(type) {
	case float64:
		// converting the value to float64 for comparison
		if a, ok := actual.(float64); ok {
			if v != a {
				return fmt.Errorf("expected %f, got %f", v, a)
			}
		} else {
			return fmt.Errorf("expected a float64, got %T", actual)
		}
	case int:
		if a, ok := actual.(int); ok {
			if v != a {
				return fmt.Errorf("expected %d, got %d", v, a)
			}
		} else {
			return fmt.Errorf("expected a int, got %T", actual)
		}
	case bool:
		if a, ok := actual.(bool); ok {
			if v != a {
				return fmt.Errorf("expected %t, got %t", expected, actual)
			}
		} else {
			return fmt.Errorf("expected a boolean, got %T", actual)
		}
	case string:
		if a, ok := actual.(string); ok {
			if v != a {
				return fmt.Errorf("expected %q, got %q", expected, actual)
			}
		} else {
			return fmt.Errorf("expected a string, got %T", actual)
		}
	case []any:
		e, ok := expected.([]any)
		if !ok {
			return fmt.Errorf("expected an array, got %T", actual)
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
	case nil:
		if actual != nil {
			return fmt.Errorf("expected nil, got %T", actual)
		}
	default:
		return fmt.Errorf("unsupported type: %T (value: %v)", actual, actual)
	}
	return nil
}

func TestNumberArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1.0},
		{"2", 2.0},
		{"1 + 2", 3.0},
		{"1 - 2", -1.0},
		{"1 * 2", 2.0},
		{"4 / 2", 2.0},
		{"50 / 2 * 2 + 10 - 5", 55.0},
		{"5 + 5 + 5 + 5 - 10", 10.0},
		{"2 * 2 * 2 * 2 * 2", 32.0},
		{"5 * 2 + 10", 20.0},
		{"5 + 2 * 10", 25.0},
		{"5 * (2 + 10)", 60.0},
		{"-5", -5.0},
		{"-10", -10.0},
		{"-50 + 100 + -50", 0.0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50.0},
		{"-(-20 + 10)", 10.0},
		{"--10", 10.0},
		{"2 ** 3", 8.0}, // Power operator test

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
		{"5 & 3", 1.0},             // 0101 & 0011 = 0001
		{"5 | 3", 7.0},             // 0101 | 0011 = 0111
		{"5 ^ 3", 6.0},             // 0101 ^ 0011 = 0110
		{"8 << 2", 32.0},           // 1000 << 2 = 100000
		{"32 >> 3", 4.0},           // 100000 >> 3 = 100
		{"15 & 7", 7.0},            // 1111 & 0111 = 0111
		{"15 | 7", 15.0},           // 1111 | 0111 = 1111
		{"15 ^ 7", 8.0},            // 1111 ^ 0111 = 1000
		{"1 << 4", 16.0},           // 0001 << 4 = 10000
		{"16 >> 2", 4.0},           // 10000 >> 2 = 100
		{"(5 & 3) | (2 ^ 1)", 3.0}, // (1) | (3) = 3
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

func TestUnaryLogicalNot(t *testing.T) {
	tests := []vmTestCase{
		{`!true`, false},
		{`!false`, true},
		{`!1`, false},
		{`!0`, true},
		{`!"hello"`, false},
		{`!""`, true},
		{`!!true`, true},
		{`!!false`, false},
		{`!!1`, true},
		{`!!0`, false},
		{`!!"hello"`, true},
		{`!!""`, false},
		{`!(!true)`, true},
		{`!(!false)`, false},
		{`!!(1 + 1)`, true},
		{`!!(0)`, false},
		{`!([1,2,3])`, false},
		{`!([])`, true},
		{`!!([1,2,3])`, true},
		{`!!([])`, false},
		{`!({})`, true},
		{`!!({})`, false},
		{`!({"a":1})`, false},
		{`!!({"a":1})`, true},

		// More complex and nested cases
		{`!([1, false, 0, ""])`, false},
		{`!([false, 0, ""])`, false},
		{`!!([false, 0, ""])`, true},
		{`!([1, 2, 3][0])`, false},
		{`!([1, 2, 3][2])`, false},
		{`!({"a": 0, "b": false}["a"])`, true},
		{`!({"a": 1, "b": false}["a"])`, false},
		{`!({"a": 1, "b": false}["b"])`, true},
		{`!!({"a": 1, "b": false}["a"])`, true},
		{`!!({"a": 1, "b": false}["b"])`, false},
		{`!([1,2,3][0] && true)`, false},
		{`!([1,2,3][0] && false)`, true},
		{`!([1,2,3][0] || false)`, false},
		{`!((1 && 0) || ("" && true))`, true},
		{`!!((1 && 0) || ("" && true))`, false},
		{`!((1 || 0) && ("foo" || ""))`, false},
		{`!!((1 || 0) && ("foo" || ""))`, true},
		{`!(!([1,2,3][1] > 1))`, true},
		{`!!(!([1,2,3][1] > 1))`, false},
		{`!(!([1,2,3][1] < 1))`, false},
		{`!!(!([1,2,3][1] < 1))`, true},
		{`!(!([1,2,3][1] == 2))`, true},
		{`!!(!([1,2,3][1] == 2))`, false},
		{`!(!([1,2,3][1] != 2))`, false},
		{`!!(!([1,2,3][1] != 2))`, true},
		{`!(!([1,2,3][1] >= 2))`, true},
		{`!!(!([1,2,3][1] >= 2))`, false},
		{`!(!([1,2,3][1] <= 2))`, true},
		{`!!(!([1,2,3][1] <= 2))`, false},
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
		{`0 || 42`, 42.0},
		{`"foo" || "bar"`, "foo"},
		{`"" || "bar"`, "bar"},
		{`1 && 2 && 3`, 3.0},
		{`1 && 0 && 3`, 0.0},
		{`0 && 1`, 0.0},

		// Nested expressions
		{`(false || 0) && (true || "baz")`, 0.0},
		{`(true && 1) || (false && "baz")`, 1.0},
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
		{`len("hello")`, 5.0},
		{`len("")`, 0.0},
		{`len("abcde")`, 5.0},
	}
	runVmTests(t, tests)
}

func TestStringSubstringFunction(t *testing.T) {
	tests := []vmTestCase{
		{`substr("hello", 0, 2)`, "he"},
		{`substr("world", 1, 3)`, "orl"},
		{`substr("foobar", 1+2, 3)`, "bar"},
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
func TestIndexing(t *testing.T) {
	tests := []vmTestCase{
		// Array indexing with []
		// {"[1, 2, 3][0]", 1.0},
		// {"[1, 2, 3][1]", 2.0},
		// {"[1, 2, 3][2]", 3.0},
		// {"[10, 20, 30, 40][3]", 40.0},
		// {"[true, false, true][1]", false},
		// {`["a", "b", "c"][2]`, "c"},
		// {"[1 + 2, 3 * 4, 5 - 1][1]", 12.0},
		// {"[1, 2, 3][0] == 1", true},
		// {"[1, 2, 3][1] + 5", 7.0},
		// {"[1, 2, 3][2] * 2", 6.0},
		// {"[1, 2, 3][0] + [4, 5, 6][2]", 7.0},
		// {"[[1,2],[3,4]][1][0]", 3.0},
		// {"[1, [2, 3+5], 4][1][1]", 8.0},

		// // String indexing with []
		// {`"hello"[0]`, "h"},
		// {`"hello"[1]`, "e"},
		// {`"hello"[4]`, "o"},
		// {`"world"[2]`, "r"},
		// {`"abcde"[3]`, "d"},
		// {`"foo"[1] == "o"`, true},

		// Array indexing with dot notation
		{"[1, 2, 3].0", 1.0},
		{"[1, 2, 3].1", 2.0},
		{"[1, 2, 3].2", 3.0},
		{"[10, 20, 30, 40].3", 40.0},
		{"[true, false, true].1", false},
		{`["a", "b", "c"].2`, "c"},
		{"[1 + 2, 3 * 4, 5 - 1].1", 12.0},
		{"[1, 2, 3].0 == 1", true},
		{"[1, 2, 3].1 + 5", 7.0},
		{"[1, 2, 3].2 * 2", 6.0},
		{"[1, 2, 3].0 + [4, 5, 6].2", 7.0},
		// {"[[1,2],[3,4]].1.2", 4}, // Fix operator being returned as 1.2 instead 1st element's 2nd element
		{"[1, [2, 3+5], 4].1.1", 8.0},

		// String indexing with dot notation
		{`"hello".0`, "h"},
		{`"hello".1`, "e"},
		{`"hello".4`, "o"},
		{`"world".2`, "r"},
		{`"abcde".3`, "d"},
		{`"foo".1 == "o"`, true},
	}
	runVmTests(t, tests)
}
func TestObjectLiterals(t *testing.T) {
	tests := []vmTestCase{
		{`{}`, map[string]any{}},
		{`{"a": 1}`, map[string]any{"a": 1.0}},
		{`{"a": 1, "b": 2}`, map[string]any{"a": 1.0, "b": 2.0}},
		{`{"x": true, "y": false, "z": "hello"}`, map[string]any{"x": true, "y": false, "z": "hello"}},
		{`{"num": 42, "arr": [1,2,3], "obj": {"nested": "yes"}}`, map[string]any{
			"num": 42.0,
			"arr": []any{1.0, 2.0, 3.0},
			"obj": map[string]any{"nested": "yes"},
		}},
	}
	runVmTests(t, tests)
}

func TestObjectIndexing(t *testing.T) {
	tests := []vmTestCase{
		{`{"a": 1, "b": 2}["a"]`, 1.0},
		{`{"a": 1, "b": 2}["b"]`, 2.0},
		{`{"x": true, "y": false}["y"]`, false},
		{`{"foo": "bar"}["foo"]`, "bar"},
		{`{"arr": [1,2,3]}["arr"][1]`, 2.0},
		{`{"obj": {"nested": 99}}["obj"]["nested"]`, 99.0},
	}
	runVmTests(t, tests)
}

func TestPipeFunction(t *testing.T) {
	tests := []vmTestCase{
		{`"foo" as $foo |pipe: $foo + "bar"`, "foobar"},
		{"[1,2] |map: $item * 2", []any{2.0, 4.0}},
		{"[1,2] |map: $item * $index", []any{0.0, 2.0}},
		{"[1,2] |map: $item * 2 |map: $item + 1", []any{3.0, 5.0}},
		{"[1,2,3,4,5,6] |filter: $item > 2", []any{3.0, 4.0, 5.0, 6.0}},

		// Reduce: sum all items
		{"[1,2,3,4] |reduce: ($acc || 0) + $item", 10.0},

		{"[1,2,3,4] |reduce: set($acc || {}, $index, $item)", map[string]any{
			"0": 1.0,
			"1": 2.0,
			"2": 3.0,
			"3": 4.0,
		}},

		// set with string key
		{`set({}, "a", 1)`, map[string]any{"a": 1.0}},
		// set with numeric key coerced to string
		{`set({}, 5, "x")`, map[string]any{"5": "x"}},

		// Find: first item greater than 2
		{"[1,2,3,4] |find: $item > 2", 3.0},

		// Some: any item is even
		{"[1,2,3,4] |some: $item % 2 == 0", true},
		{"[1,3,5] |some: $item % 2 == 0", false},

		// Every: all items are positive
		{"[1,2,3,4] |every: $item > 0", true},
		{"[1,2,3,0] |every: $item > 0", false},

		// Unique: remove duplicates
		{"[1,2,2,3,1,4] |unique: $item", []any{1.0, 2.0, 3.0, 4.0}},

		// Sort: sort by value
		{"[3,1,2] |sort: $item", []any{1.0, 2.0, 3.0}},
		// Sort: sort by computed value
		{"[3,1,2] |sort: $item * -1", []any{3.0, 2.0, 1.0}},

		// GroupBy: group by even/odd
		// {`[1,2,3,4] |groupBy: $item % 2`, map[string]any{
		// 	"1": []any{1, 3},
		// 	"0": []any{2, 4},
		// }},

		// // Window: window size 2, sum each window
		{"[1,2,3,4] |window: $window[0] + $window[1]", []any{3.0, 5.0, 7.0}},

		// Chunk: chunk size 2, sum each chunk
		// {"[1,2,3,4,5] |chunk: $chunk[0] + ($chunk[1] ?? 0)", []any{3, 7, 5}},
	}
	runVmTests(t, tests)
}

func TestNullishOperator(t *testing.T) {
	tests := []vmTestCase{
		// Simple member with existing key whose value is null
		{`{"name": null}?.name`, nil},

		// Nested existing keys
		{`{"user": {"name": "alice"}}?.user?.name`, "alice"},
		{`{"user": {"name": null}}?.user?.name`, nil},

		// Intermediate nil stops chain
		{`{"user": null}?.user?.name`, nil},
		{`{"user": null}?.user?.profile?.age`, nil},

		// Deeper nesting with existing path
		{`{"user": {"profile": {"age": 30}}}?.user?.profile?.age`, 30.0},
		{`{"user": {"profile": null}}?.user?.profile?.age`, nil},

		// Null root
		{`null?.name`, nil},
		{`null?[0]`, nil},
		{`null?.foo?.bar`, nil},
		{`null?[10]?.foo`, nil},

		// Chained null propagation
		{`{"a": null}?.a?.b`, nil},
		{`{"a": {"b": null}}?.a?.b?.c`, nil},
		{`{"a": {"b": {"c": 42}}}?.a?.b?.c`, 42.0},
		{`{"a": {"b": {"c": null}}}?.a?.b?.c`, nil},

		// Arrays + objects (indexes in‑range, properties exist when accessed)
		{`{"a": [{"b": {"c": 99}}]}?.a?[0]?.b?.c`, 99.0},
		{`{"a": [{"b": null}]}?.a?[0]?.b?.c`, nil},

		// Optional index on array of objects
		{`[{"x": 1}, {"y": 2}]?[0]?.x`, 1.0},
		{`[{"x": null}]?[0]?.x`, nil},
		{`[null]?[0]?.x`, nil},

		// Optional index on string (valid indices only)
		{`"hello"?.[0]`, "h"},
		{`"hello"?.[4]`, "o"},
	}

	runVmTests(t, tests)
}

func TestNullishCoalescing(t *testing.T) {
	tests := []vmTestCase{
		// Basic fallback
		{`null ?? 42`, 42.0},
		// {`undefinedVar ?? "default"`, "default"}, // undefined identifier should resolve to null TODO: TEST THIS LATER

		// Left side is not nullish
		{`0 ?? 99`, 0.0},
		{`false ?? true`, false},
		{`"" ?? "fallback"`, ""},

		// Chained ?? operators
		{`null ?? null ?? "x"`, "x"},
		{`null ?? 0 ?? "y"`, 0.0},
		{`null ?? false ?? "z"`, false},
		{`null ?? null ?? null ?? "last"`, "last"},
		{`1 ?? 2 ?? 3`, 1.0},

		// Array out-of-bounds and missing keys (safe mode)
		{`[1,2,3][10] ?? 99`, 99.0},
		{`{"a": 1}["b"] ?? "missing"`, "missing"},
		{`{"a": null}["a"] ?? "fallback"`, "fallback"},

		// Nested property/index with ?? fallback
		{`{"user": {"name": null}}.user.name ?? "anon"`, "anon"},
		{`{"user": {}}.user.name ?? "anon"`, "anon"},
		{`{"user": {"name": "alice"}}.user.name ?? "anon"`, "alice"},

		// Optional chaining + ??
		{`{"user": null}?.user?.name ?? "anon"`, "anon"},
		{`{"user": {"name": null}}?.user?.name ?? "anon"`, "anon"},
		{`{"user": {"name": "bob"}}?.user?.name ?? "anon"`, "bob"},

		// Right side should not be evaluated if left is not nullish
		{`1 ?? (2/0)`, 1.0}, // Should not error or panic

		// Chained with other operators
		{`(null ?? 5) + 2`, 7.0},
		{`(null ?? 0) + 2`, 2.0},
		{`(null ?? "") + "x"`, "x"},
		{`(null ?? "foo") + "bar"`, "foobar"},
	}
	runVmTests(t, tests)
}

func TestTernaryOperator(t *testing.T) {
	tests := []vmTestCase{
		{`true ? 1 : 2`, 1.0},
		{`false ? 1 : 2`, 2.0},
		{`1 ? 10 : 20`, 10.0},
		{`0 ? 10 : 20`, 20.0},
		{`null ? 1 : 2`, 2.0},

		// Numeric / boolean expressions as condition
		{`(1 + 1 == 2) ? 5 * 2 : 3 + 4`, 10.0},
		{`(1 + 1 == 3) ? 5 * 2 : 3 + 4`, 7.0},
		{`(2 * 3 > 5) ? 100 : 200`, 100.0},
		{`(2 * 3 < 5) ? 100 : 200`, 200.0},

		// Nested ternaries
		{`true ? (false ? 1 : 2) : 3`, 2.0},
		{`false ? 1 : (true ? 2 : 3)`, 2.0},
		{`false ? 1 : false ? 2 : 3`, 3.0}, // right‑associative: false ? 1 : (false ? 2 : 3)

		// Mixed with other operators (parenthesized for clarity)
		{`(true && false) ? 1 : 2`, 2.0},
		{`(true || false) ? 1 : 2`, 1.0},
		{`(null ?? 5) ? 1 : 2`, 1.0},
		{`(0 || 5) ? 1 : 2`, 1.0},
		{`(0 && 5) ? 1 : 2`, 2.0},

		// Consequent / alternate as complex expressions
		{`true ? {"a": 1} : {"b": 2}`, map[string]any{"a": 1.0}},
		{`false ? {"a": 1} : {"b": 2}`, map[string]any{"b": 2.0}},
		{`true ? [1,2] : [3,4]`, []any{1.0, 2.0}},
		{`false ? [1,2] : [3,4]`, []any{3.0, 4.0}},
		{`true ? [1, 1+1, 3*2] : 0`, []any{1.0, 2.0, 6.0}},
		{`false ? 0 : [1, 1+1, 3*2]`, []any{1.0, 2.0, 6.0}},

		// Nested inside other expressions
		{`(true ? 5 : 10) + 3`, 8.0},
		{`(false ? 5 : 10) + 3`, 13.0},
		{`(false ? 5 : 10) * (true ? 2 : 4)`, 20.0},

		// Ternary inside array / object
		{`[true ? 1 : 2, false ? 3 : 4].1`, 4.0},
		{`{"x": true ? 1 : 2, "y": false ? 3 : 4}["y"]`, 4.0},

		// Chaining with nullish and logical
		{`(null ?? 0) ? 1 : 2`, 2.0},
		{`(null ?? null) ? 1 : 2`, 2.0},
		{`(false || null) ? 1 : 2`, 2.0},
		{`(false || 7) ? 1 : 2`, 1.0},

		// Deep nesting
		{`true ? (1 ? (0 ? 10 : 20) : 30) : 40`, 20.0},
		{`false ? (1 ? (0 ? 10 : 20) : 30) : (0 ? 50 : 60)`, 60.0},
	}
	runVmTests(t, tests)
}

func TestContextValues(t *testing.T) {
	tests := []vmTestCase{
		// Accessing top-level context values
		{"value", "test"},
		{"number", 42.0},
		{"boolean", true},
		{"nullValue", nil},
		{"array[1]", 2.0},
		{"object.key1", "value1"},

		// Using context values in expressions
		{"array[1] + 1", 3.0},
		{"number * 2", 84.0},
		{"boolean && false", false},
		{"nullValue ?? 'default'", "default"},
		{"array[0] + array[1]", 3.0},
		{"object.key1 + ' is the key'", "value1 is the key"},
	}

	contextValues := map[string]any{
		"value":     "test",
		"number":    42.0,
		"boolean":   true,
		"nullValue": nil,
		"array":     []any{1.0, 2.0, 3.0},
		"object": map[string]any{
			"key1": "value1",
			"key2": "value2",
		},
	}

	runVmTests(t, tests, contextValues)
}
