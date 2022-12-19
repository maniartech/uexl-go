package parser

import (
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/ast"
)

func TestPipe(t *testing.T) {
	//fmt.Println(gotNode.Expressions[0].(ast.ArrayNode).Value[0].(ast.NumberNode).Value)
	got, _ := ParseReader("", strings.NewReader("[1, 2, 3]|map:(1 + 2)"))
	gotNode := got.(ast.PipeNode)

	gotNodeType := gotNode.Type
	if gotNodeType != "pipe" {
		t.Errorf("Token: Expected pipe, got %v", gotNodeType)
	}

	gotNodeToken := gotNode.Token
	if gotNodeToken != "[1, 2, 3]|map:(1 + 2)" {
		t.Errorf("Token: Expected [1, 2, 3]|map:(1 + 2), got %v", gotNodeToken)
	}

	gotNodePipeType := gotNode.PipeType
	if gotNodePipeType != "map" {
		t.Errorf("Pipe Type: Expected map, got %v", gotNodePipeType)
	}

	gotNodeExpressions1 := gotNode.Expressions[0].(ast.ArrayNode)
	gotNodeExpressionsToken1 := gotNodeExpressions1.Token
	if gotNodeExpressionsToken1 != "[1, 2, 3]" {
		t.Errorf("Token: Expected [1, 2, 3], got %v", gotNodeExpressionsToken1)
	}

	gotNodeExpressions1_1 := gotNodeExpressions1.Value[0].(ast.NumberNode).Token
	gotNodeExpressions1_1Value := gotNodeExpressions1.Value[0].(ast.NumberNode).Value
	if gotNodeExpressions1_1 != "1" {
		t.Errorf("Token: Expected 1, got %v", gotNodeExpressions1_1)
	}
	if gotNodeExpressions1_1Value != 1 {
		t.Errorf("Value: Expected 1, got %v", gotNodeExpressions1_1Value)
	}
	gotNodeExpressions1_2 := gotNodeExpressions1.Value[1].(ast.NumberNode).Token
	gotNodeExpressions1_2Value := gotNodeExpressions1.Value[1].(ast.NumberNode).Value
	if gotNodeExpressions1_2 != "2" {
		t.Errorf("Token: Expected 2, got %v", gotNodeExpressions1_1)
	}
	if gotNodeExpressions1_2Value != 2 {
		t.Errorf("Value: Expected 2, got %v", gotNodeExpressions1_1Value)
	}

	gotNodeExpressions1_3 := gotNodeExpressions1.Value[2].(ast.NumberNode).Token
	gotNodeExpressions1_3Value := gotNodeExpressions1.Value[2].(ast.NumberNode).Value
	if gotNodeExpressions1_3 != "3" {
		t.Errorf("Token: Expected 3, got %v", gotNodeExpressions1_1)
	}
	if gotNodeExpressions1_3Value != 3 {
		t.Errorf("Value: Expected 3, got %v", gotNodeExpressions1_1Value)
	}

	gotNodeExpressions2 := gotNode.Expressions[1].(ast.ExpressionNode)
	gotNodeExpressions2Token := gotNodeExpressions2.Token
	op := gotNodeExpressions2.Operator
	if op != "+" {
		t.Errorf("Operator: Expected +, got %v", op)
	}
	if gotNodeExpressions2Token != "1 + 2" {
		t.Errorf("Token: Expected 1 + 2, got %v", gotNodeExpressions2Token)
	}

	gotNodeExpressions2_1 := gotNodeExpressions2.Left.(ast.NumberNode).Token
	gotNodeExpressions2_1Value := gotNodeExpressions2.Left.(ast.NumberNode).Value
	if gotNodeExpressions2_1 != "1" {
		t.Errorf("Token: Expected 1, got %v", gotNodeExpressions2_1)
	}
	if gotNodeExpressions2_1Value != 1 {
		t.Errorf("Value: Expected 1, got %v", gotNodeExpressions2_1Value)
	}

	gotNodeExpressions2_2 := gotNodeExpressions2.Right.(ast.NumberNode).Token
	gotNodeExpressions2_2Value := gotNodeExpressions2.Right.(ast.NumberNode).Value
	if gotNodeExpressions2_2 != "2" {
		t.Errorf("Token: Expected 2, got %v", gotNodeExpressions2_2)
	}
	if gotNodeExpressions2_2Value != 2 {
		t.Errorf("Value: Expected 2, got %v", gotNodeExpressions2_2Value)
	}

	got, _ = ParseReader("", strings.NewReader("(25 * 4)|:'sum'"))
	gotNode = got.(ast.PipeNode)
	gotNodeType = gotNode.Type
	if gotNodeType != "pipe" {
		t.Errorf("Token: Expected pipe, got %v", gotNodeType)
	}

	gotNodeToken = gotNode.Token
	if gotNodeToken != "(25 * 4)|:'sum'" {
		t.Errorf("Token: Expected (25 * 4)|:'sum', got %v", gotNodeToken)
	}

	gotNodePipeType = gotNode.PipeType
	if gotNodePipeType != "pipe" {
		t.Errorf("Pipe Type: Expected pipe, got %v", gotNodePipeType)
	}

	gotNodeExpressions1Exp := gotNode.Expressions[0].(ast.ExpressionNode)
	gotNodeExpressionsToken1 = gotNodeExpressions1Exp.Token
	op = gotNodeExpressions1Exp.Operator
	if op != "*" {
		t.Errorf("Operator: Expected *, got %v", op)
	}
	if gotNodeExpressionsToken1 != "25 * 4" {
		t.Errorf("Token: Expected 25 * 4, got %v", gotNodeExpressionsToken1)
	}

	num1 := gotNodeExpressions1Exp.Left.(ast.NumberNode).Value
	num1Token := gotNodeExpressions1Exp.Left.(ast.NumberNode).Token
	if num1Token != "25" {
		t.Errorf("Token: Expected 25, got %v", num1Token)
	}
	if num1 != 25 {
		t.Errorf("Value: Expected 25, got %v", num1)
	}
	num2 := gotNodeExpressions1Exp.Right.(ast.NumberNode).Value
	num2Token := gotNodeExpressions1Exp.Right.(ast.NumberNode).Token
	if num2Token != "4" {
		t.Errorf("Token: Expected 4, got %v", num1Token)
	}
	if num2 != 4 {
		t.Errorf("Value: Expected 4, got %v", num1)
	}

	got, _ = ParseReader("", strings.NewReader("{'name': 'abc', 'age': 30}|filter:true"))
	gotNode = got.(ast.PipeNode)
	gotNodeType = gotNode.Type
	if gotNodeType != "pipe" {
		t.Errorf("Token: Expected pipe, got %v", gotNodeType)
	}

	gotNodeToken = gotNode.Token
	if gotNodeToken != "{'name': 'abc', 'age': 30}|filter:true" {
		t.Errorf("Token: Expected {'name': 'abc', 'age': 30}|filter:true, got %v", gotNodeToken)
	}

	gotNodePipeType = gotNode.PipeType
	if gotNodePipeType != "filter" {
		t.Errorf("Pipe Type: Expected filter, got %v", gotNodePipeType)
	}

	gotNodeExpressions1Obj := gotNode.Expressions[0].(ast.ObjectNode)
	gotNodeExpressionsToken1 = gotNodeExpressions1Obj.Token
	if gotNodeExpressionsToken1 != "{'name': 'abc', 'age': 30}" {
		t.Errorf("Token: Expected {'name': 'abc', 'age': 30}, got %v", gotNodeExpressionsToken1)
	}

	keys := gotNodeExpressions1Obj.Items
	key1 := keys[0].Key
	if key1 != "name" {
		t.Errorf("Key: Expected name, got %v", key1)
	}
	value1 := keys[0].Value.(ast.StringNode).Value
	if value1 != "abc" {
		t.Errorf("Value: Expected abc, got %v", value1)
	}

	key2 := keys[1].Key
	if key2 != "age" {
		t.Errorf("Key: Expected age, got %v", key2)
	}
	value2 := keys[1].Value.(ast.NumberNode).Value
	if value2 != 30 {
		t.Errorf("Value: Expected 30, got %v", value2)
	}

	gotNodeExpressions1Bool := gotNode.Expressions[1].(ast.BooleanNode)
	gotNodeExpressionsToken1 = gotNodeExpressions1Bool.Token
	gotNodeExpressions1BoolVal := gotNodeExpressions1Bool.Value
	if gotNodeExpressionsToken1 != "true" {
		t.Errorf("Token: Expected true, got %v", gotNodeExpressionsToken1)
	}
	if gotNodeExpressions1BoolVal != true {
		t.Errorf("Token: Expected true, got %v", gotNodeExpressions1BoolVal)
	}
}
