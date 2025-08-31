package parser

type NodeType string

const (
	NodeTypeBinaryExpression NodeType = "BinaryExpression"
	NodeTypeUnaryExpression  NodeType = "UnaryExpression"
	NodeTypeConditional      NodeType = "ConditionalExpression"
	NodeTypeNumberLiteral    NodeType = "NumberLiteral"
	NodeTypeStringLiteral    NodeType = "StringLiteral"
	NodeTypeBooleanLiteral   NodeType = "BooleanLiteral"
	NodeTypeNullLiteral      NodeType = "NullLiteral"
	NodeTypeIdentifier       NodeType = "Identifier"
	NodeTypeArrayLiteral     NodeType = "ArrayLiteral"
	NodeTypeObjectLiteral    NodeType = "ObjectLiteral"
	NodeTypeFunctionCall     NodeType = "FunctionCall"
	NodeTypeMemberAccess     NodeType = "MemberAccess"
	NodeTypePipeExpression   NodeType = "PipeExpression"
	NodeTypeProgram          NodeType = "Program"
)

type Node interface {
	Type() NodeType
	Position() (line, column int)
}

type Expression interface {
	Node
	expressionNode()
}

type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
	Line     int
	Column   int
}

func (be *BinaryExpression) expressionNode()      {}
func (be *BinaryExpression) Type() NodeType       { return NodeTypeBinaryExpression }
func (be *BinaryExpression) Position() (int, int) { return be.Line, be.Column }

// ConditionalExpression represents the ternary operator: condition ? consequent : alternate
type ConditionalExpression struct {
	Condition  Expression
	Consequent Expression
	Alternate  Expression
	Line       int
	Column     int
}

func (ce *ConditionalExpression) expressionNode()      {}
func (ce *ConditionalExpression) Type() NodeType       { return NodeTypeConditional }
func (ce *ConditionalExpression) Position() (int, int) { return ce.Line, ce.Column }

type UnaryExpression struct {
	Operator string
	Operand  Expression
	Line     int
	Column   int
}

func (ue *UnaryExpression) expressionNode()      {}
func (ue *UnaryExpression) Type() NodeType       { return NodeTypeUnaryExpression }
func (ue *UnaryExpression) Position() (int, int) { return ue.Line, ue.Column }

type NumberLiteral struct {
	Value  float64 // Parsed numeric value
	Line   int
	Column int
}

func (nl *NumberLiteral) expressionNode()      {}
func (nl *NumberLiteral) Type() NodeType       { return NodeTypeNumberLiteral }
func (nl *NumberLiteral) Position() (int, int) { return nl.Line, nl.Column }

type StringLiteral struct {
	Value          string
	Token          string // The original token, including quotes and escapes
	IsRaw          bool   // Track if this was a raw string
	IsSingleQuoted bool   // Track if this was a single-quoted string
	Line           int
	Column         int
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) Type() NodeType       { return NodeTypeStringLiteral }
func (sl *StringLiteral) Position() (int, int) { return sl.Line, sl.Column }

type BooleanLiteral struct {
	Value  bool
	Line   int
	Column int
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) Type() NodeType       { return NodeTypeBooleanLiteral }
func (bl *BooleanLiteral) Position() (int, int) { return bl.Line, bl.Column }

type NullLiteral struct {
	Line   int
	Column int
}

func (nl *NullLiteral) expressionNode()      {}
func (nl *NullLiteral) Type() NodeType       { return NodeTypeNullLiteral }
func (nl *NullLiteral) Position() (int, int) { return nl.Line, nl.Column }

type Identifier struct {
	Name   string
	Line   int
	Column int
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) Type() NodeType       { return NodeTypeIdentifier }
func (i *Identifier) Position() (int, int) { return i.Line, i.Column }

type ArrayLiteral struct {
	Elements []Expression
	Line     int
	Column   int
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) Type() NodeType       { return NodeTypeArrayLiteral }
func (al *ArrayLiteral) Position() (int, int) { return al.Line, al.Column }

type ObjectLiteral struct {
	Properties map[string]Expression
	Line       int
	Column     int
}

func (ol *ObjectLiteral) expressionNode()      {}
func (ol *ObjectLiteral) Type() NodeType       { return NodeTypeObjectLiteral }
func (ol *ObjectLiteral) Position() (int, int) { return ol.Line, ol.Column }

type FunctionCall struct {
	Function  Expression
	Arguments []Expression
	Line      int
	Column    int
}

func (fc *FunctionCall) expressionNode()      {}
func (fc *FunctionCall) Type() NodeType       { return NodeTypeFunctionCall }
func (fc *FunctionCall) Position() (int, int) { return fc.Line, fc.Column }

type MemberAccess struct {
	Target   Expression
	Property any
	Optional bool
	Line     int
	Column   int
}

func (ma *MemberAccess) expressionNode()      {}
func (ma *MemberAccess) Type() NodeType       { return NodeTypeMemberAccess }
func (ma *MemberAccess) Position() (int, int) { return ma.Line, ma.Column }

// IndexAccess represents array index access expressions
type IndexAccess struct {
	Target   Expression
	Index    Expression
	Optional bool
	Line     int
	Column   int
}

func (ia *IndexAccess) expressionNode()      {}
func (ia *IndexAccess) Position() (int, int) { return ia.Line, ia.Column }
func (ia *IndexAccess) Type() NodeType       { return NodeType("IndexAccess") }

type PipeExpression struct {
	Expression Expression // The pipe's predicate expression block
	PipeType   string
	Alias      string
	Index      int // Index of the predicate block
	Line       int
	Column     int
}

func (pe *PipeExpression) expressionNode()      {}
func (pe *PipeExpression) Type() NodeType       { return NodeTypePipeExpression }
func (pe *PipeExpression) Position() (int, int) { return pe.Line, pe.Column }

type ProgramNode struct {
	PipeExpressions []PipeExpression
	Line            int
	Column          int
}

func (pn *ProgramNode) expressionNode()      {}
func (pn *ProgramNode) Type() NodeType       { return NodeTypeProgram }
func (pn *ProgramNode) Position() (int, int) { return pn.Line, pn.Column }
