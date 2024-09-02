package core

import "github.com/maniartech/uexl_go/types"

type Evaluator interface {
	// Eval evaluates
	Eval(context types.Context) (types.Value, error)
}

type Truethy interface {

	// IsTrue returns true if the value is truethy
	IsTrue() bool
}

func IsTruthy(value interface{}) bool {
	if value == nil {
		return false
	}

	if value, ok := value.(Truethy); ok {
		return value.IsTrue()
	}

	return false
}
