package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/stretchr/testify/assert"
)

func TestDotIndexing_Number_OnArray(t *testing.T) {
	p := parser.NewParser("[10,20,30].1")
	expr, err := p.Parse()
	assert.NoError(t, err)
	_, ok := expr.(*parser.MemberAccess)
	assert.True(t, ok, "expected MemberAccess for .number on array")
}

func TestDotIndexing_Number_OnString(t *testing.T) {
	p := parser.NewParser("'abc'.2")
	expr, err := p.Parse()
	assert.NoError(t, err)
	_, ok := expr.(*parser.MemberAccess)
	assert.True(t, ok, "expected MemberAccess for .number on string")
}

func TestDotIndexing_GroupedExpr_OnArray(t *testing.T) {
	p := parser.NewParser("[1,2,3].(1+1)")
	expr, err := p.Parse()
	assert.NoError(t, err)
	_, ok := expr.(*parser.IndexAccess)
	assert.True(t, ok, "expected IndexAccess for .(expr) on array")
}

func TestDotIndexing_GroupedExpr_OnString(t *testing.T) {
	p := parser.NewParser("\"hello\".(2)")
	expr, err := p.Parse()
	assert.NoError(t, err)
	_, ok := expr.(*parser.IndexAccess)
	assert.True(t, ok, "expected IndexAccess for .(expr) on string")
}

func TestDotIndexing_Chaining_Mixed(t *testing.T) {
	p := parser.NewParser("[ ['a','b','c'], ['d','e'] ].0[1]")
	expr, err := p.Parse()
	assert.NoError(t, err)
	// Should parse as IndexAccess([1]) of MemberAccess(.0)
	outer, ok := expr.(*parser.IndexAccess)
	assert.True(t, ok)
	_, ok = outer.Target.(*parser.MemberAccess)
	assert.True(t, ok)
}

func TestDotIndexing_DoesNotClash_WithMemberAccess(t *testing.T) {
	p := parser.NewParser("obj.prop.0")
	expr, err := p.Parse()
	assert.NoError(t, err)
	// obj.prop is a MemberAccess; .0 should now be MemberAccess on that (consistent syntax-based approach)
	ma, ok := expr.(*parser.MemberAccess)
	assert.True(t, ok)
	_, ok = ma.Target.(*parser.MemberAccess)
	assert.True(t, ok)
}

func TestDotIndexing_InPipes(t *testing.T) {
	p := parser.NewParser("data |map: $1.items.0")
	expr, err := p.Parse()
	assert.NoError(t, err)
	assert.NotNil(t, expr)
}

func TestDotIndexing_IdentifierIndex_OnArray(t *testing.T) {
	p := parser.NewParser("arr.(i)")
	expr, err := p.Parse()
	assert.NoError(t, err)
	idx, ok := expr.(*parser.IndexAccess)
	assert.True(t, ok)
	_, ok = idx.Index.(*parser.Identifier)
	assert.True(t, ok, "index should be an Identifier expression")
}

func TestDotIndexing_IdentifierIndex_OnString(t *testing.T) {
	p := parser.NewParser("'abc'.(i)")
	expr, err := p.Parse()
	assert.NoError(t, err)
	idx, ok := expr.(*parser.IndexAccess)
	assert.True(t, ok)
	_, ok = idx.Index.(*parser.Identifier)
	assert.True(t, ok, "index should be an Identifier expression")
}
