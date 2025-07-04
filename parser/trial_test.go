package parser

import (
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/ast"
)

func TestPlayground(t *testing.T) {
	// This test is to be used as a playground for experimenting and debugging
	// failing tests or new features.

	// Test cases for different string types
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"Single quoted", `'hello world'`, "hello world"},
		{"Double quoted", `"hello world"`, "hello world"},
		{"Raw single quoted", `r'hello\nworld'`, "hello\\nworld"},
		{"Raw double quoted", `r"hello\nworld"`, "hello\\nworld"},
		{"Raw with embedded quotes", `r"text with ""quotes"""`, "text with \"quotes\""},
		{"Escaped unicode", `"unicode: \u00A1"`, "unicode: ¡"},
		{"Raw with backslashes", `r"path\to\file"`, "path\\to\\file"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parsed, err := ParseReaderNew("", strings.NewReader(tc.input))
			if err != nil {
				t.Errorf("Parser error: %v", err)
				return
			}

			if node, ok := parsed.(*ast.StringNode); ok {
				if string(node.Value) != tc.expected {
					t.Errorf("Expected: %s, got: %s", tc.expected, string(node.Value))
				}
			} else {
				t.Errorf("Expected StringNode, got %T", parsed)
			}
		})
	}
}
