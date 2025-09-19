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

func TestCompiler_compileAccessNode_integerProperty(t *testing.T) {
	c := New()
	// foo.123 (integer property)
	n := &parser.MemberAccess{
		Target:   &parser.Identifier{Name: "foo"},
		Property: parser.Property{I: 123},
	}
	err := c.compileAccessNode(n, false)
	if err != nil {
		t.Errorf("expected no error for integer property, got %v", err)
	}
}

func TestCompiler_compileAccessNode_softenLastMember(t *testing.T) {
	c := New()
	// foo.bar with softenLast=true
	n := &parser.MemberAccess{
		Target:   &parser.Identifier{Name: "foo"},
		Property: parser.Property{S: "bar"},
	}
	err := c.compileAccessNode(n, true)
	if err != nil {
		t.Errorf("expected no error for soften last member, got %v", err)
	}
}

func TestCompiler_compileAccessNode_softenLastIndex(t *testing.T) {
	c := New()
	// foo[1] with softenLast=true
	n := &parser.IndexAccess{
		Target: &parser.Identifier{Name: "foo"},
		Index:  &parser.NumberLiteral{Value: 1},
	}
	err := c.compileAccessNode(n, true)
	if err != nil {
		t.Errorf("expected no error for soften last index, got %v", err)
	}
}

func TestCompiler_compileNullishChain_memberAccess(t *testing.T) {
	c := New()
	// Test nullish chain with member access (non-last term)
	nodes := []parser.Node{
		&parser.MemberAccess{
			Target:   &parser.Identifier{Name: "foo"},
			Property: parser.Property{S: "bar"},
		},
		&parser.Identifier{Name: "baz"},
	}
	err := c.compileNullishChain(nodes)
	if err != nil {
		t.Errorf("expected no error for nullish chain with member access, got %v", err)
	}
}

func TestCompiler_compileNullishChain_indexAccess(t *testing.T) {
	c := New()
	// Test nullish chain with index access (non-last term)
	nodes := []parser.Node{
		&parser.IndexAccess{
			Target: &parser.Identifier{Name: "foo"},
			Index:  &parser.NumberLiteral{Value: 0},
		},
		&parser.Identifier{Name: "bar"},
	}
	err := c.compileNullishChain(nodes)
	if err != nil {
		t.Errorf("expected no error for nullish chain with index access, got %v", err)
	}
}

func TestCompiler_Compile_unsupportedBinaryOperator(t *testing.T) {
	c := New()
	// Test unsupported binary operator
	n := &parser.BinaryExpression{
		Left:     &parser.NumberLiteral{Value: 1},
		Operator: "???", // unsupported
		Right:    &parser.NumberLiteral{Value: 2},
	}
	err := c.Compile(n)
	if err == nil {
		t.Error("expected error for unsupported binary operator")
	}
}

func TestCompiler_Compile_bitwiseNot(t *testing.T) {
	c := New()
	// Test bitwise NOT operator
	n := &parser.UnaryExpression{
		Operator: "~",
		Operand:  &parser.NumberLiteral{Value: 5},
	}
	err := c.Compile(n)
	if err != nil {
		t.Errorf("expected no error for bitwise NOT, got %v", err)
	}
}

func TestCompile_binaryExpressions_comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		operator string
	}{
		{"addition", "+"},
		{"subtraction", "-"},
		{"multiplication", "*"},
		{"division", "/"},
		{"modulo", "%"},
		{"power", "**"},
		{"equal", "=="},
		{"not_equal", "!="},
		{"greater", ">"},
		{"greater_equal", ">="},
		{"less_equal", "<="},
		{"bitwise_and", "&"},
		{"bitwise_or", "|"},
		{"bitwise_xor", "^"},
		{"shift_left", "<<"},
		{"shift_right", ">>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()

			// Create binary expression: 1 operator 2
			expr := &parser.BinaryExpression{
				Left:     &parser.NumberLiteral{Value: 1.0},
				Operator: tt.operator,
				Right:    &parser.NumberLiteral{Value: 2.0},
			}

			err := c.Compile(expr)
			if err != nil {
				t.Errorf("unexpected error for operator %s: %v", tt.operator, err)
			}

			// Should have emitted some instructions
			if len(c.currentInstructions()) == 0 {
				t.Errorf("expected instructions for operator %s", tt.operator)
			}
		})
	}
}

func TestCompile_literalTypes_comprehensive(t *testing.T) {
	tests := []struct {
		name string
		node parser.Node
	}{
		{"number", &parser.NumberLiteral{Value: 42.0}},
		{"string", &parser.StringLiteral{Value: "hello"}},
		{"boolean_true", &parser.BooleanLiteral{Value: true}},
		{"boolean_false", &parser.BooleanLiteral{Value: false}},
		{"null", &parser.NullLiteral{}},
		{"identifier", &parser.Identifier{Name: "variable"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()

			err := c.Compile(tt.node)
			if err != nil {
				t.Errorf("unexpected error for %s: %v", tt.name, err)
			}
		})
	}
}

func TestCompile_arrayLiteral_comprehensive(t *testing.T) {
	c := New()

	// Test array [1, 2, 3]
	arr := &parser.ArrayLiteral{
		Elements: []parser.Expression{
			&parser.NumberLiteral{Value: 1.0},
			&parser.NumberLiteral{Value: 2.0},
			&parser.NumberLiteral{Value: 3.0},
		},
	}

	err := c.Compile(arr)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCompile_objectLiteral_comprehensive(t *testing.T) {
	c := New()

	// Test object {a: 1, b: 2}
	obj := &parser.ObjectLiteral{
		Properties: map[string]parser.Expression{
			"a": &parser.NumberLiteral{Value: 1.0},
			"b": &parser.NumberLiteral{Value: 2.0},
		},
	}

	err := c.Compile(obj)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCompile_conditionalExpression_comprehensive(t *testing.T) {
	c := New()

	// Test condition ? consequent : alternate
	cond := &parser.ConditionalExpression{
		Condition:  &parser.BooleanLiteral{Value: true},
		Consequent: &parser.NumberLiteral{Value: 1.0},
		Alternate:  &parser.NumberLiteral{Value: 2.0},
	}

	err := c.Compile(cond)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCompile_groupedExpression(t *testing.T) {
	c := New()

	// Test (42)
	grouped := &parser.GroupedExpression{
		Expression: &parser.NumberLiteral{Value: 42.0},
	}

	err := c.Compile(grouped)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCompile_functionCall(t *testing.T) {
	c := New()

	// Test func(arg1, arg2)
	call := &parser.FunctionCall{
		Function: &parser.Identifier{Name: "myFunc"},
		Arguments: []parser.Expression{
			&parser.NumberLiteral{Value: 1.0},
			&parser.NumberLiteral{Value: 2.0},
		},
	}

	err := c.Compile(call)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCompile_sliceExpression(t *testing.T) {
	tests := []struct {
		name string
		expr *parser.SliceExpression
	}{
		{
			"full_slice",
			&parser.SliceExpression{
				Target: &parser.Identifier{Name: "arr"},
				Start:  &parser.NumberLiteral{Value: 1.0},
				End:    &parser.NumberLiteral{Value: 3.0},
				Step:   &parser.NumberLiteral{Value: 1.0},
			},
		},
		{
			"slice_with_nulls",
			&parser.SliceExpression{
				Target:   &parser.Identifier{Name: "arr"},
				Start:    nil,
				End:      nil,
				Step:     nil,
				Optional: true,
			},
		},
		{
			"partial_slice",
			&parser.SliceExpression{
				Target: &parser.Identifier{Name: "arr"},
				Start:  &parser.NumberLiteral{Value: 0.0},
				End:    nil,
				Step:   &parser.NumberLiteral{Value: 2.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()

			err := c.Compile(tt.expr)
			if err != nil {
				t.Errorf("unexpected error for %s: %v", tt.name, err)
			}
		})
	}
}

func TestCompile_programNode(t *testing.T) {
	c := New()

	// Test program with pipe expressions
	prog := &parser.ProgramNode{
		PipeExpressions: []parser.PipeExpression{
			{
				Expression: &parser.NumberLiteral{Value: 42.0},
				PipeType:   "",
				Alias:      "first",
			},
			{
				Expression: &parser.NumberLiteral{Value: 1.0},
				PipeType:   "map",
				Alias:      "",
			},
		},
	}

	err := c.Compile(prog)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCompile_unsupportedUnaryOperator(t *testing.T) {
	c := New()

	// Test unsupported unary operator
	expr := &parser.UnaryExpression{
		Operator: "++", // Unsupported operator
		Operand:  &parser.NumberLiteral{Value: 5.0},
	}

	err := c.Compile(expr)
	// Should not error - it just doesn't emit an instruction for unsupported operators
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCompile_pipeLocalVariable(t *testing.T) {
	c := New()

	// Test pipe local variable (starts with $)
	expr := &parser.Identifier{Name: "$item"}

	err := c.Compile(expr)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
