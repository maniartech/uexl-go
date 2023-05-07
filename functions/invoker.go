package functions

import (
	"fmt"

	"github.com/maniartech/uexl_go/types"
)

func InvokeFunction(name string, args []any) (types.Value, error) {
	fn := Registry.Get(name)
	if fn != nil {
		return fn(args)
	}

	return nil, fmt.Errorf("function %s not found", name)
}
