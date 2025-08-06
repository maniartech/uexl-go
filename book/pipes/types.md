# Pipe Types

UExL supports several built-in pipe types for common data transformations. Each pipe type is designed for a specific pattern of data processing, especially with arrays.

| Pipe Type   | Description                                                                 | Special Variables         | Example                                   |
|-------------|-----------------------------------------------------------------------------|--------------------------|-------------------------------------------|
| `\|:`        | Simple pipe (default). Passes the value as `$last` to the next stage.       | `$last`                  | `10 \|: $last * 2`                        |
| `\|map:`     | Maps each element in an array to a new value.                               | `$item`, `$index`        | `[1, 2, 3] \|map: $item * 2`              |
| `\|filter:`  | Filters elements in an array based on a condition.                          | `$item`, `$index`        | `[1, 2, 3, 4] \|filter: $item > 2`        |
| `\|reduce:`  | Reduces an array to a single value.                                         | `$acc`, `$item`, `$index`\| `[1, 2, 3, 4, 5] \|reduce: $acc + $item`  |

## How Each Pipe Works

- **Simple Pipe (`\|:`)**: Passes the result of the left expression as `$last` to the right expression. Useful for chaining single-value transformations.
- **Map Pipe (`\|map:`)**: Iterates over each element in an array, applying the right expression to each. The right side can use `$item` (current element) and `$index` (current index). Returns a new array.
- **Filter Pipe (`\|filter:`)**: Iterates over each element, keeping only those for which the right expression (using `$item`, `$index`) evaluates to `true`.
- **Reduce Pipe (`\|reduce:`)**: Combines all elements into a single value using the right expression. The right side can use `$acc` (accumulator), `$item` (current element), and `$index` (current index).

## Edge Cases

- If the left side of a map/filter/reduce pipe is not an array, the result is `null`.
- If the filter condition never matches, the result is an empty array.
- Reduce on an empty array returns `null` unless an initial value is provided (if supported).

## Advanced Examples

```
[1, 2, 3] |map: $item * 2                       // [2, 4, 6]
[1, 2, 3, 4] |filter: $item > 2                 // [3, 4]
[1, 2, 3, 4, 5] |reduce: $acc + $item           // 15
[1, 2, 3, 4, 5] |reduce: $acc + $item * $index  // 40
["a", "b", "c"] |map: $item.toUpperCase()       // ["A", "B", "C"]
```

Understanding these pipe types is key to effective data transformation in UExL.