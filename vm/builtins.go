package vm

import (
	"fmt"
	"strings"
)

// Builtins is a map of function names to their implementations.
var Builtins = VMFunctions{
	"len":      builtinLen,
	"substr":   builtinSubstr,
	"contains": builtinContains,
	"set":      builtinSet,
	"str":      builtinStr,
}

// len("abc") or len([1,2,3])
func builtinLen(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("len expects 1 argument")
	}
	switch v := args[0].(type) {
	case string:
		return float64(len(v)), nil
	case []any:
		return float64(len(v)), nil
	default:
		return nil, fmt.Errorf("len: unsupported type %T", args[0])
	}
}

// substr("hello", 1, 3) => "ell"
func builtinSubstr(args ...any) (any, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("substr expects 3 arguments")
	}
	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("substr: first argument must be a string")
	}
	start, ok1 := args[1].(float64)
	length, ok2 := args[2].(float64)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("substr: start and length must be numbers")
	}
	s, l := start, length
	if s < 0 || l < 0 || s > float64(len(str)) {
		return nil, fmt.Errorf("substr: invalid start or length")
	}
	end := min(s+l, float64(len(str)))
	return str[int(s):int(end)], nil
}

// contains("hello", "ll") => true
func builtinContains(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("contains expects 2 arguments")
	}
	str, ok1 := args[0].(string)
	sub, ok2 := args[1].(string)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("contains: both arguments must be strings")
	}
	return strings.Contains(str, sub), nil
}

// set(obj, key, value) => obj with key set to value
func builtinSet(args ...any) (any, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("set expects 3 arguments")
	}
	obj, ok := args[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("set: first argument must be an object")
	}

	// Key should either be a string or a number. If number, it should be converted to a string.
	var key string
	switch k := args[1].(type) {
	case string:
		// Key is already a string.
		key = k
	case float64:
		// Key is a number, convert to string.
		key = fmt.Sprintf("%d", int(k))
	case int:
		// Key is an int, convert to string.
		key = fmt.Sprintf("%d", k)
	default:
		return nil, fmt.Errorf("set: key must be a string or a number")
	}

	value := args[2]
	obj[key] = value
	return obj, nil
}

// str converts a value to its string representation.
func builtinStr(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("str expects 1 argument")
	}
	return fmt.Sprintf("%v", args[0]), nil
}
