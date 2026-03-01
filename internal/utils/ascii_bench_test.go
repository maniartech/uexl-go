package utils

import "testing"

// BenchmarkIsASCII measures throughput of the isASCII fast-path under both
// the scalar build (default) and the SIMD build (GOEXPERIMENT=simd).
//
// Run scalar:  go test ./internal/utils/... -bench=BenchmarkIsASCII -benchmem
// Run SIMD:    GOEXPERIMENT=simd go test ./internal/utils/... -bench=BenchmarkIsASCII -benchmem
func BenchmarkIsASCII(b *testing.B) {
	cases := []struct {
		name string
		s    string
	}{
		// Edge cases
		{"Empty", ""},
		{"1byte_ASCII", "a"},
		{"1byte_UTF8", "\xc3\xa9"}, // é — 2 bytes, non-ASCII immediately

		// Short strings (below one SIMD chunk)
		{"Short_5_ASCII", "hello"},
		{"Short_11_ASCII", "hello world"},
		{"Short_15_ASCII", "hello, world!!!"},

		// Boundary: exactly 16 bytes (one SIMD chunk)
		{"Exact_16_ASCII", "abcdefghijklmnop"},
		{"Exact_16_UTF8", "abcdefghijklmno\xc3"}, // last byte non-ASCII

		// Medium strings (1–4 SIMD chunks)
		{"Medium_45_ASCII", "The quick brown fox jumps over the lazy dog."},
		{"Medium_45_Early_NonASCII", "Thé quick brown fox jumps over the lazy dog."},
		{"Medium_45_Late_NonASCII", "The quick brown fox jumps over the lazy Bär."},

		// Long strings (many SIMD chunks)
		{"Long_88_ASCII", "The quick brown fox jumps over the lazy dog. Pack my box with five dozen liquor jugs."},

		// All non-ASCII (worst case for early-exit scalar; SIMD exits at first chunk)
		{"Short_AllNonASCII", "ÜÖÄüöä"},
		{"Long_AllNonASCII", "ÜÖÄüöäÜÖÄüöäÜÖÄüöäÜÖÄüöäÜÖÄüöäÜÖÄüöä"},

		// Realistic UExL strings
		{"FieldAccess_ASCII", "user.firstName"},
		{"PipeExpr_ASCII", "items |map: $item.price * 1.1"},
		{"UnicodeIdent", "données.prénom"},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			s := tc.s
			for b.Loop() {
				isASCII(s)
			}
		})
	}
}
