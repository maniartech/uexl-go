package parser_test

import (
	"testing"

	p "github.com/maniartech/uexl_go/parser"
	"github.com/stretchr/testify/assert"
)

func TestOptionalMember_Basic(t *testing.T) {
	parser := p.NewParser("a?.b")
	ast, err := parser.Parse()
	assert.NoError(t, err)

	ma, ok := ast.(*p.MemberAccess)
	if assert.True(t, ok, "expected MemberAccess") {
		assert.True(t, ma.Optional, "member should be optional")
		id, ok := ma.Target.(*p.Identifier)
		if assert.True(t, ok) {
			assert.Equal(t, "a", id.Name)
		}
		assert.Equal(t, "b", ma.Property)
	}
}

func TestOptionalIndex_Basic(t *testing.T) {
	parser := p.NewParser("arr?.[0]")
	ast, err := parser.Parse()
	assert.NoError(t, err)

	ia, ok := ast.(*p.IndexAccess)
	if assert.True(t, ok, "expected IndexAccess") {
		assert.True(t, ia.Optional, "index should be optional")
		_, ok := ia.Target.(*p.Identifier)
		assert.True(t, ok)
	}
}

func TestOptional_MixedChaining(t *testing.T) {
	parser := p.NewParser("a?.b.c?.d?.[i]")
	ast, err := parser.Parse()
	assert.NoError(t, err)

	// ((a?.b).c)?.d?.[i]
	// outermost is IndexAccess (optional)
	idx, ok := ast.(*p.IndexAccess)
	if assert.True(t, ok) {
		assert.True(t, idx.Optional)
		maD, ok := idx.Target.(*p.MemberAccess)
		if assert.True(t, ok) {
			assert.True(t, maD.Optional)
			maC, ok := maD.Target.(*p.MemberAccess)
			if assert.True(t, ok) {
				assert.False(t, maC.Optional)
				maB, ok := maC.Target.(*p.MemberAccess)
				if assert.True(t, ok) {
					assert.True(t, maB.Optional)
				}
			}
		}
	}
}

func TestOptional_Precedence_WithNullish(t *testing.T) {
	parser := p.NewParser("a?.b ?? c")
	ast, err := parser.Parse()
	assert.NoError(t, err)

	be, ok := ast.(*p.BinaryExpression)
	if assert.True(t, ok) {
		assert.Equal(t, "??", be.Operator)
		_, ok := be.Left.(*p.MemberAccess)
		assert.True(t, ok, "left of ?? should be member access")
	}
}

func TestOptional_Precedence_WithLogical(t *testing.T) {
	parser := p.NewParser("a && b?.c")
	ast, err := parser.Parse()
	assert.NoError(t, err)
	be, ok := ast.(*p.BinaryExpression)
	if assert.True(t, ok) {
		assert.Equal(t, "&&", be.Operator)
		_, ok := be.Right.(*p.MemberAccess)
		assert.True(t, ok, "right of && should be member access")
	}
}

func TestOptional_Spacing_NotAllowedBetween_QAndDot(t *testing.T) {
	// 'a? .b' should be parsed as '?' operator then '.b', which will cause an error at parse time
	parser := p.NewParser("a? .b")
	_, err := parser.Parse()
	assert.Error(t, err, "space should break optional operator and likely cause error")
}

func TestOptional_AfterFunctionCall(t *testing.T) {
	parser := p.NewParser("fn(1, 2)?.prop")
	ast, err := parser.Parse()
	assert.NoError(t, err)

	// Expect a MemberAccess with Optional=true
	ma, ok := ast.(*p.MemberAccess)
	if assert.True(t, ok, "expected MemberAccess") {
		assert.True(t, ma.Optional, "member should be optional after call")
		assert.Equal(t, "prop", ma.Property)

		// And its object should be a FunctionCall to identifier 'fn'
		call, ok := ma.Target.(*p.FunctionCall)
		if assert.True(t, ok, "expected FunctionCall as receiver") {
			id, ok := call.Function.(*p.Identifier)
			if assert.True(t, ok) {
				assert.Equal(t, "fn", id.Name)
			}
			assert.Len(t, call.Arguments, 2)
		}
	}
}
