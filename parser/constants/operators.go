package constants

// Operator represents the different operators supported by the UExl language
type Operator int

// Operator constants - these represent the different operators
// that can be used in expressions
const (
	// Arithmetic operators
	OperatorPlus     Operator = iota // +
	OperatorMinus                    // -
	OperatorMultiply                 // *
	OperatorDivide                   // /
	OperatorModulo                   // %
	OperatorPower                    // ** (exponentiation)

	// Bitwise operators
	OperatorBitwiseAnd // &
	OperatorBitwiseOr  // |
	OperatorBitwiseXor // ^ (XOR)
	OperatorBitwiseNot // ~ (bitwise NOT)

	// Shift operators
	OperatorLeftShift  // <<
	OperatorRightShift // >>

	// Assignment operators
	OperatorAssign           // =
	OperatorAddAssign        // +=
	OperatorSubtractAssign   // -=
	OperatorMultiplyAssign   // *=
	OperatorDivideAssign     // /=
	OperatorModuloAssign     // %=
	OperatorBitwiseAndAssign // &=
	OperatorBitwiseOrAssign  // |=
	OperatorBitwiseXorAssign // ^=

	// Increment and decrement operators
	OperatorIncrement // ++
	OperatorDecrement // --

	// Comparison operators
	OperatorEqual          // ==
	OperatorNotEqual       // !=
	OperatorGreaterThan    // >
	OperatorLessThan       // <
	OperatorGreaterOrEqual // >=
	OperatorLessOrEqual    // <=

	// Logical operators
	OperatorLogicalAnd // &&
	OperatorLogicalOr  // ||

	// Conditional operators
	OperatorConditional // ? :

	// Pipe operators
	OperatorPipe      // |:
	OperatorNamedPipe // |map: |filter:
)

// String returns the string representation of an operator
func (o Operator) String() string {
	switch o {
	case OperatorPlus:
		return "+"
	case OperatorMinus:
		return "-"
	case OperatorMultiply:
		return "*"
	case OperatorDivide:
		return "/"
	case OperatorModulo:
		return "%"
	case OperatorPower:
		return "**"
	case OperatorBitwiseAnd:
		return "&"
	case OperatorBitwiseOr:
		return "|"
	case OperatorBitwiseXor:
		return "^"
	case OperatorBitwiseNot:
		return "~"
	case OperatorLeftShift:
		return "<<"
	case OperatorRightShift:
		return ">>"
	case OperatorAssign:
		return "="
	case OperatorAddAssign:
		return "+="
	case OperatorSubtractAssign:
		return "-="
	case OperatorMultiplyAssign:
		return "*="
	case OperatorDivideAssign:
		return "/="
	case OperatorModuloAssign:
		return "%="
	case OperatorBitwiseAndAssign:
		return "&="
	case OperatorBitwiseOrAssign:
		return "|="
	case OperatorBitwiseXorAssign:
		return "^="
	case OperatorIncrement:
		return "++"
	case OperatorDecrement:
		return "--"
	case OperatorEqual:
		return "=="
	case OperatorNotEqual:
		return "!="
	case OperatorGreaterThan:
		return ">"
	case OperatorLessThan:
		return "<"
	case OperatorGreaterOrEqual:
		return ">="
	case OperatorLessOrEqual:
		return "<="
	case OperatorLogicalAnd:
		return "&&"
	case OperatorLogicalOr:
		return "||"
	case OperatorConditional:
		return "?:"
	case OperatorPipe:
		return "|:"
	case OperatorNamedPipe:
		return "|named"
	default:
		return "Unknown"
	}
}

// StartOperator returns true if the given character is the start of an operator
// TODO: This function is currently not used in the codebase. Check and remove if not needed.
func StartOperator(ch rune) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '%' ||
		ch == '&' || ch == '|' || ch == '^' || ch == '~' ||
		ch == '=' || ch == '>' || ch == '<' || ch == '!' || ch == '?'
}
