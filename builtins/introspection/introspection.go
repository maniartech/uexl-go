// Package introspection is the UExL type-introspection family: typeOf and the is* predicates. Attach via
// uexl.WithIntrospection().
package introspection

import (
	"github.com/maniartech/uexl/builtins/fn"
	"github.com/maniartech/uexl/types"
)

// Builtins maps function names to their implementations.
var Builtins = map[string]fn.Func{
	"typeOf":     builtinTypeOf,
	"isNull":     predicate("isNull", func(v any) bool { return v == nil }),
	"isNumber":   predicate("isNumber", isNumber),
	"isString":   predicate("isString", func(v any) bool { _, ok := v.(string); return ok }),
	"isBool":     predicate("isBool", func(v any) bool { _, ok := v.(bool); return ok }),
	"isArray":    predicate("isArray", func(v any) bool { _, ok := v.([]any); return ok }),
	"isObject":   predicate("isObject", func(v any) bool { _, ok := v.(map[string]any); return ok }),
	"isDate":     predicate("isDate", func(v any) bool { _, ok := v.(types.DateTime); return ok }),
	"isDuration": predicate("isDuration", func(v any) bool { _, ok := v.(types.Duration); return ok }),
	"isEmpty":    predicate("isEmpty", isEmpty),
}

func isNumber(v any) bool {
	switch v.(type) {
	case float64, int, int64:
		return true
	}
	return false
}

func isEmpty(v any) bool {
	switch x := v.(type) {
	case nil:
		return true
	case string:
		return x == ""
	case []any:
		return len(x) == 0
	case map[string]any:
		return len(x) == 0
	}
	return false
}

// typeName returns the UExL type name of a value (datetime-spec / introspection vocabulary).
func typeName(v any) string {
	switch v.(type) {
	case nil:
		return "null"
	case float64, int, int64:
		return "number"
	case string:
		return "string"
	case bool:
		return "boolean"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	case types.DateTime:
		return "datetime"
	case types.Duration:
		return "duration"
	default:
		return "unknown"
	}
}

func builtinTypeOf(args ...any) (any, error) {
	if err := fn.Arity("typeOf", args, 1, 1); err != nil {
		return nil, err
	}
	return typeName(args[0]), nil
}

func predicate(name string, test func(any) bool) fn.Func {
	return func(args ...any) (any, error) {
		if err := fn.Arity(name, args, 1, 1); err != nil {
			return nil, err
		}
		return test(args[0]), nil
	}
}
