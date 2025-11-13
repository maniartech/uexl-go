package compiler

import (
	"testing"

	"github.com/maniartech/uexl/types"
)

func TestNewCompiler(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("expected non-nil compiler")
	}
	if len(c.scopes) == 0 {
		t.Error("expected at least one scope")
	}
}

func TestNewWithState(t *testing.T) {
	c := NewWithState([]types.Value{types.NewStringValue("foo")})
	if len(c.constants) != 1 {
		t.Error("constants not set correctly")
	}
	if str, ok := c.constants[0].AsString(); !ok || str != "foo" {
		t.Errorf("expected 'foo', got %v", c.constants[0].ToAny())
	}
}

func TestAddContextVar(t *testing.T) {
	c := New()
	idx1 := c.addContextVar("foo")
	if idx1 != 0 {
		t.Errorf("expected index 0, got %d", idx1)
	}
	idx2 := c.addContextVar("bar")
	if idx2 != 1 {
		t.Errorf("expected index 1, got %d", idx2)
	}
	idx3 := c.addContextVar("foo")
	if idx3 != 0 {
		t.Errorf("expected index 0 for duplicate, got %d", idx3)
	}
}

func TestEnterExitScope(t *testing.T) {
	c := New()
	c.enterScope()
	if len(c.scopes) != 2 || c.scopeIndex != 1 {
		t.Error("enterScope did not add new scope")
	}
	c.exitScope()
	if len(c.scopes) != 1 || c.scopeIndex != 0 {
		t.Error("exitScope did not remove scope")
	}
}

func TestExitScopePanicsAtGlobal(t *testing.T) {
	c := New()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on exitScope at global scope")
		}
	}()
	c.exitScope()
}

func TestIsPipeLocalVar(t *testing.T) {
	if !isPipeLocalVar("$foo") {
		t.Error("expected $foo to be pipe local var")
	}
	if isPipeLocalVar("foo") {
		t.Error("expected foo to not be pipe local var")
	}
}
