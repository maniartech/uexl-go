# UExL Go Implementation: Feature Overview

## 1. Expression Language Core

### 1.1. Literals

- **NumberLiteral**: Supports floating-point and integer numbers.
- **StringLiteral**: Supports string values.
- **BooleanLiteral**: Supports `true` and `false`.
- **NullLiteral**: Represents null values.
- **ArrayLiteral**: Supports array construction and manipulation.
- **ObjectLiteral**: Supports object construction with key-value pairs.

### 1.2. Binary Expressions

- **Arithmetic**: `+`, `-`, `*`, `/`, `%` for numbers.
- **Comparison**: `==`, `!=`, `<`, `>`, `<=`, `>=` for numbers and strings.
- **Logical**: `&&`, `||` for booleans, with JavaScript-like short-circuiting and normalization (first truthy for `||`, last truthy for `&&`, normalization to boolean `false` for all-falsy `||` chains).
- **String Concatenation**: `+` for strings.

### 1.3. Unary Expressions

- **Negation**: `-` for numbers.
- **Logical NOT**: `!` for booleans.

### 1.4. Member Access & Indexing

- **Object Property Access**: `obj.property`
- **Array Indexing**: `arr[index]`
- **Nested Access**: `obj.arr[0].property`

### 1.5. Function Calls

- **Built-in Functions**: e.g., `len(arr)`, `substr(str, start, end)`, `contains(arr, value)`, `set(obj, key, value)`
- **Custom Functions**: Extensible via VM/builtins.

---

## 2. Pipe Expressions

### 2.1. General Pipe Syntax

- **Chaining**:

  expr | pipe1: block1 | pipe2: block2 ...

- **Alias Support**:

  |map\:item: block

  (binds `$item` to each element)

### 2.2. Pipe Handlers _(see `pipes.go`)_

- **map**: Transforms each array element using a predicate block.
- **filter**: Filters array elements by predicate block.
- **reduce**: Reduces array to a single value using accumulator logic.
- **find**: Returns first element matching predicate.
- **some**: Returns true if any element matches predicate.
- **every**: Returns true if all elements match predicate.
- **unique**: Removes duplicate elements.
- **sort**: Sorts array by predicate result (number/string).
- **groupBy**: Groups array elements by predicate result.
- **window**: Splits array into windows of size N, applies predicate.
- **chunk**: Splits array into chunks of size N, applies predicate.
- **pipe (default)**: Pass-through or custom block evaluation.

### 2.3. Context Variables in Pipes

- `$item`: Current element in array.
- `$index`: Current index.
- `$acc`: Accumulator (for reduce).
- `$window`: Current window (for window).
- `$chunk`: Current chunk (for chunk).
- `$last`: Last value (for default pipe).

---

## 3. Virtual Machine (VM) Architecture

### 3.1. Bytecode Execution

- **InstructionBlock**: Compiled instructions for each block.
- **Frame Stack**: Isolated execution context per block.
- **Scope Stack**: Variable resolution for pipes and blocks.

### 3.2. Variable Management

- `setPipeVar`: Sets context variables for pipes.
- `pushPipeScope` / `popPipeScope`: Manages variable scope per pipe stage.

### 3.3. Error Handling

- **Type Checking**: Ensures correct input types for pipes.
- **Predicate Validation**: Ensures block is present and valid.
- **Graceful Error Propagation**: Returns errors with context, no panics.

---

## 4. Extensibility

### 4.1. Pipe Handler Registry

- **DefaultPipeHandlers**: Maps pipe names to handler functions.
- **Custom Pipes**: Easily add new pipe types by registering handlers.

### 4.2. Built-in Functions

- **Extensible**: Add new built-ins in `builtins.go`.

---

## 5. Testing & Documentation

### 5.1. Unit Tests

- **Comprehensive Coverage**: Arithmetic, boolean, string, array, object, logical short-circuiting, and all pipe operations.
- **Edge Cases**: Empty arrays, invalid types, error propagation.

### 5.2. Documentation

- **Language Guide**: Syntax, semantics, and advanced features.
- **Pipe Overview**: Detailed explanation of each pipe.
- **Advanced Concepts**: VM internals, extensibility, error handling.

---

## 6. Production Readiness

- **Code Organization**: Modular, readable, and maintainable.
- **Naming Conventions**: Consistent and clear.
- **No Legacy Code**: Cleaned up for production use.

---

## Summary Table

| Feature        | Status   | Details                                      |
| -------------- | -------- | -------------------------------------------- |
| Literals       | Complete | Number, String, Boolean, Null, Array, Object |
| Binary Expr    | Complete | Arithmetic, Comparison, Logical (JS-like)    |
| Unary Expr     | Complete | Negation, Logical NOT                        |
| Member Access  | Complete | Object/Array/Nested                          |
| Function Calls | Complete | Built-in, extensible                         |
| Pipes          | Complete | All major types, context vars                |
| VM Execution   | Complete | Bytecode, frames, scopes                     |
| Error Handling | Complete | Type, predicate, propagation                 |
| Extensibility  | Complete | Pipe registry, built-ins                     |
| Testing        | Complete | Unit tests for all features                  |
| Documentation  | Complete | Language, pipes, advanced concepts           |

---

## References

- **`pipes.go`**: Pipe handlers and context management.
- **`bytecode.go`**: Bytecode and instruction blocks.
- **`vm.go`**: VM architecture and execution.
- **`builtins.go`**: Built-in functions.
- **`pipe_compilation_and_evaluation.md`**: Pipe compilation details.
- **`LANGUAGE.md`**: Language syntax and semantics.

---

**Conclusion**  
Your implementation now includes robust, JavaScript-like logical short-circuiting, normalization, and comprehensive tests. The system is production-ready, extensible, and well
