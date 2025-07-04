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
}
