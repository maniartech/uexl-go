// Package fn holds shared argument helpers for the UExL standard-library families (numbers, conversion,
// introspection, strings, collections, json). A builtin is the universal func(args ...any) (any, error);
// these helpers kill per-function boilerplate and give consistent error messages.
package fn

import (
	"fmt"
	"math"
)

// Func is the universal builtin signature (identical to vm.VMFunction).
type Func = func(args ...any) (any, error)

// Arity validates the argument count. Pass max = -1 for "variadic / no upper bound".
func Arity(name string, args []any, min, max int) error {
	n := len(args)
	if n < min || (max >= 0 && n > max) {
		switch {
		case min == max:
			return fmt.Errorf("%s expects %d argument(s), got %d", name, min, n)
		case max < 0:
			return fmt.Errorf("%s expects at least %d argument(s), got %d", name, min, n)
		default:
			return fmt.Errorf("%s expects %d..%d arguments, got %d", name, min, max, n)
		}
	}
	return nil
}

// Num extracts a float64 argument.
func Num(name string, args []any, i int) (float64, error) {
	f, ok := args[i].(float64)
	if !ok {
		return 0, fmt.Errorf("%s: argument %d must be a number, got %T", name, i+1, args[i])
	}
	return f, nil
}

// Int extracts an integer-valued float64 argument.
func Int(name string, args []any, i int) (int, error) {
	f, err := Num(name, args, i)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(f) || math.IsInf(f, 0) || f != math.Trunc(f) {
		return 0, fmt.Errorf("%s: argument %d must be an integer, got %v", name, i+1, f)
	}
	return int(f), nil
}

// Str extracts a string argument.
func Str(name string, args []any, i int) (string, error) {
	s, ok := args[i].(string)
	if !ok {
		return "", fmt.Errorf("%s: argument %d must be a string, got %T", name, i+1, args[i])
	}
	return s, nil
}

// Bool extracts a boolean argument.
func Bool(name string, args []any, i int) (bool, error) {
	b, ok := args[i].(bool)
	if !ok {
		return false, fmt.Errorf("%s: argument %d must be a boolean, got %T", name, i+1, args[i])
	}
	return b, nil
}

// Arr extracts an array ([]any) argument.
func Arr(name string, args []any, i int) ([]any, error) {
	a, ok := args[i].([]any)
	if !ok {
		return nil, fmt.Errorf("%s: argument %d must be an array, got %T", name, i+1, args[i])
	}
	return a, nil
}

// Obj extracts an object (map[string]any) argument.
func Obj(name string, args []any, i int) (map[string]any, error) {
	m, ok := args[i].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: argument %d must be an object, got %T", name, i+1, args[i])
	}
	return m, nil
}
