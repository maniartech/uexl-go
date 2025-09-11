package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestOperatorString tests the Operator.String() method for all operators
func TestOperatorString(t *testing.T) {
	tests := []struct {
		op       Operator
		expected string
	}{
		{OperatorPlus, "+"},
		{OperatorMinus, "-"},
		{OperatorMultiply, "*"},
		{OperatorDivide, "/"},
		{OperatorModulo, "%"},
		{OperatorPower, "**"},
		{OperatorBitwiseAnd, "&"},
		{OperatorBitwiseOr, "|"},
		{OperatorBitwiseXor, "^"},
		{OperatorBitwiseNot, "~"},
		{OperatorLeftShift, "<<"},
		{OperatorRightShift, ">>"},
		{OperatorAssign, "="},
		{OperatorAddAssign, "+="},
		{OperatorSubtractAssign, "-="},
		{OperatorMultiplyAssign, "*="},
		{OperatorDivideAssign, "/="},
		{OperatorModuloAssign, "%="},
		{OperatorBitwiseAndAssign, "&="},
		{OperatorBitwiseOrAssign, "|="},
		{OperatorBitwiseXorAssign, "^="},
		{OperatorIncrement, "++"},
		{OperatorDecrement, "--"},
		{OperatorEqual, "=="},
		{OperatorNotEqual, "!="},
		{OperatorGreaterThan, ">"},
		{OperatorLessThan, "<"},
		{OperatorGreaterOrEqual, ">="},
		{OperatorLessOrEqual, "<="},
		{OperatorLogicalAnd, "&&"},
		{OperatorLogicalOr, "||"},
		{OperatorConditional, "?:"},
		{OperatorPipe, "|:"},
		{OperatorNamedPipe, "|named"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			result := test.op.String()
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestOperatorStringUnknown tests the default case for unknown operators
func TestOperatorStringUnknown(t *testing.T) {
	// Create an operator with a value that doesn't match any defined constants
	unknownOp := Operator(999)
	result := unknownOp.String()
	assert.Equal(t, "Unknown", result)
}

// TestStartOperator tests the StartOperator function
func TestStartOperator(t *testing.T) {
	tests := []struct {
		char     rune
		expected bool
	}{
		{'+', true},
		{'-', true},
		{'*', true},
		{'/', true},
		{'%', true},
		{'&', true},
		{'|', true},
		{'^', true},
		{'~', true},
		{'=', true},
		{'>', true},
		{'<', true},
		{'!', true},
		{'?', true},
		{'a', false},
		{'1', false},
		{' ', false},
		{'(', false},
		{')', false},
		{'[', false},
		{']', false},
		{'{', false},
		{'}', false},
		{',', false},
		{'.', false},
		{':', false},
		{';', false},
		{'#', false},
		{'@', false},
		{'$', false},
	}

	for _, test := range tests {
		t.Run(string(test.char), func(t *testing.T) {
			result := StartOperator(test.char)
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestTokenTypeString tests the TokenType.String() method for all token types
func TestTokenTypeString(t *testing.T) {
	tests := []struct {
		tokenType TokenType
		expected  string
	}{
		{TokenEOF, "EOF"},
		{TokenNumber, "Number"},
		{TokenIdentifier, "Identifier"},
		{TokenOperator, "Operator"},
		{TokenQuestionDot, "QuestionDot"},
		{TokenQuestionLeftBracket, "QuestionLeftBracket"},
		{TokenLeftParen, "LeftParen"},
		{TokenRightParen, "RightParen"},
		{TokenLeftBracket, "LeftBracket"},
		{TokenRightBracket, "RightBracket"},
		{TokenLeftBrace, "LeftBrace"},
		{TokenRightBrace, "RightBrace"},
		{TokenComma, "Comma"},
		{TokenDot, "Dot"},
		{TokenColon, "Colon"},
		{TokenPipe, "Pipe"},
		{TokenString, "String"},
		{TokenBoolean, "Boolean"},
		{TokenNull, "Null"},
		{TokenDollar, "Dollar"},
		{TokenAs, "As"},
		{TokenError, "Error"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			result := test.tokenType.String()
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestTokenTypeStringUnknown tests the default case for unknown token types
func TestTokenTypeStringUnknown(t *testing.T) {
	// Create a token type with a value that doesn't match any defined constants
	unknownToken := TokenType(999)
	result := unknownToken.String()
	assert.Equal(t, "Unknown", result)
}
