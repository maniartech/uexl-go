package operators

import (
	"fmt"
	"strings"

	"github.com/maniartech/uexl_go/types"
)

func Plus(a, b any) (any, error) {
	numA, numOkA := a.(types.Number)
	numB, numOkB := b.(types.Number)

	if numOkA && numOkB {
		return numA + numB, nil
	}

	return nil, fmt.Errorf("invalid argument type for plus operator: %T, %T", a, b)
}

func Minus(a, b any) (any, error) {
	numA, numOkA := a.(types.Number)
	numB, numOkB := b.(types.Number)

	if numOkA && numOkB {
		return numA - numB, nil
	}

	return nil, fmt.Errorf("invalid argument type for minus operator: %T, %T", a, b)
}

func Multiply(a, b any) (any, error) {
	numA, numOkA := a.(types.Number)
	numB, numOkB := b.(types.Number)

	if numOkA && numOkB {
		return numA * numB, nil
	}

	// If one of the arguments is a string, and the other is a number, then
	// repeat the string the number of times specified by the number.
	strA, strOkA := a.(types.String)
	strB, strOkB := b.(types.String)
	if strOkA && numOkB {
		return types.String(strings.Repeat(string(strA), int(numB))), nil
	} else if numOkA && strOkB {
		return types.String(strings.Repeat(string(strB), int(numA))), nil
	}

	return nil, fmt.Errorf("invalid argument type for multiply operator: %T, %T", a, b)
}

func Divide(a, b any) (any, error) {
	numA, numOkA := a.(types.Number)
	numB, numOkB := b.(types.Number)

	if numOkA && numOkB {
		return numA / numB, nil
	}

	return nil, fmt.Errorf("invalid argument type for divide operator: %T, %T", a, b)
}

func Modulo(a, b any) (any, error) {
	numA, numOkA := a.(types.Number)
	numB, numOkB := b.(types.Number)

	if numOkA && numOkB {
		return types.Number(int(numA) % int(numB)), nil
	}

	return nil, fmt.Errorf("invalid argument type for modulo operator: %T, %T", a, b)
}

func init() {
	Registry.Register("+", Plus)
	Registry.Register("-", Minus)
	Registry.Register("*", Multiply)
	Registry.Register("/", Divide)
	Registry.Register("%", Modulo)
}
