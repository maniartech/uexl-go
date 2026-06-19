// Package collections is the UExL collection family beyond the core built-ins (len/set): get, has, keys,
// values, remove, merge. All are pure — remove/merge return new objects, never mutating the input.
// Attach via uexl.WithCollections().
package collections

import (
	"fmt"
	"sort"

	"github.com/maniartech/uexl/builtins/fn"
)

// Builtins maps function names to their implementations.
var Builtins = map[string]fn.Func{
	"get":    builtinGet,
	"has":    builtinHas,
	"keys":   builtinKeys,
	"values": builtinValues,
	"remove": builtinRemove,
	"merge":  builtinMerge,
}

// get(container, key[, default]) returns the value at key (object) or index (array), or default (or null)
// when absent / out of bounds.
func builtinGet(args ...any) (any, error) {
	if err := fn.Arity("get", args, 2, 3); err != nil {
		return nil, err
	}
	var def any
	if len(args) == 3 {
		def = args[2]
	}
	switch c := args[0].(type) {
	case map[string]any:
		key, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("get: object key must be a string, got %T", args[1])
		}
		if v, ok := c[key]; ok {
			return v, nil
		}
		return def, nil
	case []any:
		idx, err := fn.Int("get", args, 1)
		if err != nil {
			return nil, err
		}
		if idx >= 0 && idx < len(c) {
			return c[idx], nil
		}
		return def, nil
	default:
		return nil, fmt.Errorf("get: first argument must be an object or array, got %T", args[0])
	}
}

func builtinHas(args ...any) (any, error) {
	if err := fn.Arity("has", args, 2, 2); err != nil {
		return nil, err
	}
	switch c := args[0].(type) {
	case map[string]any:
		key, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("has: object key must be a string, got %T", args[1])
		}
		_, exists := c[key]
		return exists, nil
	case []any:
		idx, err := fn.Int("has", args, 1)
		if err != nil {
			return nil, err
		}
		return idx >= 0 && idx < len(c), nil
	default:
		return nil, fmt.Errorf("has: first argument must be an object or array, got %T", args[0])
	}
}

func sortedKeys(m map[string]any) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks) // deterministic order (Go map iteration is randomized)
	return ks
}

func builtinKeys(args ...any) (any, error) {
	if err := fn.Arity("keys", args, 1, 1); err != nil {
		return nil, err
	}
	m, err := fn.Obj("keys", args, 0)
	if err != nil {
		return nil, err
	}
	ks := sortedKeys(m)
	out := make([]any, len(ks))
	for i, k := range ks {
		out[i] = k
	}
	return out, nil
}

func builtinValues(args ...any) (any, error) {
	if err := fn.Arity("values", args, 1, 1); err != nil {
		return nil, err
	}
	m, err := fn.Obj("values", args, 0)
	if err != nil {
		return nil, err
	}
	ks := sortedKeys(m)
	out := make([]any, len(ks))
	for i, k := range ks {
		out[i] = m[k]
	}
	return out, nil
}

// remove(obj, key) returns a new object without key (the input is not mutated).
func builtinRemove(args ...any) (any, error) {
	if err := fn.Arity("remove", args, 2, 2); err != nil {
		return nil, err
	}
	m, e1 := fn.Obj("remove", args, 0)
	key, e2 := fn.Str("remove", args, 1)
	if e1 != nil {
		return nil, e1
	}
	if e2 != nil {
		return nil, e2
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		if k != key {
			out[k] = v
		}
	}
	return out, nil
}

// merge(a, b) returns a new object with b's entries layered over a's (the inputs are not mutated).
func builtinMerge(args ...any) (any, error) {
	if err := fn.Arity("merge", args, 2, 2); err != nil {
		return nil, err
	}
	a, e1 := fn.Obj("merge", args, 0)
	b, e2 := fn.Obj("merge", args, 1)
	if e1 != nil {
		return nil, e1
	}
	if e2 != nil {
		return nil, e2
	}
	out := make(map[string]any, len(a)+len(b))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out, nil
}
