package parser

import (
	"testing"
)

func TestBooleans(t *testing.T) {
	// got, _ := ParseReader("", strings.NewReader("true"))
	// gotNode := got.(ast.BooleanNode)
	// nodeVal := gotNode.Value
	// nodeToken := gotNode.Token
	// if nodeToken != "true" {
	// 	t.Errorf("Token: Expected true, got, %v", nodeToken)
	// }
	// if nodeVal != true {
	// 	t.Errorf("Value: Expected true, got, %v", nodeVal)
	// }

	// got, _ = ParseReader("", strings.NewReader("false"))
	// gotNode = got.(ast.BooleanNode)
	// nodeVal = gotNode.Value
	// nodeToken = gotNode.Token
	// if nodeToken != "false" {
	// 	t.Errorf("Token: Expected false, got, %v", nodeToken)
	// }
	// if nodeVal != false {
	// 	t.Errorf("Value: Expected false, got, %v", nodeVal)
	// }

	// got, _ = ParseReader("", strings.NewReader("true || false"))
	// gotExpNode := got.(ast.ExpressionNode)
	// nodeExpToken := gotExpNode.Token
	// if nodeExpToken != "true || false" {
	// 	t.Errorf("Token: Expected true || false, got, %v", nodeToken)
	// }
	// leftNodeVal := gotExpNode.Left.(ast.BooleanNode).Value
	// leftNodeToken := gotExpNode.Left.(ast.BooleanNode).Token
	// rightNodeVal := gotExpNode.Right.(ast.BooleanNode).Value
	// rightNodeToken := gotExpNode.Right.(ast.BooleanNode).Token
	// if leftNodeToken != "true" {
	// 	t.Errorf("Token: Expected true, got %v", leftNodeToken)
	// }
	// if leftNodeVal != true {
	// 	t.Errorf("Value: Expected true, got %v", leftNodeVal)
	// }
	// if rightNodeToken != "false" {
	// 	t.Errorf("Token: Expected false, got %v", rightNodeToken)
	// }
	// if rightNodeVal != false {
	// 	t.Errorf("Value: Expected false, got %v", rightNodeVal)
	// }

	// got, _ = ParseReader("", strings.NewReader("true && (true || false)"))
	// gotExpNode = got.(ast.ExpressionNode)
	// leftNodeToken = gotExpNode.Left.(ast.BooleanNode).Token
	// leftNodeValue := gotExpNode.Left.(ast.BooleanNode).Value
	// if leftNodeToken != "true" {
	// 	t.Errorf("Token: Expected true, got %v", leftNodeToken)
	// }
	// if leftNodeValue != true {
	// 	t.Errorf("Token: Expected true, got %v", leftNodeValue)
	// }
	// rightNodeToken = gotExpNode.Right.(ast.ExpressionNode).Token
	// if rightNodeToken != "true || false" {
	// 	t.Errorf("Token: Expceted true || false, got %v", rightNodeToken)
	// }
	// rightLeftToken := gotExpNode.Right.(ast.ExpressionNode).Left.(ast.BooleanNode).Token
	// rightLeftNode := gotExpNode.Right.(ast.ExpressionNode).Left.(ast.BooleanNode).Value
	// rightRightToken := gotExpNode.Right.(ast.ExpressionNode).Right.(ast.BooleanNode).Token
	// rightRightNode := gotExpNode.Right.(ast.ExpressionNode).Right.(ast.BooleanNode).Value
	// if rightLeftToken != "true" {
	// 	t.Errorf("Token: Expected true, got %v", rightLeftToken)
	// }
	// if rightRightToken != "false" {
	// 	t.Errorf("Token: Expected false, got %v", rightRightToken)
	// }
	// if rightLeftNode != true {
	// 	t.Errorf("Value: Expected true, got %v", rightLeftNode)
	// }
	// if rightRightNode != false {
	// 	t.Errorf("Value: Expected false, got %v", rightRightNode)
	// }

	// got, _ = ParseReader("", strings.NewReader("(true || true) && (true || false)"))
	// gotExpNode = got.(ast.ExpressionNode)
	// leftNodeToken = gotExpNode.Left.(ast.ExpressionNode).Token
	// if leftNodeToken != "true || true" {
	// 	t.Errorf("Token: Expected true || true, got %v", leftNodeToken)
	// }
	// leftLeftToken := gotExpNode.Left.(ast.ExpressionNode).Left.(ast.BooleanNode).Token
	// leftLeftNode := gotExpNode.Left.(ast.ExpressionNode).Left.(ast.BooleanNode).Value
	// leftRightToken := gotExpNode.Left.(ast.ExpressionNode).Right.(ast.BooleanNode).Token
	// leftRightNode := gotExpNode.Left.(ast.ExpressionNode).Right.(ast.BooleanNode).Value
	// if leftLeftToken != "true" {
	// 	t.Errorf("Token: Expected true, got %v", rightLeftToken)
	// }
	// if leftRightToken != "true" {
	// 	t.Errorf("Token: Expected true, got %v", rightRightToken)
	// }
	// if leftLeftNode != true {
	// 	t.Errorf("Value: Expected true, got %v", rightLeftNode)
	// }
	// if leftRightNode != true {
	// 	t.Errorf("Value: Expected true, got %v", rightRightNode)
	// }
	// rightNodeToken = gotExpNode.Right.(ast.ExpressionNode).Token
	// if rightNodeToken != "true || false" {
	// 	t.Errorf("Token: Expceted true || false, got %v", rightNodeToken)
	// }
	// rightLeftToken = gotExpNode.Right.(ast.ExpressionNode).Left.(ast.BooleanNode).Token
	// rightLeftNode = gotExpNode.Right.(ast.ExpressionNode).Left.(ast.BooleanNode).Value
	// rightRightToken = gotExpNode.Right.(ast.ExpressionNode).Right.(ast.BooleanNode).Token
	// rightRightNode = gotExpNode.Right.(ast.ExpressionNode).Right.(ast.BooleanNode).Value
	// if rightLeftToken != "true" {
	// 	t.Errorf("Token: Expected true, got %v", rightLeftToken)
	// }
	// if rightRightToken != "false" {
	// 	t.Errorf("Token: Expected false, got %v", rightRightToken)
	// }
	// if rightLeftNode != true {
	// 	t.Errorf("Value: Expected true, got %v", rightLeftNode)
	// }
	// if rightRightNode != false {
	// 	t.Errorf("Value: Expected false, got %v", rightRightNode)
	// }
}
