package ast

import (
	"fmt"
	"testing"
)

func TestExpressions(t *testing.T) {
	leftNode, _ := NewNumberNode("1", 0, 0, 0)
	rightNode, _ := NewNumberNode("2", 0, 0, 0)

	node := NewExpressionNode("1 / 2", "/", ArithmeticOperator, leftNode, rightNode, 0, 1, 1)

	nodeToken := node.Token
	if nodeToken != "1 / 2" {
		t.Errorf("Token: Expected 1 / 2, got %v", nodeToken)
	}

	nodeType := node.Type
	if nodeType != "expression" {
		t.Errorf("Type: Expected expression, got %v", nodeType)
	}

	nodeLine := node.Line
	if nodeLine != 1 {
		t.Errorf("Line: Expected 1, got %v", nodeLine)
	}

	nodeCol := node.Column
	if nodeCol != 1 {
		t.Errorf("Column: Expected 1, got %v", nodeCol)
	}

	nodeOff := node.Offset
	if nodeOff != 0 {
		t.Errorf("Offset: Expected 0, got %v", nodeOff)
	}

	var tmp Map
	nodeEval, _ := node.Eval(tmp)
	expectedVal := 0
	if nodeEval.(int) != expectedVal {
		t.Errorf("Eval Function: Expected %v, got %v", expectedVal, nodeEval)
	}

	nodeExprStringer := fmt.Sprintf("%v", node)
	expectedExpr := "ExpressionNode NumberNode 1 / NumberNode 2"
	if nodeExprStringer != expectedExpr {
		t.Errorf("Stringer: Expected %v, got %v", expectedExpr, nodeExprStringer)
	}

	// firstVal, _ := NewNumberNode("1", 0, 0, 0)
	// secondVal, _ := NewNumberNode("2", 0, 0, 0)
	// v1 := []interface{}{firstVal}
	// v2 := []interface{}{secondVal}
	// val, _ := ParseExpression("1 + 2", v1, v2, 0, 1, 1)
	// fmt.Println(val)
}
