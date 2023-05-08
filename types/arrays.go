package types

import (
	"errors"
	"strings"
)

// Array represents an array in the uexl language.
type Array []Value

// Type returns the type of the array.
func (a Array) Type() string {
	return "array"
}

// IsTruthy returns true if the array is not nil or empty.
func (a Array) IsTruthy() bool {
	return len(a) > 0
}

// IsZero returns true if the array is nil or empty.
func (a Array) IsZero() bool {
	return len(a) == 0
}

// Equal returns true if the given array is equal to this array.
// All elements are comparable using the Equatable interface.
func (a Array) Equals(other Value) bool {
	if otherArray, ok := other.(Array); ok {
		if len(a) != len(otherArray) {
			return false
		}
		for i, v := range a {
			if !v.Equals(otherArray[i]) {
				return false
			}
		}
		return true
	}
	return false
}

// Compare returns -1 if the array is less than the other array, 0 if the array is equal to the other array, 1 if the array is greater than the other array.
// All elements are comparable using the Comparable interface.
func (a Array) Compare(other Value) (int, error) {

	if otherArray, ok := other.(Array); ok {
		if len(a) < len(otherArray) {
			return -1, nil
		} else if len(a) > len(otherArray) {
			return 1, nil
		}
		for i, v := range a {
			if c, err := v.(Array).Compare(otherArray[i]); err != nil {
				return 0, err
			} else if c != 0 {
				return c, nil
			}
		}
		return 0, nil
	}
	return 0, errors.New("invalid type for comparison")
}

// String returns the string representation of the array.
func (a Array) String() string {
	sb := strings.Builder{}
	sb.WriteString("[")
	for i, v := range a {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.String())
	}
	sb.WriteString("]")
	return sb.String()
}

// Add returns the concatenation of the array and the other array.
func (a Array) Add(other Value) (Value, error) {
	if otherArray, ok := other.(Array); ok {
		return append(a, otherArray...), nil
	}

	a = append(a, other)
	return a, nil
}
