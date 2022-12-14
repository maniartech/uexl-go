package parser

import (
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/ast"
)

func TestArray(t *testing.T) {
	got, _ := ParseReader("", strings.NewReader("[1, 2, 3]"))
	gotNode := got.(ast.ArrayNode)
	if gotNode.Value[0].(ast.NumberNode).Value != 1 {
		t.Errorf("Expected 1, got %v", gotNode.Value[0].(ast.NumberNode).Value)
	}
	if gotNode.Value[1].(ast.NumberNode).Value != 2 {
		t.Errorf("Expected 2, got %v", gotNode.Value[1].(ast.NumberNode).Value)
	}
	if gotNode.Value[2].(ast.NumberNode).Value != 3 {
		t.Errorf("Expected 3, got %v", gotNode.Value[2].(ast.NumberNode).Value)
	}

	got, _ = ParseReader("", strings.NewReader("[1 + 2, 3 + 4, true || false]"))
	gotNode = got.(ast.ArrayNode)
	gotExpNode1 := gotNode.Value[0].(ast.ExpressionNode)
	gotExpNode1Token := gotExpNode1.Token
	if gotExpNode1Token != "1 + 2" {
		t.Errorf("Token: Expected 1 + 2, got %v", gotExpNode1Token)
	}
	expNode1LeftToken := gotExpNode1.Left.(ast.NumberNode).Token
	expNode1LeftValue := gotExpNode1.Left.(ast.NumberNode).Value
	expNode1RightToken := gotExpNode1.Right.(ast.NumberNode).Token
	expNode1RightValue := gotExpNode1.Right.(ast.NumberNode).Value
	if expNode1LeftToken != "1" {
		t.Errorf("Token: Expected 1, got %v", expNode1LeftToken)
	}
	if expNode1LeftValue != 1 {
		t.Errorf("Value: Expected 1, got %v", expNode1LeftValue)
	}
	if expNode1RightToken != "2" {
		t.Errorf("Token: Expected 2, got %v", expNode1RightToken)
	}
	if expNode1RightValue != 2 {
		t.Errorf("Value: Expected 1, got %v", expNode1RightValue)
	}

	gotExpNode2 := gotNode.Value[1].(ast.ExpressionNode)
	gotExpNode2Token := gotExpNode2.Token
	if gotExpNode2Token != "3 + 4" {
		t.Errorf("Token: Expected 3 + 4, got %v", gotExpNode2Token)
	}
	expNode2LeftToken := gotExpNode2.Left.(ast.NumberNode).Token
	expNode2LeftValue := gotExpNode2.Left.(ast.NumberNode).Value
	expNode2RightToken := gotExpNode2.Right.(ast.NumberNode).Token
	expNode2RightValue := gotExpNode2.Right.(ast.NumberNode).Value
	if expNode2LeftToken != "3" {
		t.Errorf("Token: Expected 3, got %v", expNode2LeftToken)
	}
	if expNode2LeftValue != 3 {
		t.Errorf("Value: Expected 3, got %v", expNode2LeftValue)
	}
	if expNode2RightToken != "4" {
		t.Errorf("Token: Expected 4, got %v", expNode2RightToken)
	}
	if expNode2RightValue != 4 {
		t.Errorf("Value: Expected 4, got %v", expNode2RightValue)
	}

	gotExpNode3 := gotNode.Value[2].(ast.ExpressionNode)
	gotExpNode3Token := gotExpNode3.Token
	if gotExpNode3Token != "true || false" {
		t.Errorf("Token: Expected true || false, got %v", gotExpNode3Token)
	}
	expNode3LeftToken := gotExpNode3.Left.(ast.BooleanNode).Token
	expNode3LeftValue := gotExpNode3.Left.(ast.BooleanNode).Value
	expNode3RightToken := gotExpNode3.Right.(ast.BooleanNode).Token
	expNode3RightValue := gotExpNode3.Right.(ast.BooleanNode).Value
	if expNode3LeftToken != "true" {
		t.Errorf("Token: Expected true, got %v", expNode3LeftToken)
	}
	if expNode3LeftValue != true {
		t.Errorf("Value: Expected true, got %v", expNode3LeftValue)
	}
	if expNode3RightToken != "false" {
		t.Errorf("Token: Expected false, got %v", expNode3RightToken)
	}
	if expNode3RightValue != false {
		t.Errorf("Value: Expected false, got %v", expNode3RightValue)
	}

	got, _ = ParseReader("", strings.NewReader("[1 + 2, 3 + 4, true || false]"))
	gotNode = got.(ast.ArrayNode)

	got, _ = ParseReader("", strings.NewReader("[null, ['A', 'B']]"))
	gotNode = got.(ast.ArrayNode)
	gotNode1 := gotNode.Value[0].(ast.NullNode)
	gotNode1Token := gotNode1.Token
	if gotNode1Token != "null" {
		t.Errorf("Token: Expected null, got %v", gotExpNode1Token)
	}

	gotNode2 := gotNode.Value[1].(ast.ArrayNode)
	gotNode2Token := gotNode2.Token
	if gotNode2Token != "['A', 'B']" {
		t.Errorf("Token: Expected ['A', 'B'], got %v", gotExpNode2Token)
	}
	gotNode2LeftToken := gotNode2.Value[0].(ast.StringNode).Token
	gotNode2LeftValue := gotNode2.Value[0].(ast.StringNode).Value
	gotNode2RightToken := gotNode2.Value[1].(ast.StringNode).Token
	gotNode2RightValue := gotNode2.Value[1].(ast.StringNode).Value
	if gotNode2LeftToken != "'A'" {
		t.Errorf("Token: Expected 'A', got %v", gotNode2LeftToken)
	}
	if gotNode2LeftValue != "A" {
		t.Errorf("Value: Expected A, got %v", gotNode2LeftValue)
	}
	if gotNode2RightToken != "'B'" {
		t.Errorf("Token: Expected 'B', got %v", gotNode2RightToken)
	}
	if gotNode2RightValue != "B" {
		t.Errorf("Value: Expected B, got %v", gotNode2RightValue)
	}
}
