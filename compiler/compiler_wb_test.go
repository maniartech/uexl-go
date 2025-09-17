package compiler

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
)

func TestCompiler_compileShortCircuitChain_empty(t *testing.T) {
	c := New()
	err := c.compileShortCircuitChain([]parser.Node{}, 0)
	if err != nil {
		t.Errorf("expected no error for empty terms, got %v", err)
	}
}

func TestCompiler_compileNullishChain_empty(t *testing.T) {
	c := New()
	err := c.compileNullishChain([]parser.Node{})
	if err != nil {
		t.Errorf("expected no error for empty terms, got %v", err)
	}
}

func TestCompiler_compileAccessNode_baseOnly(t *testing.T) {
	c := New()
	base := &parser.Identifier{Name: "foo"}
	err := c.compileAccessNode(base, false)
	if err != nil {
		t.Errorf("expected no error for base only, got %v", err)
	}
}

func TestCompiler_compileAccessNode_memberAndIndex(t *testing.T) {
	c := New()
	// foo.bar[1]
	idx := &parser.IndexAccess{
		Target: &parser.MemberAccess{
			Target:   &parser.Identifier{Name: "foo"},
			Property: parser.Property{S: "bar"},
		},
		Index: &parser.NumberLiteral{Value: 1},
	}
	err := c.compileAccessNode(idx, false)
	if err != nil {
		t.Errorf("expected no error for member+index, got %v", err)
	}
}

func TestCompiler_compileAccessNode_optionalChain(t *testing.T) {
	c := New()
	// foo?.bar?.baz
	n := &parser.MemberAccess{
		Target: &parser.MemberAccess{
			Target:   &parser.Identifier{Name: "foo"},
			Property: parser.Property{S: "bar"},
			Optional: true,
		},
		Property: parser.Property{S: "baz"},
		Optional: true,
	}
	err := c.compileAccessNode(n, true)
	if err != nil {
		t.Errorf("expected no error for optional chain, got %v", err)
	}
}

func TestCompiler_compileNullishChain_various(t *testing.T) {
	c := New()
	nodes := []parser.Node{
		&parser.Identifier{Name: "foo"},
		&parser.Identifier{Name: "bar"},
	}
	err := c.compileNullishChain(nodes)
	if err != nil {
		t.Errorf("expected no error for nullish chain, got %v", err)
	}
}

func TestFlattenAccessChain_mixed(t *testing.T) {
	// foo.bar[1]?.baz
	n := &parser.MemberAccess{
		Target: &parser.IndexAccess{
			Target: &parser.MemberAccess{
				Target:   &parser.Identifier{Name: "foo"},
				Property: parser.Property{S: "bar"},
			},
			Index:    &parser.NumberLiteral{Value: 1},
			Optional: true,
		},
		Property: parser.Property{S: "baz"},
	}
	base, steps := flattenAccessChain(n)
	if id, ok := base.(*parser.Identifier); !ok || id.Name != "foo" {
		t.Errorf("expected base foo, got %v", base)
	}
	if len(steps) != 3 {
		t.Errorf("expected 3 steps, got %d", len(steps))
	}
}

func TestCompiler_compilePredicateBlock_nil(t *testing.T) {
	c := New()
	idx, err := c.compilePredicateBlock(nil)
	if err != nil {
		t.Errorf("expected no error for nil expr, got %v", err)
	}
	if idx < 0 {
		t.Errorf("expected non-negative index, got %d", idx)
	}
}
