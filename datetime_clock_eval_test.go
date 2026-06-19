package uexl_test

import (
	"context"
	"testing"

	"github.com/maniartech/uexl"
	"github.com/maniartech/uexl/types"
)

// 2024-12-01T10:30:00Z (midnight that day = 1733011200000).
const fixedNowMs = int64(1733049000000)

func TestClock_NowToday(t *testing.T) {
	env := uexl.DefaultWith(uexl.WithClock(fixedNowMs))
	eval := func(expr string) any {
		v, err := env.Eval(context.Background(), expr, nil)
		if err != nil {
			t.Fatalf("%s: %v", expr, err)
		}
		return v
	}
	trueExprs := []string{
		`now() == d"2024-12-01T10:30:00Z"`,
		`now() == now()`,           // stable within one evaluation
		`now() - now() == 0d`,      // ... so the difference is zero
		`today() == d"2024-12-01"`, // truncated to UTC midnight
		`today() <= now()`,
		`now() - today() == 37800000ms`, // 10h30m past midnight
	}
	for _, expr := range trueExprs {
		if got := eval(expr); got != true {
			t.Errorf("%s: got %v, want true", expr, got)
		}
	}
}

func TestClock_PerEvalInjection(t *testing.T) {
	env := uexl.DefaultWith() // no global clock
	for _, raw := range []any{fixedNowMs, float64(fixedNowMs), int(fixedNowMs), types.DateTime{Millis: fixedNowMs}} {
		v, err := env.Eval(context.Background(), `now() == d"2024-12-01T10:30:00Z"`, map[string]any{"$now": raw})
		if err != nil || v != true {
			t.Errorf("per-eval clock (%T): got %v, err %v", raw, v, err)
		}
	}
	// Per-eval vars shadow the env global clock.
	envG := uexl.DefaultWith(uexl.WithClock(0))
	v, err := envG.Eval(context.Background(), `now() == d"2024-12-01T10:30:00Z"`, map[string]any{"$now": fixedNowMs})
	if err != nil || v != true {
		t.Errorf("per-eval shadow: got %v, err %v", v, err)
	}
}

func TestClock_NoClockErrors(t *testing.T) {
	env := uexl.DefaultWith()
	for _, expr := range []string{`now()`, `today()`, `now() + 1d`} {
		if _, err := env.Eval(context.Background(), expr, nil); err == nil {
			t.Errorf("%s without an injected clock should error", expr)
		}
	}
	// Wrong arity and a bad clock type are errors.
	withClk := uexl.DefaultWith(uexl.WithClock(fixedNowMs))
	if _, err := withClk.Eval(context.Background(), `now(5)`, nil); err == nil {
		t.Error("now(5) should error (0 args)")
	}
	if _, err := uexl.DefaultWith().Eval(context.Background(), `now()`, map[string]any{"$now": "notatime"}); err == nil {
		t.Error("a non-numeric clock should error")
	}
}
