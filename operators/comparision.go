package operators

import (
	"fmt"

	"github.com/maniartech/uexl_go/core"
	"github.com/maniartech/uexl_go/types"
)

func comparer(op string, a, b core.Evaluator, ctx types.Context) (types.Value, error) {
	aval, err := a.Eval(ctx)
	if err != nil {
		return nil, err
	}

	if aval, ok := aval.(types.Comparer); ok {
		bval, err := b.Eval(ctx)
		if err != nil {
			return nil, err
		}

		result, err := aval.Compare(bval)
		if err != nil {
			return nil, err
		}

		switch op {
		case "<":
			return types.Boolean(result < 0), nil
		case "<=":
			return types.Boolean(result <= 0), nil
		case ">":
			return types.Boolean(result > 0), nil
		case ">=":
			return types.Boolean(result >= 0), nil
		}
	}

	// The value does not support comparison
	return nil, fmt.Errorf("invalid argument type for comparer operator: %T, %T", a, b)
}

func equaler(op string, a, b core.Evaluator, ctx types.Context) (types.Value, error) {
	aval, err := a.Eval(ctx)
	if err != nil {
		return nil, err
	}

	bval, err := b.Eval(ctx)
	if err != nil {
		return nil, err
	}

	eq := aval.Equals(bval)
	if op == "==" {
		return types.Boolean(eq), nil
	}

	return types.Boolean(!eq), nil
}

func init() {
	BinaryOpRegistry.Register("==", equaler)
	BinaryOpRegistry.Register("!=", equaler)

	BinaryOpRegistry.Register("<", comparer)
	BinaryOpRegistry.Register("<=", comparer)
	BinaryOpRegistry.Register(">", comparer)
	BinaryOpRegistry.Register(">=", comparer)
}
