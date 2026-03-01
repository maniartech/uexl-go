package vm

import (
	"fmt"
	"strings"

	"github.com/maniartech/uexl/internal/utils"
)

// Builtins is a map of function names to their implementations.
var Builtins = VMFunctions{
	"len":      builtinLen,
	"substr":   builtinSubstr,
	"contains": builtinContains,
	"set":      builtinSet,
	"str":      builtinStr,

	// Rune-level
	"runeLen":    builtinRuneLen,
	"runeSubstr": builtinRuneSubstr,

	// Grapheme-level (UAX #29)
	"graphemeLen":    builtinGraphemeLen,
	"graphemeSubstr": builtinGraphemeSubstr,

	// Explode
	"runes":     builtinRunes,
	"graphemes": builtinGraphemes,
	"bytes":     builtinBytes,

	// Reassemble
	"join": builtinJoin,
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

// ---- helpers ----------------------------------------------------------------

func requireOneString(name string, args []any) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("%s expects 1 argument", name)
	}
	s, ok := args[0].(string)
	if !ok {
		return "", fmt.Errorf("%s: argument must be a string, got %T", name, args[0])
	}
	return s, nil
}

func requireStringIntInt(name string, args []any) (string, int, int, error) {
	if len(args) != 3 {
		return "", 0, 0, fmt.Errorf("%s expects 3 arguments", name)
	}
	s, ok := args[0].(string)
	if !ok {
		return "", 0, 0, fmt.Errorf("%s: first argument must be a string, got %T", name, args[0])
	}
	startF, ok1 := args[1].(float64)
	lenF, ok2 := args[2].(float64)
	if !ok1 || !ok2 {
		return "", 0, 0, fmt.Errorf("%s: start and length must be numbers", name)
	}
	start := int(startF)
	if float64(start) != startF {
		return "", 0, 0, fmt.Errorf("%s: start must be an integer, got %g", name, startF)
	}
	n := int(lenF)
	if float64(n) != lenF {
		return "", 0, 0, fmt.Errorf("%s: length must be an integer, got %g", name, lenF)
	}
	return s, start, n, nil
}

// ---- Rune-level -------------------------------------------------------------

// runeLen(s) — number of Unicode code points in s.
func builtinRuneLen(args ...any) (any, error) {
	s, err := requireOneString("runeLen", args)
	if err != nil {
		return nil, err
	}
	return float64(utils.RuneLength(s)), nil
}

// runeSubstr(s, start, length) — substring by rune indices.
func builtinRuneSubstr(args ...any) (any, error) {
	s, start, length, err := requireStringIntInt("runeSubstr", args)
	if err != nil {
		return nil, err
	}
	return utils.RuneSlice(s, start, length)
}

// ---- Grapheme-level (UAX #29) -----------------------------------------------

// graphemeLen(s) — number of user-perceived grapheme clusters.
func builtinGraphemeLen(args ...any) (any, error) {
	s, err := requireOneString("graphemeLen", args)
	if err != nil {
		return nil, err
	}
	return float64(utils.GraphemeLength(s)), nil
}

// graphemeSubstr(s, start, length) — substring by grapheme cluster indices.
func builtinGraphemeSubstr(args ...any) (any, error) {
	s, start, length, err := requireStringIntInt("graphemeSubstr", args)
	if err != nil {
		return nil, err
	}
	return utils.GraphemeSlice(s, start, length)
}

// ---- Explode ----------------------------------------------------------------

// runes(s) — []any of single-rune strings.
func builtinRunes(args ...any) (any, error) {
	s, err := requireOneString("runes", args)
	if err != nil {
		return nil, err
	}
	return utils.CollectRunes(s), nil
}

// graphemes(s) — []any of grapheme cluster strings.
func builtinGraphemes(args ...any) (any, error) {
	s, err := requireOneString("graphemes", args)
	if err != nil {
		return nil, err
	}
	return utils.CollectGraphemes(s), nil
}

// bytes(s) — []any of byte values as float64.
func builtinBytes(args ...any) (any, error) {
	s, err := requireOneString("bytes", args)
	if err != nil {
		return nil, err
	}
	return utils.CollectBytes(s), nil
}

// ---- Reassemble -------------------------------------------------------------

// join(arr) or join(arr, sep) — concatenate []any of strings with optional separator.
func builtinJoin(args ...any) (any, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("join expects 1 or 2 arguments")
	}
	arr, ok := args[0].([]any)
	if !ok {
		return nil, fmt.Errorf("join: first argument must be an array, got %T", args[0])
	}
	sep := ""
	if len(args) == 2 {
		sep, ok = args[1].(string)
		if !ok {
			return nil, fmt.Errorf("join: separator must be a string, got %T", args[1])
		}
	}
	var sb strings.Builder
	sb.Grow(len(arr) * 4)
	for i, v := range arr {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("join: element %d must be a string, got %T", i, v)
		}
		if i > 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(s)
	}
	return sb.String(), nil
}
