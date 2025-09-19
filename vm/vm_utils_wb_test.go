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

func TestVM_Top_empty(t *testing.T) {
	vm := New(LibContext{})
	result := vm.Top()
	if result != nil {
		t.Errorf("expected nil for empty stack, got %v", result)
	}
}

func TestVM_Top_withItems(t *testing.T) {
	vm := New(LibContext{})
	vm.Push(42)
	vm.Push("hello")
	result := vm.Top()
	if result != "hello" {
		t.Errorf("expected 'hello', got %v", result)
	}
}

func TestVM_Push_stackOverflow(t *testing.T) {
	vm := &VM{
		stack: make([]any, 1024), // StackSize
		sp:    1024,              // Already at capacity
	}

	// This should return an error
	err := vm.Push(1)
	if err == nil {
		t.Error("expected stack overflow error")
	}
}

func TestVM_Pop(t *testing.T) {
	vm := New(LibContext{})
	vm.Push(1)
	vm.Push(2)
	vm.Push(3)

	popped := vm.Pop()
	if popped != 3 {
		t.Errorf("expected 3, got %v", popped)
	}
	if vm.sp != 2 {
		t.Errorf("expected sp=2, got %d", vm.sp)
	}
}

func TestVM_LastPoppedStackElem_withItems(t *testing.T) {
	vm := New(LibContext{})
	vm.Push(1)
	vm.Push(2)
	vm.Push(3)
	popped := vm.Pop() // Pop 3, sp now points to where 2 is at top

	// Should return element at sp-1 (the element at the current top, which is 2)
	elem := vm.LastPoppedStackElem()
	if elem != 2 { // The current top element after pop
		t.Errorf("expected 2, got %v", elem)
	}
	if popped != 3 {
		t.Errorf("expected popped=3, got %v", popped)
	}
}

func TestVM_LastPoppedStackElem_emptyStack(t *testing.T) {
	vm := New(LibContext{})

	// Should return nil for empty stack
	elem := vm.LastPoppedStackElem()
	if elem != nil {
		t.Errorf("expected nil, got %v", elem)
	}
}

func TestVM_getPipeVar_found(t *testing.T) {
	vm := New(LibContext{})
	vm.pushPipeScope()
	vm.setPipeVar("var1", "value1")

	value, found := vm.getPipeVar("var1")
	if !found {
		t.Error("expected to find var1")
	}
	if value != "value1" {
		t.Errorf("expected 'value1', got %v", value)
	}
}

func TestVM_getPipeVar_notFound(t *testing.T) {
	vm := New(LibContext{})
	vm.pushPipeScope()
	vm.setPipeVar("var1", "value1")

	_, found := vm.getPipeVar("nonexistent")
	if found {
		t.Error("expected not to find nonexistent variable")
	}
}

func TestVM_getPipeVar_emptyScopes(t *testing.T) {
	vm := New(LibContext{})

	_, found := vm.getPipeVar("var1")
	if found {
		t.Error("expected not to find variable in empty scopes")
	}
}
