package ast

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	node, _ := NewStringNode("hello", 0, 1, 1)
	nodeStr := node.(StringNode)

	nodeToken := nodeStr.Token
	if nodeToken != "hello" {
		t.Errorf("Token: Expected hello, got %v", nodeToken)
	}

	nodeValue := nodeStr.Value
	if nodeValue != "hello" {
		t.Errorf("Value: Expected hello, got %v", nodeValue)
	}

	nodeType := nodeStr.Type
	if nodeType != "string" {
		t.Errorf("Type: Expected string, got %v", nodeType)
	}

	nodeLine := nodeStr.Line
	if nodeLine != 1 {
		t.Errorf("Line: Expected 1, got %v", nodeLine)
	}

	nodeCol := nodeStr.Column
	if nodeCol != 1 {
		t.Errorf("Column: Expected 1, got %v", nodeCol)
	}

	nodeOff := nodeStr.Offset
	if nodeOff != 0 {
		t.Errorf("Offset: Expected 0, got %v", nodeOff)
	}

	var tmp Map
	nodeEval, _ := nodeStr.Eval(tmp)
	if nodeEval != nodeValue {
		t.Errorf("Eval Function: Expected %v, got %v", nodeValue, nodeEval)
	}

	nodeStringer := fmt.Sprintf("%v", node)
	if nodeStringer != "StringNode hello" {
		t.Errorf("Stringer: Expected StringNode hello, got %v", nodeStringer)
	}

	node, _ = NewStringNode("'hello'", 0, 1, 1)
	nodeStr = node.(StringNode)
	if nodeStr.Value != "hello" {
		t.Errorf("Value: Expected hello, got %v", nodeStr.Value)
	}
}
