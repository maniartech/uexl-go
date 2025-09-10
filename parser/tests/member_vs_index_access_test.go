package parser_test

import (
	"testing"

	p "github.com/maniartech/uexl_go/parser"
	"github.com/stretchr/testify/assert"
)

// These tests assert that after an index access (arr[0] or arr.0),
// subsequent .address or .address.street are parsed as MemberAccess nodes.

func TestIndexThenMemberAccess_BracketThenDot(t *testing.T) {
	expr, err := p.ParseString("arr[0].address")
	assert.NoError(t, err)

	ma, ok := expr.(*p.MemberAccess)
	if assert.True(t, ok, "expected MemberAccess at top-level") {
		assert.Equal(t, "address", ma.Property.S)
		_, isIndex := ma.Target.(*p.IndexAccess)
		assert.True(t, isIndex, "object of member should be IndexAccess (arr[0])")
	}
}

func TestIndexThenMemberAccess_Chain_BracketThenDot(t *testing.T) {
	expr, err := p.ParseString("arr[0].address.street")
	assert.NoError(t, err)

	// Top-level member .street
	maStreet, ok := expr.(*p.MemberAccess)
	if assert.True(t, ok, "expected MemberAccess for .street") {
		assert.Equal(t, "street", maStreet.Property.S)
		// Its object should be another MemberAccess .address
		maAddress, ok := maStreet.Target.(*p.MemberAccess)
		if assert.True(t, ok, "expected inner MemberAccess for .address") {
			assert.Equal(t, "address", maAddress.Property.S)
			_, isIndex := maAddress.Target.(*p.IndexAccess)
			assert.True(t, isIndex, "object of .address should be IndexAccess (arr[0])")
		}
	}
}

func TestIndexThenMemberAccess_DotNumberThenDot(t *testing.T) {
	expr, err := p.ParseString("arr.0.address")
	assert.NoError(t, err)

	ma, ok := expr.(*p.MemberAccess)
	if assert.True(t, ok, "expected MemberAccess at top-level") {
		assert.Equal(t, "address", ma.Property.S)
		inner, isMember := ma.Target.(*p.MemberAccess)
		assert.True(t, isMember, "object of member should be MemberAccess (arr.0)")
		if isMember {
			// inner should have integer property 0
			assert.True(t, inner.Property.IsInt(), "dot-number should produce integer property")
			assert.Equal(t, 0, inner.Property.I)
		}
	}
}

func TestIndexThenMemberAccess_Chain_DotNumberThenDot(t *testing.T) {
	expr, err := p.ParseString("arr.0.address.street")
	assert.NoError(t, err)

	// Top-level member .street
	maStreet, ok := expr.(*p.MemberAccess)
	if assert.True(t, ok, "expected MemberAccess for .street") {
		assert.Equal(t, "street", maStreet.Property.S)
		// Its object should be another MemberAccess .address
		maAddress, ok := maStreet.Target.(*p.MemberAccess)
		if assert.True(t, ok, "expected inner MemberAccess for .address") {
			assert.Equal(t, "address", maAddress.Property.S)
			inner, isMember := maAddress.Target.(*p.MemberAccess)
			assert.True(t, isMember, "object of .address should be MemberAccess (arr.0)")
			if isMember {
				assert.True(t, inner.Property.IsInt(), "dot-number should produce integer property")
				assert.Equal(t, 0, inner.Property.I)
			}
		}
	}
}
