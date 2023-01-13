package operators

import (
	"fmt"

	"github.com/maniartech/uexl_go/types"
)

func LessThan(a, b any) (any, error) {
	numA, numOkA := a.(types.Number)
	numB, numOkB := b.(types.Number)

	if numOkA && numOkB {
		return numA < numB, nil
	}

	return nil, fmt.Errorf("invalid argument type for plus operator: %T, %T", a, b)
}

func LessThanEqual(a, b any) (any, error) {
	numA, numOkA := a.(types.Number)
	numB, numOkB := b.(types.Number)

	if numOkA && numOkB {
		return numA <= numB, nil
	}

	return nil, fmt.Errorf("invalid argument type for plus operator: %T, %T", a, b)
}

func GreaterThan(a, b any) (any, error) {
	numA, numOkA := a.(types.Number)
	numB, numOkB := b.(types.Number)

	if numOkA && numOkB {
		return numA > numB, nil
	}

	return nil, fmt.Errorf("invalid argument type for plus operator: %T, %T", a, b)
}

func GreaterThanEqual(a, b any) (any, error) {
	numA, numOkA := a.(types.Number)
	numB, numOkB := b.(types.Number)

	if numOkA && numOkB {
		return numA >= numB, nil
	}

	return nil, fmt.Errorf("invalid argument type for plus operator: %T, %T", a, b)
}

func init() {
	Registry.Register("<", LessThan)
	Registry.Register("<=", LessThanEqual)
	Registry.Register(">", GreaterThan)
	Registry.Register(">=", GreaterThanEqual)
}
