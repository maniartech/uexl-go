# Appendix C: Built-in Function Reference

UExL ships with exactly 14 built-in functions. All other utility functions (`upper`, `round`, `split`, etc.) must be provided by the host application via `WithFunctions`.

---

## `len(v)`

Returns the length of a string (in bytes) or an array (in elements).

**Signature**: `len(value) → number`

| Argument | Type | Description |
|----------|------|-------------|
| `v` | string or array | The value to measure |

**Returns**: Float64 count

**Errors**: Throws if `v` is neither a string nor an array.

```uexl
len('hello')          # 5
len('héllo')          # 6  (é is 2 bytes in UTF-8)
len([1, 2, 3])        # 3
```

> Use `runeLen` or `graphemeLen` for human-visible character counts.

---

## `substr(s, start, length)`

Returns a byte-level substring.

**Signature**: `substr(s, start, length) → string`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The source string |
| `start` | number (integer) | Zero-based start byte index |
| `length` | number (integer) | Number of bytes to extract |

**Errors**: Throws if `start` or `length` are not non-negative integers, or if the range is out-of-bounds.

```uexl
substr('hello', 0, 3)    # 'hel'
substr('hello', 2, 3)    # 'llo'
```

> For code-point-safe slicing use `runeSubstr`. For grapheme-safe slicing use `graphemeSubstr`.

---

## `contains(s, sub)`

Returns `true` if `s` contains `sub` as a byte-level substring.

**Signature**: `contains(s, sub) → bool`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The string to search in |
| `sub` | string | The substring to search for |

```uexl
contains('hello world', 'world')    # true
contains('hello world', 'xyz')      # false
contains('hello', '')               # true  (empty string is contained in any string)
```

---

## `set(obj, key, value)`

Sets a key on an object (mutates in-place) and returns the object.

**Signature**: `set(obj, key, value) → obj`

| Argument | Type | Description |
|----------|------|-------------|
| `obj` | object (`map[string]any`) | The object to modify |
| `key` | string or number | The key to set |
| `value` | any | The value to assign |

> **WARNING**: `set` mutates the original Go map. If you need immutability, clone the map in Go before passing it to the expression.

```uexl
set({}, 'name', 'Alice')                    # {name: 'Alice'}
set(product, 'discountedPrice', price * 0.9)  # modifies product AND returns it
```

---

## `str(v)`

Converts any value to its string representation using Go's `fmt.Sprintf("%v", v)` format.

**Signature**: `str(v) → string`

| Argument | Type | Description |
|----------|------|-------------|
| `v` | any | The value to convert |

```uexl
str(42)        # '42'
str(3.14)      # '3.14'
str(true)      # 'true'
str(false)     # 'false'
str(null)      # '<nil>'
str([1, 2])    # '[1 2]'
```

> `str(null)` returns `'<nil>'` (Go's default nil format), not `'null'`. If you need JSON-style null representation, handle nulls before conversion: `v == null ? 'null' : str(v)`.

---

## `runeLen(s)`

Returns the number of Unicode code points (runes) in a string.

**Signature**: `runeLen(s) → number`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The string to measure |

```uexl
runeLen('hello')    # 5
runeLen('héllo')    # 5  (é is 1 code point, even though it's 2 bytes)
runeLen('日本語')   # 3
```

---

## `runeSubstr(s, start, length)`

Returns a substring measured in Unicode code points.

**Signature**: `runeSubstr(s, start, length) → string`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The source string |
| `start` | number (integer) | Zero-based code point index |
| `length` | number (integer) | Number of code points to extract |

```uexl
runeSubstr('héllo', 0, 2)    # 'hé'
runeSubstr('日本語', 1, 2)   # '本語'
```

---

## `graphemeLen(s)`

Returns the number of user-perceived characters (extended grapheme clusters) in a string.

**Signature**: `graphemeLen(s) → number`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The string to measure |

```uexl
graphemeLen('café')    # 4  (é is one grapheme cluster)
graphemeLen('👨‍👩‍👧‍👦')  # 1  (family emoji is one grapheme cluster)
```

---

## `graphemeSubstr(s, start, length)`

Returns a substring measured in grapheme clusters (user-perceived characters).

**Signature**: `graphemeSubstr(s, start, length) → string`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The source string |
| `start` | number (integer) | Zero-based grapheme index |
| `length` | number (integer) | Number of graphemes to extract |

---

## `runes(s)`

Explodes a string into an array of single-rune strings.

**Signature**: `runes(s) → array`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The string to explode |

```uexl
runes('hi')      # ['h', 'i']
runes('日本')   # ['日', '本']
```

---

## `graphemes(s)`

Explodes a string into an array of grapheme cluster strings.

**Signature**: `graphemes(s) → array`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The string to explode |

```uexl
graphemes('café')   # ['c', 'a', 'f', 'é']
```

---

## `bytes(s)`

Explodes a string into an array of byte values (as float64).

**Signature**: `bytes(s) → array`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The string to explode |

```uexl
bytes('hi')    # [104, 105]  (ASCII codes for 'h' and 'i')
```

---

## `join(arr)` / `join(arr, sep)`

Joins an array of strings into a single string, with an optional separator.

**Signature**:
- `join(arr) → string` — joins with empty string separator
- `join(arr, sep) → string` — joins with `sep` between elements

| Argument | Type | Description |
|----------|------|-------------|
| `arr` | array of strings | The strings to join |
| `sep` | string (optional) | Separator string; defaults to `""` |

**Errors**: Throws if any element of `arr` is not a string.

```uexl
join(['a', 'b', 'c'])          # 'abc'
join(['a', 'b', 'c'], ', ')    # 'a, b, c'
join(['hello', 'world'], ' ')  # 'hello world'
```

---

## Function Availability Summary

| ✅ Built-in | ❌ Not built-in (host-provided) |
|------------|-------------------------------|
| `len`, `substr`, `contains` | `upper`, `lower`, `trim`, `replace` |
| `str`, `set` | `split`, `startsWith`, `endsWith` |
| `runeLen`, `runeSubstr` | `number`, `bool`, `string`, `typeof` |
| `graphemeLen`, `graphemeSubstr` | `min`, `max`, `floor`, `ceil`, `round`, `abs` |
| `runes`, `graphemes`, `bytes` | `concat`, `sum`, `isNaN`, `clamp` |
| `join` | (any other function) |
