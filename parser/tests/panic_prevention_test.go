package parser_test

import (
	"testing"

	. "github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/parser/constants"
)

// TestPanicPrevention tests that potential panic scenarios are handled gracefully
func TestPanicPrevention(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "malformed raw string",
			input:       "r",
			expectError: false, // "r" is a valid identifier
		},
		{
			name:        "unterminated string",
			input:       "\"unterminated",
			expectError: true,
		},
		{
			name:        "invalid pipe syntax",
			input:       "|invalid:",
			expectError: true,
		},
		{
			name:        "valid expression",
			input:       "1 + 2",
			expectError: false,
		},
		{
			name:        "valid pipe",
			input:       "1 | 2",
			expectError: false,
		},
		{
			name:        "valid raw string",
			input:       "r'hello world'",
			expectError: false,
		},
		{
			name:        "valid array",
			input:       "[1, 2, 3]",
			expectError: false,
		},
		{
			name:        "valid object",
			input:       "{\"a\": 1, \"b\": 2}",
			expectError: false,
		},
		{
			name:        "complex expression",
			input:       "func(a, b) + [1, 2] | sum",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should not panic, even for malformed input
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Parser panicked on input %q: %v", tt.input, r)
				}
			}()

			// Test both the old and new parser creation methods
			parser := NewParser(tt.input)
			result, err := parser.Parse()

			// Also test the new validation method
			validatedParser, validationErr := NewParserWithValidation(tt.input)
			var validatedErr error

			if validationErr == nil {
				_, validatedErr = validatedParser.Parse()
			}

			if tt.expectError {
				// For error cases, at least one method should return an error
				hasError := (err != nil) || (validationErr != nil) || (validatedErr != nil)
				if !hasError {
					t.Errorf("Expected error for input %q, but got none", tt.input)
				}
			} else {
				// For success cases, check that old method works (new method might catch errors early)
				if err != nil {
					t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				}
				if result == nil {
					t.Errorf("Expected result for input %q, but got nil", tt.input)
				}
			}
		})
	}
}

// TestTokenizerBoundsChecking tests that tokenizer handles edge cases without panicking
func TestTokenizerBoundsChecking(t *testing.T) {
	edgeCases := []string{
		"",                    // empty
		"r",                   // incomplete raw string
		"\"",                  // single quote
		"'",                   // single quote
		"|",                   // single pipe
		":",                   // single colon
		"[",                   // single bracket
		"]",                   // single bracket
		"{",                   // single brace
		"}",                   // single brace
		"(",                   // single paren
		")",                   // single paren
		"\\",                  // backslash
		"\n",                  // newline
		"\t",                  // tab
		"   ",                 // spaces
		"r\"unterminated",     // unterminated raw string
		"\"unterminated",      // unterminated string
		"pipe:",               // pipe syntax
		"invalid:pipe",        // invalid pipe syntax
		"123.456.789",         // malformed number
		"true false",          // multiple booleans
		"null null",           // multiple nulls
		"func(",               // incomplete function
		"[1, 2, 3",            // incomplete array
		"{a: 1, b: 2",         // incomplete object
		"1 + 2 + 3 + 4 + 5",   // long expression
		"func(a, b, c, d, e)", // function with many args
	}

	for _, input := range edgeCases {
		t.Run("input_"+input, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Tokenizer panicked on input %q: %v", input, r)
				}
			}()

			tokenizer := NewTokenizer(input)
			tokens := tokenizer.PreloadTokens()

			// Verify that we get at least one token (EOF or error)
			if len(tokens) == 0 {
				t.Errorf("Expected at least one token (EOF or error) for input %q", input)
			}

			// Verify that the last token is EOF or error
			if len(tokens) > 0 {
				lastToken := tokens[len(tokens)-1]
				if lastToken.Type != constants.TokenEOF && lastToken.Type != constants.TokenError {
					t.Errorf("Expected EOF or error token at end for input %q, got %s", input, lastToken.Type)
				}
			}
		})
	}
}
