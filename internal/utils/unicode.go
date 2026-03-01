package utils

import (
	"fmt"
	"unicode/utf8"

	"github.com/rivo/uniseg"
)

// isASCII reports whether s contains only ASCII bytes (fast-path for all Unicode functions).
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

// ---- Measure ----------------------------------------------------------------

// RuneLength returns the number of Unicode code points in s.
func RuneLength(s string) int {
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
	runes := []rune(s)
	n := len(runes)
	if start < 0 {
		start = 0
	}
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
		if from == -1 {
			from, _ = gr.Positions()
		}
		_, to = gr.Positions()
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
func CollectRunes(s string) []any {
	runes := []rune(s)
	result := make([]any, len(runes))
	for i, r := range runes {
		result[i] = string(r)
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
	var clusters []string
	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		clusters = append(clusters, gr.Str())
	}
	result := make([]any, len(clusters))
	for i, c := range clusters {
		result[i] = c
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
