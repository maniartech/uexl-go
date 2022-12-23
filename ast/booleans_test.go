package ast

import (
	"fmt"
	"testing"
)

func TestBooleans(t *testing.T) {
	node, _ := NewBooleanNode("true", 0, 1, 1)
	nodeBool := node.(BooleanNode)

	nodeToken := nodeBool.Token
	if nodeToken != "true" {
		t.Errorf("Token: Expected true, got %v", nodeToken)
	}

	nodeValue := nodeBool.Value
	if nodeValue != true {
		t.Errorf("Value: Expected true, got %v", nodeValue)
	}

	nodeType := nodeBool.Type
	if nodeType != "boolean" {
		t.Errorf("Type: Expected boolean, got %v", nodeType)
	}

	nodeLine := nodeBool.Line
	if nodeLine != 1 {
		t.Errorf("Line: Expected 1, got %v", nodeLine)
	}

	nodeCol := nodeBool.Column
	if nodeCol != 1 {
		t.Errorf("Column: Expected 1, got %v", nodeCol)
	}

	nodeOff := nodeBool.Offset
	if nodeOff != 0 {
		t.Errorf("Offset: Expected 0, got %v", nodeOff)
	}

	var tmp Map
	nodeEval, _ := nodeBool.Eval(tmp)
	if nodeEval != nodeValue {
		t.Errorf("Eval Function: Expected %v, got %v", nodeValue, nodeEval)
	}

	nodeBoolStringer := fmt.Sprintf("%v", node)
	if nodeBoolStringer != "BooleanNode true" {
		t.Errorf("Stringer: Expected BooleanNode true, got %v", nodeBoolStringer)
	}

	_, err := NewBooleanNode("xyz", 0, 1, 1)
	if err == nil {
		t.Errorf("Error: Expected %v, got nil", err)
	}
}
