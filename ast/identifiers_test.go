package ast

import (
	"fmt"
	"testing"
)

func TestIdentifier(t *testing.T) {
	node, _ := NewIdentifierNode("hello.there.world", 0, 1, 1)
	nodeId := node.(IdentifierNode)

	nodeToken := nodeId.Token
	if nodeToken != "hello.there.world" {
		t.Errorf("Token: Expected hello.there.world, got %v", nodeToken)
	}

	nodeName := nodeId.Name
	if nodeName != "hello.there.world" {
		t.Errorf("Value: Expected hello.there.world, got %v", nodeName)
	}

	nodeType := nodeId.Type
	if nodeType != "identifier" {
		t.Errorf("Type: Expected identifier, got %v", nodeType)
	}

	nodeLine := nodeId.Line
	if nodeLine != 1 {
		t.Errorf("Line: Expected 1, got %v", nodeLine)
	}

	nodeCol := nodeId.Column
	if nodeCol != 1 {
		t.Errorf("Column: Expected 1, got %v", nodeCol)
	}

	nodeOff := nodeId.Offset
	if nodeOff != 0 {
		t.Errorf("Offset: Expected 0, got %v", nodeOff)
	}

	var tmp Map
	nodeEval, _ := nodeId.Eval(tmp)
	if nodeEval != nodeName {
		t.Errorf("Eval Function: Expected %v, got %v", nodeName, nodeEval)
	}

	nodeIdinger := fmt.Sprintf("%v", node)
	if nodeIdinger != "IdentifierNode hello.there.world" {
		t.Errorf("Stringer: Expected IdentifierNode hello.there.world, got %v", nodeIdinger)
	}
}
