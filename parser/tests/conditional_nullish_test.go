package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/stretchr/testify/assert"
)

func TestNullishBasic(t *testing.T) {
	p := parser.NewParser("sum(1, 2, 3) ?? b")
	expr, err := p.Parse()
	assert.NoError(t, err)

	be, ok := expr.(*parser.BinaryExpression)
	if assert.True(t, ok, "expected BinaryExpression") {
		assert.Equal(t, "??", be.Operator)
	}
}

func TestNullishWithLogicalOrAndAnd(t *testing.T) {
	p := parser.NewParser("a ?? b || c && d")
	expr, err := p.Parse()
	assert.NoError(t, err)

	// Root should be ||, left is (a ?? b), right is (c && d)
	root, ok := expr.(*parser.BinaryExpression)
	if assert.True(t, ok) {
		assert.Equal(t, "||", root.Operator)

		left, ok := root.Left.(*parser.BinaryExpression)
		if assert.True(t, ok) {
			assert.Equal(t, "??", left.Operator)
		}

		right, ok := root.Right.(*parser.BinaryExpression)
		if assert.True(t, ok) {
			assert.Equal(t, "&&", right.Operator)
		}
	}
}

func TestTernaryBasic(t *testing.T) {
	p := parser.NewParser("a ? b : c")
	expr, err := p.Parse()
	assert.NoError(t, err)

	ce, ok := expr.(*parser.ConditionalExpression)
	assert.True(t, ok, "expected ConditionalExpression, got %T", expr)
	assert.NotNil(t, ce.Condition)
	assert.NotNil(t, ce.Consequent)
	assert.NotNil(t, ce.Alternate)
}

func TestTernaryRightAssociative(t *testing.T) {
	p := parser.NewParser("a ? b : c ? d : e")
	expr, err := p.Parse()
	assert.NoError(t, err)

	ce, ok := expr.(*parser.ConditionalExpression)
	if assert.True(t, ok) {
		// alternate should be another conditional
		inner, ok := ce.Alternate.(*parser.ConditionalExpression)
		assert.True(t, ok, "expected nested ConditionalExpression on alternate")
		_ = inner
	}
}

func TestTernaryPrecedenceOverNullish(t *testing.T) {
	p := parser.NewParser("a ?? b ? c : d")
	expr, err := p.Parse()
	assert.NoError(t, err)

	// Should be parsed as (a ?? b) ? c : d
	ce, ok := expr.(*parser.ConditionalExpression)
	if assert.True(t, ok) {
		cond, ok := ce.Condition.(*parser.BinaryExpression)
		if assert.True(t, ok) {
			assert.Equal(t, "??", cond.Operator)
		}
	}
}
