package pipes

import (
	"github.com/maniartech/uexl_go/evaluators"
	"github.com/maniartech/uexl_go/types"
)

type PipeHandler func(evaluator evaluators.Evaluator, context types.Map, prevResult any) (interface{}, error)
type handlers map[string]PipeHandler

var _handlers = handlers{}

func (r handlers) Register(name string, handler PipeHandler) {
	r[name] = handler
}

func (r handlers) Unregister(name string) {
	delete(r, name)
}

func Get(name string) (PipeHandler, bool) {
	handler, ok := _handlers[name]
	return handler, ok
}
