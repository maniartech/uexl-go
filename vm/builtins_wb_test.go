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

func TestBuiltinSubstr_InvalidArgs(t *testing.T) {
	_, err := builtinSubstr("hello", 1)
	if err == nil {
		t.Error("expected error for wrong arg count")
	}
}
