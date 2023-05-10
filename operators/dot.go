package operators

import (
	"fmt"

	"github.com/maniartech/uexl_go/core"
	"github.com/maniartech/uexl_go/types"
)

func DotEval(expr core.Evaluator, key string, ctx types.Context) (types.Value, error) {
	val, err := expr.Eval(ctx)
	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, fmt.Errorf("the dot operator is not supported for nil")
	}

	if vald, ok := val.(types.Dot); ok {
		return vald.Dot(key)
	}

	// The value does not support addition
	return nil, fmt.Errorf("the dot operator is not supported for %s", val.Type())
}
