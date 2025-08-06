## Other Useful Pipe Types

In UExL, where user-defined functions are not available, certain pipe types provide powerful data manipulation capabilities that cannot be easily replaced by simple function calls. These pipes introduce new context variables or enable advanced operations directly in the expression language.

### sort
Sorts an array by a key or expression.

**Special variables:** `$item`, `$index`

**Example:**
```uexl
arr |sort: $item.age
```

---

### groupBy
Groups array elements by a key.

**Special variables:** `$item`, `$index`

**Example:**
```uexl
arr |groupBy: $item.type
```

---

### unique
Removes duplicates, optionally by a key.

**Special variables:** `$item`, `$index`

**Example:**
```uexl
arr |unique: $item.id
```

---

### find
Finds the first element matching a condition.

**Special variables:** `$item`, `$index`

**Example:**
```uexl
arr |find: $item > 10
```

---

### some / every
Checks if some or every element matches a condition.

**Special variables:** `$item`, `$index`

**Example:**
```uexl
arr |some: $item.active
arr |every: $item.valid
```

---

### flatMap
Maps and flattens arrays in one step.

**Special variables:** `$item`, `$index`

**Example:**
```uexl
arr |flatMap: $item.children
```

---

### window
Provides a sliding window of elements.

**Special variables:** `$window`, `$index`

**Example:**
```uexl
arr |window: 3 |map: $window[0] + $window[1] + $window[2]
```

---

### chunk
Splits array into chunks of a given size.

**Special variables:** `$chunk`, `$index`

**Example:**
```uexl
arr |chunk: 2
```

---

These pipe types enable expressive, reusable patterns for data transformation that would otherwise require user-defined functions. In UExL, they are essential for advanced data manipulation.
