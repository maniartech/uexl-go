package parser

import "math"

var cntCodeBlocks int

func boolToInt(isTrue bool) int {
	if isTrue {
		return 1
	}
	return 0
}

func boolCalc(a, b int) (bool, bool) {
	var aBool, bBool bool
	if a != 0 {
		aBool = true
	}
	if b != 0 {
		bBool = true
	}
	return aBool, bBool
}

var ops = map[string]func(int, int) int{
	"+": func(l, r int) int {
		return l + r
	},
	"-": func(l, r int) int {
		return l - r
	},
	"*": func(l, r int) int {
		return l * r
	},
	"/": func(l, r int) int {
		return l / r
	},
	"//": func(l, r int) int {
		return int(math.Mod(float64(l), float64(r)))
	},
	"&": func(l, r int) int {
		return l & r
	},
	"|": func(l, r int) int {
		return l | r
	},
	"^": func(l, r int) int {
		return l ^ r
	},
	"==": func(l, r int) int {
		return boolToInt(l == r)
	},
	"!=": func(l, r int) int {
		return boolToInt(l != r)
	},
	"<": func(l, r int) int {
		return boolToInt(l < r)
	},
	"<=": func(l, r int) int {
		return boolToInt(l <= r)
	},
	">": func(l, r int) int {
		return boolToInt(l > r)
	},
	">=": func(l, r int) int {
		return boolToInt(l >= r)
	},
	"&&": func(l, r int) int {
		a, b := boolCalc(l, r)
		return boolToInt(a && b)
	},
	"||": func(l, r int) int {
		a, b := boolCalc(l, r)
		return boolToInt(a || b)
	},
}

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

func eval(first, rest interface{}) int {
	l := first.(int)
	restSl := toIfaceSlice(rest)
	for _, v := range restSl {
		restExpr := toIfaceSlice(v)
		r := restExpr[3].(int)
		op := restExpr[1].(string)
		l = ops[op](l, r)
	}
	return l
}
