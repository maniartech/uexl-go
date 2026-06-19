package uexl_test

import (
	"context"
	"testing"

	"github.com/maniartech/uexl"
)

func convEnv() *uexl.Env {
	return uexl.DefaultWith(uexl.WithConversion(), uexl.WithIntrospection(), uexl.WithDatetime())
}

func TestConversion(t *testing.T) {
	env := convEnv()
	trueExprs := []string{
		`parseNum("42") == 42`,
		`parseNum("3.14") == 3.14`,
		`parseNum("  10  ") == 10`,
		`tryParseNum("42") == 42`,
		`tryParseNum("abc") == null`,
		`tryParseNum(5) == null`, // non-string -> null
		`parseBool("true") == true`,
		`parseBool("FALSE") == false`,
		`tryParseBool("nope") == null`,
		`tryParseBool("true") == true`,
	}
	for _, expr := range trueExprs {
		v, err := env.Eval(context.Background(), expr, nil)
		if err != nil || v != true {
			t.Errorf("%s: got %v, err %v", expr, v, err)
		}
	}
	for _, expr := range []string{`parseNum("abc")`, `parseBool("maybe")`, `parseNum(5)`, `parseNum()`} {
		if _, err := env.Eval(context.Background(), expr, nil); err == nil {
			t.Errorf("%s: expected an error", expr)
		}
	}
}

func TestIntrospection(t *testing.T) {
	env := convEnv()
	trueExprs := []string{
		`typeOf(42) == "number"`,
		`typeOf("x") == "string"`,
		`typeOf(true) == "boolean"`,
		`typeOf(null) == "null"`,
		`typeOf([1,2]) == "array"`,
		`typeOf({"a": 1}) == "object"`,
		`typeOf(d"2024-12-01") == "datetime"`,
		`typeOf(7d) == "duration"`,
		`isNull(null)`,
		`!isNull(0)`,
		`isNumber(3.14)`,
		`isString("x")`,
		`isBool(false)`,
		`isArray([1])`,
		`isObject({"a": 1})`,
		`isDate(d"2024-12-01")`,
		`isDuration(7d)`,
		`isEmpty("")`,
		`isEmpty([])`,
		`isEmpty({})`,
		`isEmpty(null)`,
		`!isEmpty("x")`,
		`!isEmpty(0)`,
	}
	for _, expr := range trueExprs {
		v, err := env.Eval(context.Background(), expr, nil)
		if err != nil || v != true {
			t.Errorf("%s: got %v, err %v", expr, v, err)
		}
	}
	if _, err := env.Eval(context.Background(), `typeOf()`, nil); err == nil {
		t.Error("typeOf() should error (arity)")
	}
}
