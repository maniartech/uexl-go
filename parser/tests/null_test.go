package parser_test

import (
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/ast"
	. "github.com/maniartech/uexl_go/parser"
)

func TestNull(t *testing.T) {
	got, _ := ParseReaderNew("", strings.NewReader("null"))

	// Change the type assertion to handle a pointer to NullNode
	gotNode, ok := got.(*ast.NullNode)
	if !ok {
		t.Fatalf("Expected *ast.NullNode, got %T", got)
	}

	token := gotNode.Token
	if token != "null" {
		t.Errorf("Expected null, got %v", token)
	}
}
