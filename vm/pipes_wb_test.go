package vm

import (
	"testing"
)

func TestDefaultPipeHandler_NilBlock(t *testing.T) {
	machine := New(LibContext{})
	pctx := &pipeContextImpl{vm: machine, block: nil}
	res, err := DefaultPipeHandler(pctx, 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if res != 1 {
		t.Errorf("expected input to be returned, got %v", res)
	}
}

func TestGroupByPipeHandler_nonArray(t *testing.T) {
	machine := New(LibContext{})
	pctx := &pipeContextImpl{vm: machine, block: nil}
	_, err := GroupByPipeHandler(pctx, "not an array")
	if err == nil {
		t.Error("expected error for non-array input")
	}
}

func TestGroupByPipeHandler_nilBlock(t *testing.T) {
	machine := New(LibContext{})
	pctx := &pipeContextImpl{vm: machine, block: nil}
	_, err := GroupByPipeHandler(pctx, []any{1, 2, 3})
	if err == nil {
		t.Error("expected error for nil block")
	}
}

func TestChunkPipeHandler_nonArray(t *testing.T) {
	machine := New(LibContext{})
	pctx := &pipeContextImpl{vm: machine, block: nil}
	_, err := ChunkPipeHandler(pctx, "not an array")
	if err == nil {
		t.Error("expected error for non-array input")
	}
}

func TestChunkPipeHandler_invalidSize(t *testing.T) {
	machine := New(LibContext{})
	pctx := &pipeContextImpl{vm: machine, block: nil}
	_, err := ChunkPipeHandler(pctx, []any{1, 2, 3})
	if err == nil {
		t.Error("expected error for nil block")
	}
}
