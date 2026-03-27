package parser_test

import (
	"testing"

	"github.com/maniartech/uexl/parser"
	"github.com/maniartech/uexl/parser/errors"
	"github.com/stretchr/testify/assert"
)

// parsePipeArgs is a helper that parses an expression and returns the PipeExpressions from the ProgramNode.
func parsePipeArgs(t *testing.T, input string) ([]parser.PipeExpression, error) {
	t.Helper()
	p := parser.NewParser(input)
	ast, err := p.Parse()
	if err != nil {
		return nil, err
	}
	program, ok := ast.(*parser.ProgramNode)
	if !ok {
		t.Fatalf("expected *parser.ProgramNode, got %T", ast)
	}
	return program.PipeExpressions, nil
}

// assertParseError checks that parsing input produces an error containing the given error code.
func assertParseError(t *testing.T, input string, code errors.ErrorCode) {
	t.Helper()
	p := parser.NewParser(input)
	_, err := p.Parse()
	assert.Error(t, err, "expected parse error for: %s", input)
	if parseErr, ok := err.(*errors.ParseErrors); ok {
		assert.True(t, parseErr.HasErrorCode(code),
			"expected error code %s for %q, got: %v", code, input, parseErr.Errors)
	} else if parserErr, ok := err.(errors.ParserError); ok {
		assert.Equal(t, code, parserErr.Code,
			"expected error code %s for %q, got: %s", code, input, parserErr.Code)
	} else {
		t.Fatalf("unexpected error type %T: %v", err, err)
	}
}

// TestPipeParams_HappyPath tests successful parsing of pipe arg syntax.
func TestPipeParams_HappyPath(t *testing.T) {
	tests := []struct {
		input     string
		wantType  string
		wantArgs  []any
		wantAlias string
	}{
		// Single number arg
		{`[1,2,3] |window(3): $window`, "window", []any{float64(3)}, ""},
		// Single chunk arg
		{`[1,2,3,4,5] |chunk(4): $chunk`, "chunk", []any{float64(4)}, ""},
		// String arg
		{`arr |sort("desc"): $item`, "sort", []any{"desc"}, ""},
		// Boolean arg
		{`arr |myPipe(true): $item`, "myPipe", []any{true}, ""},
		// Null arg
		{`arr |myPipe(null): $item`, "myPipe", []any{nil}, ""},
		// Multiple args
		{`arr |myPipe(3, "asc", true): $item`, "myPipe", []any{float64(3), "asc", true}, ""},
		// Float arg
		{`arr |myPipe(3.14): $item`, "myPipe", []any{float64(3.14)}, ""},
		// Negative number arg — parsed as unary minus applied to literal
		// The parser tokenizes "-3" as a number literal with value -3
		// (may vary by implementation — test passes if either the sign is included or predicate starts with it)
		// → skipped; covered by VM integration tests with full evaluation
		// Empty parens — treated as nil args
		{`arr |window(): $window`, "window", nil, ""},
		// No args — backward compat
		{`[1,2,3] |window: $window`, "window", nil, ""},
		{`[1,2] |map: $item * 2`, "map", nil, ""},
		// Whitespace variants — all spaces between tokens are non-significant
		{`arr | map : $item * 2`, "map", nil, ""},                                    // space between | and name, name and ':'
		{`arr | window ( 3 ) : $window`, "window", []any{float64(3)}, ""},            // space everywhere, including before '('
		{`arr |window( 3 , 4 ): $item`, "window", []any{float64(3), float64(4)}, ""}, // spaces inside args
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			pipes, err := parsePipeArgs(t, tt.input)
			assert.NoError(t, err)
			assert.NotEmpty(t, pipes)
			// Find the pipe segment matching the expected type
			var found *parser.PipeExpression
			for i := range pipes[1:] { // pipes[0] is the source expression segment
				if pipes[i+1].PipeType == tt.wantType {
					found = &pipes[i+1]
					break
				}
			}
			if found == nil {
				t.Fatalf("no pipe segment with type %q found in %v", tt.wantType, pipes)
			}
			assert.Equal(t, tt.wantArgs, found.Args,
				"pipe %q args mismatch for input %q", tt.wantType, tt.input)
			if tt.wantAlias != "" {
				assert.Equal(t, tt.wantAlias, found.Alias,
					"pipe %q alias mismatch for input %q", tt.wantType, tt.input)
			}
		})
	}
}

// TestPipeParams_ChainedPipes verifies that multiple chained pipes each get their own args.
func TestPipeParams_ChainedPipes(t *testing.T) {
	t.Run("only second pipe has args", func(t *testing.T) {
		pipes, err := parsePipeArgs(t, "arr |filter: $item > 0 |window(3): $window")
		assert.NoError(t, err)
		// pipes[0] = source (arr), pipes[1] = filter, pipes[2] = window
		assert.Len(t, pipes, 3)
		assert.Nil(t, pipes[1].Args, "filter pipe should have nil args")
		assert.Equal(t, []any{float64(3)}, pipes[2].Args, "window pipe should have args [3]")
	})

	t.Run("both chained pipes have args", func(t *testing.T) {
		pipes, err := parsePipeArgs(t, "arr |window(3): $window |chunk(2): $chunk")
		assert.NoError(t, err)
		assert.Len(t, pipes, 3)
		assert.Equal(t, []any{float64(3)}, pipes[1].Args, "window pipe args")
		assert.Equal(t, []any{float64(2)}, pipes[2].Args, "chunk pipe args")
	})
}

// TestPipeParams_ParseErrors verifies that invalid pipe arg syntax produces the correct error codes.
func TestPipeParams_ParseErrors(t *testing.T) {
	tests := []struct {
		input   string
		errCode errors.ErrorCode
		desc    string
	}{
		// Variable reference not allowed in args
		{`arr |window($x): $window`, errors.ErrInvalidArgument, "variable reference not allowed"},
		// Expression not allowed
		{`arr |window(1+2): $window`, errors.ErrInvalidArgument, "arithmetic expression not allowed"},
		// Function call not allowed
		{`arr |window(len(x)): $window`, errors.ErrInvalidArgument, "function call not allowed"},
		// Missing closing paren
		{`arr |window(3: $window`, errors.ErrUnclosedFunction, "missing closing paren"},
		// Missing colon after closing paren
		{`arr |window(3) $window`, errors.ErrExpectedToken, "missing colon after args"},
		// Trailing comma
		{`arr |window(3,): $window`, errors.ErrInvalidArgument, "trailing comma"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			assertParseError(t, tt.input, tt.errCode)
		})
	}
}

// TestPipeParams_DefaultPipe verifies that |: syntax is unaffected by the args feature.
func TestPipeParams_DefaultPipe(t *testing.T) {
	// Compact form
	pipes, err := parsePipeArgs(t, "x |: $item + 1")
	assert.NoError(t, err)
	assert.Len(t, pipes, 2)
	assert.Nil(t, pipes[1].Args, "default pipe should always have nil args")

	// Spaced form: '| :' should behave identically
	pipes2, err2 := parsePipeArgs(t, "x | : $item + 1")
	assert.NoError(t, err2)
	assert.Len(t, pipes2, 2)
	assert.Nil(t, pipes2[1].Args, "default pipe (spaced) should always have nil args")
}
