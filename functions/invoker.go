package functions

import "fmt"

func InvokeFunction(name string, args []any) (any, error) {
	fn := Registry.Get(name)
	if fn != nil {
		return fn(args)
	}

	return nil, fmt.Errorf("function %s not found", name)
}
