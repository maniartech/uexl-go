package uexl

import (
	"fmt"
	"sort"
	"strings"
)

// EnvInfo is a read-only snapshot of an Env's registered symbols.
// It is safe to copy and pass around freely — independent of the source Env.
type EnvInfo struct {
	Functions    []string // sorted function names
	PipeHandlers []string // sorted pipe handler names
	Globals      []string // sorted global variable names
}

// String implements fmt.Stringer with a stable human-readable multiline format:
//
//	Env:
//	  Functions (N): name1, name2
//	  PipeHandlers (N): name1, name2
//	  Globals (N): name1, name2
func (i EnvInfo) String() string {
	var b strings.Builder
	b.WriteString("Env:\n")
	fmt.Fprintf(&b, "  Functions (%d): %s\n", len(i.Functions), strings.Join(i.Functions, ", "))
	fmt.Fprintf(&b, "  PipeHandlers (%d): %s\n", len(i.PipeHandlers), strings.Join(i.PipeHandlers, ", "))
	fmt.Fprintf(&b, "  Globals (%d): %s\n", len(i.Globals), strings.Join(i.Globals, ", "))
	return b.String()
}

// sortedKeys returns a lexicographically sorted slice of map keys.
// Returns an empty slice (never nil) when m is empty.
func sortedKeys[V any](m map[string]V) []string {
	if len(m) == 0 {
		return []string{}
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
