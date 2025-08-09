package vm

import (
	"fmt"
	"strings"

	"github.com/maniartech/uexl_go/parser"
)

// Builtins is a map of function names to their implementations.
var Builtins = VMFunctions{
	"len":      builtinLen,
	"substr":   builtinSubstr,
	"contains": builtinContains,
}

// len("abc") or len([1,2,3])
func builtinLen(args ...parser.Node) (parser.Node, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("len expects 1 argument")
	}
	switch v := args[0].(type) {
	case *parser.StringLiteral:
		return &parser.NumberLiteral{Value: float64(len(v.Value))}, nil
	case *parser.ArrayLiteral:
		return &parser.NumberLiteral{Value: float64(len(v.Elements))}, nil
	default:
		return nil, fmt.Errorf("len: unsupported type %T", args[0])
	}
}

// substr("hello", 1, 3) => "ell"
func builtinSubstr(args ...parser.Node) (parser.Node, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("substr expects 3 arguments")
	}
	str, ok := args[0].(*parser.StringLiteral)
	if !ok {
		return nil, fmt.Errorf("substr: first argument must be a string")
	}
	start, ok1 := args[1].(*parser.NumberLiteral)
	length, ok2 := args[2].(*parser.NumberLiteral)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("substr: start and length must be numbers")
	}
	s, l := int(start.Value), int(length.Value)
	if s < 0 || l < 0 || s > len(str.Value) {
		return nil, fmt.Errorf("substr: invalid start or length")
	}
	end := s + l
	if end > len(str.Value) {
		end = len(str.Value)
	}
	return &parser.StringLiteral{Value: str.Value[s:end]}, nil
}

// contains("hello", "ll") => true
func builtinContains(args ...parser.Node) (parser.Node, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("contains expects 2 arguments")
	}
	str, ok1 := args[0].(*parser.StringLiteral)
	sub, ok2 := args[1].(*parser.StringLiteral)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("contains: both arguments must be strings")
	}
	return &parser.BooleanLiteral{Value: strings.Contains(str.Value, sub.Value)}, nil
}
