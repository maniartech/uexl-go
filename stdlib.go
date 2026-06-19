package uexl

import (
	"github.com/maniartech/uexl/builtins/conversion"
	"github.com/maniartech/uexl/builtins/fn"
	"github.com/maniartech/uexl/builtins/introspection"
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

// WithConversion registers the conversion family: parseNum/tryParseNum and parseBool/tryParseBool
// (strict parsers error on failure; safe parsers return null). Value-to-string is the built-in str.
func WithConversion() Option { return libOption(conversion.Builtins) }

// WithIntrospection registers the type-introspection family: typeOf and the is* predicates (isNull,
// isNumber, isString, isBool, isArray, isObject, isDate, isDuration, isEmpty).
func WithIntrospection() Option { return libOption(introspection.Builtins) }
