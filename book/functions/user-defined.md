# User-Defined Functions

UExL allows you to extend the language by registering custom functions in the host environment. These functions can then be called just like built-in functions.

## Registering Functions
- User-defined functions are registered in the host language and made available to UExL expressions.
- The function name must be unique and follow identifier rules.
- Functions can accept any number of arguments and return any data type supported by UExL.

## Calling User-Defined Functions
- Call user-defined functions just like built-in ones:
  `double(10)` // returns 20
- They can be used in pipes, maps, filters, and nested expressions.

## Argument Handling
- Arguments are evaluated before being passed to the function.
- Type conversion is applied if possible.
- If the wrong number or type of arguments is passed, an error is thrown.

## Return Values
- Functions can return numbers, strings, arrays, objects, or `null`.
- Returning `null` is valid and can be used to indicate missing or invalid data.

## Advanced Usage
- User-defined functions can be used in pipes:
  `[1, 2, 3] |map: double($1)`
- Functions can be composed:
  `sum(map([1, 2, 3], double))`
- Functions can return other functions (if supported by the host environment).

## Edge Cases
- Registering a function with a name that conflicts with a built-in function will override the built-in.
- Passing `null` or invalid arguments may result in errors or `null` results.

User-defined functions make UExL highly extensible for custom logic and domain-specific needs.

## Accessing User Defined Functions
- User-defined functions are registered in the host environment (such as an application embedding UExL).
- Once registered, they can be called like any built-in function:
```
myFunc(10, 20)
```

Refer to your host application's documentation for how to register and expose user-defined functions to UExL.

## Example
Suppose you register a function named `double`:
```
double(5) // Returns 10
```

Refer to your host application's documentation for how to register and expose user-defined functions to UExL.