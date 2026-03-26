package uexl_test

import (
	"testing"

	"github.com/maniartech/uexl"
	"github.com/stretchr/testify/assert"
)

// ── AsFloat64 ────────────────────────────────────────────────────────────────

func TestAsFloat64_float64(t *testing.T) {
	v, err := uexl.AsFloat64(3.14)
	assert.NoError(t, err)
	assert.Equal(t, 3.14, v)
}

func TestAsFloat64_float64_zero(t *testing.T) {
	v, err := uexl.AsFloat64(0.0)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, v)
}

func TestAsFloat64_int(t *testing.T) {
	v, err := uexl.AsFloat64(int(42))
	assert.NoError(t, err)
	assert.Equal(t, 42.0, v)
}

func TestAsFloat64_int_negative(t *testing.T) {
	v, err := uexl.AsFloat64(int(-7))
	assert.NoError(t, err)
	assert.Equal(t, -7.0, v)
}

func TestAsFloat64_int64(t *testing.T) {
	v, err := uexl.AsFloat64(int64(100))
	assert.NoError(t, err)
	assert.Equal(t, 100.0, v)
}

func TestAsFloat64_float32(t *testing.T) {
	v, err := uexl.AsFloat64(float32(1.5))
	assert.NoError(t, err)
	assert.InDelta(t, 1.5, v, 1e-6)
}

func TestAsFloat64_string_errors(t *testing.T) {
	_, err := uexl.AsFloat64("3.14")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "AsFloat64")
	assert.Contains(t, err.Error(), "float64")
}

func TestAsFloat64_bool_errors(t *testing.T) {
	_, err := uexl.AsFloat64(true)
	assert.Error(t, err)
}

func TestAsFloat64_nil_errors(t *testing.T) {
	_, err := uexl.AsFloat64(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "AsFloat64")
}

func TestAsFloat64_slice_errors(t *testing.T) {
	_, err := uexl.AsFloat64([]any{1.0})
	assert.Error(t, err)
}

func TestAsFloat64_map_errors(t *testing.T) {
	_, err := uexl.AsFloat64(map[string]any{"k": 1.0})
	assert.Error(t, err)
}

// Round-trip: evaluation result → AsFloat64.
func TestAsFloat64_evalRoundtrip(t *testing.T) {
	result, err := uexl.Eval("price * qty", map[string]any{"price": 4.5, "qty": 2.0})
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	v, err := uexl.AsFloat64(result)
	assert.NoError(t, err)
	assert.Equal(t, 9.0, v)
}

// ── AsBool ────────────────────────────────────────────────────────────────────

func TestAsBool_true(t *testing.T) {
	v, err := uexl.AsBool(true)
	assert.NoError(t, err)
	assert.True(t, v)
}

func TestAsBool_false(t *testing.T) {
	v, err := uexl.AsBool(false)
	assert.NoError(t, err)
	assert.False(t, v)
}

func TestAsBool_int_errors_noTruthyCoercion(t *testing.T) {
	_, err := uexl.AsBool(1)
	assert.Error(t, err, "AsBool must not coerce int to bool")
	assert.Contains(t, err.Error(), "AsBool")
}

func TestAsBool_zero_errors(t *testing.T) {
	_, err := uexl.AsBool(0)
	assert.Error(t, err, "AsBool must not coerce 0 to false")
}

func TestAsBool_float64_errors(t *testing.T) {
	_, err := uexl.AsBool(1.0)
	assert.Error(t, err)
}

func TestAsBool_string_errors(t *testing.T) {
	_, err := uexl.AsBool("true")
	assert.Error(t, err)
}

func TestAsBool_nil_errors(t *testing.T) {
	_, err := uexl.AsBool(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "AsBool")
}

func TestAsBool_evalRoundtrip(t *testing.T) {
	result, err := uexl.Eval("x > 5", map[string]any{"x": 10.0})
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	v, err := uexl.AsBool(result)
	assert.NoError(t, err)
	assert.True(t, v)
}

func TestAsBool_evalFalse(t *testing.T) {
	result, err := uexl.Eval("x > 5", map[string]any{"x": 3.0})
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	v, err := uexl.AsBool(result)
	assert.NoError(t, err)
	assert.False(t, v)
}

// ── AsString ──────────────────────────────────────────────────────────────────

func TestAsString_string(t *testing.T) {
	v, err := uexl.AsString("hello")
	assert.NoError(t, err)
	assert.Equal(t, "hello", v)
}

func TestAsString_empty(t *testing.T) {
	v, err := uexl.AsString("")
	assert.NoError(t, err)
	assert.Equal(t, "", v)
}

func TestAsString_int_errors_noConversion(t *testing.T) {
	_, err := uexl.AsString(42)
	assert.Error(t, err, "AsString must not call fmt.Sprint on non-strings")
	assert.Contains(t, err.Error(), "AsString")
}

func TestAsString_float64_errors(t *testing.T) {
	_, err := uexl.AsString(3.14)
	assert.Error(t, err)
}

func TestAsString_bool_errors(t *testing.T) {
	_, err := uexl.AsString(true)
	assert.Error(t, err)
}

func TestAsString_nil_errors(t *testing.T) {
	_, err := uexl.AsString(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "AsString")
}

func TestAsString_evalRoundtrip(t *testing.T) {
	result, err := uexl.Eval(`"hello" + " world"`, nil)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	v, err := uexl.AsString(result)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", v)
}

// ── AsSlice ───────────────────────────────────────────────────────────────────

func TestAsSlice_slice(t *testing.T) {
	input := []any{1.0, 2.0, 3.0}
	v, err := uexl.AsSlice(input)
	assert.NoError(t, err)
	assert.Equal(t, input, v)
}

func TestAsSlice_empty(t *testing.T) {
	v, err := uexl.AsSlice([]any{})
	assert.NoError(t, err)
	assert.Empty(t, v)
}

func TestAsSlice_nil_errors(t *testing.T) {
	_, err := uexl.AsSlice(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "AsSlice")
}

func TestAsSlice_string_errors(t *testing.T) {
	_, err := uexl.AsSlice("oops")
	assert.Error(t, err)
}

func TestAsSlice_float64_errors(t *testing.T) {
	_, err := uexl.AsSlice(3.14)
	assert.Error(t, err)
}

func TestAsSlice_evalRoundtrip(t *testing.T) {
	result, err := uexl.Eval("[1, 2, 3] |map: $item * 2", nil)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	v, err := uexl.AsSlice(result)
	assert.NoError(t, err)
	assert.Equal(t, []any{2.0, 4.0, 6.0}, v)
}

// ── AsMap ────────────────────────────────────────────────────────────────────

func TestAsMap_map(t *testing.T) {
	input := map[string]any{"key": "value"}
	v, err := uexl.AsMap(input)
	assert.NoError(t, err)
	assert.Equal(t, input, v)
}

func TestAsMap_empty(t *testing.T) {
	v, err := uexl.AsMap(map[string]any{})
	assert.NoError(t, err)
	assert.Empty(t, v)
}

func TestAsMap_nil_errors(t *testing.T) {
	_, err := uexl.AsMap(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "AsMap")
}

func TestAsMap_slice_errors(t *testing.T) {
	_, err := uexl.AsMap([]any{1.0})
	assert.Error(t, err)
}

func TestAsMap_string_errors(t *testing.T) {
	_, err := uexl.AsMap("oops")
	assert.Error(t, err)
}

func TestAsMap_evalRoundtrip(t *testing.T) {
	result, err := uexl.Eval(`{"key": 'value', "num": 42}`, nil)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	v, err := uexl.AsMap(result)
	assert.NoError(t, err)
	assert.Equal(t, "value", v["key"])
	assert.Equal(t, 42.0, v["num"])
}
