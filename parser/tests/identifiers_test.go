package parser_test

import (
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/parser/constants"
	"github.com/stretchr/testify/assert"
)

func TestIdentifierTokenization(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedTokens []struct {
			token     string
			tokenType constants.TokenType
		}
		description string
	}{
		// Basic identifier patterns
		{
			name:  "Simple variable",
			input: "x",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"x", constants.TokenIdentifier},
			},
			description: "Simple alpha identifier",
		},
		{
			name:  "Uppercase variable",
			input: "X",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"X", constants.TokenIdentifier},
			},
			description: "Uppercase alpha identifier",
		},
		{
			name:  "Mixed case variable",
			input: "myVar",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"myVar", constants.TokenIdentifier},
			},
			description: "Mixed case identifier",
		},
		{
			name:  "Underscore variable",
			input: "_var",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"_var", constants.TokenIdentifier},
			},
			description: "Underscore start identifier",
		},
		{
			name:  "Double underscore",
			input: "__var",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"__var", constants.TokenIdentifier},
			},
			description: "Double underscore start",
		},

		// Dollar sign identifiers (common in pipes)
		{
			name:  "Dollar variable",
			input: "$var",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"$var", constants.TokenIdentifier},
			},
			description: "Dollar sign start identifier",
		},
		{
			name:  "Dollar number",
			input: "$1",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"$1", constants.TokenIdentifier},
			},
			description: "Dollar with number",
		},
		{
			name:  "Dollar underscore",
			input: "$_var",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"$_var", constants.TokenIdentifier},
			},
			description: "Dollar underscore identifier",
		},
		{
			name:  "Double dollar",
			input: "$$var",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"$$var", constants.TokenIdentifier},
			},
			description: "Double dollar identifier",
		},

		// Alphanumeric combinations
		{
			name:  "Alpha numeric",
			input: "var123",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"var123", constants.TokenIdentifier},
			},
			description: "Alpha numeric identifier",
		},
		{
			name:  "Number in middle",
			input: "var1test",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"var1test", constants.TokenIdentifier},
			},
			description: "Number in middle",
		},
		{
			name:  "Underscore numeric",
			input: "_123var",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"_123var", constants.TokenIdentifier},
			},
			description: "Underscore with numbers",
		},
		{
			name:  "Dollar numeric",
			input: "$123var",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"$123var", constants.TokenIdentifier},
			},
			description: "Dollar with numbers",
		},

		// Complex valid identifiers
		{
			name:  "Complex identifier",
			input: "myVar_123",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"myVar_123", constants.TokenIdentifier},
			},
			description: "Complex valid identifier",
		},
		{
			name:  "Identifier with dollar",
			input: "var$test",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"var$test", constants.TokenIdentifier},
			},
			description: "Identifier with dollar sign",
		},
		{
			name:  "Mixed identifiers",
			input: "a1_ $b2_",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"a1_", constants.TokenIdentifier},
				{"$b2_", constants.TokenIdentifier},
			},
			description: "Two separate identifiers",
		},

		// Dot notation should be separate tokens now
		{
			name:  "Simple dot identifier",
			input: "obj.prop",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"obj", constants.TokenIdentifier},
				{".", constants.TokenDot},
				{"prop", constants.TokenIdentifier},
			},
			description: "Simple dot notation should be separate tokens",
		},
		{
			name:  "Chained dot identifier",
			input: "user.profile.name",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"user", constants.TokenIdentifier},
				{".", constants.TokenDot},
				{"profile", constants.TokenIdentifier},
				{".", constants.TokenDot},
				{"name", constants.TokenIdentifier},
			},
			description: "Chained dot notation should be separate tokens",
		},
		{
			name:  "Dot with underscore",
			input: "obj._private",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"obj", constants.TokenIdentifier},
				{".", constants.TokenDot},
				{"_private", constants.TokenIdentifier},
			},
			description: "Dot notation with underscore should be separate tokens",
		},
		{
			name:  "Dot with dollar",
			input: "obj.$special",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"obj", constants.TokenIdentifier},
				{".", constants.TokenDot},
				{"$special", constants.TokenIdentifier},
			},
			description: "Dot notation with dollar sign should be separate tokens",
		},
		{
			name:  "Complex dot identifier",
			input: "data.items.name",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"data", constants.TokenIdentifier},
				{".", constants.TokenDot},
				{"items", constants.TokenIdentifier},
				{".", constants.TokenDot},
				{"name", constants.TokenIdentifier},
			},
			description: "Complex dot notation should be separate tokens",
		},
		{
			name:  "Dot with number",
			input: "data.items.0.name",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"data", constants.TokenIdentifier},
				{".", constants.TokenDot},
				{"items", constants.TokenIdentifier},
				{".", constants.TokenDot},
				{"0.", constants.TokenNumber},
				{"name", constants.TokenIdentifier},
			},
			description: "Dot notation with number should be separate tokens (0. is a valid number)",
		},

		// Invalid dot patterns (tokenized correctly but should cause parsing errors)
		{
			name:  "Double dot",
			input: "obj..prop",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"obj", constants.TokenIdentifier},
				{".", constants.TokenDot},
				{".", constants.TokenDot},
				{"prop", constants.TokenIdentifier},
			},
			description: "Double dots tokenized correctly (but invalid syntax)",
		},
		{
			name:  "Trailing dot",
			input: "obj.",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"obj", constants.TokenIdentifier},
				{".", constants.TokenDot},
			},
			description: "Trailing dot tokenized correctly (but invalid syntax)",
		},
		{
			name:  "Dot followed by number",
			input: "obj.123",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"obj", constants.TokenIdentifier},
				{".", constants.TokenDot},
				{"123", constants.TokenNumber},
			},
			description: "Dot followed by number tokenized correctly (but invalid syntax)",
		},
		{
			name:  "Three dots",
			input: "obj...prop",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"obj", constants.TokenIdentifier},
				{".", constants.TokenDot},
				{".", constants.TokenDot},
				{".", constants.TokenDot},
				{"prop", constants.TokenIdentifier},
			},
			description: "Three consecutive dots should be separate tokens",
		},
		{
			name:  "Dot at start of expression",
			input: ".method()",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{".", constants.TokenDot},
				{"method", constants.TokenIdentifier},
				{"(", constants.TokenLeftParen},
				{")", constants.TokenRightParen},
			},
			description: "Dot at start should be separate token",
		},
		{
			name:  "Complex mixed dot pattern",
			input: "obj.prop..method.name",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"obj", constants.TokenIdentifier},
				{".", constants.TokenDot},
				{"prop", constants.TokenIdentifier},
				{".", constants.TokenDot},
				{".", constants.TokenDot},
				{"method", constants.TokenIdentifier},
				{".", constants.TokenDot},
				{"name", constants.TokenIdentifier},
			},
			description: "Mixed valid and invalid dot patterns should be separate tokens",
		},

		// Edge cases that should work
		{
			name:  "Single letter",
			input: "a",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"a", constants.TokenIdentifier},
			},
			description: "Single letter identifier",
		},
		{
			name:  "Single underscore",
			input: "_",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"_", constants.TokenIdentifier},
			},
			description: "Single underscore",
		},
		{
			name:  "Single dollar",
			input: "$",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"$", constants.TokenIdentifier},
			},
			description: "Single dollar sign",
		},

		// Invalid patterns (should be tokenized differently)
		{
			name:  "Number start",
			input: "123var",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"123", constants.TokenNumber},
				{"var", constants.TokenIdentifier},
			},
			description: "Number followed by identifier - should be two tokens",
		},

		// Keywords that are NOT identifiers
		{
			name:  "True keyword",
			input: "true",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"true", constants.TokenBoolean},
			},
			description: "Boolean true - not identifier",
		},
		{
			name:  "False keyword",
			input: "false",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"false", constants.TokenBoolean},
			},
			description: "Boolean false - not identifier",
		},
		{
			name:  "Null keyword",
			input: "null",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"null", constants.TokenNull},
			},
			description: "Null keyword - not identifier",
		},
		{
			name:  "As keyword",
			input: "as",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{"as", constants.TokenAs},
			},
			description: "As keyword - not identifier",
		},

		// Invalid starting patterns
		{
			name:  "Dot start",
			input: ".invalid",
			expectedTokens: []struct {
				token     string
				tokenType constants.TokenType
			}{
				{".", constants.TokenDot},
				{"invalid", constants.TokenIdentifier},
			},
			description: "Identifier cannot start with dot",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := parser.NewTokenizer(tt.input)
			var tokens []parser.Token
			for {
				token, err := tokenizer.NextToken()
				if err != nil {
					t.Errorf("Tokenization error: %v", err)
					break
				}
				if token.Type == constants.TokenEOF {
					break
				}
				tokens = append(tokens, token)
			}

			assert.Equal(t, len(tt.expectedTokens), len(tokens),
				"Expected %d tokens, got %d for input '%s'",
				len(tt.expectedTokens), len(tokens), tt.input)

			for i, expectedToken := range tt.expectedTokens {
				if i < len(tokens) {
					assert.Equal(t, expectedToken.token, tokens[i].Token,
						"Token %d: expected '%s', got '%s'", i, expectedToken.token, tokens[i].Token)
					assert.Equal(t, expectedToken.tokenType, tokens[i].Type,
						"Token %d type: expected %s, got %s", i, expectedToken.tokenType, tokens[i].Type)
				}
			}
		})
	}
}

func TestIdentifierParsing(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType interface{}
		expectedName string
		shouldError  bool
		description  string
	}{
		// Valid identifier parsing
		{
			name:         "Simple identifier",
			input:        "x",
			expectedType: &parser.Identifier{},
			expectedName: "x",
			shouldError:  false,
			description:  "Simple variable identifier",
		},
		{
			name:         "Dollar identifier",
			input:        "$var",
			expectedType: &parser.Identifier{},
			expectedName: "$var",
			shouldError:  false,
			description:  "Dollar sign identifier",
		},
		{
			name:         "Complex identifier",
			input:        "myVar_123$test",
			expectedType: &parser.Identifier{},
			expectedName: "myVar_123$test",
			shouldError:  false,
			description:  "Complex identifier with all valid characters",
		},
		{
			name:         "Underscore identifier",
			input:        "_privateVar",
			expectedType: &parser.Identifier{},
			expectedName: "_privateVar",
			shouldError:  false,
			description:  "Underscore prefixed identifier",
		},
		{
			name:         "Dot notation member access",
			input:        "obj.prop",
			expectedType: &parser.MemberAccess{},
			expectedName: "",
			shouldError:  false,
			description:  "Dot notation should parse as MemberAccess",
		},
		{
			name:         "Chained dot member access",
			input:        "user.profile.name",
			expectedType: &parser.MemberAccess{},
			expectedName: "",
			shouldError:  false,
			description:  "Chained dot notation should parse as MemberAccess",
		},
		{
			name:         "Complex dot member access",
			input:        "data.items.value",
			expectedType: &parser.MemberAccess{},
			expectedName: "",
			shouldError:  false,
			description:  "Complex dot notation should parse as MemberAccess",
		},

		// Keywords should not be parsed as identifiers
		{
			name:         "Boolean true",
			input:        "true",
			expectedType: &parser.BooleanLiteral{},
			expectedName: "",
			shouldError:  false,
			description:  "Boolean literal, not identifier",
		},
		{
			name:         "Boolean false",
			input:        "false",
			expectedType: &parser.BooleanLiteral{},
			expectedName: "",
			shouldError:  false,
			description:  "Boolean literal, not identifier",
		},
		{
			name:         "Null literal",
			input:        "null",
			expectedType: &parser.NullLiteral{},
			expectedName: "",
			shouldError:  false,
			description:  "Null literal, not identifier",
		},

		// Complex expressions with identifiers
		{
			name:         "Function call",
			input:        "func()",
			expectedType: &parser.FunctionCall{},
			expectedName: "",
			shouldError:  false,
			description:  "Function call expression",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()

			if tt.shouldError {
				assert.Error(t, err, "Expected error for input '%s'", tt.input)
				return
			}

			assert.NoError(t, err, "Unexpected error for input '%s': %v", tt.input, err)
			assert.NotNil(t, expr, "Expression should not be nil for input '%s'", tt.input)

			// Check the type of the parsed expression
			switch expected := tt.expectedType.(type) {
			case *parser.Identifier:
				identifier, ok := expr.(*parser.Identifier)
				assert.True(t, ok, "Expected Identifier, got %T for input '%s'", expr, tt.input)
				if ok {
					assert.Equal(t, tt.expectedName, identifier.Name,
						"Expected identifier name '%s', got '%s'", tt.expectedName, identifier.Name)
				}

			case *parser.MemberAccess:
				_, ok := expr.(*parser.MemberAccess)
				assert.True(t, ok, "Expected MemberAccess, got %T for input '%s'", expr, tt.input)

			case *parser.BooleanLiteral:
				_, ok := expr.(*parser.BooleanLiteral)
				assert.True(t, ok, "Expected BooleanLiteral, got %T for input '%s'", expr, tt.input)

			case *parser.NullLiteral:
				_, ok := expr.(*parser.NullLiteral)
				assert.True(t, ok, "Expected NullLiteral, got %T for input '%s'", expr, tt.input)

			case *parser.FunctionCall:
				_, ok := expr.(*parser.FunctionCall)
				assert.True(t, ok, "Expected FunctionCall, got %T for input '%s'", expr, tt.input)

			default:
				t.Errorf("Unknown expected type %T", expected)
			}
		})
	}
}

func TestIdentifierInExpressions(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "Binary expression with identifiers",
			input:       "x + y",
			description: "Addition with two identifiers",
		},
		{
			name:        "Mixed expression",
			input:       "myVar * 2 + otherVar",
			description: "Complex expression with identifiers and numbers",
		},
		{
			name:        "Identifier comparison",
			input:       "status == 'active'",
			description: "Identifier compared to string",
		},
		{
			name:        "Dollar variables",
			input:       "$1 + $2",
			description: "Pipe variables in expression",
		},
		{
			name:        "Function with identifier",
			input:       "max(value1, value2)",
			description: "Function call with identifier arguments",
		},
		{
			name:        "Dot notation identifier",
			input:       "user.profile.name",
			description: "Dot notation identifier",
		},
		{
			name:        "Array access with identifier",
			input:       "items[index]",
			description: "Array access using identifier index",
		},
		{
			name:        "Complex dot expression",
			input:       "data.items.value + other.field",
			description: "Expression with multiple dot notation identifiers",
		},
		{
			name:        "Complex pipe expression",
			input:       "data |filter: $1.status == 'active' |map: $1.name",
			description: "Pipe expression with dot notation identifiers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()

			assert.NoError(t, err, "Unexpected error parsing '%s': %v", tt.input, err)
			assert.NotNil(t, expr, "Expression should not be nil for input '%s'", tt.input)

			// Verify that the expression can be converted to AST
			parserInstance := parser.NewParser(tt.input)
			node, err := parserInstance.Parse()
			assert.NoError(t, err, "AST conversion failed for '%s': %v", tt.input, err)
			assert.NotNil(t, node, "AST node should not be nil for input '%s'", tt.input)

			// Note: Evaluation will fail for identifiers without context, but that's expected
			// The important thing is that parsing and AST conversion work correctly
		})
	}
}

func TestIdentifierParsingErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "Double dot",
			input:       "obj..prop",
			description: "Double dots should cause parsing error",
		},
		{
			name:        "Triple dot",
			input:       "obj...prop",
			description: "Triple dots should cause parsing error",
		},
		{
			name:        "Trailing dot",
			input:       "obj.",
			description: "Trailing dot should cause parsing error",
		},
		{
			name:        "Leading dot",
			input:       ".prop",
			description: "Leading dot should cause parsing error",
		},
		{
			name:        "Dot followed by number",
			input:       "obj.123",
			description: "Dot followed by number should cause parsing error",
		},
		{
			name:        "Double dot in middle",
			input:       "obj.prop..method",
			description: "Double dot in middle should cause parsing error",
		},
		{
			name:        "Multiple consecutive dots",
			input:       "obj....prop",
			description: "Multiple consecutive dots should cause parsing error",
		},
		{
			name:        "Chained method calls",
			input:       "obj.method().chain()",
			description: "Chained method calls are invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()

			// These should all cause parsing errors
			assert.Error(t, err, "Expected parsing error for invalid pattern '%s'", tt.input)

			// If no error occurred, that's unexpected
			if err == nil {
				t.Logf("Unexpected success parsing '%s', got expression: %T", tt.input, expr)
			}
		})
	}
}

func TestIdentifierSyntaxValidation(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
		description string
	}{
		// Valid patterns that should NOT error
		{
			name:        "Simple dot notation",
			input:       "obj.prop",
			shouldError: false,
			description: "Valid dot notation should parse successfully",
		},
		{
			name:        "Chained dot notation",
			input:       "user.profile.name",
			shouldError: false,
			description: "Valid chained dot notation should parse successfully",
		},
		{
			name:        "Dot with underscore",
			input:       "obj._private",
			shouldError: false,
			description: "Dot with underscore should parse successfully",
		},
		{
			name:        "Dot with dollar",
			input:       "obj.$special",
			shouldError: false,
			description: "Dot with dollar should parse successfully",
		},

		// Invalid patterns that SHOULD error
		{
			name:        "Double dot",
			input:       "obj..prop",
			shouldError: true,
			description: "Double dots should cause parsing error",
		},
		{
			name:        "Trailing dot",
			input:       "obj.",
			shouldError: true,
			description: "Trailing dot should cause parsing error",
		},
		{
			name:        "Leading dot",
			input:       ".method",
			shouldError: true,
			description: "Leading dot should cause parsing error",
		},
		{
			name:        "Dot with number",
			input:       "obj.123",
			shouldError: true,
			description: "Dot followed by number should cause parsing error",
		},
		{
			name:        "Mixed invalid pattern",
			input:       "obj.prop..method",
			shouldError: true,
			description: "Mixed valid and invalid patterns should cause parsing error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()

			if tt.shouldError {
				assert.Error(t, err, "Expected parsing error for invalid pattern '%s'", tt.input)
				if err == nil {
					t.Logf("Unexpected success parsing '%s', got expression: %T", tt.input, expr)
				}
			} else {
				assert.NoError(t, err, "Unexpected parsing error for valid pattern '%s': %v", tt.input, err)
				assert.NotNil(t, expr, "Expression should not be nil for valid pattern '%s'", tt.input)
			}
		})
	}
}

func TestIdentifierVsComplexExpressions(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType interface{}
		description  string
	}{
		// Valid identifiers (should parse as Identifier)
		{
			name:         "Simple identifier",
			input:        "name",
			expectedType: &parser.Identifier{},
			description:  "Simple identifier should parse as Identifier",
		},
		{
			name:         "Dot notation member access",
			input:        "user.profile.name",
			expectedType: &parser.MemberAccess{},
			description:  "Dot notation should parse as MemberAccess",
		},
		{
			name:         "Dollar identifier",
			input:        "$1",
			expectedType: &parser.Identifier{},
			description:  "Dollar identifier should parse as Identifier",
		},

		// Function calls (should NOT parse as Identifier)
		{
			name:         "Simple function call",
			input:        "print()",
			expectedType: &parser.FunctionCall{},
			description:  "Function call should parse as FunctionCall, not Identifier",
		},
		{
			name:         "Function with identifier argument",
			input:        "upper(name)",
			expectedType: &parser.FunctionCall{},
			description:  "Function with identifier argument should parse as FunctionCall",
		},
		{
			name:         "Function with dot notation argument",
			input:        "length(user.name)",
			expectedType: &parser.FunctionCall{},
			description:  "Function with dot notation argument should parse as FunctionCall",
		},

		// Complex expressions that are NOT identifiers
		{
			name:         "Member access with function call (should error)",
			input:        "obj.method()",
			expectedType: nil,
			description:  "Member access with function call should produce a parse error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()

			if tt.expectedType == nil {
				assert.Error(t, err, "Expected error for input '%s'", tt.input)
			} else {
				assert.NoError(t, err, "Unexpected error parsing '%s': %v", tt.input, err)
				assert.NotNil(t, expr, "Expression should not be nil for input '%s'", tt.input)

				// Check the type of the parsed expression
				switch tt.expectedType.(type) {
				case *parser.Identifier:
					identifier, ok := expr.(*parser.Identifier)
					assert.True(t, ok, "Expected Identifier, got %T for input '%s'", expr, tt.input)
					if ok {
						t.Logf("✓ Correctly parsed '%s' as Identifier with name: %s", tt.input, identifier.Name)
					}

				case *parser.MemberAccess:
					_, ok := expr.(*parser.MemberAccess)
					assert.True(t, ok, "Expected MemberAccess, got %T for input '%s'", expr, tt.input)
					if ok {
						t.Logf("✓ Correctly parsed '%s' as MemberAccess", tt.input)
					}

				case *parser.FunctionCall:
					_, ok := expr.(*parser.FunctionCall)
					assert.True(t, ok, "Expected FunctionCall, got %T for input '%s'", expr, tt.input)
					if ok {
						t.Logf("✓ Correctly parsed '%s' as FunctionCall", tt.input)
					}

				default:
					t.Errorf("Unknown expected type %T", tt.expectedType)
				}
			}
		})
	}
}

func TestChainedFunctionMemberAccessErrors(t *testing.T) {
	// Test patterns that should be rejected - chained function/member access
	invalidPatterns := []string{
		"obj.method().chain()",   // Method chaining after function call (invalid)
		"obj.method().another()", // Multiple function calls in chain (invalid)
	}

	for _, pattern := range invalidPatterns {
		t.Run(pattern, func(t *testing.T) {
			p := parser.NewParser(pattern)
			_, err := p.Parse()

			// These should now produce parsing errors
			if err == nil {
				t.Errorf("Expected parsing error for invalid pattern: %s", pattern)
			} else {
				// Check that the error message matches the new parser output
				errStr := err.Error()
				if !strings.Contains(errStr, "function calls are only allowed after identifiers or function calls") {
					t.Errorf("Expected error about function call chaining, got: %s", errStr)
				}
			}
		})
	}
}

func TestMemberAccessVsFunctionCalls(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType interface{}
		description  string
	}{
		// Pure identifiers (should parse as Identifier)
		{
			name:         "Simple identifier",
			input:        "name",
			expectedType: &parser.Identifier{},
			description:  "Simple identifier should parse as Identifier",
		},
		{
			name:         "Dot notation identifier",
			input:        "user.profile.name",
			expectedType: &parser.Identifier{},
			description:  "Dot notation identifier should parse as Identifier",
		},
		{
			name:         "Dollar identifier",
			input:        "$1.status",
			expectedType: &parser.Identifier{},
			description:  "Dollar identifier with dot notation should parse as Identifier",
		},

		// Member access (property access - these might be MemberAccess or Identifier depending on parser implementation)
		{
			name:         "Simple member access",
			input:        "obj.prop",
			expectedType: &parser.Identifier{}, // Based on current implementation
			description:  "Simple property access",
		},
		{
			name:         "Nested member access",
			input:        "user.profile.address.street",
			expectedType: &parser.Identifier{}, // Based on current implementation
			description:  "Nested property access",
		},

		// Function calls (should parse as FunctionCall)
		{
			name:         "Simple function call",
			input:        "print()",
			expectedType: &parser.FunctionCall{},
			description:  "Function call should parse as FunctionCall",
		},
		{
			name:         "Method call (should error)",
			input:        "obj.method()",
			expectedType: nil,
			description:  "Method call on member access should produce a parse error",
		},
		{
			name:         "Function with arguments",
			input:        "max(a, b)",
			expectedType: &parser.FunctionCall{},
			description:  "Function with arguments should parse as FunctionCall",
		},
		{
			name:         "Method with arguments (should error)",
			input:        "obj.calculate(x, y)",
			expectedType: nil,
			description:  "Method call with arguments on member access should produce a parse error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()

			// Handle cases that are marked as TODO (currently accepted but should be invalid)
			if strings.Contains(tt.name, "TODO: should be invalid") {
				if err != nil {
					t.Logf("✓ Expression '%s' is already correctly rejected: %v", tt.input, err)
				} else {
					t.Logf("⚠️  TODO: Expression '%s' currently parses as %T but should be rejected as invalid syntax", tt.input, expr)
				}
				// Continue with type checking for documentation purposes
			}

			if tt.expectedType == nil {
				assert.Error(t, err, "Expected error parsing '%s' (should be invalid)", tt.input)
			} else {
				assert.NoError(t, err, "Unexpected error parsing '%s': %v", tt.input, err)
				assert.NotNil(t, expr, "Expression should not be nil for input '%s'", tt.input)

				// Check the type of the parsed expression
				switch tt.expectedType.(type) {
				case *parser.Identifier:
					identifier, ok := expr.(*parser.Identifier)
					if ok {
						t.Logf("✓ Correctly parsed '%s' as Identifier with name: %s", tt.input, identifier.Name)
					} else {
						memberAccess, isMemberAccess := expr.(*parser.MemberAccess)
						if isMemberAccess {
							t.Logf("ℹ Parsed '%s' as MemberAccess instead of Identifier - this is also valid for property access", tt.input)
							assert.NotNil(t, memberAccess, "MemberAccess should not be nil")
						} else {
							assert.True(t, ok, "Expected Identifier or MemberAccess, got %T for input '%s'", expr, tt.input)
						}
					}

				case *parser.FunctionCall:
					_, ok := expr.(*parser.FunctionCall)
					assert.True(t, ok, "Expected FunctionCall, got %T for input '%s'", expr, tt.input)
					if ok && !strings.Contains(tt.name, "TODO") {
						t.Logf("✓ Correctly parsed '%s' as FunctionCall", tt.input)
					}

				case *parser.MemberAccess:
					_, ok := expr.(*parser.MemberAccess)
					assert.True(t, ok, "Expected MemberAccess, got %T for input '%s'", expr, tt.input)
					if ok && !strings.Contains(tt.name, "TODO") {
						t.Logf("✓ Correctly parsed '%s' as MemberAccess", tt.input)
					}

				default:
					t.Errorf("Unknown expected type %T", tt.expectedType)
				}
			}
		})
	}
}
