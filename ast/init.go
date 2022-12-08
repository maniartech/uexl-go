package ast

import "encoding/json"

type Map map[string]any

type Node interface {
	Eval(Map) (any, error)

	GetBaseNode() BaseNode
}

type BaseNode struct {

	// Line is the line number of the node in the source code.
	Line int `json:"line"`

	// Column is the column number of the node in the source code.
	Column int `json:"column"`

	// Offset is the position of the node in the source code.
	Offset int `json:"offset"`

	Token string `json:"token"`

	// Type is the type of the node.
	Type NodeType `json:"type"`
}

func (n BaseNode) GetBaseNode() BaseNode {
	return n
}

// Print prints the node in a json format.
func PrintNode(node Node) {
	b, _ := json.MarshalIndent(node, "", "  ")
	println(string(b))
}
