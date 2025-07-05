package parser_test

import (
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/ast"
	. "github.com/maniartech/uexl_go/parser"
)

func TestStrings(t *testing.T) {
	testCases := []struct {
		input          string
		expectedValue  string
		isSingleQuoted bool
		isRaw          bool
	}{
		// Single quoted string
		{`'This is a single quoted string where " is ignored'`, "This is a single quoted string where \" is ignored", true, false},
		// Double quoted string
		{`"This is a double quoted string where ' is ignored"`, "This is a double quoted string where ' is ignored", false, false},
		// Raw string (single quoted)
		{`r'This sentence is a raw string #$@$SD'`, "This sentence is a raw string #$@$SD", true, true},
		// Raw string (double quoted)
		{`r"This is a raw string with \\u00A1 and ""quotes"""`, "This is a raw string with \\\\u00A1 and \"quotes\"", false, true},
		// Escaped unicode in double quoted string
		{`"This contains unicode: \u00A1"`, "This contains unicode: ¡", false, false},
		// Escaped unicode in single quoted string
		{`'This contains unicode: \u00A1'`, "This contains unicode: ¡", true, false},
		// Escaped quotes in single quoted string
		{`'This has an escaped \"quote\"'`, "This has an escaped \"quote\"", true, false},
		// Escaped quotes in double quoted string
		{`"This has an escaped \"quote\""`, "This has an escaped \"quote\"", false, false},
		// Escaped backslash in single quoted string
		{`'This has a backslash: \\'`, "This has a backslash: \\", true, false},
		// Escaped backslash in double quoted string
		{`"This has a backslash: \\ \""`, "This has a backslash: \\ \"", false, false},
	}

	for _, tc := range testCases {
		testString(tc.input, tc.expectedValue, t, tc.isSingleQuoted)
	}
}

func testString(str, testValue string, t *testing.T, expectSingleQuoted ...bool) {
	parsed, err := ParseReaderNew("", strings.NewReader(str))
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	node := parsed.(*ast.StringNode)
	value := node.Value

	// Only check the value, not the token, as the parser normalizes quotes and raw strings
	if string(value) != testValue {
		t.Errorf("Value: Expected %v, got %v", testValue, value)
	}

	// Check IsSingleQuoted if expectation is provided
	if len(expectSingleQuoted) > 0 {
		if node.IsSingleQuoted != expectSingleQuoted[0] {
			t.Errorf("IsSingleQuoted: Expected %v, got %v", expectSingleQuoted[0], node.IsSingleQuoted)
		}
	}
}
