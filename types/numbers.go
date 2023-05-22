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

// String returns the string representation of the number value.
func (n Number) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

// Add returns the sum of the number and the other number.
func (n Number) Add(other Value) (Value, error) {
	if otherV, ok := other.(Number); ok {
		return n + otherV, nil
	}

	if otherV, ok := other.(Array); ok {
		return append(Array{n}, otherV...), nil
	}

	return nil, errors.New("invalid type for plus")
}

// Substract returns the difference of the number and the other number.
func (n Number) Subtract(other Value) (Value, error) {
	if other, ok := other.(Number); ok {
		return n - other, nil
	}

	return nil, errors.New("invalid type for minus")
}

// Multiply returns the product of the number and the other number.
func (n Number) Multiply(other Value) (Value, error) {

	// If number
	if num, ok := other.(Number); ok {
		return n * num, nil
	}

	// If string
	if str, ok := other.(String); ok {
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

// Mod returns the remainder of the number divided by the other number.
func (n Number) Mod(other Value) (Value, error) {
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

// And returns the logical AND of the number and the other number.
func (n Number) And(other Value) (Value, error) {
	if other, ok := other.(Number); ok {
		return Number(int64(n) & int64(other)), nil
	}

	return nil, errors.New("invalid type for and")
}

// Or returns the logical OR of the number and the other number.
func (n Number) Or(other Value) (Value, error) {
	if other, ok := other.(Number); ok {
		return Number(int64(n) | int64(other)), nil
	}

	return nil, errors.New("invalid type for or")
}

// Xor returns the logical XOR of the number and the other number.
func (n Number) Xor(other Value) (Value, error) {
	if other, ok := other.(Number); ok {
		return Number(int64(n) ^ int64(other)), nil
	}

	return nil, errors.New("invalid type for xor")
}

// Not returns the logical NOT of the number.
func (n Number) Not() (Value, error) {
	return Number(^int64(n)), nil
}

// ShiftLeft returns the number shifted left by the other number.
func (n Number) ShiftLeft(other Value) (Value, error) {
	if other, ok := other.(Number); ok {
		return Number(int64(n) << uint64(other)), nil
	}

	return nil, errors.New("invalid type for shift left")
}

// ShiftLeftUnsigned returns the number shifted left by the other number with zero extension.
func (n Number) ShiftLeftUnsigned(other Value) (Value, error) {
	if other, ok := other.(Number); ok {
		return Number(uint64(n) << uint64(other)), nil
	}

	return nil, errors.New("invalid type for shift left unsigned")
}

// ShiftRight returns the number shifted right by the other number.
func (n Number) ShiftRight(other Value) (Value, error) {
	if other, ok := other.(Number); ok {
		return Number(int64(n) >> uint64(other)), nil
	}

	return nil, errors.New("invalid type for shift right")
}

// ShiftRightUnsigned returns the number shifted right by the other number with zero extension.
func (n Number) ShiftRightUnsigned(other Value) (Value, error) {
	if other, ok := other.(Number); ok {
		return Number(uint64(n) >> uint64(other)), nil
	}

	return nil, errors.New("invalid type for shift right unsigned")
}

// Negate returns the negation of the number.
func (n Number) Negate() (Value, error) {
	return Number(-n), nil
}
