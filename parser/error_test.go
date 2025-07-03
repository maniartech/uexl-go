package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/parser/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewErrorSystem(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr errors.ErrorCode
	}{
		{"|: x + 1", errors.ErrEmptyPipe},
		{"x + 2 |: as $a", errors.ErrEmptyPipeWithAlias},
		{"[1, x as $a, 2]", errors.ErrAliasInSubExpr},
		{"func(", errors.ErrUnclosedFunction},
		{"[1, 2,", errors.ErrUnclosedArray},
		{`{"a":`, errors.ErrUnclosedObject},
		{"@", errors.ErrUnexpectedToken}, // @ is not a valid token
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			_, err := p.Parse()
			assert.Error(t, err)

			// Check if it's ParseErrors (multiple errors)
			if parseErr, ok := err.(*errors.ParseErrors); ok {
				assert.True(t, parseErr.HasErrorCode(tt.expectedErr),
					"Expected error code %s, but got %v", tt.expectedErr, parseErr.Errors)
			} else if parserErr, ok := err.(errors.ParserError); ok {
				// Single error
				assert.Equal(t, tt.expectedErr, parserErr.Code,
					"Expected error code %s, but got %s", tt.expectedErr, parserErr.Code)
			} else {
				t.Fatalf("Expected ParserError or ParseErrors, got %T", err)
			}
		})
	}
}

func TestErrorCodeDocumentation(t *testing.T) {
	// Test that we can get error messages for all error codes
	testCodes := []errors.ErrorCode{
		errors.ErrEmptyPipe,
		errors.ErrEmptyPipeWithAlias,
		errors.ErrUnexpectedToken,
		errors.ErrUnclosedArray,
		errors.ErrUnclosedObject,
		errors.ErrInvalidAlias,
	}

	for _, code := range testCodes {
		message := errors.GetErrorMessage(code)
		assert.NotEmpty(t, message, "Error code %s should have a message", code)
		assert.NotEqual(t, errors.GetErrorMessage(errors.ErrUnknown), message,
			"Error code %s should not return unknown error message", code)
	}
}

func TestStructuredErrorInformation(t *testing.T) {
	input := "x + 2 |: as $a"
	p := parser.NewParser(input)
	_, err := p.Parse()

	assert.Error(t, err)

	if parseErr, ok := err.(*errors.ParseErrors); ok {
		assert.Len(t, parseErr.Errors, 1)

		firstError := parseErr.Errors[0]
		assert.Equal(t, errors.ErrEmptyPipeWithAlias, firstError.Code)
		assert.Equal(t, 1, firstError.Line)
		assert.Equal(t, 10, firstError.Column)
		assert.Equal(t, "as", firstError.Token)
		assert.Contains(t, firstError.Message, "empty pipe expression cannot have an alias")
	}
}
