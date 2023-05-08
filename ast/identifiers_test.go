package ast_test

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/types"
)

func TestIdentifier(t *testing.T) {
	// Perform Identifier tests here
	node, _ := parser.ParseString("AVERAGE(10, 20) + x - y.z")

	result, err := node.Eval(types.Context{
		"x": types.Number(10),
		"y": types.Object{
			"z": types.Number(5),
		},
	})

	if err != nil {
		t.Errorf("Identifier.Eval() error = %v", err)
	}

	if result != types.Number(20) {
		t.Errorf("Identifier.Eval() = %v", result)
	}

}
