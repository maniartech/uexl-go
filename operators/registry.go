package operators

import (
	"github.com/maniartech/uexl_go/core"
	"github.com/maniartech/uexl_go/types"
)

type BinaryOpFn func(op string, a, b core.Evaluator, ctx types.Context) (types.Value, error)

type binaryOpRegistry map[string]BinaryOpFn

var BinaryOpRegistry binaryOpRegistry = binaryOpRegistry{}

func (r binaryOpRegistry) Register(name string, f BinaryOpFn) {
	r[name] = f
}

func (r binaryOpRegistry) Unregister(name string) {
	delete(r, name)
}

func (r binaryOpRegistry) Get(name string) BinaryOpFn {
	if f, ok := r[name]; ok {
		return f
	}
	return nil
}
