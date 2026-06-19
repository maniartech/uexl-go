package uexl

import (
	"github.com/maniartech/uexl/builtins/fn"
	"github.com/maniartech/uexl/builtins/numbers"
)

// libOption converts a standard-library family's builtin map into an Option that registers it. The
// family packages import only builtins/fn (never uexl or vm), so these per-family options live here in
// the root package — see builtins-implementation-plan.md §5.1.
func libOption(builtins map[string]fn.Func) Option {
	fns := make(Functions, len(builtins))
	for name, f := range builtins {
		fns[name] = Function(f)
	}
	return WithFunctions(fns)
}

// WithMath registers the math family: abs, sign, round, floor, ceil, trunc, sqrt, min, max, sum, avg,
// mod, pow, clamp. min/max/sum/avg accept either a single array or multiple numeric arguments.
func WithMath() Option { return libOption(numbers.Builtins) }
