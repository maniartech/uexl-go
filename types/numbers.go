package types

import (
	"errors"
	"strconv"
)

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
func (n Number) Equals(other Value) bool {
	if other, ok := other.(Number); ok {
		return n == other
	}
	return false
}

// Compare returns -1 if the number is less than the other number, 0 if the number is equal to the other number, 1 if the number is greater than the other number.
func (n Number) Compare(other Value) (int, error) {
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

// String returns the string representation of the number value.
func (n Number) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

// Plus returns the sum of the number and the other number.
func (n Number) Plus(other Value) (Value, error) {
	if other, ok := other.(Number); ok {
		return n + other, nil
	}

	return nil, errors.New("invalid type for plus")
}

// Minus returns the difference of the number and the other number.
func (n Number) Minus(other Value) (Value, error) {
	if other, ok := other.(Number); ok {
		return n - other, nil
	}

	return nil, errors.New("invalid type for minus")
}

// Multiply returns the product of the number and the other number.
func (n Number) Multiply(other Value) (Value, error) {
	if num, ok := other.(Number); ok {
		return n * num, nil
	} else if str, ok := other.(String); ok {
		// return the string repeated the number of times specified by the other number.
		var result string
		for i := 0; i < int(n); i++ {
			result += string(str)
		}
		return String(result), nil
	}

	return nil, errors.New("invalid type for multiply")
}

// Divide returns the quotient of the number and the other number.
func (n Number) Divide(other Value) (Value, error) {
	if other, ok := other.(Number); ok {
		return n / other, nil
	}

	return nil, errors.New("invalid type for divide")
}

// Modulo returns the remainder of the number divided by the other number.
func (n Number) Modulo(other Value) (Value, error) {
	if other, ok := other.(Number); ok {
		return Number(float64(n) - (float64(other) * float64(int64(n)/int64(other)))), nil
	}

	return nil, errors.New("invalid type for modulo")
}

// Power returns the number raised to the power of the other number.
func (n Number) Power(other Value) (Value, error) {
	if other, ok := other.(Number); ok {
		return Number(float64(n) * float64(other)), nil
	}

	return nil, errors.New("invalid type for power")
}
