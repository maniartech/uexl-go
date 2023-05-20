package ast

type OperatorType string

const (

	// ArithmeticOperator evaluates the expression by applying arithmatic operations
	// on the oparands.
	ArithmeticOperator = "arithmetic" // + - * / //

	// ComparisonOperator evaluates the expression by comparting the operands.
	ComparisonOperator = "comparison" // = <> < > <= >=

	// LogicalOperator evaluates the expression by applying logical operations
	LogicalOperator = "logical" // AND OR

	// BitwiseOperator evaluates the expression by applying bitwise operations
	BitwiseOperator = "bitwise" // & | ^ << >>

	// DotOperator evaluates the expression by applying dot operations
	DotOperator = "dot" // .
)

func GetOperatorType(op string) OperatorType {
	switch op {
	case "+", "-", "*", "/", "%":
		return ArithmeticOperator

	case "==", "!=", "<", ">", "<=", ">=":
		return ComparisonOperator

	case "AND", "OR":
		return LogicalOperator

	case "&", "|", "^", "<<", ">>":
		return BitwiseOperator

	case ".":
		return DotOperator

	default:
		return "unknown"
	}
}
