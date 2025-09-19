package vm

import (
	"testing"
)

func TestBuiltinLen_InvalidType(t *testing.T) {
	_, err := builtinLen(struct{}{})
	if err == nil {
		t.Error("expected error for unsupported type")
	}
}

func TestBuiltinLen_validInputs(t *testing.T) {
	// Test string
	result, err := builtinLen("hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != float64(5) {
		t.Errorf("expected 5, got %v", result)
	}

	// Test array
	result, err = builtinLen([]any{1, 2, 3})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != float64(3) {
		t.Errorf("expected 3, got %v", result)
	}

	// Test empty string
	result, err = builtinLen("")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != float64(0) {
		t.Errorf("expected 0, got %v", result)
	}

	// Test empty array
	result, err = builtinLen([]any{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != float64(0) {
		t.Errorf("expected 0, got %v", result)
	}
}

func TestBuiltinLen_wrongArgCount(t *testing.T) {
	_, err := builtinLen()
	if err == nil {
		t.Error("expected error for no arguments")
	}

	_, err = builtinLen("a", "b")
	if err == nil {
		t.Error("expected error for too many arguments")
	}
}

func TestBuiltinSubstr_InvalidArgs(t *testing.T) {
	_, err := builtinSubstr("hello", 1)
	if err == nil {
		t.Error("expected error for wrong arg count")
	}
}

func TestBuiltinSubstr_validCases(t *testing.T) {
	// Valid substring - start=1, length=3
	result, err := builtinSubstr("hello", 1.0, 3.0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != "ell" { // substr("hello", 1, 3) = "ell" (start at index 1, length 3)
		t.Errorf("expected 'ell', got %v", result)
	}

	// Start at 0 - start=0, length=2
	result, err = builtinSubstr("hello", 0.0, 2.0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != "he" { // substr("hello", 0, 2) = "he" (start at index 0, length 2)
		t.Errorf("expected 'he', got %v", result)
	}
}

func TestBuiltinSubstr_nonStringInput(t *testing.T) {
	_, err := builtinSubstr(123, 0.0, 1.0)
	if err == nil {
		t.Error("expected error for non-string input")
	}
}

func TestBuiltinSubstr_nonNumericIndices(t *testing.T) {
	_, err := builtinSubstr("hello", "start", 3.0)
	if err == nil {
		t.Error("expected error for non-numeric start index")
	}

	_, err = builtinSubstr("hello", 1.0, "end")
	if err == nil {
		t.Error("expected error for non-numeric end index")
	}
}

func TestBuiltinStr_allTypes(t *testing.T) {
	tests := []struct {
		input    any
		expected string
	}{
		{42.0, "42"},
		{"hello", "hello"},
		{true, "true"},
		{false, "false"},
		{nil, "<nil>"}, // Go's fmt.Sprintf prints nil as "<nil>"
	}

	for _, tt := range tests {
		result, err := builtinStr(tt.input)
		if err != nil {
			t.Errorf("builtinStr(%v) error: %v", tt.input, err)
		}
		if result != tt.expected {
			t.Errorf("builtinStr(%v) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestBuiltinStr_wrongArgCount(t *testing.T) {
	_, err := builtinStr()
	if err == nil {
		t.Error("expected error for no arguments")
	}

	_, err = builtinStr(1, 2)
	if err == nil {
		t.Error("expected error for too many arguments")
	}
}
