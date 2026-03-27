// Syntax error
{a: 1,, b: 2}
// => errorCode: "SyntaxError"

// Type error
"abc" * 2
// => errorCode: "TypeError"

// Reference error
x + 1
// => errorCode: "ReferenceError" (if x is not defined)

// Argument error
min()
// => errorCode: "ArgumentError"