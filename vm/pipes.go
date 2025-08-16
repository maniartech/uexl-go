package vm

import (
	"fmt"
	"sort"

	"github.com/maniartech/uexl_go/compiler"
)

var DefaultPipeHandlers = PipeHandlers{
	"map":     MapPipeHandler,
	"pipe":    DefaultPipeHandler,
	"filter":  FilterPipeHandler,
	"reduce":  ReducePipeHandler,
	"find":    FindPipeHandler,
	"some":    SomePipeHandler,
	"every":   EveryPipeHandler,
	"unique":  UniquePipeHandler,
	"sort":    SortPipeHandler,
	"groupBy": GroupByPipeHandler,
	"window":  WindowPipeHandler,
	"chunk":   ChunkPipeHandler,
}

func DefaultPipeHandler(input any, block any, alias string, vm *VM) (any, error) {
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return input, nil
	}
	vm.pushPipeScope()
	vm.setPipeVar("$last", input)
	frame := NewFrame(blk.Instructions, 0)
	vm.pushFrame(frame)
	err := vm.run()
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

func MapPipeHandler(input any, block any, alias string, vm *VM) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("map pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("map pipe expects a predicate block")
	}
	result := make([]any, len(arr))
	for i, elem := range arr {
		vm.pushPipeScope()
		if alias != "" {
			vm.setPipeVar(alias, elem)
		}
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", i)
		frame := NewFrame(blk.Instructions, 0)
		vm.pushFrame(frame)
		err := vm.run()
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
	return result, nil
}

func FilterPipeHandler(input any, block any, alias string, vm *VM) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("filter pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("filter pipe expects a predicate block")
	}
	var result []any
	for i, elem := range arr {
		vm.pushPipeScope()
		if alias != "" {
			vm.setPipeVar(alias, elem)
		}
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", i)
		frame := NewFrame(blk.Instructions, 0)
		vm.pushFrame(frame)
		err := vm.run()
		if err != nil {
			vm.popPipeScope()
			vm.popFrame()
			return nil, err
		}
		res := vm.Pop()
		vm.popFrame()
		vm.popPipeScope()
		if b, ok := res.(bool); ok && b {
			result = append(result, elem)
		}
	}
	return result, nil
}

func ReducePipeHandler(input any, block any, alias string, vm *VM) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("reduce pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("reduce pipe expects a predicate block")
	}
	if len(arr) == 0 {
		return nil, fmt.Errorf("reduce pipe cannot operate on empty array")
	}
	acc := arr[0]
	for i := 1; i < len(arr); i++ {
		vm.pushPipeScope()
		vm.setPipeVar("$acc", acc)
		vm.setPipeVar("$item", arr[i])
		vm.setPipeVar("$index", i)
		frame := NewFrame(blk.Instructions, 0)
		vm.pushFrame(frame)
		err := vm.run()
		if err != nil {
			vm.popPipeScope()
			vm.popFrame()
			return nil, err
		}
		acc = vm.Pop()
		vm.popFrame()
		vm.popPipeScope()
	}
	return acc, nil
}

func FindPipeHandler(input any, block any, alias string, vm *VM) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("find pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("find pipe expects a predicate block")
	}
	for i, elem := range arr {
		vm.pushPipeScope()
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", i)
		frame := NewFrame(blk.Instructions, 0)
		vm.pushFrame(frame)
		err := vm.run()
		if err != nil {
			vm.popPipeScope()
			vm.popFrame()
			return nil, err
		}
		res := vm.Pop()
		vm.popFrame()
		vm.popPipeScope()
		if b, ok := res.(bool); ok && b {
			return elem, nil
		}
	}
	return nil, nil
}

func SomePipeHandler(input any, block any, alias string, vm *VM) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("some pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("some pipe expects a predicate block")
	}
	for i, elem := range arr {
		vm.pushPipeScope()
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", i)
		frame := NewFrame(blk.Instructions, 0)
		vm.pushFrame(frame)
		err := vm.run()
		if err != nil {
			vm.popPipeScope()
			vm.popFrame()
			return nil, err
		}
		res := vm.Pop()
		vm.popFrame()
		vm.popPipeScope()
		if b, ok := res.(bool); ok && b {
			return true, nil
		}
	}
	return false, nil
}

func EveryPipeHandler(input any, block any, alias string, vm *VM) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("every pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("every pipe expects a predicate block")
	}
	for i, elem := range arr {
		vm.pushPipeScope()
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", i)
		frame := NewFrame(blk.Instructions, 0)
		vm.pushFrame(frame)
		err := vm.run()
		if err != nil {
			vm.popPipeScope()
			vm.popFrame()
			return nil, err
		}
		res := vm.Pop()
		vm.popFrame()
		vm.popPipeScope()
		if b, ok := res.(bool); !ok || !b {
			return false, nil
		}
	}
	return true, nil
}

func UniquePipeHandler(input any, block any, alias string, vm *VM) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("unique pipe expects array input")
	}
	seen := make(map[string]bool)
	var result []any
	for _, elem := range arr {
		key := fmt.Sprintf("%v", elem)
		if !seen[key] {
			seen[key] = true
			result = append(result, elem)
		}
	}
	return result, nil
}

func SortPipeHandler(input any, block any, alias string, vm *VM) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("sort pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("sort pipe expects a predicate block")
	}
	type sortableElem struct {
		key any
		val any
	}
	sortable := make([]sortableElem, len(arr))
	for i, elem := range arr {
		vm.pushPipeScope()
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", i)
		frame := NewFrame(blk.Instructions, 0)
		vm.pushFrame(frame)
		err := vm.run()
		if err != nil {
			vm.popPipeScope()
			vm.popFrame()
			return nil, err
		}
		key := vm.Pop()
		vm.popFrame()
		vm.popPipeScope()
		sortable[i] = sortableElem{key, elem}
	}
	sort.SliceStable(sortable, func(i, j int) bool {
		ki, kiNum := sortable[i].key.(float64)
		kj, kjNum := sortable[j].key.(float64)
		if kiNum && kjNum {
			return ki < kj
		}
		si, siStr := sortable[i].key.(string)
		sj, sjStr := sortable[j].key.(string)
		if siStr && sjStr {
			return si < sj
		}
		return false
	})
	result := make([]any, len(arr))
	for i, se := range sortable {
		result[i] = se.val
	}
	return result, nil
}

func GroupByPipeHandler(input any, block any, alias string, vm *VM) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("groupBy pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("groupBy pipe expects a predicate block")
	}
	groups := make(map[string][]any)
	for i, elem := range arr {
		vm.pushPipeScope()
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", i)
		frame := NewFrame(blk.Instructions, 0)
		vm.pushFrame(frame)
		err := vm.run()
		if err != nil {
			vm.popPipeScope()
			vm.popFrame()
			return nil, err
		}
		key := vm.Pop()
		vm.popFrame()
		vm.popPipeScope()
		keyStr := fmt.Sprintf("%v", key)
		groups[keyStr] = append(groups[keyStr], elem)
	}
	return groups, nil
}

func WindowPipeHandler(input any, block any, alias string, vm *VM) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("window pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("window pipe expects a predicate block")
	}
	windowSize := 2
	var result []any
	for i := 0; i <= len(arr)-windowSize; i++ {
		window := arr[i : i+windowSize]
		vm.pushPipeScope()
		vm.setPipeVar("$window", window)
		vm.setPipeVar("$index", i)
		frame := NewFrame(blk.Instructions, 0)
		vm.pushFrame(frame)
		err := vm.run()
		if err != nil {
			vm.popPipeScope()
			vm.popFrame()
			return nil, err
		}
		res := vm.Pop()
		vm.popFrame()
		vm.popPipeScope()
		result = append(result, res)
	}
	return result, nil
}

func ChunkPipeHandler(input any, block any, alias string, vm *VM) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("chunk pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("chunk pipe expects a predicate block")
	}
	chunkSize := 2
	var result []any
	for i := 0; i < len(arr); i += chunkSize {
		end := i + chunkSize
		if end > len(arr) {
			end = len(arr)
		}
		chunk := arr[i:end]
		vm.pushPipeScope()
		vm.setPipeVar("$chunk", chunk)
		vm.setPipeVar("$index", i/chunkSize)
		frame := NewFrame(blk.Instructions, 0)
		vm.pushFrame(frame)
		err := vm.run()
		if err != nil {
			vm.popPipeScope()
			vm.popFrame()
			return nil, err
		}
		res := vm.Pop()
		vm.popFrame()
		vm.popPipeScope()
		result = append(result, res)
	}
	return result, nil
}
