# Pipes

Pipes are one of the most powerful features of UExL, enabling you to chain operations and build expressive data transformation pipelines. It allows you to process the expression in a sequential manner, through series of pipes. Each pipe takes the output of the previous stage as input and performs a specific operation, such as mapping, filtering, or reducing the data. This allows your expressions to remain clean and readable, without the need for complex nested function calls or complex control structures.

In this chapter, you'll learn what pipes are, how they work, and how to use them effectively with practical examples.

## What Are Pipes?
A pipe takes the output of one expression and passes it as input to the next stage. This lets you build pipelines that process data step by step, making your logic more readable and maintainable.

- The value from the previous stage is accessible as `$1` in the next stage.
- Pipes can be chained to perform multiple transformations in sequence.
- Pipes are especially useful for working with arrays and collections.

## Pipe Syntax
```
expression |: next_expression
expression |map: next_expression
```
- `|:` is the default pipe, passing the value as `$1`.
- `|map:`, `|filter:`, and `|reduce:` are specialized pipes for array processing.

## Practical Examples
- **Transforming data:**
  - `[1, 2, 3] |map: $1 * 2` // Returns `[2, 4, 6]`
- **Filtering:**
  - `users |filter: $1.active` // Returns only active users
- **Aggregating:**
  - `[1, 2, 3] |reduce: $1 + $2` // Sums the array
- **Chaining:**
  - `products |filter: $1.price < 50 |map: $1.name` // Gets names of affordable products

## Tips for Using Pipes
- Use pipes to break complex logic into clear, sequential steps.
- Remember that `$1` refers to the value from the previous stage.
- Combine pipes with functions and operators for powerful transformations.

## Practice: Try It Yourself
Try these pipe expressions:
```
[10, 20, 30] |map: $1 / 10
users |filter: $1.isAdmin |map: $1.email
[1, 2, 3, 4] |filter: $1 % 2 == 0 |reduce: $1 + $2
```

Mastering pipes will help you write concise, readable, and powerful UExL code. In the next chapter, we'll explore the different types of pipes and how to chain them for advanced data processing.

The value of the first expression is accessible in the next stage as `$1` (and `$2`, etc. for reduce pipes).

See the following chapters for pipe types and chaining.