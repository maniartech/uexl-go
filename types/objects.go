package types

import "errors"

// Object represents an object in the uexl language.
type Object map[string]Value

// Type returns the type of the object.
func (o Object) Type() string {
	return "object"
}

// IsTruthy returns true if the object is not nil or empty.
func (o Object) IsTruthy() bool {
	return len(o) > 0
}

// IsZero returns true if the object is nil or empty.
func (o Object) IsZero() bool {
	return len(o) == 0
}

// Equals returns true if the object values are equal. It performs a deep comparison.
func (o Object) Equals(other Value) bool {
	if otherObject, ok := other.(Object); ok {
		if len(o) != len(otherObject) {
			return false
		}

		for k, v := range o {
			if otherV, ok := otherObject[k]; ok {
				if !otherV.Equals(v) {
					return false
				}
			} else {
				return false
			}
		}

		// All keys and values are equal
		return true
	}

	// Not an Object
	return false
}

// Compare returns 0 if the object values are equal, -1 if the object is nil or empty, and 1 if the object is not nil or empty.
func (o Object) Compare(other Value) (int, error) {
	return 0, errors.New("comparision not supported for objects")
}

// String returns the string representation of the object.
func (o Object) String() string {
	return "object"
}
