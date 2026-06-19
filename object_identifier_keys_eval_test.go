package uexl_test

import (
	"context"
	"testing"

	"github.com/maniartech/uexl"
)

// TestObjectIdentifierKeys_Eval is the end-to-end counterpart to the parser test: unquoted identifier
// object keys must compile and evaluate the same as quoted keys, including as function arguments (the
// path on which the bug was originally found, via the introspection builtins).
func TestObjectIdentifierKeys_Eval(t *testing.T) {
	env := uexl.DefaultWith(uexl.WithStdlib())
	trueExprs := []string{
		`get({a: 1}, "a") == 1`,
		`get({a:1}, "a") == 1`, // no space
		`typeOf({a:1}) == "object"`,
		`isObject({a: 1})`,
		`get({a:1, b:2}, "b") == 2`,
		`get({outer: {inner: 5}}, "outer").inner == 5`,
		`keys({a:1, b:2})[0] == "a"`,
		`{a:1}.a == 1`, // member access on identifier-keyed literal
	}
	for _, expr := range trueExprs {
		v, err := env.Eval(context.Background(), expr, nil)
		if err != nil || v != true {
			t.Errorf("%s: got %v, err %v", expr, v, err)
		}
	}
}
