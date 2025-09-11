package parser_test

import (
	"math"
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/parser/constants"
	"github.com/stretchr/testify/assert"
)

func TestIeeeSpecials_ExplicitlyDisabled(t *testing.T) {
	// When explicitly disabled, NaN and Inf should be parsed as identifiers
	opt := parser.DefaultOptions()
	opt.EnableIeeeSpecials = false
	tz := parser.NewTokenizerWithOptions("NaN", opt)
	tok, err := tz.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenIdentifier, tok.Type)
	assert.Equal(t, "NaN", tok.Token)

	tz2 := parser.NewTokenizerWithOptions("Inf", opt)
	tok2, err := tz2.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenIdentifier, tok2.Type)
	assert.Equal(t, "Inf", tok2.Token)
}

func TestIeeeSpecials_TokenizerEnabled(t *testing.T) {
	opt := parser.DefaultOptions()
	opt.EnableIeeeSpecials = true
	tz := parser.NewTokenizerWithOptions("NaN", opt)
	tok, err := tz.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenNumber, tok.Type)
	assert.True(t, math.IsNaN(tok.Value.Num))

	tz2 := parser.NewTokenizerWithOptions("Inf", opt)
	tok2, err := tz2.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenNumber, tok2.Type)
	assert.True(t, math.IsInf(tok2.Value.Num, +1))
}

func TestIeeeSpecials_DefaultTokenizer(t *testing.T) {
	tz := parser.NewTokenizer("NaN")
	tok, err := tz.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenNumber, tok.Type)
	assert.True(t, math.IsNaN(tok.Value.Num))

	tz2 := parser.NewTokenizer("Inf")
	tok2, err := tz2.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenNumber, tok2.Type)
	assert.True(t, math.IsInf(tok2.Value.Num, +1))
}

func TestIeeeSpecials_ParserUnaryMinus(t *testing.T) {
	opt := parser.DefaultOptions()
	opt.EnableIeeeSpecials = true
	p := parser.NewParserWithOptions("-Inf", opt)
	expr, err := p.Parse()
	assert.NoError(t, err)
	u, ok := expr.(*parser.UnaryExpression)
	if assert.True(t, ok, "expected unary expression") {
		assert.Equal(t, "-", u.Operator)
		// inner should be a number literal with +Inf; unary minus applied at evaluation time, but AST holds +Inf token
		n, ok := u.Operand.(*parser.NumberLiteral)
		assert.True(t, ok)
		assert.True(t, math.IsInf(n.Value, +1))
	}
}
