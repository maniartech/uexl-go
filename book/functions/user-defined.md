# User Defined Functions

UExL can be extended with user-defined functions, allowing you to add custom logic and operations.

## Accessing User Defined Functions
- User-defined functions are registered in the host environment (such as an application embedding UExL).
- Once registered, they can be called like any built-in function:
```
myFunc(10, 20)
```

## Example
Suppose you register a function named `double`:
```
double(5) // Returns 10
```

Refer to your host application's documentation for how to register and expose user-defined functions to UExL.