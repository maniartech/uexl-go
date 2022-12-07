package ast

// NodeNodeType defines the NodeTypes of the nodes in the AST.
type NodeType uint8

const (

	// value types
	NodeTypeNull    NodeType = iota // Null node
	NodeTypeBoolean                 // Boolean node
	NodeTypeNumber                  // Number node
	NodeTypeString                  // String node
	NodeTypeArray                   // Array node
	NodeTypeObject                  // Object node

	// expression types
	NodeTypeExpression      // Expression node
	NodeTypeUnaryExpression // Unary expression node

	// Other types
	NodeTypeIdentifier // Identifier node
	NodeTypeFunc       // Function node
	NodeTypeOperator   // Operator node

)

// Returns the string representation of the node type.
func (n NodeType) String() string {
	switch n {

	// Null node
	case NodeTypeNull:
		return "null"

	// Boolean node
	case NodeTypeBoolean:
		return "boolean"

	// Number node
	case NodeTypeNumber:
		return "number"

	// String node
	case NodeTypeString:
		return "string"

	// Array node
	case NodeTypeArray:
		return "array"

	// Object node
	case NodeTypeObject:
		return "object"

	// Identifier node
	case NodeTypeIdentifier:
		return "identifier"

	// Expession Node
	case NodeTypeExpression:
		return "expression"

	// Unary expression node
	case NodeTypeUnaryExpression:
		return "unary expression"

	// Function node
	case NodeTypeFunc:
		return "function"
	}

	// Unknown node
	return ErrUnknownNodeType
}
