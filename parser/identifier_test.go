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
	if name != "abc.xyz" {
		t.Errorf("Name: Expected abc.xyz, got %v", name)
	}

	got, _ = ParseReader("", strings.NewReader("aBc0.XyZ9"))
	gotNode = got.(ast.IdentifierNode)
	token = gotNode.Token
	if token != "aBc0.XyZ9" {
		t.Errorf("Token: Expected aBc0.XyZ9, got %v", token)
	}
	name = gotNode.Name
	if name != "aBc0.XyZ9" {
		t.Errorf("Name: Expected aBc0.XyZ9, got %v", name)
	}

	got, _ = ParseReader("", strings.NewReader("@employee_name.$name"))
	gotNode = got.(ast.IdentifierNode)
	token = gotNode.Token
	if token != "@employee_name.$name" {
		t.Errorf("Token: Expected @employee_name.$name, got %v", token)
	}
	name = gotNode.Name
	if name != "@employee_name.$name" {
		t.Errorf("Name: Expected aBc0.XyZ9, got %v", name)
	}
}
