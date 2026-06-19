// Package conversion is the UExL conversion family: string-to-value parsing in strict (parseX, errors on
// failure) and safe (tryParseX, null on failure) pairs — parseNum/tryParseNum, parseBool/tryParseBool.
// Value-to-string conversion is the built-in str. Attach via uexl.WithConversion(). See ADR-0001 §B.
package conversion

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/maniartech/uexl/builtins/fn"
)

// Builtins maps function names to their implementations.
var Builtins = map[string]fn.Func{
	"parseNum":     builtinParseNum,
	"tryParseNum":  builtinTryParseNum,
	"parseBool":    builtinParseBool,
	"tryParseBool": builtinTryParseBool,
}

func parseNum(s string) (float64, bool) {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

func parseBool(s string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true":
		return true, true
	case "false":
		return false, true
	}
	return false, false
}

func builtinParseNum(args ...any) (any, error) {
	if err := fn.Arity("parseNum", args, 1, 1); err != nil {
		return nil, err
	}
	s, err := fn.Str("parseNum", args, 0)
	if err != nil {
		return nil, err
	}
	f, ok := parseNum(s)
	if !ok {
		return nil, fmt.Errorf("parseNum: cannot parse %q as a number", s)
	}
	return f, nil
}

func builtinTryParseNum(args ...any) (any, error) {
	if err := fn.Arity("tryParseNum", args, 1, 1); err != nil {
		return nil, err
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, nil // non-string -> null (safe)
	}
	if f, ok := parseNum(s); ok {
		return f, nil
	}
	return nil, nil
}

func builtinParseBool(args ...any) (any, error) {
	if err := fn.Arity("parseBool", args, 1, 1); err != nil {
		return nil, err
	}
	s, err := fn.Str("parseBool", args, 0)
	if err != nil {
		return nil, err
	}
	b, ok := parseBool(s)
	if !ok {
		return nil, fmt.Errorf("parseBool: cannot parse %q as a boolean (expected true/false)", s)
	}
	return b, nil
}

func builtinTryParseBool(args ...any) (any, error) {
	if err := fn.Arity("tryParseBool", args, 1, 1); err != nil {
		return nil, err
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, nil
	}
	if b, ok := parseBool(s); ok {
		return b, nil
	}
	return nil, nil
}
