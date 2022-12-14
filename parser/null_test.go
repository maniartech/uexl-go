package parser

import (
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/ast"
)

func TestNull(t *testing.T) {
	got, _ := ParseReader("", strings.NewReader("null"))
	gotNode := got.(ast.NullNode)
	token := gotNode.Token
	if token != "null" {
		t.Errorf("Expcted null, got %v", token)
	}
}
