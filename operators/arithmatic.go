package operators

import (
	"fmt"

	"github.com/maniartech/uexl_go/core"
	"github.com/maniartech/uexl_go/types"
)

// Plus performs addition of two values
func Plus(op string, a, b core.Evaluator, ctx types.Context) (types.Value, error) {
	aval, err := a.Eval(ctx)
	if err != nil {
		return nil, err
	}

	if aval, ok := aval.(types.Adder); ok {
		bval, err := b.Eval(ctx)
		if err != nil {
			return nil, err
		}

		return aval.Add(bval)
	}

	// The value does not support addition
	return nil, fmt.Errorf("invalid argument type for plus operator: %T, %T", a, b)
}

// Minus performs subtraction of two values
func Minus(op string, a, b core.Evaluator, ctx types.Context) (types.Value, error) {
	aval, err := a.Eval(ctx)
	if err != nil {
		return nil, err
	}

	if aval, ok := aval.(types.Subtractor); ok {
		bval, err := b.Eval(ctx)
		if err != nil {
			return nil, err
		}

		return aval.Subtract(bval)
	}

	// The value does not support subtraction
	return nil, fmt.Errorf("invalid argument type for minus operator: %T, %T", a, b)
}

// Times performs multiplication of two values
func Times(op string, a, b core.Evaluator, ctx types.Context) (types.Value, error) {
	aval, err := a.Eval(ctx)
	if err != nil {
		return nil, err
	}

	if aval, ok := aval.(types.Multiplier); ok {
		bval, err := b.Eval(ctx)
		if err != nil {
			return nil, err
		}

		return aval.Multiply(bval)
	}

	// The value does not support multiplication
	return nil, fmt.Errorf("invalid argument type for Times operator: %T, %T", a, b)
}

func Divide(op string, a, b core.Evaluator, ctx types.Context) (types.Value, error) {
	aval, err := a.Eval(ctx)
	if err != nil {
		return nil, err
	}

	if aval, ok := aval.(types.Divider); ok {
		bval, err := b.Eval(ctx)
		if err != nil {
			return nil, err
		}

		return aval.Divide(bval)
	}

	// The value does not support division
	return nil, fmt.Errorf("invalid argument type for divide operator: %T, %T", a, b)
}

func Modulo(op string, a, b core.Evaluator, ctx types.Context) (types.Value, error) {
	aval, err := a.Eval(ctx)
	if err != nil {
		return nil, err
	}

	if aval, ok := aval.(types.Modulus); ok {
		bval, err := b.Eval(ctx)
		if err != nil {
			return nil, err
		}

		return aval.Mod(bval)
	}

	// The value does not support modulo
	return nil, fmt.Errorf("invalid argument type for modulo operator: %T, %T", a, b)
}

func init() {
	BinaryOpRegistry.Register("+", Plus)
	BinaryOpRegistry.Register("-", Minus)
	BinaryOpRegistry.Register("*", Times)
	BinaryOpRegistry.Register("/", Divide)
	BinaryOpRegistry.Register("//", Modulo)
}
