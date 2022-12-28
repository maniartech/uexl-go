package ast

import "encoding/json"

// BaseNode provides the base node for all nodes.
// It contains the line number, column number, offset, token, type and pipe type.
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

	PipeType string `json:"pipeType"`
}

// String returns a string representation of the node.
// It uses the json marshaler to convert the node to a string.
func (n *BaseNode) String() string {
	b, _ := json.MarshalIndent(n, "", "  ")
	return string(b)
}

// GetBaseNode returns the base node.
func (n *BaseNode) GetBaseNode() *BaseNode {
	return n
}

// GetType returns the type of the node.
// This is a helper function to make it easier to get the type of a node.
func (n *BaseNode) GetType() NodeType {
	return n.Type
}

// SetPipeType returns the line number of the node in the source code.
// This is used to determine if the node is setup under the piple
// or not. This function is called by the parser
// when it detects a pipe. It performs the recursive call
// to set the pipe type for all child nodes.
func (n *BaseNode) SetPipeType(pipeType string) {
	n.PipeType = pipeType
}
