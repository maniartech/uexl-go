package types

// Number represents a number in the uexl language.
type Number float64

// String returns a string representation of the number.
type String string

// Boolean represents a boolean in the uexl language.
type Boolean bool

// Null represents a null in the uexl language.
type Null struct{}

// Array represents an array in the uexl language.
type Array []interface{}

// Object represents an object in the uexl language.
type Object map[string]interface{}

// Function represents a function in the uexl language.
type Function func(...interface{}) (interface{}, error)

// Identifier represents an identifier in the uexl language.
type Identifier string
