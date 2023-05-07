package operators

import (
	"github.com/maniartech/uexl_go/core"
	"github.com/maniartech/uexl_go/types"
)

func Eval(op string, a, b core.Evaluator, ctx types.Context) (types.Value, error) {
	return BinaryOpRegistry.Get(op)(op, a, b, ctx)
}
