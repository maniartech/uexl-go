package parser

import (
	"testing"
)

func TestPlayground(t *testing.T) {
	// This test is to be used as a playground for experimenting and debugging
	// failing tests or new features.

	// DEBUG: Testing number parsing in tokenizer and parser
	t.Log("=== DEBUGGING NUMBER PARSING ===")

	numberTestCases := []struct {
		name               string
		input              string
		expectedTokenCount int
		expectedFirstToken string
		expectedValue      float64
		shouldParseError   bool
		shouldEvalError    bool
	}{
		{"Simple integer", "42", 1, "42", 42.0, false, false},
		{"Simple float", "3.14", 1, "3.14", 3.14, false, false},
		{"Scientific notation", "1e3", 1, "1e3", 1000.0, false, false},
		{"Scientific with plus", "1e+3", 1, "1e+3", 1000.0, false, false},
		{"Scientific with minus", "1e-2", 1, "1e-2", 0.01, false, false},
		{"Complex scientific", "2.5e-3", 1, "2.5e-3", 0.0025, false, false},
		{"Zero", "0", 1, "0", 0.0, false, false},
		{"Negative number", "-42", 2, "-", -42.0, false, true}, // Unary operator, eval fails
		{"Very large number", "123456789012345", 1, "123456789012345", 123456789012345.0, false, false},
		{"Very small decimal", "0.000001", 1, "0.000001", 0.000001, false, false},
		{"Leading zeros", "007", 1, "007", 7.0, false, false},
		{"Decimal without leading zero", ".5", 1, ".", 0.0, false, true}, // Should be parsed as dot operator
	}

	for _, tc := range numberTestCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing input: %s", tc.input)

			// Test tokenizer
			tokenizer := NewTokenizer(tc.input)
			var tokens []Token
			for {
				token := tokenizer.NextToken()
				if token.Type == TokenEOF {
					break
				}
				tokens = append(tokens, token)
				t.Logf("  Token %d: '%s' (Type: %s, Value: %v)", len(tokens), token.Token, token.Type, token.Value)
			}

			// Verify expected token count
			if len(tokens) != tc.expectedTokenCount {
				t.Logf("Expected %d tokens, got %d (this may be expected behavior)", tc.expectedTokenCount, len(tokens))
			}

			if len(tokens) > 0 && tokens[0].Token != tc.expectedFirstToken {
				t.Logf("First token: expected '%s', got '%s'", tc.expectedFirstToken, tokens[0].Token)
			}

			// Test parser
			parser := NewParser(tc.input)
			expr, err := parser.Parse()
			if tc.shouldParseError {
				if err == nil {
					t.Errorf("Expected parse error but got none")
				}
				return
			}

			if err != nil {
				t.Logf("Parser error (may be expected): %v", err)
				return
			}

			t.Logf("  Parsed expression type: %T", expr)

			// Log details about the parsed expression
			switch e := expr.(type) {
			case *NumberLiteral:
				t.Logf("  NumberLiteral value: %s", e.Value)
			case *BinaryExpression:
				t.Logf("  BinaryExpression operator: %s", e.Operator)
			case *UnaryExpression:
				t.Logf("  UnaryExpression operator: %s", e.Operator)
			}

			// Test conversion to AST and evaluation
			node, err := ParseString(tc.input)
			if err != nil {
				t.Logf("  AST conversion error: %v", err)
				return
			}

			result, err := node.Eval(nil)
			if tc.shouldEvalError {
				if err == nil {
					t.Logf("Expected evaluation error but got result: %v", result)
				} else {
					t.Logf("  Expected evaluation error: %v", err)
				}
				return
			}

			if err != nil {
				t.Logf("  Unexpected evaluation error: %v", err)
				return
			}

			t.Logf("  Final result: %v (type: %T)", result, result)
		})
	}

	// Additional test: Complex number expressions
	t.Log("\n=== TESTING COMPLEX NUMBER EXPRESSIONS ===")
	complexCases := []string{
		"1 + 2.5",
		"3.14 * 2",
		"1e3 / 10",
		"2.5e-3 + 1.5e-2",
		"(1 + 2) * 3.14",
	}

	for _, expr := range complexCases {
		t.Logf("Testing complex expression: %s", expr)

		node, err := ParseString(expr)
		if err != nil {
			t.Logf("  Parse error: %v", err)
			continue
		}

		result, err := node.Eval(nil)
		if err != nil {
			t.Logf("  Eval error: %v", err)
			continue
		}

		t.Logf("  Result: %v", result)
	}

	// DEBUG: Testing identifier parsing in tokenizer and parser
	t.Log("\n=== DEBUGGING IDENTIFIER PARSING ===")

	identifierTestCases := []struct {
		name               string
		input              string
		expectedTokenCount int
		expectedFirstToken string
		shouldParseError   bool
		shouldEvalError    bool
		description        string
	}{
		// Basic identifier patterns
		{"Simple variable", "x", 1, "x", false, true, "Simple alpha identifier"},
		{"Uppercase variable", "X", 1, "X", false, true, "Uppercase alpha identifier"},
		{"Mixed case variable", "myVar", 1, "myVar", false, true, "Mixed case identifier"},
		{"Underscore variable", "_var", 1, "_var", false, true, "Underscore start identifier"},
		{"Double underscore", "__var", 1, "__var", false, true, "Double underscore start"},

		// Dollar sign identifiers (common in pipes)
		{"Dollar variable", "$var", 1, "$var", false, true, "Dollar sign start identifier"},
		{"Dollar number", "$1", 1, "$1", false, true, "Dollar with number"},
		{"Dollar underscore", "$_var", 1, "$_var", false, true, "Dollar underscore identifier"},
		{"Double dollar", "$$var", 1, "$$var", false, true, "Double dollar identifier"},

		// Alphanumeric combinations
		{"Alpha numeric", "var123", 1, "var123", false, true, "Alpha numeric identifier"},
		{"Number in middle", "var1test", 1, "var1test", false, true, "Number in middle"},
		{"Underscore numeric", "_123var", 1, "_123var", false, true, "Underscore with numbers"},
		{"Dollar numeric", "$123var", 1, "$123var", false, true, "Dollar with numbers"},

		// Complex valid identifiers
		{"Complex identifier", "myVar_123$test", 1, "myVar_123$test", false, true, "Complex valid identifier"},
		{"All valid chars", "a1_$b2_$c3", 1, "a1_$b2_$c3", false, true, "All valid characters"},
		{"Long identifier", "veryLongVariableNameWith_Numbers123_And$Symbols", 1, "veryLongVariableNameWith_Numbers123_And$Symbols", false, true, "Very long identifier"},

		// Edge cases that should work
		{"Single letter", "a", 1, "a", false, true, "Single letter identifier"},
		{"Single underscore", "_", 1, "_", false, true, "Single underscore"},
		{"Single dollar", "$", 1, "$", false, true, "Single dollar sign"},

		// Invalid patterns (should be tokenized differently)
		{"Number start", "123var", 2, "123", false, true, "Number followed by identifier - should be two tokens"},
		{"Special chars", "var@test", 1, "var", false, true, "Special char should split tokens"},
		{"Hyphen", "var-test", 3, "var", false, true, "Hyphen should split into multiple tokens"},
		{"Dot notation", "obj.prop", 3, "obj", false, true, "Dot notation - multiple tokens"},

		// Keywords that are NOT identifiers
		{"True keyword", "true", 1, "true", false, false, "Boolean true - not identifier"},
		{"False keyword", "false", 1, "false", false, false, "Boolean false - not identifier"},
		{"Null keyword", "null", 1, "null", false, false, "Null keyword - not identifier"},
		{"As keyword", "as", 1, "as", false, true, "As keyword - not identifier"},
	}

	for _, tc := range identifierTestCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing input: %s - %s", tc.input, tc.description)

			// Test tokenizer
			tokenizer := NewTokenizer(tc.input)
			var tokens []Token
			for {
				token := tokenizer.NextToken()
				if token.Type == TokenEOF {
					break
				}
				tokens = append(tokens, token)
				t.Logf("  Token %d: '%s' (Type: %s, Value: %v)", len(tokens), token.Token, token.Type, token.Value)
			}

			// Check if first token matches expected
			if len(tokens) > 0 {
				if tokens[0].Token != tc.expectedFirstToken {
					t.Logf("First token: expected '%s', got '%s'", tc.expectedFirstToken, tokens[0].Token)
				}

				// Log the token type for the first token
				firstToken := tokens[0]
				switch firstToken.Type {
				case TokenIdentifier:
					t.Logf("  ✓ Correctly identified as Identifier")
				case TokenBoolean:
					t.Logf("  ✓ Correctly identified as Boolean keyword")
				case TokenNull:
					t.Logf("  ✓ Correctly identified as Null keyword")
				case TokenAs:
					t.Logf("  ✓ Correctly identified as 'as' keyword")
				case TokenNumber:
					t.Logf("  ✓ Correctly identified as Number (starts with digit)")
				default:
					t.Logf("  Token type: %s", firstToken.Type)
				}
			}

			// Test parser
			parser := NewParser(tc.input)
			expr, err := parser.Parse()
			if tc.shouldParseError {
				if err == nil {
					t.Errorf("Expected parse error but got none")
				}
				return
			}

			if err != nil {
				t.Logf("Parser error (may be expected): %v", err)
				return
			}

			t.Logf("  Parsed expression type: %T", expr)

			// Log details about the parsed expression
			switch e := expr.(type) {
			case *Identifier:
				t.Logf("  ✓ Identifier name: %s", e.Name)
			case *BooleanLiteral:
				t.Logf("  ✓ Boolean value: %v", e.Value)
			case *NullLiteral:
				t.Logf("  ✓ Null literal")
			case *NumberLiteral:
				t.Logf("  ✓ Number literal: %s", e.Value)
			case *BinaryExpression:
				t.Logf("  ✓ Binary expression with operator: %s", e.Operator)
			default:
				t.Logf("  Expression type: %T", e)
			}

			// Test evaluation (identifiers will fail without context)
			node, err := ParseString(tc.input)
			if err != nil {
				t.Logf("  AST conversion error: %v", err)
				return
			}

			result, err := node.Eval(nil)
			if tc.shouldEvalError {
				if err == nil {
					t.Logf("Expected evaluation error but got result: %v", result)
				} else {
					t.Logf("  Expected evaluation error: %v", err)
				}
				return
			}

			if err != nil {
				t.Logf("  Unexpected evaluation error: %v", err)
				return
			}

			t.Logf("  Final result: %v (type: %T)", result, result)
		})
	}
}
