package vm

import (
	"fmt"
	"sort"

	"github.com/maniartech/uexl/compiler"
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
	"flatMap": FlatMapPipeHandler,
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

	// Fast path optimization removed: tryFastMapArithmetic was a benchmark cheat

	result := make([]any, len(arr))

	// Optimization: Reuse pipe scope and frame for all iterations
	vm.pushPipeScope()
	frame := NewFrame(blk.Instructions, 0)

	// Pre-set alias and $index keys once (if needed)
	if alias != "" {
		vm.setPipeVar(alias, nil) // Initialize key
	}
	vm.setPipeVar("$item", nil) // Initialize key
	vm.setPipeVar("$index", 0)  // Initialize key

	for i, elem := range arr {
		// Update scope variables for this iteration using fast-path
		if alias != "" {
			vm.setPipeVar(alias, elem)
		}
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", i)

		// Reset frame for reuse
		frame.ip = 0
		frame.basePointer = vm.sp

		vm.pushFrame(frame)
		err := vm.run()
		if err != nil {
			vm.popFrame()
			vm.popPipeScope()
			return nil, err
		}
		res := vm.Pop()
		vm.popFrame()
		result[i] = res
	}

	vm.popPipeScope()
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

	// Optimization: Reuse pipe scope and frame for all iterations
	vm.pushPipeScope()
	frame := NewFrame(blk.Instructions, 0)

	// Pre-initialize scope keys once
	if alias != "" {
		vm.setPipeVar(alias, nil) // Initialize key
	}
	vm.setPipeVar("$item", nil) // Initialize key
	vm.setPipeVar("$index", 0)  // Initialize key

	for i, elem := range arr {
		// Update scope variables for this iteration
		if alias != "" {
			vm.setPipeVar(alias, elem)
		}
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", i)

		// Reset frame for reuse
		frame.ip = 0
		frame.basePointer = vm.sp
		vm.pushFrame(frame)

		err := vm.run()
		if err != nil {
			vm.popFrame()
			vm.popPipeScope()
			return nil, err
		}

		res := vm.Pop()
		vm.popFrame()

		if b, ok := res.(bool); ok && b {
			result = append(result, elem)
		}
	}

	vm.popPipeScope()
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

	var acc any

	// Optimization: Reuse pipe scope and frame for all iterations
	vm.pushPipeScope()
	frame := NewFrame(blk.Instructions, 0)

	// Pre-initialize scope keys once
	vm.setPipeVar("$acc", nil)  // Initialize key
	vm.setPipeVar("$item", nil) // Initialize key
	vm.setPipeVar("$index", 0)  // Initialize key

	for i := 0; i < len(arr); i++ {
		// Update scope variables for this iteration
		vm.setPipeVar("$acc", acc)
		vm.setPipeVar("$item", arr[i])
		vm.setPipeVar("$index", i)

		// Reset frame for reuse
		frame.ip = 0
		frame.basePointer = vm.sp
		vm.pushFrame(frame)

		err := vm.run()
		if err != nil {
			vm.popFrame()
			vm.popPipeScope()
			return nil, err
		}

		acc = vm.Pop()
		vm.popFrame()
	}

	vm.popPipeScope()
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

	// Optimization: Reuse pipe scope and frame for all iterations
	vm.pushPipeScope()
	frame := NewFrame(blk.Instructions, 0)

	// Pre-initialize scope keys once
	vm.setPipeVar("$item", nil) // Initialize key
	vm.setPipeVar("$index", 0)  // Initialize key

	for i, elem := range arr {
		// Update scope variables for this iteration
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", i)

		// Reset frame for reuse
		frame.ip = 0
		frame.basePointer = vm.sp
		vm.pushFrame(frame)

		err := vm.run()
		if err != nil {
			vm.popFrame()
			vm.popPipeScope()
			return nil, err
		}

		res := vm.Pop()
		vm.popFrame()

		if b, ok := res.(bool); ok && b {
			vm.popPipeScope()
			return elem, nil
		}
	}

	vm.popPipeScope()
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

	// Optimization: Reuse pipe scope and frame for all iterations
	vm.pushPipeScope()
	frame := NewFrame(blk.Instructions, 0)

	// Pre-initialize scope keys once
	vm.setPipeVar("$item", nil) // Initialize key
	vm.setPipeVar("$index", 0)  // Initialize key

	for i, elem := range arr {
		// Update scope variables for this iteration
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", i)

		// Reset frame for reuse
		frame.ip = 0
		frame.basePointer = vm.sp
		vm.pushFrame(frame)

		err := vm.run()
		if err != nil {
			vm.popFrame()
			vm.popPipeScope()
			return nil, err
		}

		res := vm.Pop()
		vm.popFrame()

		if b, ok := res.(bool); ok && b {
			vm.popPipeScope()
			return true, nil
		}
	}

	vm.popPipeScope()
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

	// Optimization: Reuse pipe scope and frame for all iterations
	vm.pushPipeScope()
	frame := NewFrame(blk.Instructions, 0)

	// Pre-initialize scope keys once
	vm.setPipeVar("$item", nil) // Initialize key
	vm.setPipeVar("$index", 0)  // Initialize key

	for i, elem := range arr {
		// Update scope variables for this iteration
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", i)

		// Reset frame for reuse
		frame.ip = 0
		frame.basePointer = vm.sp
		vm.pushFrame(frame)

		err := vm.run()
		if err != nil {
			vm.popFrame()
			vm.popPipeScope()
			return nil, err
		}

		res := vm.Pop()
		vm.popFrame()

		if b, ok := res.(bool); !ok || !b {
			vm.popPipeScope()
			return false, nil
		}
	}

	vm.popPipeScope()
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

	// Optimization: Reuse pipe scope and frame for all iterations
	vm.pushPipeScope()
	frame := NewFrame(blk.Instructions, 0)

	// Pre-initialize scope keys once
	vm.setPipeVar("$item", nil) // Initialize key
	vm.setPipeVar("$index", 0)  // Initialize key

	for i, elem := range arr {
		// Update scope variables for this iteration
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", i)

		// Reset frame for reuse
		frame.ip = 0
		frame.basePointer = vm.sp
		vm.pushFrame(frame)

		err := vm.run()
		if err != nil {
			vm.popFrame()
			vm.popPipeScope()
			return nil, err
		}

		key := vm.Pop()
		vm.popFrame()
		sortable[i] = sortableElem{key, elem}
	}

	vm.popPipeScope()

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

	// Optimization: Reuse pipe scope and frame for all iterations
	vm.pushPipeScope()
	frame := NewFrame(blk.Instructions, 0)

	// Pre-initialize scope keys once
	vm.setPipeVar("$item", nil) // Initialize key
	vm.setPipeVar("$index", 0)  // Initialize key

	for i, elem := range arr {
		// Update scope variables for this iteration
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", i)

		// Reset frame for reuse
		frame.ip = 0
		frame.basePointer = vm.sp
		vm.pushFrame(frame)

		err := vm.run()
		if err != nil {
			vm.popFrame()
			vm.popPipeScope()
			return nil, err
		}

		key := vm.Pop()
		vm.popFrame()
		keyStr := fmt.Sprintf("%v", key)
		groups[keyStr] = append(groups[keyStr], elem)
	}

	vm.popPipeScope()
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

	// Optimization: Reuse pipe scope and frame for all iterations
	vm.pushPipeScope()
	frame := NewFrame(blk.Instructions, 0)

	// Pre-initialize scope keys once
	vm.setPipeVar("$window", nil) // Initialize key
	vm.setPipeVar("$index", 0)    // Initialize key

	for i := 0; i <= len(arr)-windowSize; i++ {
		window := arr[i : i+windowSize]

		// Update scope variables for this iteration
		vm.setPipeVar("$window", window)
		vm.setPipeVar("$index", i)

		// Reset frame for reuse
		frame.ip = 0
		frame.basePointer = vm.sp
		vm.pushFrame(frame)

		err := vm.run()
		if err != nil {
			vm.popFrame()
			vm.popPipeScope()
			return nil, err
		}

		res := vm.Pop()
		vm.popFrame()
		result = append(result, res)
	}

	vm.popPipeScope()
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

	// Optimization: Reuse pipe scope and frame for all iterations
	vm.pushPipeScope()
	frame := NewFrame(blk.Instructions, 0)

	// Pre-initialize scope keys once
	vm.setPipeVar("$chunk", nil) // Initialize key
	vm.setPipeVar("$index", 0)   // Initialize key

	for i := 0; i < len(arr); i += chunkSize {
		end := i + chunkSize
		if end > len(arr) {
			end = len(arr)
		}
		chunk := arr[i:end]

		// Update scope variables for this iteration
		vm.setPipeVar("$chunk", chunk)
		vm.setPipeVar("$index", i/chunkSize)

		// Reset frame for reuse
		frame.ip = 0
		frame.basePointer = vm.sp
		vm.pushFrame(frame)

		err := vm.run()
		if err != nil {
			vm.popFrame()
			vm.popPipeScope()
			return nil, err
		}

		res := vm.Pop()
		vm.popFrame()
		result = append(result, res)
	}

	vm.popPipeScope()
	return result, nil
}

// FlatMapPipeHandler maps each element and flattens results in one operation
func FlatMapPipeHandler(input any, block any, alias string, vm *VM) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("flatMap pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("flatMap pipe expects a predicate block")
	}

	var result []any

	// Optimization: Reuse pipe scope and frame for all iterations
	vm.pushPipeScope()
	frame := NewFrame(blk.Instructions, 0)

	// Pre-initialize scope keys once
	if alias != "" {
		vm.setPipeVar(alias, nil) // Initialize key
	}
	vm.setPipeVar("$item", nil) // Initialize key
	vm.setPipeVar("$index", 0)  // Initialize key

	for i, elem := range arr {
		// Update scope variables for this iteration
		if alias != "" {
			vm.setPipeVar(alias, elem)
		}
		vm.setPipeVar("$item", elem)
		vm.setPipeVar("$index", i)

		// Reset frame for reuse
		frame.ip = 0
		frame.basePointer = vm.sp
		vm.pushFrame(frame)

		err := vm.run()
		if err != nil {
			vm.popFrame()
			vm.popPipeScope()
			return nil, err
		}

		res := vm.Pop()
		vm.popFrame()

		// Flatten: if result is an array, append its elements
		if resArr, ok := res.([]any); ok {
			result = append(result, resArr...)
		} else {
			// If not an array, append the value itself
			result = append(result, res)
		}
	}

	vm.popPipeScope()
	return result, nil
}
