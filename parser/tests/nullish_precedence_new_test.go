package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/stretchr/testify/assert"
)

func TestNullishVsAdditive_Shift(t *testing.T) {
	// a + b ?? c => (a + b) ?? c
	p := parser.NewParser("a + b ?? c")
	expr, err := p.Parse()
	assert.NoError(t, err)
	be, ok := expr.(*parser.BinaryExpression)
	if assert.True(t, ok) {
		assert.Equal(t, "??", be.Operator)
	}

	// a << k ?? 0 => (a << k) ?? 0
	p2 := parser.NewParser("a << k ?? 0")
	expr2, err2 := p2.Parse()
	assert.NoError(t, err2)
	be2, ok := expr2.(*parser.BinaryExpression)
	if assert.True(t, ok) {
		assert.Equal(t, "??", be2.Operator)
	}
}

func TestNullishVsComparison_Equality(t *testing.T) {
	// x == y ?? z => x == (y ?? z)
	p := parser.NewParser("x == y ?? z")
	expr, err := p.Parse()
	assert.NoError(t, err)

	root, ok := expr.(*parser.BinaryExpression)
	if assert.True(t, ok) {
		assert.Equal(t, "==", root.Operator)
		right, ok := root.Right.(*parser.BinaryExpression)
		if assert.True(t, ok) {
			assert.Equal(t, "??", right.Operator)
		}
	}

	// a < b ?? c => a < (b ?? c)
	p2 := parser.NewParser("a < b ?? c")
	expr2, err2 := p2.Parse()
	assert.NoError(t, err2)

	root2, ok := expr2.(*parser.BinaryExpression)
	if assert.True(t, ok) {
		assert.Equal(t, "<", root2.Operator)
		right2, ok := root2.Right.(*parser.BinaryExpression)
		if assert.True(t, ok) {
			assert.Equal(t, "??", right2.Operator)
		}
	}
}

func TestNullishVsBitwise_Logical(t *testing.T) {
	// a & b ?? c => a & (b ?? c)
	p := parser.NewParser("a & b ?? c")
	expr, err := p.Parse()
	assert.NoError(t, err)

	root, ok := expr.(*parser.BinaryExpression)
	if assert.True(t, ok) {
		assert.Equal(t, "&", root.Operator)
		right, ok := root.Right.(*parser.BinaryExpression)
		if assert.True(t, ok) {
			assert.Equal(t, "??", right.Operator)
		}
	}

	// a | b ?? c => a | (b ?? c)
	p2 := parser.NewParser("a | b ?? c")
	expr2, err2 := p2.Parse()
	assert.NoError(t, err2)

	root2, ok := expr2.(*parser.BinaryExpression)
	if assert.True(t, ok) {
		assert.Equal(t, "|", root2.Operator)
		right2, ok := root2.Right.(*parser.BinaryExpression)
		if assert.True(t, ok) {
			assert.Equal(t, "??", right2.Operator)
		}
	}

	// a && b ?? c => a && (b ?? c)
	p3 := parser.NewParser("a && b ?? c")
	expr3, err3 := p3.Parse()
	assert.NoError(t, err3)

	root3, ok := expr3.(*parser.BinaryExpression)
	if assert.True(t, ok) {
		assert.Equal(t, "&&", root3.Operator)
		right3, ok := root3.Right.(*parser.BinaryExpression)
		if assert.True(t, ok) {
			assert.Equal(t, "??", right3.Operator)
		}
	}

	// a || b ?? c => a || (b ?? c)
	p4 := parser.NewParser("a || b ?? c")
	expr4, err4 := p4.Parse()
	assert.NoError(t, err4)

	root4, ok := expr4.(*parser.BinaryExpression)
	if assert.True(t, ok) {
		assert.Equal(t, "||", root4.Operator)
		right4, ok := root4.Right.(*parser.BinaryExpression)
		if assert.True(t, ok) {
			assert.Equal(t, "??", right4.Operator)
		}
	}
}
