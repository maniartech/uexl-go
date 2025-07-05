package constants

// Language constants and magic strings used throughout the parser

// Common pipe types
const (
	// DefaultPipeType is the default type for pipe expressions
	DefaultPipeType = "pipe"
)

// Parser state constants
const (
	// These constants define various parser states and behaviors
	DefaultParserState = "default"
)

// Operator precedence levels
const (
	PrecedenceLowest  = 0
	PrecedenceOr      = 1  // ||
	PrecedenceAnd     = 2  // &&
	PrecedenceBitOr   = 3  // |
	PrecedenceBitXor  = 4  // ^
	PrecedenceBitAnd  = 5  // &
	PrecedenceEquals  = 6  // == !=
	PrecedenceCompare = 7  // > < >= <=
	PrecedenceShift   = 8  // << >>
	PrecedenceSum     = 9  // + -
	PrecedenceProduct = 10 // * / %
	PrecedencePrefix  = 11 // -x !x ~x
	PrecedenceCall    = 12 // myFunction(x)
	PrecedenceIndex   = 13 // array[index]
	PrecedenceHighest = 14
)

// Operator symbols - centralized string constants
const (
	// Arithmetic operators
	SymbolPlus     = "+"
	SymbolMinus    = "-"
	SymbolMultiply = "*"
	SymbolDivide   = "/"
	SymbolModulo   = "%"

	// Bitwise operators
	SymbolBitwiseAnd = "&"
	SymbolBitwiseOr  = "|"
	SymbolBitwiseXor = "^"
	SymbolBitwiseNot = "~"

	// Shift operators
	SymbolLeftShift  = "<<"
	SymbolRightShift = ">>"

	// Assignment operators
	SymbolAssign           = "="
	SymbolAddAssign        = "+="
	SymbolSubtractAssign   = "-="
	SymbolMultiplyAssign   = "*="
	SymbolDivideAssign     = "/="
	SymbolModuloAssign     = "%="
	SymbolBitwiseAndAssign = "&="
	SymbolBitwiseOrAssign  = "|="
	SymbolBitwiseXorAssign = "^="

	// Increment and decrement operators
	SymbolIncrement = "++"
	SymbolDecrement = "--"

	// Comparison operators
	SymbolEqual          = "=="
	SymbolNotEqual       = "!="
	SymbolGreaterThan    = ">"
	SymbolLessThan       = "<"
	SymbolGreaterOrEqual = ">="
	SymbolLessOrEqual    = "<="

	// Logical operators
	SymbolLogicalAnd = "&&"
	SymbolLogicalOr  = "||"

	// Conditional operators
	SymbolConditional = "?:"

	// Pipe operators
	SymbolPipe      = "|:"
	SymbolNamedPipe = "|"
)

// Special characters and delimiters
const (
	SymbolLeftParen    = "("
	SymbolRightParen   = ")"
	SymbolLeftBracket  = "["
	SymbolRightBracket = "]"
	SymbolLeftBrace    = "{"
	SymbolRightBrace   = "}"
	SymbolComma        = ","
	SymbolDot          = "."
	SymbolColon        = ":"
	SymbolDollar       = "$"
	SymbolAs           = "as"
)

// Literal constants
const (
	LiteralTrue  = "true"
	LiteralFalse = "false"
	LiteralNull  = "null"
)

// Error message constants
const (
	MsgEmptyExpression   = "expression is nil"
	MsgUnexpectedToken   = "unexpected token at end"
	MsgMissingDollarSign = "missing dollar sign"
	MsgAliasInSubExpr    = "alias not allowed in sub-expression"
)
