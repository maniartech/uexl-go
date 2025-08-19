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
	ma, ok := expr.(*parser.MemberAccess)
	assert.True(t, ok, "expected MemberAccess for .number on array")
	if ok {
		// Property should be an int for dot-number
		v, isInt := ma.Property.(int)
		assert.True(t, isInt, "expected integer property for dot-number")
		assert.Equal(t, 1, v)
	}
}

func TestDotIndexing_Number_OnString(t *testing.T) {
	p := parser.NewParser("'abc'.2")
	expr, err := p.Parse()
	assert.NoError(t, err)
	ma, ok := expr.(*parser.MemberAccess)
	assert.True(t, ok, "expected MemberAccess for .number on string")
	if ok {
		v, isInt := ma.Property.(int)
		assert.True(t, isInt, "expected integer property for dot-number")
		assert.Equal(t, 2, v)
	}
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

func TestDotIndexing_DecimalSplit_OnObject(t *testing.T) {
	p := parser.NewParser("obj.1.5")
	expr, err := p.Parse()
	assert.NoError(t, err)

	// Expect nested MemberAccess: ((obj).1).5
	outer, ok := expr.(*parser.MemberAccess)
	assert.True(t, ok)
	if ok {
		v2, ok2 := outer.Property.(int)
		assert.True(t, ok2)
		assert.Equal(t, 5, v2)
		inner, ok3 := outer.Target.(*parser.MemberAccess)
		assert.True(t, ok3)
		if ok3 {
			v1, ok4 := inner.Property.(int)
			assert.True(t, ok4)
			assert.Equal(t, 1, v1)
		}
	}
}

func TestDotIndexing_DecimalSplit_Chaining(t *testing.T) {
	p := parser.NewParser("obj.1.5.name")
	expr, err := p.Parse()
	assert.NoError(t, err)

	// Expect MemberAccess .name on top of ((obj).1).5
	nameMA, ok := expr.(*parser.MemberAccess)
	assert.True(t, ok)
	if ok {
		assert.Equal(t, "name", nameMA.Property)
		outer, ok2 := nameMA.Target.(*parser.MemberAccess)
		assert.True(t, ok2)
		if ok2 {
			v2, ok3 := outer.Property.(int)
			assert.True(t, ok3)
			assert.Equal(t, 5, v2)
			inner, ok4 := outer.Target.(*parser.MemberAccess)
			assert.True(t, ok4)
			if ok4 {
				v1, ok5 := inner.Property.(int)
				assert.True(t, ok5)
				assert.Equal(t, 1, v1)
			}
		}
	}
}

// Helper to collect integer properties from a nested MemberAccess chain
func collectIntProps(t *testing.T, expr parser.Expression) (props []int, base parser.Expression) {
	t.Helper()
	current := expr
	for {
		ma, ok := current.(*parser.MemberAccess)
		if !ok {
			base = current
			break
		}
		if v, ok := ma.Property.(int); ok {
			props = append(props, v)
		} else if _, ok := ma.Property.(string); ok {
			// Stop if we encounter a string property (used in trailing identifier test)
			// Push a sentinel by returning current with remaining structure
			// Caller can validate the string separately
			base = ma
			break
		} else {
			t.Fatalf("unexpected property type %T", ma.Property)
		}
		current = ma.Target
	}
	return props, base
}

func TestDotIndexing_LongChain_NumberLevels(t *testing.T) {
	inputs := []string{
		"arr.0",
		"arr.0.1",
		"arr.0.1.2",
		"arr.0.1.2.3",
		"arr.0.1.2.3.4",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			p := parser.NewParser(input)
			expr, err := p.Parse()
			assert.NoError(t, err)

			props, base := collectIntProps(t, expr)
			// props are collected from outermost to innermost (e.g., 4,3,2,1,0)
			// Reverse to compare ascending
			for i, j := 0, len(props)-1; i < j; i, j = i+1, j-1 {
				props[i], props[j] = props[j], props[i]
			}
			// Build expected ascending sequence from 0..n
			expected := make([]int, 0)
			// Count dots in input to know how many numeric segments
			// Or simpler: split on '.' and ignore first token 'arr'
			segs := 0
			for _, r := range input {
				if r == '.' {
					segs++
				}
			}
			for i := 0; i < segs; i++ {
				expected = append(expected, i)
			}
			assert.Equal(t, expected, props, "properties sequence should match")

			id, ok := base.(*parser.Identifier)
			assert.True(t, ok, "base should be Identifier")
			if ok {
				assert.Equal(t, "arr", id.Name)
			}
		})
	}
}

func TestDotIndexing_LongChain_WithTrailingIdentifier(t *testing.T) {
	p := parser.NewParser("arr.0.1.2.3.4.name")
	expr, err := p.Parse()
	assert.NoError(t, err)

	// Top-level should be MemberAccess with Property "name"
	top, ok := expr.(*parser.MemberAccess)
	assert.True(t, ok)
	if ok {
		assert.Equal(t, "name", top.Property)
		// Collect numeric properties from the target chain
		props, base := collectIntProps(t, top.Target)
		for i, j := 0, len(props)-1; i < j; i, j = i+1, j-1 {
			props[i], props[j] = props[j], props[i]
		}
		assert.Equal(t, []int{0, 1, 2, 3, 4}, props)
		id, ok := base.(*parser.Identifier)
		assert.True(t, ok)
		if ok {
			assert.Equal(t, "arr", id.Name)
		}
	}
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
