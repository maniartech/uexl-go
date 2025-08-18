package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/stretchr/testify/assert"
)

// TestUnaryNotBasics verifies parsing of ! and !! with literals and identifiers
func TestUnaryNotBasics(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"not_true", "!true"},
		{"not_false", "!false"},
		{"double_not_identifier", "!!name"},
		{"not_identifier", "!age"},
		{"triple_not_identifier", "!!!x"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := parser.NewParser(c.input)
			expr, err := p.Parse()
			assert.NoError(t, err, "parse should succeed: %s", c.input)
			assert.NotNil(t, expr)

			// At least the outermost should be a UnaryExpression
			_, ok := expr.(*parser.UnaryExpression)
			assert.True(t, ok, "expected UnaryExpression at root for %s", c.input)
		})
	}
}

// TestUnaryNotWithLogical ensures precedence with logical operators and short-circuiting parse shape
func TestUnaryNotWithLogical(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		rootOp       string
		leftIsUnary  bool
		rightIsUnary bool
	}{
		{"not_and", "!a && b", "&&", true, false},
		{"or_not", "a || !b", "||", false, true},
		{"double_not_or", "!!a || b", "||", true, false},
		{"and_double_not", "a && !!b", "&&", false, true},
		{"grouped_not", "!(a && b)", "!", false, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := parser.NewParser(c.input)
			expr, err := p.Parse()
			assert.NoError(t, err, "parse should succeed: %s", c.input)
			assert.NotNil(t, expr)

			if c.rootOp == "!" {
				// Root should be a unary NOT for grouped_not
				un, ok := expr.(*parser.UnaryExpression)
				assert.True(t, ok, "expected UnaryExpression at root for %s", c.input)
				assert.Equal(t, "!", un.Operator)
				return
			}

			// Otherwise expect a binary root with the given operator
			bin, ok := expr.(*parser.BinaryExpression)
			assert.True(t, ok, "expected BinaryExpression at root for %s", c.input)
			assert.Equal(t, c.rootOp, bin.Operator)

			// Check left unary flag
			_, lIsUnary := bin.Left.(*parser.UnaryExpression)
			assert.Equal(t, c.leftIsUnary, lIsUnary, "left unary mismatch for %s", c.input)

			// Check right unary flag
			_, rIsUnary := bin.Right.(*parser.UnaryExpression)
			assert.Equal(t, c.rightIsUnary, rIsUnary, "right unary mismatch for %s", c.input)
		})
	}
}
