package types

import "errors"

// Number represents a number in the uexl language. It is a float64.
// It implements the Value interface.
type Number float64

// Type returns the type of the number.
func (n Number) Type() string {
	return "number"
}

// IsTruthy returns true if the number is not zero.
func (n Number) IsTruthy() bool {
	return n != 0
}

// IsZero returns true if the number is zero.
func (n Number) IsZero() bool {
	return n == 0
}

// Equals returns true if the number values are equal.
func (n Number) Equals(other interface{}) bool {
	if other, ok := other.(Number); ok {
		return n == other
	}
	return false
}

// Compare returns -1 if the number is less than the other number, 0 if the number is equal to the other number, 1 if the number is greater than the other number.
func (n Number) Compare(other interface{}) (int, error) {
	if other, ok := other.(Number); ok {
		if n < other {
			return -1, nil
		} else if n > other {
			return 1, nil
		}
		return 0, nil
	}

	return 0, errors.New("invalid type for comparison")
}
