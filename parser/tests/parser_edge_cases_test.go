package parser_test

import (
	"testing"

	"github.com/maniartech/uexl/parser"
	"github.com/maniartech/uexl/parser/errors"
	"github.com/stretchr/testify/assert"
)

// TestParserEdgeCasesAndErrorPaths tests various edge cases and error paths in the parser
func TestParserEdgeCasesAndErrorPaths(t *testing.T) {
	t.Run("NewParserWithValidation error cases", func(t *testing.T) {
		// Test empty input
		p, err := parser.NewParserWithValidation("")
		assert.Error(t, err)
		assert.Nil(t, p)

		// Test valid input
		p, err = parser.NewParserWithValidation("1 + 2")
		assert.NoError(t, err)
		assert.NotNil(t, p)
	})

	t.Run("Parser error accumulation", func(t *testing.T) {
		// Test invalid syntax that should produce errors
		p := parser.NewParser("1 + + 2")
		expr, err := p.Parse()
		assert.Error(t, err)
		assert.Nil(t, expr)

		// Should have errors
		parseErr, ok := err.(*errors.ParseErrors)
		assert.True(t, ok)
		assert.True(t, len(parseErr.Errors) > 0)
	})

	t.Run("Bitwise OR parsing", func(t *testing.T) {
		// Test that | is parsed as bitwise OR in normal expressions
		p := parser.NewParser("1 | 2")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		binary, ok := expr.(*parser.BinaryExpression)
		assert.True(t, ok)
		assert.Equal(t, "|", binary.Operator)
	})

	t.Run("Unexpected EOF in conditional", func(t *testing.T) {
		p := parser.NewParser("true ?")
		expr, err := p.Parse()
		assert.Error(t, err)
		assert.Nil(t, expr)
	})

	t.Run("Missing colon in conditional", func(t *testing.T) {
		p := parser.NewParser("true ? 1")
		expr, err := p.Parse()
		assert.Error(t, err)
		assert.Nil(t, expr)
	})

	t.Run("Unexpected token after conditional", func(t *testing.T) {
		p := parser.NewParser("true ? 1 :")
		expr, err := p.Parse()
		assert.Error(t, err)
		assert.Nil(t, expr)
	})
}

// TestParserPowerOperatorEdgeCases tests edge cases for the power operator
func TestParserPowerOperatorEdgeCases(t *testing.T) {
	t.Run("Power with unary operators", func(t *testing.T) {
		// Test -2 ** 3 with Excel-style precedence (should be (-2)**3 = -8)
		p := parser.NewParser("-2 ** 3")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		// Should be a binary expression with unary expression as left operand (Excel-style)
		power, ok := expr.(*parser.BinaryExpression)
		assert.True(t, ok)
		assert.Equal(t, "**", power.Operator)

		unary, ok := power.Left.(*parser.UnaryExpression)
		assert.True(t, ok)
		assert.Equal(t, "-", unary.Operator)
	})

	t.Run("Chained power operators", func(t *testing.T) {
		// Test 2 ** 3 ** 2 (should be right-associative: 2 ** (3 ** 2))
		p := parser.NewParser("2 ** 3 ** 2")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		// Should be a binary expression with ** operator
		binary, ok := expr.(*parser.BinaryExpression)
		assert.True(t, ok)
		assert.Equal(t, "**", binary.Operator)

		// Right side should also be a power expression
		rightPower, ok := binary.Right.(*parser.BinaryExpression)
		assert.True(t, ok)
		assert.Equal(t, "**", rightPower.Operator)
	})
}

// TestParserMemberAccessEdgeCases tests edge cases for member access
func TestParserMemberAccessEdgeCases(t *testing.T) {
	t.Run("Simple member access", func(t *testing.T) {
		p := parser.NewParser("obj.prop")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		memberAccess, ok := expr.(*parser.MemberAccess)
		assert.True(t, ok)
		assert.False(t, memberAccess.Optional)
	})

	t.Run("Chained member access", func(t *testing.T) {
		p := parser.NewParser("obj.prop.method")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		// Should be nested member access
		outerAccess, ok := expr.(*parser.MemberAccess)
		assert.True(t, ok)

		innerAccess, ok := outerAccess.Target.(*parser.MemberAccess)
		assert.True(t, ok)
		assert.False(t, innerAccess.Optional)
	})

	t.Run("Index access", func(t *testing.T) {
		p := parser.NewParser("arr[0]")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		// Should be an index access
		indexAccess, ok := expr.(*parser.IndexAccess)
		assert.True(t, ok)
		assert.False(t, indexAccess.Optional)
	})
}

// TestParserSliceExpressionEdgeCases tests edge cases for slice expressions
func TestParserSliceExpressionEdgeCases(t *testing.T) {
	t.Run("Slice with missing end", func(t *testing.T) {
		p := parser.NewParser("arr[1:]")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		slice, ok := expr.(*parser.SliceExpression)
		assert.True(t, ok)
		assert.NotNil(t, slice.Start)
		assert.Nil(t, slice.End)
	})

	t.Run("Slice with missing start", func(t *testing.T) {
		p := parser.NewParser("arr[:5]")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		slice, ok := expr.(*parser.SliceExpression)
		assert.True(t, ok)
		assert.Nil(t, slice.Start)
		assert.NotNil(t, slice.End)
	})

	t.Run("Slice with both start and end", func(t *testing.T) {
		p := parser.NewParser("arr[1:5]")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		slice, ok := expr.(*parser.SliceExpression)
		assert.True(t, ok)
		assert.NotNil(t, slice.Start)
		assert.NotNil(t, slice.End)
	})

	t.Run("Empty slice", func(t *testing.T) {
		p := parser.NewParser("arr[:]")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		slice, ok := expr.(*parser.SliceExpression)
		assert.True(t, ok)
		assert.Nil(t, slice.Start)
		assert.Nil(t, slice.End)
	})
}

// TestParserArrayAndObjectEdgeCases tests edge cases for arrays and objects
func TestParserArrayAndObjectEdgeCases(t *testing.T) {
	t.Run("Empty array", func(t *testing.T) {
		p := parser.NewParser("[]")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		array, ok := expr.(*parser.ArrayLiteral)
		assert.True(t, ok)
		assert.Equal(t, 0, len(array.Elements))
	})

	t.Run("Empty object", func(t *testing.T) {
		p := parser.NewParser("{}")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		obj, ok := expr.(*parser.ObjectLiteral)
		assert.True(t, ok)
		assert.Equal(t, 0, len(obj.Properties))
	})

	t.Run("Array with multiple elements", func(t *testing.T) {
		p := parser.NewParser("[1, 2, 3]")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		array, ok := expr.(*parser.ArrayLiteral)
		assert.True(t, ok)
		assert.Equal(t, 3, len(array.Elements))
	})

	t.Run("Object with multiple properties", func(t *testing.T) {
		p := parser.NewParser(`{"a": 1, "b": 2}`)
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		obj, ok := expr.(*parser.ObjectLiteral)
		assert.True(t, ok)
		assert.Equal(t, 2, len(obj.Properties))
	})

	t.Run("Unclosed array error", func(t *testing.T) {
		p := parser.NewParser("[1, 2, 3")
		expr, err := p.Parse()
		assert.Error(t, err)
		assert.Nil(t, expr)
	})

	t.Run("Unclosed object error", func(t *testing.T) {
		p := parser.NewParser(`{"a": 1, "b": 2`)
		expr, err := p.Parse()
		assert.Error(t, err)
		assert.Nil(t, expr)
	})
}

// TestParserFunctionCallEdgeCases tests edge cases for function calls
func TestParserFunctionCallEdgeCases(t *testing.T) {
	t.Run("Function with no arguments", func(t *testing.T) {
		p := parser.NewParser("func()")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		funcCall, ok := expr.(*parser.FunctionCall)
		assert.True(t, ok)
		assert.Equal(t, 0, len(funcCall.Arguments))
	})

	t.Run("Function with multiple arguments", func(t *testing.T) {
		p := parser.NewParser("func(1, 2, 3)")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		funcCall, ok := expr.(*parser.FunctionCall)
		assert.True(t, ok)
		assert.Equal(t, 3, len(funcCall.Arguments))
	})

	t.Run("Unclosed function call error", func(t *testing.T) {
		p := parser.NewParser("func(1, 2, 3")
		expr, err := p.Parse()
		assert.Error(t, err)
		assert.Nil(t, expr)
	})
}

// TestParserInternalFunctionCoverage tests internal parser functions for coverage
func TestParserInternalFunctionCoverage(t *testing.T) {
	t.Run("Parser with options", func(t *testing.T) {
		opts := parser.DefaultOptions()
		opts.EnableNullish = true
		opts.EnableOptionalChaining = true

		p := parser.NewParserWithOptions("1 ?? 2", opts)
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		binary, ok := expr.(*parser.BinaryExpression)
		assert.True(t, ok)
		assert.Equal(t, "??", binary.Operator)
	})

	t.Run("Grouped expressions", func(t *testing.T) {
		p := parser.NewParser("(1 + 2) * 3")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		binary, ok := expr.(*parser.BinaryExpression)
		assert.True(t, ok)
		assert.Equal(t, "*", binary.Operator)
	})

	t.Run("Nested expressions", func(t *testing.T) {
		p := parser.NewParser("((1 + 2) * 3) / 4")
		expr, err := p.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, expr)

		binary, ok := expr.(*parser.BinaryExpression)
		assert.True(t, ok)
		assert.Equal(t, "/", binary.Operator)
	})
}
