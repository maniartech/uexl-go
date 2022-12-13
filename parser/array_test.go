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
}
