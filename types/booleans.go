package types

import "errors"

// Boolean represents a boolean in the uexl language.
type Boolean bool

// Type returns the type of the boolean.
func (b Boolean) Type() string {
	return "boolean"
}

// IsTruthy returns the boolean value.
func (b Boolean) IsTruthy() bool {
	return bool(b)
}

// IsZero returns true if the boolean value is false.
func (b Boolean) IsZero() bool {
	return !bool(b)
}

// Equals returns true if the boolean values are equal.
func (b Boolean) Equals(other Value) bool {
	if other, ok := other.(Boolean); ok {
		return b == other
	}
	return false
}

func (b Boolean) Add(other Value) (Value, error) {
	if otherV, ok := other.(Array); ok {
		return append(Array{b}, otherV...), nil
	}

	if otherV, ok := other.(String); ok {
		return String(b.String()) + otherV, nil
	}

	return nil, errors.New("invalid type for addition")
}

// Compare returns 0 if the boolean values are equal, -1 if the boolean value is false, and 1 if the boolean value is true.
func (b Boolean) Compare(other Value) (int, error) {
	if other, ok := other.(Boolean); ok {
		if b == other {
			return 0, nil
		} else if b {
			return 1, nil
		}
		return -1, nil
	}

	return 0, errors.New("invalid type for comparison")
}

// String returns the string representation of the boolean value.
func (b Boolean) String() string {
	if b {
		return "true"
	}
	return "false"
}
