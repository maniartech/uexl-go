package compiler

import (
	"testing"
)

func TestByteCode_StructFields(t *testing.T) {
	b := &ByteCode{}
	if b.Instructions != nil {
		t.Errorf("expected Instructions to be nil")
	}
	if b.Constants != nil {
		t.Errorf("expected Constants to be nil")
	}
	if b.ContextVars != nil {
		t.Errorf("expected ContextVars to be nil")
	}
	if b.SystemVars != nil {
		t.Errorf("expected SystemVars to be nil")
	}
}
