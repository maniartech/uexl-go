package types_test

import (
	"testing"

	"github.com/maniartech/uexl_go/types"
)

func TestMapGet(t *testing.T) {

	// Perform Map.Get() tests here
	m := types.Map{
		"a": types.Map{
			"b": types.Map{
				"c": 10,
				"d": []interface{}{
					1,
					2,
					types.Map{
						"e": 3,
						"f": []interface{}{
							4,
							5,
						},
					},
				},
			},
		},
	}

	if v, _ := m.Get("a.b.c"); v != 10 {
		t.Errorf("Map.Get() = %v", v)
	}

	// test for slice
	if v, _ := m.Get("a.b.d.1"); v != 2 {
		t.Errorf("Map.Get() = %v", v)
	}

	// test for nested slice
	if v, _ := m.Get("a.b.d.2.f.1"); v != 5 {
		t.Errorf("Map.Get() = %v", v)
	}
}
