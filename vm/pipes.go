package vm

import (
	"fmt"

	"github.com/maniartech/uexl_go/compiler"
	"github.com/maniartech/uexl_go/parser"
)

func DefaultPipeHandlers() PipeHandlers {
	return PipeHandlers{
		"map":    MapPipeHandler,
		"pipe":   DefaultPipeHandler,
		"filter": FilterPipeHandler,
		"reduce": ReducePipeHandler,
		"find":   FindPipeHandler,
		"some":   SomePipeHandler,
		"every":  EveryPipeHandler,
		"unique": UniquePipeHandler,
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

func ReducePipeHandler(input parser.Node, block any, alias string, vm *VM) (parser.Node, error) {
	arr, ok := input.(*parser.ArrayLiteral)
	if !ok {
		return nil, fmt.Errorf("reduce pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("reduce pipe expects a predicate block")
	}
	if len(arr.Elements) == 0 {
		return nil, fmt.Errorf("reduce pipe cannot operate on empty array")
	}

	acc := arr.Elements[0]
	for i := 1; i < len(arr.Elements); i++ {
		vm.pushPipeScope()
		vm.setPipeVar("$acc", acc)
		vm.setPipeVar("$item", arr.Elements[i])
		vm.setPipeVar("$index", &parser.NumberLiteral{Value: float64(i)})
		frame := NewFrame(blk.Instructions, 0)
		vm.pushFrame(frame)
		err := vm.Run()
		if err != nil {
			vm.popPipeScope()
			vm.popFrame()
			return nil, err
		}
		acc = vm.Pop().(parser.Expression)
		vm.popFrame()
		vm.popPipeScope()
	}
	return acc, nil
}

func FindPipeHandler(input parser.Node, block any, alias string, vm *VM) (parser.Node, error) {
	arr, ok := input.(*parser.ArrayLiteral)
	if !ok {
		return nil, fmt.Errorf("find pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("find pipe expects a predicate block")
	}

	for i, elem := range arr.Elements {
		vm.pushPipeScope()
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", &parser.NumberLiteral{Value: float64(i)})
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
		if boolLit, ok := res.(*parser.BooleanLiteral); ok && boolLit.Value {
			return elem, nil
		}
	}
	return &parser.NullLiteral{}, nil
}

func SomePipeHandler(input parser.Node, block any, alias string, vm *VM) (parser.Node, error) {
	arr, ok := input.(*parser.ArrayLiteral)
	if !ok {
		return nil, fmt.Errorf("some pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("some pipe expects a predicate block")
	}

	for i, elem := range arr.Elements {
		vm.pushPipeScope()
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", &parser.NumberLiteral{Value: float64(i)})
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
		if boolLit, ok := res.(*parser.BooleanLiteral); ok && boolLit.Value {
			return &parser.BooleanLiteral{Value: true}, nil
		}
	}
	return &parser.BooleanLiteral{Value: false}, nil
}

func EveryPipeHandler(input parser.Node, block any, alias string, vm *VM) (parser.Node, error) {
	arr, ok := input.(*parser.ArrayLiteral)
	if !ok {
		return nil, fmt.Errorf("every pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("every pipe expects a predicate block")
	}

	for i, elem := range arr.Elements {
		vm.pushPipeScope()
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", &parser.NumberLiteral{Value: float64(i)})
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
		if boolLit, ok := res.(*parser.BooleanLiteral); !ok || !boolLit.Value {
			return &parser.BooleanLiteral{Value: false}, nil
		}
	}
	return &parser.BooleanLiteral{Value: true}, nil
}

func UniquePipeHandler(input parser.Node, block any, alias string, vm *VM) (parser.Node, error) {
	arr, ok := input.(*parser.ArrayLiteral)
	if !ok {
		return nil, fmt.Errorf("unique pipe expects array input")
	}
	seen := make(map[string]bool)
	var result []parser.Node
	for _, elem := range arr.Elements {
		key := fmt.Sprintf("%v", elem)
		if !seen[key] {
			seen[key] = true
			result = append(result, elem)
		}
	}
	exprs := make([]parser.Expression, len(result))
	for i, n := range result {
		exprs[i] = n.(parser.Expression)
	}
	return &parser.ArrayLiteral{Elements: exprs}, nil
}
