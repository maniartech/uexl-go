package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---- RuneLength -------------------------------------------------------------

func TestRuneLength(t *testing.T) {
	assert.Equal(t, 5, RuneLength("hello"))      // ASCII: bytes == runes
	assert.Equal(t, 5, RuneLength("naïve"))      // ï is 2 bytes but 1 rune
	assert.Equal(t, 4, RuneLength("café"))       // precomposed é (U+00E9)
	assert.Equal(t, 5, RuneLength("café\u0301")) // decomposed é = e + combining accent
	assert.Equal(t, 7, RuneLength("👨‍👩‍👧‍👦"))    // family emoji: 4 people + 3 ZWJ = 7 runes
	assert.Equal(t, 0, RuneLength(""))
}

// ---- GraphemeLength ---------------------------------------------------------

func TestGraphemeLength(t *testing.T) {
	assert.Equal(t, 5, GraphemeLength("hello"))      // ASCII fast-path
	assert.Equal(t, 5, GraphemeLength("naïve"))      // ï = 1 grapheme
	assert.Equal(t, 4, GraphemeLength("café"))       // precomposed é = 1 grapheme
	assert.Equal(t, 4, GraphemeLength("café\u0301")) // decomposed é = 1 grapheme (e + combining)
	assert.Equal(t, 1, GraphemeLength("👨‍👩‍👧‍👦"))    // family emoji = 1 grapheme cluster
	assert.Equal(t, 0, GraphemeLength(""))
}

// ---- RuneSlice --------------------------------------------------------------

func TestRuneSlice(t *testing.T) {
	s, err := RuneSlice("naïve", 0, 3)
	assert.NoError(t, err)
	assert.Equal(t, "naï", s)

	s, err = RuneSlice("naïve", 2, 1)
	assert.NoError(t, err)
	assert.Equal(t, "ï", s)

	// Clamp past end.
	s, err = RuneSlice("hello", 3, 100)
	assert.NoError(t, err)
	assert.Equal(t, "lo", s)

	// Start beyond string.
	s, err = RuneSlice("hello", 10, 3)
	assert.NoError(t, err)
	assert.Equal(t, "", s)

	// Negative start: clamped to 0.
	s, err = RuneSlice("hello", -1, 3)
	assert.NoError(t, err)
	assert.Equal(t, "hel", s)

	// Negative length: error.
	_, err = RuneSlice("hello", 0, -1)
	assert.Error(t, err)

	// Empty string.
	s, err = RuneSlice("", 0, 5)
	assert.NoError(t, err)
	assert.Equal(t, "", s)
}

// ---- GraphemeSlice ----------------------------------------------------------

func TestGraphemeSlice(t *testing.T) {
	// ASCII fast-path.
	s, err := GraphemeSlice("hello", 1, 3)
	assert.NoError(t, err)
	assert.Equal(t, "ell", s)

	// Precomposed accented character.
	s, err = GraphemeSlice("café", 0, 3)
	assert.NoError(t, err)
	assert.Equal(t, "caf", s)

	// Decomposed: combining accent stays with its base letter.
	s, err = GraphemeSlice("café\u0301", 0, 4)
	assert.NoError(t, err)
	assert.Equal(t, "café\u0301", s)

	s, err = GraphemeSlice("café\u0301", 0, 3)
	assert.NoError(t, err)
	assert.Equal(t, "caf", s)

	// Emoji cluster stays whole.
	s, err = GraphemeSlice("👨‍👩‍👧‍👦 hi", 0, 2)
	assert.NoError(t, err)
	assert.Equal(t, "👨‍👩‍👧‍👦 ", s)

	// Start beyond string.
	s, err = GraphemeSlice("hi", 10, 3)
	assert.NoError(t, err)
	assert.Equal(t, "", s)

	// Length clamp.
	s, err = GraphemeSlice("hi", 1, 100)
	assert.NoError(t, err)
	assert.Equal(t, "i", s)

	// Negative length: error.
	_, err = GraphemeSlice("hi", 0, -1)
	assert.Error(t, err)

	// Empty string.
	s, err = GraphemeSlice("", 0, 5)
	assert.NoError(t, err)
	assert.Equal(t, "", s)
}

// ---- CollectRunes -----------------------------------------------------------

func TestCollectRunes(t *testing.T) {
	got := CollectRunes("naïve")
	assert.Equal(t, []any{"n", "a", "ï", "v", "e"}, got)

	got = CollectRunes("hi")
	assert.Equal(t, []any{"h", "i"}, got)

	got = CollectRunes("")
	assert.Equal(t, []any{}, got)
}

// ---- CollectGraphemes -------------------------------------------------------

func TestCollectGraphemes(t *testing.T) {
	// ASCII fast-path.
	got := CollectGraphemes("hi")
	assert.Equal(t, []any{"h", "i"}, got)

	// Decomposed accent clusters with base letter.
	got = CollectGraphemes("e\u0301")
	assert.Equal(t, []any{"e\u0301"}, got) // one grapheme cluster

	// Emoji counts as one.
	got = CollectGraphemes("👨‍👩‍👧‍👦")
	assert.Equal(t, []any{"👨‍👩‍👧‍👦"}, got)

	got = CollectGraphemes("")
	assert.Equal(t, []any{}, got)
}

// ---- CollectBytes -----------------------------------------------------------

func TestCollectBytes(t *testing.T) {
	got := CollectBytes("hi")
	assert.Equal(t, []any{float64('h'), float64('i')}, got)

	// ï = 0xC3 0xAF
	got = CollectBytes("ï")
	assert.Equal(t, []any{float64(0xC3), float64(0xAF)}, got)

	got = CollectBytes("")
	assert.Equal(t, []any{}, got)
}
