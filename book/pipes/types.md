# Pipe Types

UExL supports several built-in pipe types for common data transformations. Each pipe type is designed for a specific pattern of data processing, especially with arrays.

| Pipe Type | Description | Example |
|-----------|-------------|---------|
| `|:` | Simple pipe (default). Passes the value as `$1` to the next stage. | `10 |: $1 * 2` |
| `|map:` | Maps each element in an array to a new value. `$1` is the current element. | `[1, 2, 3] |map: $1 * 2` |
| `|filter:` | Filters elements in an array based on a condition. `$1` is the current element. | `[1, 2, 3, 4] |filter: $1 > 2` |
| `|reduce:` | Reduces an array to a single value. `$1` is the accumulator, `$2` is the current element. | `[1, 2, 3, 4, 5] |reduce: $1 + $2` |

## How Each Pipe Works
- **Simple Pipe (`|:`)**: Passes the result of the left expression as `$1` to the right expression. Useful for chaining single-value transformations.
- **Map Pipe (`|map:`)**: Iterates over each element in an array, applying the right expression to each. Returns a new array.
- **Filter Pipe (`|filter:`)**: Iterates over each element, keeping only those for which the right expression evaluates to `true`.
- **Reduce Pipe (`|reduce:`)**: Combines all elements into a single value using the right expression. `$1` is the accumulator, `$2` is the current element.

## Edge Cases
- If the left side of a map/filter/reduce pipe is not an array, the result is `null`.
- If the filter condition never matches, the result is an empty array.
- Reduce on an empty array returns `null` unless an initial value is provided (if supported).

## Advanced Examples
```
[1, 2, 3] |map: $1 * 2                  // [2, 4, 6]
[1, 2, 3, 4] |filter: $1 > 2            // [3, 4]
[1, 2, 3, 4, 5] |reduce: $1 + $2        // 15
["a", "b", "c"] |map: $1.toUpperCase() // ["A", "B", "C"]
```

Understanding these pipe types is key to effective data transformation in UExL.