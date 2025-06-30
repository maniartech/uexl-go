# Chaining Pipes

Multiple pipes can be chained together to perform complex data transformations in a single expression.

## Example
```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// Filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

Each stage receives the output of the previous stage as its input. This enables powerful and readable pipelines for processing data.