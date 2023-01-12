package functions

import (
	"fmt"

	"github.com/maniartech/uexl_go/types"
)

func sum(args []interface{}) (interface{}, error) {
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

func average(args []interface{}) (interface{}, error) {
	var sum types.Number
	for _, arg := range args {
		switch v := arg.(type) {
		case types.Number:
			sum += types.Number(v)
		default:
			return nil, fmt.Errorf("invalid argument type for average function: %T", v)
		}
	}

	return sum / types.Number(len(args)), nil
}
