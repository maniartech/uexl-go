package parser_test

import (
	"testing"

	"github.com/maniartech/uexl/parser"
	"github.com/stretchr/testify/assert"
)

func TestBooleans(t *testing.T) {
	// Test simple boolean true
	p := parser.NewParser("true")
	expr, err := p.Parse()
	assert.NoError(t, err)
	boolLit, ok := expr.(*parser.BooleanLiteral)
	assert.True(t, ok, "Expected BooleanLiteral")
	assert.True(t, boolLit.Value, "Expected true value")

	// Test simple boolean false
	p = parser.NewParser("false")
	expr, err = p.Parse()
	assert.NoError(t, err)
	boolLit, ok = expr.(*parser.BooleanLiteral)
	assert.True(t, ok, "Expected BooleanLiteral")
	assert.False(t, boolLit.Value, "Expected false value")

	// Test boolean OR expression: true || false
	p = parser.NewParser("true || false")
	expr, err = p.Parse()
	assert.NoError(t, err)
	binExpr, ok := expr.(*parser.BinaryExpression)
	assert.True(t, ok, "Expected BinaryExpression")
	assert.Equal(t, "||", binExpr.Operator)

	leftBool, ok := binExpr.Left.(*parser.BooleanLiteral)
	assert.True(t, ok, "Expected left operand to be BooleanLiteral")
	assert.True(t, leftBool.Value, "Expected left operand to be true")

	rightBool, ok := binExpr.Right.(*parser.BooleanLiteral)
	assert.True(t, ok, "Expected right operand to be BooleanLiteral")
	assert.False(t, rightBool.Value, "Expected right operand to be false")

	// Test boolean AND expression with parentheses: true && (true || false)
	p = parser.NewParser("true && (true || false)")
	expr, err = p.Parse()
	assert.NoError(t, err)
	binExpr, ok = expr.(*parser.BinaryExpression)
	assert.True(t, ok, "Expected BinaryExpression")
	assert.Equal(t, "&&", binExpr.Operator)

	leftBool, ok = binExpr.Left.(*parser.BooleanLiteral)
	assert.True(t, ok, "Expected left operand to be BooleanLiteral")
	assert.True(t, leftBool.Value, "Expected left operand to be true")

	// The right operand should be a GroupedExpression containing a BinaryExpression
	rightGrouped, ok := binExpr.Right.(*parser.GroupedExpression)
	assert.True(t, ok, "Expected right operand to be GroupedExpression")

	rightBinExpr, ok := rightGrouped.Expression.(*parser.BinaryExpression)
	assert.True(t, ok, "Expected grouped expression to contain BinaryExpression")
	assert.Equal(t, "||", rightBinExpr.Operator)

	rightLeftBool, ok := rightBinExpr.Left.(*parser.BooleanLiteral)
	assert.True(t, ok, "Expected right-left operand to be BooleanLiteral")
	assert.True(t, rightLeftBool.Value, "Expected right-left operand to be true")

	rightRightBool, ok := rightBinExpr.Right.(*parser.BooleanLiteral)
	assert.True(t, ok, "Expected right-right operand to be BooleanLiteral")
	assert.False(t, rightRightBool.Value, "Expected right-right operand to be false")

	// Test complex boolean expression: (true || true) && (true || false)
	p = parser.NewParser("(true || true) && (true || false)")
	expr, err = p.Parse()
	assert.NoError(t, err)
	binExpr, ok = expr.(*parser.BinaryExpression)
	assert.True(t, ok, "Expected BinaryExpression")
	assert.Equal(t, "&&", binExpr.Operator)

	// The left operand should be a GroupedExpression containing a BinaryExpression
	leftGrouped, ok := binExpr.Left.(*parser.GroupedExpression)
	assert.True(t, ok, "Expected left operand to be GroupedExpression")

	leftBinExpr, ok := leftGrouped.Expression.(*parser.BinaryExpression)
	assert.True(t, ok, "Expected grouped expression to contain BinaryExpression")
	assert.Equal(t, "||", leftBinExpr.Operator)

	leftLeftBool, ok := leftBinExpr.Left.(*parser.BooleanLiteral)
	assert.True(t, ok, "Expected left-left operand to be BooleanLiteral")
	assert.True(t, leftLeftBool.Value, "Expected left-left operand to be true")

	leftRightBool, ok := leftBinExpr.Right.(*parser.BooleanLiteral)
	assert.True(t, ok, "Expected left-right operand to be BooleanLiteral")
	assert.True(t, leftRightBool.Value, "Expected left-right operand to be true")

	// The right operand should also be a GroupedExpression containing a BinaryExpression
	rightGrouped, ok = binExpr.Right.(*parser.GroupedExpression)
	assert.True(t, ok, "Expected right operand to be GroupedExpression")

	rightBinExpr, ok = rightGrouped.Expression.(*parser.BinaryExpression)
	assert.True(t, ok, "Expected grouped expression to contain BinaryExpression")
	assert.Equal(t, "||", rightBinExpr.Operator)

	rightLeftBool, ok = rightBinExpr.Left.(*parser.BooleanLiteral)
	assert.True(t, ok, "Expected right-left operand to be BooleanLiteral")
	assert.True(t, rightLeftBool.Value, "Expected right-left operand to be true")

	rightRightBool, ok = rightBinExpr.Right.(*parser.BooleanLiteral)
	assert.True(t, ok, "Expected right-right operand to be BooleanLiteral")
	assert.False(t, rightRightBool.Value, "Expected right-right operand to be false")
}
