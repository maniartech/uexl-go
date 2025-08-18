package parser_test

import (
	"fmt"
	"testing"

	"github.com/maniartech/uexl_go/compiler"
	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/vm"
)

func TestPlayground(t *testing.T) {
	input := "[1, 2, 3, 4].2" // Example input to test
	// input := "[1,2,3,4] |reduce: ($acc || 0) +  $item"
	// input := "[1,2,3,4] |reduce: ($acc || '') + str($item) + ',' |: substr($last, 0, len($last)-1)"

	// Reusability Approaches

	// When Arg is passed from context
	// input := "arg |reduce: ($acc || '') + str($item) + ',' |: substr($last, 0, len($last)-1)"

	// When reusable expressions are defined
	// concatStr := "|reduce: ($acc || '') + str($item) + ',' |: substr($last, 0, len($last)-1)"
	// input := "[1,2, 3, 4, 5] " + concatStr

	// Dynamic Function Expressions
	// concatStr := "|reduce: ($acc || '') + str($item) + ',' |: substr($last, 0, len($last)-1)"
	// uexl.RegisterFunctionExpression("concatStr", concatStr)
	// input := "concatStr([1,2, 3, 4, 5])"

	// Dynamic Pipe Expressions, defines predicate expression that works value
	// ?

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
