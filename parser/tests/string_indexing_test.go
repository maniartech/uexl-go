package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/stretchr/testify/assert"
)

func TestStringIndexing_Basic(t *testing.T) {
	p := parser.NewParser(`"hello"[0]`)
	ast, err := p.Parse()
	assert.NoError(t, err)

	ia, ok := ast.(*parser.IndexAccess)
	if !ok {
		t.Fatalf("expected IndexAccess, got %T", ast)
	}

	// Array should be a StringLiteral("hello")
	sl, ok := ia.Array.(*parser.StringLiteral)
	if !ok {
		t.Fatalf("expected StringLiteral as array, got %T", ia.Array)
	}
	assert.Equal(t, "hello", sl.Value)

	// Index should be NumberLiteral(0)
	nl, ok := ia.Index.(*parser.NumberLiteral)
	if !ok {
		t.Fatalf("expected NumberLiteral as index, got %T", ia.Index)
	}
	assert.Equal(t, 0.0, nl.Value)
}

func TestStringIndexing_WithExpressionIndex(t *testing.T) {
	p := parser.NewParser(`"hello"[1+1]`)
	ast, err := p.Parse()
	assert.NoError(t, err)

	ia, ok := ast.(*parser.IndexAccess)
	if !ok {
		t.Fatalf("expected IndexAccess, got %T", ast)
	}

	// Index should be a BinaryExpression 1 + 1
	be, ok := ia.Index.(*parser.BinaryExpression)
	if !ok {
		t.Fatalf("expected BinaryExpression as index, got %T", ia.Index)
	}
	left, lok := be.Left.(*parser.NumberLiteral)
	right, rok := be.Right.(*parser.NumberLiteral)
	assert.True(t, lok && rok, "index expression should be number + number")
	assert.Equal(t, 1.0, left.Value)
	assert.Equal(t, 1.0, right.Value)
	assert.Equal(t, "+", be.Operator)
}

func TestStringIndexing_WithIdentifierIndex(t *testing.T) {
	p := parser.NewParser(`"abc"[i]`)
	ast, err := p.Parse()
	assert.NoError(t, err)

	ia, ok := ast.(*parser.IndexAccess)
	if !ok {
		t.Fatalf("expected IndexAccess, got %T", ast)
	}
	_, ok = ia.Index.(*parser.Identifier)
	assert.True(t, ok, "index should be an identifier")
}

func TestStringIndexing_Chained(t *testing.T) {
	p := parser.NewParser(`"hi"[0][0]`)
	ast, err := p.Parse()
	assert.NoError(t, err)

	// Outer IndexAccess
	outer, ok := ast.(*parser.IndexAccess)
	if !ok {
		t.Fatalf("expected outer IndexAccess, got %T", ast)
	}
	// Inner IndexAccess
	inner, ok := outer.Array.(*parser.IndexAccess)
	if !ok {
		t.Fatalf("expected inner IndexAccess, got %T", outer.Array)
	}
	// Base should be StringLiteral("hi")
	sl, ok := inner.Array.(*parser.StringLiteral)
	if !ok {
		t.Fatalf("expected StringLiteral base, got %T", inner.Array)
	}
	assert.Equal(t, "hi", sl.Value)
}

func TestStringIndexing_InGroupedExpression(t *testing.T) {
	p := parser.NewParser(`("hello")[2]`)
	ast, err := p.Parse()
	assert.NoError(t, err)
	_, ok := ast.(*parser.IndexAccess)
	assert.True(t, ok, "expected IndexAccess on grouped string literal")
}

func TestStringIndexing_InPipes(t *testing.T) {
	p := parser.NewParser(`"world" |: $1[0]`)
	ast, err := p.Parse()
	assert.NoError(t, err)
	// Should be a ProgramNode with pipe expressions
	_, ok := ast.(*parser.ProgramNode)
	assert.True(t, ok, "expected ProgramNode with pipe where $1[0] is valid")
}

func TestStringIndexing_MissingRightBracketError(t *testing.T) {
	p := parser.NewParser(`"x"[1`)
	_, err := p.Parse()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "]", "should mention expected ']'")
}
