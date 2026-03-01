//go:build !goexperiment.simd

package utils

// isASCII reports whether s contains only ASCII bytes.
// Scalar path: byte-by-byte scan with early exit.
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}
