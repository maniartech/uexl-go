package parser

import (
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/ast"
)

func TestObject(t *testing.T) {
	// Test case for {'abc': 'xyz'}
	got, _ := ParseReader("", strings.NewReader("{'abc': 'xyz'}"))
	gotNode := got.(ast.ObjectNode)
	gotToken := gotNode.Token
	if gotToken != "{'abc': 'xyz'}" {
		t.Errorf("Token: Expected {'abc': 'xyz'}, got %v", gotToken)
	}
	key := gotNode.Items[0].Key
	value := gotNode.Items[0].Value.(ast.StringNode).Value
	if key != "abc" {
		t.Errorf("Key: Expected abc, got %v", key)
	}
	if value != "xyz" {
		t.Errorf("Value: Expected xzy, got %v", value)
	}

	// Test case for {'data': 30}
	got, _ = ParseReader("", strings.NewReader("{'data': 30}"))
	gotNode = got.(ast.ObjectNode)
	gotToken = gotNode.Token
	if gotToken != "{'data': 30}" {
		t.Errorf("Token: Expected {'data': 30}, got %v", gotToken)
	}
	key = gotNode.Items[0].Key
	valueNum := gotNode.Items[0].Value.(ast.NumberNode).Value
	if key != "data" {
		t.Errorf("Key: Expected data, got %v", key)
	}
	if valueNum != 30 {
		t.Errorf("Value: Expected 30, got %v", value)
	}

	// Test case for {'type': 'books', 'amount': 100}
	got, _ = ParseReader("", strings.NewReader("{'type': 'books', 'amount': 100}"))
	gotNode = got.(ast.ObjectNode)
	gotToken = gotNode.Token
	if gotToken != "{'type': 'books', 'amount': 100}" {
		t.Errorf("Token: Expected {'type': 'books', 'amount': 100}, got %v", gotToken)
	}
	key = gotNode.Items[0].Key
	value = gotNode.Items[0].Value.(ast.StringNode).Value
	if key != "type" {
		t.Errorf("Key: Expected type, got %v", key)
	}
	if value != "books" {
		t.Errorf("Value: Expected books, got %v", value)
	}
	key = gotNode.Items[1].Key
	valueNum = gotNode.Items[1].Value.(ast.NumberNode).Value
	if key != "amount" {
		t.Errorf("Key: Expected amount, got %v", key)
	}
	if valueNum != 100 {
		t.Errorf("Value: Expected 100, got %v", value)
	}

	got, _ = ParseReader("", strings.NewReader("{'nums': [3.14, 9.8, 2.71]}"))
	gotNode = got.(ast.ObjectNode)
	gotToken = gotNode.Token
	if gotToken != "{'nums': [3.14, 9.8, 2.71]}" {
		t.Errorf("Token: Expected {'nums': [3.14, 9.8, 2.71]}, got %v", gotToken)
	}
	key = gotNode.Items[0].Key
	if key != "nums" {
		t.Errorf("Key: Expected nums, got %v", key)
	}
	val1 := gotNode.Items[0].Value.(ast.ArrayNode).Value[0].(ast.NumberNode).Value
	if val1 != 3.14 {
		t.Errorf("Value: Expected 3.14, got %v", val1)
	}
	val2 := gotNode.Items[0].Value.(ast.ArrayNode).Value[1].(ast.NumberNode).Value
	if val2 != 9.8 {
		t.Errorf("Value: Expected 9.8, got %v", val2)
	}
	val3 := gotNode.Items[0].Value.(ast.ArrayNode).Value[2].(ast.NumberNode).Value
	if val3 != 2.71 {
		t.Errorf("Value: Expected 2.71, got %v", val3)
	}

	// Test case for {'details': {'name': 'abc', 'age': 30}}
	got, _ = ParseReader("", strings.NewReader("{'details': {'name': 'abc', 'age': 30}}"))
	gotNode = got.(ast.ObjectNode)
	gotToken = gotNode.Token
	if gotToken != "{'details': {'name': 'abc', 'age': 30}}" {
		t.Errorf("Token: Expected {'details': {'name': 'abc', 'age': 30}}, got %v", gotToken)
	}
	node := gotNode.Items[0]
	key = node.Key
	if key != "details" {
		t.Errorf("Key: Expected details, got %v", key)
	}
	valuesToken := node.Value.(ast.ObjectNode).Token
	if valuesToken != "{'name': 'abc', 'age': 30}" {
		t.Errorf("Token: Expected {'name': 'abc', 'age': 30}, got %v", valuesToken)
	}
	valueKey1 := node.Value.(ast.ObjectNode).Items[0].Key
	valueValue1 := node.Value.(ast.ObjectNode).Items[0].Value.(ast.StringNode).Value
	if valueKey1 != "name" {
		t.Errorf("Key: Expected name, got %v", valueKey1)
	}
	if valueValue1 != "abc" {
		t.Errorf("Value: Expected abc, got %v", valueValue1)
	}
	valueKey2 := node.Value.(ast.ObjectNode).Items[1].Key
	valueValue2Num := node.Value.(ast.ObjectNode).Items[1].Value.(ast.NumberNode).Value
	if valueKey2 != "age" {
		t.Errorf("Key: Expected age, got %v", valueKey2)
	}
	if valueValue2Num != 30 {
		t.Errorf("Value: Expected 30, got %v", valueValue1)
	}
}
