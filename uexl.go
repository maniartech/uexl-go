package uexl_go

import (
	"github.com/maniartech/uexl_go/compiler"
	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/vm"
)

func EvalExpr(expr string) (any, error) {
	node, err := parser.ParseString(expr)
	if err != nil {
		return nil, err
	}

	comp := compiler.New()
	err = comp.Compile(node)
	if err != nil {
		return nil, err
	}

	machine := vm.New(vm.LibContext{})
	result, err := machine.Run(comp.ByteCode(), nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}
