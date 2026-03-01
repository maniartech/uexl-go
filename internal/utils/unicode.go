package utils

import (
	"fmt"
	"unicode/utf8"

	"github.com/rivo/uniseg"
)

// ---- Measure ----------------------------------------------------------------

// RuneLength returns the number of Unicode code points in s.
// Uses an ASCII fast-path: for pure ASCII, runes == bytes.
func RuneLength(s string) int {
	if isASCII(s) {
		return len(s)
	}
	return utf8.RuneCountInString(s)
}

// GraphemeLength returns the number of user-perceived grapheme clusters in s.
// Uses an ASCII fast-path: for pure ASCII, graphemes == bytes.
func GraphemeLength(s string) int {
	if isASCII(s) {
		return len(s)
	}
	n := 0
	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		n++
	}
	return n
}

// ---- Cut --------------------------------------------------------------------

// RuneSlice returns the substring of s spanning runes [start : start+length].
// Negative start is clamped to 0. If start >= rune count, returns "".
// length is clamped to available runes.
func RuneSlice(s string, start, length int) (string, error) {
	if length < 0 {
		return "", fmt.Errorf("runeSubstr: length must be non-negative, got %d", length)
	}
	if start < 0 {
		start = 0
	}
	if isASCII(s) {
		// Byte == rune for ASCII — substring is a zero-copy slice.
		n := len(s)
		if start >= n {
			return "", nil
		}
		end := start + length
		if end > n {
			end = n
		}
		return s[start:end], nil
	}
	runes := []rune(s)
	n := len(runes)
	if start >= n {
		return "", nil
	}
	end := start + length
	if end > n {
		end = n
	}
	return string(runes[start:end]), nil
}

// GraphemeSlice returns the substring of s spanning grapheme clusters [start : start+length].
// Negative start is clamped to 0. If start >= grapheme count, returns "".
// length is clamped to available graphemes.
func GraphemeSlice(s string, start, length int) (string, error) {
	if length < 0 {
		return "", fmt.Errorf("graphemeSubstr: length must be non-negative, got %d", length)
	}
	if start < 0 {
		start = 0
	}
	if isASCII(s) {
		// Byte == rune == grapheme for ASCII.
		n := len(s)
		if start >= n {
			return "", nil
		}
		end := start + length
		if end > n {
			end = n
		}
		return s[start:end], nil
	}

	gr := uniseg.NewGraphemes(s)
	i := 0
	// Skip to start.
	for i < start && gr.Next() {
		i++
	}
	if i < start {
		return "", nil // start beyond end of string
	}
	// Collect length graphemes.
	from := -1
	to := 0
	for gr.Next() {
		s1, e1 := gr.Positions()
		if from == -1 {
			from = s1
		}
		to = e1
		length--
		if length == 0 {
			break
		}
	}
	if from == -1 {
		return "", nil
	}
	return s[from:to], nil
}

// ---- Explode ----------------------------------------------------------------

// CollectRunes returns each Unicode code point in s as a single-rune string,
// packed into a []any slice ready for the VM.
// ASCII fast-path uses zero-copy string slices. Unicode path uses range directly,
// avoiding an intermediate []rune allocation.
func CollectRunes(s string) []any {
	if isASCII(s) {
		result := make([]any, len(s))
		for i := range s {
			result[i] = s[i : i+1] // zero-copy: slice into original string
		}
		return result
	}
	// range s decodes runes directly — no []rune intermediate allocation.
	// Over-allocate by byte length (upper bound on rune count).
	result := make([]any, 0, len(s))
	for _, r := range s {
		result = append(result, string(r))
	}
	return result
}

// CollectGraphemes returns each grapheme cluster in s as a string,
// packed into a []any slice ready for the VM. Uses an ASCII fast-path.
func CollectGraphemes(s string) []any {
	if isASCII(s) {
		result := make([]any, len(s))
		for i := range s {
			result[i] = s[i : i+1]
		}
		return result
	}
	// Append directly to []any — avoids building an intermediate []string.
	result := make([]any, 0, len(s)/2) // heuristic: ≥2 UTF-8 bytes per non-ASCII grapheme
	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		result = append(result, gr.Str())
	}
	return result
}

// CollectBytes returns each UTF-8 byte in s as a float64,
// packed into a []any slice ready for the VM. Single-pass, no intermediate slice.
func CollectBytes(s string) []any {
	result := make([]any, len(s))
	for i := 0; i < len(s); i++ {
		result[i] = float64(s[i])
	}
	return result
}
