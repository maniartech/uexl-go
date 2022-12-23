package ast

import (
	"fmt"
	"testing"
)

func TestObject(t *testing.T) {
	nameNode, _ := NewStringNode("abc", 0, 1, 1)
	ageNode, _ := NewNumberNode("20", 0, 1, 1)
	objNodes := []ObjectItem{
		{String("name"), nameNode},
		{String("age"), ageNode},
	}

	node, _ := NewObjectNode("{'name': 'abc', 'age': 20}", objNodes, 0, 1, 1)
	nodeObj := node.(ObjectNode)

	nodeToken := nodeObj.Token
	if nodeToken != "{'name': 'abc', 'age': 20}" {
		t.Errorf("Token: Expected {'name': 'abc', 'age': 20}, got %v", nodeToken)
	}

	nodeType := nodeObj.Type
	if nodeType != "object" {
		t.Errorf("Type: Expected object, got %v", nodeType)
	}

	nodeLine := nodeObj.Line
	if nodeLine != 1 {
		t.Errorf("Line: Expected 1, got %v", nodeLine)
	}

	nodeCol := nodeObj.Column
	if nodeCol != 1 {
		t.Errorf("Column: Expected 1, got %v", nodeCol)
	}

	nodeOff := nodeObj.Offset
	if nodeOff != 0 {
		t.Errorf("Offset: Expected 0, got %v", nodeOff)
	}

	var tmp Map
	nodeEval, err := nodeObj.Eval(tmp)
	expectedVal1 := nodeEval.(map[String]interface{})["name"].(String)
	if expectedVal1 != "abc" {
		t.Errorf("Value 1: Expected abc, got %v", expectedVal1)
	}
	expectedVal2 := nodeEval.(map[String]interface{})["age"].(Number)
	if expectedVal2 != 20 {
		t.Errorf("Value 2: Expected 30, got %v", expectedVal1)
	}

	nodeObjStringer := fmt.Sprintf("%v", node)
	if nodeObjStringer != "ObjectNode: {name: StringNode abc, age: NumberNode 20}" {
		t.Errorf("Stringer: Expected ObjectNode: {name: StringNode abc, age: NumberNode 20}, got %v", nodeObjStringer)
	}

	if err != nil {
		t.Errorf("Error: Expected nil, got %v", err)
	}
}
