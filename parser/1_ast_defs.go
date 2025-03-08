package parser

type Node interface {
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
func (be *BinaryExpression) Position() (int, int) { return be.Line, be.Column }

type UnaryExpression struct {
	Operator string
	Operand  Expression
	Line     int
	Column   int
}

func (ue *UnaryExpression) expressionNode()      {}
func (ue *UnaryExpression) Position() (int, int) { return ue.Line, ue.Column }

type NumberLiteral struct {
	Value  string
	Line   int
	Column int
}

func (nl *NumberLiteral) expressionNode()      {}
func (nl *NumberLiteral) Position() (int, int) { return nl.Line, nl.Column }

type StringLiteral struct {
	Value  string
	Line   int
	Column int
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) Position() (int, int) { return sl.Line, sl.Column }

type BooleanLiteral struct {
	Value  bool
	Line   int
	Column int
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) Position() (int, int) { return bl.Line, bl.Column }

type NullLiteral struct {
	Line   int
	Column int
}

func (nl *NullLiteral) expressionNode()      {}
func (nl *NullLiteral) Position() (int, int) { return nl.Line, nl.Column }

type Identifier struct {
	Name   string
	Line   int
	Column int
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) Position() (int, int) { return i.Line, i.Column }

type ArrayLiteral struct {
	Elements []Expression
	Line     int
	Column   int
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) Position() (int, int) { return al.Line, al.Column }

type ObjectLiteral struct {
	Properties map[string]Expression
	Line       int
	Column     int
}

func (ol *ObjectLiteral) expressionNode()      {}
func (ol *ObjectLiteral) Position() (int, int) { return ol.Line, ol.Column }

type FunctionCall struct {
	Function  Expression
	Arguments []Expression
	Line      int
	Column    int
}

func (fc *FunctionCall) expressionNode()      {}
func (fc *FunctionCall) Position() (int, int) { return fc.Line, fc.Column }

type MemberAccess struct {
	Object   Expression
	Property string
	Line     int
	Column   int
}

func (ma *MemberAccess) expressionNode()      {}
func (ma *MemberAccess) Position() (int, int) { return ma.Line, ma.Column }

type PipeExpression struct {
	Expressions []Expression
	PipeTypes   []string
	Aliases     []string // New field to store aliases
	Line        int
	Column      int
}

func (pe *PipeExpression) expressionNode()      {}
func (pe *PipeExpression) Position() (int, int) { return pe.Line, pe.Column }
