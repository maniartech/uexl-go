# Variables and Identifiers

Variables and identifiers in UExL are used to name values, reference data, and access properties. Understanding the rules for naming and using identifiers is essential for writing clear expressions.

## Naming Rules
- Identifiers can contain letters (a-z, A-Z), numbers (0-9), underscores (_), and the dollar sign ($).
- Identifiers cannot start with a number.
- Identifiers are case-sensitive (`Value` and `value` are different).
- Reserved words (such as `true`, `false`, `null`, and built-in function names) cannot be used as identifiers.

## Examples
```
x
count
user_name
$value
_data1
```

## Object Properties and Dot Notation
Dot notation is used to access object properties:
```
user.name
config.settings.theme
```
Bracket notation can also be used for dynamic or non-standard keys:
```
data["values"]
user["first-name"]
```

## Arrays and Indexing
Arrays are accessed by zero-based index:
```
arr[0]
users[2].name
```

## Special Identifiers in Pipes
Within pipe operations, `$1` refers to the input value passed to the pipe. For reduce pipes, `$2` and higher refer to additional arguments.
```
[1, 2, 3] |map: $1 * 2
[1, 2, 3, 4, 5] |reduce: $1 + $2
```

Following these rules will help you avoid naming conflicts and write maintainable UExL expressions.