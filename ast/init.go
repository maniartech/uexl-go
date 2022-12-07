package ast

import "encoding/json"

type Map map[string]any

type Node interface {
	Eval(Map) (any, error)
}

type BaseNode struct {

	// Line is the line number of the node in the source code.
	Line int `json:"line"`

	// Column is the column number of the node in the source code.
	Column int `json:"column"`

	// Offset is the position of the node in the source code.
	Offset int `json:"offset"`

	// Type is the type of the node.
	Type NodeType `json:"type"`
}

// Print prints the node in a json format.
func (n BaseNode) Print() {
	b, _ := json.MarshalIndent(n, "", "  ")
	println(string(b))
}
