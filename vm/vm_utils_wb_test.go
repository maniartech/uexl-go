package vm

import (
	"testing"
)

func TestNewVM(t *testing.T) {
	vm := New(LibContext{})
	if vm == nil {
		t.Fatal("expected non-nil VM")
	}
}
