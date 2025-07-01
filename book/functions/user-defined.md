# User-Defined Functions

User-defined functions are a powerful feature of UExL, allowing you to extend the language with custom logic tailored to your application's needs. In this chapter, you'll learn what user-defined functions are, how to register and use them, and best practices for leveraging them in your own projects.

## What Are User-Defined Functions?
User-defined functions let you add new operations to UExL by writing code in your host environment (such as Go, JavaScript, or Python). Once registered, these functions can be called from any UExL expression, just like built-in functions.

## Why Use User-Defined Functions?
- **Custom Logic:** Implement domain-specific calculations or business rules.
- **Reusability:** Encapsulate complex operations for reuse across multiple expressions.
- **Integration:** Bridge UExL with your application's data and services.

## Registering Functions
To make a function available in UExL, you must register it in your host environment. The registration process depends on your programming language or framework. The function name must be unique and follow identifier rules. Functions can accept any number of arguments and return any data type supported by UExL.

## Using User-Defined Functions in Expressions
Once registered, you can call your function just like any built-in function:
```
double(10) // returns 20
```
User-defined functions can be used in pipes, maps, filters, and nested expressions:
```
[1, 2, 3] |map: double($1)
```

## Argument Handling and Return Values
- Arguments are evaluated before being passed to the function.
- Type conversion is applied if possible.
- If the wrong number or type of arguments is passed, an error is thrown.
- Functions can return numbers, strings, arrays, objects, or `null`.
- Returning `null` is valid and can be used to indicate missing or invalid data.

## Advanced Usage
- Compose user-defined functions with built-in ones:
  `sum(map([1, 2, 3], double))`
- Functions can return other functions (if supported by the host environment).

## Edge Cases and Best Practices
- Registering a function with a name that conflicts with a built-in function will override the built-in.
- Passing `null` or invalid arguments may result in errors or `null` results.
- Use clear, descriptive names for your functions to avoid confusion.

User-defined functions make UExL highly extensible for custom logic and domain-specific needs. Refer to your host application's documentation for details on how to register and expose user-defined functions to UExL.

## Example
Suppose you register a function named `double`:
```
double(5) // Returns 10
```

With user-defined functions, you can unlock the full power of UExL and tailor it to your unique requirements.

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
