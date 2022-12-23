package ast

import (
	"fmt"
	"testing"
)

func TestArray(t *testing.T) {
	numNode1, _ := NewNumberNode("1", 1, 1, 2)
	numNode2, _ := NewNumberNode("2", 4, 1, 5)
	numNode3, _ := NewNumberNode("3", 7, 1, 8)
	nodes := []Node{numNode1, numNode2, numNode3}

	node, _ := NewArrayNode("[1, 2, 3]", nodes, 0, 1, 1)
	nodeArr := node.(ArrayNode)

	nodeToken := nodeArr.Token
	if nodeToken != "[1, 2, 3]" {
		t.Errorf("Token: Expected [1, 2, 3], got %v", nodeToken)
	}

	nodeType := nodeArr.Type
	if nodeType != "array" {
		t.Errorf("Type: Expected array, got %v", nodeType)
	}

	nodeLine := nodeArr.Line
	if nodeLine != 1 {
		t.Errorf("Line: Expected 1, got %v", nodeLine)
	}

	nodeCol := nodeArr.Column
	if nodeCol != 1 {
		t.Errorf("Column: Expected 1, got %v", nodeCol)
	}

	nodeOff := nodeArr.Offset
	if nodeOff != 0 {
		t.Errorf("Offset: Expected 0, got %v", nodeOff)
	}

	var tmp Map
	nodeEval, _ := nodeArr.Eval(tmp)
	expectedVals := []float64{1, 2, 3}
	for i := 0; i < 3; i++ {
		num := nodeEval.(Array)[i].(NumberNode).Value
		if num != Number(expectedVals[i]) {
			t.Errorf("Value %v: Expected %v, got %v", i+1, Number(expectedVals[i]), num)
		}
	}

	nodeArrStringer := fmt.Sprintf("%v", node)
	expected := "ArrayNode [NumberNode 1 NumberNode 2 NumberNode 3]"
	if nodeArrStringer != expected {
		t.Errorf("Stringer: Expected %v, got %v", expected, nodeArrStringer)
	}
}
