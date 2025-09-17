package vm

import (
	"testing"
)

func TestDefaultPipeHandler_NilBlock(t *testing.T) {
	vm := &VM{}
	res, err := DefaultPipeHandler(1, nil, "", vm)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if res != 1 {
		t.Errorf("expected input to be returned")
	}
}
