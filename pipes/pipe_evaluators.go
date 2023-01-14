package pipes

import (
	"fmt"

	"github.com/maniartech/uexl_go/core"
	"github.com/maniartech/uexl_go/types"
)

// firstEvaluator evalues the first node in the pipe.
func firstEvaluator(evaluator core.Evaluator, context types.Map, prevResult any) (interface{}, error) {
	if context == nil {
		context = make(types.Map)
	}

	return evaluator.Eval(context)
}

// passEvaluator evalues the passes the result of the previous node to the current node.
func passEvaluator(evaluator core.Evaluator, context types.Map, prevResult any) (interface{}, error) {
	defer delete(context, "$1")

	context["$1"] = prevResult
	return evaluator.Eval(context)
}

// mapEvaluator evalues the map node in the pipe.
func mapEvaluator(evaluator core.Evaluator, context types.Map, prevResult any) (interface{}, error) {
	defer delete(context, "$1")

	array, ok := prevResult.([]interface{})
	if !ok {
		return nil, fmt.Errorf("filter expects an array")
	}

	newArray := make([]interface{}, 0, len(array))
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
func filterEvaluator(evaluator core.Evaluator, context types.Map, prevResult any) (interface{}, error) {
	defer delete(context, "$1")

	array, ok := prevResult.([]interface{})
	if !ok {
		return nil, fmt.Errorf("filter expects an array")
	}

	newArray := make([]interface{}, 0, len(array))
	for i := 0; i < len(array); i++ {
		context["$1"] = array[i]
		result, err := evaluator.Eval(context)
		if err != nil {
			return nil, err
		}
		if result.(bool) {
			newArray = append(newArray, array[i])
		}
	}

	return newArray, nil
}

// findEvaluator evalues the find node in the pipe.
func findEvaluator(evaluator core.Evaluator, context types.Map, prevResult any) (interface{}, error) {
	defer delete(context, "$1")

	array, ok := prevResult.([]interface{})
	if !ok {
		return nil, fmt.Errorf("filter expects an array")
	}

	for i := 0; i < len(array); i++ {
		context["$1"] = array[i]
		result, err := evaluator.Eval(context)
		if err != nil {
			return nil, err
		}
		if result.(bool) {
			return array[i], nil
		}
	}

	return nil, nil
}

func init() {
	_handlers.Register("first", firstEvaluator)
	_handlers.Register("pass", passEvaluator)

	_handlers.Register("map", mapEvaluator)
	_handlers.Register("filter", filterEvaluator)
	_handlers.Register("find", findEvaluator)

}
