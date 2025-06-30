# Pipe Types

UExL supports several built-in pipe types for common data transformations:

| Pipe Type | Description | Example |
|-----------|-------------|---------|
| `|:` | Simple pipe (default) | `10 |: $1 * 2` |
| `|map:` | Maps each element in an array | `[1, 2, 3] |map: $1 * 2` |
| `|filter:` | Filters elements in an array | `[1, 2, 3, 4] |filter: $1 > 2` |
| `|reduce:` | Reduces an array to a single value | `[1, 2, 3, 4, 5] |reduce: $1 + $2` |

## Examples
```
[1, 2, 3] |map: $1 * 2                  // [2, 4, 6]
[1, 2, 3, 4] |filter: $1 > 2            // [3, 4]
[1, 2, 3, 4, 5] |reduce: $1 + $2        // 15
```