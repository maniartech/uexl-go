package uexl_test

import (
	"context"
	"testing"

	"github.com/maniartech/uexl"
)

func mathEval(t *testing.T, expr string) any {
	t.Helper()
	v, err := uexl.DefaultWith(uexl.WithMath()).Eval(context.Background(), expr, nil)
	if err != nil {
		t.Fatalf("%s: %v", expr, err)
	}
	return v
}

func TestMath_Functions(t *testing.T) {
	num := map[string]float64{
		`abs(-5)`:           5,
		`abs(5)`:            5,
		`sign(-3)`:          -1,
		`sign(0)`:           0,
		`sign(7)`:           1,
		`round(2.5)`:        3,
		`round(2.4)`:        2,
		`floor(2.9)`:        2,
		`ceil(2.1)`:         3,
		`trunc(-2.9)`:       -2,
		`sqrt(9)`:           3,
		`min(3, 1, 2)`:      1,
		`max(3, 1, 2)`:      3,
		`min([3, 1, 2])`:    1,
		`max([3, 1, 2])`:    3,
		`sum(1, 2, 3)`:      6,
		`sum([1, 2, 3])`:    6,
		`sum()`:             0,
		`avg(2, 4)`:         3,
		`avg([1, 2, 3, 4])`: 2.5,
		`mod(7, 3)`:         1,
		`pow(2, 10)`:        1024,
		`clamp(5, 0, 3)`:    3,
		`clamp(-1, 0, 3)`:   0,
		`clamp(2, 0, 3)`:    2,
	}
	for expr, want := range num {
		if got, ok := mathEval(t, expr).(float64); !ok || got != want {
			t.Errorf("%s = %v, want %v", expr, mathEval(t, expr), want)
		}
	}
}

func TestMath_Errors(t *testing.T) {
	env := uexl.DefaultWith(uexl.WithMath())
	for _, expr := range []string{
		`abs("x")`,
		`abs()`,
		`abs(1, 2)`,
		`min()`,
		`avg()`,
		`min([1, "x"])`,
		`mod(1)`,
		`clamp(1, 3, 0)`, // lower > upper
		`sum(1, "x")`,
	} {
		if _, err := env.Eval(context.Background(), expr, nil); err == nil {
			t.Errorf("%s: expected an error", expr)
		}
	}
}
