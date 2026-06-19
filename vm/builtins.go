package vm

import (
	"github.com/maniartech/uexl/builtins/core"
)

// Builtins is the runtime's default function set: the always-on core built-ins. Their implementations
// live in the builtins/core package (alongside the attachable standard-library families in builtins/*);
// here they are adapted to the vm's VMFunction type. The import direction is vm -> builtins/core; nothing
// under builtins/ imports vm.
var Builtins = coreBuiltins()

func coreBuiltins() VMFunctions {
	m := make(VMFunctions, len(core.Builtins))
	for name, f := range core.Builtins {
		m[name] = VMFunction(f)
	}
	return m
}
