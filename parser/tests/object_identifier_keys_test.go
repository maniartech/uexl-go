package parser_test

import (
	"testing"

	"github.com/maniartech/uexl/parser"
	"github.com/stretchr/testify/assert"
)

// TestObjectIdentifierKeys verifies that bareword identifier keys parse identically to quoted string
// keys: {a: 1} and {a:1} (no space) are sugar for {"a": 1}. Regression test for the parser rejecting
// unquoted object keys (it required a TokenString key), which broke documented forms like {id: 123} and
// object arguments to functions such as typeOf({a:1}) / isObject({a:1}).
func TestObjectIdentifierKeys(t *testing.T) {
	// Each input must parse to an ObjectLiteral whose key `wantKey` holds the NumberLiteral `wantVal`.
	cases := []struct {
		name    string
		input   string
		wantKey string
		wantVal float64
	}{
		{"ident no space", `{a:1}`, "a", 1},
		{"ident with space", `{a: 1}`, "a", 1},
		{"quoted no space", `{"a":1}`, "a", 1},
		{"quoted with space", `{"a": 1}`, "a", 1},
		{"ident underscore", `{a_1: 7}`, "a_1", 7},
		{"ident dollar", `{$x: 8}`, "$x", 8},
		{"ident leading underscore", `{_priv: 9}`, "_priv", 9},
		{"ident alnum", `{abc123: 3}`, "abc123", 3},
		{"contextual keyword as", `{as: 4}`, "as", 4},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			expr, err := parser.NewParser(c.input).Parse()
			assert.NoError(t, err, "parse %q", c.input)
			obj, ok := expr.(*parser.ObjectLiteral)
			assert.True(t, ok, "%q -> %T, want ObjectLiteral", c.input, expr)
			val, exists := obj.Properties[c.wantKey]
			assert.True(t, exists, "%q: key %q missing", c.input, c.wantKey)
			num, ok := val.(*parser.NumberLiteral)
			assert.True(t, ok, "%q: value -> %T, want NumberLiteral", c.input, val)
			if ok {
				assert.Equal(t, c.wantVal, num.Value, "%q: value", c.input)
			}
		})
	}
}

// TestObjectIdentifierKeys_Mixed covers multiple/mixed string+identifier keys, nesting, ternary values
// (whose `:` must not be mistaken for a key separator), and function-argument position.
func TestObjectIdentifierKeys_Mixed(t *testing.T) {
	t.Run("mixed keys", func(t *testing.T) {
		expr, err := parser.NewParser(`{a:1, "b":2, c:3}`).Parse()
		assert.NoError(t, err)
		obj, ok := expr.(*parser.ObjectLiteral)
		assert.True(t, ok)
		assert.Len(t, obj.Properties, 3)
		for _, k := range []string{"a", "b", "c"} {
			_, exists := obj.Properties[k]
			assert.True(t, exists, "key %q should exist", k)
		}
	})

	t.Run("nested identifier keys", func(t *testing.T) {
		expr, err := parser.NewParser(`{outer: {inner: 5}}`).Parse()
		assert.NoError(t, err)
		obj, ok := expr.(*parser.ObjectLiteral)
		assert.True(t, ok)
		inner, exists := obj.Properties["outer"]
		assert.True(t, exists)
		innerObj, ok := inner.(*parser.ObjectLiteral)
		assert.True(t, ok, "nested value should be ObjectLiteral")
		_, exists = innerObj.Properties["inner"]
		assert.True(t, exists)
	})

	t.Run("ternary value does not break key", func(t *testing.T) {
		// The ternary `:` inside the value must be consumed by the value expression, not the object.
		expr, err := parser.NewParser(`{x: true ? 1 : 2}`).Parse()
		assert.NoError(t, err)
		obj, ok := expr.(*parser.ObjectLiteral)
		assert.True(t, ok)
		assert.Len(t, obj.Properties, 1)
		_, exists := obj.Properties["x"]
		assert.True(t, exists)
	})

	t.Run("object as function argument", func(t *testing.T) {
		for _, in := range []string{`typeOf({a:1})`, `typeOf({a: 1})`, `isObject({k: v})`} {
			expr, err := parser.NewParser(in).Parse()
			assert.NoError(t, err, "parse %q", in)
			_, ok := expr.(*parser.FunctionCall)
			assert.True(t, ok, "%q -> %T, want FunctionCall", in, expr)
		}
	})
}

// TestObjectInvalidKeys confirms that non-identifier, non-string key tokens are still rejected, so the
// fix only widens keys to identifiers (not arbitrary value literals).
func TestObjectInvalidKeys(t *testing.T) {
	for _, in := range []string{
		`{1: 2}`,    // number key
		`{true: 1}`, // boolean literal key
		`{null: 1}`, // null literal key
		`{a 1}`,     // missing colon
		`{a:}`,      // missing value
		`{a:1 b:2}`, // missing comma between pairs
	} {
		_, err := parser.NewParser(in).Parse()
		assert.Error(t, err, "expected a parse error for %q", in)
	}
}
