package uexl

import (
	"github.com/maniartech/uexl/builtins/collections"
	"github.com/maniartech/uexl/builtins/conversion"
	"github.com/maniartech/uexl/builtins/datetime"
	"github.com/maniartech/uexl/builtins/fn"
	"github.com/maniartech/uexl/builtins/introspection"
	"github.com/maniartech/uexl/builtins/json"
	"github.com/maniartech/uexl/builtins/numbers"
	strs "github.com/maniartech/uexl/builtins/strings"
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

// WithStrings registers the string family beyond the core built-ins: upper, lower, trim/trimStart/
// trimEnd, replace, split, startsWith, endsWith, indexOf, repeat, padStart, padEnd.
func WithStrings() Option { return libOption(strs.Builtins) }

// WithCollections registers the collection family beyond the core built-ins: get, has, keys, values,
// remove, merge (all pure — remove/merge return new objects).
func WithCollections() Option { return libOption(collections.Builtins) }

// WithJSON registers the JSON family: parseJson (string -> value) and toJson (value -> string;
// datetime/duration serialize to ISO 8601). toJson takes an optional second bool for pretty-printing.
func WithJSON() Option { return libOption(json.Builtins) }

// WithStdlib registers every standard-library family at once — math, conversion, introspection, strings,
// collections, json, and datetime. now()/today() additionally require an injected clock (uexl.WithClock).
func WithStdlib() Option {
	all := map[string]fn.Func{}
	for _, m := range []map[string]fn.Func{
		numbers.Builtins, conversion.Builtins, introspection.Builtins,
		strs.Builtins, collections.Builtins, json.Builtins, datetime.Builtins,
	} {
		for k, v := range m {
			all[k] = v
		}
	}
	return libOption(all)
}
