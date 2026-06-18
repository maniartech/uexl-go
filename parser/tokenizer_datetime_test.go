package parser

import (
	"strings"
	"testing"

	"github.com/maniartech/uexl/parser/constants"
)

func allTokens(input string) ([]Token, error) {
	tz := NewTokenizer(input)
	var toks []Token
	for {
		tok, err := tz.NextToken()
		if err != nil {
			return toks, err
		}
		if tok.Type == constants.TokenEOF {
			return toks, nil
		}
		toks = append(toks, tok)
	}
}

func TestTokenize_DateTimeLiteral(t *testing.T) {
	cases := map[string]int64{
		`d"2024-12-01"`:                1733011200000,
		`d"1970-01-01T00:00:00Z"`:      0,
		`d'2024-12-01T10:30:00Z'`:      1733049000000,
		`d"1969-12-31T00:00:00Z"`:      -86400000,
		`d"0001-01-01T00:00:00Z"`:      -62135596800000,
		`d"2024-12-01T10:30:00+05:30"`: 1733029200000,
	}
	for in, want := range cases {
		toks, err := allTokens(in)
		if err != nil {
			t.Errorf("%s: unexpected error %v", in, err)
			continue
		}
		if len(toks) != 1 || toks[0].Type != constants.TokenDateTime {
			t.Errorf("%s: expected one DateTime token, got %v", in, toks)
			continue
		}
		if got, ok := toks[0].AsInt(); !ok || got != want {
			t.Errorf("%s: got %d (ok=%v), want %d", in, got, ok, want)
		}
	}
}

func TestTokenize_DurationLiteral(t *testing.T) {
	cases := map[string]int64{
		"7d":    604800000,
		"30ms":  30,
		"1.5h":  5400000,
		"500ms": 500,
		"2w":    1209600000,
		"45s":   45000,
		"10m":   600000,
		"0d":    0,
		"3h":    10800000,
	}
	for in, want := range cases {
		toks, err := allTokens(in)
		if err != nil {
			t.Errorf("%s: unexpected error %v", in, err)
			continue
		}
		if len(toks) != 1 || toks[0].Type != constants.TokenDuration {
			t.Errorf("%s: expected one Duration token, got %v", in, toks)
			continue
		}
		if got, ok := toks[0].AsInt(); !ok || got != want {
			t.Errorf("%s: got %d (ok=%v), want %d", in, got, ok, want)
		}
	}
}

func TestTokenize_DateTimeLiteralErrors(t *testing.T) {
	bad := []string{
		`d"2024-13-01"`,           // month out of range (dt-err-001)
		`d"2024-02-30"`,           // day out of range (dt-err-002)
		`d"2024-12-01T10:30:60Z"`, // leap second (dt-err-003)
		`d"0000-01-01"`,           // year 0000 (dt-err-004)
		`d"10000-01-01"`,          // 5-digit year (dt-err-005)
		`d"not-a-date"`,           // malformed
		`d"2024-12-01`,            // unterminated
	}
	for _, in := range bad {
		if _, err := allTokens(in); err == nil {
			t.Errorf("%s: expected a parse-time error", in)
		}
	}
}

// TestTokenize_DurationQuoteAdjacency covers dt-err-009: a duration immediately followed by a quote
// (no operator) is a parse error.
func TestTokenize_DurationQuoteAdjacency(t *testing.T) {
	for _, in := range []string{`7d"2024-12-01"`, `30ms'x'`} {
		if _, err := allTokens(in); err == nil {
			t.Errorf("%s: expected a parse-time error (dt-err-009)", in)
		}
	}
}

// TestTokenize_DurationOutOfRange covers §3.3.6: a duration whose value exceeds the representable range
// is a parse error (exercises the bound check in makeDuration).
func TestTokenize_DurationOutOfRange(t *testing.T) {
	for _, in := range []string{"1000000000000d", "99999999w"} {
		if _, err := allTokens(in); err == nil {
			t.Errorf("%q: expected an out-of-range duration error", in)
		}
	}
}

// TestToken_AsIntNonTemporal verifies AsInt reports false for non-temporal tokens.
func TestToken_AsIntNonTemporal(t *testing.T) {
	toks, err := allTokens("42")
	if err != nil || len(toks) != 1 {
		t.Fatalf("42: got %v err %v", toks, err)
	}
	if v, ok := toks[0].AsInt(); ok {
		t.Errorf("AsInt on a number token should be (0,false), got (%d,true)", v)
	}
}

// TestTokenize_DurationHugeMagnitude covers the overflow guard: a magnitude so large that ParseFloat
// returns +Inf must be rejected, not silently produce a garbage int64.
func TestTokenize_DurationHugeMagnitude(t *testing.T) {
	if _, err := allTokens(strings.Repeat("9", 400) + "d"); err == nil {
		t.Error("an overflowing duration magnitude should be a parse error")
	}
}

// TestTokenize_DateTimeMultiline covers the newline branch of the datetime scanner: a newline inside the
// literal is consumed (line/column tracking) and the content then fails ISO parsing.
func TestTokenize_DateTimeMultiline(t *testing.T) {
	if _, err := allTokens("d\"2024-\n12-01\""); err == nil {
		t.Error("a datetime literal with an embedded newline should error")
	}
}

// TestTokenize_IdentifiersWithDPrefix verifies the dispatch fallthrough chain: identifiers that begin
// with d or r (and the bare letters) still lex as identifiers, not datetime literals.
func TestTokenize_IdentifiersWithDPrefix(t *testing.T) {
	for _, in := range []string{"day", "dur", "d", "raw", "r", "date", "duration"} {
		toks, err := allTokens(in)
		if err != nil {
			t.Errorf("%s: unexpected error %v", in, err)
			continue
		}
		if len(toks) != 1 || toks[0].Type != constants.TokenIdentifier {
			t.Errorf("%s: expected one Identifier token, got %v", in, toks)
		}
	}
}

// TestTokenize_DurationWordBoundary verifies longest-match and that a unit followed by an identifier
// char is NOT a duration (30ms is one duration; 5month is number 5 + identifier "month").
func TestTokenize_DurationWordBoundary(t *testing.T) {
	// 30ms must be a single Duration token of 30 (ms), not 30m + s.
	toks, err := allTokens("30ms")
	if err != nil || len(toks) != 1 || toks[0].Type != constants.TokenDuration {
		t.Fatalf("30ms: got %v err %v", toks, err)
	}
	// 5month: number 5 then identifier "month" (the 'm' is not a duration unit here).
	toks2, err := allTokens("5month")
	if err != nil {
		t.Fatalf("5month: unexpected error %v", err)
	}
	if len(toks2) != 2 || toks2[0].Type != constants.TokenNumber || toks2[1].Type != constants.TokenIdentifier {
		t.Errorf("5month: expected [Number, Identifier], got %v", toks2)
	}
	// 1e3ms: scientific notation is not a duration magnitude -> number 1e3 then identifier ms.
	toks3, err := allTokens("1e3ms")
	if err != nil {
		t.Fatalf("1e3ms: unexpected error %v", err)
	}
	if len(toks3) != 2 || toks3[0].Type != constants.TokenNumber || toks3[1].Type != constants.TokenIdentifier {
		t.Errorf("1e3ms: expected [Number, Identifier], got %v", toks3)
	}
}
