package functions

import (
	"fmt"

	"github.com/maniartech/uexl_go/types"
)

func sum(args []any) (any, error) {
	var sum types.Number
	for _, arg := range args {
		switch v := arg.(type) {
		case types.Number:
			sum += types.Number(v)
		default:
			return nil, fmt.Errorf("invalid argument type for add function: %T", v)
		}
	}

	return sum, nil
}

func subtract(args []any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments for subtract function: %d", len(args))
	}

	v1 := args[0]
	v2 := args[1]

	switch v1.(type) {
	case types.Number:
		switch v2.(type) {
		case types.Number:
			return v1.(types.Number) - v2.(types.Number), nil
		}
	}

	return nil, fmt.Errorf("invalid argument type for subtract function: %T, %T", v1, v2)
}
