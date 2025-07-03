package errors

import (
	"fmt"
	"strings"
)

// ErrorCode represents a unique identifier for parser errors
type ErrorCode string

const (
	// Syntax Errors
	ErrUnexpectedToken ErrorCode = "unexpected-token"
	ErrExpectedToken   ErrorCode = "expected-token"
	ErrUnexpectedEOF   ErrorCode = "unexpected-eof"
	ErrInvalidSyntax   ErrorCode = "invalid-syntax"

	// Expression Errors
	ErrEmptyExpression   ErrorCode = "empty-expression"
	ErrInvalidExpression ErrorCode = "invalid-expression"
	ErrExpressionTooLong ErrorCode = "expression-too-long"

	// Pipe Errors
	ErrEmptyPipe           ErrorCode = "empty-pipe"
	ErrEmptyPipeWithAlias  ErrorCode = "empty-pipe-with-alias"
	ErrPipeInSubExpression ErrorCode = "pipe-in-sub-expression"
	ErrInvalidPipeType     ErrorCode = "invalid-pipe-type"
	ErrMissingPipeType     ErrorCode = "missing-pipe-type"

	// Identifier/Alias Errors
	ErrExpectedIdentifier ErrorCode = "expected-identifier"
	ErrInvalidAlias       ErrorCode = "invalid-alias"
	ErrMissingDollarSign  ErrorCode = "missing-dollar-sign"
	ErrAliasInSubExpr     ErrorCode = "alias-in-sub-expression"

	// Literal Errors
	ErrInvalidNumber      ErrorCode = "invalid-number"
	ErrInvalidString      ErrorCode = "invalid-string"
	ErrUnterminatedString ErrorCode = "unterminated-string"

	// Collection Errors
	ErrUnclosedArray       ErrorCode = "unclosed-array"
	ErrUnclosedObject      ErrorCode = "unclosed-object"
	ErrInvalidArrayElement ErrorCode = "invalid-array-element"
	ErrInvalidObjectKey    ErrorCode = "invalid-object-key"
	ErrInvalidObjectValue  ErrorCode = "invalid-object-value"

	// Function Errors
	ErrUnclosedFunction ErrorCode = "unclosed-function"
	ErrInvalidArgument  ErrorCode = "invalid-argument"

	// Tokenizer Errors
	ErrConsecutiveDots   ErrorCode = "consecutive-dots"
	ErrInvalidToken      ErrorCode = "invalid-token"
	ErrInvalidCharacter  ErrorCode = "invalid-character"
	ErrUnterminatedQuote ErrorCode = "unterminated-quote"
	ErrTooManyArguments  ErrorCode = "too-many-arguments"

	// Operator Errors
	ErrInvalidOperator ErrorCode = "invalid-operator"
	ErrMissingOperand  ErrorCode = "missing-operand"
	ErrInvalidUnaryOp  ErrorCode = "invalid-unary-operator"

	// General Errors
	ErrInternal ErrorCode = "internal-error"
	ErrUnknown  ErrorCode = "unknown-error"
)

// ParserError represents a structured parser error
type ParserError struct {
	Code     ErrorCode `json:"code"`
	Message  string    `json:"message"`
	Line     int       `json:"line"`
	Column   int       `json:"column"`
	Token    string    `json:"token,omitempty"`
	Expected string    `json:"expected,omitempty"`
	Context  string    `json:"context,omitempty"`
}

func (e ParserError) Error() string {
	if e.Token != "" {
		return fmt.Sprintf("[%s] Line %d, Column %d: %s (token: %s)",
			e.Code, e.Line, e.Column, e.Message, e.Token)
	}
	return fmt.Sprintf("[%s] Line %d, Column %d: %s",
		e.Code, e.Line, e.Column, e.Message)
}

// ParseErrors represents multiple parser errors
type ParseErrors struct {
	Errors []ParserError `json:"errors"`
}

func (pe ParseErrors) Error() string {
	if len(pe.Errors) == 1 {
		return pe.Errors[0].Error()
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("parsing failed with %d errors:\n", len(pe.Errors)))
	for i, err := range pe.Errors {
		result.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err.Error()))
	}
	return result.String()
}

// HasErrorCode checks if any error has the specified code
func (pe ParseErrors) HasErrorCode(code ErrorCode) bool {
	for _, err := range pe.Errors {
		if err.Code == code {
			return true
		}
	}
	return false
}

// GetErrorsByCode returns all errors with the specified code
func (pe ParseErrors) GetErrorsByCode(code ErrorCode) []ParserError {
	var result []ParserError
	for _, err := range pe.Errors {
		if err.Code == code {
			result = append(result, err)
		}
	}
	return result
}

// Error message templates for consistent messaging
var errorMessages = map[ErrorCode]string{
	ErrUnexpectedToken: "unexpected token",
	ErrExpectedToken:   "expected token",
	ErrUnexpectedEOF:   "unexpected end of input",
	ErrInvalidSyntax:   "invalid syntax",

	ErrEmptyExpression:   "empty expression",
	ErrInvalidExpression: "invalid expression",
	ErrExpressionTooLong: "expression too long",

	ErrEmptyPipe:           "empty pipe expression is not allowed",
	ErrEmptyPipeWithAlias:  "empty pipe expression cannot have an alias",
	ErrPipeInSubExpression: "pipe expressions cannot be sub-expressions",
	ErrInvalidPipeType:     "invalid pipe type",
	ErrMissingPipeType:     "missing pipe type or empty pipe segment",

	ErrExpectedIdentifier: "expected identifier",
	ErrInvalidAlias:       "invalid alias",
	ErrMissingDollarSign:  "expected identifier starting with $",
	ErrAliasInSubExpr:     "aliases cannot be used in sub-expressions",

	ErrInvalidNumber:      "invalid number format",
	ErrInvalidString:      "invalid string format",
	ErrUnterminatedString: "unterminated string",

	ErrUnclosedArray:       "unclosed array, expected ']'",
	ErrUnclosedObject:      "unclosed object, expected '}'",
	ErrInvalidArrayElement: "invalid array element",
	ErrInvalidObjectKey:    "invalid object key",
	ErrInvalidObjectValue:  "invalid object value",

	ErrUnclosedFunction: "unclosed function call, expected ')'",
	ErrInvalidArgument:  "invalid function argument",
	ErrTooManyArguments: "too many function arguments",

	ErrConsecutiveDots:   "consecutive dots in identifier",
	ErrInvalidToken:      "invalid token",
	ErrInvalidCharacter:  "invalid character",
	ErrUnterminatedQuote: "unterminated quoted string",

	ErrInvalidOperator: "invalid operator",
	ErrMissingOperand:  "missing operand",
	ErrInvalidUnaryOp:  "invalid unary operator",

	ErrInternal: "internal parser error",
	ErrUnknown:  "unknown error",
}

// GetErrorMessage returns the default message for an error code
func GetErrorMessage(code ErrorCode) string {
	if msg, exists := errorMessages[code]; exists {
		return msg
	}
	return errorMessages[ErrUnknown]
}

// NewParserError creates a new parser error with the given details
func NewParserError(code ErrorCode, line, column int, message string) ParserError {
	return ParserError{
		Code:    code,
		Message: message,
		Line:    line,
		Column:  column,
	}
}

// NewParserErrorWithToken creates a parser error with token information
func NewParserErrorWithToken(code ErrorCode, line, column int, message, token string) ParserError {
	return ParserError{
		Code:    code,
		Message: message,
		Line:    line,
		Column:  column,
		Token:   token,
	}
}

// NewParserErrorWithExpected creates a parser error with expected token information
func NewParserErrorWithExpected(code ErrorCode, line, column int, message, token, expected string) ParserError {
	return ParserError{
		Code:     code,
		Message:  message,
		Line:     line,
		Column:   column,
		Token:    token,
		Expected: expected,
	}
}

// NewParseErrors creates a new ParseErrors with the given errors
func NewParseErrors(errors ...ParserError) *ParseErrors {
	return &ParseErrors{Errors: errors}
}
