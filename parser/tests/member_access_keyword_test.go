package parser_test

import (
	"testing"

	"github.com/maniartech/uexl/parser"
	"github.com/stretchr/testify/assert"
)

// TestMemberAccessKeywordAs verifies that the contextual keyword `as` is accepted as a member name after
// `.` (and `?.`), mirroring the object-key rule that {as: 1} is a valid object. Previously `{as:1}.as`
// parse-errored ("expected identifier after .") even though the bracket form `{as:1}["as"]` worked.
func TestMemberAccessKeywordAs(t *testing.T) {
	// inputs that must parse to a MemberAccess whose property is "as".
	memberCases := []struct {
		name  string
		input string
	}{
		{"object literal dot as", `{as: 1}.as`},
		{"object literal dot as no space", `{as:1}.as`},
		{"identifier dot as", `obj.as`},
		{"optional dot as", `obj?.as`},
		{"grouped dot as", `({as: 1}).as`},
	}
	for _, c := range memberCases {
		t.Run(c.name, func(t *testing.T) {
			expr, err := parser.NewParser(c.input).Parse()
			assert.NoError(t, err, "parse %q", c.input)
			ma, ok := expr.(*parser.MemberAccess)
			assert.True(t, ok, "%q -> %T, want MemberAccess", c.input, expr)
			if ok {
				assert.Equal(t, "as", ma.Property.S, "%q: property name", c.input)
			}
		})
	}

	// chained access with `as` in the middle: obj.as.value
	t.Run("chained as in middle", func(t *testing.T) {
		expr, err := parser.NewParser(`obj.as.value`).Parse()
		assert.NoError(t, err)
		outer, ok := expr.(*parser.MemberAccess)
		assert.True(t, ok, "outer should be MemberAccess")
		if ok {
			assert.Equal(t, "value", outer.Property.S)
			inner, ok := outer.Target.(*parser.MemberAccess)
			assert.True(t, ok, "inner should be MemberAccess")
			if ok {
				assert.Equal(t, "as", inner.Property.S)
			}
		}
	})

	// optional-chain property name is `as`.
	t.Run("optional chain as", func(t *testing.T) {
		expr, err := parser.NewParser(`obj?.as`).Parse()
		assert.NoError(t, err)
		ma, ok := expr.(*parser.MemberAccess)
		assert.True(t, ok)
		if ok {
			assert.True(t, ma.Optional, "?. should set Optional")
			assert.Equal(t, "as", ma.Property.S)
		}
	})
}

// TestMemberAccessAs_NoPipeAliasRegression ensures the change does not disturb `as` in its original role
// as a pipe alias, including the tricky combination of member access `.as` immediately followed by an
// `as` alias in the same pipe stage.
func TestMemberAccessAs_NoPipeAliasRegression(t *testing.T) {
	for _, in := range []string{
		`arr |map: $item as $x`,          // `as` is a pipe alias (unchanged)
		`[{as: 1}] |map: $item.as`,       // `.as` is member access inside a pipe body
		`[{as: 1}] |map: $item.as as $a`, // member access `.as` THEN an `as` alias — both roles, one stage
		`obj.name`,                       // ordinary identifier member access still works
		`obj.0`,                          // numeric member access still works
		`obj.as["x"]`,                    // dot-as then bracket
	} {
		expr, err := parser.NewParser(in).Parse()
		assert.NoError(t, err, "parse %q -> %v", in, err)
		assert.NotNil(t, expr, "parse %q produced nil", in)
	}
}
