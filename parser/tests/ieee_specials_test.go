package parser_test

import (
	"math"
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/parser/constants"
	"github.com/stretchr/testify/assert"
)

func TestIeeeSpecials_ExplicitlyDisabled(t *testing.T) {
	// When explicitly disabled, NaN and Inf should be parsed as identifiers
	opt := parser.DefaultOptions()
	opt.EnableIeeeSpecials = false
	tz := parser.NewTokenizerWithOptions("NaN", opt)
	tok, err := tz.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenIdentifier, tok.Type)
	assert.Equal(t, "NaN", tok.Token)

	tz2 := parser.NewTokenizerWithOptions("Inf", opt)
	tok2, err := tz2.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenIdentifier, tok2.Type)
	assert.Equal(t, "Inf", tok2.Token)
}

func TestIeeeSpecials_TokenizerEnabled(t *testing.T) {
	opt := parser.DefaultOptions()
	opt.EnableIeeeSpecials = true
	tz := parser.NewTokenizerWithOptions("NaN", opt)
	tok, err := tz.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenNumber, tok.Type)
	assert.True(t, math.IsNaN(tok.Value.Num))

	tz2 := parser.NewTokenizerWithOptions("Inf", opt)
	tok2, err := tz2.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenNumber, tok2.Type)
	assert.True(t, math.IsInf(tok2.Value.Num, +1))
}

func TestIeeeSpecials_DefaultTokenizer(t *testing.T) {
	tz := parser.NewTokenizer("NaN")
	tok, err := tz.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenNumber, tok.Type)
	assert.True(t, math.IsNaN(tok.Value.Num))

	tz2 := parser.NewTokenizer("Inf")
	tok2, err := tz2.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenNumber, tok2.Type)
	assert.True(t, math.IsInf(tok2.Value.Num, +1))
}

func TestIeeeSpecials_ParserUnaryMinus(t *testing.T) {
	opt := parser.DefaultOptions()
	opt.EnableIeeeSpecials = true
	p := parser.NewParserWithOptions("-Inf", opt)
	expr, err := p.Parse()
	assert.NoError(t, err)
	u, ok := expr.(*parser.UnaryExpression)
	if assert.True(t, ok, "expected unary expression") {
		assert.Equal(t, "-", u.Operator)
		// inner should be a number literal with +Inf; unary minus applied at evaluation time, but AST holds +Inf token
		n, ok := u.Operand.(*parser.NumberLiteral)
		assert.True(t, ok)
		assert.True(t, math.IsInf(n.Value, +1))
	}
}

func TestIeeeSpecials_PowerOperatorPrecedence(t *testing.T) {
	// This tests Excel-style power operator precedence for all numbers (including IEEE-754 specials)
	// Excel model: unary operators bind tighter than power operator
	// -2**2 should be parsed as (-2)**2 = 4, not -(2**2) = -4
	// This is more intuitive for expression language users
	opt := parser.DefaultOptions()
	opt.EnableIeeeSpecials = true

	tests := []struct {
		input    string
		expected string // AST structure description
	}{
		{
			input:    "-Inf ** 2",
			expected: "BinaryExpr with UnaryExpr(-Inf) as left operand", // (-Inf) ** 2 - Excel style
		},
		{
			input:    "-(Inf ** 2)",
			expected: "UnaryExpr with GroupedExpr containing BinaryExpr", // -(Inf ** 2) - explicit grouping
		},
		{
			input:    "(-Inf) ** 2",
			expected: "BinaryExpr with GroupedExpr containing UnaryExpr as left operand", // (-Inf) ** 2 - explicit parentheses
		},
		{
			input:    "-2 ** 3",
			expected: "BinaryExpr with UnaryExpr(-2) as left operand", // (-2) ** 3 - consistent with Excel
		},
		{
			input:    "-NaN ** 2",
			expected: "BinaryExpr with UnaryExpr(-NaN) as left operand", // (-NaN) ** 2 - consistent with other numbers
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.NewParserWithOptions(tt.input, opt)
			expr, err := p.Parse()
			assert.NoError(t, err)

			switch tt.input {
			case "-Inf ** 2", "-2 ** 3", "-NaN ** 2":
				// Should be BinaryExpr with UnaryExpr as left operand (Excel-style precedence)
				binExpr, ok := expr.(*parser.BinaryExpression)
				if assert.True(t, ok, "expected binary expression for %s", tt.input) {
					assert.Equal(t, "**", binExpr.Operator)

					// Left should be UnaryExpr
					unaryExpr, ok := binExpr.Left.(*parser.UnaryExpression)
					if assert.True(t, ok, "expected unary expression as left operand") {
						assert.Equal(t, "-", unaryExpr.Operator)

						// Operand should be NumberLiteral
						leftNum, ok := unaryExpr.Operand.(*parser.NumberLiteral)
						assert.True(t, ok, "expected number literal as unary operand")

						if tt.input == "-Inf ** 2" {
							assert.True(t, math.IsInf(leftNum.Value, 1), "expected +Inf")
						} else if tt.input == "-NaN ** 2" {
							assert.True(t, math.IsNaN(leftNum.Value), "expected NaN")
						} else {
							assert.Equal(t, 2.0, leftNum.Value)
						}
					}

					// Right should be NumberLiteral
					rightNum, ok := binExpr.Right.(*parser.NumberLiteral)
					if assert.True(t, ok, "expected number literal as right operand") {
						if tt.input == "-2 ** 3" {
							assert.Equal(t, 3.0, rightNum.Value)
						} else {
							assert.Equal(t, 2.0, rightNum.Value)
						}
					}
				}

			case "-(Inf ** 2)":
				// Should be UnaryExpr with GroupedExpr containing BinaryExpr
				unaryExpr, ok := expr.(*parser.UnaryExpression)
				if assert.True(t, ok, "expected unary expression for %s", tt.input) {
					assert.Equal(t, "-", unaryExpr.Operator)

					// Operand should be GroupedExpr
					groupedExpr, ok := unaryExpr.Operand.(*parser.GroupedExpression)
					if assert.True(t, ok, "expected grouped expression as unary operand") {
						// Inner expression should be BinaryExpr
						binExpr, ok := groupedExpr.Expression.(*parser.BinaryExpression)
						if assert.True(t, ok, "expected binary expression inside grouped expression") {
							assert.Equal(t, "**", binExpr.Operator)
						}
					}
				}

			case "(-Inf) ** 2":
				// Should be BinaryExpr with GroupedExpr containing UnaryExpr as left operand
				binExpr, ok := expr.(*parser.BinaryExpression)
				if assert.True(t, ok, "expected binary expression for %s", tt.input) {
					assert.Equal(t, "**", binExpr.Operator)

					// Left should be GroupedExpr
					groupedExpr, ok := binExpr.Left.(*parser.GroupedExpression)
					if assert.True(t, ok, "expected grouped expression as left operand") {
						// Inner expression should be UnaryExpr
						unaryExpr, ok := groupedExpr.Expression.(*parser.UnaryExpression)
						if assert.True(t, ok, "expected unary expression inside grouped expression") {
							assert.Equal(t, "-", unaryExpr.Operator)
						}
					}

					// Right should be NumberLiteral
					rightNum, ok := binExpr.Right.(*parser.NumberLiteral)
					if assert.True(t, ok, "expected number literal as right operand") {
						assert.Equal(t, 2.0, rightNum.Value)
					}
				}
			}
		})
	}
}
