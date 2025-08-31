# Slicing Operator

The slicing operator in UExL provides a powerful and concise way to extract subsequences from arrays and strings. Its design is heavily inspired by Python's slicing mechanism, offering flexibility with positive and negative indexing, and optional start, end, and step parameters.

## Syntax

The basic syntax for the slicing operator is:

```uexl
sequence[start:end:step]
```

- `sequence`: The array or string to be sliced.
- `start`: The index of the first element to include in the slice. If omitted, it defaults to `0` for positive steps and the end of the sequence for negative steps.
- `end`: The index of the first element to exclude from the slice. If omitted, it defaults to the end of the sequence for positive steps and the beginning of the sequence for negative steps.
- `step`: The increment between indices. It can be positive or negative. If omitted, it defaults to `1`.

All three parameters (`start`, `end`, and `step`) are optional.

## Behavior

The slicing operator is applicable to both `Array` and `String` types.

- **For Arrays**: It returns a new array containing the elements from the original array within the specified range.
- **For Strings**: It returns a new string containing the characters from the original string within the specified range.

The original sequence is not modified.

### Negative Indexing

Negative indices allow for indexing from the end of the sequence. `-1` refers to the last element, `-2` to the second-to-last, and so on.

## Examples

Let's consider an array `arr = [10, 20, 30, 40, 50]` and a string `str = "UExL"`.

### Basic Slicing

- `arr[1:4]` → `[20, 30, 40]`
- `str[1:3]` → `"Ex"`

### Omitting `start` or `end`

- `arr[:3]` → `[10, 20, 30]` (from the beginning up to index 3)
- `str[1:]` → `"ExL"` (from index 1 to the end)
- `arr[:]` → `[10, 20, 30, 40, 50]` (a shallow copy of the entire array)

### Negative Indexing

- `arr[-1]` → `50` (accessing the last element)
- `arr[:-1]` → `[10, 20, 30, 40]` (all elements except the last one)
- `str[-3:-1]` → `"Ex"` (from the 3rd last to the 1st last element)

### Using `step`

- `arr[0:5:2]` → `[10, 30, 50]` (every second element)
- `str[::2]` → `"UL"`

### Negative `step` (Reversing)

A negative step value reverses the direction of the slicing.

- `arr[::-1]` → `[50, 40, 30, 20, 10]` (the entire array, reversed)
- `str[::-1]` → `"LxEU"`
- `arr[4:1:-1]` → `[50, 40, 30]`

## Edge Cases

- If `start` or `end` are out of bounds, they are clamped to the valid range of indices for the sequence.
- If `start` is greater than or equal to `end` with a positive `step`, an empty sequence (`[]` or `""`) is returned.
- If `start` is less than or equal to `end` with a negative `step`, an empty sequence is returned.

## Implementation Plan

To support the slicing operator in UExL, we will introduce a new Abstract Syntax Tree (AST) node and enhance the parser to handle the slicing syntax.

### AST Node

A new AST node, `SliceExpression`, will be defined to represent the slicing operation. This node will capture the target of the slice, as well as the optional `start`, `end`, and `step` expressions.

The proposed structure for the `SliceExpression` node is as follows:

```go
// SliceExpression represents a slicing operation on a sequence (array or string).
type SliceExpression struct {
    Target Expression // The sequence being sliced
    Start  Expression // The start index (optional, can be nil)
    End    Expression // The end index (optional, can be nil)
    Step   Expression // The step value (optional, can be nil)
    Line   int
    Column   int
}
```

This new node will be added to the `ast/expressions.go` file.

### Parser Update

The parser will be updated to distinguish between a simple index access (`[index]`) and a slice expression (`[start:end:step]`). This logic will be implemented within the `parseMemberAccess` function in `parser/parser.go`.

When the parser encounters a `[` token, it will look ahead for a `:` token to determine whether it is parsing an index or a slice.

- If a `:` is present, the parser will proceed to parse the `start`, `end`, and `step` expressions. It will correctly handle cases where any of these expressions are omitted.
- If no `:` is present before the closing `]`, the expression will be parsed as a standard `IndexAccess` node.

This approach ensures that the new slicing syntax is integrated smoothly with the existing array and string access logic, maintaining backward compatibility while extending the language's capabilities.
