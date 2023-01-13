package core

import "github.com/maniartech/uexl_go/types"

type Evaluator interface {

	// Eval evaluates
	Eval(context types.Map) (interface{}, error)
}
