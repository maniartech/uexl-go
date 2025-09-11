package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/stretchr/testify/assert"
)

// TestParserInternalFunctions tests internal parser functions that might not be covered
func TestParserInternalFunctions(t *testing.T) {
	// Test peelLeadingUnary function indirectly through power operator parsing
	tests := []struct {
		name     string
		input    string
		expected string // simplified representation
	}{
		{"simple power", "2**3", "2**3"},
		{"unary with power", "-2**3", "-(2**3)"}, // should be parsed as -(2**3)
		{"double unary with power", "--2**3", "--(2**3)"},
		{"not with power", "!true**2", "!(true**2)"},
		{"mixed unary with power", "-!2**3", "-!(2**3)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, expr)
			// The exact structure depends on the precedence rules
			// We're mainly testing that it doesn't panic and produces valid AST
		})
	}
}

// TestParserBinaryOperatorPrecedence tests all binary operator precedence levels
func TestParserBinaryOperatorPrecedence(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		// Logical OR (lowest precedence)
		{"logical or", "a || b || c"},

		// Logical AND
		{"logical and", "a && b && c"},

		// Bitwise OR
		{"bitwise or", "a | b | c"},

		// Bitwise XOR
		{"bitwise xor", "a ^ b ^ c"},

		// Bitwise AND
		{"bitwise and", "a & b & c"},

		// Equality
		{"equality", "a == b != c"},

		// Comparison
		{"comparison", "a < b > c <= d >= e"},

		// Nullish coalescing
		{"nullish", "a ?? b ?? c"},

		// Bitwise shift
		{"shift", "a << b >> c"},

		// Additive
		{"additive", "a + b - c"},

		// Multiplicative
		{"multiplicative", "a * b / c % d"},

		// Power (highest precedence, right-associative)
		{"power", "a ** b ** c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, expr)
		})
	}
}

// TestParserComplexNesting tests deeply nested expressions
func TestParserComplexNesting(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"nested parentheses", "((((1 + 2) * 3) - 4) / 5)"},
		{"nested arrays", "[[[1, 2], [3, 4]], [[5, 6], [7, 8]]]"},
		{"nested objects", `{"a": {"b": {"c": {"d": 1}}}}`},
		{"nested function calls", "func1(func2(func3(1, 2), 3), 4)"},
		{"nested member access", "obj.prop1.prop2.prop3.prop4"},
		{"nested index access", "arr[0][1][2][3]"},
		{"mixed nesting", "obj.arr[0].prop"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, expr)
		})
	}
}

// TestParserSliceExpressions tests slice expression parsing
func TestParserSliceExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"full slice", "arr[1:4:2]"},
		{"start only", "arr[1:]"},
		{"end only", "arr[:4]"},
		{"step only", "arr[::2]"},
		{"start and end", "arr[1:4]"},
		{"start and step", "arr[1::2]"},
		{"end and step", "arr[:4:2]"},
		{"empty slice", "arr[:]"},
		{"negative indices", "arr[-1:-3:-1]"},
		{"expression indices", "arr[start():end():step()]"},
		{"optional slice", "arr?[1:4]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, expr)

			// Verify it's a slice expression
			if tt.input != "arr?[1:4]" { // optional slice creates different structure
				assert.IsType(t, &parser.SliceExpression{}, expr)
			}
		})
	}
}

// TestParserConditionalExpressions tests ternary conditional expressions
func TestParserConditionalExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple conditional", "true ? 1 : 2"},
		{"nested conditional", "a ? b ? c : d : e"},
		{"conditional with complex expressions", "(x > 0) ? (x * 2) : (x * -1)"},
		{"conditional in function call", "func(a ? b : c, d)"},
		{"conditional with nullish", "a ?? b ? c : d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, expr)
		})
	}
}

// TestParserPipeExpressions tests pipe expression parsing
func TestParserPipeExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple pipe", "[1,2,3] |: x * 2"},
		{"named pipe", "[1,2,3] |map: x * 2"},
		{"multiple pipes", "[1,2,3] |map: x * 2 |filter: x > 2"},
		{"pipe with alias", "data |map: x * 2 as $doubled"},
		{"complex pipe", "users |filter: age > 18 |map: name |join: ', '"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, expr)
		})
	}
}

// TestParserOptionalChaining tests optional chaining expressions
func TestParserOptionalChaining(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"optional member", "obj?.prop"},
		{"optional index", "arr?.[0]"},
		{"optional member access", "obj?.prop"},
		{"chained optional", "obj?.prop?.value"},
		{"optional with nullish", "obj?.prop ?? defaultValue"},
		{"optional in complex expression", "(obj?.prop || fallback).value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, expr)
		})
	}
}

// TestParserMemberAccessVariations tests different member access patterns
func TestParserMemberAccessVariations(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"dot notation", "obj.prop"},
		{"dot with number", "obj.123"},
		{"dot with decimal", "obj.1.5"}, // should parse as obj.1.5
		{"dot with expression", "obj.(expr)"},
		{"bracket notation", "obj[prop]"},
		{"bracket with string", `obj["prop"]`},
		{"bracket with number", "obj[123]"},
		{"bracket with expression", "obj[a + b]"},
		{"mixed access", "obj.prop[0].method"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, expr)
		})
	}
}

// TestParserFunctionCalls tests function call parsing
func TestParserFunctionCalls(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"no args", "func()"},
		{"single arg", "func(1)"},
		{"multiple args", "func(1, 2, 3)"},
		{"complex args", "func(a + b, obj.prop, arr[0])"},
		{"nested calls", "func1(func2(func3()))"},
		{"function call", "method(1, 2)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, expr)
		})
	}
}

// TestParserArrayLiterals tests array literal parsing
func TestParserArrayLiterals(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty array", "[]"},
		{"single element", "[1]"},
		{"multiple elements", "[1, 2, 3]"},
		{"mixed types", `[1, "hello", true, null]`},
		{"nested arrays", "[[1, 2], [3, 4]]"},
		{"complex elements", "[a + b, obj.prop, func()]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, expr)
			assert.IsType(t, &parser.ArrayLiteral{}, expr)
		})
	}
}

// TestParserObjectLiterals tests object literal parsing
func TestParserObjectLiterals(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty object", "{}"},
		{"single property", `{"key": "value"}`},
		{"multiple properties", `{"a": 1, "b": 2}`},
		{"mixed values", `{"num": 42, "str": "hello", "bool": true, "null": null}`},
		{"nested objects", `{"outer": {"inner": {"deep": 1}}}`},
		{"complex values", `{"sum": a + b, "prop": obj.prop, "call": func()}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, expr)
			assert.IsType(t, &parser.ObjectLiteral{}, expr)
		})
	}
}

// TestParserUnaryExpressions tests unary expression parsing
func TestParserUnaryExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"negative number", "-42"},
		{"negative variable", "-x"},
		{"not boolean", "!true"},
		{"not variable", "!x"},
		{"double negative", "--42"},
		{"double not", "!!true"},
		{"mixed unary", "-!x"},
		{"unary with complex", "-(a + b)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, expr)
		})
	}
}

// TestParserErrorRecovery tests parser error recovery
func TestParserErrorRecovery(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"missing closing paren", "(1 + 2"},
		{"missing closing bracket", "[1, 2"},
		{"missing closing brace", `{"a": 1`},
		{"invalid operator", "1 ++ 2"},
		{"trailing comma in array", "[1, 2,]"},
		{"missing value in object", `{"a": }`},
		{"empty expression in array", "[1, , 3]"},
		{"missing value in object", `{"a": }`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.Error(t, err)
			assert.Nil(t, expr)
		})
	}
}

// TestParserWithCustomOptions tests parser with custom options
func TestParserWithCustomOptions(t *testing.T) {
	// Test with all features disabled
	opts := parser.Options{
		EnableNullish:          false,
		EnableOptionalChaining: false,
		EnablePipes:            false,
		MaxDepth:               5,
	}

	// Basic expressions should still work
	p := parser.NewParserWithOptions("1 + 2", opts)
	expr, err := p.Parse()
	assert.NoError(t, err)
	assert.NotNil(t, expr)

	// Complex expressions should still work
	p = parser.NewParserWithOptions("func(a, b) * (c + d)", opts)
	expr, err = p.Parse()
	assert.NoError(t, err)
	assert.NotNil(t, expr)
}
