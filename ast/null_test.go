package ast

import (
	"fmt"
	"testing"
)

func TestNull(t *testing.T) {
	node, _ := NewNullNode("null", 0, 1, 1)
	nodeNull := node.(NullNode)

	nodeToken := nodeNull.Token
	if nodeToken != "null" {
		t.Errorf("Token: Expected null, got %v", nodeToken)
	}

	nodeType := nodeNull.Type
	if nodeType != "null" {
		t.Errorf("Type: Expected null, got %v", nodeType)
	}

	nodeLine := nodeNull.Line
	if nodeLine != 1 {
		t.Errorf("Line: Expected 1, got %v", nodeLine)
	}

	nodeCol := nodeNull.Column
	if nodeCol != 1 {
		t.Errorf("Column: Expected 1, got %v", nodeCol)
	}

	nodeOff := nodeNull.Offset
	if nodeOff != 0 {
		t.Errorf("Offset: Expected 0, got %v", nodeOff)
	}

	var tmp Map
	nodeEval, _ := nodeNull.Eval(tmp)
	if nodeEval != nil {
		t.Errorf("Eval Function: Expected nil, got %v", nodeEval)
	}

	nodeNullStringer := fmt.Sprintf("%v", node)
	if nodeNullStringer != "NullNode null" {
		t.Errorf("Stringer: Expected NodeNode null, got %v", nodeNullStringer)
	}
}
