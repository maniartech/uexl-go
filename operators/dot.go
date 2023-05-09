package operators

import (
	"fmt"

	"github.com/maniartech/uexl_go/core"
	"github.com/maniartech/uexl_go/types"
)

func dot(op string, a, b core.Evaluator, ctx types.Context) (types.Value, error) {
	aval, err := a.Eval(ctx)
	if err != nil {
		return nil, err
	}

	if _, ok := aval.(types.Dot); ok {
		actx := types.Context(aval.(types.Object))
		res, err := b.Eval(actx)
		if err != nil {
			return nil, err
		}

		return res, nil
		// return aval.Dot(bval)
	}

	// The value does not support addition
	return nil, fmt.Errorf("the dot operator is not supported for %s", aval.Type())
}

func init() {
	BinaryOpRegistry.Register(".", dot)
}
