package core

// Number represents a number in the uexl language.
type Number float64

// String returns a string representation of the number.
type String string

// Boolean represents a boolean in the uexl language.
type Boolean bool

// Null represents a null in the uexl language.
type Null struct{}

// Array represents an array in the uexl language.
type Array []any

// Object represents an object in the uexl language.
type Object map[string]any

// Function represents a function in the uexl language.
type Function func(...any) (any, error)

// Identifier represents an identifier in the uexl language.
type Identifier string

// Value represents a value in the uexl language.
type Value interface {
	Number | String | Boolean | Null | Array | Object
}
