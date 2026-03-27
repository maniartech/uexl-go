package vm

import (
	"testing"
)

// Compile-time check: *pipeContextImpl satisfies PipeContext.
var _ PipeContext = (*pipeContextImpl)(nil)

func TestPipeContextImpl_Args_nilWhenNotSet(t *testing.T) {
	pctx := &pipeContextImpl{}
	if got := pctx.Args(); got != nil {
		t.Errorf("expected nil args, got %v", got)
	}
}

func TestPipeContextImpl_Args_returnsSetSlice(t *testing.T) {
	args := []any{float64(3), "hello", true}
	pctx := &pipeContextImpl{args: args}
	got := pctx.Args()
	if len(got) != 3 {
		t.Fatalf("expected 3 args, got %d", len(got))
	}
	if got[0] != float64(3) {
		t.Errorf("args[0]: expected 3.0, got %v", got[0])
	}
	if got[1] != "hello" {
		t.Errorf("args[1]: expected \"hello\", got %v", got[1])
	}
	if got[2] != true {
		t.Errorf("args[2]: expected true, got %v", got[2])
	}
}

// WindowPipeHandler: explicit size 3 from args — input len=2 < windowSize=3 → no iterations,
// no EvalWith called, no nil-block error. This proves windowSize=3 was used instead of default 2.
func TestWindowPipeHandler_argsSize3_shortInput(t *testing.T) {
	machine := New(LibContext{})
	pctx := &pipeContextImpl{vm: machine, block: nil, args: []any{float64(3)}}
	res, err := WindowPipeHandler(pctx, []any{1.0, 2.0})
	if err != nil {
		t.Errorf("expected no error (no iterations for size=3 with len=2 input), got %v", err)
	}
	arr, _ := res.([]any)
	if len(arr) != 0 {
		t.Errorf("expected empty result for input shorter than window size, got %v", res)
	}
}

// WindowPipeHandler: arg below minimum (n=1) falls back to size 2,
// which triggers EvalWith on a 2-element input → nil block error proves fallback occurred.
func TestWindowPipeHandler_argsBelowMin_fallsBackTo2(t *testing.T) {
	machine := New(LibContext{})
	pctx := &pipeContextImpl{vm: machine, block: nil, args: []any{float64(1)}}
	_, err := WindowPipeHandler(pctx, []any{1.0, 2.0})
	if err == nil {
		t.Error("expected error from nil block when fallback size 2 triggers iteration")
	}
}

// WindowPipeHandler: non-float64 arg falls back to size 2 → same observable behaviour.
func TestWindowPipeHandler_argsInvalidType_fallsBackTo2(t *testing.T) {
	machine := New(LibContext{})
	pctx := &pipeContextImpl{vm: machine, block: nil, args: []any{"bad"}}
	_, err := WindowPipeHandler(pctx, []any{1.0, 2.0})
	if err == nil {
		t.Error("expected error from nil block when invalid arg type falls back to size 2")
	}
}

// WindowPipeHandler: empty args slice falls back to size 2.
func TestWindowPipeHandler_argsEmpty_fallsBackTo2(t *testing.T) {
	machine := New(LibContext{})
	pctx := &pipeContextImpl{vm: machine, block: nil, args: []any{}}
	_, err := WindowPipeHandler(pctx, []any{1.0, 2.0})
	if err == nil {
		t.Error("expected error from nil block when empty args falls back to size 2")
	}
}

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
