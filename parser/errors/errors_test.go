package errors

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestErrorList tests the ErrorList functionality
func TestErrorList(t *testing.T) {
	t.Run("Empty ErrorList", func(t *testing.T) {
		var errList ErrorList

		assert.Equal(t, 0, errList.Len())
		assert.Equal(t, "no errors", errList.Error())
		assert.Nil(t, errList.Err())
	})

	t.Run("Single Error", func(t *testing.T) {
		var errList ErrorList
		err := NewParserError(ErrGeneric, 1, 5, "test error")
		errList.AddError(err)

		assert.Equal(t, 1, errList.Len())
		assert.Equal(t, err.Error(), errList.Error())
		assert.NotNil(t, errList.Err())
	})

	t.Run("Multiple Errors", func(t *testing.T) {
		var errList ErrorList
		err1 := NewParserError(ErrGeneric, 1, 5, "first error")
		err2 := NewParserError(ErrGeneric, 2, 10, "second error")

		errList.AddError(err1)
		errList.AddError(err2)

		assert.Equal(t, 2, errList.Len())

		errorStr := errList.Error()
		assert.True(t, strings.Contains(errorStr, "2 errors:"))
		assert.True(t, strings.Contains(errorStr, "first error"))
		assert.True(t, strings.Contains(errorStr, "second error"))
	})

	t.Run("Add method", func(t *testing.T) {
		var errList ErrorList
		errList.Add(3, 15, "added error")

		assert.Equal(t, 1, errList.Len())
		assert.True(t, strings.Contains(errList.Error(), "added error"))
	})

	t.Run("Reset method", func(t *testing.T) {
		var errList ErrorList
		errList.Add(1, 1, "error")
		assert.Equal(t, 1, errList.Len())

		errList.Reset()
		assert.Equal(t, 0, errList.Len())
	})
}

// TestErrorListSorting tests the sorting functionality of ErrorList
func TestErrorListSorting(t *testing.T) {
	var errList ErrorList

	// Add errors in random order
	errList.AddError(NewParserError(ErrGeneric, 3, 5, "third line"))
	errList.AddError(NewParserError(ErrGeneric, 1, 10, "first line, second column"))
	errList.AddError(NewParserError(ErrGeneric, 1, 5, "first line, first column"))
	errList.AddError(NewParserError(ErrGeneric, 2, 1, "second line"))

	// Test Less method indirectly through sorting
	errList.Sort()

	// Verify order: line 1 col 5, line 1 col 10, line 2 col 1, line 3 col 5
	assert.Equal(t, 1, errList[0].Line)
	assert.Equal(t, 5, errList[0].Column)

	assert.Equal(t, 1, errList[1].Line)
	assert.Equal(t, 10, errList[1].Column)

	assert.Equal(t, 2, errList[2].Line)
	assert.Equal(t, 1, errList[2].Column)

	assert.Equal(t, 3, errList[3].Line)
	assert.Equal(t, 5, errList[3].Column)
}

// TestErrorListSwap tests the Swap method
func TestErrorListSwap(t *testing.T) {
	var errList ErrorList
	err1 := NewParserError(ErrGeneric, 1, 1, "first")
	err2 := NewParserError(ErrGeneric, 2, 2, "second")

	errList.AddError(err1)
	errList.AddError(err2)

	// Verify initial order
	assert.Equal(t, "first", errList[0].Message)
	assert.Equal(t, "second", errList[1].Message)

	// Swap
	errList.Swap(0, 1)

	// Verify swapped order
	assert.Equal(t, "second", errList[0].Message)
	assert.Equal(t, "first", errList[1].Message)
}

// TestErrorListRemoveMultiples tests the RemoveMultiples functionality
func TestErrorListRemoveMultiples(t *testing.T) {
	t.Run("No duplicates", func(t *testing.T) {
		var errList ErrorList
		errList.AddError(NewParserError(ErrGeneric, 1, 1, "first"))
		errList.AddError(NewParserError(ErrGeneric, 2, 2, "second"))

		errList.RemoveMultiples()
		assert.Equal(t, 2, errList.Len())
	})

	t.Run("With duplicates", func(t *testing.T) {
		var errList ErrorList
		errList.AddError(NewParserError(ErrGeneric, 1, 1, "first"))
		errList.AddError(NewParserError(ErrGeneric, 1, 1, "duplicate"))
		errList.AddError(NewParserError(ErrGeneric, 2, 2, "second"))
		errList.AddError(NewParserError(ErrGeneric, 1, 1, "another duplicate"))

		errList.RemoveMultiples()
		assert.Equal(t, 2, errList.Len())

		// Should keep first occurrence of each position
		assert.Equal(t, 1, errList[0].Line)
		assert.Equal(t, 1, errList[0].Column)
		assert.Equal(t, 2, errList[1].Line)
		assert.Equal(t, 2, errList[1].Column)
	})

	t.Run("Single error", func(t *testing.T) {
		var errList ErrorList
		errList.AddError(NewParserError(ErrGeneric, 1, 1, "single"))

		errList.RemoveMultiples()
		assert.Equal(t, 1, errList.Len())
	})

	t.Run("Empty list", func(t *testing.T) {
		var errList ErrorList
		errList.RemoveMultiples()
		assert.Equal(t, 0, errList.Len())
	})
}

// TestParserErrorCreation tests various error creation functions
func TestParserErrorCreation(t *testing.T) {
	t.Run("NewParserError", func(t *testing.T) {
		err := NewParserError(ErrUnexpectedToken, 5, 10, "unexpected token")
		assert.Equal(t, ErrUnexpectedToken, err.Code)
		assert.Equal(t, 5, err.Line)
		assert.Equal(t, 10, err.Column)
		assert.Equal(t, "unexpected token", err.Message)
		assert.Equal(t, "", err.Token)
		assert.Equal(t, "", err.Expected)
	})

	t.Run("NewParserErrorWithToken", func(t *testing.T) {
		err := NewParserErrorWithToken(ErrUnexpectedToken, 3, 7, "bad token", "badtok")
		assert.Equal(t, ErrUnexpectedToken, err.Code)
		assert.Equal(t, 3, err.Line)
		assert.Equal(t, 7, err.Column)
		assert.Equal(t, "bad token", err.Message)
		assert.Equal(t, "badtok", err.Token)
		assert.Equal(t, "", err.Expected)
	})

	t.Run("NewParserErrorWithExpected", func(t *testing.T) {
		err := NewParserErrorWithExpected(ErrExpectedToken, 2, 4, "expected something", "got", "expected")
		assert.Equal(t, ErrExpectedToken, err.Code)
		assert.Equal(t, 2, err.Line)
		assert.Equal(t, 4, err.Column)
		assert.Equal(t, "expected something", err.Message)
		assert.Equal(t, "got", err.Token)
		assert.Equal(t, "expected", err.Expected)
	})
}

// TestParserErrorString tests the Error() method of ParserError
func TestParserErrorString(t *testing.T) {
	t.Run("Basic error", func(t *testing.T) {
		err := NewParserError(ErrGeneric, 1, 1, "basic error")
		errorStr := err.Error()
		assert.True(t, strings.Contains(errorStr, "basic error"))
		assert.True(t, strings.Contains(errorStr, "Line 1, Column 1"))
		assert.True(t, strings.Contains(errorStr, "[generic-error]"))
	})

	t.Run("Error with token", func(t *testing.T) {
		err := NewParserErrorWithToken(ErrUnexpectedToken, 2, 5, "unexpected", "token")
		errorStr := err.Error()
		assert.True(t, strings.Contains(errorStr, "unexpected"))
		assert.True(t, strings.Contains(errorStr, "(token: token)"))
		assert.True(t, strings.Contains(errorStr, "Line 2, Column 5"))
		assert.True(t, strings.Contains(errorStr, "[unexpected-token]"))
	})

	t.Run("Error with expected", func(t *testing.T) {
		err := NewParserErrorWithExpected(ErrExpectedToken, 3, 8, "expected", "got", "wanted")
		errorStr := err.Error()
		assert.True(t, strings.Contains(errorStr, "expected"))
		assert.True(t, strings.Contains(errorStr, "(token: got)"))
		assert.True(t, strings.Contains(errorStr, "Line 3, Column 8"))
		assert.True(t, strings.Contains(errorStr, "[expected-token]"))
	})
}

// TestParseErrors tests the ParseErrors functionality
func TestParseErrors(t *testing.T) {
	t.Run("Single error", func(t *testing.T) {
		err := NewParserError(ErrGeneric, 1, 1, "single error")
		parseErrors := NewParseErrors(err)

		errorStr := parseErrors.Error()
		assert.True(t, strings.Contains(errorStr, "single error"))
		assert.True(t, strings.Contains(errorStr, "[generic-error]"))
	})

	t.Run("Multiple errors", func(t *testing.T) {
		err1 := NewParserError(ErrGeneric, 1, 1, "first error")
		err2 := NewParserError(ErrUnexpectedToken, 2, 5, "second error")
		parseErrors := NewParseErrors(err1, err2)

		errorStr := parseErrors.Error()
		assert.True(t, strings.Contains(errorStr, "parsing failed with 2 errors"))
		assert.True(t, strings.Contains(errorStr, "first error"))
		assert.True(t, strings.Contains(errorStr, "second error"))
	})

	t.Run("HasErrorCode", func(t *testing.T) {
		err1 := NewParserError(ErrGeneric, 1, 1, "error")
		err2 := NewParserError(ErrUnexpectedToken, 2, 5, "error")
		parseErrors := NewParseErrors(err1, err2)

		assert.True(t, parseErrors.HasErrorCode(ErrGeneric))
		assert.True(t, parseErrors.HasErrorCode(ErrUnexpectedToken))
		assert.False(t, parseErrors.HasErrorCode(ErrExpectedToken))
	})

	t.Run("GetErrorsByCode", func(t *testing.T) {
		err1 := NewParserError(ErrGeneric, 1, 1, "first generic")
		err2 := NewParserError(ErrUnexpectedToken, 2, 5, "unexpected")
		err3 := NewParserError(ErrGeneric, 3, 10, "second generic")
		parseErrors := NewParseErrors(err1, err2, err3)

		genericErrors := parseErrors.GetErrorsByCode(ErrGeneric)
		assert.Equal(t, 2, len(genericErrors))
		assert.Equal(t, "first generic", genericErrors[0].Message)
		assert.Equal(t, "second generic", genericErrors[1].Message)

		unexpectedErrors := parseErrors.GetErrorsByCode(ErrUnexpectedToken)
		assert.Equal(t, 1, len(unexpectedErrors))
		assert.Equal(t, "unexpected", unexpectedErrors[0].Message)

		nonExistentErrors := parseErrors.GetErrorsByCode(ErrExpectedToken)
		assert.Equal(t, 0, len(nonExistentErrors))
	})
}

// TestGetErrorMessage tests the GetErrorMessage function
func TestGetErrorMessage(t *testing.T) {
	t.Run("Known error codes", func(t *testing.T) {
		tests := []struct {
			code     ErrorCode
			expected string
		}{
			{ErrUnexpectedToken, "unexpected token"},
			{ErrExpectedToken, "expected token"},
			{ErrEmptyExpression, "empty expression"},
			{ErrInvalidSyntax, "invalid syntax"},
			{ErrGeneric, "generic error"},
		}

		for _, test := range tests {
			t.Run(string(test.code), func(t *testing.T) {
				result := GetErrorMessage(test.code)
				assert.Equal(t, test.expected, result)
			})
		}
	})

	t.Run("Unknown error code", func(t *testing.T) {
		unknownCode := ErrorCode("unknown-code")
		result := GetErrorMessage(unknownCode)
		assert.Equal(t, "unknown error", result)
	})
}
