package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/parser/errors"
)

func TestTokenizerErrors(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError errors.ErrorCode
		description   string
	}{
		{
			name:          "consecutive dots",
			input:         "a..b",
			expectedError: errors.ErrConsecutiveDots,
			description:   "should detect consecutive dots in identifier",
		},
		{
			name:          "unterminated string double quote",
			input:         `"hello world`,
			expectedError: errors.ErrUnterminatedQuote,
			description:   "should detect unterminated double quoted string",
		},
		{
			name:          "unterminated string single quote",
			input:         `'hello world`,
			expectedError: errors.ErrUnterminatedQuote,
			description:   "should detect unterminated single quoted string",
		},
		{
			name:          "unterminated raw string",
			input:         `r"hello world`,
			expectedError: errors.ErrUnterminatedQuote,
			description:   "should detect unterminated raw string",
		},
		{
			name:          "invalid character",
			input:         "hello @ world",
			expectedError: errors.ErrInvalidCharacter,
			description:   "should detect invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := parser.NewTokenizer(tt.input)

			// Tokenize until we find an error token
			var errorToken parser.Token
			found := false
			for {
				token := tokenizer.NextToken()
				if token.Type == parser.TokenError {
					errorToken = token
					found = true
					break
				}
				if token.Type == parser.TokenEOF {
					break
				}
			}

			if !found {
				t.Errorf("Expected to find error token for %s, but no error was found", tt.description)
				return
			}

			if errorCode, ok := errorToken.Value.(errors.ErrorCode); !ok || errorCode != tt.expectedError {
				t.Errorf("Expected error code %s for %s, but got %v",
					tt.expectedError, tt.description, errorToken.Value)
			}
		})
	}
}

func TestTokenizerErrorIntegrationWithParser(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError errors.ErrorCode
		description   string
	}{
		{
			name:          "consecutive dots in parser",
			input:         "a..b + 1",
			expectedError: errors.ErrConsecutiveDots,
			description:   "parser should catch consecutive dots error from tokenizer",
		},
		{
			name:          "unterminated string in parser",
			input:         `"hello + 1`,
			expectedError: errors.ErrUnterminatedQuote,
			description:   "parser should catch unterminated string error from tokenizer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := parser.NewParser(tt.input)
			_, err := parser.Parse()

			if err == nil {
				t.Errorf("Expected error for %s, but parsing succeeded", tt.description)
				return
			}

			if parseErrors, ok := err.(*errors.ParseErrors); ok {
				found := false
				for _, parseError := range parseErrors.Errors {
					if parseError.Code == tt.expectedError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error code %s for %s, but got errors: %v",
						tt.expectedError, tt.description, parseErrors.Errors)
				}
			} else {
				t.Errorf("Expected ParseErrors for %s, but got: %T", tt.description, err)
			}
		})
	}
}
