package vm

import (
	"testing"
)

func TestGetContextValue_NilContext(t *testing.T) {
	vm := &VM{}
	_, err := vm.getContextValue("foo")
	if err == nil {
		t.Error("expected error for nil contextVarsValues")
	}
}
