// Package numbers is the UExL "math" standard-library family: abs, sign, round, floor, ceil, trunc, min,
// max, sum, avg, mod, pow, sqrt, clamp. Attach via uexl.WithMath().
package numbers

import (
	"fmt"
	"math"

	"github.com/maniartech/uexl/builtins/fn"
)

// Builtins maps function names to their implementations.
var Builtins = map[string]fn.Func{
	"abs":   unary("abs", math.Abs),
	"sign":  unary("sign", func(x float64) float64 { return float64(sign(x)) }),
	"round": unary("round", math.Round),
	"floor": unary("floor", math.Floor),
	"ceil":  unary("ceil", math.Ceil),
	"trunc": unary("trunc", math.Trunc),
	"sqrt":  unary("sqrt", math.Sqrt),
	"min":   reduceNums("min", math.Inf(1), math.Min),
	"max":   reduceNums("max", math.Inf(-1), math.Max),
	"sum":   builtinSum,
	"avg":   builtinAvg,
	"mod":   binary("mod", math.Mod),
	"pow":   binary("pow", math.Pow),
	"clamp": builtinClamp,
}

func sign(x float64) int {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}

func unary(name string, f func(float64) float64) fn.Func {
	return func(args ...any) (any, error) {
		if err := fn.Arity(name, args, 1, 1); err != nil {
			return nil, err
		}
		x, err := fn.Num(name, args, 0)
		if err != nil {
			return nil, err
		}
		return f(x), nil
	}
}

func binary(name string, f func(a, b float64) float64) fn.Func {
	return func(args ...any) (any, error) {
		if err := fn.Arity(name, args, 2, 2); err != nil {
			return nil, err
		}
		a, e1 := fn.Num(name, args, 0)
		b, e2 := fn.Num(name, args, 1)
		if e1 != nil {
			return nil, e1
		}
		if e2 != nil {
			return nil, e2
		}
		return f(a, b), nil
	}
}

// collectNums gathers numbers from either a single array argument (min([1,2,3])) or multiple numeric
// arguments (min(1, 2, 3)).
func collectNums(name string, args []any) ([]float64, error) {
	if len(args) == 1 {
		if arr, ok := args[0].([]any); ok {
			out := make([]float64, len(arr))
			for i, e := range arr {
				f, ok := e.(float64)
				if !ok {
					return nil, fmt.Errorf("%s: array element %d is not a number, got %T", name, i, e)
				}
				out[i] = f
			}
			return out, nil
		}
	}
	out := make([]float64, len(args))
	for i, a := range args {
		f, ok := a.(float64)
		if !ok {
			return nil, fmt.Errorf("%s: argument %d is not a number, got %T", name, i+1, a)
		}
		out[i] = f
	}
	return out, nil
}

func reduceNums(name string, _ float64, f func(a, b float64) float64) fn.Func {
	return func(args ...any) (any, error) {
		nums, err := collectNums(name, args)
		if err != nil {
			return nil, err
		}
		if len(nums) == 0 {
			return nil, fmt.Errorf("%s expects at least one number", name)
		}
		acc := nums[0]
		for _, x := range nums[1:] {
			acc = f(acc, x)
		}
		return acc, nil
	}
}

func builtinSum(args ...any) (any, error) {
	nums, err := collectNums("sum", args)
	if err != nil {
		return nil, err
	}
	total := 0.0
	for _, x := range nums {
		total += x
	}
	return total, nil
}

func builtinAvg(args ...any) (any, error) {
	nums, err := collectNums("avg", args)
	if err != nil {
		return nil, err
	}
	if len(nums) == 0 {
		return nil, fmt.Errorf("avg expects at least one number")
	}
	total := 0.0
	for _, x := range nums {
		total += x
	}
	return total / float64(len(nums)), nil
}

func builtinClamp(args ...any) (any, error) {
	if err := fn.Arity("clamp", args, 3, 3); err != nil {
		return nil, err
	}
	x, e1 := fn.Num("clamp", args, 0)
	lo, e2 := fn.Num("clamp", args, 1)
	hi, e3 := fn.Num("clamp", args, 2)
	for _, e := range []error{e1, e2, e3} {
		if e != nil {
			return nil, e
		}
	}
	if lo > hi {
		return nil, fmt.Errorf("clamp: lower bound %v exceeds upper bound %v", lo, hi)
	}
	return math.Max(lo, math.Min(x, hi)), nil
}
