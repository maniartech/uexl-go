package ast

import (
	"encoding/json"
)

// Print prints the node in a json format.
func PrintNode(node Node) {
	b, _ := json.MarshalIndent(node, "", "  ")
	println(string(b))
}
