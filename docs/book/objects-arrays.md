# Objects, Arrays, and String Indexing

Objects and arrays are fundamental data structures in UExL, enabling you to represent and manipulate structured data. In addition, strings support bracket-based indexing to access individual characters by position.

## Objects
- Objects are collections of key-value pairs.
- Keys are strings (quoted or unquoted if valid identifiers).
- Values can be any UExL expression (number, string, array, object, etc.).
- Objects can be nested.

### Syntax
```
{
  key1: value1,
  key2: value2,
  "key 3": value3
}
```

### Accessing Object Properties
- Use dot notation: `obj.key1`
- Use bracket notation: `obj["key 3"]`
- Bracket notation is required for keys with spaces or special characters.
- The bracket key can be any expression that evaluates to a string — enabling computed/dynamic key access.

### Example
```
user = {
  name: "Alice",
  age: 30,
  "favorite color": "blue"
}
user.name                    // "Alice"
user["favorite color"]       // "blue"

// Computed key — any expression works as the key:
field = "name"
user[field]                  // "Alice"
user["favor" + "ite color"]  // "blue"
```

## Arrays
- Arrays are ordered collections of values.
- Elements can be any UExL expression.
- Arrays can be nested.

### Syntax
```
[1, 2, 3]
["a", {x: 1}, [2, 3]]
```

### Accessing Array Elements
- Use zero-based indexing: `arr[0]`
- Negative indices count from the end: `arr[-1]` is the last element, `arr[-2]` is second-to-last.
- Out-of-bounds access raises an error.

### Example
```
arr = [10, 20, 30]
arr[0]    // 10
arr[1]    // 20
arr[-1]   // 30  (last element)
arr[-2]   // 20
arr[10]   // error: array index out of bounds
```

## Strings: Index Access
- Strings can be indexed with square brackets to access a byte at a given zero-based position.
- Negative indices count from the end: `"hello"[-1]` is `"o"`.
- Out-of-bounds access raises an error.
- For Unicode strings, indexing is at the byte level. Use `runes(s)` or `graphemes(s)` to index safely at the code-point or grapheme level.

### Examples
```
"hello"[0]   // "h"
"hello"[4]   // "o"
"hello"[-1]  // "o"  (last byte)
"hello"[10]  // error: string index out of bounds
```

Note:
- Strings are immutable; indexing does not allow assignment.
- Indexing is by position. If your strings include multi-byte characters, indexing returns the unit at that position and may not correspond to a full user-perceived glyph.

## Advanced Usage
- Objects and arrays can be deeply nested:
  `{user: {profile: {name: "Bob"}}}`
- Arrays can contain objects, and vice versa.
- Use pipes to process arrays:
  `[1, 2, 3] |map: $item * 2`
- Use pipes to extract properties:
  `users |map: $1.name`
 - Use indexing on strings when mapping or filtering:
   `words |map: $1[0]` // take first character of each word

## Edge Cases
- Accessing a missing property returns `null`.
- Accessing an array with a non-integer index returns `null`.
- Modifying objects/arrays is not supported (expressions are immutable).
 - Accessing a string with a non-integer index returns `null`.

Objects and arrays are essential for modeling and transforming data in UExL.