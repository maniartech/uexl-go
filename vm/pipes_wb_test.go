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

func TestGroupByPipeHandler_nonArray(t *testing.T) {
	vm := &VM{}
	_, err := GroupByPipeHandler("not an array", nil, "", vm)
	if err == nil {
		t.Error("expected error for non-array input")
	}
}

func TestGroupByPipeHandler_nilBlock(t *testing.T) {
	vm := &VM{}
	_, err := GroupByPipeHandler([]any{1, 2, 3}, nil, "", vm)
	if err == nil {
		t.Error("expected error for nil block")
	}
}

func TestChunkPipeHandler_nonArray(t *testing.T) {
	vm := &VM{}
	_, err := ChunkPipeHandler("not an array", nil, "", vm)
	if err == nil {
		t.Error("expected error for non-array input")
	}
}

func TestChunkPipeHandler_invalidSize(t *testing.T) {
	vm := &VM{}
	_, err := ChunkPipeHandler([]any{1, 2, 3}, "invalid", "", vm)
	if err == nil {
		t.Error("expected error for invalid chunk size")
	}
}
