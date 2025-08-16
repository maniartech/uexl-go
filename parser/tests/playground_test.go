package parser_test

import (
	"fmt"
	"testing"

	"github.com/maniartech/uexl_go/compiler"
	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/vm"
)

func TestPlayground(t *testing.T) {
	input := "[0,2,3,4] |reduce: set($acc || {}, $index, $item)"

	parserNode, err := parser.ParseString(input)
	if err != nil {
		t.Fatalf("Failed to parse expression: %v", err)
	}

	code := compiler.New()
	err = code.Compile(parserNode)
	if err != nil {
		t.Fatalf("Failed to compile expression: %v", err)
	}

	virtualMachine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	output, err := virtualMachine.Run(code.ByteCode())
	if err != nil {
		t.Fatalf("Failed to create VM: %v", err)
	}

	fmt.Println(output)

	// utils.PrintJSON(expr)
}
