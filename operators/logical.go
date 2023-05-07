package operators

import (
	"github.com/maniartech/uexl_go/core"
	"github.com/maniartech/uexl_go/types"
)

func LogicalAnd(op string, a, b core.Evaluator, ctx types.Context) (types.Value, error) {
	aval, err := a.Eval(ctx)
	if err != nil {
		return nil, err
	}

	if !aval.IsTruthy() {
		return types.Boolean(false), nil
	}

	bval, err := b.Eval(ctx)
	if err != nil {
		return nil, err
	}

	return types.Boolean(bval.IsTruthy()), nil
}

func LogicalOr(op string, a, b core.Evaluator, ctx types.Context) (types.Value, error) {
	aval, err := a.Eval(ctx)
	if err != nil {
		return nil, err
	}

	if aval.IsTruthy() {
		return aval, nil
	}

	bval, err := b.Eval(ctx)
	if err != nil {
		return nil, err
	}

	if core.IsTruthy(bval) {
		return bval, nil
	}

	return types.Boolean(false), nil
}

func init() {
	BinaryOpRegistry.Register("&&", LogicalAnd)
	BinaryOpRegistry.Register("||", LogicalOr)
}
