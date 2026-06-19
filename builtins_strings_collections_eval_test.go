package uexl_test

import (
	"context"
	"testing"

	"github.com/maniartech/uexl"
)

func stdEnv() *uexl.Env { return uexl.DefaultWith(uexl.WithStdlib()) }

func stdTrue(t *testing.T, exprs []string) {
	t.Helper()
	env := stdEnv()
	for _, expr := range exprs {
		v, err := env.Eval(context.Background(), expr, nil)
		if err != nil || v != true {
			t.Errorf("%s: got %v, err %v", expr, v, err)
		}
	}
}

func TestStrings(t *testing.T) {
	stdTrue(t, []string{
		`upper("abc") == "ABC"`,
		`lower("ABC") == "abc"`,
		`trim("  hi  ") == "hi"`,
		`trim("xxhixx", "x") == "hi"`,
		`trimStart("  hi") == "hi"`,
		`trimEnd("hi  ") == "hi"`,
		`replace("a-b-c", "-", "_") == "a_b_c"`,
		`startsWith("hello", "he")`,
		`endsWith("hello", "lo")`,
		`!startsWith("hello", "x")`,
		`indexOf("hello", "ll") == 2`,
		`indexOf("hello", "z") == -1`,
		`repeat("ab", 3) == "ababab"`,
		`padStart("7", 3, "0") == "007"`,
		`padEnd("7", 3, "0") == "700"`,
		`padStart("toolong", 3, "0") == "toolong"`,
		`len(split("a,b,c", ",")) == 3`,
		`split("a,b,c", ",")[1] == "b"`,
	})
	env := stdEnv()
	for _, expr := range []string{`upper(5)`, `repeat("a", -1)`, `replace("a", "b")`} {
		if _, err := env.Eval(context.Background(), expr, nil); err == nil {
			t.Errorf("%s: expected an error", expr)
		}
	}
}

func TestCollections(t *testing.T) {
	stdTrue(t, []string{
		`get({"a": 1}, "a") == 1`,
		`get({"a": 1}, "z") == null`,
		`get({"a": 1}, "z", 99) == 99`,
		`get([10, 20, 30], 1) == 20`,
		`get([10], 5) == null`,
		`has({"a": 1}, "a")`,
		`!has({"a": 1}, "z")`,
		`has([1, 2], 0)`,
		`!has([1, 2], 5)`,
		`len(keys({"a": 1, "b": 2})) == 2`,
		`keys({"b": 2, "a": 1})[0] == "a"`, // sorted
		`values({"b": 2, "a": 1})[0] == 1`,
		`has(remove({"a": 1, "b": 2}, "a"), "a") == false`,
		`get(merge({"a": 1}, {"b": 2}), "b") == 2`,
		`get(merge({"a": 1}, {"a": 9}), "a") == 9`, // right wins
	})
	// merge/remove must not mutate the input.
	env := stdEnv()
	v, err := env.Eval(context.Background(), `has(remove({"a": 1}, "a"), "a")`, nil)
	if err != nil || v != false {
		t.Errorf("remove result: %v %v", v, err)
	}
	for _, expr := range []string{`get(5, "a")`, `keys([1,2])`, `merge({"a":1}, 5)`} {
		if _, err := env.Eval(context.Background(), expr, nil); err == nil {
			t.Errorf("%s: expected an error", expr)
		}
	}
}
