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
		// `as` as a member name (mirrors {as: 1} being a valid object); dot form must agree with bracket form.
		`{as: 1}.as == 1`,
		`{as: 1}.as == {as: 1}["as"]`,
		`{a: {as: 5}}.a.as == 5`, // `as` mid-chain
		`({as: 1}).as == 1`,      // grouped
	}
	for _, expr := range trueExprs {
		v, err := env.Eval(context.Background(), expr, nil)
		if err != nil || v != true {
			t.Errorf("%s: got %v, err %v", expr, v, err)
		}
	}
}
