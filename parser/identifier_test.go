package parser

import (
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/ast"
)

func TestIdentifiers(t *testing.T) {
	got, _ := ParseReader("", strings.NewReader("abc.xyz"))
	gotNode := got.(ast.IdentifierNode)
	token := gotNode.Token
	if token != "abc.xyz" {
		t.Errorf("Token: Expected abc.xyz, got %v", token)
	}
	name := gotNode.Name
	if name != "abc" {
		t.Errorf("Name: Expected abc, got %v", name)
	}

	value := gotNode.Value
	if value != "xyz" {
		t.Errorf("Value: Expected xyz, got %v", value)
	}

	got, _ = ParseReader("", strings.NewReader("aBc0.XyZ9"))
	gotNode = got.(ast.IdentifierNode)
	token = gotNode.Token
	if token != "aBc0.XyZ9" {
		t.Errorf("Token: Expected aBc0.XyZ9, got %v", token)
	}
	name = gotNode.Name
	if name != "aBc0" {
		t.Errorf("Name: Expected aBc0, got %v", name)
	}

	value = gotNode.Value
	if value != "XyZ9" {
		t.Errorf("Value: Expected XyZ9, got %v", value)
	}
}
