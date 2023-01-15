package types

// Number represents a number in the uexl language.
type Number float64

// IsTrue returns true if the number is not zero.
func (n Number) IsTrue() bool {
	return n != 0
}

// String returns a string representation of the number.
type String string

// IsTrue returns true if the string is not empty.
func (s String) IsTrue() bool {
	return s != ""
}

// Boolean represents a boolean in the uexl language.
type Boolean bool

// IsTrue returns the boolean value.
func (b Boolean) IsTrue() bool {
	return bool(b)
}

// Null represents a null in the uexl language.
type Null struct{}

// IsTrue returns false.
func (n Null) IsTrue() bool {
	return false
}

// Array represents an array in the uexl language.
type Array []interface{}

// IsTrue returns true if the array is not nil or empty.
func (a Array) IsTrue() bool {
	return len(a) > 0
}

// Object represents an object in the uexl language.
type Object map[string]interface{}

// IsTrue returns true if the object is not nil or empty.
func (o Object) IsTrue() bool {
	return len(o) > 0
}

// Function represents a function in the uexl language.
type Function func(...interface{}) (interface{}, error)

// Identifier represents an identifier in the uexl language.
type Identifier string
