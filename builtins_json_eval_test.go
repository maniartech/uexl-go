package uexl_test

import (
	"context"
	"testing"

	"github.com/maniartech/uexl"
)

func TestJSON(t *testing.T) {
	env := uexl.DefaultWith(uexl.WithStdlib())
	trueExprs := []string{
		`parseJson("42") == 42`,
		`parseJson("\"hi\"") == "hi"`,
		`parseJson("true") == true`,
		`parseJson("null") == null`,
		`parseJson("[1,2,3]")[1] == 2`,
		`get(parseJson("{\"a\": 1}"), "a") == 1`,
		`toJson(42) == "42"`,
		`toJson("hi") == "\"hi\""`,
		`toJson([1, 2]) == "[1,2]"`,
		`toJson(d"2024-12-01T10:30:00Z") == "\"2024-12-01T10:30:00Z\""`,
		`toJson(90m) == "\"PT1H30M\""`,
		`parseJson(toJson([1, 2, 3]))[2] == 3`, // round-trip
		`formatNum(3.14159, 2) == "3.14"`,
		`formatNum(42) == "42"`,
	}
	for _, expr := range trueExprs {
		v, err := env.Eval(context.Background(), expr, nil)
		if err != nil || v != true {
			t.Errorf("%s: got %v, err %v", expr, v, err)
		}
	}
	for _, expr := range []string{`parseJson("{bad json}")`, `parseJson(5)`, `formatNum("x")`, `formatNum(1, -1)`} {
		if _, err := env.Eval(context.Background(), expr, nil); err == nil {
			t.Errorf("%s: expected an error", expr)
		}
	}
}
