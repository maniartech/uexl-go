package parser_test

import (
	"testing"

	"github.com/maniartech/uexl/parser"
	"github.com/maniartech/uexl/parser/constants"
	"github.com/stretchr/testify/assert"
)

// TestParseStringFunction tests the ParseString function from init.go
func TestParseStringFunction(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "simple expression",
			input:       "1 + 2",
			expectError: false,
		},
		{
			name:        "identifier",
			input:       "x",
			expectError: false,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "invalid syntax",
			input:       "1 + +",
			expectError: true,
		},
		{
			name:        "complex expression",
			input:       "func(a, b) * (c + d)",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parser.ParseString(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, node)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, node)
			}
		})
	}
}

// TestNewParserWithValidationFunction tests the NewParserWithValidation function
func TestNewParserWithValidationFunction(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "valid input",
			input:       "1 + 2",
			expectError: false,
		},
		{
			name:        "empty input",
			input:       "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			input:       "   ",
			expectError: false, // whitespace is valid, just results in no meaningful tokens
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := parser.NewParserWithValidation(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, p)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, p)
			}
		})
	}
}

// TestDefaultOptionsFunction tests the DefaultOptions function
func TestDefaultOptionsFunction(t *testing.T) {
	opts := parser.DefaultOptions()
	assert.True(t, opts.EnableNullish)
	assert.True(t, opts.EnableOptionalChaining)
	assert.True(t, opts.EnablePipes)
	assert.Equal(t, 0, opts.MaxDepth)
}

// TestPropertyHelperFunctions tests the Property helper functions
func TestPropertyHelperFunctions(t *testing.T) {
	// Test PropS
	strProp := parser.PropS("test")
	assert.True(t, strProp.IsString())
	assert.False(t, strProp.IsInt())
	assert.Equal(t, "test", strProp.S)

	// Test PropI
	intProp := parser.PropI(42)
	assert.False(t, intProp.IsString())
	assert.True(t, intProp.IsInt())
	assert.Equal(t, 42, intProp.I)
}

// TestTokenValueHelperMethods tests the TokenValue helper methods
func TestTokenValueHelperMethods(t *testing.T) {
	// Test AsFloat
	numToken := parser.Token{
		Type:  constants.TokenNumber,
		Value: parser.TokenValue{Kind: parser.TVKNumber, Num: 3.14},
	}
	val, ok := numToken.AsFloat()
	assert.True(t, ok)
	assert.Equal(t, 3.14, val)

	// Test AsString with string token
	strToken := parser.Token{
		Type:  constants.TokenString,
		Value: parser.TokenValue{Kind: parser.TVKString, Str: "hello"},
	}
	str, ok := strToken.AsString()
	assert.True(t, ok)
	assert.Equal(t, "hello", str)

	// Test AsString with identifier token
	idToken := parser.Token{
		Type:  constants.TokenIdentifier,
		Value: parser.TokenValue{Kind: parser.TVKIdentifier, Str: "variable"},
	}
	str, ok = idToken.AsString()
	assert.True(t, ok)
	assert.Equal(t, "variable", str)

	// Test AsString with operator token
	opToken := parser.Token{
		Type:  constants.TokenOperator,
		Value: parser.TokenValue{Kind: parser.TVKOperator, Str: "+"},
	}
	str, ok = opToken.AsString()
	assert.True(t, ok)
	assert.Equal(t, "+", str)

	// Test AsBool
	boolToken := parser.Token{
		Type:  constants.TokenBoolean,
		Value: parser.TokenValue{Kind: parser.TVKBoolean, Bool: true},
	}
	b, ok := boolToken.AsBool()
	assert.True(t, ok)
	assert.True(t, b)

	// Test failed conversions
	_, ok = strToken.AsFloat()
	assert.False(t, ok)

	_, ok = numToken.AsString()
	assert.False(t, ok)

	_, ok = strToken.AsBool()
	assert.False(t, ok)
}

// TestAllNodeTypesAndPositions tests that all AST nodes implement the Node interface correctly
func TestAllNodeTypesAndPositions(t *testing.T) {
	tests := []struct {
		name         string
		node         parser.Expression
		expectedType parser.NodeType
		expectedLine int
		expectedCol  int
	}{
		{
			name:         "NumberLiteral",
			node:         &parser.NumberLiteral{Value: 42, Line: 1, Column: 5},
			expectedType: parser.NodeTypeNumberLiteral,
			expectedLine: 1,
			expectedCol:  5,
		},
		{
			name:         "StringLiteral",
			node:         &parser.StringLiteral{Value: "test", Line: 2, Column: 10},
			expectedType: parser.NodeTypeStringLiteral,
			expectedLine: 2,
			expectedCol:  10,
		},
		{
			name:         "BooleanLiteral",
			node:         &parser.BooleanLiteral{Value: true, Line: 3, Column: 15},
			expectedType: parser.NodeTypeBooleanLiteral,
			expectedLine: 3,
			expectedCol:  15,
		},
		{
			name:         "NullLiteral",
			node:         &parser.NullLiteral{Line: 4, Column: 20},
			expectedType: parser.NodeTypeNullLiteral,
			expectedLine: 4,
			expectedCol:  20,
		},
		{
			name:         "Identifier",
			node:         &parser.Identifier{Name: "x", Line: 5, Column: 25},
			expectedType: parser.NodeTypeIdentifier,
			expectedLine: 5,
			expectedCol:  25,
		},
		{
			name: "BinaryExpression",
			node: &parser.BinaryExpression{
				Left:     &parser.NumberLiteral{Value: 1},
				Operator: "+",
				Right:    &parser.NumberLiteral{Value: 2},
				Line:     6,
				Column:   30,
			},
			expectedType: parser.NodeTypeBinaryExpression,
			expectedLine: 6,
			expectedCol:  30,
		},
		{
			name: "UnaryExpression",
			node: &parser.UnaryExpression{
				Operator: "-",
				Operand:  &parser.NumberLiteral{Value: 5},
				Line:     7,
				Column:   35,
			},
			expectedType: parser.NodeTypeUnaryExpression,
			expectedLine: 7,
			expectedCol:  35,
		},
		{
			name: "ConditionalExpression",
			node: &parser.ConditionalExpression{
				Condition:  &parser.BooleanLiteral{Value: true},
				Consequent: &parser.NumberLiteral{Value: 1},
				Alternate:  &parser.NumberLiteral{Value: 2},
				Line:       8,
				Column:     40,
			},
			expectedType: parser.NodeTypeConditional,
			expectedLine: 8,
			expectedCol:  40,
		},
		{
			name: "ArrayLiteral",
			node: &parser.ArrayLiteral{
				Elements: []parser.Expression{&parser.NumberLiteral{Value: 1}},
				Line:     9,
				Column:   45,
			},
			expectedType: parser.NodeTypeArrayLiteral,
			expectedLine: 9,
			expectedCol:  45,
		},
		{
			name: "ObjectLiteral",
			node: &parser.ObjectLiteral{
				Properties: map[string]parser.Expression{"key": &parser.StringLiteral{Value: "value"}},
				Line:       10,
				Column:     50,
			},
			expectedType: parser.NodeTypeObjectLiteral,
			expectedLine: 10,
			expectedCol:  50,
		},
		{
			name: "FunctionCall",
			node: &parser.FunctionCall{
				Function:  &parser.Identifier{Name: "func"},
				Arguments: []parser.Expression{&parser.NumberLiteral{Value: 1}},
				Line:      11,
				Column:    55,
			},
			expectedType: parser.NodeTypeFunctionCall,
			expectedLine: 11,
			expectedCol:  55,
		},
		{
			name: "MemberAccess",
			node: &parser.MemberAccess{
				Target:   &parser.Identifier{Name: "obj"},
				Property: parser.PropS("prop"),
				Line:     12,
				Column:   60,
			},
			expectedType: parser.NodeTypeMemberAccess,
			expectedLine: 12,
			expectedCol:  60,
		},
		{
			name: "IndexAccess",
			node: &parser.IndexAccess{
				Target: &parser.Identifier{Name: "arr"},
				Index:  &parser.NumberLiteral{Value: 0},
				Line:   13,
				Column: 65,
			},
			expectedType: parser.NodeType("IndexAccess"),
			expectedLine: 13,
			expectedCol:  65,
		},
		{
			name: "SliceExpression",
			node: &parser.SliceExpression{
				Target: &parser.Identifier{Name: "arr"},
				Start:  &parser.NumberLiteral{Value: 1},
				End:    &parser.NumberLiteral{Value: 3},
				Line:   14,
				Column: 70,
			},
			expectedType: parser.NodeTypeSliceExpression,
			expectedLine: 14,
			expectedCol:  70,
		},
		{
			name: "PipeExpression",
			node: &parser.PipeExpression{
				Expression: &parser.NumberLiteral{Value: 1},
				PipeType:   "map",
				Line:       15,
				Column:     75,
			},
			expectedType: parser.NodeTypePipeExpression,
			expectedLine: 15,
			expectedCol:  75,
		},
		{
			name: "ProgramNode",
			node: &parser.ProgramNode{
				PipeExpressions: []parser.PipeExpression{
					{Expression: &parser.NumberLiteral{Value: 1}},
				},
				Line:   16,
				Column: 80,
			},
			expectedType: parser.NodeTypeProgram,
			expectedLine: 16,
			expectedCol:  80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Type() method
			assert.Equal(t, tt.expectedType, tt.node.Type())

			// Test Position() method
			line, col := tt.node.Position()
			assert.Equal(t, tt.expectedLine, line)
			assert.Equal(t, tt.expectedCol, col)

			// Ensure the node implements the Expression interface
			// (expressionNode is unexported, so we just verify it's an Expression)
			var _ parser.Expression = tt.node
		})
	}
}

// TestTokenizerCreationAndBasicFunctionality tests tokenizer creation and basic methods
func TestTokenizerCreationAndBasicFunctionality(t *testing.T) {
	tokenizer := parser.NewTokenizer("hello 123 + world")

	// Test that tokenizer is created properly
	assert.NotNil(t, tokenizer)

	// Test NextToken functionality
	token, err := tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenIdentifier, token.Type)
	assert.Equal(t, "hello", token.Value.Str)

	// Test PreloadTokens
	tokenizer2 := parser.NewTokenizer("1 + 2")
	tokens := tokenizer2.PreloadTokens()
	assert.Len(t, tokens, 4) // number, operator, number, EOF

	expectedTypes := []constants.TokenType{
		constants.TokenNumber,
		constants.TokenOperator,
		constants.TokenNumber,
		constants.TokenEOF,
	}

	for i, expectedType := range expectedTypes {
		assert.Equal(t, expectedType, tokens[i].Type)
	}
}

// TestTokenStringMethod tests the Token String() method
func TestTokenStringMethod(t *testing.T) {
	token := parser.Token{
		Type:   constants.TokenNumber,
		Value:  parser.TokenValue{Kind: parser.TVKNumber, Num: 42},
		Token:  "42",
		Line:   1,
		Column: 5,
	}

	str := token.String()
	assert.Contains(t, str, "Number")
	assert.Contains(t, str, "42")
	assert.Contains(t, str, "1:5")
}

// TestParserConstructors tests all parser constructor functions
func TestParserConstructors(t *testing.T) {
	// Test NewParser
	p1 := parser.NewParser("1 + 2")
	assert.NotNil(t, p1)

	// Test NewParserWithOptions
	opts := parser.Options{
		EnableNullish:          false,
		EnableOptionalChaining: false,
		EnablePipes:            false,
		MaxDepth:               10,
	}
	p2 := parser.NewParserWithOptions("1 + 2", opts)
	assert.NotNil(t, p2)

	// Test that both can parse successfully
	expr1, err1 := p1.Parse()
	assert.NoError(t, err1)
	assert.NotNil(t, expr1)

	expr2, err2 := p2.Parse()
	assert.NoError(t, err2)
	assert.NotNil(t, expr2)
}

// TestParserComplexExpressionParsing tests complex expression parsing
func TestParserComplexExpressionParsing(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"power with unary", "-2**3"},
		{"nested conditionals", "a ? b ? c : d : e"},
		{"complex member access", "obj.prop1.prop2[0]"},
		{"mixed operators", "a + b * c - d / e % f"},
		{"optional chaining", "obj?.prop?.value"},
		{"nullish with conditional", "a ?? b ? c : d"},
		{"pipe with complex expressions", "[1,2,3] |map: x * 2 + 1 |filter: x > 5"},
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
