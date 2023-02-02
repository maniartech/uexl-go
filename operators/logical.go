package operators

import "github.com/maniartech/uexl_go/evaluators"

func LogicalAnd(a, b any) (any, error) {
	return evaluators.IsTruthy(a) && evaluators.IsTruthy(b), nil
}

func LogicalOr(a, b any) (any, error) {
	if evaluators.IsTruthy(a) {
		return a, nil
	}

	if evaluators.IsTruthy(b) {
		return b, nil
	}

	return false, nil
}

func LogicalXor(a, b any) (any, error) {
	return evaluators.IsTruthy(a) != evaluators.IsTruthy(b), nil
}

func LogicalNand(a, b any) (any, error) {
	return !(evaluators.IsTruthy(a) && evaluators.IsTruthy(b)), nil
}

// Urinary operators
func LogicalNot(a any) (any, error) {
	return !evaluators.IsTruthy(a), nil
}

func init() {
	Registry.Register("&&", LogicalAnd)
	Registry.Register("||", LogicalOr)
}
