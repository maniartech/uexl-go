package vm

import (
	"fmt"

	"github.com/maniartech/uexl_go/compiler"
	"github.com/maniartech/uexl_go/parser"
)

func DefaultPipeHandlers() PipeHandlers {
	return PipeHandlers{
		"map":  MapPipeHandler,
		"pipe": DefaultPipeHandler, // Default pipe handler
		"filter": FilterPipeHandler,
	}
}

func DefaultPipeHandler(input parser.Node, block any, alias string, vm *VM) (parser.Node, error) {
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		// Pass-through if no block
		return input, nil
	}
	vm.pushPipeScope()
	vm.setPipeVar("$last", input)
	// setting all the system variables as pipe variables
	for i, v := range vm.aliasVars {
		vm.setPipeVar(i, v)
	}
	frame := NewFrame(blk.Instructions, 0)
	vm.pushFrame(frame)
	err := vm.Run()
	if err != nil {
		vm.popPipeScope()
		vm.popFrame()
		return nil, err
	}
	res := vm.Pop()
	vm.popFrame()
	vm.popPipeScope()
	return res, nil
}

func MapPipeHandler(input parser.Node, block any, alias string, vm *VM) (parser.Node, error) {
	arr, ok := input.(*parser.ArrayLiteral)
	if !ok {
		return nil, fmt.Errorf("map pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("map pipe expects a predicate block")
	}

	result := make([]parser.Node, len(arr.Elements))
	for i, elem := range arr.Elements {
		vm.pushPipeScope()
		if alias != "" {
			vm.setPipeVar(alias, elem)
		}
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", &parser.NumberLiteral{Value: float64(i)})
		for i, v := range vm.aliasVars {
			vm.setPipeVar(i, v)
		}
		frame := NewFrame(blk.Instructions, 0)
		vm.pushFrame(frame)
		err := vm.Run()
		if err != nil {
			vm.popPipeScope()
			vm.popFrame()
			return nil, err
		}
		res := vm.Pop()
		vm.popFrame()
		vm.popPipeScope()
		result[i] = res
	}

	exprs := make([]parser.Expression, len(result))
	for i, n := range result {
		exprs[i] = n.(parser.Expression)
	}

	return &parser.ArrayLiteral{Elements: exprs}, nil
}

func FilterPipeHandler(input parser.Node, block any, alias string, vm *VM) (parser.Node, error) {
	arr, ok := input.(*parser.ArrayLiteral)
	if !ok {
		return nil, fmt.Errorf("filter pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("filter pipe expects a predicate block")
	}

	var result []parser.Node
	for i, elem := range arr.Elements {
		vm.pushPipeScope()
		if alias != "" {
			vm.setPipeVar(alias, elem)
		}
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", &parser.NumberLiteral{Value: float64(i)})
		for j, v := range vm.aliasVars {
			vm.setPipeVar(j, v)
		}
		frame := NewFrame(blk.Instructions, 0)
		vm.pushFrame(frame)
		err := vm.Run()
		if err != nil {
			vm.popPipeScope()
			vm.popFrame()
			return nil, err
		}
		res := vm.Pop()
		vm.popFrame()
		vm.popPipeScope()

		// Only include element if predicate is true
		if boolLit, ok := res.(*parser.BooleanLiteral); ok && boolLit.Value {
			result = append(result, elem)
		}
	}

	exprs := make([]parser.Expression, len(result))
	for i, n := range result {
		exprs[i] = n.(parser.Expression)
	}

	return &parser.ArrayLiteral{Elements: exprs}, nil
}
