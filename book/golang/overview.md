# Using UExL in Golang

UExL can be embedded in Go applications to evaluate expressions dynamically.

## Installation
Install the UExL Go package (when available):
```
go get github.com/maniartech/uexl-go
```

## Importing the Library
```go
import "github.com/maniartech/uexl-go"
```

## Basic Usage
```go
result, err := uexl.Eval("10 + 20 |: $1 * 2")
if err != nil {
    // handle error
}
fmt.Println(result) // Output: 60
```

Refer to the Go package documentation for advanced integration and custom function registration.