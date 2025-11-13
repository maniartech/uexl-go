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
	PrecedenceBitAnd  = 4  // &
	PrecedenceEquals  = 5  // == !=
	PrecedenceCompare = 6  // > < >= <=
	PrecedenceShift   = 7  // << >>
	PrecedenceSum     = 8  // + -
	PrecedenceProduct = 9  // * / %
	PrecedencePower   = 10 // ** (right-associative)
	PrecedencePrefix  = 11 // -x !x ~x
	PrecedenceCall    = 13 // myFunction(x)
	PrecedenceIndex   = 14 // array[index]
	PrecedenceHighest = 15
)

// Operator symbols - centralized string constants
const (
	// Arithmetic operators
	SymbolPlus     = "+"
	SymbolMinus    = "-"
	SymbolMultiply = "*"
	SymbolDivide   = "/"
	SymbolModulo   = "%"
	SymbolPower    = "**" // Legacy alias, also supports "^" (Excel-compatible)
	SymbolPowerAlt = "^"  // Excel-compatible power operator (was XOR in v1.x)

	// Bitwise operators
	SymbolBitwiseAnd = "&"
	SymbolBitwiseOr  = "|"
	SymbolBitwiseXor = "~" // Changed from "^" - Lua-style context-dependent (binary XOR)
	SymbolBitwiseNot = "~" // Same symbol as XOR, context-dependent (unary NOT)

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
	SymbolNotEqualExcel  = "<>" // Excel-compatible not-equals alias
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
