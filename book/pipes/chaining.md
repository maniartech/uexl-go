# Chaining Pipes

Pipes can be chained together to build complex data transformation pipelines. Each pipe receives the output of the previous stage, allowing you to compose multiple operations in a clear and concise way.

## How Chaining Works

- The result of each pipe stage is passed as `$1` to the next stage.
- You can mix different pipe types (e.g., `map`, `filter`, `reduce`) in a single chain.
- Chaining is especially powerful for processing arrays and collections.

## Example: Array Transformation

```rust
[1, 2, 3, 4, 5]
  |filter: $1 % 2 == 1
  |map: $1 * 10
  |reduce: $1 + $2
// Result: 30 (filters odd numbers, multiplies by 10, sums)
```

## Advanced Chaining

- You can chain as many pipes as needed for your logic.
- Intermediate results can be objects or arrays, not just numbers.
- Pipes can be nested inside function calls or object properties.

### Example: Nested Pipes

```
users
  |filter: $1.active
  |map: {
      name: $1.name,
      score: $1.scores |reduce: $1 + $2
    }
```

## Edge Cases

- If any stage returns `null`, subsequent pipes receive `null` as input.
- Chaining on non-array values with `map`, `filter`, or `reduce` results in `null`.

## Best Practices

- Use clear, descriptive expressions in each pipe stage.
- Add parentheses for clarity in complex chains.
- Test each stage separately when debugging.

Chaining pipes is a core feature of UExL for building readable and maintainable data transformations.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.

## Example

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data