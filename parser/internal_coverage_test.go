package parser

import (
	"testing"
)

// TestInternalCoverage provides coverage measurement for the parser package
// This test exercises the main public APIs to measure coverage
func TestInternalCoverage(t *testing.T) {
	// Test basic parser functionality
	p := NewParser("1 + 2 * 3")
	if p == nil {
		t.Fatal("NewParser returned nil")
	}

	expr, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if expr == nil {
		t.Fatal("Parse returned nil expression")
	}

	// Test tokenizer
	tokenizer := NewTokenizer("hello world 123")
	if tokenizer == nil {
		t.Fatal("NewTokenizer returned nil")
	}

	token, err := tokenizer.NextToken()
	if err != nil {
		t.Fatalf("NextToken failed: %v", err)
	}
	if token.Type == 0 {
		t.Fatal("NextToken returned invalid token")
	}

	// Test options
	opts := DefaultOptions()
	if !opts.EnableNullish {
		t.Fatal("DefaultOptions should enable nullish")
	}

	// Test ParseString
	node, err := ParseString("42")
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}
	if node == nil {
		t.Fatal("ParseString returned nil")
	}

	// Test property helpers
	strProp := PropS("test")
	if !strProp.IsString() {
		t.Fatal("PropS should create string property")
	}

	intProp := PropI(42)
	if !intProp.IsInt() {
		t.Fatal("PropI should create int property")
	}
}
