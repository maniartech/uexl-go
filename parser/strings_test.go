package parser

import (
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
		// sdq := "\"" + s + "\""

		// Test case for double quoted string
		ssq := "'" + s + "'"
		parsed, err := ParseReader("", strings.NewReader(ssq))
		if err != nil {
			t.Errorf("Error: %v", err)
			continue
		}

		node := parsed.(*ast.StringNode)
		token := node.Token
		value := node.Value

		// Check if the token is same as the string (s)
		if token != ssq {
			t.Errorf("Token: Expected %v, got %v", ssq, token)
		}

		// Check if the value is escaped string
		escaped := escappedValues[i]

		// Check if the value is same as the escaped string
		if string(value) != escaped {
			t.Errorf("Value: Expected %v, got %v", escaped, value)
		}
	}
}
