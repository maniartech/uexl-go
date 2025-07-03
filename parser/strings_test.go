package parser

import (
	"strconv"
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/ast"
)

func TestStrings(t *testing.T) {

	// testStrs is list of random strings to test the parser.
	// It contains all the possible combinations of single and double quotes.
	// It also contains strings with escaped characters, strings with
	// escaped quotes and unicode characters. It contains around 50
	// strings to test the parser.
	testStrs := []string{
		`This \u00A1sentence has a Unicode character!`,
		`This sentence has \"escaped quotes\"`,
		`This sentence has an \\escaped character`,
		`\u00A1This sentence has a Unicode character and \"escaped quotes\"`,
		`This sentence has an \\escaped character and \"escaped quotes\"`,
		`This sentence has an \\escaped character and \u00A1a Unicode character!`,
		`\u00A1This sentence has a Unicode character and an \\escaped character`,
		`This sentence has \"escaped quotes\" and \u00A1a Unicode character!`,
		`\u00A1This sentence has a Unicode character and \"escaped quotes\" and an \\escaped character`,
		`This sentence has an \\escaped character and \"escaped quotes\" and \u00A1a Unicode character!`,
		`This sentence has an \\escaped character and \u00A1a Unicode character! and \"escaped quotes\"`,
		`\u00A1This sentence has a Unicode character and an \\escaped character and \"escaped quotes\"`,
		`This sentence has \"escaped quotes\" and an \\escaped character and \u00A1a Unicode character!`,
		`\u00A1This sentence has a Unicode character and \"escaped quotes\" and an \\escaped character and \u00A1a Unicode character!`,
		`This sentence has an \\escaped character and \"escaped quotes\" and \u00A1a Unicode character! and an \\escaped character`,
		`This sentence has an encoded emoji \ud83d\ude00, and an encoded character \u00a1,  escaped character \\, escaped quotes \" and an escaped quote character \"`,
	}

	escappedValues := []string{
		`This ¡sentence has a Unicode character!`,
		`This sentence has "escaped quotes"`,
		`This sentence has an \escaped character`,
		`¡This sentence has a Unicode character and "escaped quotes"`,
		`This sentence has an \escaped character and "escaped quotes"`,
		`This sentence has an \escaped character and ¡a Unicode character!`,
		`¡This sentence has a Unicode character and an \escaped character`,
		`This sentence has "escaped quotes" and ¡a Unicode character!`,
		`¡This sentence has a Unicode character and "escaped quotes" and an \escaped character`,
		`This sentence has an \escaped character and "escaped quotes" and ¡a Unicode character!`,
		`This sentence has an \escaped character and ¡a Unicode character! and "escaped quotes"`,
		`¡This sentence has a Unicode character and an \escaped character and "escaped quotes"`,
		`This sentence has "escaped quotes" and an \escaped character and ¡a Unicode character!`,
		`¡This sentence has a Unicode character and "escaped quotes" and an \escaped character and ¡a Unicode character!`,
		`This sentence has an \escaped character and "escaped quotes" and ¡a Unicode character! and an \escaped character`,
		`This sentence has an encoded emoji 😀, and an encoded character ¡,  escaped character \, escaped quotes " and an escaped quote character "`,
	}

	for i, s := range testStrs {
		// Test case for double quoted string (Go/JSON compliant)
		str := strconv.Quote(s)
		testString(str, escappedValues[i], t)

		// Test case for single quoted string (convert escapes for Go compatibility)
		single := s
		// Replace \" with ", \' with ', leave \\ as is
		single = strings.ReplaceAll(single, "\\\"", "\"")
		single = strings.ReplaceAll(single, "\\'", "'")
		str = "'" + single + "'"
		testString(str, escappedValues[i], t)

		// Test case for raw string (r"..." and r'...')
		rawDouble := "r\"" + s + "\""
		testString(rawDouble, s, t)
		rawSingle := "r'" + s + "'"
		testString(rawSingle, s, t)
	}
}

func testString(str, testValue string, t *testing.T) {
	parsed, err := ParseReaderNew("", strings.NewReader(str))
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	node := parsed.(*ast.StringNode)
	value := node.Value

	// Debug output
	t.Logf("Input: %s", str)
	t.Logf("Expected: %s", testValue)
	t.Logf("Actual: %s", value)
	t.Logf("Node Type: %T", parsed)

	// Only check the value, not the token, as the parser normalizes quotes and raw strings
	if string(value) != testValue {
		t.Errorf("Value: Expected %v, got %v", testValue, value)
	}
}
