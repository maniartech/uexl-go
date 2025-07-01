# Using UExL in Golang

UExL can be embedded in Go applications to evaluate expressions dynamically.

## Installation
Install the UExL Go package (when available):
```
go get github.com/maniartech/uexl-go
```

## Importing the Library
```go
import "github.com/maniartech/uexlgo"
```

## Basic Usage
```go
result, err := uexl.Eval("10 + 20 |: $1 * 2")
if err != nil {
    // handle error
}
fmt.Println(result) // Output: 60
```

## Registering User-Defined Functions

You can extend UExL in Go by registering custom functions. For example, to register a function named "double":

```go
import "github.com/maniartech/uexlgo"

// Register a function named "double" that multiplies its argument by 2
uexlgo.RegisterFunction("double", func(x float64) float64 {
    return x * 2
})
```

Once registered, you can call this function in UExL expressions:

```uexl
double(10) // returns 20
[1, 2, 3] |map: double($1) // returns [2, 4, 6]
```

Arguments are evaluated before being passed to the function, and type conversion is applied if possible. Functions can return numbers, strings, arrays, objects, or null. If the wrong number or type of arguments is passed, an error is thrown.

Refer to the Go package documentation for advanced integration and custom function registration.
