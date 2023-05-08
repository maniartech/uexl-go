package types_test

import (
	"testing"

	"github.com/maniartech/uexl_go/types"
)

func TestMapGet(t *testing.T) {

	// Perform Map.Get() tests here
	m := types.Context{
		"a": types.Object{
			"b": types.Object{
				"c": types.Number(10),
				"d": types.Array{
					types.Number(1),
					types.Number(2),
					types.Object{
						"e": types.Number(3),
						"f": types.Array{
							types.Number(4),
							types.Number(5),
						},
					},
				},
			},
		},
	}

	if v, _ := m.ValueAtPath("a.b.c"); v != types.Number(10) {
		t.Errorf("Map.ValueAtPath() = %v", v)
	}

	// test for slice
	if v, _ := m.ValueAtPath("a.b.d.1"); v != types.Number(2) {
		t.Errorf("Map.ValueAtPath() = %v", v)
	}

	// test for nested slice
	if v, _ := m.ValueAtPath("a.b.d.2.f.1"); v != types.Number(5) {
		t.Errorf("Map.ValueAtPath() = %v", v)
	}
}
