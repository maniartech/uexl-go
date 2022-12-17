package ast

// NodeNodeType defines the NodeTypes of the nodes in the AST.
type NodeType string

const (

	// value types
	NodeTypeNull    = "null"    // Null node
	NodeTypeBoolean = "boolean" // Boolean node
	NodeTypeNumber  = "number"  // Number node
	NodeTypeString  = "string"  // String node
	NodeTypeArray   = "array"   // Array node
	NodeTypeObject  = "object"  // Object node

	// expression types
	NodeTypeExpression      = "expression"       // Expression node
	NodeTypeUnaryExpression = "unary expression" // Unary expression node

	// Other types
	NodeTypeIdentifier = "identifier" // Identifier node
	NodeTypeFunc       = "function"   // Function node
	NodeTypeOperator   = "operator"   // Operator node
	NodeTypePipe       = "pipe"       // Pipe node

)
