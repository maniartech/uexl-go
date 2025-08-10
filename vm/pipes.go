package vm

import (
	"fmt"

	"github.com/maniartech/uexl_go/parser"
)

func DefaultPipeHandlers() PipeHandlers {
	return PipeHandlers{
		"map":  MapPipeHandler,
		"pipe": DefaultPipeHandler, // Default pipe handler
		// "filter": FilterPipeHandler,
	}
}

func DefaultPipeHandler(input parser.Node, lambda parser.Node, alias string, vm *VM) (parser.Node, error) {
	// in the default pipe type the lambda is the input itself, but we will check if the input is nil if yes then we return the lambda
	if input == nil {
		if lambda == nil {
			return nil, fmt.Errorf("default pipe requires input or lambda")
		}
		return lambda, nil
	}
	return input, nil
}

func MapPipeHandler(input parser.Node, lambda parser.Node, alias string, vm *VM) (parser.Node, error) {
	arr, ok := input.(*parser.ArrayLiteral)
	if !ok {
		return nil, fmt.Errorf("map pipe expects array input")
	}

	result := make([]parser.Node, len(arr.Elements))
	for i, elem := range arr.Elements {
		vm.pushPipeScope()
		if alias != "" {
			vm.setPipeVar(alias, elem)
		}
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", &parser.NumberLiteral{Value: float64(i)})

		// Execute the lambda bytecode (already on stack or in instruction stream)
		// The lambda instructions will run and leave result on stack
		res := vm.Pop() // Get the result that lambda instructions produced

		vm.popPipeScope()
		result[i] = res
	}

	exprs := make([]parser.Expression, len(result))
	for i, n := range result {
		exprs[i] = n.(parser.Expression)
	}

	return &parser.ArrayLiteral{Elements: exprs}, nil
}
