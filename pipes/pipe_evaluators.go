package pipes

import (
	"github.com/maniartech/uexl_go/core"
	"github.com/maniartech/uexl_go/types"
)

// firstEvaluator evalues the first node in the pipe.
func firstEvaluator(evaluator core.Evaluator, context types.Context, prevResult types.Value) (types.Value, error) {
	if context == nil {
		context = make(types.Context)
	}

	return evaluator.Eval(context)
}

// pipeEvaluator evalues the passes the result of the previous node to the current node.
func pipeEvaluator(evaluator core.Evaluator, context types.Context, prevResult types.Value) (types.Value, error) {
	defer delete(context, "$1")

	context["$1"] = prevResult
	result, err := evaluator.Eval(context)
	return result, err
}

// mapEvaluator evalues the map node in the pipe.
func mapEvaluator(evaluator core.Evaluator, context types.Context, prevResult types.Value) (types.Value, error) {
	defer delete(context, "$1")

	array, ok := prevResult.(types.Array)
	// If the previous result is not an array, then make it an array.
	if !ok {
		array = types.Array{
			prevResult.(types.Value),
		}
	}

	newArray := make(types.Array, 0, len(array))
	for i := 0; i < len(array); i++ {
		context["$1"] = array[i]
		result, err := evaluator.Eval(context)
		if err != nil {
			return nil, err
		}
		newArray = append(newArray, result)
	}

	return newArray, nil
}

// filterEvaluator evalues the filter node in the pipe.
func filterEvaluator(evaluator core.Evaluator, context types.Context, prevResult types.Value) (types.Value, error) {
	defer delete(context, "$1")

	array, ok := prevResult.(types.Array)

	// If the previous result is not an array, then make it an array.
	if !ok {
		array = types.Array{
			prevResult.(types.Value),
		}
	}

	newArray := make(types.Array, 0, len(array))
	for i := 0; i < len(array); i++ {
		context["$1"] = array[i]
		result, err := evaluator.Eval(context)
		if err != nil {
			return nil, err
		}
		if result.IsTruthy() {
			newArray = append(newArray, array[i])
		}
	}

	return newArray, nil
}

// findEvaluator evalues the find node in the pipe.
func findEvaluator(evaluator core.Evaluator, context types.Context, prevResult types.Value) (types.Value, error) {
	defer delete(context, "$1")

	array, ok := prevResult.(types.Array)
	// If the previous result is not an array, then make it an array.
	if !ok {
		array = types.Array{
			prevResult.(types.Value),
		}
	}

	for i := 0; i < len(array); i++ {
		context["$1"] = array[i]
		result, err := evaluator.Eval(context)
		if err != nil {
			return nil, err
		}
		if result.IsTruthy() {
			return array[i], nil
		}
	}

	return nil, nil
}

func init() {
	_handlers.Register("first", firstEvaluator)
	_handlers.Register("pipe", pipeEvaluator)
	_handlers.Register("map", mapEvaluator)
	_handlers.Register("filter", filterEvaluator)
	_handlers.Register("find", findEvaluator)
}
