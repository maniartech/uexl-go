package vm

import (
	"testing"
)

func TestExecuteIndexExpression_NilLeft(t *testing.T) {
	vm := &VM{}
	err := vm.executeIndexExpression(nil, 0, false)
	if err == nil {
		t.Error("expected error for nil left")
	}
}
