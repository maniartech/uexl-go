package parser

type NodeType string

const (
	NodeTypeBinaryExpression NodeType = "BinaryExpression"
	NodeTypeUnaryExpression  NodeType = "UnaryExpression"
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
	Object   Expression
	Property string
	Line     int
	Column   int
}

func (ma *MemberAccess) expressionNode()      {}
func (ma *MemberAccess) Type() NodeType       { return NodeTypeMemberAccess }
func (ma *MemberAccess) Position() (int, int) { return ma.Line, ma.Column }

type PipeExpression struct {
	Expressions []Expression
	PipeTypes   []string
	Aliases     []string // New field to store aliases
	Line        int
	Column      int
}

func (pe *PipeExpression) expressionNode()      {}
func (pe *PipeExpression) Type() NodeType       { return NodeTypePipeExpression }
func (pe *PipeExpression) Position() (int, int) { return pe.Line, pe.Column }
