package types

import "errors"

// Null represents a null in the uexl language.
type Null struct{}

// Type returns the type of the null.
func (n Null) Type() string {
	return "null"
}

// IsTruthy returns false.
func (n Null) IsTruthy() bool {
	return false
}

// Equals returns true if the null values are equal.
func (n Null) Equals(other interface{}) bool {
	_, ok := other.(Null)
	return ok
}

// Compare returns 0 if the null values are equal.
func (n Null) Compare(other interface{}) (int, error) {
	if _, ok := other.(Null); ok {
		return 0, nil
	}

	return 0, errors.New("invalid type for comparison")
}

// String returns the string representation of the null value.
func (n Null) String() string {
	return "null"
}
