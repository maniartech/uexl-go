package vm_test

import (
	"fmt"
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/vm"
)

// parse is a helper function to parse the input string
// this runs the lexer and parser to create an AST
func parse(input string) parser.Node {
	p := parser.NewParser(input)
	node, err := p.Parse()
	if err != nil {
		fmt.Printf("Parse error: %s\n", err)
		return nil
	}
	if node == nil {
		fmt.Println("Parse error: no node returned")
		return nil
	}
	return node
}

type vmTestCases struct {
	input    string
	expected any
}

func TestVM(t *testing.T) {
	vm := vm.New()
	if vm == nil {
		t.Fatal("failed to create VM")
	}
}
