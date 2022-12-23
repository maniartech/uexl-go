package ast

import (
	"fmt"
	"testing"
)

func TestNumber(t *testing.T) {
	node, _ := NewNumberNode("3.14", 0, 1, 1)
	nodeNumber := node.(NumberNode)

	nodeToken := nodeNumber.Token
	if nodeToken != "3.14" {
		t.Errorf("Token: Expected 3.14, got %v", nodeToken)
	}

	nodeValue := nodeNumber.Value
	if nodeValue != 3.14 {
		t.Errorf("Value: Expected 3.14, got %v", nodeValue)
	}

	nodeType := nodeNumber.Type
	if nodeType != "number" {
		t.Errorf("Type: Expected number, got %v", nodeType)
	}

	nodeLine := nodeNumber.Line
	if nodeLine != 1 {
		t.Errorf("Line: Expected 1, got %v", nodeLine)
	}

	nodeCol := nodeNumber.Column
	if nodeCol != 1 {
		t.Errorf("Column: Expected 1, got %v", nodeCol)
	}

	nodeOff := nodeNumber.Offset
	if nodeOff != 0 {
		t.Errorf("Offset: Expected 0, got %v", nodeOff)
	}

	var tmp Map
	nodeEval, _ := nodeNumber.Eval(tmp)
	if nodeEval != nodeValue {
		t.Errorf("Eval Function: Expected %v, got %v", nodeValue, nodeEval)
	}

	NodeNumberStringer := fmt.Sprintf("%v", node)
	if NodeNumberStringer != "NumberNode 3.14" {
		t.Errorf("Stringer: Expected NumberNode 3.14, got %v", NodeNumberStringer)
	}

	_, err := NewNumberNode("xyz", 0, 1, 1)
	if err == nil {
		t.Errorf("Error: Expected %v, got nil", err)
	}
}
