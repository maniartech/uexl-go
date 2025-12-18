# Grapheme-Aware String Implementation Strategy

## Design Goal

Implement grapheme-aware string operations (#file:accessing-graphemes.md) in a way that is:
1. **Easy to implement**: Minimal changes to existing codebase
2. **User-extensible**: Works seamlessly with user-defined functions (they only see regular strings)
3. **Performance-conscious**: Fast-path for ASCII, correct for Unicode
4. **Zero-panic robustness**: Maintains UExL's error handling philosophy

## Core Design: Tagged String Pattern (Metadata Wrapper)

### Architecture Overview

**CRITICAL CONSTRAINT**: User-defined functions receive `string` type, not custom types. Therefore, we CANNOT use a StringView type that propagates through the system.

**Solution**: Use **metadata tagging** at VM operation boundaries:

```
String + Metadata Tag (view level) → VM Operations Apply Tag → Result String
                                                              ↓
                           User Functions Receive Plain Strings (no tag awareness needed)
```

### Key Insight: Metadata at Operation Boundaries, Not in Types

The design keeps strings as regular Go `string` type, but **VM operations consult metadata** to determine processing level:

1. **View functions** (`char()`, `utf8()`) return regular strings with attached metadata
2. **VM operations** (indexing, slicing, len) check metadata before processing
3. **User functions** receive and return plain strings (no awareness needed)
4. **Metadata propagates** through VM stack, not through function calls

## Implementation Plan

### Phase 1: String Metadata System (Foundation)

#### 1.1 Create `vm/string_metadata.go`

```go
package vm

// StringLevel represents the Unicode processing level for string operations
type StringLevel uint8

const (
	LevelGrapheme StringLevel = iota // Default: user-perceived characters
	LevelCodePoint                    // Runes/code points
	LevelUTF8                         // UTF-8 bytes
	LevelUTF16                        // UTF-16 code units
)

// stringMetadata is attached to strings on the VM stack to track their processing level
// This metadata is INTERNAL to the VM and never exposed to user functions
type stringMetadata struct {
	level StringLevel
}

// Default metadata for regular strings (grapheme-aware by default)
var defaultStringMetadata = stringMetadata{level: LevelGrapheme}

// stringWithMeta pairs a string with its processing metadata
// Used internally in VM stack operations
type stringWithMeta struct {
	str  string
	meta stringMetadata
}

// Helper functions for metadata management

func (vm *VM) getStringMeta(stackValue any) (string, stringMetadata) {
	switch v := stackValue.(type) {
	case stringWithMeta:
		return v.str, v.meta
	case string:
		// Regular string gets default metadata (grapheme-aware)
		return v, defaultStringMetadata
	default:
		return "", defaultStringMetadata
	}
}

func (vm *VM) pushStringWithMeta(str string, meta stringMetadata) error {
	// Only wrap with metadata if it's non-default
	// This keeps the common case (plain strings) efficient
	if meta.level == LevelGrapheme {
		return vm.Push(str) // Plain string, no wrapper
	}
	return vm.Push(stringWithMeta{str: str, meta: meta})
}
```

#### 1.2 Add String Processing Utilities (`internal/utils/grapheme.go`)

```go
package utils

// StringLevel and processing functions
type StringLevel uint8

const (
	LevelGrapheme StringLevel = iota
	LevelCodePoint
	LevelUTF8
	LevelUTF16
)

// StringLength returns the length of a string at the specified level
func StringLength(s string, level StringLevel) int {
	switch level {
	case LevelGrapheme:
		return GraphemeLength(s)
	case LevelCodePoint:
		return len([]rune(s))
	case LevelUTF8:
		return len(s)
	case LevelUTF16:
		return UTF16Length(s)
	default:
		return 0
	}
}

// StringIndex returns the element at index in the string at the specified level
func StringIndex(s string, idx int, level StringLevel) (string, error) {
	switch level {
	case LevelGrapheme:
		return GraphemeIndex(s, idx)
	case LevelCodePoint:
		return CodePointIndex(s, idx)
	case LevelUTF8:
		return UTF8Index(s, idx)
	case LevelUTF16:
		return UTF16Index(s, idx)
	default:
		return "", fmt.Errorf("unknown string level")
	}
}

// StringSlice returns a substring at the specified level
func StringSlice(s string, start, end, step int, level StringLevel) (string, error) {
	switch level {
	case LevelGrapheme:
		return GraphemeSlice(s, start, end, step)
	case LevelCodePoint:
		return CodePointSlice(s, start, end, step)
	case LevelUTF8:
		return UTF8Slice(s, start, end, step)
	case LevelUTF16:
		return UTF16Slice(s, start, end, step)
	default:
		return "", fmt.Errorf("unknown string level")
	}
}

// GraphemeLength returns the number of grapheme clusters (with ASCII fast-path)
func GraphemeLength(s string) int {
	if isASCII(s) {
		return len(s) // Fast path
	}
	return len(segmentGraphemes(s))
}

// GraphemeIndex returns the grapheme at the given index
func GraphemeIndex(s string, idx int) (string, error) {
	if isASCII(s) {
		// Fast path: ASCII strings
		if idx < 0 {
			idx = len(s) + idx
		}
		if idx < 0 || idx >= len(s) {
			return "", fmt.Errorf("index out of bounds")
		}
		return s[idx : idx+1], nil
	}

	// Complex path: Unicode with grapheme segmentation
	graphemes := segmentGraphemes(s)
	length := len(graphemes)
	if idx < 0 {
		idx = length + idx
	}
	if idx < 0 || idx >= length {
		return "", fmt.Errorf("index out of bounds")
	}
	return graphemes[idx], nil
}

// segmentGraphemes uses Unicode UAX #29 to segment string into grapheme clusters
func segmentGraphemes(s string) []string {
	// Fast-path: if all ASCII, just split bytes
	if isASCII(s) {
		result := make([]string, len(s))
		for i := 0; i < len(s); i++ {
			result[i] = s[i : i+1]
		}
		return result
	}

	// Complex path: Use Unicode segmentation
	// TODO: Implement using github.com/rivo/uniseg or similar
	return segmentGraphemesComplex(s)
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 128 {
			return false
		}
	}
	return true
}

// CodePointIndex, UTF8Index, UTF16Index implementations...
// (Similar pattern to GraphemeIndex)
```

### Phase 2: View Functions (Return Tagged Strings)

#### 3.1 Update String Indexing (`vm/indexing.go`)

```go
func (vm *VM) executeStringIndex(str string, index any) error {
	idxVal, ok := index.(float64)
	if !ok {
		return fmt.Errorf("string index must be a number")
	}

	intIdx := int(idxVal)
	if float64(intIdx) != idxVal {
		return fmt.Errorf("string index must be an integer, got %f", idxVal)
	}

	// NEW: Use grapheme-aware indexing by default
	result, err := utils.GraphemeIndex(str, intIdx)
	if err != nil {
		return err
	}

	// Push plain string result (no metadata)
	return vm.Push(result)
}
```

#### 3.2 Update String Slicing (`vm/slicing.go`)

```go
func (vm *VM) sliceString(str string, start, end, step any) error {
	// Parse step
	st, err := vm.parseSliceStep(step)
	if err != nil {
		return err
	}

	// Set defaults based on step direction
	length := utils.GraphemeLength(str) // NEW: Grapheme-aware length
	var defaultStart, defaultEnd int
	if st > 0 {
		defaultStart = 0
		defaultEnd = length
	} else {
		defaultStart = length - 1
		defaultEnd = -1
	}

	s, err := vm.parseSliceIndex(start, defaultStart)
	if err != nil {
		return err
	}

	e, err := vm.parseSliceIndex(end, defaultEnd)
	if err != nil {
		return err
	}

	s = vm.adjustSliceIndex(s, length)
	if e != -1 {
		e = vm.adjustSliceIndex(e, length)
	}

	// NEW: Use grapheme-aware slicing
	result, err := utils.GraphemeSlice(str, s, e, st)
	if err != nil {
		return err
	}

	// Push plain string result (no metadata)
	return vm.Push(result)
}
```

#### 3.3 Update Built-in Functions (`vm/builtins.go`)

```go
func builtinLen(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("len expects 1 argument")
	}

	switch v := args[0].(type) {
	case string:
		// Grapheme-aware length by default
		return float64(utils.GraphemeLength(v)), nil
	case stringWithMeta:
		// Tagged string: use its specified level
		return float64(utils.StringLength(v.str, v.meta.level)), nil
	case []any:
		return float64(len(v)), nil
	default:
		return nil, fmt.Errorf("len: unsupported type %T", args[0])
	}
}

func builtinContains(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("contains expects 2 arguments")
	}

	// Extract underlying strings (strip metadata)
	str1, _ := extractString(args[0])
	str2, _ := extractString(args[1])

	if str1 == "" && args[0] != "" {
		return nil, fmt.Errorf("contains: first argument must be a string")
	}
	if str2 == "" && args[1] != "" {
		return nil, fmt.Errorf("contains: second argument must be a string")
	}

	// Grapheme-aware contains (search for grapheme sequence)
	return utils.GraphemeContains(str1, str2), nil
}
```### Phase 4: StringView Support in VM Operations

#### 4.1 Extend Index Operation (`vm/indexing.go`)

```go
func (vm *VM) executeIndexExpression(left, index any, optional bool) error {
	if left == nil {
		if optional {
			return vm.Push(nil)
		}
		return fmt.Errorf("cannot index a null value")
	}

	switch typedLeft := left.(type) {
	case []any:
		return vm.executeArrayIndex(typedLeft, index)
	case map[string]any:
		return vm.executeObjectKey(typedLeft, index)
	case string:
		return vm.executeStringIndex(typedLeft, index)
	case *types.StringView: // NEW
		return vm.executeStringViewIndex(typedLeft, index)
	default:
		return fmt.Errorf("invalid type for index: %s", reflect.TypeOf(left).String())
	}
}

func (vm *VM) executeStringViewIndex(view *types.StringView, index any) error {
	idxVal, ok := index.(float64)
	if !ok {
		return fmt.Errorf("string view index must be a number")
	}

	result, err := view.Index(int(idxVal))
	if err != nil {
		return err
	}

	return vm.Push(result)
}
```

#### 4.2 Extend Slice Operation (`vm/slicing.go`)

```go
func (vm *VM) executeSliceExpression(target, start, end, step any, optional bool) error {
	if target == nil {
		if optional {
			return vm.Push(nil)
		}
		return fmt.Errorf("cannot slice a null value")
	}

	switch typedTarget := target.(type) {
	case []any:
		return vm.sliceArray(typedTarget, start, end, step)
	case string:
		return vm.sliceString(typedTarget, start, end, step)
	case *types.StringView: // NEW
		return vm.sliceStringView(typedTarget, start, end, step)
	default:
		return fmt.Errorf("invalid type for slice: %s", reflect.TypeOf(target).String())
	}
}

func (vm *VM) sliceStringView(view *types.StringView, start, end, step any) error {
	// Parse step
	st, err := vm.parseSliceStep(step)
	if err != nil {
		return err
	}

	// Parse indices
	// ... existing index parsing logic ...

	// Delegate to view
	result, err := view.Slice(s, e, st)
	if err != nil {
		return err
	}

	return vm.Push(result)
}
```

### Phase 5: Pipe Integration (Converts Tagged Strings to Arrays)

#### 5.1 String-to-Array Conversion for Pipes

When a string (tagged or plain) enters a pipe, it should be converted to an array of its elements:

```go
// In vm/pipes.go, before entering pipe handlers

func (vm *VM) preparePipeInput(input any) ([]any, error) {
	switch v := input.(type) {
	case []any:
		// Already an array
		return v, nil
	case string:
		// Convert to array of graphemes (default behavior)
		return stringToArray(v, LevelGrapheme), nil
	case stringWithMeta:
		// Convert to array at the specified level
		return stringToArray(v.str, v.meta.level), nil
	default:
		// Not iterable
		return nil, fmt.Errorf("pipe input must be array or string")
	}
}

func stringToArray(s string, level StringLevel) []any {
	length := utils.StringLength(s, utils.StringLevel(level))
	result := make([]any, length)
	for i := 0; i < length; i++ {
		elem, _ := utils.StringIndex(s, i, utils.StringLevel(level))
		result[i] = elem // Plain strings, no metadata
	}
	return result
}
```

#### 5.2 Update Pipe Handlers to Accept Pre-Converted Arrays

```go
func MapPipeHandler(input any, block any, alias string, vm *VM) (any, error) {
	// Input should already be prepared as array by the caller
	arr, ok := input.([]any)
	if !ok {
		// Fallback: convert if needed
		prepared, err := vm.preparePipeInput(input)
		if err != nil {
			return nil, err
		}
		arr = prepared
	}

	// ... existing map logic using arr ...
	// Result is always an array of values

	return result, nil
}
```

**Key Point**: After pipe processing, result is an array. User can join if needed:
```javascript
"éclair" |map: upper($item) |join: ""  // → "ÉCLAIR"
```### Phase 6: User-Defined Function Integration

**KEY INSIGHT**: User-defined functions **automatically work** because:

1. Functions receive `any` type arguments
2. StringView is just another `any` type
3. VM handles StringView transparently in operations

Example user-defined function:

```go
// User registers custom function
customFunctions["customSlice"] = func(args ...any) (any, error) {
	str := args[0].(string) // Works with both string and StringView
	start := int(args[1].(float64))
	end := int(args[2].(float64))

	// User's code just works - VM handles the view transparently
	// If str is StringView, it has String() method
	// If operations need indexing, VM handles it

	return str[start:end], nil // Go slice syntax
}
```

**Advanced**: If users want view-aware functions:

```go
customFunctions["advancedFunc"] = func(args ...any) (any, error) {
	// User can check for StringView explicitly
	if view, ok := args[0].(*types.StringView); ok {
		// Specialized logic for views
		return view.Index(0)
	}

	// Fallback for regular strings
	str := args[0].(string)
	return string(str[0]), nil
}
```

## Testing Strategy

### Test 1: Basic View Functions
```go
func TestStringViews_Basic(t *testing.T) {
	tests := []struct {
		expr     string
		expected any
	}{
		{`len("café")`, 4.0},              // 4 graphemes
		{`len(char("café"))`, 4.0},        // 4 code points (if "é" is U+00E9)
		{`len(char("caf\u0065\u0301"))`, 5.0}, // 5 code points (decomposed)
		{`"éclair"[0]`, "é"},              // Complete grapheme
		{`char("éclair")[0]`, "e"},        // First code point (if decomposed)
	}
	// ... run tests ...
}
```

### Test 2: Pipe Integration
```go
func TestStringViews_Pipes(t *testing.T) {
	tests := []struct {
		expr     string
		expected any
	}{
		{`"abc" |map: upper($item)`, []any{"A", "B", "C"}},
		{`char("abc") |map: upper($item)`, []any{"A", "B", "C"}},
		{`"abc" |filter: $item != "b"`, []any{"a", "c"}},
	}
	// ... run tests ...
}
```

### Test 3: User Functions
```go
func TestStringViews_UserFunctions(t *testing.T) {
	// Register user function
	userFuncs := vm.VMFunctions{
		"first": func(args ...any) (any, error) {
			str := args[0].(string)
			if len(str) == 0 {
				return "", nil
			}
			return string(str[0]), nil
		},
	}

	machine := vm.New(vm.LibContext{
		Functions: userFuncs,
	})

	// Test that user function works with views
	result, err := evaluateWithVM(machine, `first("hello")`)
	assert.NoError(t, err)
	assert.Equal(t, "h", result)
}
```

## Implementation Checklist

### Phase 1: Foundation (Types)
- [ ] Create `types/stringview.go` with StringView struct
- [ ] Implement Length(), Index(), Slice() methods
- [ ] Add segmentGraphemes() with ASCII fast-path
- [ ] Update `types/value.go` with TypeStringView
- [ ] Add NewStringViewValue(), AsStringView() methods

### Phase 2: View Functions
- [ ] Create `vm/string_views.go` with constructor functions
- [ ] Register char(), utf8(), utf16() in Builtins
- [ ] Add parser tests for view function calls

### Phase 3: Default Grapheme Behavior
- [ ] Update vm/indexing.go to use grapheme views by default
- [ ] Update vm/slicing.go to use grapheme views by default
- [ ] Update builtinLen() for grapheme-aware length
- [ ] Update builtinContains() for grapheme-aware search

### Phase 4: VM Operations
- [ ] Add executeStringViewIndex() to vm/indexing.go
- [ ] Add sliceStringView() to vm/slicing.go
- [ ] Update executeIndexExpression() with StringView case
- [ ] Update executeSliceExpression() with StringView case

### Phase 5: Pipe Support
- [ ] Add StringView case to MapPipeHandler
- [ ] Add StringView case to FilterPipeHandler
- [ ] Implement ToArray() method on StringView
- [ ] Add optional string joining for map results

### Phase 6: Testing
- [ ] Write unit tests for StringView type
- [ ] Write integration tests for view functions
- [ ] Write tests for default grapheme behavior
- [ ] Write tests for pipe operations with views
- [ ] Write tests for user function integration

### Phase 7: Performance Optimization
- [ ] Implement ASCII fast-path in segmentGraphemes()
- [ ] Add caching for repeated view operations
- [ ] Profile and optimize hot paths

### Phase 8: Documentation
- [ ] Update book/v2/accessing-graphemes.md with examples
- [ ] Add internal documentation for StringView API
- [ ] Document user function integration patterns

## Advantages of This Approach

### 1. Easy to Implement
- **Minimal changes**: Only add metadata handling in VM operations (indexing, slicing, len)
- **No type system changes**: Strings remain `string` type in function signatures
- **Incremental rollout**: Can implement phase-by-phase
- **Existing code mostly works**: Plain strings continue working as before

### 2. User-Extensible (SOLVED)
- **Zero awareness needed**: User functions receive plain `string` type
- **No special handling**: Users write normal Go string code
- **Metadata is internal**: Only VM operations see/use metadata
- **Automatic compatibility**: All user functions work without changes

### 3. Performance-Conscious
- **ASCII fast-path**: Detects ASCII-only strings, uses byte indexing
- **Lazy computation**: Only segments when first accessed
- **No wrapper overhead**: Plain strings are plain strings (no wrapping)
- **Metadata only when needed**: View functions create tagged strings, others use plain strings

### 4. Correct by Default
- **Grapheme-aware**: Default behavior is safe for Unicode
- **Explicit optimization**: Users can opt into code-point/byte level via view functions
- **No silent corruption**: Never breaks user-perceived characters
- **Transparent to users**: Complexity hidden in VM operations

## Migration Path

### Existing UExL Code
All existing code continues to work:
```javascript
len("hello")  // Still returns 5 (ASCII fast-path)
"abc"[1]      // Still returns "b" (ASCII fast-path)
```

### Unicode-Aware Code
New code gets correct behavior automatically:
```javascript
len("café")    // Returns 4 (grapheme count)
"éclair"[0]    // Returns complete "é" grapheme
```

### Performance-Critical Code
Can opt into lower levels explicitly:
```javascript
len(char("café"))  // Code point count for identifiers
len(utf8("café"))  // Byte count for protocols
```

## Performance Characteristics

### ASCII Fast-Path
- **Detection**: O(n) scan on first access, cached
- **Indexing**: O(1) byte access
- **Slicing**: O(1) substring operation
- **Memory**: Zero overhead (no segmentation)

### Unicode Grapheme Path
- **Segmentation**: O(n) on first access, cached
- **Indexing**: O(1) array access after segmentation
- **Slicing**: O(k) where k = slice length
- **Memory**: ~2x string size (stores grapheme offsets)

### View Construction
- **char()**: O(n) rune conversion (lazy)
- **utf8()**: O(1) wrapper (no conversion)
- **utf16()**: O(n) conversion (lazy)

## Comparison with Alternatives

### Alternative 1: StringView Type (Original Proposal)
**Problem**: User functions can't receive custom types, only `string`
**Solution**: Use metadata tagging instead of type wrapping

### Alternative 2: Modify String Type System
**Problem**: Would break all existing user function signatures
**Solution**: Keep strings as strings, add internal metadata

### Alternative 3: Rune-Based Default
**Problem**: Breaks grapheme clusters (emoji, combining marks)
**Our approach**: Grapheme-aware by default, opt-in for code points

### Alternative 4: Always Grapheme-Segment Everything
**Problem**: Performance overhead for ASCII strings
**Our approach**: ASCII fast-path detection + lazy segmentation

## Conclusion

This design provides the best balance of:
- **Ease of implementation**: Minimal code changes (only VM operations)
- **User extensibility**: User functions receive plain `string` (no awareness needed)
- **Performance**: ASCII fast-path + lazy Unicode processing
- **Correctness**: Grapheme-aware by default

The key insight is that **metadata at boundaries, not in types**: by keeping strings as `string` type and using internal metadata tags only at VM operation boundaries, we achieve grapheme-awareness without breaking user function compatibility.
