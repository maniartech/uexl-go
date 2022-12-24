package ast

import (
	"fmt"
	"testing"
)

func TestPipe(t *testing.T) {
	pType := []string{"pipe", "empty"}

	arrItem1, _ := NewNumberNode("1", 0, 0, 0)
	arrItem2, _ := NewNumberNode("2", 0, 0, 0)
	arrItem3, _ := NewNumberNode("3", 0, 0, 0)
	arrItems := []Node{arrItem1, arrItem2, arrItem3}
	leftNode, _ := NewArrayNode("[1, 2, 3]", arrItems, 0, 0, 0)

	rightNode1, _ := NewStringNode("hello", 0, 0, 0)
	rightNode2, _ := NewNullNode("null", 0, 0, 0)
	rightNodes := []Node{rightNode1, rightNode2}

	node, _ := NewPipeNode("[1, 2, 3]|:'hello'|empty:null", pType, leftNode, rightNodes, 0, 1, 1)
	nodePipe := node.(PipeNode)

	nodeToken := nodePipe.Token
	if nodeToken != "[1, 2, 3]|:'hello'|empty:null" {
		t.Errorf("Token: Expected [1, 2, 3]|:'hello'|:null, got %v", nodeToken)
	}

	nodeType := nodePipe.Type
	if nodeType != "pipe" {
		t.Errorf("Type: Expected pipe, got %v", nodeType)
	}

	nodeLine := nodePipe.Line
	if nodeLine != 1 {
		t.Errorf("Line: Expected 1, got %v", nodeLine)
	}

	nodeCol := nodePipe.Column
	if nodeCol != 1 {
		t.Errorf("Column: Expected 1, got %v", nodeCol)
	}

	nodeOff := nodePipe.Offset
	if nodeOff != 0 {
		t.Errorf("Offset: Expected 0, got %v", nodeOff)
	}

	var tmp Map
	nodeEval, _ := nodePipe.Eval(tmp)
	nodeEvalArr := nodeEval.([]Node)
	exp1 := nodeEvalArr[0].(ArrayNode)
	expectedVals := []float64{1, 2, 3}
	for i := 0; i < 3; i++ {
		num := exp1.Value[i].(NumberNode).Value
		if num != Number(expectedVals[i]) {
			t.Errorf("Value %v: Expected %v, got %v", i+1, Number(expectedVals[i]), num)
		}
	}

	exp2 := nodeEvalArr[1].(StringNode)
	if exp2.Token != "hello" {
		t.Errorf("Token: Expected hello, got %v", exp2.Token)
	}
	if exp2.Value != "hello" {
		t.Errorf("Token: Expected hello, got %v", exp2.Value)
	}
	exp3 := nodeEvalArr[2].(NullNode)
	if exp3.Token != "null" {
		t.Errorf("Token: Expected null, got %v", exp3.Token)
	}

	nodeObjStringer := fmt.Sprintf("%v", node)
	expected := "PipeNode [ArrayNode [NumberNode 1 NumberNode 2 NumberNode 3] StringNode hello NullNode null]"
	if nodeObjStringer != expected {
		t.Errorf("Stringer: Expected %v, got %v", expected, nodeObjStringer)
	}

}
