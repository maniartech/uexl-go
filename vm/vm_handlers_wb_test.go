package vm

import (
	"testing"

	"github.com/maniartech/uexl/code"
)

func TestGetContextValue_NilContext(t *testing.T) {
	vm := &VM{}
	_, err := vm.getContextValue("foo")
	if err == nil {
		t.Error("expected error for nil context")
	}
}

func TestExecuteBooleanBinaryOperation_logicalAnd(t *testing.T) {
	vm := &VM{stack: make([]Value, 10)}
	err := vm.executeBooleanBinaryOperation(code.OpLogicalAnd, true, false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	result := vm.Top()
	if result != false {
		t.Errorf("expected false, got %v", result)
	}
}

func TestExecuteBooleanBinaryOperation_logicalOr(t *testing.T) {
	vm := &VM{stack: make([]Value, 10)}
	err := vm.executeBooleanBinaryOperation(code.OpLogicalOr, false, true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	result := vm.Top()
	if result != true {
		t.Errorf("expected true, got %v", result)
	}
}

func TestExecuteBooleanBinaryOperation_unsupported(t *testing.T) {
	vm := &VM{stack: make([]Value, 10)}
	err := vm.executeBooleanBinaryOperation(code.OpEqual, true, false) // OpEqual (unsupported for booleans)
	if err == nil {
		t.Error("expected error for unsupported boolean operation")
	}
}

func TestExecuteIndex_array(t *testing.T) {
	vm := &VM{stack: make([]Value, 10)}
	err := vm.executeIndex([]any{1, 2, 3}, 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestExecuteIndex_map(t *testing.T) {
	vm := &VM{stack: make([]Value, 10)}
	err := vm.executeIndex(map[string]any{"key": "value"}, "key")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestExecuteIndex_nil(t *testing.T) {
	vm := &VM{stack: make([]Value, 10)}
	err := vm.executeIndex(nil, 0)
	if err == nil {
		t.Error("expected error for nil target")
	}
}

func TestExecuteIndexValue_array(t *testing.T) {
	vm := &VM{stack: make([]Value, 10)}
	err := vm.executeIndexValue([]any{1, 2, 3}, 1.0)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestExecuteIndexValue_string(t *testing.T) {
	vm := &VM{stack: make([]Value, 10)}
	err := vm.executeIndexValue("hello", 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestExecuteIndexValue_outOfBounds(t *testing.T) {
	vm := &VM{stack: make([]Value, 10)}
	err := vm.executeIndexValue([]any{1, 2}, 5)
	if err == nil {
		t.Error("expected error for out of bounds index")
	}
}
