package vm

import (
	"fmt"
	"sort"

	"github.com/maniartech/uexl_go/compiler"
	"github.com/maniartech/uexl_go/parser"
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

// SORT: sorts array by predicate result (e.g. $item.property)
func SortPipeHandler(input parser.Node, block any, alias string, vm *VM) (parser.Node, error) {
	arr, ok := input.(*parser.ArrayLiteral)
	if !ok {
		return nil, fmt.Errorf("sort pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("sort pipe expects a predicate block")
	}

	type sortableElem struct {
		key parser.Node
		val parser.Node
	}
	sortable := make([]sortableElem, len(arr.Elements))
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
		key := vm.Pop()
		vm.popFrame()
		vm.popPipeScope()
		sortable[i] = sortableElem{key, elem}
	}
	sort.SliceStable(sortable, func(i, j int) bool {
		ki, ok1 := sortable[i].key.(*parser.NumberLiteral)
		kj, ok2 := sortable[j].key.(*parser.NumberLiteral)
		if ok1 && ok2 {
			return ki.Value < kj.Value
		}
		si, ok1 := sortable[i].key.(*parser.StringLiteral)
		sj, ok2 := sortable[j].key.(*parser.StringLiteral)
		if ok1 && ok2 {
			return si.Value < sj.Value
		}
		return false // fallback: keep original order
	})
	result := make([]parser.Expression, len(arr.Elements))
	for i, se := range sortable {
		result[i] = se.val.(parser.Expression)
	}
	return &parser.ArrayLiteral{Elements: result}, nil
}

// GROUPBY: groups array by predicate result (e.g. $item.property)
func GroupByPipeHandler(input parser.Node, block any, alias string, vm *VM) (parser.Node, error) {
	arr, ok := input.(*parser.ArrayLiteral)
	if !ok {
		return nil, fmt.Errorf("groupBy pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("groupBy pipe expects a predicate block")
	}
	groups := make(map[string][]parser.Expression)
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
		key := vm.Pop()
		vm.popFrame()
		vm.popPipeScope()
		keyStr := fmt.Sprintf("%v", key)
		groups[keyStr] = append(groups[keyStr], elem)
	}
	// Convert to object literal: key -> array
	obj := make(map[string]parser.Expression)
	for k, v := range groups {
		obj[k] = &parser.ArrayLiteral{Elements: v}
	}
	return &parser.ObjectLiteral{Properties: obj}, nil
}

// WINDOW: splits array into windows of size N, runs predicate on each window
func WindowPipeHandler(input parser.Node, block any, alias string, vm *VM) (parser.Node, error) {
	arr, ok := input.(*parser.ArrayLiteral)
	if !ok {
		return nil, fmt.Errorf("window pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("window pipe expects a predicate block")
	}
	// For demo, use window size 2 (could be parameterized)
	windowSize := 2
	var result []parser.Node
	for i := 0; i <= len(arr.Elements)-windowSize; i++ {
		window := arr.Elements[i : i+windowSize]
		vm.pushPipeScope()
		vm.setPipeVar("$window", &parser.ArrayLiteral{Elements: window})
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
		result = append(result, res)
	}
	exprs := make([]parser.Expression, len(result))
	for i, n := range result {
		exprs[i] = n.(parser.Expression)
	}
	return &parser.ArrayLiteral{Elements: exprs}, nil
}

// CHUNK: splits array into chunks of size N, runs predicate on each chunk
func ChunkPipeHandler(input parser.Node, block any, alias string, vm *VM) (parser.Node, error) {
	arr, ok := input.(*parser.ArrayLiteral)
	if !ok {
		return nil, fmt.Errorf("chunk pipe expects array input")
	}
	blk, ok := block.(*compiler.InstructionBlock)
	if !ok || blk == nil || blk.Instructions == nil {
		return nil, fmt.Errorf("chunk pipe expects a predicate block")
	}
	// For demo, use chunk size 2 (could be parameterized)
	chunkSize := 2
	var result []parser.Node
	for i := 0; i < len(arr.Elements); i += chunkSize {
		end := i + chunkSize
		if end > len(arr.Elements) {
			end = len(arr.Elements)
		}
		chunk := arr.Elements[i:end]
		vm.pushPipeScope()
		vm.setPipeVar("$chunk", &parser.ArrayLiteral{Elements: chunk})
		vm.setPipeVar("$index", &parser.NumberLiteral{Value: float64(i / chunkSize)})
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
		result = append(result, res)
	}
	exprs := make([]parser.Expression, len(result))
	for i, n := range result {
		exprs[i] = n.(parser.Expression)
	}
	return &parser.ArrayLiteral{Elements: exprs}, nil
}
