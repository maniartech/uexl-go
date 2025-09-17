package vm

import (
	"testing"
)

func TestExecuteSliceExpression_NilTarget(t *testing.T) {
	vm := &VM{}
	err := vm.executeSliceExpression(nil, 0, 0, 1, false)
	if err == nil {
		t.Error("expected error for nil target")
	}
}
