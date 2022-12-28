package ast

import "encoding/json"

type Map map[string]any

// Print prints the node in a json format.
func PrintNode(node Node) {
	b, _ := json.MarshalIndent(node, "", "  ")
	println(string(b))
}
