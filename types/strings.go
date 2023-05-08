package types

import "errors"

// String returns a string representation of the number.
type String string

// Type returns the type of the string.
func (s String) Type() string {
	return "string"
}

// IsTruthy returns true if the string is not empty.
func (s String) IsTruthy() bool {
	return s != ""
}

// IsZero returns true if the string is empty.
func (s String) IsZero() bool {
	return s == ""
}

// Equals returns true if the string values are equal.
func (s String) Equals(other Value) bool {
	if other, ok := other.(String); ok {
		return s == other
	}
	return false
}

// Compare returns -1 if the string is less than the other string, 0 if the string is equal to the other string, 1 if the string is greater than the other string.
func (s String) Compare(other Value) (int, error) {
	if other, ok := other.(String); ok {
		if s < other {
			return -1, nil
		} else if s > other {
			return 1, nil
		}
		return 0, nil
	}

	return 0, errors.New("invalid type for comparison")
}

// String returns the string representation of the string value.
func (s String) String() string {
	return string(s)
}

// Add returns the concatenation of the string and the other string.
func (s String) Add(other Value) (Value, error) {
	if other, ok := other.(String); ok {
		return s + other, nil
	}

	return nil, errors.New("invalid type for plus")
}

// Multiply returns the string repeated the number of times specified by the other number.
func (s String) Multiply(other Value) (Value, error) {
	if num, ok := other.(Number); ok {
		result := ""
		for i := 0; i < int(num); i++ {
			result += string(s)
		}
		return String(result), nil
	}

	return nil, errors.New("invalid type for multiply")
}
