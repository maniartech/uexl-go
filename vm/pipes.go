package vm

import (
	"context"
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

// pipeContextImpl is the internal implementation of PipeContext.
// One instance is created per OpPipe dispatch; it holds the predicate block,
// the alias name, and a back-pointer to the executing VM.
//
// Frame reuse: the *Frame is created lazily on the first EvalItem/EvalWith call
// and then reused across all iterations by resetting ip and basePointer only.
// This eliminates per-iteration frame allocations on hot paths.
type pipeContextImpl struct {
	vm    *VM
	block *compiler.InstructionBlock
	alias string
	frame *Frame // lazily created, reused across iterations
}

// EvalItem sets $item, $index (and the alias if declared), then runs the predicate.
// Zero-allocation hot path for map / filter / find / some / every / sort / groupBy / flatMap.
func (p *pipeContextImpl) EvalItem(item any, index int) (any, error) {
	if p.alias != "" {
		p.vm.setPipeVar(p.alias, item)
	}
	p.vm.setPipeVar("$item", item)
	p.vm.setPipeVar("$index", index)
	return p.runFrame()
}

// EvalWith sets arbitrary scope variables, then runs the predicate.
// For reduce/window/chunk: allocate the map once outside the loop and reuse it.
func (p *pipeContextImpl) EvalWith(scopeVars map[string]any) (any, error) {
	for k, v := range scopeVars {
		p.vm.setPipeVar(k, v)
	}
	return p.runFrame()
}

// Alias returns the user-declared alias (e.g. "$x" from "|map as $x:"),
// or empty string when none was declared.
func (p *pipeContextImpl) Alias() string { return p.alias }

// Context returns the evaluation context for cancellation and deadline checks.
func (p *pipeContextImpl) Context() context.Context { return p.vm.ctx }

// runFrame resets and executes the pipe predicate block.
// The frame is allocated once and reused across all iterations.
func (p *pipeContextImpl) runFrame() (any, error) {
	if p.block == nil || p.block.Instructions == nil {
		return nil, fmt.Errorf("pipe predicate block is required")
	}
	if p.frame == nil {
		p.frame = NewFrame(p.block.Instructions, 0)
	}
	p.frame.ip = 0
	p.frame.basePointer = p.vm.sp
	p.vm.pushFrame(p.frame)
	err := p.vm.run()
	if err != nil {
		p.vm.popFrame()
		return nil, err
	}
	res := p.vm.Pop()
	p.vm.popFrame()
	return res, nil
}

func DefaultPipeHandler(ctx PipeContext, input any) (any, error) {
	pctx := ctx.(*pipeContextImpl)
	if pctx.block == nil || pctx.block.Instructions == nil {
		return input, nil
	}
	return ctx.EvalWith(map[string]any{"$last": input})
}

func MapPipeHandler(ctx PipeContext, input any) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("map pipe expects array input")
	}
	result := make([]any, len(arr))
	for i, elem := range arr {
		val, err := ctx.EvalItem(elem, i)
		if err != nil {
			return nil, err
		}
		result[i] = val
	}
	return result, nil
}

func FilterPipeHandler(ctx PipeContext, input any) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("filter pipe expects array input")
	}
	var result []any
	for i, elem := range arr {
		keep, err := ctx.EvalItem(elem, i)
		if err != nil {
			return nil, err
		}
		if b, ok := keep.(bool); ok && b {
			result = append(result, elem)
		}
	}
	return result, nil
}

func ReducePipeHandler(ctx PipeContext, input any) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("reduce pipe expects array input")
	}
	if len(arr) == 0 {
		return nil, fmt.Errorf("reduce pipe cannot operate on empty array")
	}
	var acc any
	// Allocate scope map once and reuse across iterations — avoids per-iteration allocation.
	scope := make(map[string]any, 3)
	for i, elem := range arr {
		scope["$acc"] = acc
		scope["$item"] = elem
		scope["$index"] = i
		var err error
		acc, err = ctx.EvalWith(scope)
		if err != nil {
			return nil, err
		}
	}
	return acc, nil
}

func FindPipeHandler(ctx PipeContext, input any) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("find pipe expects array input")
	}
	for i, elem := range arr {
		matched, err := ctx.EvalItem(elem, i)
		if err != nil {
			return nil, err
		}
		if b, ok := matched.(bool); ok && b {
			return elem, nil
		}
	}
	return nil, nil
}

func SomePipeHandler(ctx PipeContext, input any) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("some pipe expects array input")
	}
	for i, elem := range arr {
		matched, err := ctx.EvalItem(elem, i)
		if err != nil {
			return nil, err
		}
		if b, ok := matched.(bool); ok && b {
			return true, nil
		}
	}
	return false, nil
}

func EveryPipeHandler(ctx PipeContext, input any) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("every pipe expects array input")
	}
	for i, elem := range arr {
		matched, err := ctx.EvalItem(elem, i)
		if err != nil {
			return nil, err
		}
		if b, ok := matched.(bool); !ok || !b {
			return false, nil
		}
	}
	return true, nil
}

func UniquePipeHandler(ctx PipeContext, input any) (any, error) {
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

func SortPipeHandler(ctx PipeContext, input any) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("sort pipe expects array input")
	}
	type sortableElem struct {
		key any
		val any
	}
	sortable := make([]sortableElem, len(arr))
	for i, elem := range arr {
		key, err := ctx.EvalItem(elem, i)
		if err != nil {
			return nil, err
		}
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

func GroupByPipeHandler(ctx PipeContext, input any) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("groupBy pipe expects array input")
	}
	groups := make(map[string][]any)
	for i, elem := range arr {
		key, err := ctx.EvalItem(elem, i)
		if err != nil {
			return nil, err
		}
		keyStr := fmt.Sprintf("%v", key)
		groups[keyStr] = append(groups[keyStr], elem)
	}
	return groups, nil
}

func WindowPipeHandler(ctx PipeContext, input any) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("window pipe expects array input")
	}
	windowSize := 2
	var result []any
	scope := make(map[string]any, 2) // allocated once, reused across iterations
	for i := 0; i <= len(arr)-windowSize; i++ {
		scope["$window"] = arr[i : i+windowSize]
		scope["$index"] = i
		res, err := ctx.EvalWith(scope)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	return result, nil
}

func ChunkPipeHandler(ctx PipeContext, input any) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("chunk pipe expects array input")
	}
	chunkSize := 2
	var result []any
	scope := make(map[string]any, 2) // allocated once, reused across iterations
	for i := 0; i < len(arr); i += chunkSize {
		end := i + chunkSize
		if end > len(arr) {
			end = len(arr)
		}
		scope["$chunk"] = arr[i:end]
		scope["$index"] = i / chunkSize
		res, err := ctx.EvalWith(scope)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	return result, nil
}

// FlatMapPipeHandler maps each element and flattens results in one operation
func FlatMapPipeHandler(ctx PipeContext, input any) (any, error) {
	arr, ok := input.([]any)
	if !ok {
		return nil, fmt.Errorf("flatMap pipe expects array input")
	}
	var result []any
	for i, elem := range arr {
		res, err := ctx.EvalItem(elem, i)
		if err != nil {
			return nil, err
		}
		if resArr, ok := res.([]any); ok {
			result = append(result, resArr...)
		} else {
			result = append(result, res)
		}
	}
	return result, nil
}
